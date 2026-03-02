# Getting Started

## Installation

### Go Package

```bash
go get github.com/plexusone/omniserp@latest
```

### CLI Tool

```bash
go install github.com/plexusone/omniserp/cmd/omniserp@latest
```

### MCP Server

```bash
go install github.com/plexusone/omniserp/cmd/mcp-omniserp@latest
```

## Configuration

### Environment Variables

OmniSerp uses environment variables for configuration:

```bash
# Choose which engine to use (optional, defaults to "serper")
export SEARCH_ENGINE="serper"  # or "serpapi"

# API keys for respective engines
export SERPER_API_KEY="your_serper_key"
export SERPAPI_API_KEY="your_serpapi_key"
```

### Getting API Keys

#### Serper.dev

1. Visit [serper.dev](https://serper.dev)
2. Sign up for an account
3. Copy your API key from the dashboard
4. Set the environment variable: `export SERPER_API_KEY="your_key"`

#### SerpAPI

1. Visit [serpapi.com](https://serpapi.com)
2. Sign up for an account
3. Copy your API key from your account settings
4. Set the environment variable: `export SERPAPI_API_KEY="your_key"`

## Basic Usage

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
    // Create client (auto-selects engine based on SEARCH_ENGINE env var)
    c, err := client.New()
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    log.Printf("Using engine: %s v%s", c.GetName(), c.GetVersion())

    // Perform a search
    result, err := c.Search(context.Background(), omniserp.SearchParams{
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

## Selecting a Specific Engine

```go
// Create client with a specific engine
c, err := client.NewWithEngine("serpapi")
if err != nil {
    log.Fatal(err)
}

// Or switch engines at runtime
c.SetEngine("serper")
```

## Testing

Run tests without API keys (tests will skip gracefully):

```bash
go test ./...
```

Run tests with API calls (requires API keys):

```bash
export SERPER_API_KEY="your_key"
export SERPAPI_API_KEY="your_key"
go test -v ./client
```
