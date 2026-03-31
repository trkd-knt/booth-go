package parser

import (
	"errors"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/trkd-knt/booth-go/internal/domain"
)

// ParseShopPage はショップページ HTML からショップ情報を抽出します。
func ParseShopPage(reader io.Reader) (*domain.Shop, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	shop := parseShopFromStructuredData(doc)
	if shop == nil {
		shop = &domain.Shop{}
	}
	fillShopFromDOM(doc, shop)

	if shop.Name == "" {
		return nil, errors.New("shop name not found")
	}

	return shop, nil
}

type shopStructuredData struct {
	Type string `json:"@type"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

func parseShopFromStructuredData(doc *goquery.Document) *domain.Shop {
	for _, raw := range scriptJSONCandidates(doc) {
		var single shopStructuredData
		if decodeJSON(raw, &single) && isStoreType(single.Type) {
			return &domain.Shop{
				Name: strings.TrimSpace(single.Name),
				URL:  strings.TrimSpace(single.URL),
				Host: parseURLHost(single.URL),
			}
		}

		var list []shopStructuredData
		if decodeJSON(raw, &list) {
			for _, entry := range list {
				if isStoreType(entry.Type) {
					return &domain.Shop{
						Name: strings.TrimSpace(entry.Name),
						URL:  strings.TrimSpace(entry.URL),
						Host: parseURLHost(entry.URL),
					}
				}
			}
		}
	}
	return nil
}

func isStoreType(value string) bool {
	return strings.EqualFold(value, "store") || strings.EqualFold(value, "organization")
}

func fillShopFromDOM(doc *goquery.Document, shop *domain.Shop) {
	if shop.URL == "" {
		shop.URL = canonicalURL(doc)
	}
	if shop.URL == "" {
		shop.URL = metaContent(doc, `meta[property="og:url"]`)
	}
	if shop.Name == "" {
		shop.Name = metaContent(doc, `meta[property="og:site_name"]`)
	}
	if shop.Name == "" {
		shop.Name = strings.TrimSpace(doc.Find("title").First().Text())
	}
	if shop.Name == "" {
		shop.Name = strings.TrimSpace(doc.Find("h1").First().Text())
	}
	if shop.Host == "" {
		shop.Host = parseURLHost(shop.URL)
	}
}
