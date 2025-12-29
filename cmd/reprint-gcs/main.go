package main

import (
	"context"
	"fmt"
	"io"
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
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "reprint-gcs",
	Short: "External image uploader CLI for deck using Google Cloud Storage",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&bucket, "bucket", "", "GCS bucket name")
	rootCmd.PersistentFlags().StringVar(&prefix, "prefix", "", "Object prefix")
	rootCmd.PersistentFlags().StringVar(&credentials, "credentials", "", "Service account key file path")

	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(deleteCmd)
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload image to GCS",
	RunE:  runUpload,
}

var mime string

func init() {
	uploadCmd.Flags().StringVar(&mime, "mime", "", "Image MIME type")
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

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete image from GCS",
	RunE:  runDelete,
}

var filename string

func init() {
	deleteCmd.Flags().StringVar(&filename, "filename", "", "Filename to delete")
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
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Override with CLI flags
	if bucket != "" {
		cfg.Bucket = bucket
	}
	if prefix != "" {
		cfg.Prefix = prefix
	}
	if credentials != "" {
		cfg.Credentials = credentials
	}

	// Validate required fields
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("bucket is required (--bucket, REPRINT_BUCKET, or config file)")
	}

	return cfg, nil
}

// Ensure stdin is not used elsewhere
var _ io.Reader = os.Stdin
