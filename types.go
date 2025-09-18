package metasearch

import (
	"context"
)

// SearchParams represents common search parameters across all engines
type SearchParams struct {
	Query      string `json:"query" jsonschema:"description:Search query"`
	Location   string `json:"location,omitempty" jsonschema:"description:Search location"`
	Language   string `json:"language,omitempty" jsonschema:"description:Search language (e.g., 'en')"`
	Country    string `json:"country,omitempty" jsonschema:"description:Country code (e.g., 'us')"`
	NumResults int    `json:"num_results,omitempty" jsonschema:"description:Number of results (1-100),default:10"`
}

// ScrapeParams represents parameters for web scraping
type ScrapeParams struct {
	URL string `json:"url" jsonschema:"description:URL to scrape"`
}

// SearchResult represents a common search result structure
type SearchResult struct {
	Data interface{} `json:"data"`
	Raw  string      `json:"raw,omitempty"`
}

// Engine defines the interface that all search engines must implement
type Engine interface {
	// GetName returns the name of the search engine
	GetName() string

	// GetVersion returns the version of the engine implementation
	GetVersion() string

	// GetSupportedTools returns a list of tool names supported by this engine
	GetSupportedTools() []string

	// Search performs a general web search
	Search(ctx context.Context, params SearchParams) (*SearchResult, error)

	// SearchNews performs a news search
	SearchNews(ctx context.Context, params SearchParams) (*SearchResult, error)

	// SearchImages performs an image search
	SearchImages(ctx context.Context, params SearchParams) (*SearchResult, error)

	// SearchVideos performs a video search
	SearchVideos(ctx context.Context, params SearchParams) (*SearchResult, error)

	// SearchPlaces performs a places search
	SearchPlaces(ctx context.Context, params SearchParams) (*SearchResult, error)

	// SearchMaps performs a maps search
	SearchMaps(ctx context.Context, params SearchParams) (*SearchResult, error)

	// SearchReviews performs a reviews search
	SearchReviews(ctx context.Context, params SearchParams) (*SearchResult, error)

	// SearchShopping performs a shopping search
	SearchShopping(ctx context.Context, params SearchParams) (*SearchResult, error)

	// SearchScholar performs a scholar search
	SearchScholar(ctx context.Context, params SearchParams) (*SearchResult, error)

	// SearchLens performs a visual search (if supported)
	SearchLens(ctx context.Context, params SearchParams) (*SearchResult, error)

	// SearchAutocomplete gets search suggestions
	SearchAutocomplete(ctx context.Context, params SearchParams) (*SearchResult, error)

	// ScrapeWebpage scrapes content from a webpage
	ScrapeWebpage(ctx context.Context, params ScrapeParams) (*SearchResult, error)
}

// Registry manages available search engines
type Registry struct {
	engines map[string]Engine
}

// NewRegistry creates a new engine registry
func NewRegistry() *Registry {
	return &Registry{
		engines: make(map[string]Engine),
	}
}

// Register adds a search engine to the registry
func (r *Registry) Register(engine Engine) {
	r.engines[engine.GetName()] = engine
}

// Get retrieves a search engine by name
func (r *Registry) Get(name string) (Engine, bool) {
	engine, exists := r.engines[name]
	return engine, exists
}

// List returns all registered engine names
func (r *Registry) List() []string {
	names := make([]string, 0, len(r.engines))
	for name := range r.engines {
		names = append(names, name)
	}
	return names
}

// GetAll returns all registered engines
func (r *Registry) GetAll() map[string]Engine {
	result := make(map[string]Engine)
	for name, engine := range r.engines {
		result[name] = engine
	}
	return result
}
