package gcs

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

const (
	testBucket   = "test-bucket"
	testEndpoint = "http://localhost:4443/storage/v1/"
)

func TestIntegration_UploadAndDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Wait for emulator to be ready
	if err := waitForEmulator(30 * time.Second); err != nil {
		t.Fatalf("emulator not ready: %v", err)
	}

	ctx := context.Background()

	// Create bucket via HTTP API
	if err := createBucket(ctx, testBucket); err != nil {
		t.Fatalf("failed to create bucket: %v", err)
	}

	// Create client with emulator endpoint
	client, err := NewClientWithEndpoint(ctx, testBucket, "test-prefix/", "", testEndpoint)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	// Test Upload
	testData := []byte("test image data")
	filename := "test-file-123"
	contentType := "image/png"

	url, err := client.Upload(ctx, filename, bytes.NewReader(testData), contentType)
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}

	expectedURL := "http://localhost:4443/" + testBucket + "/test-prefix/" + filename
	if url != expectedURL {
		t.Errorf("Upload() URL = %q, want %q", url, expectedURL)
	}

	// Test Delete
	if err := client.Delete(ctx, filename); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify object is deleted (should fail to delete again)
	if err := client.Delete(ctx, filename); err == nil {
		t.Error("Delete() should fail for non-existent object")
	}
}

func waitForEmulator(timeout time.Duration) error {
	client := &http.Client{Timeout: 1 * time.Second}
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		resp, err := client.Get("http://localhost:4443/storage/v1/b")
		if err == nil {
			resp.Body.Close()
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("emulator not available after %v", timeout)
}

func createBucket(ctx context.Context, bucket string) error {
	client, err := storage.NewClient(ctx,
		option.WithEndpoint(testEndpoint),
		option.WithoutAuthentication(),
	)
	if err != nil {
		return err
	}
	defer client.Close()

	// Ignore error if bucket already exists
	_ = client.Bucket(bucket).Create(ctx, "test-project", nil)
	return nil
}
