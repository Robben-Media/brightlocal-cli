package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/builtbyrobben/brightlocal-cli/internal/brightlocal"
	"github.com/builtbyrobben/brightlocal-cli/internal/outfmt"
)

type CitationsCmd struct {
	Audit CitationsAuditCmd `cmd:"" help:"Run a citation audit for a business"`
}

type CitationsAuditCmd struct {
	Business string `required:"" help:"Business name"`
	Location string `required:"" help:"Location (e.g. 'Columbia, MO')"`
}

func (cmd *CitationsAuditCmd) Run(ctx context.Context) error {
	client, err := getBrightLocalClient()
	if err != nil {
		return err
	}

	result, err := client.Citations().Audit(ctx, brightlocal.CitationAuditRequest{
		BusinessName: cmd.Business,
		Location:     cmd.Location,
	})
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	fmt.Fprintf(os.Stderr, "Citation audit submitted\n\n")
	fmt.Printf("Report ID: %d\n", result.ReportID)

	if len(result.Citations) > 0 {
		fmt.Println()
		for _, c := range result.Citations {
			fmt.Printf("Directory: %s\n", c.Directory)
			if c.URL != "" {
				fmt.Printf("  URL: %s\n", c.URL)
			}
			if c.Status != "" {
				fmt.Printf("  Status: %s\n", c.Status)
			}
			fmt.Printf("  NAP Match: %v\n", c.NAPMatch)
		}
	}

	return nil
}
