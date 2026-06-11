// Package config centralises application configuration loaded from the
// environment. Values mirror the original Flask backend so behaviour is
// identical across the Python and Go implementations.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds every runtime setting for the backend.
type Config struct {
	Env       string
	Debug     bool
	HTTPHost  string
	HTTPPort  string
	SecretKey string

	// JWT
	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration

	// PostgreSQL
	DatabaseURL string

	// Redis
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	RedisEnabled  bool

	// Storage directories
	DataDir   string
	UploadDir string
	BackupDir string
	ExportDir string

	// Business rules
	MinClientAge            int
	MinDrivingExperienceYrs int
	MaxUploadBytes          int64
	LateFeeMultiplier       float64

	// Backups
	EnableScheduledBackups bool
	BackupRetentionDays    int

	// SMS (zudsms.tj)
	SMSEnabled bool
	SMSURL     string
	SMSLogin   string
	SMSSender  string
	SMSSecret  string

	// PublicURL is the public base URL of the frontend (no trailing slash),
	// used to build links encoded into QR codes (e.g. the rental status page
	// a client opens after scanning). Override with your domain/IP in prod.
	PublicURL string

	// CORS
	CORSOrigins []string
}

// getenv returns the environment variable or a fallback default.
func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getenvBool(key string, fallback bool) bool {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return strings.EqualFold(v, "true") || v == "1" || strings.EqualFold(v, "yes")
	}
	return fallback
}

// Load reads configuration from the environment (and an optional .env file).
func Load() *Config {
	// Best-effort load of a local .env file; ignored if it does not exist.
	_ = godotenv.Load()

	dataDir := getenv("RENTACAR_DATA_DIR", "./data")
	cfg := &Config{
		Env:       getenv("RENTACAR_ENV", "development"),
		HTTPHost:  getenv("RENTACAR_HOST", "0.0.0.0"),
		HTTPPort:  getenv("RENTACAR_PORT", "5000"),
		SecretKey: getenv("RENTACAR_SECRET_KEY", "change-me-in-production"),

		JWTAccessExpiry:  time.Duration(getenvInt("RENTACAR_JWT_HOURS", 12)) * time.Hour,
		JWTRefreshExpiry: 30 * 24 * time.Hour,

		DatabaseURL: getenv(
			"RENTACAR_DATABASE_URI",
			"postgres://rentacar:rentacar@localhost:5432/rentacar?sslmode=disable",
		),

		RedisAddr:     getenv("RENTACAR_REDIS_ADDR", "localhost:6379"),
		RedisPassword: getenv("RENTACAR_REDIS_PASSWORD", ""),
		RedisDB:       getenvInt("RENTACAR_REDIS_DB", 0),
		RedisEnabled:  getenvBool("RENTACAR_REDIS_ENABLED", true),

		DataDir:   dataDir,
		UploadDir: getenv("RENTACAR_UPLOAD_DIR", filepath.Join(dataDir, "uploads")),
		BackupDir: getenv("RENTACAR_BACKUP_DIR", filepath.Join(dataDir, "backups")),
		ExportDir: getenv("RENTACAR_EXPORT_DIR", filepath.Join(dataDir, "exports")),

		MinClientAge:            getenvInt("RENTACAR_MIN_CLIENT_AGE", 21),
		MinDrivingExperienceYrs: getenvInt("RENTACAR_MIN_DRIVING_EXPERIENCE", 1),
		MaxUploadBytes:          32 * 1024 * 1024,
		LateFeeMultiplier:       1.5,

		EnableScheduledBackups: getenvBool("RENTACAR_ENABLE_BACKUPS", true),
		BackupRetentionDays:    getenvInt("RENTACAR_BACKUP_RETENTION", 30),

		SMSEnabled: getenvBool("RENTACAR_SMS_ENABLED", false),
		SMSURL:     getenv("RENTACAR_SMS_URL", "https://api.zudsms.tj/api/v1/send"),
		SMSLogin:   getenv("RENTACAR_SMS_LOGIN", ""),
		SMSSender:  getenv("RENTACAR_SMS_SENDER", ""),
		SMSSecret:  getenv("RENTACAR_SMS_SECRET", ""),

		PublicURL: strings.TrimRight(getenv("RENTACAR_PUBLIC_URL", "http://localhost:3000"), "/"),
	}

	cfg.JWTSecret = getenv("RENTACAR_JWT_SECRET", cfg.SecretKey)
	cfg.Debug = cfg.Env != "production"

	origins := getenv("RENTACAR_CORS_ORIGINS", "*")
	cfg.CORSOrigins = splitAndTrim(origins)

	return cfg
}

// EnsureDirectories creates the runtime directories used for uploads, backups
// and exports. Errors are returned so the caller can fail fast on start-up.
func (c *Config) EnsureDirectories() error {
	for _, dir := range []string{c.DataDir, c.UploadDir, c.BackupDir, c.ExportDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

func splitAndTrim(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}
