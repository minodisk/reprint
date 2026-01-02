package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_FromEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("REPRINT_BUCKET", "test-bucket")
	os.Setenv("REPRINT_PREFIX", "test-prefix/")
	os.Setenv("REPRINT_CREDENTIALS", "/path/to/creds.json")
	defer func() {
		os.Unsetenv("REPRINT_BUCKET")
		os.Unsetenv("REPRINT_PREFIX")
		os.Unsetenv("REPRINT_CREDENTIALS")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Bucket != "test-bucket" {
		t.Errorf("Bucket = %q, want %q", cfg.Bucket, "test-bucket")
	}
	if cfg.Prefix != "test-prefix/" {
		t.Errorf("Prefix = %q, want %q", cfg.Prefix, "test-prefix/")
	}
	if cfg.Credentials != "/path/to/creds.json" {
		t.Errorf("Credentials = %q, want %q", cfg.Credentials, "/path/to/creds.json")
	}
}

func TestLoad_EmptyConfig(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("REPRINT_BUCKET")
	os.Unsetenv("REPRINT_PREFIX")
	os.Unsetenv("REPRINT_CREDENTIALS")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Bucket != "" {
		t.Errorf("Bucket = %q, want empty", cfg.Bucket)
	}
	if cfg.Prefix != "" {
		t.Errorf("Prefix = %q, want empty", cfg.Prefix)
	}
	if cfg.Credentials != "" {
		t.Errorf("Credentials = %q, want empty", cfg.Credentials)
	}
}

func TestLoad_WithOptions(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("REPRINT_BUCKET")
	os.Unsetenv("REPRINT_PREFIX")
	os.Unsetenv("REPRINT_CREDENTIALS")

	cfg, err := Load(
		WithBucket("cli-bucket"),
		WithPrefix("cli-prefix/"),
		WithCredentials("/cli/creds.json"),
	)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Bucket != "cli-bucket" {
		t.Errorf("Bucket = %q, want %q", cfg.Bucket, "cli-bucket")
	}
	if cfg.Prefix != "cli-prefix/" {
		t.Errorf("Prefix = %q, want %q", cfg.Prefix, "cli-prefix/")
	}
	if cfg.Credentials != "/cli/creds.json" {
		t.Errorf("Credentials = %q, want %q", cfg.Credentials, "/cli/creds.json")
	}
}

func TestLoad_OptionsPriorityOverEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("REPRINT_BUCKET", "env-bucket")
	os.Setenv("REPRINT_PREFIX", "env-prefix/")
	os.Setenv("REPRINT_CREDENTIALS", "/env/creds.json")
	defer func() {
		os.Unsetenv("REPRINT_BUCKET")
		os.Unsetenv("REPRINT_PREFIX")
		os.Unsetenv("REPRINT_CREDENTIALS")
	}()

	// CLI flags should override env vars
	cfg, err := Load(
		WithBucket("cli-bucket"),
		WithPrefix("cli-prefix/"),
		WithCredentials("/cli/creds.json"),
	)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Bucket != "cli-bucket" {
		t.Errorf("Bucket = %q, want %q", cfg.Bucket, "cli-bucket")
	}
	if cfg.Prefix != "cli-prefix/" {
		t.Errorf("Prefix = %q, want %q", cfg.Prefix, "cli-prefix/")
	}
	if cfg.Credentials != "/cli/creds.json" {
		t.Errorf("Credentials = %q, want %q", cfg.Credentials, "/cli/creds.json")
	}
}

func TestLoad_EmptyOptionsDoNotOverride(t *testing.T) {
	// Set environment variables
	os.Setenv("REPRINT_BUCKET", "env-bucket")
	os.Setenv("REPRINT_PREFIX", "env-prefix/")
	defer func() {
		os.Unsetenv("REPRINT_BUCKET")
		os.Unsetenv("REPRINT_PREFIX")
	}()

	// Empty CLI flags should not override env vars
	cfg, err := Load(
		WithBucket(""),
		WithPrefix(""),
	)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Bucket != "env-bucket" {
		t.Errorf("Bucket = %q, want %q", cfg.Bucket, "env-bucket")
	}
	if cfg.Prefix != "env-prefix/" {
		t.Errorf("Prefix = %q, want %q", cfg.Prefix, "env-prefix/")
	}
}

func TestDefaultCredentialsPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot get home directory")
	}

	tests := []struct {
		name    string
		appName string
		want    string
	}{
		{
			name:    "reprint-gcs",
			appName: "reprint-gcs",
			want:    filepath.Join(home, ".config", "reprint-gcs", "credentials.json"),
		},
		{
			name:    "reprint-s3",
			appName: "reprint-s3",
			want:    filepath.Join(home, ".config", "reprint-s3", "credentials.json"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultCredentialsPath(tt.appName)
			if got != tt.want {
				t.Errorf("DefaultCredentialsPath(%q) = %q, want %q", tt.appName, got, tt.want)
			}
		})
	}
}

func TestLoad_WithAppName_DefaultCredentials(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("REPRINT_BUCKET")
	os.Unsetenv("REPRINT_PREFIX")
	os.Unsetenv("REPRINT_CREDENTIALS")

	// Create a temporary directory for default credentials
	tmpDir := t.TempDir()
	appName := "test-app"
	credDir := filepath.Join(tmpDir, ".config", appName)
	if err := os.MkdirAll(credDir, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	credFile := filepath.Join(credDir, DefaultCredentialsFilename)
	if err := os.WriteFile(credFile, []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to create credentials file: %v", err)
	}

	// Override home directory for test
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg, err := Load(WithAppName(appName))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Credentials != credFile {
		t.Errorf("Credentials = %q, want %q", cfg.Credentials, credFile)
	}
}

func TestLoad_WithAppName_NoDefaultCredentials(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("REPRINT_BUCKET")
	os.Unsetenv("REPRINT_PREFIX")
	os.Unsetenv("REPRINT_CREDENTIALS")

	// Create a temporary directory without credentials file
	tmpDir := t.TempDir()

	// Override home directory for test
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg, err := Load(WithAppName("test-app"))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Should be empty since no default credentials file exists
	if cfg.Credentials != "" {
		t.Errorf("Credentials = %q, want empty", cfg.Credentials)
	}
}

func TestLoad_ExplicitCredentialsPriorityOverDefault(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("REPRINT_BUCKET")
	os.Unsetenv("REPRINT_PREFIX")
	os.Unsetenv("REPRINT_CREDENTIALS")

	// Create a temporary directory with default credentials
	tmpDir := t.TempDir()
	appName := "test-app"
	credDir := filepath.Join(tmpDir, ".config", appName)
	if err := os.MkdirAll(credDir, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	defaultCredFile := filepath.Join(credDir, DefaultCredentialsFilename)
	if err := os.WriteFile(defaultCredFile, []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to create credentials file: %v", err)
	}

	// Override home directory for test
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// Explicit credentials should take priority over default
	explicitCreds := "/explicit/path/to/creds.json"
	cfg, err := Load(
		WithAppName(appName),
		WithCredentials(explicitCreds),
	)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Credentials != explicitCreds {
		t.Errorf("Credentials = %q, want %q", cfg.Credentials, explicitCreds)
	}
}
