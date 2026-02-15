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
			Success: true,
			Locations: []Location{
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

	if len(result.Locations) != 1 {
		t.Fatalf("expected 1 location, got %d", len(result.Locations))
	}

	if result.Locations[0].ID != "loc-1" {
		t.Errorf("expected location ID 'loc-1', got %q", result.Locations[0].ID)
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
		json.NewEncoder(w).Encode(LocationSearchResponse{Success: true, Locations: []Location{}})
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

		if r.URL.Path != "/rankings/check" {
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
			Success:  true,
			ReportID: 42,
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

	if result.ReportID != 42 {
		t.Errorf("expected report ID 42, got %d", result.ReportID)
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
		if r.URL.Path != "/rankings/42" {
			t.Errorf("expected path /rankings/42, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(RankingsGetResponse{
			Success:  true,
			ReportID: 42,
			Status:   "completed",
			Results: []RankingResult{
				{SearchTerm: "plumber", Rank: 3},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(server)

	result, err := client.Rankings().Get(context.Background(), 42)
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

func TestRankings_Get_ZeroID(t *testing.T) {
	client := &Client{Client: api.NewClient("key")}

	_, err := client.Rankings().Get(context.Background(), 0)
	if err == nil {
		t.Fatal("expected error for zero report ID")
	}
}

func TestCitations_Audit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/citations/audit" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var body CitationAuditRequest
		json.NewDecoder(r.Body).Decode(&body)

		if body.BusinessName != "Test Biz" {
			t.Errorf("expected business 'Test Biz', got %q", body.BusinessName)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CitationAuditResponse{
			Success:  true,
			ReportID: 99,
			Citations: []Citation{
				{Directory: "Yelp", Status: "found", NAPMatch: true},
			},
		})
	}))
	defer server.Close()

	client := newTestClient(server)

	result, err := client.Citations().Audit(context.Background(), CitationAuditRequest{
		BusinessName: "Test Biz",
		Location:     "Columbia, MO",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ReportID != 99 {
		t.Errorf("expected report ID 99, got %d", result.ReportID)
	}

	if len(result.Citations) != 1 {
		t.Fatalf("expected 1 citation, got %d", len(result.Citations))
	}

	if result.Citations[0].Directory != "Yelp" {
		t.Errorf("expected directory 'Yelp', got %q", result.Citations[0].Directory)
	}
}

func TestCitations_Audit_MissingBusiness(t *testing.T) {
	client := &Client{Client: api.NewClient("key")}

	_, err := client.Citations().Audit(context.Background(), CitationAuditRequest{
		Location: "Columbia, MO",
	})
	if err == nil {
		t.Fatal("expected error for missing business name")
	}
}

func TestCitations_Audit_MissingLocation(t *testing.T) {
	client := &Client{Client: api.NewClient("key")}

	_, err := client.Citations().Audit(context.Background(), CitationAuditRequest{
		BusinessName: "Test",
	})
	if err == nil {
		t.Fatal("expected error for missing location")
	}
}

func TestReports_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/reports" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if r.URL.Query().Get("page") != "1" {
			t.Errorf("expected page=1, got %s", r.URL.Query().Get("page"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ReportsListResponse{
			Success: true,
			Reports: []Report{
				{ID: 1, Name: "Test Report", Type: "rankings", Status: "active"},
			},
			Page:       1,
			PageSize:   10,
			TotalItems: 1,
		})
	}))
	defer server.Close()

	client := newTestClient(server)

	result, err := client.Reports().List(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Reports) != 1 {
		t.Fatalf("expected 1 report, got %d", len(result.Reports))
	}

	if result.Reports[0].Name != "Test Report" {
		t.Errorf("expected name 'Test Report', got %q", result.Reports[0].Name)
	}
}

func TestReports_List_DefaultPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "1" {
			t.Errorf("expected page=1 default, got %s", r.URL.Query().Get("page"))
		}

		if r.URL.Query().Get("page_size") != "10" {
			t.Errorf("expected page_size=10 default, got %s", r.URL.Query().Get("page_size"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ReportsListResponse{Success: true, Reports: []Report{}})
	}))
	defer server.Close()

	client := newTestClient(server)

	_, err := client.Reports().List(context.Background(), 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReports_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body ReportCreateRequest
		json.NewDecoder(r.Body).Decode(&body)

		if body.Name != "My Report" {
			t.Errorf("expected name 'My Report', got %q", body.Name)
		}

		if body.Type != "rankings" {
			t.Errorf("expected type 'rankings', got %q", body.Type)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ReportCreateResponse{
			Success:  true,
			ReportID: 55,
		})
	}))
	defer server.Close()

	client := newTestClient(server)

	result, err := client.Reports().Create(context.Background(), ReportCreateRequest{
		Name: "My Report",
		Type: "rankings",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ReportID != 55 {
		t.Errorf("expected report ID 55, got %d", result.ReportID)
	}
}

func TestReports_Create_MissingName(t *testing.T) {
	client := &Client{Client: api.NewClient("key")}

	_, err := client.Reports().Create(context.Background(), ReportCreateRequest{
		Type: "rankings",
	})
	if err == nil {
		t.Fatal("expected error for missing report name")
	}
}

func TestReports_Create_MissingType(t *testing.T) {
	client := &Client{Client: api.NewClient("key")}

	_, err := client.Reports().Create(context.Background(), ReportCreateRequest{
		Name: "Test",
	})
	if err == nil {
		t.Fatal("expected error for missing report type")
	}
}
