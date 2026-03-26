package parser

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// TestParseShopPage はショップページからショップ情報を抽出できることを確認します。
func TestParseShopPage(t *testing.T) {
	t.Parallel()

	file, err := os.Open(filepath.Join("testdata", "shop.html"))
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer file.Close()

	shop, err := ParseShopPage(file)
	if err != nil {
		t.Fatalf("ParseShopPage() error = %v", err)
	}

	if shop.Name != "サンプルショップ" {
		t.Fatalf("shop.Name = %q", shop.Name)
	}
	if shop.Host != "sample.booth.pm" {
		t.Fatalf("shop.Host = %q", shop.Host)
	}
}

// TestParseShopPageDOMFallback は構造化データが無くても DOM から抽出できることを確認します。
func TestParseShopPageDOMFallback(t *testing.T) {
	t.Parallel()

	file, err := os.Open(filepath.Join("testdata", "shop_dom.html"))
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer file.Close()

	shop, err := ParseShopPage(file)
	if err != nil {
		t.Fatalf("ParseShopPage() error = %v", err)
	}

	if shop.Name != "DOMショップ" {
		t.Fatalf("shop.Name = %q", shop.Name)
	}
	if shop.Host != "dom.booth.pm" {
		t.Fatalf("shop.Host = %q", shop.Host)
	}
}

// TestParseShopPageNoName はショップ名が無い場合に失敗することを確認します。
func TestParseShopPageNoName(t *testing.T) {
	t.Parallel()

	file, err := os.Open(filepath.Join("testdata", "shop_invalid.html"))
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer file.Close()

	_, err = ParseShopPage(file)
	if err == nil {
		t.Fatal("ParseShopPage() error = nil, want error")
	}
	if !errors.Is(err, err) {
		t.Fatalf("ParseShopPage() error = %v", err)
	}
}
