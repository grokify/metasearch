package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/grokify/metasearch"
	"github.com/grokify/metasearch/client"
)

func main() {
	// Create client
	c, err := client.New()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Printf("Using engine: %s v%s\n\n", c.GetName(), c.GetVersion())

	query := "golang programming"
	if len(os.Args) > 1 {
		query = os.Args[1]
	}

	// Demonstrate both raw and normalized responses
	fmt.Println("=== Raw Response (Engine-Specific) ===")
	rawResult, err := c.Search(context.Background(), metasearch.SearchParams{
		Query:      query,
		NumResults: 3,
	})
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	// Raw response structure differs between engines
	if data, ok := rawResult.Data.(map[string]any); ok {
		fmt.Printf("Raw response keys: %v\n\n", getKeys(data))
	}

	// Now get normalized response
	fmt.Println("=== Normalized Response (Engine-Agnostic) ===")
	normalized, err := c.SearchNormalized(context.Background(), metasearch.SearchParams{
		Query:      query,
		NumResults: 3,
	})
	if err != nil {
		log.Fatalf("Normalized search failed: %v", err)
	}

	// Print metadata
	fmt.Printf("Engine: %s\n", normalized.SearchMetadata.Engine)
	fmt.Printf("Query: %s\n", normalized.SearchMetadata.Query)
	fmt.Printf("Total Organic Results: %d\n\n", len(normalized.OrganicResults))

	// Print organic results (same structure regardless of engine)
	for i, result := range normalized.OrganicResults {
		fmt.Printf("%d. %s\n", i+1, result.Title)
		fmt.Printf("   URL: %s\n", result.Link)
		fmt.Printf("   Snippet: %s\n", truncate(result.Snippet, 100))
		fmt.Println()
	}

	// Print answer box if present
	if normalized.AnswerBox != nil {
		fmt.Println("=== Answer Box ===")
		fmt.Printf("Title: %s\n", normalized.AnswerBox.Title)
		fmt.Printf("Answer: %s\n", truncate(normalized.AnswerBox.Answer, 200))
		fmt.Println()
	}

	// Print knowledge graph if present
	if normalized.KnowledgeGraph != nil {
		fmt.Println("=== Knowledge Graph ===")
		fmt.Printf("Title: %s\n", normalized.KnowledgeGraph.Title)
		fmt.Printf("Type: %s\n", normalized.KnowledgeGraph.Type)
		fmt.Printf("Description: %s\n", truncate(normalized.KnowledgeGraph.Description, 200))
		fmt.Println()
	}

	// Print related searches if present
	if len(normalized.RelatedSearches) > 0 {
		fmt.Println("=== Related Searches ===")
		for _, related := range normalized.RelatedSearches {
			fmt.Printf("- %s\n", related.Query)
		}
		fmt.Println()
	}

	// Demonstrate engine switching with normalized responses
	fmt.Println("=== Switching Engines ===")
	availableEngines := c.ListEngines()
	fmt.Printf("Available engines: %v\n\n", availableEngines)

	for _, engineName := range availableEngines {
		if engineName == c.GetName() {
			continue // Skip current engine
		}

		fmt.Printf("Switching to: %s\n", engineName)
		if err := c.SetEngine(engineName); err != nil {
			fmt.Printf("  Error: %v\n\n", err)
			continue
		}

		// Same normalized API works with different engine!
		result, err := c.SearchNormalized(context.Background(), metasearch.SearchParams{
			Query:      query,
			NumResults: 2,
		})
		if err != nil {
			fmt.Printf("  Error: %v\n\n", err)
			continue
		}

		fmt.Printf("  Found %d results\n", len(result.OrganicResults))
		if len(result.OrganicResults) > 0 {
			fmt.Printf("  First result: %s\n", result.OrganicResults[0].Title)
		}
		fmt.Println()
	}

	// Show full JSON structure of normalized response
	fmt.Println("=== Full Normalized JSON ===")
	jsonData, _ := json.MarshalIndent(normalized, "", "  ")
	fmt.Println(string(jsonData))
}

func getKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
