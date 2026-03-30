---
name: booth-go-sdk
description: Use when a task needs to consume the booth-go Go library from another project, including writing sample code, building BOOTH search queries with SearchOptions, calling GetItem or GetShop, handling results, or explaining the public API.
---

# booth-go SDK

Use this skill when the task is about using the `booth-go` library as a dependency in application code.

## Use Cases

- Write Go code that calls `booth.NewClient`
- Search BOOTH items with `SearchItems`
- Fetch item details with `GetItem`
- Fetch shop information with `GetShop`
- Explain `SearchOptions`, `Sort`, `ItemType`, and `AdultFilter`
- Show how to handle `Item`, `Shop`, and `SearchResult`

## Public API Guidance

Create a client with:

```go
client, err := booth.NewClient(
    booth.WithLang("ja"),
)
```

Primary methods:

- `GetItem(ctx, itemID)`
- `GetItemDescription(ctx, itemID)`
- `SearchItems(ctx, opts)`
- `GetShop(ctx, shopHost)`

## SearchOptions Mapping

Map user intent to `SearchOptions` like this:

- keyword search -> `Query`
- browse category slug/path segment -> `Category`
- excluded words -> `ExceptWords`
- tags -> `Tags`
- event slug -> `Event`
- product kind -> `Type`
- adult filter -> `Adult`
- minimum price -> `MinPrice`
- maximum price -> `MaxPrice`
- sort order -> `Sort`
- page number -> `Page`

Use typed values where available:

- `SortDefault`
- `SortNewest`
- `SortPopular`
- `SortPriceAsc`
- `SortPriceDesc`
- `ItemTypeDigital`
- `ItemTypePhysical`
- `AdultFilterDefault`
- `AdultFilterOnly`
- `AdultFilterInclude`

## Response Handling

Common result shapes:

- `SearchItems` returns `*SearchResult`
- `GetItem` returns `*Item`
- `GetShop` returns `*Shop`

Useful fields on `Item` include:

- `Title`
- `Price`
- `PriceText`
- `ShopHost`
- `Images`
- `Category`
- `Shop`
- `IsAdult`
- `Likes`
- `Downloadables`

## Output Style

When helping a caller use this SDK:

- Prefer concise example code over long prose
- Use the public API only
- Avoid referencing internals such as `internal/parser`
- If the task is about filtering/search, show the exact `SearchOptions` struct to build
- If error handling matters, mention exported errors such as `ErrItemNotFound`, `ErrShopNotFound`, `ErrTooManyRequests`, `ErrParseFailed`, and `ErrRequestFailed`
