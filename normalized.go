package metasearch

// NormalizedSearchResult represents a unified search result structure across all engines
type NormalizedSearchResult struct {
	// Organic search results
	OrganicResults []OrganicResult `json:"organic_results,omitempty"`

	// Featured content
	AnswerBox      *AnswerBox      `json:"answer_box,omitempty"`
	KnowledgeGraph *KnowledgeGraph `json:"knowledge_graph,omitempty"`

	// Additional result types
	RelatedSearches []RelatedSearch `json:"related_searches,omitempty"`
	PeopleAlsoAsk   []PeopleAlsoAsk `json:"people_also_ask,omitempty"`

	// News-specific (for SearchNews)
	NewsResults []NewsResult `json:"news_results,omitempty"`

	// Image-specific (for SearchImages)
	ImageResults []ImageResult `json:"image_results,omitempty"`

	// Video-specific (for SearchVideos)
	VideoResults []VideoResult `json:"video_results,omitempty"`

	// Places-specific (for SearchPlaces)
	PlaceResults []PlaceResult `json:"place_results,omitempty"`

	// Shopping-specific (for SearchShopping)
	ShoppingResults []ShoppingResult `json:"shopping_results,omitempty"`

	// Scholar-specific (for SearchScholar)
	ScholarResults []ScholarResult `json:"scholar_results,omitempty"`

	// Autocomplete-specific (for SearchAutocomplete)
	Suggestions []string `json:"suggestions,omitempty"`

	// Metadata
	SearchMetadata SearchMetadata `json:"search_metadata"`

	// Original response (for debugging or fallback)
	Raw *SearchResult `json:"raw,omitempty"`
}

// OrganicResult represents a standard web search result
type OrganicResult struct {
	Position int    `json:"position"`
	Title    string `json:"title"`
	Link     string `json:"link"`
	URL      string `json:"url"` // Alias for Link
	Snippet  string `json:"snippet"`
	Domain   string `json:"domain,omitempty"`
	Date     string `json:"date,omitempty"`
}

// AnswerBox represents a featured answer at the top of results
type AnswerBox struct {
	Type    string `json:"type,omitempty"`
	Title   string `json:"title,omitempty"`
	Answer  string `json:"answer,omitempty"`
	Snippet string `json:"snippet,omitempty"`
	Source  string `json:"source,omitempty"`
	Link    string `json:"link,omitempty"`
}

// KnowledgeGraph represents a knowledge panel
type KnowledgeGraph struct {
	Title       string            `json:"title,omitempty"`
	Type        string            `json:"type,omitempty"`
	Description string            `json:"description,omitempty"`
	Source      string            `json:"source,omitempty"`
	ImageURL    string            `json:"image_url,omitempty"`
	Attributes  map[string]string `json:"attributes,omitempty"`
}

// RelatedSearch represents a related search suggestion
type RelatedSearch struct {
	Query string `json:"query"`
	Link  string `json:"link,omitempty"`
}

// PeopleAlsoAsk represents a "People Also Ask" question
type PeopleAlsoAsk struct {
	Question string `json:"question"`
	Answer   string `json:"answer,omitempty"`
	Title    string `json:"title,omitempty"`
	Link     string `json:"link,omitempty"`
	Source   string `json:"source,omitempty"`
}

// NewsResult represents a news article result
type NewsResult struct {
	Position  int    `json:"position"`
	Title     string `json:"title"`
	Link      string `json:"link"`
	Source    string `json:"source"`
	Date      string `json:"date,omitempty"`
	Snippet   string `json:"snippet,omitempty"`
	ImageURL  string `json:"image_url,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

// ImageResult represents an image search result
type ImageResult struct {
	Position   int    `json:"position"`
	Title      string `json:"title,omitempty"`
	ImageURL   string `json:"image_url"`
	Thumbnail  string `json:"thumbnail,omitempty"`
	Source     string `json:"source,omitempty"`
	SourceURL  string `json:"source_url,omitempty"`
	Width      int    `json:"width,omitempty"`
	Height     int    `json:"height,omitempty"`
	IsProduct  bool   `json:"is_product,omitempty"`
}

// VideoResult represents a video search result
type VideoResult struct {
	Position  int    `json:"position"`
	Title     string `json:"title"`
	Link      string `json:"link"`
	Channel   string `json:"channel,omitempty"`
	Platform  string `json:"platform,omitempty"` // e.g., "youtube", "vimeo"
	Duration  string `json:"duration,omitempty"`
	Date      string `json:"date,omitempty"`
	Views     string `json:"views,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Snippet   string `json:"snippet,omitempty"`
}

// PlaceResult represents a local business or place result
type PlaceResult struct {
	Position    int               `json:"position"`
	Title       string            `json:"title"`
	PlaceID     string            `json:"place_id,omitempty"`
	DataID      string            `json:"data_id,omitempty"`
	Address     string            `json:"address,omitempty"`
	Phone       string            `json:"phone,omitempty"`
	Website     string            `json:"website,omitempty"`
	Rating      float64           `json:"rating,omitempty"`
	Reviews     int               `json:"reviews,omitempty"`
	Type        string            `json:"type,omitempty"`
	Hours       string            `json:"hours,omitempty"`
	Price       string            `json:"price,omitempty"`
	Latitude    float64           `json:"latitude,omitempty"`
	Longitude   float64           `json:"longitude,omitempty"`
	Thumbnail   string            `json:"thumbnail,omitempty"`
	Attributes  map[string]string `json:"attributes,omitempty"`
}

// ShoppingResult represents a shopping/product result
type ShoppingResult struct {
	Position      int      `json:"position"`
	Title         string   `json:"title"`
	Link          string   `json:"link"`
	ProductID     string   `json:"product_id,omitempty"`
	Price         string   `json:"price,omitempty"`
	OriginalPrice string   `json:"original_price,omitempty"`
	Currency      string   `json:"currency,omitempty"`
	Rating        float64  `json:"rating,omitempty"`
	Reviews       int      `json:"reviews,omitempty"`
	Source        string   `json:"source,omitempty"`
	Delivery      string   `json:"delivery,omitempty"`
	Thumbnail     string   `json:"thumbnail,omitempty"`
	Images        []string `json:"images,omitempty"`
	InStock       bool     `json:"in_stock,omitempty"`
}

// ScholarResult represents a scholarly article result
type ScholarResult struct {
	Position       int      `json:"position"`
	Title          string   `json:"title"`
	Link           string   `json:"link"`
	PublicationURL string   `json:"publication_url,omitempty"`
	Authors        []string `json:"authors,omitempty"`
	Year           string   `json:"year,omitempty"`
	Source         string   `json:"source,omitempty"` // Journal/Conference name
	Citations      int      `json:"citations,omitempty"`
	Snippet        string   `json:"snippet,omitempty"`
	PDF            string   `json:"pdf,omitempty"`
}

// SearchMetadata contains metadata about the search itself
type SearchMetadata struct {
	Engine     string  `json:"engine"`      // "serper", "serpapi", etc.
	Query      string  `json:"query"`
	Location   string  `json:"location,omitempty"`
	Language   string  `json:"language,omitempty"`
	Country    string  `json:"country,omitempty"`
	TotalResults int64 `json:"total_results,omitempty"`
	TimeTaken  float64 `json:"time_taken,omitempty"` // seconds
}
