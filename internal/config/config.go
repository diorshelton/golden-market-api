package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL        string
	JWTSecret          string
	RefreshSecret      string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	AllowedOrigins     []string
	Port               string
	Environment        string
}

const redacted = "[REDACTED]"

// String implements fmt.Stringer, redacting DatabaseURL (which embeds
// credentials) and the JWT/refresh secrets so an accidental
// log.Printf("%v", cfg) or similar doesn't leak them.
func (c *Config) String() string {
	return fmt.Sprintf(
		"Config{DatabaseURL:%s JWTSecret:%s RefreshSecret:%s AccessTokenExpiry:%s RefreshTokenExpiry:%s AllowedOrigins:%v Port:%s Environment:%s}",
		redacted, redacted, redacted, c.AccessTokenExpiry, c.RefreshTokenExpiry, c.AllowedOrigins, c.Port, c.Environment,
	)
}

// GoString implements fmt.GoStringer, so "%#v" is redacted the same way.
func (c *Config) GoString() string {
	return c.String()
}

var defaultAllowedOrigins = []string{
	"http://localhost:5173", // Vite default
	"http://localhost:3000", // React default
	"http://localhost:8080",
}

// Load reads configuration from environment variables (and .env, if present)
// and returns a populated, validated Config.
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found, using environment variables")
	}

	required := map[string]string{
		"DATABASE_URL":         "",
		"JWT_SECRET":           "",
		"REFRESH_SECRET":       "",
		"ACCESS_TOKEN_EXPIRY":  "",
		"REFRESH_TOKEN_EXPIRY": "",
	}

	missing := []string{}
	for k := range required {
		v := os.Getenv(k)
		if v == "" {
			missing = append(missing, k)
		}
		required[k] = v
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("required environment variable(s) missing: %s", strings.Join(missing, ", "))
	}

	accessTokenExpiry, err := time.ParseDuration(required["ACCESS_TOKEN_EXPIRY"])
	if err != nil {
		return nil, fmt.Errorf("invalid ACCESS_TOKEN_EXPIRY: %w", err)
	}

	refreshTokenExpiry, err := time.ParseDuration(required["REFRESH_TOKEN_EXPIRY"])
	if err != nil {
		return nil, fmt.Errorf("invalid REFRESH_TOKEN_EXPIRY: %w", err)
	}

	allowedOrigins := defaultAllowedOrigins
	if raw := os.Getenv("ALLOWED_ORIGINS"); raw != "" {
		allowedOrigins = nil
		for _, origin := range strings.Split(raw, ",") {
			allowedOrigins = append(allowedOrigins, strings.TrimSpace(origin))
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	return &Config{
		DatabaseURL:        required["DATABASE_URL"],
		JWTSecret:          required["JWT_SECRET"],
		RefreshSecret:      required["REFRESH_SECRET"],
		AccessTokenExpiry:  accessTokenExpiry,
		RefreshTokenExpiry: refreshTokenExpiry,
		AllowedOrigins:     allowedOrigins,
		Port:               port,
		Environment:        environment,
	}, nil
}
