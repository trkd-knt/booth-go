package booth

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type mockHTTPClient struct {
	do func(req *http.Request) (*http.Response, error)
}

func (m mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.do(req)
}

// TestGetItem は商品詳細取得の基本動作を確認します。
func TestGetItem(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		WithHTTPClient(mockHTTPClient{
			do: func(req *http.Request) (*http.Response, error) {
				switch req.URL.String() {
				case "https://booth.pm/ja/items/12345.json":
					return &http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(strings.NewReader(`{
							"id": 12345,
							"name": "サンプル商品",
							"description": "説明",
							"price": "¥ 1,200",
							"url": "https://booth.pm/ja/items/12345",
							"is_sold_out": false,
							"is_adult": false,
							"wish_lists_count": 42,
							"category": {"id": 209, "name": "3D衣装"},
							"shop": {"name": "サンプルショップ", "subdomain": "sample", "thumbnail_url": "https://example.com/thumb.jpg", "url": "https://sample.booth.pm/"},
							"tags": [{"name": "VRChat"}, {"name": "衣装"}],
							"images": [
								{"original": "https://example.com/original.jpg", "resized": "https://example.com/resized.jpg"}
							]
						}`)),
					}, nil
				case "https://booth.pm/ja/items/12345":
					return &http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(strings.NewReader(`<!doctype html><html><body>
							<ul id="variations">
								<li><div class="variation-name">フルパック - FullPack</div></li>
								<li><div class="variation-name">ルミナ - LUMINA</div></li>
								<li><div class="variation-name">ショコラ - Chocolat ※共通素体あり</div></li>
							</ul>
							</body></html>`)),
					}, nil
				default:
					t.Fatalf("unexpected url: %s", req.URL.String())
				}
				return nil, nil
			},
		}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	item, err := client.GetItem(context.Background(), "12345")
	if err != nil {
		t.Fatalf("GetItem() error = %v", err)
	}
	if item.Title != "サンプル商品" {
		t.Fatalf("item.Title = %q", item.Title)
	}
	if item.Summary != "説明" {
		t.Fatalf("item.Summary = %q", item.Summary)
	}
	if item.PriceText != "¥ 1,200" {
		t.Fatalf("item.PriceText = %q", item.PriceText)
	}
	if item.Price != 1200 {
		t.Fatalf("item.Price = %d", item.Price)
	}
	if len(item.Tags) != 2 || item.Tags[0] != "VRChat" {
		t.Fatalf("item.Tags = %#v", item.Tags)
	}
	if len(item.Images) != 1 || item.Images[0] != "https://example.com/resized.jpg" {
		t.Fatalf("item.Images = %#v", item.Images)
	}
	if len(item.Avatars) != 3 || item.Avatars[0] != "フルパック - FullPack" || item.Avatars[1] != "ルミナ - LUMINA" || item.Avatars[2] != "ショコラ - Chocolat ※共通素体あり" {
		t.Fatalf("item.Avatars = %#v", item.Avatars)
	}
}

// TestGetItemUsesPayloadShop は item JSON の shop 情報をそのまま利用することを確認します。
func TestGetItemUsesPayloadShop(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		WithHTTPClient(mockHTTPClient{
			do: func(req *http.Request) (*http.Response, error) {
				switch req.URL.String() {
				case "https://booth.pm/ja/items/12345.json":
					return &http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(strings.NewReader(`{
							"id": 12345,
							"name": "サンプル商品",
							"url": "https://booth.pm/ja/items/12345",
							"price": "¥ 1,200",
							"shop": {"name": "サンプルショップ", "url": "https://booth.pm/"}
						}`)),
					}, nil
				case "https://booth.pm/ja/items/12345":
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`<!doctype html><html><body></body></html>`)),
					}, nil
				default:
					t.Fatalf("unexpected url: %s", req.URL.String())
				}
				return nil, nil
			},
		}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	item, err := client.GetItem(context.Background(), "12345")
	if err != nil {
		t.Fatalf("GetItem() error = %v", err)
	}
	if item.URL != "https://booth.pm/ja/items/12345" {
		t.Fatalf("item.URL = %q", item.URL)
	}
	if item.ShopHost != "booth.pm" {
		t.Fatalf("item.ShopHost = %q", item.ShopHost)
	}
	if item.Shop == nil || item.Shop.Host != "booth.pm" || item.Shop.URL != "https://booth.pm/" {
		t.Fatalf("item.Shop = %+v", item.Shop)
	}
}

