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
	Get   RankingsGetCmd   `cmd:"" help:"Get a rankings report by ID"`
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

	fmt.Fprintf(os.Stderr, "Rankings check submitted\n\n")
	fmt.Printf("Report ID: %d\n", result.ReportID)

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
	ReportID int `arg:"" required:"" help:"Report ID"`
}

func (cmd *RankingsGetCmd) Run(ctx context.Context) error {
	client, err := getBrightLocalClient()
	if err != nil {
		return err
	}

	result, err := client.Rankings().Get(ctx, cmd.ReportID)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	fmt.Printf("Report ID: %d\n", result.ReportID)
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
