package gcs

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

const (
	// DefaultSignedURLExpiration is the default expiration time for signed URLs.
	DefaultSignedURLExpiration = 15 * time.Minute
)

// Client wraps the GCS client.
type Client struct {
	client   *storage.Client
	bucket   string
	prefix   string
	endpoint string // custom endpoint for emulator
}

// NewClient creates a new GCS client.
// credentials must be a path to a service account key file (required for signed URLs).
func NewClient(ctx context.Context, bucket, prefix, credentials string) (*Client, error) {
	return NewClientWithEndpoint(ctx, bucket, prefix, credentials, "")
}

// NewClientWithEndpoint creates a new GCS client with a custom endpoint.
// This is useful for testing with emulators like fake-gcs-server.
func NewClientWithEndpoint(ctx context.Context, bucket, prefix, credentials, endpoint string) (*Client, error) {
	var opts []option.ClientOption
	if credentials != "" {
		opts = append(opts, option.WithCredentialsFile(credentials))
	}
	if endpoint != "" {
		opts = append(opts, option.WithEndpoint(endpoint), option.WithoutAuthentication())
	}

	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	return &Client{
		client:   client,
		bucket:   bucket,
		prefix:   prefix,
		endpoint: endpoint,
	}, nil
}

// Close closes the GCS client.
func (c *Client) Close() error {
	return c.client.Close()
}

// Upload uploads data to GCS and returns a signed URL.
func (c *Client) Upload(ctx context.Context, filename string, data io.Reader, contentType string) (string, error) {
	objectName := c.objectName(filename)
	obj := c.client.Bucket(c.bucket).Object(objectName)

	w := obj.NewWriter(ctx)
	w.ContentType = contentType

	if _, err := io.Copy(w, data); err != nil {
		return "", fmt.Errorf("failed to write to GCS: %w", err)
	}

	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close GCS writer: %w", err)
	}

	return c.SignedURL(filename, DefaultSignedURLExpiration)
}

// SignedURL returns a signed URL for an object with the specified expiration.
// Requires a service account key file to be configured via credentials.
func (c *Client) SignedURL(filename string, expiration time.Duration) (string, error) {
	// For emulator, return public URL (signed URLs don't work with emulator)
	if c.endpoint != "" {
		return c.PublicURL(filename), nil
	}

	objectName := c.objectName(filename)
	opts := &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(expiration),
	}

	url, err := c.client.Bucket(c.bucket).SignedURL(objectName, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return url, nil
}

// Delete deletes an object from GCS.
func (c *Client) Delete(ctx context.Context, filename string) error {
	objectName := c.objectName(filename)
	obj := c.client.Bucket(c.bucket).Object(objectName)

	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete from GCS: %w", err)
	}

	return nil
}

// CheckBucket checks if the bucket exists and is accessible.
func (c *Client) CheckBucket(ctx context.Context) error {
	_, err := c.client.Bucket(c.bucket).Attrs(ctx)
	if err != nil {
		return fmt.Errorf("failed to access bucket %q: %w", c.bucket, err)
	}
	return nil
}

// PublicURL returns the public URL for an object.
func (c *Client) PublicURL(filename string) string {
	objectName := c.objectName(filename)
	if c.endpoint != "" {
		// For emulator, use the endpoint URL
		return fmt.Sprintf("http://localhost:4443/%s/%s", c.bucket, objectName)
	}
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.bucket, objectName)
}

// objectName returns the full object name with prefix.
func (c *Client) objectName(filename string) string {
	if c.prefix == "" {
		return filename
	}
	return c.prefix + filename
}