// TestGetItemNotFound は 404 が商品未検出へ変換されることを確認します。
func TestGetItemNotFound(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		WithHTTPClient(mockHTTPClient{
			do: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader("not found")),
				}, nil
			},
		}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetItem(context.Background(), "12345")
	if !errors.Is(err, ErrItemNotFound) {
		t.Fatalf("GetItem() error = %v, want ErrItemNotFound", err)
	}
}

func TestGetItemDescription(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		WithHTTPClient(mockHTTPClient{
			do: func(req *http.Request) (*http.Response, error) {
				if req.URL.String() != "https://booth.pm/ja/items/12345" {
					t.Fatalf("unexpected url: %s", req.URL.String())
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`<!doctype html><html><body><article>
						<section class="main-info-column">
							<div class="js-market-item-detail-description description">
								<p class="autolink whitespace-pre-line">冒頭本文</p>
							</div>
						</section>
						<section class="shop__text"><h2>見出し</h2><p class="js-autolink whitespace-pre-line">後続本文</p><p class="js-autolink whitespace-pre-line">補足段落</p></section>
						<div class="shop__text"><p>無関係</p></div>
						</article></body></html>`)),
				}, nil
			},
		}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	description, err := client.GetItemDescription(context.Background(), "12345")
	if err != nil {
		t.Fatalf("GetItemDescription() error = %v", err)
	}
	if !strings.Contains(description, "冒頭本文") || !strings.Contains(description, "見出し\n後続本文") {
		t.Fatalf("description = %q", description)
	}
}

func TestGetItemDescriptionParseFailed(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		WithHTTPClient(mockHTTPClient{
			do: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`<!doctype html><html><body></body></html>`)),
				}, nil
			},
		}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetItemDescription(context.Background(), "12345")
	if !errors.Is(err, ErrItemDescriptionParseFailed) {
		t.Fatalf("GetItemDescription() error = %v, want ErrItemDescriptionParseFailed", err)
	}
}

// TestSearchItems は検索結果取得の基本動作を確認します。
func TestSearchItems(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		WithHTTPClient(mockHTTPClient{
			do: func(req *http.Request) (*http.Response, error) {
				if req.URL.String() != "https://booth.pm/ja/search?q=sample" {
					t.Fatalf("unexpected url: %s", req.URL.String())
				}
				if got := req.Header.Get("User-Agent"); got == "" {
					t.Fatalf("User-Agent is empty")
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`<!doctype html><html><body>
						<p>対象商品 1 件</p>
						<article><a href="https://sample.booth.pm/items/12345">サンプル商品1</a><a href="https://sample.booth.pm">サンプルショップ</a><span>¥ 500</span></article>
						</body></html>`)),
				}, nil
			},
		}),
		WithSearchBaseURL("https://booth.pm/ja/search"),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	result, err := client.SearchItems(context.Background(), SearchOptions{Query: "sample", Page: 1})
	if err != nil {
		t.Fatalf("SearchItems() error = %v", err)
	}
	if len(result.Items) != 1 {
		t.Fatalf("len(result.Items) = %d, want 1", len(result.Items))
	}
	if result.Items[0].Shop == nil || result.Items[0].Shop.Host != "sample.booth.pm" {
		t.Fatalf("result.Items[0].Shop = %+v", result.Items[0].Shop)
	}
	if result.Items[0].ShopHost != "" {
		t.Fatalf("result.Items[0].ShopHost = %q, want empty", result.Items[0].ShopHost)
	}
}

// TestWithLang は検索言語設定が URL とヘッダーに反映されることを確認します。
func TestWithLang(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		WithLang("en"),
		WithHTTPClient(mockHTTPClient{
			do: func(req *http.Request) (*http.Response, error) {
				if req.URL.String() != "https://booth.pm/en/search?q=sample" {
					t.Fatalf("unexpected url: %s", req.URL.String())
				}
				if got := req.Header.Get("Accept-Language"); got != "en" {
					t.Fatalf("Accept-Language = %q, want %q", got, "en")
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`<!doctype html><html><body>
						<p>対象商品 1 件</p>
						<article><a href="https://sample.booth.pm/items/12345">サンプル商品1</a><a href="https://sample.booth.pm">サンプルショップ</a><span>¥ 500</span></article>
						</body></html>`)),
				}, nil
			},
		}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.SearchItems(context.Background(), SearchOptions{Query: "sample", Page: 1})
	if err != nil {
		t.Fatalf("SearchItems() error = %v", err)
	}
}

