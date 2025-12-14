package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/metasearch"
	"github.com/grokify/metasearch/client/serpapi"
	"github.com/grokify/metasearch/client/serper"
)

func main() {
	// Initialize engine registry with all available engines
	registry := metasearch.NewRegistry()

	// Register available search engines
	if serperEngine, err := serper.New(); err == nil {
		registry.Register(serperEngine)
		log.Printf("Registered Serper engine")
	} else {
		log.Printf("Failed to initialize Serper engine: %v", err)
	}

	if serpApiEngine, err := serpapi.New(); err == nil {
		registry.Register(serpApiEngine)
		log.Printf("Registered SerpAPI engine")
	} else {
		log.Printf("Failed to initialize SerpAPI engine: %v", err)
	}

	// Get the default/selected engine
	searchEngine, err := metasearch.GetDefaultEngine(registry)
	if err != nil {
		log.Printf("Warning: %v", err)
	}
	if searchEngine == nil {
		log.Fatal("No search engines available. Please ensure API keys are set.")
	}

	log.Printf("Using search engine: %s v%s", searchEngine.GetName(), searchEngine.GetVersion())
	log.Printf("Available engines: %v", registry.List())

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "multi-search-server",
		Version: "2.0.0",
	}, nil)

	// Register search tools dynamically based on supported tools
	registerSearchTool := func(toolName, description string, searchFunc func(context.Context, metasearch.SearchParams) (*metasearch.SearchResult, error)) {
		mcp.AddTool(server, &mcp.Tool{
			Name:        toolName,
			Description: description,
		}, func(ctx context.Context, req *mcp.CallToolRequest, args metasearch.SearchParams) (*mcp.CallToolResult, any, error) {
			result, err := searchFunc(ctx, args)
			if err != nil {
				return nil, nil, fmt.Errorf("%s failed: %w", toolName, err)
			}

			resultJSON, _ := json.MarshalIndent(result.Data, "", "  ")
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: string(resultJSON)},
				},
			}, nil, nil
		})
	}

	// Register all search tools
	registerSearchTool("google_search", "Perform a Google web search", searchEngine.Search)
	registerSearchTool("google_search_news", "Search for news articles using Google News", searchEngine.SearchNews)
	registerSearchTool("google_search_images", "Search for images using Google Images", searchEngine.SearchImages)
	registerSearchTool("google_search_videos", "Search for videos using Google Videos", searchEngine.SearchVideos)
	registerSearchTool("google_search_places", "Search for places using Google Places", searchEngine.SearchPlaces)
	registerSearchTool("google_search_maps", "Search for locations using Google Maps", searchEngine.SearchMaps)
	registerSearchTool("google_search_reviews", "Search for reviews", searchEngine.SearchReviews)
	registerSearchTool("google_search_shopping", "Search for products using Google Shopping", searchEngine.SearchShopping)
	registerSearchTool("google_search_scholar", "Search for academic papers using Google Scholar", searchEngine.SearchScholar)
	registerSearchTool("google_search_lens", "Perform visual search using Google Lens", searchEngine.SearchLens)
	registerSearchTool("google_search_autocomplete", "Get search suggestions using Google Autocomplete", searchEngine.SearchAutocomplete)

	// Web scraping tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "webpage_scrape",
		Description: "Scrape content from a webpage",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args metasearch.ScrapeParams) (*mcp.CallToolResult, any, error) {
		result, err := searchEngine.ScrapeWebpage(ctx, args)
		if err != nil {
			return nil, nil, fmt.Errorf("scraping failed: %w", err)
		}

		resultJSON, _ := json.MarshalIndent(result.Data, "", "  ")
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(resultJSON)},
			},
		}, nil, nil
	})

	log.Printf("Starting Multi-Search MCP Server with %s engine...", searchEngine.GetName())
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Printf("Server failed: %v", err)
	}
}
