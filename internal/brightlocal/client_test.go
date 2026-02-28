package brightlocal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/builtbyrobben/brightlocal-cli/internal/api"
)

func newTestClient(server *httptest.Server) *Client {
	return &Client{
		Client: api.NewClient("test-key",
			api.WithBaseURL(server.URL),
			api.WithUserAgent("brightlocal-cli/test"),
			api.WithAuthHeader("x-api-key"),
		),
	}
}

func TestLocations_Search(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/locations/search" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var body LocationSearchRequest
		json.NewDecoder(r.Body).Decode(&body)

		if body.Query != "Columbia, MO" {
			t.Errorf("expected query 'Columbia, MO', got %q", body.Query)
		}

		if body.Country != "USA" {
			t.Errorf("expected country 'USA', got %q", body.Country)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LocationSearchResponse{
			TotalCount: 1,
			Items: []Location{
				{ID: "loc-1", Name: "Columbia", City: "Columbia", State: "MO", Country: "USA"},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(server)

	result, err := client.Locations().Search(context.Background(), LocationSearchRequest{
		Query:   "Columbia, MO",
		Country: "USA",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Items) != 1 {
		t.Fatalf("expected 1 location, got %d", len(result.Items))
	}

	if result.Items[0].ID != "loc-1" {
		t.Errorf("expected location ID 'loc-1', got %q", result.Items[0].ID)
	}
}

func TestLocations_Search_EmptyQuery(t *testing.T) {
	client := &Client{Client: api.NewClient("key")}

	_, err := client.Locations().Search(context.Background(), LocationSearchRequest{})
	if err == nil {
		t.Fatal("expected error for empty query")
	}
}

func TestLocations_Search_DefaultLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body LocationSearchRequest
		json.NewDecoder(r.Body).Decode(&body)

		if body.Limit != 10 {
			t.Errorf("expected default limit 10, got %d", body.Limit)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LocationSearchResponse{TotalCount: 0, Items: []Location{}})
	}))
	defer server.Close()

	client := newTestClient(server)

	_, err := client.Locations().Search(context.Background(), LocationSearchRequest{
		Query: "test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRankings_Check(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/rankings/search" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var body RankingsCheckRequest
		json.NewDecoder(r.Body).Decode(&body)

		if body.BusinessName != "Test Biz" {
			t.Errorf("expected business 'Test Biz', got %q", body.BusinessName)
		}

		if len(body.SearchTerms) != 2 {
			t.Errorf("expected 2 search terms, got %d", len(body.SearchTerms))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(RankingsCheckResponse{
			Success:   true,
			RequestID: "abc-123",
		})
	}))
	defer server.Close()

	client := newTestClient(server)

	result, err := client.Rankings().Check(context.Background(), RankingsCheckRequest{
		BusinessName: "Test Biz",
		Location:     "Columbia, MO",
		SearchTerms:  []string{"plumber", "hvac"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.RequestID != "abc-123" {
		t.Errorf("expected request ID 'abc-123', got %q", result.RequestID)
	}
}

func TestRankings_Check_MissingBusiness(t *testing.T) {
	client := &Client{Client: api.NewClient("key")}

	_, err := client.Rankings().Check(context.Background(), RankingsCheckRequest{
		Location:    "Columbia, MO",
		SearchTerms: []string{"plumber"},
	})
	if err == nil {
		t.Fatal("expected error for missing business name")
	}
}

func TestRankings_Check_MissingLocation(t *testing.T) {
	client := &Client{Client: api.NewClient("key")}

	_, err := client.Rankings().Check(context.Background(), RankingsCheckRequest{
		BusinessName: "Test",
		SearchTerms:  []string{"plumber"},
	})
	if err == nil {
		t.Fatal("expected error for missing location")
	}
}

func TestRankings_Check_MissingTerms(t *testing.T) {
	client := &Client{Client: api.NewClient("key")}

	_, err := client.Rankings().Check(context.Background(), RankingsCheckRequest{
		BusinessName: "Test",
		Location:     "Columbia, MO",
	})
	if err == nil {
		t.Fatal("expected error for missing search terms")
	}
}

func TestRankings_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rankings/results/abc-123" {
			t.Errorf("expected path /rankings/results/abc-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(RankingsGetResponse{
			Success:   true,
			RequestID: "abc-123",
			Status:    "completed",
			Results: []RankingResult{
				{SearchTerm: "plumber", Rank: 3},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(server)

	result, err := client.Rankings().Get(context.Background(), "abc-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "completed" {
		t.Errorf("expected status 'completed', got %q", result.Status)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}

	if result.Results[0].Rank != 3 {
		t.Errorf("expected rank 3, got %d", result.Results[0].Rank)
	}
}

func TestRankings_Get_EmptyID(t *testing.T) {
	client := &Client{Client: api.NewClient("key")}

	_, err := client.Rankings().Get(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty request ID")
	}
}
