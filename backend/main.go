// Command suggestit is the SuggestIT backend.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"chalmers.it/suggestit/internal/api"
	"chalmers.it/suggestit/internal/auth"
	"chalmers.it/suggestit/internal/config"
	"chalmers.it/suggestit/internal/retry"
	"chalmers.it/suggestit/internal/store"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	if err := run(logger); err != nil {
		logger.Error("startup failed", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load(os.Getenv)
	if err != nil {
		return fmt.Errorf("invalid configuration:\n%w", err)
	}
	if cfg.GeneratedSecret {
		logger.Warn("SESSION_SECRET is not set; using a random secret for local development (sessions will not survive restarts)")
	}
	keys, err := auth.DeriveKeys(cfg.SessionSecret)
	if err != nil {
		return err
	}

	db, err := store.Open(cfg.Postgres.DSN())
	if err != nil {
		return err
	}
	defer db.Close()
	logger.Info("connecting to postgres", "host", cfg.Postgres.Host, "port", cfg.Postgres.Port, "db", cfg.Postgres.DB)
	if err := retry.Startup.Do(ctx, logger, "connect to postgres", db.Ping); err != nil {
		return interrupted(ctx, logger, err)
	}
	if err := db.Migrate(ctx); err != nil {
		return fmt.Errorf("create suggestion table: %w", err)
	}

	cookieOpts := auth.CookieOptions{Secure: cfg.SecureCookies, Domain: cfg.CookieDomain}
	sessions := auth.NewSessions(keys, cookieOpts)
	oidcCfg := auth.OIDCConfig{
		Issuer:       cfg.OIDCIssuer,
		ClientID:     cfg.OIDCClientID,
		ClientSecret: cfg.OIDCClientSecret,
		RedirectURL:  cfg.CallbackURL,
	}
	logger.Info("discovering OIDC provider", "issuer", cfg.OIDCIssuer)
	var provider rp.RelyingParty
	err = retry.Startup.Do(ctx, logger, "OIDC discovery", func(ctx context.Context) error {
		p, err := auth.NewRelyingParty(ctx, oidcCfg, auth.NewFlowCookieHandler(keys, cfg.SecureCookies), logger)
		if errors.Is(err, oidc.ErrIssuerInvalid) {
			return retry.Permanent(fmt.Errorf("%w (OAUTH_URL must equal the provider's issuer exactly)", err))
		}
		provider = p
		return err
	})
	if err != nil {
		return interrupted(ctx, logger, err)
	}

	handler, err := api.New(api.Options{
		Store:          db,
		Sessions:       sessions,
		Login:          auth.LoginHandler(provider, sessions),
		Callback:       auth.CallbackHandler(provider, sessions, logger),
		Logger:         logger,
		TrustedOrigins: []string{cfg.PublicOrigin},
	})
	if err != nil {
		return fmt.Errorf("trusted origin from CALLBACK_URL: %w", err)
	}
	srv := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      60 * time.Second, // the OIDC callback calls Gamma twice
		IdleTimeout:       2 * time.Minute,
		MaxHeaderBytes:    16 << 10,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
	}

	serveErr := make(chan error, 1)
	go func() {
		logger.Info("listening", "addr", cfg.ListenAddr, "secure_cookies", cfg.SecureCookies, "session_cookie", sessions.CookieName())
		serveErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
	}
	stop()
	logger.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	logger.Info("stopped")
	return nil
}

// interrupted turns an error caused by a shutdown signal during startup into
// a clean exit.
func interrupted(ctx context.Context, logger *slog.Logger, err error) error {
	if ctx.Err() != nil {
		logger.Info("interrupted during startup")
		return nil
	}
	return err
}
