package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/grokify/metasearch"
	"github.com/grokify/metasearch/serpapi"
	"github.com/grokify/metasearch/serper"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: metasearch <query> [engine]")
		fmt.Println("Available engines: serper, serpapi")
		fmt.Println("Set SEARCH_ENGINE environment variable to specify default engine")
		os.Exit(1)
	}

	query := os.Args[1]
	var selectedEngine string
	if len(os.Args) > 2 {
		selectedEngine = os.Args[2]
	}

	// Create registry and register all available engines
	registry := createRegistry()

	// Get the engine to use
	var engine metasearch.Engine
	var err error

	if selectedEngine != "" {
		engineSelected, exists := registry.Get(selectedEngine)
		if !exists {
			log.Fatalf("Engine '%s' not found. Available engines: %v", selectedEngine, registry.List())
		}
		engine = engineSelected
	} else {
		engine, err = metasearch.GetDefaultEngine(registry)
		if err != nil {
			log.Printf("Warning: %v", err)
		}
	}

	if engine == nil {
		log.Fatal("No search engine available")
	}

	// Perform search
	params := metasearch.SearchParams{
		Query:      query,
		NumResults: 10,
	}

	result, err := engine.Search(context.Background(), params)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	// Output results
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal results: %v", err)
	}

	fmt.Println(string(output))
}

// createRegistry creates a new registry with all available engines pre-registered
func createRegistry() *metasearch.Registry {
	registry := metasearch.NewRegistry()

	// Register available engines
	if serperEngine, err := serper.New(); err == nil {
		registry.Register(serperEngine)
	} else {
		log.Printf("Failed to initialize Serper engine: %v", err)
	}

	if serpApiEngine, err := serpapi.New(); err == nil {
		registry.Register(serpApiEngine)
	} else {
		log.Printf("Failed to initialize SerpAPI engine: %v", err)
	}

	return registry
}
