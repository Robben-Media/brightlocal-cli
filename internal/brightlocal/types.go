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
	Success   bool       `json:"success"`
	Locations []Location `json:"locations"`
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
	Success  bool            `json:"success"`
	ReportID int             `json:"report_id"`
	Results  []RankingResult `json:"results,omitempty"`
}

// RankingsGetResponse is the response for getting a rankings report.
type RankingsGetResponse struct {
	Success  bool            `json:"success"`
	ReportID int             `json:"report_id"`
	Status   string          `json:"status"`
	Results  []RankingResult `json:"results,omitempty"`
}

// CitationAuditRequest is the request body for a citation audit.
type CitationAuditRequest struct {
	BusinessName string `json:"business_name"`
	Location     string `json:"location"`
}

// Citation represents a single citation entry.
type Citation struct {
	Directory string `json:"directory"`
	URL       string `json:"url,omitempty"`
	Status    string `json:"status,omitempty"`
	NAPMatch  bool   `json:"nap_match"`
}

// CitationAuditResponse is the response for a citation audit.
type CitationAuditResponse struct {
	Success   bool       `json:"success"`
	ReportID  int        `json:"report_id"`
	Citations []Citation `json:"citations,omitempty"`
}

// Report represents a BrightLocal report.
type Report struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Status    string `json:"status,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

// ReportsListResponse is the paginated response for listing reports.
type ReportsListResponse struct {
	Success    bool     `json:"success"`
	Reports    []Report `json:"reports"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
	TotalItems int      `json:"total_items"`
}

// ReportCreateRequest is the request body for creating a report.
type ReportCreateRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// ReportCreateResponse is the response for creating a report.
type ReportCreateResponse struct {
	Success  bool `json:"success"`
	ReportID int  `json:"report_id"`
}
