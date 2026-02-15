package cmd

import (
	"fmt"
	"os"

	"github.com/builtbyrobben/brightlocal-cli/internal/brightlocal"
	"github.com/builtbyrobben/brightlocal-cli/internal/secrets"
)

func getBrightLocalClient() (*brightlocal.Client, error) {
	// Check for environment variable override first
	apiKey := os.Getenv("BRIGHTLOCAL_API_KEY")

	if apiKey == "" {
		// Try to get from keyring
		store, err := secrets.OpenDefault()
		if err != nil {
			return nil, fmt.Errorf("open credential store: %w", err)
		}

		apiKey, err = store.GetAPIKey()
		if err != nil {
			return nil, fmt.Errorf("get API key: %w (set BRIGHTLOCAL_API_KEY or run 'brightlocal-cli auth set-key --stdin')", err)
		}
	}

	return brightlocal.NewClient(apiKey), nil
}
