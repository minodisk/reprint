package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/minodisk/reprint/internal/config"
	"github.com/minodisk/reprint/internal/gcs"
	"github.com/spf13/cobra"
)

var (
	bucket      string
	prefix      string
	credentials string
	mime        string
	filename    string
)

var rootCmd = &cobra.Command{
	Use:   "reprint-gcs",
	Short: "External image uploader CLI for deck using Google Cloud Storage",
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload image to GCS",
	RunE:  runUpload,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete image from GCS",
	RunE:  runDelete,
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose configuration and credentials",
	RunE:  runDoctor,
}

func init() {
	// Root flags
	rootCmd.PersistentFlags().StringVar(&bucket, "bucket", "", "GCS bucket name")
	rootCmd.PersistentFlags().StringVar(&prefix, "prefix", "", "Object prefix")
	rootCmd.PersistentFlags().StringVar(&credentials, "credentials", "", "Service account key file path (required)")

	// Upload flags
	uploadCmd.Flags().StringVar(&mime, "mime", "", "Image MIME type")

	// Delete flags
	deleteCmd.Flags().StringVar(&filename, "filename", "", "Filename to delete")

	// Add subcommands
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(doctorCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

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

func runDelete(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	// Get filename from flag or environment variable
	if filename == "" {
		filename = os.Getenv("DECK_DELETE_FILENAME")
	}
	if filename == "" {
		return fmt.Errorf("filename is required (--filename or DECK_DELETE_FILENAME)")
	}

	ctx := context.Background()
	client, err := gcs.NewClient(ctx, cfg.Bucket, cfg.Prefix, cfg.Credentials)
	if err != nil {
		return err
	}
	defer client.Close()

	return client.Delete(ctx, filename)
}

func loadConfig() (*config.Config, error) {
	cfg, err := config.Load(
		config.WithBucket(bucket),
		config.WithPrefix(prefix),
		config.WithCredentials(credentials),
	)
	if err != nil {
		return nil, err
	}

	// Validate required fields
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("bucket is required (--bucket, REPRINT_BUCKET, or config file)")
	}
	if cfg.Credentials == "" {
		return nil, fmt.Errorf("credentials is required (--credentials, REPRINT_CREDENTIALS, or config file)")
	}

	return cfg, nil
}

func runDoctor(cmd *cobra.Command, args []string) error {
	fmt.Println("Checking reprint-gcs configuration...")
	fmt.Println()

	allOK := true

	// Check 1: Load config
	fmt.Print("[Config] Loading configuration... ")
	cfg, err := config.Load(
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
	if cfg.Credentials == "" {
		fmt.Println("ERROR: credentials is not configured")
		fmt.Println("  Set via: --credentials, REPRINT_CREDENTIALS, or ~/.config/reprint/config.yaml")
		allOK = false
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
