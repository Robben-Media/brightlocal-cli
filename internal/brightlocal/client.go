package brightlocal

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/builtbyrobben/brightlocal-cli/internal/api"
)

var (
	errQueryRequired        = errors.New("query is required")
	errBusinessNameRequired = errors.New("business name is required")
	errLocationRequired     = errors.New("location is required")
	errSearchTermsRequired  = errors.New("search terms are required")
	errRequestIDRequired    = errors.New("request ID is required")
)

const defaultBaseURL = "https://api.brightlocal.com/manage/v1"

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

// LRT provides methods for the Local Rank Tracker API.
func (c *Client) LRT() *LRTService {
	return &LRTService{client: c}
}

// Clients provides methods for the Clients API.
func (c *Client) Clients() *ClientsService {
	return &ClientsService{client: c}
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

// Check submits a rankings search request.
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
	if err := s.client.Post(ctx, "/rankings/search", req, &result); err != nil {
		return nil, fmt.Errorf("check rankings: %w", err)
	}

	return &result, nil
}

// Get retrieves rankings results by request ID.
func (s *RankingsService) Get(ctx context.Context, requestID string) (*RankingsGetResponse, error) {
	if requestID == "" {
		return nil, errRequestIDRequired
	}

	var result RankingsGetResponse

	path := fmt.Sprintf("/rankings/results/%s", url.PathEscape(requestID))
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("get rankings: %w", err)
	}

	return &result, nil
}

// LRTService handles Local Rank Tracker operations.
type LRTService struct {
	client *Client
}

// ListReports lists all LRT reports.
func (s *LRTService) ListReports(ctx context.Context) (*LRTReportsListResponse, error) {
	var result LRTReportsListResponse
	if err := s.client.Get(ctx, "/lrt/reports", &result); err != nil {
		return nil, fmt.Errorf("list LRT reports: %w", err)
	}
	return &result, nil
}

// GetReport retrieves a single LRT report by ID.
func (s *LRTService) GetReport(ctx context.Context, reportID string) (*LRTReportDetail, error) {
	if reportID == "" {
		return nil, errRequestIDRequired
	}
	var result LRTReportDetail
	path := fmt.Sprintf("/lrt/reports/%s", url.PathEscape(reportID))
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("get LRT report: %w", err)
	}
	return &result, nil
}

// GetResult retrieves ranking results for an LRT report.
func (s *LRTService) GetResult(ctx context.Context, reportID string) (*LRTResultResponse, error) {
	if reportID == "" {
		return nil, errRequestIDRequired
	}
	var result LRTResultResponse
	path := fmt.Sprintf("/lrt/reports/%s/result", url.PathEscape(reportID))
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("get LRT result: %w", err)
	}
	return &result, nil
}

// GetHistory retrieves historical data for an LRT report.
func (s *LRTService) GetHistory(ctx context.Context, reportID string) (*LRTHistoryResponse, error) {
	if reportID == "" {
		return nil, errRequestIDRequired
	}
	var result LRTHistoryResponse
	path := fmt.Sprintf("/lrt/reports/%s/history", url.PathEscape(reportID))
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("get LRT history: %w", err)
	}
	return &result, nil
}

// ClientsService handles client management operations.
type ClientsService struct {
	client *Client
}

// List retrieves all BrightLocal clients.
func (s *ClientsService) List(ctx context.Context) (*ClientsListResponse, error) {
	var result ClientsListResponse
	if err := s.client.Get(ctx, "/clients", &result); err != nil {
		return nil, fmt.Errorf("list clients: %w", err)
	}
	return &result, nil
}

// Get retrieves a single client by ID.
func (s *ClientsService) Get(ctx context.Context, clientID string) (*BLClient, error) {
	if clientID == "" {
		return nil, errors.New("client ID is required")
	}
	var result BLClient
	path := fmt.Sprintf("/clients/%s", url.PathEscape(clientID))
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("get client: %w", err)
	}
	return &result, nil
}
