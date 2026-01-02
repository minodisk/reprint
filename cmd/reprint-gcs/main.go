package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const appName = "reprint-gcs"

var (
	bucket      string
	prefix      string
	credentials string
	mime        string
	objectID    string
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
	rootCmd.PersistentFlags().StringVar(&credentials, "credentials", "", "Service account key file path")

	// Upload flags
	uploadCmd.Flags().StringVar(&mime, "mime", "", "Image MIME type")

	// Delete flags
	deleteCmd.Flags().StringVar(&objectID, "object-id", "", "Object ID to delete")

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
