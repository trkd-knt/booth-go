package parser

import (
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/trkd-knt/booth-go/internal/model"
)

var totalCountPattern = regexp.MustCompile(`対象商品\s*([\d,]+)\s*件`)

// ParseSearchPage は検索結果ページ HTML から検索結果を抽出します。
func ParseSearchPage(reader io.Reader) (*model.SearchResult, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	result := parseSearchFromStructuredData(doc)
	if result == nil {
		result = &model.SearchResult{}
	}
	fillSearchFromDOM(doc, result)

	if len(result.Items) == 0 {
		return nil, errors.New("search result items not found")
	}

	return result, nil
}

type searchStructuredData struct {
	Type            string `json:"@type"`
	NumberOfItems   int    `json:"numberOfItems"`
	ItemListElement []struct {
		URL    string `json:"url"`
		Name   string `json:"name"`
		Offers []struct {
			Price string `json:"price"`
		} `json:"offers"`
	} `json:"itemListElement"`
}

func parseSearchFromStructuredData(doc *goquery.Document) *model.SearchResult {
	for _, raw := range scriptJSONCandidates(doc) {
		var payload searchStructuredData
		if !decodeJSON(raw, &payload) || !strings.EqualFold(payload.Type, "ItemList") {
			continue
		}

		result := &model.SearchResult{}
		if payload.NumberOfItems > 0 {
			total := payload.NumberOfItems
			result.TotalCount = &total
		}
		for _, entry := range payload.ItemListElement {
			item := model.Item{
				Title: strings.TrimSpace(entry.Name),
				URL:   strings.TrimSpace(entry.URL),
				ID:    parseItemIDFromURL(entry.URL),
			}
			for _, offer := range entry.Offers {
				if item.Price == 0 {
					item.Price = parsePrice(offer.Price)
				}
				if item.PriceText == "" {
					item.PriceText = offer.Price
				}
			}
			if item.Title != "" && item.URL != "" {
				result.Items = append(result.Items, item)
			}
		}
		return result
	}
	return nil
}

func fillSearchFromDOM(doc *goquery.Document, result *model.SearchResult) {
	if result.TotalCount == nil {
		if matches := totalCountPattern.FindStringSubmatch(doc.Text()); len(matches) == 2 {
			value, err := strconv.Atoi(strings.ReplaceAll(matches[1], ",", ""))
			if err == nil {
				result.TotalCount = &value
			}
		}
	}

	if result.Page == 0 {
		result.Page = parseCurrentPage(doc)
	}
	result.HasNext = doc.Find(`a[rel="next"], .pagination .next a`).Length() > 0

	if len(result.Items) == 0 {
		result.Items = parseSearchItemsFromDOM(doc)
	}
}

func parseCurrentPage(doc *goquery.Document) int {
	current := strings.TrimSpace(doc.Find(`.pagination .current, [aria-current="page"]`).First().Text())
	if current == "" {
		return 1
	}
	value, err := strconv.Atoi(current)
	if err != nil || value <= 0 {
		return 1
	}
	return value
}

func parseSearchItemsFromDOM(doc *goquery.Document) []model.Item {
	seen := map[string]struct{}{}
	items := make([]model.Item, 0)

	doc.Find(`a[href*="/items/"]`).Each(func(_ int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok || href == "" {
			return
		}
		if _, exists := seen[href]; exists {
			return
		}

		title := strings.TrimSpace(s.Text())
		container := s.Parent()
		for i := 0; i < 3 && container.Length() > 0; i++ {
			if title != "" {
				break
			}
			title = strings.TrimSpace(container.Find(`a[href*="/items/"]`).First().Text())
			container = container.Parent()
		}
		if title == "" {
			return
		}

		item := model.Item{
			ID:    parseItemIDFromURL(href),
			Title: title,
			URL:   href,
		}

		node := s.Parent()
		for i := 0; i < 4 && node.Length() > 0; i++ {
			if item.Price != 0 && item.Shop != nil {
				break
			}
			text := strings.TrimSpace(node.Text())
			if item.Price == 0 {
				item.Price = parsePrice(text)
				if item.Price > 0 && item.PriceText == "" {
					item.PriceText = text
				}
			}
			if item.Shop == nil {
				node.Find("a").Each(func(_ int, candidate *goquery.Selection) {
					link, ok := candidate.Attr("href")
					host := parseURLHost(link)
					if ok && host != "" && !strings.Contains(link, "/items/") && item.Shop == nil {
						item.Shop = &model.ShopPreview{
							Host: host,
							URL:  link,
						}
					}
				})
			}
			node = node.Parent()
		}

		item.IsSoldOut = strings.Contains(strings.ToLower(s.Parent().Text()), "在庫なし")
		items = append(items, item)
		seen[href] = struct{}{}
	})

	return items
}
