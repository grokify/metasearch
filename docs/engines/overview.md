# Supported Engines

OmniSerp supports multiple search engine backends through a unified interface.

## Available Engines

### Serper

- **Package**: `github.com/plexusone/omniserp/client/serper`
- **Environment Variable**: `SERPER_API_KEY`
- **Website**: [serper.dev](https://serper.dev)
- **Supported Operations**: All 12 search types including Lens

### SerpAPI

- **Package**: `github.com/plexusone/omniserp/client/serpapi`
- **Environment Variable**: `SERPAPI_API_KEY`
- **Website**: [serpapi.com](https://serpapi.com)
- **Supported Operations**: 11 search types (no Lens support)

!!! note
    `SearchLens()` is not supported by SerpAPI and will return `ErrOperationNotSupported`

## Feature Comparison

| Operation | Serper | SerpAPI |
|-----------|:------:|:-------:|
| Web Search | ✓ | ✓ |
| News Search | ✓ | ✓ |
| Image Search | ✓ | ✓ |
| Video Search | ✓ | ✓ |
| Places Search | ✓ | ✓ |
| Maps Search | ✓ | ✓ |
| Reviews Search | ✓ | ✓ |
| Shopping Search | ✓ | ✓ |
| Scholar Search | ✓ | ✓ |
| **Lens Search** | **✓** | **✗** |
| Autocomplete | ✓ | ✓ |
| Webpage Scrape | ✓ | ✓ |

## Engine Interface

All engines implement the `Engine` interface:

```go
type Engine interface {
    // Metadata
    GetName() string
    GetVersion() string
    GetSupportedTools() []string

    // Search methods
    Search(ctx context.Context, params SearchParams) (*SearchResult, error)
    SearchNews(ctx context.Context, params SearchParams) (*SearchResult, error)
    SearchImages(ctx context.Context, params SearchParams) (*SearchResult, error)
    SearchVideos(ctx context.Context, params SearchParams) (*SearchResult, error)
    SearchPlaces(ctx context.Context, params SearchParams) (*SearchResult, error)
    SearchMaps(ctx context.Context, params SearchParams) (*SearchResult, error)
    SearchReviews(ctx context.Context, params SearchParams) (*SearchResult, error)
    SearchShopping(ctx context.Context, params SearchParams) (*SearchResult, error)
    SearchScholar(ctx context.Context, params SearchParams) (*SearchResult, error)
    SearchLens(ctx context.Context, params SearchParams) (*SearchResult, error)
    SearchAutocomplete(ctx context.Context, params SearchParams) (*SearchResult, error)

    // Utility
    ScrapeWebpage(ctx context.Context, params ScrapeParams) (*SearchResult, error)
}
```

## Selecting an Engine

### Via Environment Variable

```bash
export SEARCH_ENGINE="serper"  # or "serpapi"
```

### Programmatically

```go
// Use default (from SEARCH_ENGINE env var)
c, err := client.New()

// Specify explicitly
c, err := client.NewWithEngine("serpapi")

// Switch at runtime
c.SetEngine("serper")
```

## Checking Engine Capabilities

```go
c, _ := client.New()

// Check current engine
fmt.Printf("Engine: %s v%s\n", c.GetName(), c.GetVersion())

// List supported tools
tools := c.GetSupportedTools()
fmt.Printf("Supported: %v\n", tools)

// Check specific operation
if c.SupportsOperation(client.OpSearchLens) {
    // Lens is supported
}
```
