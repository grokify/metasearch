package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/grokify/metaserp"
	"github.com/grokify/metaserp/client"
)

func main() {
	// Create client with default engine
	c, err := client.New()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Printf("Using engine: %s v%s\n", c.GetName(), c.GetVersion())
	fmt.Printf("Supported operations: %v\n\n", c.GetSupportedTools())

	// Check if current engine supports Lens search
	fmt.Println("=== Checking Lens Support ===")
	if c.SupportsOperation(client.OpSearchLens) {
		fmt.Printf("✓ %s supports Google Lens search\n", c.GetName())

		// Try to perform a lens search
		result, err := c.SearchLens(context.Background(), metaserp.SearchParams{
			Query: "red apple",
		})
		if err != nil {
			fmt.Printf("  Error performing lens search: %v\n", err)
		} else {
			fmt.Printf("  Lens search succeeded: %d bytes\n", len(result.Raw))
		}
	} else {
		fmt.Printf("✗ %s does NOT support Google Lens search\n", c.GetName())

		// Try anyway to demonstrate error handling
		_, err := c.SearchLens(context.Background(), metaserp.SearchParams{
			Query: "red apple",
		})
		if errors.Is(err, client.ErrOperationNotSupported) {
			fmt.Printf("  Got expected error: %v\n", err)
		}
	}

	// Demonstrate switching engines
	fmt.Println("\n=== Switching Engines ===")
	availableEngines := c.ListEngines()
	fmt.Printf("Available engines: %v\n", availableEngines)

	for _, engineName := range availableEngines {
		if engineName == c.GetName() {
			continue // Skip current engine
		}

		fmt.Printf("\nSwitching to: %s\n", engineName)
		if err := c.SetEngine(engineName); err != nil {
			fmt.Printf("  Error switching: %v\n", err)
			continue
		}

		fmt.Printf("  Engine: %s v%s\n", c.GetName(), c.GetVersion())
		fmt.Printf("  Lens support: %v\n", c.SupportsOperation(client.OpSearchLens))
	}

	// Example: Check multiple operations
	fmt.Println("\n=== Operation Support Matrix ===")
	operations := []struct {
		name string
		op   string
	}{
		{"Web Search", client.OpSearch},
		{"News Search", client.OpSearchNews},
		{"Image Search", client.OpSearchImages},
		{"Lens Search", client.OpSearchLens},
		{"Shopping Search", client.OpSearchShopping},
		{"Scholar Search", client.OpSearchScholar},
	}

	for _, engineName := range c.ListEngines() {
		if err := c.SetEngine(engineName); err != nil {
			continue
		}

		fmt.Printf("\n%s:\n", engineName)
		for _, op := range operations {
			supported := "✓"
			if !c.SupportsOperation(op.op) {
				supported = "✗"
			}
			fmt.Printf("  %s %s\n", supported, op.name)
		}
	}
}
