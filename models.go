package booth

import "github.com/trkd-knt/booth-go/internal/model"

// Item は BOOTH の商品情報を表します。
type Item = model.Item

// Image は商品画像を表します。
type Image = model.Image

// Category は商品カテゴリを表します。
type Category = model.Category

// ShopPreview は商品に紐づくショップ概要を表します。
type ShopPreview = model.ShopPreview

// Downloadable はダウンロード可能なファイル情報を表します。
type Downloadable = model.Downloadable

// Shop は BOOTH のショップ情報を表します。
type Shop = model.Shop

// SearchOptions は BOOTH 全体検索の条件を表します。
type SearchOptions = model.SearchOptions

// Sort は BOOTH の検索ソート条件を表します。
type Sort = model.Sort

// ItemType は BOOTH の商品種別フィルタです。
type ItemType = model.ItemType

// AdultFilter は BOOTH の成人向けフィルタです。
type AdultFilter = model.AdultFilter

// SearchResult は BOOTH 全体検索の結果を表します。
type SearchResult = model.SearchResult

const (
	// SortDefault は BOOTH 側の既定ソートを利用します。
	SortDefault = model.SortDefault
	// SortNewest は新着順を表します。
	SortNewest = model.SortNewest
	// SortPopular は人気順を表します。
	SortPopular = model.SortPopular
	// SortPriceAsc は価格昇順を表します。
	SortPriceAsc = model.SortPriceAsc
	// SortPriceDesc は価格降順を表します。
	SortPriceDesc = model.SortPriceDesc
)

const (
	// ItemTypeDigital はデジタル商品を表します。
	ItemTypeDigital = model.ItemTypeDigital
	// ItemTypePhysical は物理商品を表します。
	ItemTypePhysical = model.ItemTypePhysical
)

const (
	// AdultFilterDefault は成人向けフィルタを指定しません。
	AdultFilterDefault = model.AdultFilterDefault
	// AdultFilterOnly は成人向け商品のみを対象にします。
	AdultFilterOnly = model.AdultFilterOnly
	// AdultFilterInclude は成人向け商品を含めます。
	AdultFilterInclude = model.AdultFilterInclude
)
