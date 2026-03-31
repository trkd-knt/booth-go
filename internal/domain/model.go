package domain

// Image は商品画像を表します。
type Image struct {
	Original string
	Resized  string
}

// Category は商品カテゴリを表します。
type Category struct {
	ID   int
	Name string
}

// ShopPreview は商品に紐づくショップ概要を表します。
type ShopPreview struct {
	Name      string
	Host      string
	URL       string
	Thumbnail string
}

// Downloadable はダウンロード可能なファイル情報を表します。
type Downloadable struct {
	FileName      string
	FileExtension string
	FileSize      string
	Name          string
	URL           string
}

// Item は BOOTH の商品情報を表します。
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

// Shop は BOOTH のショップ情報を表します。
type Shop struct {
	Name string
	Host string
	URL  string
}

// SearchOptions は BOOTH 全体検索の条件を表します。
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

// Sort は BOOTH の検索ソート条件を表します。
type Sort string

// ItemType は BOOTH の商品種別フィルタです。
type ItemType string

// AdultFilter は BOOTH の成人向けフィルタです。
type AdultFilter string

// SearchResult は BOOTH 全体検索の結果を表します。
type SearchResult struct {
	Items      []Item
	Page       int
	HasNext    bool
	TotalCount *int
}

const (
	SortDefault   Sort = ""
	SortNewest    Sort = "new"
	SortPopular   Sort = "popular"
	SortPriceAsc  Sort = "price_asc"
	SortPriceDesc Sort = "price_desc"
)

const (
	ItemTypeDigital  ItemType = "digital"
	ItemTypePhysical ItemType = "physical"
)

const (
	AdultFilterDefault AdultFilter = ""
	AdultFilterOnly    AdultFilter = "only"
	AdultFilterInclude AdultFilter = "include"
)
