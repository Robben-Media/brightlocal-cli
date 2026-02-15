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
