package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	flags "github.com/jessevdk/go-flags"

	"github.com/agentplexus/omniserp"
	"github.com/agentplexus/omniserp/client"
)

type Options struct {
	Engine string `short:"e" long:"engine" description:"Search engine (serper, serpapi)" required:"true"`
	Query  string `short:"q" long:"query" description:"Query" required:"true"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	query := opts.Query

	// Create client SDK
	var c *client.Client

	c, err = client.NewWithEngine(opts.Engine)
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	// Perform search
	params := omniserp.SearchParams{
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
