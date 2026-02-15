package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/builtbyrobben/brightlocal-cli/internal/brightlocal"
	"github.com/builtbyrobben/brightlocal-cli/internal/outfmt"
)

type ReportsCmd struct {
	List   ReportsListCmd   `cmd:"" help:"List all reports"`
	Create ReportsCreateCmd `cmd:"" help:"Create a new report"`
}

type ReportsListCmd struct {
	Page     int `help:"Page number" default:"1"`
	PageSize int `help:"Page size" default:"10"`
}

func (cmd *ReportsListCmd) Run(ctx context.Context) error {
	client, err := getBrightLocalClient()
	if err != nil {
		return err
	}

	result, err := client.Reports().List(ctx, cmd.Page, cmd.PageSize)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	if len(result.Reports) == 0 {
		fmt.Fprintln(os.Stderr, "No reports found")
		return nil
	}

	fmt.Fprintf(os.Stderr, "Showing %d of %d reports (page %d)\n\n", len(result.Reports), result.TotalItems, result.Page)
	for _, report := range result.Reports {
		fmt.Printf("ID: %d\n", report.ID)
		if report.Name != "" {
			fmt.Printf("  Name: %s\n", report.Name)
		}
		if report.Type != "" {
			fmt.Printf("  Type: %s\n", report.Type)
		}
		if report.Status != "" {
			fmt.Printf("  Status: %s\n", report.Status)
		}
		if report.CreatedAt != "" {
			fmt.Printf("  Created: %s\n", report.CreatedAt)
		}
		fmt.Println()
	}

	return nil
}

type ReportsCreateCmd struct {
	Name string `required:"" help:"Report name"`
	Type string `required:"" help:"Report type (e.g. rankings, citations)"`
}

func (cmd *ReportsCreateCmd) Run(ctx context.Context) error {
	client, err := getBrightLocalClient()
	if err != nil {
		return err
	}

	result, err := client.Reports().Create(ctx, brightlocal.ReportCreateRequest{
		Name: cmd.Name,
		Type: cmd.Type,
	})
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	fmt.Fprintf(os.Stderr, "Report created\n\n")
	fmt.Printf("Report ID: %d\n", result.ReportID)

	return nil
}
