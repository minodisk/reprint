package main

import (
	"context"
	"fmt"
	"os"

	"github.com/minodisk/reprint/internal/config"
	"github.com/minodisk/reprint/internal/gcs"
	"github.com/spf13/cobra"
)

func runDoctor(cmd *cobra.Command, args []string) error {
	fmt.Println("Checking reprint-gcs configuration...")
	fmt.Println()

	allOK := true

	// Check 1: Load config
	fmt.Print("[Config] Loading configuration... ")
	cfg, err := config.Load(
		config.WithAppName(appName),
		config.WithBucket(bucket),
		config.WithPrefix(prefix),
		config.WithCredentials(credentials),
	)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return nil
	}
	fmt.Println("OK")

	// Check 2: Bucket configured
	fmt.Print("[Config] Bucket configured... ")
	if cfg.Bucket == "" {
		fmt.Println("ERROR: bucket is not configured")
		fmt.Println("  Set via: --bucket, REPRINT_BUCKET, or ~/.config/reprint/config.yaml")
		allOK = false
	} else {
		fmt.Printf("OK (%s)\n", cfg.Bucket)
	}

	// Check 3: Prefix (optional)
	fmt.Print("[Config] Prefix configured... ")
	if cfg.Prefix == "" {
		fmt.Println("(not set)")
	} else {
		fmt.Printf("OK (%s)\n", cfg.Prefix)
	}

	// Check 4: Credentials configured
	fmt.Print("[Auth] Credentials configured... ")
	defaultCredPath := config.DefaultCredentialsPath(appName)
	if cfg.Credentials == "" {
		fmt.Println("ERROR: credentials is not configured")
		fmt.Println("  Set via:")
		fmt.Println("    - --credentials flag")
		fmt.Println("    - REPRINT_CREDENTIALS environment variable")
		fmt.Println("    - credentials in ~/.config/reprint/config.yaml")
		fmt.Printf("    - Place file at %s\n", defaultCredPath)
		allOK = false
	} else if cfg.Credentials == defaultCredPath {
		fmt.Printf("OK (using default: %s)\n", cfg.Credentials)
	} else {
		fmt.Printf("OK (%s)\n", cfg.Credentials)
	}

	// Check 5: Credentials file exists
	if cfg.Credentials != "" {
		fmt.Print("[Auth] Credentials file exists... ")
		if _, err := os.Stat(cfg.Credentials); os.IsNotExist(err) {
			fmt.Printf("ERROR: file not found: %s\n", cfg.Credentials)
			allOK = false
		} else if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			allOK = false
		} else {
			fmt.Println("OK")
		}
	}

	// Check 6: GCS connection (only if bucket and credentials are configured)
	if cfg.Bucket != "" && cfg.Credentials != "" {
		fmt.Print("[GCS] Connecting to GCS... ")
		ctx := context.Background()
		client, err := gcs.NewClient(ctx, cfg.Bucket, cfg.Prefix, cfg.Credentials)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			allOK = false
		} else {
			defer client.Close()
			fmt.Println("OK")

			// Check 7: Bucket access
			fmt.Print("[GCS] Checking bucket access... ")
			if err := client.CheckBucket(ctx); err != nil {
				fmt.Printf("ERROR: %v\n", err)
				allOK = false
			} else {
				fmt.Println("OK")
			}
		}
	}

	fmt.Println()
	if allOK {
		fmt.Println("All checks passed!")
	} else {
		fmt.Println("Some checks failed. Please fix the issues above.")
	}

	return nil
}
