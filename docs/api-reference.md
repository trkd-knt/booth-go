# API Reference

## Client

### `NewClient(opts ...Option) (*Client, error)`

`Client` を生成します。

主な option:

- `WithLang(lang string)`
- `WithRateLimit(requestsPerSecond float64)`
- `WithTimeout(timeout time.Duration)`
- `WithUserAgent(userAgent string)`
- `WithHTTPClient(httpClient HTTPClient)`
- `WithSearchBaseURL(rawURL string)`

## Public Methods

### `GetItem(ctx context.Context, itemID string) (*Item, error)`

商品 ID を指定して商品詳細を取得します。

取得元:

- 商品本体: 商品 JSON
- 対応アバター: 商品 HTML

### `GetItemDescription(ctx context.Context, itemID string) (string, error)`

商品 ID を指定して商品ページ HTML から全文説明を取得します。

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

各フィールドは BOOTH 検索クエリに対応します。

- `Category` は `/browse/{category}` のパスに使われます
- `Query` は `q`
- `ExceptWords` は `except_words[]`
- `Tags` は `tags[]`
- `Event` は `event`
- `Type` は `type`
- `Adult` は `adult`
- `MinPrice` は `min_price`
- `MaxPrice` は `max_price`
- `OnlyAvailable` は `in_stock`

### `Sort`

```go
const (
	SortDefault   Sort = ""
	SortNewest    Sort = "new"
	SortPopular   Sort = "popular"
	SortPriceAsc  Sort = "price_asc"
	SortPriceDesc Sort = "price_desc"
)
```

### `ItemType`

```go
const (
	ItemTypeDigital  ItemType = "digital"
	ItemTypePhysical ItemType = "physical"
)
```

### `AdultFilter`

```go
const (
	AdultFilterDefault AdultFilter = ""
	AdultFilterOnly    AdultFilter = "only"
	AdultFilterInclude AdultFilter = "include"
)
```

### `Item`

```go
type Item struct {
	ID            string
	Title         string
	Price         int
	PriceText     string
	ShopHost      string
	Tags          []string
	Avatars       []string
	Images        []string
	ImageDetails  []Image
	Summary       string
	Description   string
	URL           string
	IsSoldOut     bool
	Category      *Category
	Shop          *ShopPreview
	IsAdult       bool
	Likes         int
	Downloadables []Downloadable
}
```

補足:

- `Summary` は商品 JSON に含まれる説明文です
- `Description` は `GetItem` では埋まりません
- 全文説明が必要な場合は `GetItemDescription` を使います
- `SearchItems` の `Items` では `ShopHost` は保証しません
- `Avatars` は商品ページ HTML に存在する場合のみ返ります

### `Image`

```go
type Image struct {
	Original string
	Resized  string
}
```

### `Category`

```go
type Category struct {
	ID   int
	Name string
}
```

### `ShopPreview`

```go
type ShopPreview struct {
	Name      string
	Host      string
	URL       string
	Thumbnail string
}
```

### `Downloadable`

```go
type Downloadable struct {
	FileName      string
	FileExtension string
	FileSize      string
	Name          string
	URL           string
}
```

### `Shop`

```go
type Shop struct {
	Name string
	Host string
	URL  string
}
```

### `SearchResult`

```go
type SearchResult struct {
	Items      []Item
	Page       int
	HasNext    bool
	TotalCount *int
}
```

## Errors

主な公開エラー:

- `ErrItemNotFound`
- `ErrShopNotFound`
- `ErrSearchNotFound`
- `ErrTooManyRequests`
- `ErrParseFailed`
- `ErrItemAvatarsParseFailed`
- `ErrItemDescriptionParseFailed`
- `ErrRequestFailed`
