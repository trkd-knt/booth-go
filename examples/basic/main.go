package main

import (
	"context"
	"fmt"
	"log"

	booth "github.com/trkd-knt/booth-go"
)

// main は booth-go の基本的な利用例を示します。
func main() {
	client, err := booth.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.SearchItems(context.Background(), booth.SearchOptions{
		Query: "sample",
		Page:  1,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range result.Items {
		fmt.Printf("%s (%s) %d\n", item.Title, item.ShopHost, item.Price)
	}
}
