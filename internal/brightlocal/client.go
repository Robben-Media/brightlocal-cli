package brightlocal

import (
	"context"
	"errors"
	"fmt"

	"github.com/builtbyrobben/brightlocal-cli/internal/api"
)

var (
	errQueryRequired        = errors.New("query is required")
	errBusinessNameRequired = errors.New("business name is required")
	errLocationRequired     = errors.New("location is required")
	errSearchTermsRequired  = errors.New("search terms are required")
	errReportIDRequired     = errors.New("report ID is required")
	errReportNameRequired   = errors.New("report name is required")
	errReportTypeRequired   = errors.New("report type is required")
)

const defaultBaseURL = "https://api.brightlocal.com/data/v1"

// Client wraps the API client with BrightLocal-specific methods.
type Client struct {
	*api.Client
}

// NewClient creates a new BrightLocal API client.
func NewClient(apiKey string) *Client {
	return &Client{
		Client: api.NewClient(apiKey,
			api.WithBaseURL(defaultBaseURL),
			api.WithUserAgent("brightlocal-cli/1.0"),
			api.WithAuthHeader("x-api-key"),
		),
	}
}

// Locations provides methods for the Locations API.
func (c *Client) Locations() *LocationsService {
	return &LocationsService{client: c}
}

// Rankings provides methods for the Rankings API.
func (c *Client) Rankings() *RankingsService {
	return &RankingsService{client: c}
}

// Citations provides methods for the Citations API.
func (c *Client) Citations() *CitationsService {
	return &CitationsService{client: c}
}

// Reports provides methods for the Reports API.
func (c *Client) Reports() *ReportsService {
	return &ReportsService{client: c}
}

// LocationsService handles location operations.
type LocationsService struct {
	client *Client
}

// Search searches for locations by query.
func (s *LocationsService) Search(ctx context.Context, req LocationSearchRequest) (*LocationSearchResponse, error) {
	if req.Query == "" {
		return nil, errQueryRequired
	}

	if req.Limit < 1 {
		req.Limit = 10
	}

	var result LocationSearchResponse
	if err := s.client.Post(ctx, "/locations/search", req, &result); err != nil {
		return nil, fmt.Errorf("search locations: %w", err)
	}

	return &result, nil
}

// RankingsService handles ranking operations.
type RankingsService struct {
	client *Client
}

// Check submits a rankings check request.
func (s *RankingsService) Check(ctx context.Context, req RankingsCheckRequest) (*RankingsCheckResponse, error) {
	if req.BusinessName == "" {
		return nil, errBusinessNameRequired
	}

	if req.Location == "" {
		return nil, errLocationRequired
	}

	if len(req.SearchTerms) == 0 {
		return nil, errSearchTermsRequired
	}

	var result RankingsCheckResponse
	if err := s.client.Post(ctx, "/rankings/check", req, &result); err != nil {
		return nil, fmt.Errorf("check rankings: %w", err)
	}

	return &result, nil
}

// Get retrieves a rankings report by ID.
func (s *RankingsService) Get(ctx context.Context, reportID int) (*RankingsGetResponse, error) {
	if reportID == 0 {
		return nil, errReportIDRequired
	}

	var result RankingsGetResponse

	path := fmt.Sprintf("/rankings/%d", reportID)
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("get rankings: %w", err)
	}

	return &result, nil
}

// CitationsService handles citation operations.
type CitationsService struct {
	client *Client
}

// Audit submits a citation audit request.
func (s *CitationsService) Audit(ctx context.Context, req CitationAuditRequest) (*CitationAuditResponse, error) {
	if req.BusinessName == "" {
		return nil, errBusinessNameRequired
	}

	if req.Location == "" {
		return nil, errLocationRequired
	}

	var result CitationAuditResponse
	if err := s.client.Post(ctx, "/citations/audit", req, &result); err != nil {
		return nil, fmt.Errorf("audit citations: %w", err)
	}

	return &result, nil
}

// ReportsService handles report operations.
type ReportsService struct {
	client *Client
}

// List returns all reports (paginated).
func (s *ReportsService) List(ctx context.Context, page, pageSize int) (*ReportsListResponse, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	path := fmt.Sprintf("/reports?page=%d&page_size=%d", page, pageSize)

	var result ReportsListResponse
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("list reports: %w", err)
	}

	return &result, nil
}

// Create creates a new report.
func (s *ReportsService) Create(ctx context.Context, req ReportCreateRequest) (*ReportCreateResponse, error) {
	if req.Name == "" {
		return nil, errReportNameRequired
	}

	if req.Type == "" {
		return nil, errReportTypeRequired
	}

	var result ReportCreateResponse
	if err := s.client.Post(ctx, "/reports", req, &result); err != nil {
		return nil, fmt.Errorf("create report: %w", err)
	}

	return &result, nil
}
