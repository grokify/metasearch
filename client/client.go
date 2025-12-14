package client

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/grokify/metasearch"
	"github.com/grokify/metasearch/client/serpapi"
	"github.com/grokify/metasearch/client/serper"
)

// Operation names that map to Engine interface methods
const (
	OpSearch             = "google_search"
	OpSearchNews         = "google_search_news"
	OpSearchImages       = "google_search_images"
	OpSearchVideos       = "google_search_videos"
	OpSearchPlaces       = "google_search_places"
	OpSearchMaps         = "google_search_maps"
	OpSearchReviews      = "google_search_reviews"
	OpSearchShopping     = "google_search_shopping"
	OpSearchScholar      = "google_search_scholar"
	OpSearchLens         = "google_search_lens"
	OpSearchAutocomplete = "google_search_autocomplete"
	OpScrapeWebpage      = "webpage_scrape"
)

// ErrOperationNotSupported is returned when an operation is not supported by the current engine
var ErrOperationNotSupported = errors.New("operation not supported by current engine")

// Client is a unified SDK that fronts multiple search engine backends
type Client struct {
	registry *metasearch.Registry
	engine   metasearch.Engine
}

// New creates a new client with all available engines auto-registered
// It automatically selects an engine based on the SEARCH_ENGINE environment variable
func New() (*Client, error) {
	return NewWithOptions(nil)
}

// NewWithEngine creates a new client with a specific engine name selected
func NewWithEngine(engineName string) (*Client, error) {
	return NewWithOptions(&Options{
		EngineName: engineName,
	})
}

// Options configures the client behavior
type Options struct {
	// EngineName specifies which engine to use (e.g., "serper", "serpapi")
	// If empty, uses SEARCH_ENGINE env var or defaults to "serper"
	EngineName string

	// Silent suppresses initialization logs
	Silent bool
}

// NewWithOptions creates a new client with custom options
func NewWithOptions(opts *Options) (*Client, error) {
	if opts == nil {
		opts = &Options{}
	}

	registry := metasearch.NewRegistry()

	// Register all available engines
	if serperEngine, err := serper.New(); err == nil {
		registry.Register(serperEngine)
		if !opts.Silent {
			log.Printf("Registered Serper engine")
		}
	} else {
		if !opts.Silent {
			log.Printf("Failed to initialize Serper engine: %v", err)
		}
	}

	if serpApiEngine, err := serpapi.New(); err == nil {
		registry.Register(serpApiEngine)
		if !opts.Silent {
			log.Printf("Registered SerpAPI engine")
		}
	} else {
		if !opts.Silent {
			log.Printf("Failed to initialize SerpAPI engine: %v", err)
		}
	}

	client := &Client{
		registry: registry,
	}

	// Select the engine
	var engine metasearch.Engine
	var err error

	if opts.EngineName != "" {
		engine, err = client.GetEngine(opts.EngineName)
		if err != nil {
			return nil, err
		}
	} else {
		engine, err = metasearch.GetDefaultEngine(registry)
		if err != nil && !opts.Silent {
			log.Printf("Warning: %v", err)
		}
	}

	if engine == nil {
		return nil, fmt.Errorf("no search engines available. Please ensure API keys are set")
	}

	client.engine = engine

	if !opts.Silent {
		log.Printf("Using search engine: %s v%s", engine.GetName(), engine.GetVersion())
	}

	return client, nil
}

// GetEngine retrieves a specific engine by name
func (c *Client) GetEngine(name string) (metasearch.Engine, error) {
	engine, exists := c.registry.Get(name)
	if !exists {
		return nil, fmt.Errorf("engine '%s' not found. Available engines: %v", name, c.registry.List())
	}
	return engine, nil
}

// SetEngine sets the active engine by name
func (c *Client) SetEngine(name string) error {
	engine, err := c.GetEngine(name)
	if err != nil {
		return err
	}
	c.engine = engine
	return nil
}

// GetRegistry returns the underlying engine registry
func (c *Client) GetRegistry() *metasearch.Registry {
	return c.registry
}

// ListEngines returns all registered engine names
func (c *Client) ListEngines() []string {
	return c.registry.List()
}

// GetCurrentEngine returns the currently selected engine
func (c *Client) GetCurrentEngine() metasearch.Engine {
	return c.engine
}

// SupportsOperation checks if the current engine supports a specific operation
func (c *Client) SupportsOperation(operation string) bool {
	supportedTools := c.engine.GetSupportedTools()
	for _, tool := range supportedTools {
		if tool == operation {
			return true
		}
	}
	return false
}

// checkSupport returns an error if the operation is not supported by the current engine
func (c *Client) checkSupport(operation string) error {
	if !c.SupportsOperation(operation) {
		return fmt.Errorf("%w: '%s' (engine: %s, supported: %v)",
			ErrOperationNotSupported, operation, c.engine.GetName(), c.engine.GetSupportedTools())
	}
	return nil
}

