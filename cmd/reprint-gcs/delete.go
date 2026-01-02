package main

import (
	"context"
	"fmt"
	"os"

	"github.com/minodisk/reprint/internal/gcs"
	"github.com/spf13/cobra"
)

func runDelete(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	// Get object ID from flag or environment variable
	if objectID == "" {
		objectID = os.Getenv("DECK_DELETE_ID")
	}
	if objectID == "" {
		return fmt.Errorf("object-id is required (--object-id or DECK_DELETE_ID)")
	}

	ctx := context.Background()
	client, err := gcs.NewClient(ctx, cfg.Bucket, cfg.Prefix, cfg.Credentials)
	if err != nil {
		return err
	}
	defer client.Close()

	return client.Delete(ctx, objectID)
}
