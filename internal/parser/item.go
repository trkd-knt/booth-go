package parser

import (
	"errors"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/trkd-knt/booth-go/internal/model"
)

// ParseItemPage は商品詳細ページ HTML から商品情報を抽出します。
func ParseItemPage(reader io.Reader) (*model.Item, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	item := parseItemFromStructuredData(doc)
	if item == nil {
		item = &model.Item{}
	}
	fillItemFromDOM(doc, item)

	if item.Title == "" {
		return nil, errors.New("item title not found")
	}

	return item, nil
}

type itemStructuredData struct {
	Type        string `json:"@type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       any    `json:"image"`
	URL         string `json:"url"`
	Offers      any    `json:"offers"`
}

func parseItemFromStructuredData(doc *goquery.Document) *model.Item {
	for _, raw := range scriptJSONCandidates(doc) {
		var generic any
		if decodeJSON(raw, &generic) {
			if item := itemFromGenericStructuredData(generic); item != nil {
				return item
			}
		}

		var single itemStructuredData
		if decodeJSON(raw, &single) && isProductType(single.Type) {
			return itemFromStructuredData(single)
		}

		var list []itemStructuredData
		if decodeJSON(raw, &list) {
			for _, entry := range list {
				if isProductType(entry.Type) {
					return itemFromStructuredData(entry)
				}
			}
		}
	}
	return nil
}

func isProductType(value string) bool {
	return strings.EqualFold(value, "product")
}

func itemFromStructuredData(data itemStructuredData) *model.Item {
	item := &model.Item{
		Title:       strings.TrimSpace(data.Name),
		Description: strings.TrimSpace(data.Description),
		URL:         strings.TrimSpace(data.URL),
	}

	switch images := data.Image.(type) {
	case string:
		item.Images = []string{images}
		item.ImageDetails = []model.Image{{Original: images}}
	case []any:
		for _, image := range images {
			if value, ok := image.(string); ok {
				item.Images = append(item.Images, value)
				item.ImageDetails = append(item.ImageDetails, model.Image{Original: value})
			}
		}
	}
	item.Images = uniqueStrings(item.Images)

	fillItemOfferDetails(item, data.Offers)

	item.ShopHost = parseURLHost(item.URL)
	item.ID = parseItemIDFromURL(item.URL)
	return item
}

func itemFromGenericStructuredData(data any) *model.Item {
	rootMap, ok := data.(map[string]any)
	if !ok {
		return nil
	}

	typeValue, _ := findJSONValueByKey(rootMap, "@type", "type")
	if !isProductType(asString(typeValue)) {
		return nil
	}

	item := &model.Item{}
	if value, ok := findJSONValueByKey(rootMap, "name", "title"); ok {
		item.Title = asString(value)
	}
	if value, ok := findJSONValueByKey(rootMap, "description"); ok {
		item.Description = asString(value)
	}
	if value, ok := findJSONValueByKey(rootMap, "url"); ok {
		item.URL = asString(value)
	}
	if value, ok := findJSONValueByKey(rootMap, "liked", "likes", "loveCount", "likeCount"); ok {
		item.Likes = asInt(value)
	}
	if value, ok := findJSONValueByKey(rootMap, "adult", "isAdult", "adultContent"); ok {
		if boolValue, ok := value.(bool); ok {
			item.IsAdult = boolValue
		}
	}
	if value, ok := findJSONValueByKey(rootMap, "category"); ok {
		item.Category = parseCategoryValue(value)
	}
	if value, ok := findJSONValueByKey(rootMap, "shop"); ok {
		item.Shop = parseShopPreviewValue(value)
		if item.Shop != nil && item.ShopHost == "" {
			item.ShopHost = item.Shop.Host
		}
	}
	if value, ok := findJSONValueByKey(rootMap, "downloadables", "downloadable"); ok {
		item.Downloadables = parseDownloadables(value)
	}
	if value, ok := findJSONValueByKey(rootMap, "images", "image"); ok {
		item.Images, item.ImageDetails = parseImagesValue(value)
	}
	if value, ok := findJSONValueByKey(rootMap, "offers", "offer"); ok {
		fillItemOfferDetails(item, value)
	}

	if item.URL != "" {
		item.ID = parseItemIDFromURL(item.URL)
		if item.ShopHost == "" {
			item.ShopHost = parseURLHost(item.URL)
		}
	}
	if item.Title == "" {
		return nil
	}
	return item
}

func priceFromOffer(offer map[string]any) int {
	if price, ok := offer["price"].(string); ok {
		return parsePrice(price)
	}
	if price, ok := offer["price"].(float64); ok {
		return int(price)
	}
	return 0
}

func soldOutFromOffer(offer map[string]any) bool {
	availability, _ := offer["availability"].(string)
	return strings.Contains(strings.ToLower(availability), "soldout") || strings.Contains(strings.ToLower(availability), "outofstock")
}

func firstPriceFromOffers(offers []any) int {
	for _, offer := range offers {
		offerMap, ok := offer.(map[string]any)
		if !ok {
			continue
		}
		if price := priceFromOffer(offerMap); price > 0 {
			return price
		}
	}
	return 0
}

func allOffersSoldOut(offers []any) bool {
	if len(offers) == 0 {
		return false
	}
	for _, offer := range offers {
		offerMap, ok := offer.(map[string]any)
		if !ok {
			return false
		}
		if !soldOutFromOffer(offerMap) {
			return false
		}
	}
	return true
}

func fillItemFromDOM(doc *goquery.Document, item *model.Item) {
	if item.URL == "" {
		item.URL = canonicalURL(doc)
	}
	if item.URL == "" {
		item.URL = metaContent(doc, `meta[property="og:url"]`)
	}
	if item.Title == "" {
		item.Title = metaContent(doc, `meta[property="og:title"]`)
	}
	if item.Title == "" {
		item.Title = strings.TrimSpace(doc.Find("h1").First().Text())
	}
	if item.Description == "" {
		item.Description = metaContent(doc, `meta[property="og:description"]`)
	}
	if item.Description == "" {
		item.Description = strings.TrimSpace(doc.Find(".item-description, .js-item-description, [data-test='item-description']").First().Text())
	}
	if len(item.Images) == 0 {
		doc.Find(`meta[property="og:image"]`).Each(func(_ int, s *goquery.Selection) {
			if content, ok := s.Attr("content"); ok {
				item.Images = append(item.Images, content)
				item.ImageDetails = append(item.ImageDetails, model.Image{Original: content})
			}
		})
	}
	if item.Price == 0 {
		item.Price = parsePrice(doc.Find(".price, .js-price, [data-test='item-price']").First().Text())
	}
	if item.PriceText == "" {
		item.PriceText = strings.TrimSpace(doc.Find(".price, .js-price, [data-test='item-price']").First().Text())
	}
	if item.Price == 0 {
		item.Price = parsePrice(doc.Text())
	}
	if !item.IsSoldOut {
		text := strings.ToLower(doc.Text())
		item.IsSoldOut = strings.Contains(text, "在庫なし") || strings.Contains(text, "sold out")
	}
	if item.ShopHost == "" {
		item.ShopHost = parseURLHost(item.URL)
	}
	if item.ID == "" {
		item.ID = parseItemIDFromURL(item.URL)
	}
	item.Images = uniqueStrings(item.Images)
	if item.Shop == nil && item.ShopHost != "" {
		item.Shop = &model.ShopPreview{
			Host: item.ShopHost,
			URL:  "https://" + item.ShopHost,
		}
	}
}

func fillItemOfferDetails(item *model.Item, value any) {
	switch offers := value.(type) {
	case map[string]any:
		if item.Price == 0 {
			item.Price = priceFromOffer(offers)
		}
		if item.PriceText == "" {
			item.PriceText = asString(offers["price"])
		}
		item.IsSoldOut = item.IsSoldOut || soldOutFromOffer(offers)
	case []any:
		if item.Price == 0 {
			item.Price = firstPriceFromOffers(offers)
		}
		if item.PriceText == "" {
			for _, offer := range offers {
				if offerMap, ok := offer.(map[string]any); ok {
					item.PriceText = asString(offerMap["price"])
					if item.PriceText != "" {
						break
					}
				}
			}
		}
		item.IsSoldOut = item.IsSoldOut || allOffersSoldOut(offers)
	}
}

func parseCategoryValue(value any) *model.Category {
	categoryMap, ok := value.(map[string]any)
	if !ok {
		return nil
	}
	category := &model.Category{
		ID:   asInt(categoryMap["id"]),
		Name: asString(categoryMap["name"]),
	}
	if category.ID == 0 && category.Name == "" {
		return nil
	}
	return category
}

func parseShopPreviewValue(value any) *model.ShopPreview {
	shopMap, ok := value.(map[string]any)
	if !ok {
		return nil
	}
	shop := &model.ShopPreview{
		Name:      asString(shopMap["name"]),
		Host:      asString(shopMap["subdomain"]),
		URL:       asString(shopMap["url"]),
		Thumbnail: asString(shopMap["thumbnail"]),
	}
	if shop.Host == "" {
		shop.Host = parseURLHost(shop.URL)
	}
	if shop.Name == "" && shop.Host == "" && shop.URL == "" && shop.Thumbnail == "" {
		return nil
	}
	return shop
}

func parseDownloadables(value any) []model.Downloadable {
	entries, ok := value.([]any)
	if !ok {
		return nil
	}
	result := make([]model.Downloadable, 0, len(entries))
	for _, entry := range entries {
		entryMap, ok := entry.(map[string]any)
		if !ok {
			continue
		}
		downloadable := model.Downloadable{
			FileName:      asString(entryMap["file_name"]),
			FileExtension: asString(entryMap["file_extension"]),
			FileSize:      asString(entryMap["file_size"]),
			Name:          asString(entryMap["name"]),
			URL:           asString(entryMap["url"]),
		}
		if downloadable.Name != "" || downloadable.URL != "" {
			result = append(result, downloadable)
		}
	}
	return result
}

func parseImagesValue(value any) ([]string, []model.Image) {
	switch images := value.(type) {
	case string:
		return []string{images}, []model.Image{{Original: images}}
	case []any:
		urls := make([]string, 0, len(images))
		details := make([]model.Image, 0, len(images))
		for _, entry := range images {
			switch image := entry.(type) {
			case string:
				urls = append(urls, image)
				details = append(details, model.Image{Original: image})
			case map[string]any:
				detail := model.Image{
					Original: asString(image["original"]),
					Resized:  asString(image["resized"]),
				}
				if detail.Original == "" {
					detail.Original = asString(image["url"])
				}
				if detail.Original != "" {
					urls = append(urls, detail.Original)
					details = append(details, detail)
				}
			}
		}
		return uniqueStrings(urls), details
	default:
		return nil, nil
	}
}
