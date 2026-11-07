// Package config loads the SuggestIT backend configuration from environment
// variables.
package config

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// MinSessionSecretLen is the minimum length in bytes of SESSION_SECRET.
const MinSessionSecretLen = 32

// Postgres holds the database connection settings.
type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
	SSLMode  string
}

// DSN returns a lib/pq connection URL. All parts are escaped by net/url, so
// passwords and names containing special characters are safe.
func (p Postgres) DSN() string {
	q := url.Values{}
	q.Set("sslmode", p.SSLMode)
	q.Set("connect_timeout", "5")
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(p.User, p.Password),
		Host:     net.JoinHostPort(p.Host, p.Port),
		Path:     "/" + p.DB,
		RawQuery: q.Encode(),
	}
	return u.String()
}

// Config is the complete backend configuration.
type Config struct {
	ListenAddr string
	Postgres   Postgres

	OIDCIssuer       string
	OIDCClientID     string
	OIDCClientSecret string
	CallbackURL      string
	// PublicOrigin is the scheme://host[:port] of CALLBACK_URL, i.e. the
	// origin browsers use for the site. It is trusted for cross-origin
	// protection so that browsers without Sec-Fetch-Site still work when a
	// reverse proxy rewrites the Host header.
	PublicOrigin string

	// CookieDomain is empty for host-only cookies.
	CookieDomain string
	// SecureCookies is false only when CALLBACK_URL uses plain http (dev).
	SecureCookies bool

	SessionSecret []byte
	// GeneratedSecret is true when SESSION_SECRET was unset in a local dev
	// setup and a random, per-process secret was generated instead.
	GeneratedSecret bool
}

// Load reads the configuration using getenv (normally os.Getenv).
func Load(getenv func(string) string) (Config, error) {
	env := func(k string) string { return strings.TrimSpace(getenv(k)) }

	cfg := Config{
		ListenAddr: env("LISTEN_ADDR"),
		Postgres: Postgres{
			Host:     env("SUGGESTIT_POSTGRES_HOST"),
			Port:     env("SUGGESTIT_POSTGRES_PORT"),
			User:     env("SUGGESTIT_POSTGRES_USER"),
			Password: getenv("SUGGESTIT_POSTGRES_PASSWORD"), // not trimmed: may legitimately contain spaces
			DB:       env("SUGGESTIT_POSTGRES_DB"),
			SSLMode:  env("SUGGESTIT_POSTGRES_SSLMODE"),
		},
		OIDCIssuer:       strings.TrimRight(env("OAUTH_URL"), "/"),
		OIDCClientID:     env("OAUTH_CLIENT_ID"),
		OIDCClientSecret: env("OAUTH_CLIENT_SECRET"),
		CallbackURL:      env("CALLBACK_URL"),
		CookieDomain:     env("COOKIE_DOMAIN"),
	}
	if cfg.ListenAddr == "" {
		cfg.ListenAddr = ":3001"
	}
	if cfg.Postgres.Port == "" {
		cfg.Postgres.Port = "5432"
	}
	if cfg.Postgres.SSLMode == "" {
		cfg.Postgres.SSLMode = "disable"
	}

	var errs []error
	for _, req := range []struct{ name, val string }{
		{"SUGGESTIT_POSTGRES_HOST", cfg.Postgres.Host},
		{"SUGGESTIT_POSTGRES_USER", cfg.Postgres.User},
		{"SUGGESTIT_POSTGRES_DB", cfg.Postgres.DB},
		{"OAUTH_URL", cfg.OIDCIssuer},
		{"OAUTH_CLIENT_ID", cfg.OIDCClientID},
		{"OAUTH_CLIENT_SECRET", cfg.OIDCClientSecret},
		{"CALLBACK_URL", cfg.CallbackURL},
	} {
		if req.val == "" {
			errs = append(errs, fmt.Errorf("%s is required", req.name))
		}
	}
	if n, err := strconv.Atoi(cfg.Postgres.Port); err != nil || n < 1 || n > 65535 {
		errs = append(errs, fmt.Errorf("SUGGESTIT_POSTGRES_PORT %q is not a valid port", cfg.Postgres.Port))
	}
	if cfg.OIDCIssuer != "" {
		if _, err := parseHTTPURL(cfg.OIDCIssuer); err != nil {
			errs = append(errs, fmt.Errorf("OAUTH_URL: %w", err))
		}
	}

	dev := false
	if cfg.CallbackURL != "" {
		cb, err := parseHTTPURL(cfg.CallbackURL)
		if err != nil {
			errs = append(errs, fmt.Errorf("CALLBACK_URL: %w", err))
		} else {
			cfg.SecureCookies = cb.Scheme != "http"
			cfg.PublicOrigin = origin(cb)
			dev = cb.Scheme == "http" && isLoopback(cb.Hostname())
		}
	}

	switch secret := getenv("SESSION_SECRET"); {
	case secret != "":
		if len(secret) < MinSessionSecretLen {
			errs = append(errs, fmt.Errorf("SESSION_SECRET must be at least %d bytes (got %d)", MinSessionSecretLen, len(secret)))
		}
		cfg.SessionSecret = []byte(secret)
	case dev:
		cfg.SessionSecret = make([]byte, MinSessionSecretLen)
		_, _ = rand.Read(cfg.SessionSecret) // crypto/rand.Read never fails
		cfg.GeneratedSecret = true
	default:
		errs = append(errs, fmt.Errorf("SESSION_SECRET is required (at least %d bytes, e.g. `openssl rand -base64 48`)", MinSessionSecretLen))
	}

	if err := errors.Join(errs...); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func parseHTTPURL(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return nil, fmt.Errorf("%q must be an absolute http(s) URL", raw)
	}
	return u, nil
}

// origin returns the serialized origin of u the way browsers send it in the
// Origin header: lowercase host and no default port.
func origin(u *url.URL) string {
	host := strings.ToLower(u.Host)
	if p := u.Port(); (u.Scheme == "http" && p == "80") || (u.Scheme == "https" && p == "443") {
		host = strings.ToLower(u.Hostname())
		if strings.Contains(host, ":") {
			host = "[" + host + "]"
		}
	}
	return u.Scheme + "://" + host
}

func isLoopback(host string) bool {
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}
