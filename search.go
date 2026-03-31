package booth

import (
	"context"
	"fmt"

	"github.com/trkd-knt/booth-go/internal/boothhttp"
	"github.com/trkd-knt/booth-go/internal/parser"
)

// SearchItems は BOOTH 全体検索結果を取得します。
func (c *Client) SearchItems(ctx context.Context, opts SearchOptions) (*SearchResult, error) {
	if err := validateSearchOptions(opts); err != nil {
		return nil, err
	}
	opts = normalizeSearchOptions(opts)

	rawURL, err := boothhttp.NewSearchURL(c.searchBaseURL, c.lang, opts)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequestFailed, err)
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
