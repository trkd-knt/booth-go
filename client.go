package booth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/trkd-knt/booth-go/internal/boothhttp"
	"github.com/trkd-knt/booth-go/internal/domain"
	"golang.org/x/time/rate"
)

const (
	defaultUserAgent      = "booth-go/0.1 (+https://github.com/trkd-knt/booth-go)"
	defaultTimeout        = 10 * time.Second
	defaultRequestsPerSec = 1
	defaultLang           = "ja"
)

var (
	// ErrItemNotFound は商品が見つからなかったことを表します。
	ErrItemNotFound = errors.New("booth: item not found")
	// ErrShopNotFound はショップが見つからなかったことを表します。
	ErrShopNotFound = errors.New("booth: shop not found")
	// ErrSearchNotFound は検索結果ページが見つからなかったことを表します。
	ErrSearchNotFound = errors.New("booth: search not found")
	// ErrTooManyRequests は BOOTH 側からレート超過を返されたことを表します。
	ErrTooManyRequests = errors.New("booth: too many requests")
	// ErrParseFailed はレスポンスの解析に失敗したことを表します。
	ErrParseFailed = errors.New("booth: parse failed")
	// ErrItemAvatarsParseFailed は商品詳細の対応アバター補完に失敗したことを表します。
	ErrItemAvatarsParseFailed = errors.New("booth: item avatars parse failed")
	// ErrItemDescriptionParseFailed は商品詳細の本文説明抽出に失敗したことを表します。
	ErrItemDescriptionParseFailed = errors.New("booth: item description parse failed")
	// ErrRequestFailed は HTTP リクエスト自体に失敗したことを表します。
	ErrRequestFailed = errors.New("booth: request failed")
)

// HTTPClient は HTTP リクエストを実行するための最小インターフェースです。
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client は BOOTH へのアクセスをまとめるクライアントです。
type Client struct {
	httpClient    HTTPClient
	limiter       *rate.Limiter
	userAgent     string
	lang          string
	searchBaseURL *url.URL
}

// Option は Client の生成時設定です。
type Option func(*Client) error

// NewClient はデフォルト設定を持つ Client を生成します。
func NewClient(opts ...Option) (*Client, error) {
	searchBaseURL, err := boothhttp.BuildSearchBaseURL(defaultLang)
	if err != nil {
		return nil, fmt.Errorf("%w: build default search url: %v", ErrRequestFailed, err)
	}

	client := &Client{
		httpClient:    &http.Client{Timeout: defaultTimeout},
		limiter:       rate.NewLimiter(rate.Limit(defaultRequestsPerSec), 1),
		userAgent:     defaultUserAgent,
		lang:          defaultLang,
		searchBaseURL: searchBaseURL,
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	if client.httpClient == nil {
		return nil, fmt.Errorf("%w: http client is nil", ErrRequestFailed)
	}
	if client.limiter == nil {
		client.limiter = rate.NewLimiter(rate.Limit(defaultRequestsPerSec), 1)
	}
	if client.userAgent == "" {
		client.userAgent = defaultUserAgent
	}
	if client.lang == "" {
		client.lang = defaultLang
	}
	if client.searchBaseURL == nil {
		client.searchBaseURL, err = boothhttp.BuildSearchBaseURL(client.lang)
		if err != nil {
			return nil, fmt.Errorf("%w: build search url: %v", ErrRequestFailed, err)
		}
	}

	return client, nil
}

// MustNewClient はエラー時に panic する Client 生成ヘルパーです。
func MustNewClient(opts ...Option) *Client {
	client, err := NewClient(opts...)
	if err != nil {
		panic(err)
	}
	return client
}

// WithHTTPClient は HTTP 実装を差し替えます。
func WithHTTPClient(httpClient HTTPClient) Option {
	return func(c *Client) error {
		if httpClient == nil {
			return fmt.Errorf("%w: http client is nil", ErrRequestFailed)
		}
		c.httpClient = httpClient
		return nil
	}
}

// WithTimeout はデフォルト http.Client のタイムアウトを設定します。
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) error {
		if timeout <= 0 {
			return fmt.Errorf("%w: timeout must be positive", ErrRequestFailed)
		}

		httpClient, ok := c.httpClient.(*http.Client)
		if !ok {
			return fmt.Errorf("%w: timeout option requires *http.Client", ErrRequestFailed)
		}
		httpClient.Timeout = timeout
		return nil
	}
}

// WithUserAgent は送信する User-Agent を設定します。
func WithUserAgent(userAgent string) Option {
	return func(c *Client) error {
		userAgent = strings.TrimSpace(userAgent)
		if userAgent == "" {
			return fmt.Errorf("%w: user agent must not be empty", ErrRequestFailed)
		}
		c.userAgent = userAgent
		return nil
	}
}

