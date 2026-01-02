package main

import (
	"fmt"

	"github.com/minodisk/reprint/internal/config"
)

func loadConfig() (*config.Config, error) {
	cfg, err := config.Load(
		config.WithAppName(appName),
		config.WithBucket(bucket),
		config.WithPrefix(prefix),
		config.WithCredentials(credentials),
	)
	if err != nil {
		return nil, err
	}

	// Validate required fields
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("bucket is required (--bucket, REPRINT_BUCKET, or config file)")
	}
	if cfg.Credentials == "" {
		return nil, fmt.Errorf("credentials is required (--credentials, REPRINT_CREDENTIALS, config file, or place at %s)", config.DefaultCredentialsPath(appName))
	}

	return cfg, nil
}
