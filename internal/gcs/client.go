package gcs

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// Client wraps the GCS client.
type Client struct {
	client   *storage.Client
	bucket   string
	prefix   string
	endpoint string // custom endpoint for emulator
}

// NewClient creates a new GCS client.
// If credentials is provided, it will be used for authentication.
// Otherwise, default credentials will be used.
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

// Upload uploads data to GCS and returns the public URL.
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

	return c.PublicURL(filename), nil
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
