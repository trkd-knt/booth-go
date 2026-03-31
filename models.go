package booth

import "github.com/trkd-knt/booth-go/internal/domain"

// Item は BOOTH の商品情報を表します。
type Item = domain.Item

// Image は商品画像を表します。
type Image = domain.Image

// Category は商品カテゴリを表します。
type Category = domain.Category

// ShopPreview は商品に紐づくショップ概要を表します。
type ShopPreview = domain.ShopPreview

// Downloadable はダウンロード可能なファイル情報を表します。
type Downloadable = domain.Downloadable

// Shop は BOOTH のショップ情報を表します。
type Shop = domain.Shop

// SearchOptions は BOOTH 全体検索の条件を表します。
type SearchOptions = domain.SearchOptions

// Sort は BOOTH の検索ソート条件を表します。
type Sort = domain.Sort

// ItemType は BOOTH の商品種別フィルタです。
type ItemType = domain.ItemType

// AdultFilter は BOOTH の成人向けフィルタです。
type AdultFilter = domain.AdultFilter

// SearchResult は BOOTH 全体検索の結果を表します。
type SearchResult = domain.SearchResult

const (
	// SortDefault は BOOTH 側の既定ソートを利用します。
	SortDefault = domain.SortDefault
	// SortNewest は新着順を表します。
	SortNewest = domain.SortNewest
	// SortPopular は人気順を表します。
	SortPopular = domain.SortPopular
	// SortPriceAsc は価格昇順を表します。
	SortPriceAsc = domain.SortPriceAsc
	// SortPriceDesc は価格降順を表します。
	SortPriceDesc = domain.SortPriceDesc
)

const (
	// ItemTypeDigital はデジタル商品を表します。
	ItemTypeDigital = domain.ItemTypeDigital
	// ItemTypePhysical は物理商品を表します。
	ItemTypePhysical = domain.ItemTypePhysical
)

const (
	// AdultFilterDefault は成人向けフィルタを指定しません。
	AdultFilterDefault = domain.AdultFilterDefault
	// AdultFilterOnly は成人向け商品のみを対象にします。
	AdultFilterOnly = domain.AdultFilterOnly
	// AdultFilterInclude は成人向け商品を含めます。
	AdultFilterInclude = domain.AdultFilterInclude
)
