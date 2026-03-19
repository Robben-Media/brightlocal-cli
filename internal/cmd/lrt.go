package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/builtbyrobben/brightlocal-cli/internal/outfmt"
)

type LRTCmd struct {
	Reports LRTReportsCmd `cmd:"" help:"LRT report operations"`
}

type LRTReportsCmd struct {
	List    LRTReportsListCmd    `cmd:"" help:"List all LRT reports"`
	Get     LRTReportsGetCmd     `cmd:"" help:"Get report details"`
	Result  LRTReportsResultCmd  `cmd:"" help:"Get ranking results for a report"`
	History LRTReportsHistoryCmd `cmd:"" help:"Get historical runs for a report"`
}

// --- lrt reports list ---

type LRTReportsListCmd struct{}

func (cmd *LRTReportsListCmd) Run(ctx context.Context) error {
	client, err := getBrightLocalClient()
	if err != nil {
		return err
	}

	resp, err := client.LRT().ListReports(ctx)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, resp)
	}
	if outfmt.IsPlain(ctx) {
		headers := []string{"REPORT_ID", "NAME", "LOCATION_ID", "SCHEDULE"}
		var rows [][]string
		for _, r := range resp.Items {
			rows = append(rows, []string{
				fmt.Sprintf("%d", r.ID),
				r.Name,
				fmt.Sprintf("%d", r.LocationID),
				r.Schedule,
			})
		}
		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	if len(resp.Items) == 0 {
		fmt.Fprintln(os.Stderr, "No LRT reports found")
		return nil
	}

	fmt.Fprintf(os.Stderr, "Found %d LRT reports\n\n", resp.TotalCount)
	for _, r := range resp.Items {
		fmt.Printf("Report ID: %d\n", r.ID)
		fmt.Printf("  Name: %s\n", r.Name)
		fmt.Printf("  Location ID: %d\n", r.LocationID)
		fmt.Printf("  Schedule: %s\n", r.Schedule)
		fmt.Println()
	}

	return nil
}

// --- lrt reports get ---

type LRTReportsGetCmd struct {
	ID string `arg:"" required:"" help:"Report ID"`
}

func (cmd *LRTReportsGetCmd) Run(ctx context.Context) error {
	client, err := getBrightLocalClient()
	if err != nil {
		return err
	}

	report, err := client.LRT().GetReport(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, report)
	}
	if outfmt.IsPlain(ctx) {
		headers := []string{"REPORT_ID", "NAME", "LOCATION_ID", "SCHEDULE", "LAST_PROCESSED"}
		rows := [][]string{{
			fmt.Sprintf("%d", report.ID),
			report.Name,
			fmt.Sprintf("%d", report.LocationID),
			report.Schedule,
			report.LastProcessedAt,
		}}
		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	fmt.Printf("Report ID: %d\n", report.ID)
	fmt.Printf("Name: %s\n", report.Name)
	fmt.Printf("Customer ID: %d\n", report.CustomerID)
	fmt.Printf("Location ID: %d\n", report.LocationID)
	fmt.Printf("Country: %s\n", report.Country)
	fmt.Printf("Search Location: %s\n", report.SearchLocation)
	if len(report.Keywords) > 0 {
		fmt.Printf("Keywords: %s\n", strings.Join(report.Keywords, ", "))
	}
	if len(report.SearchEngines) > 0 {
		fmt.Printf("Search Engines: %s\n", strings.Join(report.SearchEngines, ", "))
	}
	if len(report.BusinessNames) > 0 {
		fmt.Printf("Business Names: %s\n", strings.Join(report.BusinessNames, ", "))
	}
	fmt.Printf("Schedule: %s\n", report.Schedule)
	fmt.Printf("Last Processed: %s\n", report.LastProcessedAt)
	fmt.Printf("Created: %s\n", report.CreatedAt)
	fmt.Printf("Running: %v\n", report.IsRunning)

	return nil
}

// --- lrt reports result ---

type LRTReportsResultCmd struct {
	ID string `arg:"" required:"" help:"Report ID"`
}

func (cmd *LRTReportsResultCmd) Run(ctx context.Context) error {
	client, err := getBrightLocalClient()
	if err != nil {
		return err
	}

	result, err := client.LRT().GetResult(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	// Flatten the nested structure for plain/human output
	type flatEntry struct {
		keyword      string
		searchEngine string
		rank         int
		entryType    string
		last         int
		url          string
		date         string
	}

	var entries []flatEntry
	for _, kw := range result.Rankings.ByKeyword {
		for engine, rankings := range kw.Results {
			for _, r := range rankings {
				entries = append(entries, flatEntry{
					keyword:      r.Keyword,
					searchEngine: engine,
					rank:         r.Rank,
					entryType:    r.Type,
					last:         r.Last,
					url:          r.URL,
					date:         r.Date,
				})
			}
		}
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"KEYWORD", "SEARCH_ENGINE", "RANK", "TYPE", "LAST", "URL", "DATE"}
		var rows [][]string
		for _, e := range entries {
			rows = append(rows, []string{
				e.keyword,
				e.searchEngine,
				fmt.Sprintf("%d", e.rank),
				e.entryType,
				fmt.Sprintf("%d", e.last),
				e.url,
				e.date,
			})
		}
		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	if len(entries) == 0 {
		fmt.Fprintln(os.Stderr, "No ranking results found")
		return nil
	}

	fmt.Fprintf(os.Stderr, "%d ranking entries\n\n", len(entries))
	for _, e := range entries {
		fmt.Printf("Keyword: %s\n", e.keyword)
		fmt.Printf("  Search Engine: %s\n", e.searchEngine)
		fmt.Printf("  Rank: %d\n", e.rank)
		fmt.Printf("  Type: %s\n", e.entryType)
		if e.last != 0 {
			fmt.Printf("  Previous: %d\n", e.last)
		}
		if e.url != "" {
			fmt.Printf("  URL: %s\n", e.url)
		}
		if e.date != "" {
			fmt.Printf("  Date: %s\n", e.date)
		}
		fmt.Println()
	}

	return nil
}

// --- lrt reports history ---

type LRTReportsHistoryCmd struct {
	ID string `arg:"" required:"" help:"Report ID"`
}

func (cmd *LRTReportsHistoryCmd) Run(ctx context.Context) error {
	client, err := getBrightLocalClient()
	if err != nil {
		return err
	}

	result, err := client.LRT().GetHistory(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}
	if outfmt.IsPlain(ctx) {
		headers := []string{"HISTORY_ID", "REPORT_ID", "TYPE", "DATE"}
		var rows [][]string
		for _, h := range result.Items {
			rows = append(rows, []string{
				fmt.Sprintf("%d", h.ReportHistoryID),
				fmt.Sprintf("%d", h.ReportID),
				h.HistoryType,
				h.GenerationDate,
			})
		}
		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	if len(result.Items) == 0 {
		fmt.Fprintln(os.Stderr, "No history entries found")
		return nil
	}

	fmt.Fprintf(os.Stderr, "%d history entries\n\n", result.TotalCount)
	for _, h := range result.Items {
		fmt.Printf("History ID: %d  Type: %s  Date: %s\n", h.ReportHistoryID, h.HistoryType, h.GenerationDate)
	}

	return nil
}
