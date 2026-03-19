package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/builtbyrobben/brightlocal-cli/internal/outfmt"
)

type ClientsCmd struct {
	List ClientsListCmd `cmd:"" help:"List all clients"`
	Get  ClientsGetCmd  `cmd:"" help:"Get client details"`
}

// --- clients list ---

type ClientsListCmd struct{}

func (cmd *ClientsListCmd) Run(ctx context.Context) error {
	client, err := getBrightLocalClient()
	if err != nil {
		return err
	}

	resp, err := client.Clients().List(ctx)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, resp)
	}
	if outfmt.IsPlain(ctx) {
		headers := []string{"CLIENT_ID", "NAME", "TYPE", "WEBSITE_URL"}
		var rows [][]string
		for _, c := range resp.Items {
			rows = append(rows, []string{
				fmt.Sprintf("%d", c.ClientID),
				c.Name,
				c.Type,
				c.WebsiteURL,
			})
		}
		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	if len(resp.Items) == 0 {
		fmt.Fprintln(os.Stderr, "No clients found")
		return nil
	}

	fmt.Fprintf(os.Stderr, "Found %d clients\n\n", resp.TotalCount)
	for _, c := range resp.Items {
		fmt.Printf("Client ID: %d\n", c.ClientID)
		fmt.Printf("  Name: %s\n", c.Name)
		fmt.Printf("  Type: %s\n", c.Type)
		if c.WebsiteURL != "" {
			fmt.Printf("  Website: %s\n", c.WebsiteURL)
		}
		fmt.Println()
	}

	return nil
}

// --- clients get ---

type ClientsGetCmd struct {
	ID string `arg:"" required:"" help:"Client ID"`
}

func (cmd *ClientsGetCmd) Run(ctx context.Context) error {
	client, err := getBrightLocalClient()
	if err != nil {
		return err
	}

	blClient, err := client.Clients().Get(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, blClient)
	}
	if outfmt.IsPlain(ctx) {
		headers := []string{"CLIENT_ID", "NAME", "TYPE", "WEBSITE_URL"}
		rows := [][]string{{
			fmt.Sprintf("%d", blClient.ClientID),
			blClient.Name,
			blClient.Type,
			blClient.WebsiteURL,
		}}
		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	fmt.Printf("Client ID: %d\n", blClient.ClientID)
	fmt.Printf("Name: %s\n", blClient.Name)
	fmt.Printf("Type: %s\n", blClient.Type)
	if blClient.WebsiteURL != "" {
		fmt.Printf("Website: %s\n", blClient.WebsiteURL)
	}

	return nil
}
