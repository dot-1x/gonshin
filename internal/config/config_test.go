package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("HOYOLAB_API_BASE", "")
	t.Setenv("ADDR", "")
	t.Setenv("CACHE_TTL_SECONDS", "")
	t.Setenv("PUBLIC_BASE_URL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIBase != "https://hoyo.dotcchix.dev" {
		t.Errorf("APIBase = %q", cfg.APIBase)
	}
	if cfg.Addr != ":8080" {
		t.Errorf("Addr = %q", cfg.Addr)
	}
	if cfg.CacheTTL != 20*time.Minute {
		t.Errorf("CacheTTL = %v", cfg.CacheTTL)
	}
	if cfg.PublicBaseURL != "" {
		t.Errorf("PublicBaseURL = %q", cfg.PublicBaseURL)
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("HOYOLAB_API_BASE", "http://localhost:8504")
	t.Setenv("ADDR", ":9000")
	t.Setenv("CACHE_TTL_SECONDS", "60")
	t.Setenv("PUBLIC_BASE_URL", "https://genshin.example.com")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIBase != "http://localhost:8504" || cfg.Addr != ":9000" || cfg.CacheTTL != time.Minute {
		t.Errorf("unexpected config: %+v", cfg)
	}
	if cfg.PublicBaseURL != "https://genshin.example.com" {
		t.Errorf("PublicBaseURL = %q", cfg.PublicBaseURL)
	}
}

func TestLoadInvalidTTL(t *testing.T) {
	t.Setenv("CACHE_TTL_SECONDS", "abc")
	if _, err := Load(); err == nil {
		t.Fatal("expected error for invalid CACHE_TTL_SECONDS")
	}
}
