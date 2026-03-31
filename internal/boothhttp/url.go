package boothhttp

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/trkd-knt/booth-go/internal/domain"
)

func BuildItemURL(lang, itemID string) (string, error) {
	itemID = strings.TrimSpace(itemID)
	if itemID == "" {
		return "", fmt.Errorf("item id must not be empty")
	}
	lang = strings.Trim(strings.TrimSpace(lang), "/")
	if lang == "" {
		lang = "ja"
	}
	return "https://booth.pm/" + lang + "/items/" + itemID, nil
}

func BuildItemJSONURL(lang, itemID string) (string, error) {
	itemURL, err := BuildItemURL(lang, itemID)
	if err != nil {
		return "", err
	}
	return itemURL + ".json", nil
}

func BuildShopURL(shopHost string) (string, error) {
	if strings.TrimSpace(shopHost) == "" {
		return "", fmt.Errorf("shop host must not be empty")
	}
	return "https://" + strings.TrimSpace(shopHost), nil
}

func BuildSearchBaseURL(lang string) (*url.URL, error) {
	return url.Parse("https://booth.pm" + buildSearchPath(lang))
}

func NewSearchURL(baseURL *url.URL, lang string, opts domain.SearchOptions) (string, error) {
	if baseURL == nil {
		return "", fmt.Errorf("search base url is nil")
	}

	base := *baseURL
	if opts.Category != "" {
		base.Path = buildBrowsePath(lang, opts.Category)
	} else {
		base.Path = buildSearchPath(lang)
	}

	query := base.Query()
	if opts.Query != "" {
		query.Set("q", opts.Query)
	}
	for _, word := range opts.ExceptWords {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}
		query.Add("except_words[]", word)
	}
	for _, tag := range opts.Tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		query.Add("tags[]", tag)
	}
	if opts.Event != "" {
		query.Set("event", opts.Event)
	}
	if opts.Type != "" {
		query.Set("type", string(opts.Type))
	}
	if opts.Adult != "" {
		query.Set("adult", string(opts.Adult))
	}
	if opts.MinPrice > 0 {
		query.Set("min_price", fmt.Sprintf("%d", opts.MinPrice))
	}
	if opts.MaxPrice > 0 {
		query.Set("max_price", fmt.Sprintf("%d", opts.MaxPrice))
	}
	if opts.Sort != domain.SortDefault {
		query.Set("sort", string(opts.Sort))
	}
	if opts.Page > 1 {
		query.Set("page", fmt.Sprintf("%d", opts.Page))
	}
	if opts.OnlyAvailable {
		query.Set("only_available", "1")
	}
	base.RawQuery = query.Encode()

	return base.String(), nil
}

func ParseURLHost(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return parsed.Host
}

func buildSearchPath(lang string) string {
	return "/" + path.Join(lang, "search")
}

func buildBrowsePath(lang, category string) string {
	return "/" + path.Join(lang, "browse", category)
}
