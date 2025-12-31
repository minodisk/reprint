package config

import (
	"os"
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
