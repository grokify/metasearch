package client

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/agentplexus/omniserp"
)

// TestCapabilityChecking tests that the client properly validates operation support
func TestCapabilityChecking(t *testing.T) {
	// Note: These tests require actual API keys to initialize engines
	// They will be skipped if the API keys are not available

	// Test with Serper (supports Lens)
	t.Run("Serper supports Lens", func(t *testing.T) {
		c, err := NewWithEngine("serper")
		if err != nil {
			t.Skip("Serper engine not available (likely missing API key)")
		}

		if !c.SupportsOperation(OpSearchLens) {
			t.Error("Serper should support google_search_lens")
		}

		// This should NOT return ErrOperationNotSupported with Serper
		_, err = c.SearchLens(context.Background(), omniserp.SearchParams{Query: "test"})
		if errors.Is(err, ErrOperationNotSupported) {
			t.Error("Serper should support SearchLens but got unsupported error")
		}
		// Note: May still get API errors, but not unsupported errors
	})

	// Test with SerpAPI (does NOT support Lens)
	t.Run("SerpAPI does not support Lens", func(t *testing.T) {
		c, err := NewWithEngine("serpapi")
		if err != nil {
			t.Skip("SerpAPI engine not available (likely missing API key)")
		}

		if c.SupportsOperation(OpSearchLens) {
			t.Error("SerpAPI should NOT support google_search_lens")
		}

		// This should return ErrOperationNotSupported with SerpAPI
		_, err = c.SearchLens(context.Background(), omniserp.SearchParams{Query: "test"})
		if !errors.Is(err, ErrOperationNotSupported) {
			t.Errorf("Expected ErrOperationNotSupported, got: %v", err)
		}
	})

	// Test all engines support basic search
	t.Run("All engines support basic search", func(t *testing.T) {
		c, err := New()
		if err != nil {
			t.Skip("No engines available (likely missing API keys)")
		}

		if !c.SupportsOperation(OpSearch) {
			t.Error("All engines should support basic google_search")
		}
	})
}

// TestActualSearchWithSerper performs an actual search if SERPER_API_KEY is set
func TestActualSearchWithSerper(t *testing.T) {
	if os.Getenv("SERPER_API_KEY") == "" {
		t.Skip("SERPER_API_KEY not set, skipping live API test")
	}

	c, err := NewWithEngine("serper")
	if err != nil {
		t.Fatalf("Failed to create Serper client: %v", err)
	}

	t.Run("Basic web search", func(t *testing.T) {
		result, err := c.Search(context.Background(), omniserp.SearchParams{
			Query:      "golang programming",
			NumResults: 5,
		})
		if err != nil {
			t.Fatalf("Search failed: %v", err)
		}
		if result == nil {
			t.Error("Expected non-nil result")
		}
		if result.Data == nil {
			t.Error("Expected non-nil result data")
		}
		t.Logf("Search succeeded, got %d bytes of data", len(result.Raw))
	})

	t.Run("Lens search (supported by Serper)", func(t *testing.T) {
		if !c.SupportsOperation(OpSearchLens) {
			t.Fatal("Serper should support Lens")
		}

		result, err := c.SearchLens(context.Background(), omniserp.SearchParams{
			Query:      "red apple",
			NumResults: 5,
		})
		if err != nil {
			t.Logf("Lens search failed (may be expected): %v", err)
			// Don't fail the test - Lens might have API issues
			return
		}
		if result == nil {
			t.Error("Expected non-nil result")
		}
		t.Logf("Lens search succeeded, got %d bytes of data", len(result.Raw))
	})
}

// TestActualSearchWithSerpAPI performs an actual search if SERPAPI_API_KEY is set
func TestActualSearchWithSerpAPI(t *testing.T) {
	if os.Getenv("SERPAPI_API_KEY") == "" {
		t.Skip("SERPAPI_API_KEY not set, skipping live API test")
	}

	c, err := NewWithEngine("serpapi")
	if err != nil {
		t.Fatalf("Failed to create SerpAPI client: %v", err)
	}

	t.Run("Basic web search", func(t *testing.T) {
		result, err := c.Search(context.Background(), omniserp.SearchParams{
			Query:      "golang programming",
			NumResults: 5,
		})
		if err != nil {
			t.Fatalf("Search failed: %v", err)
		}
		if result == nil {
			t.Error("Expected non-nil result")
		}
		if result.Data == nil {
			t.Error("Expected non-nil result data")
		}
		t.Logf("Search succeeded, got %d bytes of data", len(result.Raw))
	})

	t.Run("Lens search (NOT supported by SerpAPI)", func(t *testing.T) {
		if c.SupportsOperation(OpSearchLens) {
			t.Fatal("SerpAPI should NOT support Lens")
		}

		_, err := c.SearchLens(context.Background(), omniserp.SearchParams{
			Query: "red apple",
		})
		if !errors.Is(err, ErrOperationNotSupported) {
			t.Errorf("Expected ErrOperationNotSupported, got: %v", err)
		} else {
			t.Logf("Correctly returned unsupported error: %v", err)
		}
	})
}

// TestEngineSwitching tests switching between engines at runtime
func TestEngineSwitching(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Skip("No engines available (likely missing API keys)")
	}

	initialEngine := c.GetName()
	t.Logf("Initial engine: %s", initialEngine)

	availableEngines := c.ListEngines()
	if len(availableEngines) < 2 {
		t.Skip("Need at least 2 engines to test switching")
	}

	for _, engineName := range availableEngines {
		if engineName != initialEngine {
			err := c.SetEngine(engineName)
			if err != nil {
				t.Errorf("Failed to switch to %s: %v", engineName, err)
			}
			if c.GetName() != engineName {
				t.Errorf("Expected engine %s, got %s", engineName, c.GetName())
			}
			t.Logf("Switched to engine: %s", c.GetName())
			break
		}
	}
}

// TestOperationConstants verifies that all operation constants are defined
func TestOperationConstants(t *testing.T) {
	operations := []string{
		OpSearch,
		OpSearchNews,
		OpSearchImages,
		OpSearchVideos,
		OpSearchPlaces,
		OpSearchMaps,
		OpSearchReviews,
		OpSearchShopping,
		OpSearchScholar,
		OpSearchLens,
		OpSearchAutocomplete,
		OpScrapeWebpage,
	}

	for _, op := range operations {
		if op == "" {
			t.Error("Found empty operation constant")
		}
	}
}
