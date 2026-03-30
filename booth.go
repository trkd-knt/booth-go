package booth

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/trkd-knt/booth-go/internal/parser"
)

// GetItem は商品 ID から商品詳細を取得します。
func (c *Client) GetItem(ctx context.Context, itemID string) (*Item, error) {
	rawURL, err := buildItemURL(c.lang, itemID)
	if err != nil {
		return nil, err
	}
	jsonURL, err := buildItemJSONURL(c.lang, itemID)
	if err != nil {
		return nil, err
	}

	body, err := c.doGet(ctx, jsonURL, "item")
	if err != nil {
		return nil, err
	}
	defer body.Close()

	item, err := decodeItemJSON(body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseFailed, err)
	}
	if err := validateDecodedItem(item); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseFailed, err)
	}
	normalizeItem(item, rawURL)

	htmlBody, err := c.doGet(ctx, rawURL, "item")
	if err == nil {
		defer htmlBody.Close()
		avatars, parseErr := parser.ParseItemAvatarsPage(htmlBody)
		if parseErr != nil && !errors.Is(parseErr, parser.ErrItemAvatarsNotFound()) {
			return nil, fmt.Errorf("%w: %v", ErrItemAvatarsParseFailed, parseErr)
		}
		item.Avatars = avatars
	}

	return item, nil
}

// GetItemDescription は商品 ID から HTML 本文説明を取得します。
func (c *Client) GetItemDescription(ctx context.Context, itemID string) (string, error) {
	rawURL, err := buildItemURL(c.lang, itemID)
	if err != nil {
		return "", err
	}

	body, err := c.doGet(ctx, rawURL, "item")
	if err != nil {
		return "", err
	}
	defer body.Close()

	description, err := parser.ParseItemDescriptionPage(body)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrItemDescriptionParseFailed, err)
	}

	return description, nil
}

// SearchItems は BOOTH 全体検索結果を取得します。
func (c *Client) SearchItems(ctx context.Context, opts SearchOptions) (*SearchResult, error) {
	opts = normalizeSearchOptions(opts)

	rawURL, err := c.newSearchURL(opts)
	if err != nil {
		return nil, err
	}

	body, err := c.doGet(ctx, rawURL, "search")
	if err != nil {
		return nil, err
	}
	defer body.Close()

	result, err := parser.ParseSearchPage(body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseFailed, err)
	}
	if result.Page == 0 {
		result.Page = opts.Page
	}

	return result, nil
}

// GetShop は指定したショップホストからショップ情報を取得します。
func (c *Client) GetShop(ctx context.Context, shopHost string) (*Shop, error) {
	rawURL, err := buildShopURL(shopHost)
	if err != nil {
		return nil, err
	}

	body, err := c.doGet(ctx, rawURL, "shop")
	if err != nil {
		return nil, err
	}
	defer body.Close()

	shop, err := parser.ParseShopPage(body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseFailed, err)
	}
	if shop.URL == "" {
		shop.URL = rawURL
	}
	if shop.Host == "" {
		shop.Host = shopHost
	}

	return shop, nil
}

// buildItemURL は商品詳細ページの URL を構築します。
func buildItemURL(lang, itemID string) (string, error) {
	itemID = strings.TrimSpace(itemID)
	if itemID == "" {
		return "", fmt.Errorf("%w: item id must not be empty", ErrRequestFailed)
	}
	lang = strings.Trim(strings.TrimSpace(lang), "/")
	if lang == "" {
		lang = defaultLang
	}
	return "https://booth.pm/" + lang + "/items/" + itemID, nil
}

// buildShopURL はショップページの URL を構築します。
func buildShopURL(shopHost string) (string, error) {
	if shopHost == "" {
		return "", fmt.Errorf("%w: shop host must not be empty", ErrRequestFailed)
	}
	return "https://" + shopHost, nil
}

func parseURLHost(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return parsed.Host
}
