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

// ToolDefinition defines a search tool with its metadata
type ToolDefinition struct {
	Name        string
	Description string
	SearchFunc  func(context.Context, metasearch.SearchParams) (*metasearch.SearchResult, error)
}

func main() {
	// Initialize client SDK with all available engines
	searchClient, err := client.New()
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	log.Printf("Using engine: %s v%s", searchClient.GetName(), searchClient.GetVersion())
	log.Printf("Available engines: %v", searchClient.ListEngines())

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "multi-search-server",
		Version: "2.0.0",
	}, nil)

	// Define all possible search tools with their operation names
	allTools := []ToolDefinition{
		{client.OpSearch, "Perform a Google web search", searchClient.Search},
		{client.OpSearchNews, "Search for news articles using Google News", searchClient.SearchNews},
		{client.OpSearchImages, "Search for images using Google Images", searchClient.SearchImages},
		{client.OpSearchVideos, "Search for videos using Google Videos", searchClient.SearchVideos},
		{client.OpSearchPlaces, "Search for places using Google Places", searchClient.SearchPlaces},
		{client.OpSearchMaps, "Search for locations using Google Maps", searchClient.SearchMaps},
		{client.OpSearchReviews, "Search for reviews", searchClient.SearchReviews},
		{client.OpSearchShopping, "Search for products using Google Shopping", searchClient.SearchShopping},
		{client.OpSearchScholar, "Search for academic papers using Google Scholar", searchClient.SearchScholar},
		{client.OpSearchLens, "Perform visual search using Google Lens", searchClient.SearchLens},
		{client.OpSearchAutocomplete, "Get search suggestions using Google Autocomplete", searchClient.SearchAutocomplete},
	}

	// Register tools only if supported by the current engine
	registeredTools := []string{}
	skippedTools := []string{}

	for _, tool := range allTools {
		if searchClient.SupportsOperation(tool.Name) {
			// Register this tool
			toolName := tool.Name
			toolDesc := tool.Description
			searchFunc := tool.SearchFunc

			mcp.AddTool(server, &mcp.Tool{
				Name:        toolName,
				Description: toolDesc,
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

			registeredTools = append(registeredTools, tool.Name)
		} else {
			skippedTools = append(skippedTools, tool.Name)
		}
	}

	// Register web scraping tool if supported
	if searchClient.SupportsOperation(client.OpScrapeWebpage) {
		mcp.AddTool(server, &mcp.Tool{
			Name:        client.OpScrapeWebpage,
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
		registeredTools = append(registeredTools, client.OpScrapeWebpage)
	} else {
		skippedTools = append(skippedTools, client.OpScrapeWebpage)
	}

	// Log tool registration summary
	log.Printf("Registered %d tools: %v", len(registeredTools), registeredTools)
	if len(skippedTools) > 0 {
		log.Printf("Skipped %d unsupported tools: %v", len(skippedTools), skippedTools)
	}

	log.Printf("Starting Multi-Search MCP Server with %s engine...", searchClient.GetName())
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Printf("Server failed: %v", err)
	}
}
