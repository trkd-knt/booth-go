package parser

import (
	"encoding/json"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var pricePattern = regexp.MustCompile(`\d[\d,]*`)

// canonicalURL は canonical link を返します。
func canonicalURL(doc *goquery.Document) string {
	if href, ok := doc.Find(`link[rel="canonical"]`).Attr("href"); ok {
		return strings.TrimSpace(href)
	}
	return ""
}

// metaContent は meta タグの content を返します。
func metaContent(doc *goquery.Document, selector string) string {
	if content, ok := doc.Find(selector).Attr("content"); ok {
		return strings.TrimSpace(content)
	}
	return ""
}

// scriptJSONCandidates は JSON を含む script テキストを列挙します。
func scriptJSONCandidates(doc *goquery.Document) []string {
	var scripts []string
	doc.Find("script").Each(func(_ int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text == "" {
			return
		}
		if strings.HasPrefix(text, "{") || strings.HasPrefix(text, "[") {
			scripts = append(scripts, text)
			return
		}
		if strings.Contains(text, "{") && strings.Contains(text, "}") {
			scripts = append(scripts, text)
		}
	})
	return scripts
}

// parsePrice は文字列から価格を抽出します。
func parsePrice(text string) int {
	match := pricePattern.FindString(text)
	if match == "" {
		return 0
	}
	value, err := strconv.Atoi(strings.ReplaceAll(match, ",", ""))
	if err != nil {
		return 0
	}
	return value
}

// parseURLHost は URL のホスト名を返します。
func parseURLHost(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return parsed.Host
}

// parseItemIDFromURL は URL から item ID を抽出します。
func parseItemIDFromURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	base := path.Base(parsed.Path)
	if base == "." || base == "/" {
		return ""
	}
	return base
}

// uniqueStrings は文字列スライスの重複を除きます。
func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

// decodeJSON は任意 JSON をデコードします。
func decodeJSON(raw string, dest any) bool {
	return json.Unmarshal([]byte(raw), dest) == nil
}

// findJSONValueByKey は JSON 互換データから指定キーを再帰的に探します。
func findJSONValueByKey(data any, keys ...string) (any, bool) {
	keySet := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		keySet[strings.ToLower(key)] = struct{}{}
	}
	return findJSONValueByKeySet(data, keySet)
}

func findJSONValueByKeySet(data any, keys map[string]struct{}) (any, bool) {
	switch value := data.(type) {
	case map[string]any:
		for key, child := range value {
			if _, ok := keys[strings.ToLower(key)]; ok {
				return child, true
			}
		}
		for _, child := range value {
			if found, ok := findJSONValueByKeySet(child, keys); ok {
				return found, true
			}
		}
	case []any:
		for _, child := range value {
			if found, ok := findJSONValueByKeySet(child, keys); ok {
				return found, true
			}
		}
	}
	return nil, false
}

// asString は値を文字列へ変換します。
func asString(value any) string {
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case json.Number:
		return typed.String()
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case int:
		return strconv.Itoa(typed)
	default:
		return ""
	}
}

// asInt は値を整数へ変換します。
func asInt(value any) int {
	switch typed := value.(type) {
	case int:
		return typed
	case float64:
		return int(typed)
	case json.Number:
		number, err := typed.Int64()
		if err == nil {
			return int(number)
		}
	case string:
		return parsePrice(typed)
	}
	return 0
}
