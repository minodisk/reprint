package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/minodisk/reprint/internal/config"
	"github.com/minodisk/reprint/internal/gcs"
	"github.com/spf13/cobra"
)

func runDoctor(cmd *cobra.Command, args []string) error {
	fmt.Println("Checking reprint-gcs configuration...")
	fmt.Println()

	ctx := context.Background()
	allOK := true

	cfg, ok := checkConfig()
	if !ok {
		allOK = false
	}

	var client *gcs.Client
	if cfg != nil && cfg.Bucket != "" && cfg.Credentials != "" {
		var ok bool
		client, ok = checkGCSConnection(ctx, cfg)
		if !ok {
			allOK = false
		}
		if client != nil {
			defer client.Close()
		}
	}

	if client != nil {
		if !checkBucketAccess(ctx, client) {
			allOK = false
		}

		objectID, ok := checkUploadPermission(ctx, client)
		if !ok {
			allOK = false
		}

		if objectID != "" {
			if !checkDeletePermission(ctx, client, objectID) {
				allOK = false
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

func checkConfig() (*config.Config, bool) {
	allOK := true

	fmt.Print("[Config] Loading configuration... ")
	cfg, err := config.Load(
		config.WithAppName(appName),
		config.WithBucket(bucket),
		config.WithPrefix(prefix),
		config.WithCredentials(credentials),
	)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return nil, false
	}
	fmt.Println("OK")

	fmt.Print("[Config] Bucket configured... ")
	if cfg.Bucket == "" {
		fmt.Println("ERROR: bucket is not configured")
		fmt.Println("  Set via: --bucket, REPRINT_BUCKET, or ~/.config/reprint/config.yaml")
		allOK = false
	} else {
		fmt.Printf("OK (%s)\n", cfg.Bucket)
	}

	fmt.Print("[Config] Prefix configured... ")
	if cfg.Prefix == "" {
		fmt.Println("(not set)")
	} else {
		fmt.Printf("OK (%s)\n", cfg.Prefix)
	}

	defaultCredPath := config.DefaultCredentialsPath(appName)
	fmt.Print("[Auth] Credentials configured... ")
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

	return cfg, allOK
}

func checkGCSConnection(ctx context.Context, cfg *config.Config) (*gcs.Client, bool) {
	fmt.Print("[GCS] Connecting to GCS... ")
	client, err := gcs.NewClient(ctx, cfg.Bucket, cfg.Prefix, cfg.Credentials)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return nil, false
	}
	fmt.Println("OK")
	return client, true
}

func checkBucketAccess(ctx context.Context, client *gcs.Client) bool {
	fmt.Print("[GCS] Checking bucket access... ")
	if err := client.CheckBucket(ctx); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		fmt.Println("  Required permission: storage.buckets.get")
		fmt.Println("  Recommended role: roles/storage.bucketViewer")
		return false
	}
	fmt.Println("OK")
	return true
}

func checkUploadPermission(ctx context.Context, client *gcs.Client) (string, bool) {
	testObjectID := ".reprint-doctor-test-" + uuid.New().String()
	testData := strings.NewReader("reprint-gcs doctor test")

	fmt.Print("[GCS] Testing upload permission... ")
	_, err := client.Upload(ctx, testObjectID, testData, "text/plain")
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		fmt.Println("  Required permission: storage.objects.create, storage.objects.get")
		fmt.Println("  Recommended role: roles/storage.objectAdmin")
		return "", false
	}
	fmt.Println("OK")
	return testObjectID, true
}

func checkDeletePermission(ctx context.Context, client *gcs.Client, objectID string) bool {
	fmt.Print("[GCS] Testing delete permission... ")
	if err := client.Delete(ctx, objectID); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		fmt.Println("  Required permission: storage.objects.delete")
		fmt.Println("  Recommended role: roles/storage.objectAdmin")
		fmt.Printf("  Note: Test object %q was left in the bucket\n", objectID)
		return false
	}
	fmt.Println("OK")
	return true
}