// TestSearchItemsBrowseCategoryAndFilters は browse URL と各種検索条件が反映されることを確認します。
func TestSearchItemsBrowseCategoryAndFilters(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		WithHTTPClient(mockHTTPClient{
			do: func(req *http.Request) (*http.Response, error) {
				got := req.URL.String()
				want := "https://booth.pm/ja/browse/3D%E8%A1%A3%E8%A3%85?adult=only&event=osakafes-mar2026&except_words%5B%5D=b&max_price=2500&min_price=1000&q=a&sort=new&tags%5B%5D=t&type=digital"
				if got != want {
					t.Fatalf("unexpected url: %s", got)
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`<!doctype html><html><body>
						<p>対象商品 1 件</p>
						<article><a href="https://sample.booth.pm/items/12345">サンプル商品1</a><a href="https://sample.booth.pm">サンプルショップ</a><span>¥ 500</span></article>
						</body></html>`)),
				}, nil
			},
		}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.SearchItems(context.Background(), SearchOptions{
		Category:    "3D衣装",
		Query:       "a",
		ExceptWords: []string{"b"},
		Tags:        []string{"t"},
		Event:       "osakafes-mar2026",
		Type:        ItemTypeDigital,
		Adult:       AdultFilterOnly,
		MinPrice:    1000,
		MaxPrice:    2500,
		Sort:        SortNewest,
	})
	if err != nil {
		t.Fatalf("SearchItems() error = %v", err)
	}
}

// TestSearchItemsInvalidPriceRange は不正価格範囲を弾くことを確認します。
func TestSearchItemsInvalidPriceRange(t *testing.T) {
	t.Parallel()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.SearchItems(context.Background(), SearchOptions{
		MinPrice: 3000,
		MaxPrice: 1000,
	})
	if !errors.Is(err, ErrRequestFailed) {
		t.Fatalf("SearchItems() error = %v, want ErrRequestFailed", err)
	}
}

// TestSearchItemsInvalidItemType は不正な type を弾くことを確認します。
func TestSearchItemsInvalidItemType(t *testing.T) {
	t.Parallel()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.SearchItems(context.Background(), SearchOptions{
		Type: ItemType("unknown"),
	})
	if !errors.Is(err, ErrRequestFailed) {
		t.Fatalf("SearchItems() error = %v, want ErrRequestFailed", err)
	}
}

// TestSearchItemsInvalidAdultFilter は不正な adult を弾くことを確認します。
func TestSearchItemsInvalidAdultFilter(t *testing.T) {
	t.Parallel()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.SearchItems(context.Background(), SearchOptions{
		Adult: AdultFilter("invalid"),
	})
	if !errors.Is(err, ErrRequestFailed) {
		t.Fatalf("SearchItems() error = %v, want ErrRequestFailed", err)
	}
}

// TestGetShopNotFound は 404 がショップ未検出へ変換されることを確認します。
func TestGetShopNotFound(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		WithHTTPClient(mockHTTPClient{
			do: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader("not found")),
				}, nil
			},
		}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.GetShop(context.Background(), "sample.booth.pm")
	if !errors.Is(err, ErrShopNotFound) {
		t.Fatalf("GetShop() error = %v, want ErrShopNotFound", err)
	}
}

// TestGetShop はショップ取得の正常系を確認します。
func TestGetShop(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		WithHTTPClient(mockHTTPClient{
			do: func(req *http.Request) (*http.Response, error) {
				if req.URL.String() != "https://sample.booth.pm" {
					t.Fatalf("unexpected url: %s", req.URL.String())
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`<!doctype html><html><head>
						<link rel="canonical" href="https://sample.booth.pm">
						<meta property="og:site_name" content="サンプルショップ">
						</head><body><h1>サンプルショップ</h1></body></html>`)),
				}, nil
			},
		}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	shop, err := client.GetShop(context.Background(), "sample.booth.pm")
	if err != nil {
		t.Fatalf("GetShop() error = %v", err)
	}
	if shop.Name != "サンプルショップ" || shop.Host != "sample.booth.pm" {
		t.Fatalf("shop = %+v", shop)
	}
}

