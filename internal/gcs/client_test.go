package gcs

import (
	"testing"
	"time"
)

func TestClient_objectName(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		filename string
		want     string
	}{
		{
			name:     "without prefix",
			prefix:   "",
			filename: "test-file",
			want:     "test-file",
		},
		{
			name:     "with prefix",
			prefix:   "images/",
			filename: "test-file",
			want:     "images/test-file",
		},
		{
			name:     "with prefix no trailing slash",
			prefix:   "images",
			filename: "test-file",
			want:     "imagestest-file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{prefix: tt.prefix}
			if got := c.objectName(tt.filename); got != tt.want {
				t.Errorf("objectName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestClient_PublicURL(t *testing.T) {
	tests := []struct {
		name     string
		bucket   string
		prefix   string
		filename string
		want     string
	}{
		{
			name:     "without prefix",
			bucket:   "my-bucket",
			prefix:   "",
			filename: "abc-123",
			want:     "https://storage.googleapis.com/my-bucket/abc-123",
		},
		{
			name:     "with prefix",
			bucket:   "my-bucket",
			prefix:   "deck/",
			filename: "abc-123",
			want:     "https://storage.googleapis.com/my-bucket/deck/abc-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{bucket: tt.bucket, prefix: tt.prefix}
			if got := c.PublicURL(tt.filename); got != tt.want {
				t.Errorf("PublicURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestClient_SignedURL_Emulator(t *testing.T) {
	// When using emulator, SignedURL should return public URL
	c := &Client{
		bucket:   "test-bucket",
		prefix:   "test-prefix/",
		endpoint: "http://localhost:4443/storage/v1/",
	}

	url, err := c.SignedURL("test-file", 15*time.Minute)
	if err != nil {
		t.Fatalf("SignedURL() error = %v", err)
	}

	want := "http://localhost:4443/test-bucket/test-prefix/test-file"
	if url != want {
		t.Errorf("SignedURL() = %q, want %q", url, want)
	}
}

func TestClient_PublicURL_Emulator(t *testing.T) {
	c := &Client{
		bucket:   "test-bucket",
		prefix:   "test-prefix/",
		endpoint: "http://localhost:4443/storage/v1/",
	}

	got := c.PublicURL("test-file")
	want := "http://localhost:4443/test-bucket/test-prefix/test-file"
	if got != want {
		t.Errorf("PublicURL() = %q, want %q", got, want)
	}
}
