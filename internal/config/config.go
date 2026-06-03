package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port    int
	DBPath  string
	AppName string
	AppBaseURL string

	// Email. PostmarkServerToken empty → dev log sender (writes mail
	// to slog instead of sending). PostmarkFrom is the From: header on
	// outgoing mail; Postmark accepts either "user@example.com" or
	// "Name <user@example.com>", and the address must match a Sender
	// Signature or verified domain on the Postmark server.
	PostmarkServerToken string
	PostmarkFrom        string

	// Cleanup worker interval. Defaults to 1h; tunable via
	// CLEANUP_INTERVAL for shorter-cycle testing.
	CleanupInterval time.Duration

	// Tracking pixels (omit to disable).
	PixelID string
	GtagID  string

	// Cloudflare Turnstile. Both keys must be set for verification to
	// run. When SecretKey is empty the contact form skips the widget
	// and the server-side verifier no-ops, so local dev works without
	// a Cloudflare account.
	TurnstileSiteKey   string
	TurnstileSecretKey string
}

func Load() (*Config, error) {
	port, err := parseInt("PORT", 8080)
	if err != nil {
		return nil, err
	}

	cleanupInterval, err := parseDuration("CLEANUP_INTERVAL", 1*time.Hour)
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:                port,
		DBPath:              envDefault("DB_PATH", "./data/app.db"),
		AppName:             envDefault("APP_NAME", "Standard Template"),
		AppBaseURL:          envDefault("APP_BASE_URL", "http://localhost:8080"),
		PostmarkServerToken: os.Getenv("POSTMARK_SERVER_TOKEN"),
		PostmarkFrom:        envDefault("POSTMARK_FROM", "admin@localhost"),
		CleanupInterval:     cleanupInterval,
		PixelID:             os.Getenv("PIXEL_ID"),
		GtagID:              os.Getenv("GTAG_ID"),
		TurnstileSiteKey:    os.Getenv("TURNSTILE_SITE_KEY"),
		TurnstileSecretKey:  os.Getenv("TURNSTILE_SECRET_KEY"),
	}, nil
}

// Addr returns the server address string in the format expected by
// http.ListenAndServe.
func (c *Config) Addr() string {
	return fmt.Sprintf(":%d", c.Port)
}

func envDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func parseInt(key string, fallback int) (int, error) {
	val := os.Getenv(key)
	if val == "" {
		return fallback, nil
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid value for %s: %q", key, val)
	}
	return n, nil
}

func parseDuration(key string, fallback time.Duration) (time.Duration, error) {
	val := os.Getenv(key)
	if val == "" {
		return fallback, nil
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return 0, fmt.Errorf("invalid value for %s: %q", key, val)
	}
	return d, nil
}
