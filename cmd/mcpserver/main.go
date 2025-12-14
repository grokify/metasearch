package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/metasearch"
	"github.com/grokify/metasearch/client"
)

func main() {
	// Initialize client SDK with all available engines
	searchClient, err := client.New()
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	log.Printf("Available engines: %v", searchClient.ListEngines())

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
	registerSearchTool("google_search", "Perform a Google web search", searchClient.Search)
	registerSearchTool("google_search_news", "Search for news articles using Google News", searchClient.SearchNews)
	registerSearchTool("google_search_images", "Search for images using Google Images", searchClient.SearchImages)
	registerSearchTool("google_search_videos", "Search for videos using Google Videos", searchClient.SearchVideos)
	registerSearchTool("google_search_places", "Search for places using Google Places", searchClient.SearchPlaces)
	registerSearchTool("google_search_maps", "Search for locations using Google Maps", searchClient.SearchMaps)
	registerSearchTool("google_search_reviews", "Search for reviews", searchClient.SearchReviews)
	registerSearchTool("google_search_shopping", "Search for products using Google Shopping", searchClient.SearchShopping)
	registerSearchTool("google_search_scholar", "Search for academic papers using Google Scholar", searchClient.SearchScholar)
	registerSearchTool("google_search_lens", "Perform visual search using Google Lens", searchClient.SearchLens)
	registerSearchTool("google_search_autocomplete", "Get search suggestions using Google Autocomplete", searchClient.SearchAutocomplete)

	// Web scraping tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "webpage_scrape",
		Description: "Scrape content from a webpage",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args metasearch.ScrapeParams) (*mcp.CallToolResult, any, error) {
		result, err := searchClient.ScrapeWebpage(ctx, args)
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

	log.Printf("Starting Multi-Search MCP Server with %s engine...", searchClient.GetName())
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Printf("Server failed: %v", err)
	}
}
