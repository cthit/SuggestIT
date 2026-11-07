package config

import (
	"net/url"
	"strings"
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"SUGGESTIT_POSTGRES_HOST":     "db",
		"SUGGESTIT_POSTGRES_PORT":     "5432",
		"SUGGESTIT_POSTGRES_USER":     "suggestit",
		"SUGGESTIT_POSTGRES_PASSWORD": "p@ss word/#?",
		"SUGGESTIT_POSTGRES_DB":       "suggestit",
		"OAUTH_URL":                   "https://auth.chalmers.it/",
		"OAUTH_CLIENT_ID":             "id",
		"OAUTH_CLIENT_SECRET":         "secret",
		"CALLBACK_URL":                "https://suggestit.chalmers.it/callback",
		"SESSION_SECRET":              strings.Repeat("s", 32),
	}
}

func load(env map[string]string) (Config, error) {
	return Load(func(k string) string { return env[k] })
}

func TestLoadDefaultsAndProd(t *testing.T) {
	cfg, err := load(baseEnv())
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ListenAddr != ":3001" {
		t.Errorf("ListenAddr = %q", cfg.ListenAddr)
	}
	if cfg.Postgres.SSLMode != "disable" {
		t.Errorf("SSLMode = %q", cfg.Postgres.SSLMode)
	}
	if cfg.OIDCIssuer != "https://auth.chalmers.it" {
		t.Errorf("issuer = %q, want trailing slash trimmed", cfg.OIDCIssuer)
	}
	if !cfg.SecureCookies {
		t.Error("https callback must give secure cookies")
	}
	if cfg.GeneratedSecret {
		t.Error("secret was provided, must not be generated")
	}
	if cfg.CookieDomain != "" {
		t.Errorf("CookieDomain = %q", cfg.CookieDomain)
	}
	if cfg.PublicOrigin != "https://suggestit.chalmers.it" {
		t.Errorf("PublicOrigin = %q", cfg.PublicOrigin)
	}
}

func TestPublicOrigin(t *testing.T) {
	for cb, want := range map[string]string{
		"https://SuggestIT.chalmers.it:443/callback": "https://suggestit.chalmers.it",
		"http://localhost:3000/callback":             "http://localhost:3000",
		"http://[::1]:80/api/auth/callback":          "http://[::1]",
		"https://example.com:8443/cb?x=1":            "https://example.com:8443",
	} {
		env := baseEnv()
		env["CALLBACK_URL"] = cb
		cfg, err := load(env)
		if err != nil {
			t.Fatalf("%s: %v", cb, err)
		}
		if cfg.PublicOrigin != want {
			t.Errorf("CALLBACK_URL %s: PublicOrigin = %q, want %q", cb, cfg.PublicOrigin, want)
		}
	}
}

func TestLoadMissingRequired(t *testing.T) {
	_, err := load(map[string]string{})
	if err == nil {
		t.Fatal("expected error")
	}
	for _, name := range []string{"SUGGESTIT_POSTGRES_HOST", "SUGGESTIT_POSTGRES_USER", "SUGGESTIT_POSTGRES_DB",
		"OAUTH_URL", "OAUTH_CLIENT_ID", "OAUTH_CLIENT_SECRET", "CALLBACK_URL", "SESSION_SECRET"} {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("error does not mention %s: %v", name, err)
		}
	}
}

func TestLoadInvalidValues(t *testing.T) {
	cases := map[string]map[string]string{
		"short secret":      {"SESSION_SECRET": "too-short"},
		"bad port":          {"SUGGESTIT_POSTGRES_PORT": "abc"},
		"relative issuer":   {"OAUTH_URL": "auth.chalmers.it"},
		"bad callback":      {"CALLBACK_URL": "/callback"},
		"ftp callback":      {"CALLBACK_URL": "ftp://example.com/cb"},
		"missing secret":    {"SESSION_SECRET": ""},
		"http non-local":    {"SESSION_SECRET": "", "CALLBACK_URL": "http://suggestit.example.com/callback"},
		"port out of range": {"SUGGESTIT_POSTGRES_PORT": "70000"},
	}
	for name, over := range cases {
		t.Run(name, func(t *testing.T) {
			env := baseEnv()
			for k, v := range over {
				env[k] = v
			}
			if _, err := load(env); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestLoadDevGeneratesSecret(t *testing.T) {
	for _, cb := range []string{"http://localhost:3000/callback", "http://127.0.0.1:3001/api/auth/callback", "http://[::1]:3000/callback"} {
		env := baseEnv()
		env["SESSION_SECRET"] = ""
		env["CALLBACK_URL"] = cb
		cfg, err := load(env)
		if err != nil {
			t.Fatalf("%s: %v", cb, err)
		}
		if !cfg.GeneratedSecret || len(cfg.SessionSecret) < MinSessionSecretLen {
			t.Errorf("%s: expected generated secret", cb)
		}
		if cfg.SecureCookies {
			t.Errorf("%s: http callback must give insecure cookies", cb)
		}
	}
}

func TestLoadOptional(t *testing.T) {
	env := baseEnv()
	env["SUGGESTIT_POSTGRES_PORT"] = ""
	env["SUGGESTIT_POSTGRES_SSLMODE"] = "require"
	env["LISTEN_ADDR"] = "127.0.0.1:8080"
	env["COOKIE_DOMAIN"] = "chalmers.it"
	cfg, err := load(env)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Postgres.Port != "5432" || cfg.Postgres.SSLMode != "require" || cfg.ListenAddr != "127.0.0.1:8080" || cfg.CookieDomain != "chalmers.it" {
		t.Errorf("unexpected config: %+v", cfg)
	}
}

func TestDSNEscapes(t *testing.T) {
	p := Postgres{Host: "db", Port: "5432", User: "us er", Password: "p@ss word/#?", DB: "suggest it", SSLMode: "disable"}
	u, err := url.Parse(p.DSN())
	if err != nil {
		t.Fatal(err)
	}
	pw, _ := u.User.Password()
	if u.Scheme != "postgres" || u.User.Username() != "us er" || pw != "p@ss word/#?" ||
		u.Host != "db:5432" || u.Path != "/suggest it" || u.Query().Get("sslmode") != "disable" {
		t.Errorf("DSN round trip mismatch: %s", p.DSN())
	}

	v6 := Postgres{Host: "::1", Port: "5432", User: "u", DB: "d", SSLMode: "disable"}
	if !strings.Contains(v6.DSN(), "[::1]:5432") {
		t.Errorf("IPv6 host not bracketed: %s", v6.DSN())
	}
}
