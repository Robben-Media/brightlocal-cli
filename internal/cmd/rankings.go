package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/builtbyrobben/brightlocal-cli/internal/brightlocal"
	"github.com/builtbyrobben/brightlocal-cli/internal/outfmt"
)

type RankingsCmd struct {
	Check RankingsCheckCmd `cmd:"" help:"Check rankings for a business"`
	Get   RankingsGetCmd   `cmd:"" help:"Get rankings results by request ID"`
}

type RankingsCheckCmd struct {
	Business string `required:"" help:"Business name"`
	Location string `required:"" help:"Location (e.g. 'Columbia, MO')"`
	Terms    string `required:"" help:"Comma-separated search terms"`
}

func (cmd *RankingsCheckCmd) Run(ctx context.Context) error {
	client, err := getBrightLocalClient()
	if err != nil {
		return err
	}

	var terms []string
	for _, term := range strings.Split(cmd.Terms, ",") {
		t := strings.TrimSpace(term)
		if t != "" {
			terms = append(terms, t)
		}
	}

	result, err := client.Rankings().Check(ctx, brightlocal.RankingsCheckRequest{
		BusinessName: cmd.Business,
		Location:     cmd.Location,
		SearchTerms:  terms,
	})
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}
	if outfmt.IsPlain(ctx) {
		headers := []string{"REQUEST_ID", "TERM", "RANK", "URL", "SOURCE"}
		var rows [][]string
		for _, r := range result.Results {
			rows = append(rows, []string{result.RequestID, r.SearchTerm, fmt.Sprintf("%d", r.Rank), r.URL, r.Source})
		}
		if len(rows) == 0 {
			rows = append(rows, []string{result.RequestID, "", "", "", ""})
		}
		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	fmt.Fprintf(os.Stderr, "Rankings check submitted\n\n")
	fmt.Printf("Request ID: %s\n", result.RequestID)

	if len(result.Results) > 0 {
		fmt.Println()
		for _, r := range result.Results {
			fmt.Printf("Term: %s\n", r.SearchTerm)
			fmt.Printf("  Rank: %d\n", r.Rank)
			if r.URL != "" {
				fmt.Printf("  URL: %s\n", r.URL)
			}
			if r.Source != "" {
				fmt.Printf("  Source: %s\n", r.Source)
			}
		}
	}

	return nil
}

type RankingsGetCmd struct {
	RequestID string `arg:"" required:"" help:"Request ID (UUID)"`
}

func (cmd *RankingsGetCmd) Run(ctx context.Context) error {
	client, err := getBrightLocalClient()
	if err != nil {
		return err
	}

	result, err := client.Rankings().Get(ctx, cmd.RequestID)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}
	if outfmt.IsPlain(ctx) {
		headers := []string{"REQUEST_ID", "STATUS", "TERM", "RANK", "URL", "SOURCE"}
		var rows [][]string
		for _, r := range result.Results {
			rows = append(rows, []string{result.RequestID, result.Status, r.SearchTerm, fmt.Sprintf("%d", r.Rank), r.URL, r.Source})
		}
		if len(rows) == 0 {
			rows = append(rows, []string{result.RequestID, result.Status, "", "", "", ""})
		}
		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	fmt.Printf("Request ID: %s\n", result.RequestID)
	fmt.Printf("Status: %s\n", result.Status)

	if len(result.Results) > 0 {
		fmt.Println()
		for _, r := range result.Results {
			fmt.Printf("Term: %s\n", r.SearchTerm)
			fmt.Printf("  Rank: %d\n", r.Rank)
			if r.URL != "" {
				fmt.Printf("  URL: %s\n", r.URL)
			}
			if r.Source != "" {
				fmt.Printf("  Source: %s\n", r.Source)
			}
		}
	}

	return nil
}
