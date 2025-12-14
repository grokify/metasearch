package client

import (
	"context"
	"errors"
	"testing"

	"github.com/grokify/metasearch"
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

		// This should succeed with Serper
		_, err = c.SearchLens(context.Background(), metasearch.SearchParams{Query: "test"})
		if err != nil && errors.Is(err, ErrOperationNotSupported) {
			t.Error("Serper should support SearchLens but got unsupported error")
		}
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
		_, err = c.SearchLens(context.Background(), metasearch.SearchParams{Query: "test"})
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
