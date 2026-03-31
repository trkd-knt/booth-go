package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/trkd-knt/booth-go/internal/domain"
)

type itemJSONResponse struct {
	Description    string            `json:"description"`
	ID             int64             `json:"id"`
	IsAdult        bool              `json:"is_adult"`
	IsSoldOut      bool              `json:"is_sold_out"`
	Name           string            `json:"name"`
	Price          string            `json:"price"`
	URL            string            `json:"url"`
	WishListsCount int               `json:"wish_lists_count"`
	Category       *itemJSONCategory `json:"category"`
	Images         []itemJSONImage   `json:"images"`
	Shop           *itemJSONShop     `json:"shop"`
	Tags           []itemJSONTag     `json:"tags"`
	Variations     []itemJSONVariant `json:"variations"`
}

type itemJSONCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type itemJSONImage struct {
	Original string `json:"original"`
	Resized  string `json:"resized"`
}

type itemJSONShop struct {
	Name         string `json:"name"`
	Subdomain    string `json:"subdomain"`
	ThumbnailURL string `json:"thumbnail_url"`
	URL          string `json:"url"`
}

type itemJSONTag struct {
	Name string `json:"name"`
}

type itemJSONVariant struct {
	Price int `json:"price"`
}

func DecodeItemJSON(body io.Reader) (*domain.Item, error) {
	var payload itemJSONResponse
	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return nil, err
	}

	item := &domain.Item{
		ID:          strconv.FormatInt(payload.ID, 10),
		Title:       strings.TrimSpace(payload.Name),
		Price:       deriveItemPrice(payload.Price, payload.Variations),
		PriceText:   strings.TrimSpace(payload.Price),
		Summary:     payload.Description,
		Description: payload.Description,
		URL:         strings.TrimSpace(payload.URL),
		IsSoldOut:   payload.IsSoldOut,
		IsAdult:     payload.IsAdult,
		Likes:       payload.WishListsCount,
	}

	if item.ID == "0" {
		item.ID = ""
	}

	if payload.Category != nil {
		item.Category = &domain.Category{
			ID:   payload.Category.ID,
			Name: strings.TrimSpace(payload.Category.Name),
		}
	}

	if payload.Shop != nil {
		item.Shop = &domain.ShopPreview{
			Name:      strings.TrimSpace(payload.Shop.Name),
			Host:      buildShopHostFromJSON(payload.Shop.Subdomain, payload.Shop.URL),
			URL:       strings.TrimSpace(payload.Shop.URL),
			Thumbnail: strings.TrimSpace(payload.Shop.ThumbnailURL),
		}
		item.ShopHost = item.Shop.Host
	}

	for _, image := range payload.Images {
		detail := domain.Image{
			Original: strings.TrimSpace(image.Original),
			Resized:  strings.TrimSpace(image.Resized),
		}
		if detail.Original == "" && detail.Resized == "" {
			continue
		}
		item.ImageDetails = append(item.ImageDetails, detail)
		if detail.Resized != "" {
			item.Images = append(item.Images, detail.Resized)
			continue
		}
		item.Images = append(item.Images, detail.Original)
	}

	for _, tag := range payload.Tags {
		name := strings.TrimSpace(tag.Name)
		if name == "" {
			continue
		}
		item.Tags = append(item.Tags, name)
	}

	return item, nil
}

func NormalizeItem(item *domain.Item, itemURL string) {
	if item.URL == "" {
		item.URL = itemURL
	}
	if item.ShopHost == "" {
		item.ShopHost = parseURLHost(item.URL)
	}
	if item.Shop == nil {
		item.Shop = &domain.ShopPreview{}
	}
	if item.Shop.Host == "" {
		item.Shop.Host = item.ShopHost
	}
	if item.Shop.URL == "" && item.Shop.Host != "" {
		item.Shop.URL = "https://" + item.Shop.Host
	}
	if item.ID == "" {
		if parsed, err := url.Parse(item.URL); err == nil {
			item.ID = path.Base(parsed.Path)
		}
	}
}

func ValidateDecodedItem(item *domain.Item) error {
	if item == nil {
		return fmt.Errorf("decoded item is nil")
	}
	if item.ID == "" && item.Title == "" {
		return fmt.Errorf("decoded item is empty")
	}
	return nil
}

func buildShopHostFromJSON(subdomain, rawURL string) string {
	subdomain = strings.TrimSpace(subdomain)
	if subdomain != "" {
		return subdomain + ".booth.pm"
	}
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return ""
	}
	return parsed.Host
}

func deriveItemPrice(priceText string, variations []itemJSONVariant) int {
	if price := parsePriceNumber(priceText); price > 0 {
		return price
	}
	for _, variation := range variations {
		if variation.Price > 0 {
			return variation.Price
		}
	}
	return 0
}

var priceDigitsPattern = regexp.MustCompile(`\d+`)

func parsePriceNumber(priceText string) int {
	matches := priceDigitsPattern.FindAllString(priceText, -1)
	if len(matches) == 0 {
		return 0
	}
	value, err := strconv.Atoi(strings.Join(matches, ""))
	if err != nil {
		return 0
	}
	return value
}
