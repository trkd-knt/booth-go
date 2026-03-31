package booth

import (
	"context"
	"errors"
	"fmt"

	"github.com/trkd-knt/booth-go/internal/boothhttp"
	"github.com/trkd-knt/booth-go/internal/parser"
)

// GetItem は商品 ID から商品詳細を取得します。
func (c *Client) GetItem(ctx context.Context, itemID string) (*Item, error) {
	rawURL, err := boothhttp.BuildItemURL(c.lang, itemID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequestFailed, err)
	}
	jsonURL, err := boothhttp.BuildItemJSONURL(c.lang, itemID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequestFailed, err)
	}

	body, err := c.doGet(ctx, jsonURL, "item")
	if err != nil {
		return nil, err
	}
	defer body.Close()

	item, err := parser.DecodeItemJSON(body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseFailed, err)
	}
	if err := parser.ValidateDecodedItem(item); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseFailed, err)
	}
	parser.NormalizeItem(item, rawURL)

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
	rawURL, err := boothhttp.BuildItemURL(c.lang, itemID)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrRequestFailed, err)
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
