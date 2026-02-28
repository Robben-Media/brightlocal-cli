package brightlocal

// LocationSearchRequest is the request body for searching locations.
type LocationSearchRequest struct {
	Query   string `json:"query"`
	Country string `json:"country,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}

// Location represents a BrightLocal location result.
type Location struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Country string `json:"country,omitempty"`
	Zip     string `json:"zip,omitempty"`
}

// LocationSearchResponse is the response for location search.
type LocationSearchResponse struct {
	TotalCount int        `json:"total_count"`
	Items      []Location `json:"items"`
}

// RankingsCheckRequest is the request body for checking rankings.
type RankingsCheckRequest struct {
	BusinessName string   `json:"business_name"`
	Location     string   `json:"location"`
	SearchTerms  []string `json:"search_terms"`
}

// RankingResult represents a single ranking result.
type RankingResult struct {
	SearchTerm string `json:"search_term"`
	Rank       int    `json:"rank"`
	URL        string `json:"url,omitempty"`
	Source     string `json:"source,omitempty"`
}

// RankingsCheckResponse is the response for rankings check.
type RankingsCheckResponse struct {
	Success   bool            `json:"success"`
	RequestID string          `json:"request_id"`
	Results   []RankingResult `json:"results,omitempty"`
}

// RankingsGetResponse is the response for getting a rankings report.
type RankingsGetResponse struct {
	Success   bool            `json:"success"`
	RequestID string          `json:"request_id"`
	Status    string          `json:"status"`
	Results   []RankingResult `json:"results,omitempty"`
}