// WithRateLimit は 1 秒あたりのリクエスト数を設定します。
func WithRateLimit(requestsPerSecond float64) Option {
	return func(c *Client) error {
		if requestsPerSecond <= 0 {
			return fmt.Errorf("%w: requests per second must be positive", ErrRequestFailed)
		}
		c.limiter = rate.NewLimiter(rate.Limit(requestsPerSecond), 1)
		return nil
	}
}

// WithSearchBaseURL は検索ページのベース URL を差し替えます。
func WithSearchBaseURL(rawURL string) Option {
	return func(c *Client) error {
		parsed, err := url.Parse(rawURL)
		if err != nil {
			return fmt.Errorf("%w: parse search base url: %v", ErrRequestFailed, err)
		}
		if parsed.Scheme == "" || parsed.Host == "" {
			return fmt.Errorf("%w: search base url must be absolute", ErrRequestFailed)
		}
		c.searchBaseURL = parsed
		return nil
	}
}

// WithLang は BOOTH ページの言語を設定します。
func WithLang(lang string) Option {
	return func(c *Client) error {
		lang = strings.TrimSpace(lang)
		if lang == "" {
			return fmt.Errorf("%w: lang must not be empty", ErrRequestFailed)
		}
		c.lang = lang

		searchBaseURL, err := boothhttp.BuildSearchBaseURL(lang)
		if err != nil {
			return fmt.Errorf("%w: build search url: %v", ErrRequestFailed, err)
		}
		c.searchBaseURL = searchBaseURL
		return nil
	}
}

// doGet は GET リクエストを実行してレスポンスボディを返します。
func (c *Client) doGet(ctx context.Context, rawURL string, resourceKind string) (io.ReadCloser, error) {
	resp, err := boothhttp.ExecuteGET(ctx, c.httpClient, c.limiter, c.userAgent, c.lang, rawURL)
	if err != nil {
		if _, ok := err.(*url.Error); !ok && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("%w: build request: %v", ErrRequestFailed, err)
		}
		return nil, wrapRequestError(err)
	}

	if err := mapStatusError(resourceKind, resp.StatusCode); err != nil {
		defer resp.Body.Close()
		return nil, err
	}

	return resp.Body, nil
}

// validateSearchOptions は検索条件を検証します。
func validateSearchOptions(opts SearchOptions) error {
	if opts.Page < 0 {
		return fmt.Errorf("%w: page must be greater than or equal to 0", ErrRequestFailed)
	}
	if opts.MinPrice < 0 {
		return fmt.Errorf("%w: min price must be greater than or equal to 0", ErrRequestFailed)
	}
	if opts.MaxPrice < 0 {
		return fmt.Errorf("%w: max price must be greater than or equal to 0", ErrRequestFailed)
	}
	if opts.MaxPrice > 0 && opts.MinPrice > opts.MaxPrice {
		return fmt.Errorf("%w: min price must be less than or equal to max price", ErrRequestFailed)
	}
	if opts.Page == 0 {
		opts.Page = 1
	}

	switch opts.Sort {
	case SortDefault, SortNewest, SortPopular, SortPriceAsc, SortPriceDesc:
	default:
		return fmt.Errorf("%w: invalid sort value %q", ErrRequestFailed, opts.Sort)
	}

	switch opts.Type {
	case "", ItemTypeDigital, ItemTypePhysical:
	default:
		return fmt.Errorf("%w: invalid item type %q", ErrRequestFailed, opts.Type)
	}

	switch opts.Adult {
	case AdultFilterDefault, AdultFilterOnly, AdultFilterInclude:
		return nil
	default:
		return fmt.Errorf("%w: invalid adult filter %q", ErrRequestFailed, opts.Adult)
	}

	return nil
}

// mapStatusError は HTTP ステータスを独自エラーへ変換します。
func mapStatusError(resourceKind string, statusCode int) error {
	switch statusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		if resourceKind == "shop" {
			return ErrShopNotFound
		}
		if resourceKind == "search" {
			return ErrSearchNotFound
		}
		return ErrItemNotFound
	case http.StatusTooManyRequests:
		return ErrTooManyRequests
	default:
		if statusCode >= 400 {
			return fmt.Errorf("%w: status=%d", ErrRequestFailed, statusCode)
		}
		return nil
	}
}

// wrapRequestError は context 系エラーを保ったままリクエストエラーを包みます。
func wrapRequestError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	return fmt.Errorf("%w: %v", ErrRequestFailed, err)
}

// normalizeSearchOptions は不足値を既定値で補います。
func normalizeSearchOptions(opts SearchOptions) domain.SearchOptions {
	if opts.Page == 0 {
		opts.Page = 1
	}
	return opts
}
