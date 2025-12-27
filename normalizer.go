package omniserp

import (
	"fmt"
	"strings"
)

// Normalizer converts engine-specific responses to normalized format
type Normalizer struct {
	engineName string
}

// NewNormalizer creates a new normalizer for the specified engine
func NewNormalizer(engineName string) *Normalizer {
	return &Normalizer{engineName: strings.ToLower(engineName)}
}

// NormalizeSearch normalizes a web search result
func (n *Normalizer) NormalizeSearch(result *SearchResult, query string) (*NormalizedSearchResult, error) {
	if result == nil || result.Data == nil {
		return nil, fmt.Errorf("nil result or data")
	}

	data, ok := result.Data.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unexpected data type: %T", result.Data)
	}

	normalized := &NormalizedSearchResult{
		SearchMetadata: SearchMetadata{
			Engine: n.engineName,
			Query:  query,
		},
		Raw: result,
	}

	switch n.engineName {
	case "serper":
		n.normalizeSerperSearch(data, normalized)
	case "serpapi":
		n.normalizeSerpAPISearch(data, normalized)
	default:
		return nil, fmt.Errorf("unsupported engine: %s", n.engineName)
	}

	return normalized, nil
}

// NormalizeNews normalizes a news search result
func (n *Normalizer) NormalizeNews(result *SearchResult, query string) (*NormalizedSearchResult, error) {
	if result == nil || result.Data == nil {
		return nil, fmt.Errorf("nil result or data")
	}

	data, ok := result.Data.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unexpected data type: %T", result.Data)
	}

	normalized := &NormalizedSearchResult{
		SearchMetadata: SearchMetadata{
			Engine: n.engineName,
			Query:  query,
		},
		Raw: result,
	}

	switch n.engineName {
	case "serper":
		n.normalizeSerperNews(data, normalized)
	case "serpapi":
		n.normalizeSerpAPINews(data, normalized)
	default:
		return nil, fmt.Errorf("unsupported engine: %s", n.engineName)
	}

	return normalized, nil
}

// NormalizeImages normalizes an image search result
func (n *Normalizer) NormalizeImages(result *SearchResult, query string) (*NormalizedSearchResult, error) {
	if result == nil || result.Data == nil {
		return nil, fmt.Errorf("nil result or data")
	}

	data, ok := result.Data.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unexpected data type: %T", result.Data)
	}

	normalized := &NormalizedSearchResult{
		SearchMetadata: SearchMetadata{
			Engine: n.engineName,
			Query:  query,
		},
		Raw: result,
	}

	switch n.engineName {
	case "serper":
		n.normalizeSerperImages(data, normalized)
	case "serpapi":
		n.normalizeSerpAPIImages(data, normalized)
	default:
		return nil, fmt.Errorf("unsupported engine: %s", n.engineName)
	}

	return normalized, nil
}

// Helper functions for Serper normalization

func (n *Normalizer) normalizeSerperSearch(data map[string]any, normalized *NormalizedSearchResult) {
	// Extract organic results
	if organic, ok := data["organic"].([]any); ok {
		for i, item := range organic {
			if itemMap, ok := item.(map[string]any); ok {
				normalized.OrganicResults = append(normalized.OrganicResults, OrganicResult{
					Position: i + 1,
					Title:    getString(itemMap, "title"),
					Link:     getString(itemMap, "link"),
					URL:      getString(itemMap, "link"),
					Snippet:  getString(itemMap, "snippet"),
					Date:     getString(itemMap, "date"),
				})
			}
		}
	}

	// Extract answer box
	if answerBox, ok := data["answerBox"].(map[string]any); ok {
		normalized.AnswerBox = &AnswerBox{
			Type:    getString(answerBox, "type"),
			Title:   getString(answerBox, "title"),
			Answer:  getString(answerBox, "answer"),
			Snippet: getString(answerBox, "snippet"),
			Source:  getString(answerBox, "source"),
			Link:    getString(answerBox, "link"),
		}
	}

	// Extract knowledge graph
	if kg, ok := data["knowledgeGraph"].(map[string]any); ok {
		normalized.KnowledgeGraph = &KnowledgeGraph{
			Title:       getString(kg, "title"),
			Type:        getString(kg, "type"),
			Description: getString(kg, "description"),
			ImageURL:    getString(kg, "imageUrl"),
		}
	}

	// Extract related searches
	if related, ok := data["relatedSearches"].([]any); ok {
		for _, item := range related {
			if itemMap, ok := item.(map[string]any); ok {
				normalized.RelatedSearches = append(normalized.RelatedSearches, RelatedSearch{
					Query: getString(itemMap, "query"),
				})
			}
		}
	}

	// Extract people also ask
	if paa, ok := data["peopleAlsoAsk"].([]any); ok {
		for _, item := range paa {
			if itemMap, ok := item.(map[string]any); ok {
				normalized.PeopleAlsoAsk = append(normalized.PeopleAlsoAsk, PeopleAlsoAsk{
					Question: getString(itemMap, "question"),
					Answer:   getString(itemMap, "answer"),
					Title:    getString(itemMap, "title"),
					Link:     getString(itemMap, "link"),
				})
			}
		}
	}

	// Extract search metadata
	if searchParams, ok := data["searchParameters"].(map[string]any); ok {
		normalized.SearchMetadata.Query = getString(searchParams, "q")
		normalized.SearchMetadata.Location = getString(searchParams, "location")
		normalized.SearchMetadata.Language = getString(searchParams, "hl")
		normalized.SearchMetadata.Country = getString(searchParams, "gl")
	}
}

