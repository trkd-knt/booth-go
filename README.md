# booth-go

`booth-go` は BOOTH の商品詳細、商品説明、ショップ情報、全体検索結果を取得するための Go ライブラリです。

- 商品詳細は `https://booth.pm/{lang}/items/{id}.json` を基準に取得します
- 全文説明と対応アバターは商品ページ HTML から補完します
- 検索結果は BOOTH の検索ページを解析して取得します

詳細な API は [docs/api-reference.md](/home/trkdknt/git/Cozym/booth-go/docs/api-reference.md) を参照してください。

## Requirements

- Go `1.26.1`

## Install

```bash
go get github.com/trkd-knt/booth-go
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	booth "github.com/trkd-knt/booth-go"
)

func main() {
	client, err := booth.NewClient(
		booth.WithLang("ja"),
		booth.WithRateLimit(1),
	)
	if err != nil {
		log.Fatal(err)
	}

	item, err := client.GetItem(context.Background(), "7472126")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(item.Title)
	fmt.Println(item.PriceText)
	fmt.Println(item.Tags)
	fmt.Println(item.Avatars)
}
```

検索例:

```go
result, err := client.SearchItems(context.Background(), booth.SearchOptions{
	Category:    "3D衣装",
	Query:       "VRChat",
	ExceptWords: []string{"R18"},
	Tags:        []string{"アバター"},
	Type:        booth.ItemTypeDigital,
	Adult:       booth.AdultFilterDefault,
	Sort:        booth.SortPopular,
	Page:        1,
})
if err != nil {
	log.Fatal(err)
}

for _, item := range result.Items {
	fmt.Printf("%s %d %s\n", item.Title, item.Price, item.URL)
}
```

## Main Methods

- `GetItem(ctx, itemID)`  
  商品 JSON を基準に商品詳細を返します。`Tags`、`Images`、`ImageDetails`、`Summary`、`Shop` などを含みます。`Avatars` は HTML から補完します。

- `GetItemDescription(ctx, itemID)`  
  商品ページ HTML から全文説明を返します。JSON の `Summary` より長い本文が必要な場合に使います。

- `SearchItems(ctx, opts)`  
  BOOTH 全体検索結果を返します。

- `GetShop(ctx, shopHost)`  
  ショップ情報を返します。

## Notes

- `Summary` は商品 JSON の説明文です
- `Description` は `GetItem` では埋まりません。全文が必要な場合は `GetItemDescription` を使ってください
- 検索結果の `Items` では `ShopHost` は保証しません。ショップ情報が拾えた場合だけ `Shop` が入ります
- `Avatars` は商品ページ HTML に存在する場合のみ返ります
