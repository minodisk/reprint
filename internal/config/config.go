package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the configuration for reprint CLIs.
type Config struct {
	Bucket      string `mapstructure:"bucket"`
	Prefix      string `mapstructure:"prefix"`
	Credentials string `mapstructure:"credentials"`
}

// Load loads configuration from config file and environment variables.
// Priority (highest to lowest): CLI flags > Environment variables > Config file
func Load() (*Config, error) {
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

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