func (n *Normalizer) normalizeSerperNews(data map[string]any, normalized *NormalizedSearchResult) {
	if news, ok := data["news"].([]any); ok {
		for i, item := range news {
			if itemMap, ok := item.(map[string]any); ok {
				normalized.NewsResults = append(normalized.NewsResults, NewsResult{
					Position:  i + 1,
					Title:     getString(itemMap, "title"),
					Link:      getString(itemMap, "link"),
					Source:    getString(itemMap, "source"),
					Date:      getString(itemMap, "date"),
					Snippet:   getString(itemMap, "snippet"),
					ImageURL:  getString(itemMap, "imageUrl"),
					Thumbnail: getString(itemMap, "imageUrl"),
				})
			}
		}
	}
}

func (n *Normalizer) normalizeSerperImages(data map[string]any, normalized *NormalizedSearchResult) {
	if images, ok := data["images"].([]any); ok {
		for i, item := range images {
			if itemMap, ok := item.(map[string]any); ok {
				normalized.ImageResults = append(normalized.ImageResults, ImageResult{
					Position:  i + 1,
					Title:     getString(itemMap, "title"),
					ImageURL:  getString(itemMap, "imageUrl"),
					Thumbnail: getString(itemMap, "imageUrl"),
					Source:    getString(itemMap, "source"),
					SourceURL: getString(itemMap, "link"),
				})
			}
		}
	}
}

// Helper functions for SerpAPI normalization

func (n *Normalizer) normalizeSerpAPISearch(data map[string]any, normalized *NormalizedSearchResult) {
	// Extract organic results
	if organic, ok := data["organic_results"].([]any); ok {
		for i, item := range organic {
			if itemMap, ok := item.(map[string]any); ok {
				normalized.OrganicResults = append(normalized.OrganicResults, OrganicResult{
					Position: i + 1,
					Title:    getString(itemMap, "title"),
					Link:     getString(itemMap, "link"),
					URL:      getString(itemMap, "link"),
					Snippet:  getString(itemMap, "snippet"),
					Date:     getString(itemMap, "date"),
				})
			}
		}
	}

	// Extract answer box
	if answerBox, ok := data["answer_box"].(map[string]any); ok {
		normalized.AnswerBox = &AnswerBox{
			Type:    getString(answerBox, "type"),
			Title:   getString(answerBox, "title"),
			Answer:  getString(answerBox, "answer"),
			Snippet: getString(answerBox, "snippet"),
			Link:    getString(answerBox, "link"),
		}
	}

	// Extract knowledge graph
	if kg, ok := data["knowledge_graph"].(map[string]any); ok {
		normalized.KnowledgeGraph = &KnowledgeGraph{
			Title:       getString(kg, "title"),
			Type:        getString(kg, "type"),
			Description: getString(kg, "description"),
			ImageURL:    getString(kg, "image"),
		}
	}

	// Extract related searches
	if related, ok := data["related_searches"].([]any); ok {
		for _, item := range related {
			if itemMap, ok := item.(map[string]any); ok {
				normalized.RelatedSearches = append(normalized.RelatedSearches, RelatedSearch{
					Query: getString(itemMap, "query"),
					Link:  getString(itemMap, "link"),
				})
			}
		}
	}

	// Extract people also ask
	if paa, ok := data["related_questions"].([]any); ok {
		for _, item := range paa {
			if itemMap, ok := item.(map[string]any); ok {
				normalized.PeopleAlsoAsk = append(normalized.PeopleAlsoAsk, PeopleAlsoAsk{
					Question: getString(itemMap, "question"),
					Answer:   getString(itemMap, "answer"),
					Title:    getString(itemMap, "title"),
					Link:     getString(itemMap, "link"),
					Source:   getString(itemMap, "displayed_link"),
				})
			}
		}
	}

	// Extract search metadata
	if searchParams, ok := data["search_parameters"].(map[string]any); ok {
		normalized.SearchMetadata.Query = getString(searchParams, "q")
		normalized.SearchMetadata.Location = getString(searchParams, "location")
		normalized.SearchMetadata.Language = getString(searchParams, "hl")
		normalized.SearchMetadata.Country = getString(searchParams, "gl")
	}
}

func (n *Normalizer) normalizeSerpAPINews(data map[string]any, normalized *NormalizedSearchResult) {
	if news, ok := data["news_results"].([]any); ok {
		for i, item := range news {
			if itemMap, ok := item.(map[string]any); ok {
				normalized.NewsResults = append(normalized.NewsResults, NewsResult{
					Position: i + 1,
					Title:    getString(itemMap, "title"),
					Link:     getString(itemMap, "link"),
					Source:   getString(itemMap, "source"),
					Date:     getString(itemMap, "date"),
					Snippet:  getString(itemMap, "snippet"),
				})

				// SerpAPI may have thumbnail object
				if thumb, ok := itemMap["thumbnail"].(string); ok {
					normalized.NewsResults[len(normalized.NewsResults)-1].Thumbnail = thumb
				}
			}
		}
	}
}

func (n *Normalizer) normalizeSerpAPIImages(data map[string]any, normalized *NormalizedSearchResult) {
	if images, ok := data["images_results"].([]any); ok {
		for i, item := range images {
			if itemMap, ok := item.(map[string]any); ok {
				normalized.ImageResults = append(normalized.ImageResults, ImageResult{
					Position:  i + 1,
					Title:     getString(itemMap, "title"),
					ImageURL:  getString(itemMap, "original"),
					Thumbnail: getString(itemMap, "thumbnail"),
					Source:    getString(itemMap, "source"),
					SourceURL: getString(itemMap, "link"),
				})
			}
		}
	}
}

// Helper function to safely extract string values from maps
func getString(m map[string]any, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
