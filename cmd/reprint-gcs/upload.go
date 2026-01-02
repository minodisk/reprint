package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/minodisk/reprint/internal/gcs"
	"github.com/spf13/cobra"
)

func runUpload(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	// Get MIME type from flag or environment variable
	if mime == "" {
		mime = os.Getenv("DECK_UPLOAD_MIME")
	}
	if mime == "" {
		return fmt.Errorf("MIME type is required (--mime or DECK_UPLOAD_MIME)")
	}

	ctx := context.Background()
	client, err := gcs.NewClient(ctx, cfg.Bucket, cfg.Prefix, cfg.Credentials)
	if err != nil {
		return err
	}
	defer client.Close()

	// Generate UUID filename
	filename := uuid.New().String()

	// Read from stdin and upload
	url, err := client.Upload(ctx, filename, os.Stdin, mime)
	if err != nil {
		return err
	}

	// Output URL and filename
	fmt.Println(url)
	fmt.Println(filename)

	return nil
}
