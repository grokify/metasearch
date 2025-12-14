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
	if len(os.Args) < 2 {
		fmt.Println("Usage: metasearch <query> [engine]")
		fmt.Println("Available engines: serper, serpapi")
		fmt.Println("Set SEARCH_ENGINE environment variable to specify default engine")
		os.Exit(1)
	}

	query := os.Args[1]
	var engineName string
	if len(os.Args) > 2 {
		engineName = os.Args[2]
	}

	// Create client SDK
	var c *client.Client
	var err error

	if engineName != "" {
		c, err = client.NewWithEngine(engineName)
	} else {
		c, err = client.New()
	}

	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	// Perform search
	params := metasearch.SearchParams{
		Query:      query,
		NumResults: 10,
	}

	result, err := c.Search(context.Background(), params)
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
