package booth

import (
	"context"
	"fmt"

	"github.com/trkd-knt/booth-go/internal/boothhttp"
	"github.com/trkd-knt/booth-go/internal/parser"
)

// GetShop は指定したショップホストからショップ情報を取得します。
func (c *Client) GetShop(ctx context.Context, shopHost string) (*Shop, error) {
	rawURL, err := boothhttp.BuildShopURL(shopHost)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequestFailed, err)
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
