package parser

import (
	"errors"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var errItemAvatarsNotFound = errors.New("item avatars not found")

// ErrItemAvatarsNotFound は対応アバター一覧が HTML 上に見つからないことを表します。
func ErrItemAvatarsNotFound() error {
	return errItemAvatarsNotFound
}

// ParseItemDescriptionPage は商品詳細ページ HTML から本文説明を抽出します。
func ParseItemDescriptionPage(reader io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", err
	}

	description := parseItemDescriptionFromDOM(doc)
	if description == "" {
		return "", errors.New("item description not found")
	}

	return description, nil
}

// ParseItemAvatarsPage は商品詳細ページ HTML から対応アバター一覧を抽出します。
func ParseItemAvatarsPage(reader io.Reader) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	avatars := parseItemAvatarsFromDOM(doc)
	if len(avatars) == 0 {
		return nil, errItemAvatarsNotFound
	}

	return avatars, nil
}

func parseItemDescriptionFromDOM(doc *goquery.Document) string {
	sections := make([]string, 0)

	primary := extractCleanText(doc.Find(".js-market-item-detail-description").First(), "script,style")
	if primary != "" {
		sections = append(sections, primary)
	}

	doc.Find("article section.shop__text").Each(func(_ int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h2").First().Text())
		body := extractCleanText(s, "h2,script,style")
		if body == "" {
			return
		}
		if title == "" && primary != "" && body == primary {
			return
		}
		if title != "" {
			sections = append(sections, title+"\n"+body)
			return
		}
		sections = append(sections, body)
	})

	return strings.TrimSpace(strings.Join(sections, "\n\n"))
}

func extractCleanText(selection *goquery.Selection, removeSelectors string) string {
	if selection == nil || selection.Length() == 0 {
		return ""
	}
	clone := selection.Clone()
	if removeSelectors != "" {
		clone.Find(removeSelectors).Remove()
	}
	return strings.TrimSpace(clone.Text())
}

func parseItemAvatarsFromDOM(doc *goquery.Document) []string {
	avatars := make([]string, 0)

	doc.Find("#variations .variation-name").Each(func(_ int, s *goquery.Selection) {
		name := rawAvatarName(s.Text())
		if name == "" {
			return
		}
		avatars = append(avatars, name)
	})

	return uniqueStrings(avatars)
}

func rawAvatarName(raw string) string {
	return strings.TrimSpace(raw)
}
