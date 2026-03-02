# Client SDK

The `client` package provides a high-level SDK that simplifies working with multiple search engines.

## Key Features

- **Auto-registration**: Automatically discovers and registers all available engines
- **Smart selection**: Uses `SEARCH_ENGINE` environment variable or defaults to Serper
- **Runtime switching**: Switch between engines without recreating the client
- **Capability checking**: Validates operations before calling backends
- **Error handling**: Returns `ErrOperationNotSupported` for unsupported operations
- **Clean API**: Implements the same `Engine` interface, proxying to the selected backend

## Quick Start

```go
import "github.com/plexusone/omniserp/client"

// Create client - auto-selects engine based on SEARCH_ENGINE env var
c, err := client.New()

// Or specify engine explicitly
c, err := client.NewWithEngine("serper")

// Check support before calling
if c.SupportsOperation(client.OpSearchLens) {
    result, _ := c.SearchLens(ctx, params)
}

// Switch engines at runtime
c.SetEngine("serpapi")
```

## Operation Constants

The SDK provides constants for all operations:

| Constant | Operation |
|----------|-----------|
| `client.OpSearch` | Web search |
| `client.OpSearchNews` | News search |
| `client.OpSearchImages` | Image search |
| `client.OpSearchVideos` | Video search |
| `client.OpSearchPlaces` | Places search |
| `client.OpSearchMaps` | Maps search |
| `client.OpSearchReviews` | Reviews search |
| `client.OpSearchShopping` | Shopping search |
| `client.OpSearchScholar` | Scholar search |
| `client.OpSearchLens` | Lens search (Serper only) |
| `client.OpSearchAutocomplete` | Autocomplete |
| `client.OpScrapeWebpage` | Webpage scraping |

## Capability Checking

The client SDK automatically checks if operations are supported by the current backend:

```go
c, _ := client.New()

// Check if an operation is supported
if c.SupportsOperation(client.OpSearchLens) {
    result, err := c.SearchLens(ctx, params)
    // ...
} else {
    log.Println("Current engine doesn't support Lens search")
}

// Or let the client return an error
result, err := c.SearchLens(ctx, params)
if errors.Is(err, client.ErrOperationNotSupported) {
    log.Printf("Operation not supported: %v", err)
}
```

## Engine Switching

```go
// Create client with default engine
c, _ := client.New()
fmt.Printf("Using: %s\n", c.GetName())  // "serper"

// Switch to a different engine
c.SetEngine("serpapi")
fmt.Printf("Now using: %s\n", c.GetName())  // "serpapi"

// List all available engines
engines := c.ListEngines()
fmt.Printf("Available: %v\n", engines)  // ["serper", "serpapi"]
```

## Advanced: Registry Access

For direct registry access:

```go
import (
    "github.com/plexusone/omniserp"
    "github.com/plexusone/omniserp/client/serper"
    "github.com/plexusone/omniserp/client/serpapi"
)

func main() {
    // Create registry and manually register engines
    registry := omniserp.NewRegistry()

    // Register engines (handle errors as needed)
    if serperEngine, err := serper.New(); err == nil {
        registry.Register(serperEngine)
    }
    if serpApiEngine, err := serpapi.New(); err == nil {
        registry.Register(serpApiEngine)
    }

    // Get default engine (based on SEARCH_ENGINE env var)
    engine, err := omniserp.GetDefaultEngine(registry)
    if err != nil {
        log.Printf("Warning: %v", err)
    }

    // Perform a search
    result, err := engine.Search(context.Background(), omniserp.SearchParams{
        Query: "golang programming",
    })
}
```

## Registry Operations

```go
// Create new registry and register engines
registry := omniserp.NewRegistry()

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

## Engine Information

```go
// Get info about specific engine
engine, _ := registry.Get("serper")
info := omniserp.GetEngineInfo(engine)
log.Printf("Engine: %s v%s, Tools: %v", info.Name, info.Version, info.SupportedTools)

// Get info about all engines
allInfo := omniserp.GetAllEngineInfo(registry)
```

## Error Handling

```go
engine, err := omniserp.GetDefaultEngine(registry)
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
