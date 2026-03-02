# Normalized Responses

The client SDK provides **optional normalized response methods** that return unified structures across all search engines.

## Overview

Different search engines return results in different formats. Normalized responses provide a consistent structure regardless of which backend engine is used.

## Usage

```go
// Use *Normalized() methods for engine-agnostic response structures
normalized, err := c.SearchNormalized(ctx, params)

// Access results in a consistent format regardless of engine
for _, result := range normalized.OrganicResults {
    fmt.Printf("%s: %s\n", result.Title, result.Link)
}

// Switch engines without changing your code!
c.SetEngine("serpapi")
normalized, err = c.SearchNormalized(ctx, params) // Same structure!
```

## Available Normalized Methods

| Method | Description |
|--------|-------------|
| `SearchNormalized()` | Web search with normalized results |
| `SearchNewsNormalized()` | News search with normalized results |
| `SearchImagesNormalized()` | Image search with normalized results |

## Normalized Structure

```go
type NormalizedSearchResult struct {
    OrganicResults  []OrganicResult    // Standard search results
    AnswerBox       *AnswerBox         // Featured answer
    KnowledgeGraph  *KnowledgeGraph    // Knowledge panel
    RelatedSearches []RelatedSearch    // Related queries
    PeopleAlsoAsk   []PeopleAlsoAsk    // PAA questions
    NewsResults     []NewsResult       // News articles
    ImageResults    []ImageResult      // Images
    SearchMetadata  SearchMetadata     // Search info
    Raw             *SearchResult      // Original response
}
```

### OrganicResult

```go
type OrganicResult struct {
    Title    string
    Link     string
    Snippet  string
    Position int
}
```

### AnswerBox

```go
type AnswerBox struct {
    Title  string
    Answer string
    Link   string
}
```

### KnowledgeGraph

```go
type KnowledgeGraph struct {
    Title       string
    Type        string
    Description string
}
```

## Comparison: Raw vs Normalized

| Aspect | Raw Response | Normalized Response |
|--------|--------------|---------------------|
| Field names | Engine-specific | Unified |
| Structure | Varies by engine | Consistent |
| Engine switching | Requires code changes | No changes needed |
| Type safety | `interface{}` | Strongly typed |
| Use case | Engine-specific features | Engine-agnostic apps |

## Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/plexusone/omniserp"
    "github.com/plexusone/omniserp/client"
)

func main() {
    c, err := client.New()
    if err != nil {
        log.Fatal(err)
    }

    // Get normalized response
    normalized, err := c.SearchNormalized(context.Background(), omniserp.SearchParams{
        Query:      "golang programming",
        NumResults: 5,
    })
    if err != nil {
        log.Fatal(err)
    }

    // Print metadata
    fmt.Printf("Engine: %s\n", normalized.SearchMetadata.Engine)
    fmt.Printf("Query: %s\n", normalized.SearchMetadata.Query)
    fmt.Printf("Total Results: %d\n\n", len(normalized.OrganicResults))

    // Print organic results (same structure regardless of engine)
    for i, result := range normalized.OrganicResults {
        fmt.Printf("%d. %s\n", i+1, result.Title)
        fmt.Printf("   URL: %s\n", result.Link)
        fmt.Printf("   Snippet: %s\n\n", result.Snippet)
    }

    // Print answer box if present
    if normalized.AnswerBox != nil {
        fmt.Println("=== Answer Box ===")
        fmt.Printf("Title: %s\n", normalized.AnswerBox.Title)
        fmt.Printf("Answer: %s\n", normalized.AnswerBox.Answer)
    }

    // Print knowledge graph if present
    if normalized.KnowledgeGraph != nil {
        fmt.Println("=== Knowledge Graph ===")
        fmt.Printf("Title: %s\n", normalized.KnowledgeGraph.Title)
        fmt.Printf("Type: %s\n", normalized.KnowledgeGraph.Type)
    }
}
```

## Benefits

- **Engine-Agnostic**: Same code works with any backend
- **Type-Safe**: Strongly-typed result structures
- **Optional**: Raw responses still available via standard methods
- **Complete**: Preserves original response in `Raw` field
