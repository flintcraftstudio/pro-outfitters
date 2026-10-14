package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/flintcraftstudio/standard-template/internal/cleanup"
	"github.com/flintcraftstudio/standard-template/internal/config"
	"github.com/flintcraftstudio/standard-template/internal/email"
	"github.com/flintcraftstudio/standard-template/internal/handler"
	"github.com/flintcraftstudio/standard-template/internal/middleware"
	"github.com/flintcraftstudio/standard-template/internal/migrations"
	"github.com/flintcraftstudio/standard-template/internal/session"
	"github.com/flintcraftstudio/standard-template/internal/store"
	"github.com/flintcraftstudio/standard-template/internal/turnstile"
	"github.com/flintcraftstudio/standard-template/internal/view"

	"github.com/pressly/goose/v3"
	"golang.org/x/time/rate"
	_ "modernc.org/sqlite"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := loadEnv(".env"); err != nil {
		slog.Error("env error", "err", err)
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config error", "err", err)
		os.Exit(1)
	}

	// Make tracking + Turnstile site key + app name available to templates.
	view.SiteName = cfg.AppName
	view.GtagID = cfg.GtagID
	view.PixelID = cfg.PixelID
	view.TurnstileSiteKey = cfg.TurnstileSiteKey
	view.BaseURL = cfg.AppBaseURL
	view.Indexable = cfg.AllowIndexing
	if !cfg.AllowIndexing {
		slog.Info("ALLOW_INDEXING not set, site is noindex")
	}
	if cfg.GtagID == "" {
		slog.Info("GTAG_ID not set, Google Analytics disabled")
	}
	if cfg.PixelID == "" {
		slog.Info("PIXEL_ID not set, Facebook Pixel disabled")
	}
	if cfg.TurnstileSiteKey == "" || cfg.TurnstileSecretKey == "" {
		slog.Info("TURNSTILE_*_KEY not set, Turnstile disabled")
	}

	// Database — modernc/sqlite is pure-Go so CGO stays off.
	if err := os.MkdirAll(filepath.Dir(cfg.DBPath), 0o755); err != nil {
		slog.Error("create database directory", "err", err)
		os.Exit(1)
	}
	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		slog.Error("database open", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	if _, err := db.Exec("PRAGMA journal_mode=WAL; PRAGMA synchronous=NORMAL; PRAGMA foreign_keys=ON; PRAGMA busy_timeout=5000;"); err != nil {
		slog.Error("database pragma", "err", err)
		os.Exit(1)
	}

	// Run migrations from the embedded FS so a deployed binary doesn't
	// need the migrations/ folder shipped alongside it.
	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		slog.Error("goose dialect", "err", err)
		os.Exit(1)
	}
	if err := goose.Up(db, "."); err != nil {
		slog.Error("migration", "err", err)
		os.Exit(1)
	}
	slog.Info("migrations applied")

	st := store.New(db)

	// Email sender. Postmark in prod when POSTMARK_SERVER_TOKEN is set;
	// otherwise the dev log sender writes outgoing mail (including
	// magic-link URLs) to slog so local testing works without API keys.
	var mailer email.Sender
	if cfg.PostmarkServerToken != "" {
		mailer = email.NewPostmarkSender(cfg.PostmarkServerToken, cfg.PostmarkFrom)
		slog.Info("email sender", "kind", "postmark", "from", cfg.PostmarkFrom)
	} else {
		mailer = email.LogSender{}
		slog.Info("email sender", "kind", "log")
	}

	turnstileVerifier := turnstile.New(cfg.TurnstileSecretKey)

	mux := http.NewServeMux()

	// Per-IP limits on the two forms that send email. They key on
	// RemoteAddr, which middleware.RealIP has already set to the client
	// IP (see the middleware composition below). Limits are in memory and reset on restart, which is fine at this
	// scale. Contact allows a burst of 8 (validation retries count) and
	// refills one every 10 minutes; login allows 5, then one a minute.
	// Public form bodies are capped well above any real submission (the
	// contact form's fields total about 6 KB at their maxlength) so a bot
	// can't post megabytes for the server to parse and store.
	const maxFormBytes = 64 << 10

	contactLimit := middleware.RateLimit(middleware.RateLimitConfig{
		Rate:    rate.Every(10 * time.Minute),
		Burst:   8,
		OnLimit: handler.ContactRateLimited,
	})
	loginLimit := middleware.RateLimit(middleware.RateLimitConfig{
		Rate:    rate.Every(time.Minute),
		Burst:   5,
		OnLimit: handler.LoginRateLimited,
	})

	// Static files. Cache-Control is set per asset class so repeat visits
	// don't re-download fonts/images; the unhashed CSS/JS get a short TTL so
	// a deploy propagates quickly.
	mux.Handle("GET /static/", http.StripPrefix("/static/", cacheStatic(http.FileServer(http.Dir("web/static")))))

	// Favicon + PWA manifest assets. Browsers and crawlers request these at
	// the site root (/favicon.ico, /apple-touch-icon.png) regardless of the
	// <link> tags, and the manifest's own icon srcs are root-relative — so
	// serve them at their conventional root paths. The files live under
	// web/static/favicon so the existing web/ COPY ships them.
	for urlPath, name := range map[string]string{
		"/favicon.ico":                  "favicon.ico",
		"/favicon.svg":                  "favicon.svg",
		"/favicon-96x96.png":            "favicon-96x96.png",
		"/apple-touch-icon.png":         "apple-touch-icon.png",
		"/site.webmanifest":             "site.webmanifest",
		"/web-app-manifest-192x192.png": "web-app-manifest-192x192.png",
		"/web-app-manifest-512x512.png": "web-app-manifest-512x512.png",
	} {
		file := filepath.Join("web/static/favicon", name)
		mux.HandleFunc("GET "+urlPath, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "public, max-age=604800")
			// Go doesn't know the .webmanifest extension; set the spec type
			// so ServeFile won't sniff it to text/plain.
			if strings.HasSuffix(name, ".webmanifest") {
				w.Header().Set("Content-Type", "application/manifest+json")
			}
			http.ServeFile(w, r, file)
		})
	}

	// Crawler files. Both follow ALLOW_INDEXING.
	mux.Handle("GET /robots.txt", handler.Robots(cfg.AppBaseURL, cfg.AllowIndexing))
	mux.Handle("GET /sitemap.xml", handler.Sitemap(cfg.AppBaseURL))

	// Public pages.
	mux.Handle("GET /{$}", handler.Home())
	mux.Handle("GET /lodges", handler.Lodges())
	mux.Handle("GET /lodges/{slug}", handler.Lodge())
	mux.Handle("GET /smith-river", handler.SmithRiver())
	mux.Handle("GET /about", handler.About())
	// Draft pages (noindex, content pending from the client).
	mux.Handle("GET /fly-fishing", handler.Page(view.FlyFishingPage()))
	mux.Handle("GET /bird-hunting", handler.Page(view.BirdHuntingPage()))
	mux.Handle("GET /kids-camp", handler.Page(view.KidsCampPage()))
	mux.Handle("GET /blackfoot-river", handler.Page(view.BlackfootRiverPage()))
	mux.Handle("GET /packing-lists/{slug}", handler.PackingList())
	mux.Handle("GET /contact", handler.Contact())
	mux.Handle("POST /contact", http.MaxBytesHandler(contactLimit(handler.ContactSubmit(handler.ContactDeps{
		Store:     st,
		Mailer:    mailer,
		Turnstile: turnstileVerifier,
		AppName:   cfg.AppName,
	})), maxFormBytes))

	// Auth.
	authDeps := handler.AuthDeps{
		Store:   st,
		Mailer:  mailer,
		BaseURL: cfg.AppBaseURL,
		AppName: cfg.AppName,
	}
	mux.Handle("GET /login", handler.LoginPage())
	mux.Handle("POST /login", http.MaxBytesHandler(loginLimit(handler.RequestMagicLink(authDeps)), maxFormBytes))
	mux.Handle("GET /login/magic/{token}", handler.ConsumeMagicLink(st))
	mux.Handle("POST /logout", handler.Logout(st))

	// Admin.
	mux.Handle("GET /admin", session.RequireAuth(handler.Dashboard(st)))
	mux.Handle("GET /admin/inquiries", session.RequireAuth(handler.Inquiries(st)))
	mux.Handle("GET /admin/inquiries/{id}", session.RequireAuth(handler.InquiryDetail(st)))
	mux.Handle("POST /admin/inquiries/{id}/respond", session.RequireAuth(handler.ToggleInquiry(st, "responded", true)))
	mux.Handle("POST /admin/inquiries/{id}/unrespond", session.RequireAuth(handler.ToggleInquiry(st, "responded", false)))
	mux.Handle("POST /admin/inquiries/{id}/archive", session.RequireAuth(handler.ToggleInquiry(st, "archived", true)))
	mux.Handle("POST /admin/inquiries/{id}/unarchive", session.RequireAuth(handler.ToggleInquiry(st, "archived", false)))
	mux.Handle("GET /admin/settings", session.RequireAuth(handler.Settings(st)))
	mux.Handle("POST /admin/settings", session.RequireAuth(handler.UpdateSettings(st)))

	// 404 catch-all (registered last so explicit routes win).
	mux.Handle("GET /", handler.NotFound())

	// Middleware composition. Inside-out: session attaches user
	// context → strip trailing slash for canonical paths → old-site
	// redirects (ahead of the slash strip, so each is one 301) → request
	// logging, so it sees the final status code → real client IP
	// outermost. TrustedProxies is 1: Caddy sits in front and appends the
	// client IP to X-Forwarded-For. With no header (local dev) the
	// connection's own address is kept.
	srv := session.Middleware(st)(mux)
	srv = withPath(srv)
	srv = middleware.StripTrailingSlash(srv)
	srv = handler.LegacyRedirects(srv)
	srv = middleware.Logging(logger)(srv)
	srv = middleware.RealIP(1)(srv)

	server := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      srv,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Root context cancels on SIGINT/SIGTERM. Background workers
	// observe this ctx to stop cleanly.
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cleanupDone := make(chan struct{})
	cleanupWorker := cleanup.NewWorker(cleanup.Config{Interval: cfg.CleanupInterval}, st, logger)
	go func() {
		defer close(cleanupDone)
		cleanupWorker.Run(rootCtx)
	}()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("server starting", "addr", cfg.Addr())
		fmt.Printf("listening on %s\n", cfg.Addr())
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-rootCtx.Done():
		slog.Info("shutdown signal received")
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}

	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("shutdown deadline exceeded, forcing close", "err", err)
		server.Close()
		os.Exit(1)
	}

	select {
	case <-cleanupDone:
	case <-time.After(5 * time.Second):
		slog.Warn("cleanup worker did not stop in time")
	}

	slog.Info("server stopped gracefully")
}

// loadEnv reads a .env file and sets environment variables if not
// already set.
func loadEnv(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimRight(line, "\r")
		if line == "" || line[0] == '#' {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
	return nil
}

// cacheStatic wraps a static file handler and stamps a Cache-Control header
// chosen by file extension. Fonts and images are content-stable and get a
// one-year lifetime; the unhashed CSS/JS change on every deploy, so they get
// a short TTL to bound how long a client can serve stale assets.
// withPath records the request path for the layout's canonical link. It
// runs after StripTrailingSlash, so the path is already canonical.
func withPath(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(view.WithPath(r.Context(), r.URL.Path)))
	})
}

func cacheStatic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.ToLower(filepath.Ext(r.URL.Path)) {
		case ".woff2", ".woff", ".ttf", ".otf":
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		case ".jpg", ".jpeg", ".png", ".webp", ".avif", ".gif", ".svg", ".ico":
			w.Header().Set("Cache-Control", "public, max-age=31536000")
		case ".css", ".js":
			w.Header().Set("Cache-Control", "public, max-age=3600")
		default:
			w.Header().Set("Cache-Control", "public, max-age=3600")
		}
		next.ServeHTTP(w, r)
	})
}