// TestSearchItemsTooManyRequests は 429 が変換されることを確認します。
func TestSearchItemsTooManyRequests(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		WithHTTPClient(mockHTTPClient{
			do: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusTooManyRequests,
					Body:       io.NopCloser(strings.NewReader("rate limited")),
				}, nil
			},
		}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.SearchItems(context.Background(), SearchOptions{Query: "sample"})
	if !errors.Is(err, ErrTooManyRequests) {
		t.Fatalf("SearchItems() error = %v, want ErrTooManyRequests", err)
	}
}

// TestSearchItemsNotFound は 404 が検索未検出へ変換されることを確認します。
func TestSearchItemsNotFound(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		WithHTTPClient(mockHTTPClient{
			do: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader("not found")),
				}, nil
			},
		}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.SearchItems(context.Background(), SearchOptions{Query: "sample"})
	if !errors.Is(err, ErrSearchNotFound) {
		t.Fatalf("SearchItems() error = %v, want ErrSearchNotFound", err)
	}
}

// TestSearchItemsRequestFailed は 5xx が ErrRequestFailed へ変換されることを確認します。
func TestSearchItemsRequestFailed(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		WithHTTPClient(mockHTTPClient{
			do: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader("server error")),
				}, nil
			},
		}),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.SearchItems(context.Background(), SearchOptions{Query: "sample"})
	if !errors.Is(err, ErrRequestFailed) {
		t.Fatalf("SearchItems() error = %v, want ErrRequestFailed", err)
	}
}

// TestSearchItemsInvalidSort は不正ソート値を弾くことを確認します。
func TestSearchItemsInvalidSort(t *testing.T) {
	t.Parallel()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.SearchItems(context.Background(), SearchOptions{Query: "sample", Sort: Sort("invalid")})
	if !errors.Is(err, ErrRequestFailed) {
		t.Fatalf("SearchItems() error = %v, want ErrRequestFailed", err)
	}
}

// TestDoGetCanceledContext は context cancel をそのまま返すことを確認します。
func TestDoGetCanceledContext(t *testing.T) {
	t.Parallel()

	client, err := NewClient(WithRateLimit(1000))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = client.SearchItems(ctx, SearchOptions{Query: "sample"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("SearchItems() error = %v, want context.Canceled", err)
	}
}

// TestWithTimeoutRejectsInvalidValue は不正タイムアウトを弾くことを確認します。
func TestWithTimeoutRejectsInvalidValue(t *testing.T) {
	t.Parallel()

	_, err := NewClient(WithTimeout(0))
	if !errors.Is(err, ErrRequestFailed) {
		t.Fatalf("NewClient() error = %v, want ErrRequestFailed", err)
	}
}

// TestWithRateLimitRejectsInvalidValue は不正レート制限を弾くことを確認します。
func TestWithRateLimitRejectsInvalidValue(t *testing.T) {
	t.Parallel()

	_, err := NewClient(WithRateLimit(0))
	if !errors.Is(err, ErrRequestFailed) {
		t.Fatalf("NewClient() error = %v, want ErrRequestFailed", err)
	}
}

// TestWithUserAgentRejectsInvalidValue は空 User-Agent を弾くことを確認します。
func TestWithUserAgentRejectsInvalidValue(t *testing.T) {
	t.Parallel()

	_, err := NewClient(WithUserAgent(" "))
	if !errors.Is(err, ErrRequestFailed) {
		t.Fatalf("NewClient() error = %v, want ErrRequestFailed", err)
	}
}

// TestWithLangRejectsInvalidValue は空言語設定を弾くことを確認します。
func TestWithLangRejectsInvalidValue(t *testing.T) {
	t.Parallel()

	_, err := NewClient(WithLang(" "))
	if !errors.Is(err, ErrRequestFailed) {
		t.Fatalf("NewClient() error = %v, want ErrRequestFailed", err)
	}
}

// TestWithTimeoutOnCustomHTTPClient は *http.Client 以外で timeout option を弾くことを確認します。
func TestWithTimeoutOnCustomHTTPClient(t *testing.T) {
	t.Parallel()

	_, err := NewClient(
		WithHTTPClient(mockHTTPClient{
			do: func(req *http.Request) (*http.Response, error) {
				return nil, nil
			},
		}),
		WithTimeout(time.Second),
	)
	if !errors.Is(err, ErrRequestFailed) {
		t.Fatalf("NewClient() error = %v, want ErrRequestFailed", err)
	}
}
