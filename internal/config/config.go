// Package config loads runtime configuration from the environment.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all runtime settings for the genshin-impact app.
type Config struct {
	// APIBase is the root of the Hoyolab data API (no trailing slash).
	APIBase string
	// Addr is the TCP address the HTTP server listens on.
	Addr string
	// CacheTTL is how long cached list endpoints are considered fresh.
	CacheTTL time.Duration
}

// Load reads configuration from the environment, applying defaults.
func Load() (Config, error) {
	cfg := Config{
		APIBase:  envOr("HOYOLAB_API_BASE", "https://hoyo.dotcchix.dev"),
		Addr:     envOr("ADDR", ":8080"),
		CacheTTL: 20 * time.Minute,
	}

	if raw := os.Getenv("CACHE_TTL_SECONDS"); raw != "" {
		secs, err := strconv.Atoi(raw)
		if err != nil {
			return Config{}, fmt.Errorf("invalid CACHE_TTL_SECONDS %q: %w", raw, err)
		}
		cfg.CacheTTL = time.Duration(secs) * time.Second
	}

	return cfg, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
