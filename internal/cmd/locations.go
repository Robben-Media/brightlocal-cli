package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/builtbyrobben/brightlocal-cli/internal/brightlocal"
	"github.com/builtbyrobben/brightlocal-cli/internal/outfmt"
)

type LocationsCmd struct {
	Search LocationsSearchCmd `cmd:"" help:"Search for locations"`
}

type LocationsSearchCmd struct {
	Query   string `required:"" help:"Search query (e.g. 'Columbia, MO')"`
	Country string `help:"Country code (e.g. USA, GBR)" default:"USA"`
	Limit   int    `help:"Maximum number of results" default:"10"`
}

func (cmd *LocationsSearchCmd) Run(ctx context.Context) error {
	client, err := getBrightLocalClient()
	if err != nil {
		return err
	}

	result, err := client.Locations().Search(ctx, brightlocal.LocationSearchRequest{
		Query:   cmd.Query,
		Country: cmd.Country,
		Limit:   cmd.Limit,
	})
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	if len(result.Locations) == 0 {
		fmt.Fprintln(os.Stderr, "No locations found")
		return nil
	}

	fmt.Fprintf(os.Stderr, "Found %d locations\n\n", len(result.Locations))
	for _, loc := range result.Locations {
		fmt.Printf("ID: %s\n", loc.ID)
		if loc.Name != "" {
			fmt.Printf("  Name: %s\n", loc.Name)
		}
		if loc.Address != "" {
			fmt.Printf("  Address: %s\n", loc.Address)
		}
		if loc.City != "" {
			fmt.Printf("  City: %s\n", loc.City)
		}
		if loc.State != "" {
			fmt.Printf("  State: %s\n", loc.State)
		}
		if loc.Country != "" {
			fmt.Printf("  Country: %s\n", loc.Country)
		}
		fmt.Println()
	}

	return nil
}
