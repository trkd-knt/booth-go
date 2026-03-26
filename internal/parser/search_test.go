package parser

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// TestParseSearchPage は検索結果ページから商品一覧を抽出できることを確認します。
func TestParseSearchPage(t *testing.T) {
	t.Parallel()

	file, err := os.Open(filepath.Join("testdata", "search.html"))
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer file.Close()

	result, err := ParseSearchPage(file)
	if err != nil {
		t.Fatalf("ParseSearchPage() error = %v", err)
	}

	if len(result.Items) != 2 {
		t.Fatalf("len(result.Items) = %d, want 2", len(result.Items))
	}
	if result.Page != 2 {
		t.Fatalf("result.Page = %d, want 2", result.Page)
	}
	if !result.HasNext {
		t.Fatalf("result.HasNext = false, want true")
	}
	if result.TotalCount == nil || *result.TotalCount != 25 {
		t.Fatalf("result.TotalCount = %v, want 25", result.TotalCount)
	}
}

// TestParseSearchPageStructuredData は ItemList JSON からも抽出できることを確認します。
func TestParseSearchPageStructuredData(t *testing.T) {
	t.Parallel()

	file, err := os.Open(filepath.Join("testdata", "search_json.html"))
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer file.Close()

	result, err := ParseSearchPage(file)
	if err != nil {
		t.Fatalf("ParseSearchPage() error = %v", err)
	}

	if len(result.Items) != 1 {
		t.Fatalf("len(result.Items) = %d, want 1", len(result.Items))
	}
	if result.Items[0].PriceText != "1500" {
		t.Fatalf("result.Items[0].PriceText = %q", result.Items[0].PriceText)
	}
}

// TestParseSearchPageNoResults は検索結果が無い場合に失敗することを確認します。
func TestParseSearchPageNoResults(t *testing.T) {
	t.Parallel()

	file, err := os.Open(filepath.Join("testdata", "search_empty.html"))
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer file.Close()

	_, err = ParseSearchPage(file)
	if err == nil {
		t.Fatal("ParseSearchPage() error = nil, want error")
	}
	if !errors.Is(err, err) {
		t.Fatalf("ParseSearchPage() error = %v", err)
	}
}
