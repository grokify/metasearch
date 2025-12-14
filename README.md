# Metasearch Package

[![Build Status][build-status-svg]][build-status-url]
[![Lint Status][lint-status-svg]][lint-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

A modular, plugin-based search engine abstraction package for Go that provides a unified interface for multiple search engines.

## Overview

The `metasearch` package provides:

- **Unified Interface**: Common `Engine` interface for all search providers
- **Plugin Architecture**: Easy addition of new search engines
- **Multiple Providers**: Built-in support for Serper and SerpAPI
- **Type Safety**: Structured parameter and result types
- **Registry System**: Automatic discovery and management of engines
- **MCP Server**: Model Context Protocol server for AI integration (`cmd/mcpserver`)

## Applications

### CLI Tool

#### Installation
```bash
go build ./cmd/metasearch
```

#### Basic Usage
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

### MCP Server

The Model Context Protocol (MCP) server enables AI assistants to perform web searches through this package.

#### Installation
```bash
go install github.com/grokify/metasearch/cmd/mcpserver@latest
```

Or build from source:
```bash
go build ./cmd/mcpserver
```

#### Configuration

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "metasearch": {
      "command": "mcpserver",
      "env": {
        "SERPER_API_KEY": "your_serper_api_key",
        "SEARCH_ENGINE": "serper"
      }
    }
  }
}
```

#### Features

The MCP server provides tools for:
- **Web Search**: General web searches with customizable parameters
- **News Search**: Search news articles
- **Image Search**: Search for images
- **Video Search**: Search for videos
- **Places Search**: Search for locations and businesses
- **Maps Search**: Search maps data
- **Reviews Search**: Search reviews
- **Shopping Search**: Search shopping/product listings
- **Scholar Search**: Search academic papers
- **Lens Search**: Visual search capabilities
- **Autocomplete**: Get search suggestions

All searches support parameters like location, language, country, and number of results.

## Library Usage

```go
package main

import (
    "context"
    "log"

    "github.com/grokify/metasearch"
    "github.com/grokify/metasearch/client/serper"
    "github.com/grokify/metasearch/client/serpapi"
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
- **Package**: `github.com/grokify/metasearch/client/serper`
- **Environment Variable**: `SERPER_API_KEY`
- **Website**: [serper.dev](https://serper.dev)
- **All search types supported**

### SerpAPI
- **Package**: `github.com/grokify/metasearch/client/serpapi`
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
import (
    "github.com/grokify/metasearch"
    "github.com/grokify/metasearch/client/serper"
)

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

1. **Create engine package under `client/`**:
```go
// client/newengine/newengine.go
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
import (
    "github.com/grokify/metasearch"
    "github.com/grokify/metasearch/client/newengine"
    "github.com/grokify/metasearch/client/serper"
)

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

 [build-status-svg]: https://github.com/grokify/metasearch/actions/workflows/ci.yaml/badge.svg?branch=main
 [build-status-url]: https://github.com/grokify/metasearch/actions/workflows/ci.yaml
 [lint-status-svg]: https://github.com/grokify/metasearch/actions/workflows/lint.yaml/badge.svg?branch=main
 [lint-status-url]: https://github.com/grokify/metasearch/actions/workflows/lint.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/metasearch
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/metasearch
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/metasearch
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/metasearch
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fmetasearch
 [loc-svg]: https://tokei.rs/b1/github/grokify/metasearch
 [repo-url]: https://github.com/grokify/metasearch
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/metasearch/blob/master/LICENSE
