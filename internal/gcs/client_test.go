package gcs

import "testing"

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
