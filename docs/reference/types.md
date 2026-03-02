# Types Reference

## Core Types

### SearchParams

Parameters for search operations.

```go
type SearchParams struct {
    Query      string `json:"query"`                 // Required: search query
    Location   string `json:"location,omitempty"`    // Optional: search location
    Language   string `json:"language,omitempty"`    // Optional: language code (e.g., "en")
    Country    string `json:"country,omitempty"`     // Optional: country code (e.g., "us")
    NumResults int    `json:"num_results,omitempty"` // Optional: number of results (1-100)
}
```

#### Fields

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `Query` | `string` | The search query (required) | `"golang programming"` |
| `Location` | `string` | Geographic location for localized results | `"New York, NY"` |
| `Language` | `string` | Language code (ISO 639-1) | `"en"`, `"es"`, `"fr"` |
| `Country` | `string` | Country code (ISO 3166-1 alpha-2) | `"us"`, `"gb"`, `"de"` |
| `NumResults` | `int` | Number of results to return (1-100) | `10` |

### ScrapeParams

Parameters for webpage scraping.

```go
type ScrapeParams struct {
    URL string `json:"url"` // Required: URL to scrape
}
```

### SearchResult

Result returned by all search operations.

```go
type SearchResult struct {
    Data interface{} `json:"data"`          // Parsed response data
    Raw  string      `json:"raw,omitempty"` // Raw response (optional)
}
```

## Normalized Types

### NormalizedSearchResult

Unified response structure across all engines.

```go
type NormalizedSearchResult struct {
    OrganicResults  []OrganicResult
    AnswerBox       *AnswerBox
    KnowledgeGraph  *KnowledgeGraph
    RelatedSearches []RelatedSearch
    PeopleAlsoAsk   []PeopleAlsoAsk
    NewsResults     []NewsResult
    ImageResults    []ImageResult
    SearchMetadata  SearchMetadata
    Raw             *SearchResult
}
```

### OrganicResult

A single organic search result.

```go
type OrganicResult struct {
    Title    string
    Link     string
    Snippet  string
    Position int
}
```

### AnswerBox

Featured answer snippet.

```go
type AnswerBox struct {
    Title  string
    Answer string
    Link   string
}
```

### KnowledgeGraph

Knowledge panel information.

```go
type KnowledgeGraph struct {
    Title       string
    Type        string
    Description string
}
```

### SearchMetadata

Metadata about the search request.

```go
type SearchMetadata struct {
    Engine string
    Query  string
}
```

## Error Types

### ErrOperationNotSupported

Returned when an operation is not supported by the current engine.

```go
var ErrOperationNotSupported = errors.New("operation not supported by current engine")
```

**Usage:**

```go
result, err := c.SearchLens(ctx, params)
if errors.Is(err, client.ErrOperationNotSupported) {
    log.Printf("Lens search not supported: %v", err)
}
```

## Engine Interface

The interface that all search engines must implement.

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
