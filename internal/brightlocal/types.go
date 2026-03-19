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

// --- Management API types ---

// LRTReport represents a Local Rank Tracker report (list item).
type LRTReport struct {
	ID            int    `json:"report_id"`
	LocationID    int    `json:"location_id"`
	Name          string `json:"report_name"`
	Schedule      string `json:"schedule"`
	RunDayOfWeek  string `json:"run_day_of_week"`
	RunDayOfMonth *int   `json:"run_day_of_month"`
}

// LRTReportsListResponse is the response for listing LRT reports.
type LRTReportsListResponse struct {
	TotalCount int         `json:"total_count"`
	Items      []LRTReport `json:"items"`
}

// LRTReportDetail represents a full LRT report (single report view).
type LRTReportDetail struct {
	ID              int      `json:"report_id"`
	CustomerID      int      `json:"customer_id"`
	LocationID      int      `json:"location_id"`
	Name            string   `json:"report_name"`
	SearchEngines   []string `json:"search_engines"`
	Keywords        []string `json:"keywords"`
	Country         string   `json:"country"`
	SearchLocation  string   `json:"search_location"`
	BusinessNames   []string `json:"business_names"`
	Postcode        string   `json:"postcode"`
	Telephone       string   `json:"telephone"`
	WebsiteAddresses []string `json:"website_addresses"`
	Schedule        string   `json:"schedule"`
	RunDayOfWeek    string   `json:"run_day_of_week"`
	RunDayOfMonth   *int     `json:"run_day_of_month"`
	RunTime         string   `json:"run_time"`
	RunTimeZone     string   `json:"run_time_zone"`
	IsRunning       bool     `json:"is_running"`
	IsPublic        bool     `json:"is_public"`
	CreatedAt       string   `json:"created_at"`
	LastProcessedAt string   `json:"last_processed_at"`
}

// LRTRankingEntry represents a single ranking entry from an LRT result.
type LRTRankingEntry struct {
	ID             int    `json:"id"`
	Keyword        string `json:"keyword"`
	SearchEngine   string `json:"search_engine"`
	Rank           int    `json:"rank"`
	UnblendedRank  int    `json:"unblended_rank"`
	Page           int    `json:"page"`
	Type           string `json:"type"`
	Match          string `json:"match"`
	Title          string `json:"title"`
	URL            string `json:"url"`
	OrigURL        string `json:"orig_url"`
	Date           string `json:"date"`
	Last           int    `json:"last"`
	LastUnblended  int    `json:"last_unblended"`
	LastDate       string `json:"last_date"`
	SubRank        int    `json:"sub_rank"`
}

// LRTResultURLs contains report URL links.
type LRTResultURLs struct {
	PDFURL         string `json:"pdf_url"`
	CSVURL         string `json:"csv_url"`
	InteractiveURL string `json:"interactive_url"`
}

// LRTResultResponse is the response for LRT report results.
// The actual API returns: {"urls": {...}, "rankings": {"by_keyword": [...]}}
type LRTResultResponse struct {
	URLs     LRTResultURLs `json:"urls"`
	Rankings struct {
		ByKeyword []LRTKeywordResultRaw `json:"by_keyword"`
	} `json:"rankings"`
}

// LRTKeywordResultRaw matches the actual API structure where results
// is a map of search engine to ranking entries.
type LRTKeywordResultRaw struct {
	Keyword string                          `json:"keyword"`
	Results map[string][]LRTRankingEntry    `json:"results"`
}

// LRTHistoryEntry represents a single history run entry.
type LRTHistoryEntry struct {
	ReportHistoryID int    `json:"report_history_id"`
	ReportID        int    `json:"report_id"`
	LocationID      int    `json:"location_id"`
	HistoryType     string `json:"history_type"`
	GenerationDate  string `json:"generation_date"`
}

// LRTHistoryResponse is the response for LRT report history.
type LRTHistoryResponse struct {
	TotalCount int                `json:"total_count"`
	Items      []LRTHistoryEntry  `json:"items"`
}

// BLClient represents a BrightLocal client.
type BLClient struct {
	ClientID        int    `json:"client_id"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	UniqueReference string `json:"unique_reference"`
	WebsiteURL      string `json:"website_url"`
}

// ClientsListResponse is the response for listing clients.
type ClientsListResponse struct {
	TotalCount int        `json:"total_count"`
	Items      []BLClient `json:"items"`
}
