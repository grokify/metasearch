package metasearch

import (
	"testing"
)

func TestNormalizeSerperSearch(t *testing.T) {
	// Mock Serper response
	serperData := map[string]any{
		"organic": []any{
			map[string]any{
				"title":   "Go Programming Language",
				"link":    "https://golang.org",
				"snippet": "Build simple, secure, scalable systems with Go",
				"date":    "2023-01-01",
			},
			map[string]any{
				"title":   "Go Documentation",
				"link":    "https://golang.org/doc",
				"snippet": "The Go programming language documentation",
			},
		},
		"answerBox": map[string]any{
			"title":   "What is Go?",
			"answer":  "Go is a statically typed, compiled programming language",
			"snippet": "Designed at Google",
		},
		"relatedSearches": []any{
			map[string]any{"query": "golang tutorial"},
			map[string]any{"query": "golang vs python"},
		},
	}

	result := &SearchResult{
		Data: serperData,
		Raw:  `{"organic":[...]}`,
	}

	normalizer := NewNormalizer("serper")
	normalized, err := normalizer.NormalizeSearch(result, "golang programming")

	if err != nil {
		t.Fatalf("NormalizeSearch failed: %v", err)
	}

	// Verify organic results
	if len(normalized.OrganicResults) != 2 {
		t.Errorf("Expected 2 organic results, got %d", len(normalized.OrganicResults))
	}

	if normalized.OrganicResults[0].Title != "Go Programming Language" {
		t.Errorf("Expected title 'Go Programming Language', got '%s'", normalized.OrganicResults[0].Title)
	}

	if normalized.OrganicResults[0].Link != "https://golang.org" {
		t.Errorf("Expected link 'https://golang.org', got '%s'", normalized.OrganicResults[0].Link)
	}

	// Verify answer box
	if normalized.AnswerBox == nil {
		t.Fatal("Expected answer box to be present")
	}

	if normalized.AnswerBox.Title != "What is Go?" {
		t.Errorf("Expected answer box title 'What is Go?', got '%s'", normalized.AnswerBox.Title)
	}

	// Verify related searches
	if len(normalized.RelatedSearches) != 2 {
		t.Errorf("Expected 2 related searches, got %d", len(normalized.RelatedSearches))
	}

	// Verify metadata
	if normalized.SearchMetadata.Engine != "serper" {
		t.Errorf("Expected engine 'serper', got '%s'", normalized.SearchMetadata.Engine)
	}

	if normalized.SearchMetadata.Query != "golang programming" {
		t.Errorf("Expected query 'golang programming', got '%s'", normalized.SearchMetadata.Query)
	}
}

func TestNormalizeSerpAPISearch(t *testing.T) {
	// Mock SerpAPI response (note different field names)
	serpAPIData := map[string]any{
		"organic_results": []any{
			map[string]any{
				"title":   "Go Programming Language",
				"link":    "https://golang.org",
				"snippet": "Build simple, secure, scalable systems with Go",
			},
			map[string]any{
				"title":   "Go Documentation",
				"link":    "https://golang.org/doc",
				"snippet": "The Go programming language documentation",
			},
		},
		"answer_box": map[string]any{
			"title":   "What is Go?",
			"answer":  "Go is a statically typed, compiled programming language",
			"snippet": "Designed at Google",
		},
		"related_searches": []any{
			map[string]any{"query": "golang tutorial"},
			map[string]any{"query": "golang vs python"},
		},
	}

	result := &SearchResult{
		Data: serpAPIData,
		Raw:  `{"organic_results":[...]}`,
	}

	normalizer := NewNormalizer("serpapi")
	normalized, err := normalizer.NormalizeSearch(result, "golang programming")

	if err != nil {
		t.Fatalf("NormalizeSearch failed: %v", err)
	}

	// Verify organic results
	if len(normalized.OrganicResults) != 2 {
		t.Errorf("Expected 2 organic results, got %d", len(normalized.OrganicResults))
	}

	if normalized.OrganicResults[0].Title != "Go Programming Language" {
		t.Errorf("Expected title 'Go Programming Language', got '%s'", normalized.OrganicResults[0].Title)
	}

	// Verify answer box
	if normalized.AnswerBox == nil {
		t.Fatal("Expected answer box to be present")
	}

	// Verify metadata
	if normalized.SearchMetadata.Engine != "serpapi" {
		t.Errorf("Expected engine 'serpapi', got '%s'", normalized.SearchMetadata.Engine)
	}
}

