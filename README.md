# Metasearch Package

A modular, plugin-based search engine abstraction package for Go that provides a unified interface for multiple search engines.

## Overview

The `metasearch` package provides:

- **Unified Interface**: Common `Engine` interface for all search providers
- **Plugin Architecture**: Easy addition of new search engines
- **Multiple Providers**: Built-in support for Serper and SerpAPI
- **Type Safety**: Structured parameter and result types
- **Registry System**: Automatic discovery and management of engines

## CLI Usage

### Installation
```bash
go build ./cmd/metasearch
```

### Basic Usage
```bash
# Set API key
export SERPER_API_KEY="your_api_key"

# Basic search
./metasearch "golang programming"

# Specify engine
./metasearch "golang programming" serper
./metasearch "golang programming" serpapi

# Use environment variable for default engine
export SEARCH_ENGINE="serpapi"
./metasearch "golang programming"
```

## Library Usage

```go
package main

import (
    "context"
    "log"
    
    "github.com/grokify/metasearch"
    "github.com/grokify/metasearch/serper"
    "github.com/grokify/metasearch/serpapi"
)

func main() {
    // Create registry and manually register engines
    registry := metasearch.NewRegistry()
    
    // Register engines (handle errors as needed)
    if serperEngine, err := serper.New(); err == nil {
        registry.Register(serperEngine)
    }
    if serpApiEngine, err := serpapi.New(); err == nil {
        registry.Register(serpApiEngine)
    }
    
    // Get default engine (based on SEARCH_ENGINE env var)
    engine, err := metasearch.GetDefaultEngine(registry)
    if err != nil {
        log.Printf("Warning: %v", err) // May still have a fallback engine
    }
    
    if engine == nil {
        log.Fatal("No search engines available")
    }
    
    // Perform a search
    result, err := engine.Search(context.Background(), metasearch.SearchParams{
        Query:      "golang programming",
        NumResults: 10,
        Language:   "en",
        Country:    "us",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Use the result
    log.Printf("Search completed: %+v", result.Data)
}
```

## Supported Engines

### Serper
- **Package**: `metasearch/serper`
- **Environment Variable**: `SERPER_API_KEY`
- **Website**: [serper.dev](https://serper.dev)
- **All search types supported**

### SerpAPI
- **Package**: `metasearch/serpapi`
- **Environment Variable**: `SERPAPI_API_KEY`
- **Website**: [serpapi.com](https://serpapi.com)
- **Most search types supported**

## Available Search Methods

All engines implement these methods:

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

## Types

### SearchParams
```go
type SearchParams struct {
    Query      string `json:"query"`                    // Required: search query
    Location   string `json:"location,omitempty"`       // Optional: search location
    Language   string `json:"language,omitempty"`       // Optional: language code (e.g., "en")
    Country    string `json:"country,omitempty"`        // Optional: country code (e.g., "us")
    NumResults int    `json:"num_results,omitempty"`    // Optional: number of results (1-100)
}
```

### ScrapeParams
```go
type ScrapeParams struct {
    URL string `json:"url"` // Required: URL to scrape
}
```

### SearchResult
```go
type SearchResult struct {
    Data interface{} `json:"data"`          // Parsed response data
    Raw  string      `json:"raw,omitempty"` // Raw response (optional)
}
```

## Registry Usage

### Basic Registry Operations
```go
// Create new registry and register engines
registry := metasearch.NewRegistry()

// Register engines manually
if serperEngine, err := serper.New(); err == nil {
    registry.Register(serperEngine)
}

// List available engines
engines := registry.List()
log.Printf("Available engines: %v", engines)

// Get specific engine
if engine, exists := registry.Get("serper"); exists {
    log.Printf("Using engine: %s v%s", engine.GetName(), engine.GetVersion())
}

// Get all engines
allEngines := registry.GetAll()
```

### Engine Information
```go
// Get info about specific engine
engine, _ := registry.Get("serper")
info := metasearch.GetEngineInfo(engine)
log.Printf("Engine: %s v%s, Tools: %v", info.Name, info.Version, info.SupportedTools)

// Get info about all engines
allInfo := metasearch.GetAllEngineInfo(registry)
```

## Environment Configuration

The package uses environment variables for configuration:

```bash
# Choose which engine to use (optional, defaults to "serper")
export SEARCH_ENGINE="serper"  # or "serpapi"

# API keys for respective engines
export SERPER_API_KEY="your_serper_key"
export SERPAPI_API_KEY="your_serpapi_key"
```

## Adding New Engines

To add a new search engine:

1. **Create engine package**:
```go
// metasearch/newengine/newengine.go
package newengine

import (
    "context"
    "fmt"
    "os"
    "github.com/grokify/metasearch"
)

type Engine struct {
    apiKey string
    // other fields
}

func New() (*Engine, error) {
    apiKey := os.Getenv("NEWENGINE_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("NEWENGINE_API_KEY required")
    }
    return &Engine{apiKey: apiKey}, nil
}

func (e *Engine) GetName() string { return "newengine" }
func (e *Engine) GetVersion() string { return "1.0.0" }
func (e *Engine) GetSupportedTools() []string { /* return supported tools */ }

// Implement all other metasearch.Engine methods...
func (e *Engine) Search(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
    // Implementation
}
// ... implement all other interface methods
```

2. **Register in your application**:
```go
// In your application code (e.g., cmd/yourapp/main.go)
import "github.com/grokify/metasearch/newengine"

func createRegistry() *metasearch.Registry {
    registry := metasearch.NewRegistry()
    
    // Register existing engines
    if serperEngine, err := serper.New(); err == nil {
        registry.Register(serperEngine)
    }
    
    // Register new engine
    if newEng, err := newengine.New(); err == nil {
        registry.Register(newEng)
    }
    
    return registry
}
```

3. **Update CLI (optional)**:
   Add the new engine import and registration to `cmd/metasearch/main.go`

## Error Handling

The package provides consistent error handling:

```go
engine, err := metasearch.GetDefaultEngine(registry)
if err != nil {
    // Handle engine selection error
    log.Printf("Engine selection warning: %v", err)
}

result, err := engine.Search(ctx, params)
if err != nil {
    // Handle search error
    log.Printf("Search failed: %v", err)
}
```

## Thread Safety

The registry is safe for concurrent read operations. Engine implementations should be thread-safe for concurrent use.

## License

This package is designed to be self-contained and can be used independently of the MCP server implementation.