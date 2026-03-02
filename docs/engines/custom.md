# Adding Custom Engines

OmniSerp's plugin architecture makes it easy to add new search engine backends.

## Implementation Steps

### 1. Create Engine Package

Create a new package under `client/`:

```go
// client/newengine/newengine.go
package newengine

import (
    "context"
    "fmt"
    "os"
    "github.com/plexusone/omniserp"
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
func (e *Engine) GetSupportedTools() []string {
    return []string{
        "google_search",
        "google_search_news",
        // ... list supported operations
    }
}
```

### 2. Implement Search Methods

Implement all methods from the `Engine` interface:

```go
func (e *Engine) Search(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
    // Make API request to your search engine
    // Parse response
    // Return normalized result

    return &omniserp.SearchResult{
        Data: result,
        Raw:  rawResponse,
    }, nil
}

func (e *Engine) SearchNews(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
    // Implementation
}

// ... implement all other interface methods

// For unsupported operations, return an error:
func (e *Engine) SearchLens(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
    return nil, fmt.Errorf("google_search_lens is not supported by newengine")
}
```

### 3. Register in Your Application

```go
// In your application code (e.g., cmd/yourapp/main.go)
import (
    "github.com/plexusone/omniserp"
    "github.com/plexusone/omniserp/client/newengine"
    "github.com/plexusone/omniserp/client/serper"
)

func createRegistry() *omniserp.Registry {
    registry := omniserp.NewRegistry()

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

### 4. Update CLI (Optional)

Add the new engine import and registration to `cmd/omniserp/main.go`.

## Full Example

Here's a complete example for a hypothetical "CustomSearch" engine:

```go
package customsearch

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "github.com/plexusone/omniserp"
)

const (
    baseURL       = "https://api.customsearch.com"
    engineName    = "customsearch"
    engineVersion = "1.0.0"
)

type Engine struct {
    apiKey string
    client *http.Client
}

func New() (*Engine, error) {
    apiKey := os.Getenv("CUSTOMSEARCH_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("CUSTOMSEARCH_API_KEY environment variable is required")
    }
    return &Engine{
        apiKey: apiKey,
        client: &http.Client{},
    }, nil
}

func (e *Engine) GetName() string    { return engineName }
func (e *Engine) GetVersion() string { return engineVersion }

func (e *Engine) GetSupportedTools() []string {
    return []string{
        "google_search",
        "google_search_news",
        "google_search_images",
    }
}

func (e *Engine) Search(ctx context.Context, params omniserp.SearchParams) (*omniserp.SearchResult, error) {
    // Build request
    url := fmt.Sprintf("%s/search?q=%s&key=%s", baseURL, params.Query, e.apiKey)

    resp, err := e.client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &omniserp.SearchResult{
        Data: result,
    }, nil
}

// Implement other methods...
```

## Best Practices

1. **Environment Variables**: Use environment variables for API keys
2. **Error Messages**: Provide clear error messages for missing configuration
3. **Supported Tools**: Only list tools that are actually implemented
4. **Graceful Failures**: Return descriptive errors for unsupported operations
5. **Thread Safety**: Ensure your engine is safe for concurrent use
