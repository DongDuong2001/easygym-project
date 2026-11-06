// Command api is the easyGYM modular monolith: one process, one database, one module per business area.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"easygym/internal/modules/accounts"
	"easygym/internal/modules/catalog"
	"easygym/internal/modules/content"
	"easygym/internal/modules/files"
	"easygym/internal/modules/leads"
	"easygym/internal/modules/orders"
	"easygym/internal/platform"
)

func main() {
	if err := run(); err != nil {
		slog.Error("api stopped", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := platform.LoadConfig()
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := platform.OpenDB(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	// Wiring is the only place modules meet: orders gets the catalog as its Pricer,
	// and every module asks accounts (through auth) who is signed in and with which role.
	auth := &platform.Auth{JWKSURL: cfg.AuthJWKSURL, Issuer: cfg.AuthIssuer, Audience: cfg.AuthAudience}
	acc := accounts.New(db, auth, cfg.AdminEmails)
	auth.Roles = acc
	if cfg.AuthJWKSURL == "" {
		slog.Warn("AUTH_JWKS_URL is not set: sign-in, accounts and admin routes answer 401")
	}
	cat := catalog.New(db, auth)
	mux := http.NewServeMux()
	for _, m := range []platform.Module{acc, cat, content.New(db, auth), leads.New(db, auth),
		orders.New(db, cat, auth), files.New(cfg.FilesDir, auth)} {
		m.Register(mux)
	}
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(r.Context()); err != nil {
			platform.WriteError(w, r, err)
			return
		}
		platform.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	srv := &http.Server{
		Addr: cfg.Addr,
		// 10 orders and 10 sign-ups per IP per 10 minutes: enough for a real customer, not for a spam script.
		Handler:           platform.Logging(platform.RateLimit(mux, 10, 10*time.Minute, "/api/orders", "/api/leads")),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	errc := make(chan error, 1)
	go func() {
		slog.Info("api listening", "addr", cfg.Addr)
		errc <- srv.ListenAndServe()
	}()
	select {
	case err := <-errc:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
	}
	shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdown)
}
