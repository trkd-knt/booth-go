package parser

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// TestParseItemPage は JSON-LD を優先して商品情報を抽出できることを確認します。
func TestParseItemPage(t *testing.T) {
	t.Parallel()

	file, err := os.Open(filepath.Join("testdata", "item.html"))
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer file.Close()

	item, err := ParseItemPage(file)
	if err != nil {
		t.Fatalf("ParseItemPage() error = %v", err)
	}

	if item.ID != "12345" {
		t.Fatalf("item.ID = %q, want %q", item.ID, "12345")
	}
	if item.ShopHost != "sample.booth.pm" {
		t.Fatalf("item.ShopHost = %q", item.ShopHost)
	}
	if item.Title != "サンプル商品" {
		t.Fatalf("item.Title = %q", item.Title)
	}
	if item.Price != 1200 {
		t.Fatalf("item.Price = %d", item.Price)
	}
	if item.PriceText != "1200" {
		t.Fatalf("item.PriceText = %q", item.PriceText)
	}
	if item.Category == nil || item.Category.ID != 212 || item.Category.Name != "VRoid" {
		t.Fatalf("item.Category = %+v", item.Category)
	}
	if item.Shop == nil || item.Shop.Name != "サンプルショップ" || item.Shop.Host != "sample.booth.pm" {
		t.Fatalf("item.Shop = %+v", item.Shop)
	}
	if item.Likes != 42 {
		t.Fatalf("item.Likes = %d", item.Likes)
	}
	if len(item.Downloadables) != 1 {
		t.Fatalf("len(item.Downloadables) = %d, want 1", len(item.Downloadables))
	}
	if len(item.ImageDetails) != 2 {
		t.Fatalf("len(item.ImageDetails) = %d, want 2", len(item.ImageDetails))
	}
	if item.IsSoldOut {
		t.Fatalf("item.IsSoldOut = true, want false")
	}
}

// TestParseItemPageDOMFallback は構造化データが無い場合でも DOM から抽出できることを確認します。
func TestParseItemPageDOMFallback(t *testing.T) {
	t.Parallel()

	file, err := os.Open(filepath.Join("testdata", "item_dom.html"))
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer file.Close()

	item, err := ParseItemPage(file)
	if err != nil {
		t.Fatalf("ParseItemPage() error = %v", err)
	}

	if item.Title != "DOM商品" {
		t.Fatalf("item.Title = %q", item.Title)
	}
	if item.Price != 980 {
		t.Fatalf("item.Price = %d", item.Price)
	}
	if !item.IsSoldOut {
		t.Fatalf("item.IsSoldOut = false, want true")
	}
}

// TestParseItemPageNoTitle はタイトルが無い場合に失敗することを確認します。
func TestParseItemPageNoTitle(t *testing.T) {
	t.Parallel()

	file, err := os.Open(filepath.Join("testdata", "item_invalid.html"))
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer file.Close()

	_, err = ParseItemPage(file)
	if err == nil {
		t.Fatal("ParseItemPage() error = nil, want error")
	}
	if !errors.Is(err, err) {
		t.Fatalf("ParseItemPage() error = %v", err)
	}
}
