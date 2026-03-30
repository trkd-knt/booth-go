# API Reference

## Client

### `NewClient(opts ...Option) (*Client, error)`

`Client` を生成します。

利用可能な主な option:

- `WithLang(lang string)`
- `WithRateLimit(requestsPerSecond float64)`
- `WithTimeout(timeout time.Duration)`
- `WithUserAgent(userAgent string)`
- `WithHTTPClient(httpClient HTTPClient)`

## Public Methods

### `GetItem(ctx context.Context, itemID string) (*Item, error)`

商品 ID を指定して商品詳細を取得します。

### `GetItemDescription(ctx context.Context, itemID string) (string, error)`

商品 ID を指定して HTML 上の全文説明を取得します。

### `SearchItems(ctx context.Context, opts SearchOptions) (*SearchResult, error)`

BOOTH 全体検索結果を取得します。

### `GetShop(ctx context.Context, shopHost string) (*Shop, error)`

ショップ情報を取得します。

## Main Models

### `SearchOptions`

```go
type SearchOptions struct {
	Query         string
	Category      string
	ExceptWords   []string
	Tags          []string
	Event         string
	Type          ItemType
	Adult         AdultFilter
	MinPrice      int
	MaxPrice      int
	Sort          Sort
	Page          int
	OnlyAvailable bool
}
```

- `Category` は `/browse/{category}` のパス部分に使われます
- `Query` は `q`
- `ExceptWords` は `except_words[]`
- `Tags` は `tags[]`
- `Event` は `event`
- `Type` は `type`
- `Adult` は `adult`
- `MinPrice` は `min_price`
- `MaxPrice` は `max_price`

### `Sort`

```go
const (
	SortDefault
	SortNewest
	SortPopular
	SortPriceAsc
	SortPriceDesc
)
```

### `ItemType`

```go
const (
	ItemTypeDigital
	ItemTypePhysical
)
```

### `AdultFilter`

```go
const (
	AdultFilterDefault
	AdultFilterOnly
	AdultFilterInclude
)
```

### `Item`

主なフィールド:

- `ID`
- `Title`
- `Price`
- `PriceText`
- `ShopHost`
- `Images`
- `ImageDetails`
- `Description`
- `URL`
- `IsSoldOut`
- `Category`
- `Shop`
- `IsAdult`
- `Likes`
- `Downloadables`

補足:

- `GetItem` は JSON エンドポイントを基準に返却します
- `GetItemDescription` は HTML 本文抽出専用です
- `SearchItems` の `Items` では `ShopHost` は保証せず、ショップ情報がある場合は `Shop` を参照します

### `Shop`

主なフィールド:

- `Name`
- `Host`
- `URL`

### `SearchResult`

主なフィールド:

- `Items`
- `Page`
- `HasNext`
- `TotalCount`
