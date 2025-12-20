package serpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/grokify/metaserp"
)

const (
	baseURL       = "https://serpapi.com"
	engineName    = "serpapi"
	engineVersion = "1.0.0"
	searchURL     = "https://serpapi.com/search.json"
)

// Engine implements the metaserp.Engine interface for SerpAPI
type Engine struct {
	apiKey string
	client *http.Client
}

// New creates a new SerpAPI engine instance
func New() (*Engine, error) {
	apiKey := os.Getenv("SERPAPI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("SERPAPI_API_KEY environment variable is required")
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
		// Note: google_search_lens is NOT supported by SerpAPI
		"google_search_autocomplete",
		"webpage_scrape",
	}
}

// makeRequest performs HTTP request to SerpAPI
func (e *Engine) makeRequest(params map[string]string) (*metaserp.SearchResult, error) {
	// Build URL with query parameters
	reqURL, err := url.Parse(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// Add API key and other parameters
	q := reqURL.Query()
	q.Set("api_key", e.apiKey)
	for key, value := range params {
		q.Set(key, value)
	}
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

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

	return &metaserp.SearchResult{
		Data: result,
		Raw:  string(body),
	}, nil
}

// buildParams converts SearchParams to SerpAPI parameters
func (e *Engine) buildParams(params metaserp.SearchParams, engine string) map[string]string {
	apiParams := map[string]string{
		"q":      params.Query,
		"engine": engine,
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
		apiParams["num"] = fmt.Sprintf("%d", params.NumResults)
	}

	return apiParams
}

// Search performs a general web search
func (e *Engine) Search(ctx context.Context, params metaserp.SearchParams) (*metaserp.SearchResult, error) {
	return e.makeRequest(e.buildParams(params, "google"))
}

// SearchNews performs a news search
func (e *Engine) SearchNews(ctx context.Context, params metaserp.SearchParams) (*metaserp.SearchResult, error) {
	return e.makeRequest(e.buildParams(params, "google_news"))
}

// SearchImages performs an image search
func (e *Engine) SearchImages(ctx context.Context, params metaserp.SearchParams) (*metaserp.SearchResult, error) {
	return e.makeRequest(e.buildParams(params, "google_images"))
}

// SearchVideos performs a video search
func (e *Engine) SearchVideos(ctx context.Context, params metaserp.SearchParams) (*metaserp.SearchResult, error) {
	return e.makeRequest(e.buildParams(params, "google_videos"))
}

// SearchPlaces performs a places search
func (e *Engine) SearchPlaces(ctx context.Context, params metaserp.SearchParams) (*metaserp.SearchResult, error) {
	// For places, we use Google Maps search with type parameter
	apiParams := e.buildParams(params, "google_maps")
	apiParams["type"] = "search"
	return e.makeRequest(apiParams)
}

// SearchMaps performs a maps search
func (e *Engine) SearchMaps(ctx context.Context, params metaserp.SearchParams) (*metaserp.SearchResult, error) {
	return e.makeRequest(e.buildParams(params, "google_maps"))
}

// SearchReviews performs a reviews search
func (e *Engine) SearchReviews(ctx context.Context, params metaserp.SearchParams) (*metaserp.SearchResult, error) {
	// Reviews can be searched through Google with specific query modification
	apiParams := e.buildParams(params, "google")
	apiParams["q"] = params.Query + " reviews"
	return e.makeRequest(apiParams)
}

// SearchShopping performs a shopping search
func (e *Engine) SearchShopping(ctx context.Context, params metaserp.SearchParams) (*metaserp.SearchResult, error) {
	return e.makeRequest(e.buildParams(params, "google_shopping"))
}

// SearchScholar performs a scholar search
func (e *Engine) SearchScholar(ctx context.Context, params metaserp.SearchParams) (*metaserp.SearchResult, error) {
	apiParams := map[string]string{
		"q":      params.Query,
		"engine": "google_scholar",
	}

	if params.Language != "" {
		apiParams["hl"] = params.Language
	}
	if params.NumResults > 0 {
		apiParams["num"] = fmt.Sprintf("%d", params.NumResults)
	}

	return e.makeRequest(apiParams)
}

// SearchLens performs a visual search (not supported by SerpAPI)
func (e *Engine) SearchLens(ctx context.Context, params metaserp.SearchParams) (*metaserp.SearchResult, error) {
	return nil, fmt.Errorf("google_search_lens is not supported by SerpAPI")
}

// SearchAutocomplete gets search suggestions
func (e *Engine) SearchAutocomplete(ctx context.Context, params metaserp.SearchParams) (*metaserp.SearchResult, error) {
	apiParams := map[string]string{
		"q":      params.Query,
		"engine": "google_autocomplete",
	}

	if params.Language != "" {
		apiParams["hl"] = params.Language
	}
	if params.Country != "" {
		apiParams["gl"] = params.Country
	}

	return e.makeRequest(apiParams)
}

// ScrapeWebpage scrapes content from a webpage (using SerpAPI's custom scraping)
func (e *Engine) ScrapeWebpage(ctx context.Context, params metaserp.ScrapeParams) (*metaserp.SearchResult, error) {
	// Validate URL
	if _, err := url.Parse(params.URL); err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// SerpAPI doesn't have a direct scraping endpoint like Serper
	// We'll implement a basic HTTP scraping here
	req, err := http.NewRequest(http.MethodGet, params.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape webpage: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("scraping error: status %d", resp.StatusCode)
	}

	// Return the raw HTML content
	result := map[string]any{
		"url":     params.URL,
		"content": string(body),
		"status":  resp.StatusCode,
		"headers": resp.Header,
	}

	return &metaserp.SearchResult{
		Data: result,
		Raw:  string(body),
	}, nil
}