// Engine interface methods - proxy to the selected engine

// GetName returns the name of the current search engine
func (c *Client) GetName() string {
	return c.engine.GetName()
}

// GetVersion returns the version of the current engine implementation
func (c *Client) GetVersion() string {
	return c.engine.GetVersion()
}

// GetSupportedTools returns a list of tool names supported by the current engine
func (c *Client) GetSupportedTools() []string {
	return c.engine.GetSupportedTools()
}

// Search performs a general web search
func (c *Client) Search(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	if err := c.checkSupport(OpSearch); err != nil {
		return nil, err
	}
	return c.engine.Search(ctx, params)
}

// SearchNews performs a news search
func (c *Client) SearchNews(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	if err := c.checkSupport(OpSearchNews); err != nil {
		return nil, err
	}
	return c.engine.SearchNews(ctx, params)
}

// SearchImages performs an image search
func (c *Client) SearchImages(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	if err := c.checkSupport(OpSearchImages); err != nil {
		return nil, err
	}
	return c.engine.SearchImages(ctx, params)
}

// SearchVideos performs a video search
func (c *Client) SearchVideos(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	if err := c.checkSupport(OpSearchVideos); err != nil {
		return nil, err
	}
	return c.engine.SearchVideos(ctx, params)
}

// SearchPlaces performs a places search
func (c *Client) SearchPlaces(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	if err := c.checkSupport(OpSearchPlaces); err != nil {
		return nil, err
	}
	return c.engine.SearchPlaces(ctx, params)
}

// SearchMaps performs a maps search
func (c *Client) SearchMaps(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	if err := c.checkSupport(OpSearchMaps); err != nil {
		return nil, err
	}
	return c.engine.SearchMaps(ctx, params)
}

// SearchReviews performs a reviews search
func (c *Client) SearchReviews(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	if err := c.checkSupport(OpSearchReviews); err != nil {
		return nil, err
	}
	return c.engine.SearchReviews(ctx, params)
}

// SearchShopping performs a shopping search
func (c *Client) SearchShopping(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	if err := c.checkSupport(OpSearchShopping); err != nil {
		return nil, err
	}
	return c.engine.SearchShopping(ctx, params)
}

// SearchScholar performs a scholar search
func (c *Client) SearchScholar(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	if err := c.checkSupport(OpSearchScholar); err != nil {
		return nil, err
	}
	return c.engine.SearchScholar(ctx, params)
}

// SearchLens performs a visual search (if supported)
func (c *Client) SearchLens(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	if err := c.checkSupport(OpSearchLens); err != nil {
		return nil, err
	}
	return c.engine.SearchLens(ctx, params)
}

// SearchAutocomplete gets search suggestions
func (c *Client) SearchAutocomplete(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	if err := c.checkSupport(OpSearchAutocomplete); err != nil {
		return nil, err
	}
	return c.engine.SearchAutocomplete(ctx, params)
}

// ScrapeWebpage scrapes content from a webpage
func (c *Client) ScrapeWebpage(ctx context.Context, params metasearch.ScrapeParams) (*metasearch.SearchResult, error) {
	if err := c.checkSupport(OpScrapeWebpage); err != nil {
		return nil, err
	}
	return c.engine.ScrapeWebpage(ctx, params)
}

// Normalized response methods - these return unified response structures across all engines

// SearchNormalized performs a web search and returns a normalized response
func (c *Client) SearchNormalized(ctx context.Context, params metasearch.SearchParams) (*metasearch.NormalizedSearchResult, error) {
	result, err := c.Search(ctx, params)
	if err != nil {
		return nil, err
	}

	normalizer := metasearch.NewNormalizer(c.GetName())
	return normalizer.NormalizeSearch(result, params.Query)
}

// SearchNewsNormalized performs a news search and returns a normalized response
func (c *Client) SearchNewsNormalized(ctx context.Context, params metasearch.SearchParams) (*metasearch.NormalizedSearchResult, error) {
	result, err := c.SearchNews(ctx, params)
	if err != nil {
		return nil, err
	}

	normalizer := metasearch.NewNormalizer(c.GetName())
	return normalizer.NormalizeNews(result, params.Query)
}

// SearchImagesNormalized performs an image search and returns a normalized response
func (c *Client) SearchImagesNormalized(ctx context.Context, params metasearch.SearchParams) (*metasearch.NormalizedSearchResult, error) {
	result, err := c.SearchImages(ctx, params)
	if err != nil {
		return nil, err
	}

	normalizer := metasearch.NewNormalizer(c.GetName())
	return normalizer.NormalizeImages(result, params.Query)
}
