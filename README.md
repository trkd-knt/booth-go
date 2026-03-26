# booth-go

`booth-go` は BOOTH の商品詳細、ショップ情報、全体検索結果を取得するための Go ライブラリです。  
スクレイピング対象の HTML から、商品情報や検索結果を Go の構造体として扱える形に変換します。

詳細な API は [docs/api-reference.md](/home/trkdknt/git/Cozym/booth-go/docs/api-reference.md) を参照してください。

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

	result, err := client.SearchItems(context.Background(), booth.SearchOptions{
		Category:    "3D衣装",
		Query:       "VRoid",
		ExceptWords: []string{"R18"},
		Tags:        []string{"アバター"},
		Event:       "osakafes-mar2026",
		Type:        booth.ItemTypeDigital,
		Adult:       booth.AdultFilterOnly,
		MinPrice:    1000,
		MaxPrice:    2500,
		Sort:        booth.SortPopular,
		Page:        1,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range result.Items {
		fmt.Printf("%s (%s) %d\n", item.Title, item.ShopHost, item.Price)
	}
}
```

## Response Example

`GetItem` のレスポンス例:

```go
&booth.Item{
	ID:        "3652121",
	Title:     "[無料/Free] VRoid Hair Texture",
	Price:     0,
	PriceText: "0 JPY",
	ShopHost:  "honeyrosy.booth.pm",
	Images: []string{
		"https://booth.pximg.net/example-1.jpg",
		"https://booth.pximg.net/example-2.jpg",
	},
	ImageDetails: []booth.Image{
		{
			Original: "https://booth.pximg.net/example-1.jpg",
			Resized:  "https://booth.pximg.net/c/72x72_a2_g5/example-1.jpg",
		},
	},
	Description: "Pastel Balayage Color Hair Textures...",
	URL:         "https://honeyrosy.booth.pm/items/3652121",
	IsSoldOut:   false,
	Category: &booth.Category{
		ID:   212,
		Name: "VRoid",
	},
	Shop: &booth.ShopPreview{
		Name:      "HoneyRosy",
		Host:      "honeyrosy.booth.pm",
		URL:       "https://honeyrosy.booth.pm",
		Thumbnail: "https://booth.pximg.net/icon.jpg",
	},
	IsAdult: false,
	Likes:   10415,
	Downloadables: []booth.Downloadable{
		{
			FileName:      "Pastel_Balayage_Color_Hair_Texture",
			FileExtension: ".zip",
			FileSize:      "3.92 MB",
			Name:          "Pastel_Balayage_Color_Hair_Texture.zip",
			URL:           "https://booth.pm/downloadables/2336025",
		},
	},
}
```