func TestNormalizeNews(t *testing.T) {
	// Mock Serper news response
	serperNews := map[string]any{
		"news": []any{
			map[string]any{
				"title":    "Breaking: Go 1.21 Released",
				"link":     "https://blog.golang.org/go1.21",
				"source":   "Go Blog",
				"date":     "2023-08-08",
				"snippet":  "The Go team is happy to announce...",
				"imageUrl": "https://example.com/image.jpg",
			},
		},
	}

	result := &SearchResult{
		Data: serperNews,
	}

	normalizer := NewNormalizer("serper")
	normalized, err := normalizer.NormalizeNews(result, "golang news")

	if err != nil {
		t.Fatalf("NormalizeNews failed: %v", err)
	}

	if len(normalized.NewsResults) != 1 {
		t.Errorf("Expected 1 news result, got %d", len(normalized.NewsResults))
	}

	if normalized.NewsResults[0].Title != "Breaking: Go 1.21 Released" {
		t.Errorf("Expected title 'Breaking: Go 1.21 Released', got '%s'", normalized.NewsResults[0].Title)
	}

	if normalized.NewsResults[0].Source != "Go Blog" {
		t.Errorf("Expected source 'Go Blog', got '%s'", normalized.NewsResults[0].Source)
	}
}

func TestNormalizeImages(t *testing.T) {
	// Mock Serper images response
	serperImages := map[string]any{
		"images": []any{
			map[string]any{
				"title":    "Gopher Mascot",
				"imageUrl": "https://example.com/gopher.png",
				"link":     "https://golang.org",
				"source":   "golang.org",
			},
		},
	}

	result := &SearchResult{
		Data: serperImages,
	}

	normalizer := NewNormalizer("serper")
	normalized, err := normalizer.NormalizeImages(result, "golang gopher")

	if err != nil {
		t.Fatalf("NormalizeImages failed: %v", err)
	}

	if len(normalized.ImageResults) != 1 {
		t.Errorf("Expected 1 image result, got %d", len(normalized.ImageResults))
	}

	if normalized.ImageResults[0].Title != "Gopher Mascot" {
		t.Errorf("Expected title 'Gopher Mascot', got '%s'", normalized.ImageResults[0].Title)
	}

	if normalized.ImageResults[0].ImageURL != "https://example.com/gopher.png" {
		t.Errorf("Expected imageUrl 'https://example.com/gopher.png', got '%s'", normalized.ImageResults[0].ImageURL)
	}
}

func TestNormalizerUnifiedStructure(t *testing.T) {
	// This test demonstrates that both Serper and SerpAPI produce the same normalized structure

	serperData := map[string]any{
		"organic": []any{
			map[string]any{
				"title":   "Test Result",
				"link":    "https://example.com",
				"snippet": "Test snippet",
			},
		},
	}

	serpAPIData := map[string]any{
		"organic_results": []any{
			map[string]any{
				"title":   "Test Result",
				"link":    "https://example.com",
				"snippet": "Test snippet",
			},
		},
	}

	serperResult := &SearchResult{Data: serperData}
	serpAPIResult := &SearchResult{Data: serpAPIData}

	serperNormalizer := NewNormalizer("serper")
	serpAPINormalizer := NewNormalizer("serpapi")

	serperNormalized, _ := serperNormalizer.NormalizeSearch(serperResult, "test")
	serpAPINormalized, _ := serpAPINormalizer.NormalizeSearch(serpAPIResult, "test")

	// Both should have the same structure and content (except engine name)
	if len(serperNormalized.OrganicResults) != len(serpAPINormalized.OrganicResults) {
		t.Errorf("Normalized results have different lengths")
	}

	if serperNormalized.OrganicResults[0].Title != serpAPINormalized.OrganicResults[0].Title {
		t.Errorf("Normalized titles don't match")
	}

	if serperNormalized.OrganicResults[0].Link != serpAPINormalized.OrganicResults[0].Link {
		t.Errorf("Normalized links don't match")
	}
}
