package serper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	metasearch "github.com/grokify/metaserp"
)

const (
	baseURL       = "https://google.serper.dev"
	engineName    = "serper"
	engineVersion = "1.0.0"
)

// Engine implements the metasearch.Engine interface for Serper API
type Engine struct {
	apiKey string
	client *http.Client
}

// New creates a new Serper engine instance
func New() (*Engine, error) {
	apiKey := os.Getenv("SERPER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("SERPER_API_KEY environment variable is required")
	}

	return &Engine{
		apiKey: apiKey,
		client: &http.Client{},
	}, nil
}

// GetName returns the engine name
func (e *Engine) GetName() string {
	return engineName
}

// GetVersion returns the engine version
func (e *Engine) GetVersion() string {
	return engineVersion
}

// GetSupportedTools returns the list of supported tools
func (e *Engine) GetSupportedTools() []string {
	return []string{
		"google_search",
		"google_search_news",
		"google_search_images",
		"google_search_videos",
		"google_search_places",
		"google_search_maps",
		"google_search_reviews",
		"google_search_shopping",
		"google_search_scholar",
		"google_search_lens",
		"google_search_autocomplete",
		"webpage_scrape",
	}
}

// makeRequest performs HTTP request to Serper API
func (e *Engine) makeRequest(endpoint string, params map[string]interface{}) (*metasearch.SearchResult, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+endpoint, strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-KEY", e.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &metasearch.SearchResult{
		Data: result,
		Raw:  string(body),
	}, nil
}

// buildParams converts SearchParams to API parameters
func (e *Engine) buildParams(params metasearch.SearchParams) map[string]any {
	apiParams := map[string]any{
		"q": params.Query,
	}

	if params.Location != "" {
		apiParams["location"] = params.Location
	}
	if params.Language != "" {
		apiParams["hl"] = params.Language
	}
	if params.Country != "" {
		apiParams["gl"] = params.Country
	}
	if params.NumResults > 0 {
		apiParams["num"] = params.NumResults
	}

	return apiParams
}

// Search performs a general web search
func (e *Engine) Search(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	return e.makeRequest("/search", e.buildParams(params))
}

// SearchNews performs a news search
func (e *Engine) SearchNews(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	return e.makeRequest("/news", e.buildParams(params))
}

// SearchImages performs an image search
func (e *Engine) SearchImages(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	return e.makeRequest("/images", e.buildParams(params))
}

// SearchVideos performs a video search
func (e *Engine) SearchVideos(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	return e.makeRequest("/videos", e.buildParams(params))
}

// SearchPlaces performs a places search
func (e *Engine) SearchPlaces(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	return e.makeRequest("/places", e.buildParams(params))
}

// SearchMaps performs a maps search
func (e *Engine) SearchMaps(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	return e.makeRequest("/maps", e.buildParams(params))
}

// SearchReviews performs a reviews search
func (e *Engine) SearchReviews(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	return e.makeRequest("/reviews", e.buildParams(params))
}

// SearchShopping performs a shopping search
func (e *Engine) SearchShopping(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	return e.makeRequest("/shopping", e.buildParams(params))
}

// SearchScholar performs a scholar search
func (e *Engine) SearchScholar(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	// Scholar search has different parameter requirements
	apiParams := map[string]interface{}{
		"q": params.Query,
	}
	if params.Language != "" {
		apiParams["hl"] = params.Language
	}
	if params.NumResults > 0 {
		apiParams["num"] = params.NumResults
	}

	return e.makeRequest("/scholar", apiParams)
}

// SearchLens performs a visual search
func (e *Engine) SearchLens(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	// Lens search has limited parameters
	apiParams := map[string]interface{}{
		"q": params.Query,
	}
	if params.Language != "" {
		apiParams["hl"] = params.Language
	}
	if params.Country != "" {
		apiParams["gl"] = params.Country
	}
	if params.NumResults > 0 {
		apiParams["num"] = params.NumResults
	}

	return e.makeRequest("/lens", apiParams)
}

// SearchAutocomplete gets search suggestions
func (e *Engine) SearchAutocomplete(ctx context.Context, params metasearch.SearchParams) (*metasearch.SearchResult, error) {
	// Autocomplete has limited parameters
	apiParams := map[string]interface{}{
		"q": params.Query,
	}
	if params.Language != "" {
		apiParams["hl"] = params.Language
	}
	if params.Country != "" {
		apiParams["gl"] = params.Country
	}

	return e.makeRequest("/autocomplete", apiParams)
}

// ScrapeWebpage scrapes content from a webpage
func (e *Engine) ScrapeWebpage(ctx context.Context, params metasearch.ScrapeParams) (*metasearch.SearchResult, error) {
	// Validate URL
	if _, err := url.Parse(params.URL); err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	apiParams := map[string]interface{}{
		"url": params.URL,
	}

	return e.makeRequest("/scrape", apiParams)
}
