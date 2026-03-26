package booth

import (
	"context"
	"fmt"

	"github.com/trkd-knt/booth-go/internal/parser"
)

// GetItem は指定したショップホストと商品 ID から商品詳細を取得します。
func (c *Client) GetItem(ctx context.Context, shopHost, itemID string) (*Item, error) {
	rawURL, err := buildItemURL(shopHost, itemID)
	if err != nil {
		return nil, err
	}

	body, err := c.doGet(ctx, rawURL, "item")
	if err != nil {
		return nil, err
	}
	defer body.Close()

	item, err := parser.ParseItemPage(body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseFailed, err)
	}
	if item.URL == "" {
		item.URL = rawURL
	}
	if item.ShopHost == "" {
		item.ShopHost = shopHost
	}

	return item, nil
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
func buildItemURL(shopHost, itemID string) (string, error) {
	shopURL, err := buildShopURL(shopHost)
	if err != nil {
		return "", err
	}
	return shopURL + "/items/" + itemID, nil
}

// buildShopURL はショップページの URL を構築します。
func buildShopURL(shopHost string) (string, error) {
	if shopHost == "" {
		return "", fmt.Errorf("%w: shop host must not be empty", ErrRequestFailed)
	}
	return "https://" + shopHost, nil
}
