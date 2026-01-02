package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// DefaultCredentialsFilename is the default credentials file name.
const DefaultCredentialsFilename = "credentials.json"

// Config holds the configuration for reprint CLIs.
type Config struct {
	Bucket      string `mapstructure:"bucket"`
	Prefix      string `mapstructure:"prefix"`
	Credentials string `mapstructure:"credentials"`
	appName     string // internal: used for default credentials path
}

// DefaultCredentialsPath returns the default path for credentials file.
// Returns empty string if home directory cannot be determined.
// appName should be the CLI name (e.g., "reprint-gcs", "reprint-s3").
func DefaultCredentialsPath(appName string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", appName, DefaultCredentialsFilename)
}

// Option is a function that modifies Config.
type Option func(*Config)

// WithBucket sets the bucket from CLI flag.
func WithBucket(bucket string) Option {
	return func(c *Config) {
		if bucket != "" {
			c.Bucket = bucket
		}
	}
}

// WithPrefix sets the prefix from CLI flag.
func WithPrefix(prefix string) Option {
	return func(c *Config) {
		if prefix != "" {
			c.Prefix = prefix
		}
	}
}

// WithCredentials sets the credentials from CLI flag.
func WithCredentials(credentials string) Option {
	return func(c *Config) {
		if credentials != "" {
			c.Credentials = credentials
		}
	}
}

// WithAppName sets the app name for default credentials path.
func WithAppName(appName string) Option {
	return func(c *Config) {
		c.appName = appName
	}
}

// Load loads configuration from config file, environment variables, and CLI flags.
// Priority (highest to lowest): CLI flags > Environment variables > Config file
func Load(opts ...Option) (*Config, error) {
	v := viper.New()

	// Config file settings
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// Add config paths
	if home, err := os.UserHomeDir(); err == nil {
		v.AddConfigPath(filepath.Join(home, ".config", "reprint"))
	}

	// Read config file (ignore if not found)
	_ = v.ReadInConfig()

	// Environment variables
	v.SetEnvPrefix("REPRINT")
	v.AutomaticEnv()

	// Bind environment variables explicitly
	v.BindEnv("bucket")
	v.BindEnv("prefix")
	v.BindEnv("credentials")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Apply CLI flag options (highest priority)
	for _, opt := range opts {
		opt(&cfg)
	}

	// If credentials not set, check default path
	if cfg.Credentials == "" && cfg.appName != "" {
		defaultPath := DefaultCredentialsPath(cfg.appName)
		if defaultPath != "" {
			if _, err := os.Stat(defaultPath); err == nil {
				cfg.Credentials = defaultPath
			}
		}
	}

	return &cfg, nil
}
