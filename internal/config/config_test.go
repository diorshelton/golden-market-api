package config

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	validEnv := map[string]string{
		"DATABASE_URL":         "postgres://user:pass@localhost:5432/db",
		"JWT_SECRET":           "test-jwt-secret",
		"REFRESH_SECRET":       "test-refresh-secret",
		"ACCESS_TOKEN_EXPIRY":  "15m",
		"REFRESH_TOKEN_EXPIRY": "168h",
	}

	tests := []struct {
		name      string
		overrides map[string]string
		wantErr   bool
	}{
		{
			name:    "all required vars set",
			wantErr: false,
		},
		{
			name:      "missing REFRESH_TOKEN_EXPIRY",
			overrides: map[string]string{"REFRESH_TOKEN_EXPIRY": ""},
			wantErr:   true,
		},
		{
			name:      "invalid ACCESS_TOKEN_EXPIRY format",
			overrides: map[string]string{"ACCESS_TOKEN_EXPIRY": "not-a-duration"},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range validEnv {
				t.Setenv(k, v)
			}
			for k, v := range tt.overrides {
				t.Setenv(k, v)
			}

			cfg, err := Load()

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cfg.DatabaseURL != validEnv["DATABASE_URL"] {
				t.Errorf("DatabaseURL = %q, want %q", cfg.DatabaseURL, validEnv["DATABASE_URL"])
			}
			if cfg.AccessTokenExpiry != 15*time.Minute {
				t.Errorf("AccessTokenExpiry = %v, want %v", cfg.AccessTokenExpiry, 15*time.Minute)
			}
			if cfg.RefreshTokenExpiry != 168*time.Hour {
				t.Errorf("RefreshTokenExpiry = %v, want %v", cfg.RefreshTokenExpiry, 168*time.Hour)
			}
		})
	}
}

func TestConfigStringRedactsSecrets(t *testing.T) {
	cfg := &Config{
		DatabaseURL:   "postgres://user:pass@localhost:5432/db",
		JWTSecret:     "super-secret-jwt",
		RefreshSecret: "super-secret-refresh",
		Environment:   "development",
	}

	for _, got := range []string{fmt.Sprintf("%v", cfg), fmt.Sprintf("%+v", cfg), fmt.Sprintf("%#v", cfg)} {
		for _, secret := range []string{cfg.DatabaseURL, cfg.JWTSecret, cfg.RefreshSecret} {
			if strings.Contains(got, secret) {
				t.Errorf("formatted output leaked a secret: %q contains %q", got, secret)
			}
		}
		if !strings.Contains(got, redacted) {
			t.Errorf("formatted output %q does not contain redaction marker %q", got, redacted)
		}
	}
}
