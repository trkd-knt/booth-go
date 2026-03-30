package model

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
	// Query は検索キーワードです。空文字を許可します。
	Query string
	// Category は BOOTH の browse URL に使用するカテゴリ識別子です。
	Category string
	// ExceptWords は除外ワードです。except_words[] クエリとして送信されます。
	ExceptWords []string
	// Tags はタグ絞り込みです。tags[] クエリとして送信されます。
	Tags []string
	// Event はイベント識別子です。event クエリとして送信されます。
	Event string
	// Type は商品種別です。type クエリとして送信されます。
	Type ItemType
	// Adult は成人向けフィルタです。adult クエリとして送信されます。
	Adult AdultFilter
	// MinPrice は最低価格です。min_price クエリとして送信されます。
	MinPrice int
	// MaxPrice は最高価格です。max_price クエリとして送信されます。
	MaxPrice int
	// Sort は BOOTH 側で有効な値のみを受け付けます。
	Sort Sort
	// Page は 1 以上を指定します。0 は既定値として 1 扱いです。
	Page int
	// OnlyAvailable が true の場合は在庫あり商品のみを対象にします。
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
	// SortDefault は BOOTH 側の既定ソートを利用します。
	SortDefault Sort = ""
	// SortNewest は新着順を表します。
	SortNewest Sort = "new"
	// SortPopular は人気順を表します。
	SortPopular Sort = "popular"
	// SortPriceAsc は価格昇順を表します。
	SortPriceAsc Sort = "price_asc"
	// SortPriceDesc は価格降順を表します。
	SortPriceDesc Sort = "price_desc"
)

const (
	// ItemTypeDigital はデジタル商品を表します。
	ItemTypeDigital ItemType = "digital"
	// ItemTypePhysical は物理商品を表します。
	ItemTypePhysical ItemType = "physical"
)

const (
	// AdultFilterDefault は成人向けフィルタを指定しません。
	AdultFilterDefault AdultFilter = ""
	// AdultFilterOnly は成人向け商品のみを対象にします。
	AdultFilterOnly AdultFilter = "only"
	// AdultFilterInclude は成人向け商品を含めます。
	AdultFilterInclude AdultFilter = "include"
)
