# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository or in projects derived from it.

## What this repo is

The FlintCraft "standard" starter — a Go/templ/Tailwind brochure site with built-in CRM (contact-form inbox) and CMS plumbing (admin auth, editable site settings, SQLite persistence). Designed to be cloned for new client projects and grown into a managed-content site without rewiring the foundation.

If you're reading this in a project that *used* this template, the bones below are still here; per-project content models (e.g. portfolio pieces, products, blog posts) live in additional migrations and packages the project added itself.

## Tech stack

- **Go** (stdlib `net/http`, no router framework) — Go 1.22+ ServeMux for path-based routing
- **templ** — type-safe HTML templating; `*_templ.go` files are gitignored, regenerated via `templ generate`
- **Tailwind CSS v3** (standalone CLI) — no Node toolchain, downloaded by `mage installTailwind`
- **HTMX + Alpine.js** — minimal client-side interactivity
- **SQLite** via `modernc.org/sqlite` (pure Go, CGO disabled) — one file per site
- **goose** — migration tool, migrations embedded into the binary via `//go:embed`
- **Postmark** — transactional email; falls back to a slog "log sender" in dev
- **Cloudflare Turnstile** — bot protection on the contact form (optional)
- **Mage** — task runner (`magefile.go`)

## Domains

The admin breaks into three parallel areas. Each has its own store file, handler file, and templ views; they share layout/auth/middleware but otherwise don't interleave. When adding a new feature, figure out which domain it belongs to first.

1. **Inquiries** (CRM) — contact-form submissions. `store/inquiries.go`, `handler/inquiries.go`, `view/admin_inquiries.templ`. Public site posts to `/contact`; admin manages status via `responded_at` / `archived_at` toggle columns.
2. **Settings** — single-row `site_settings` table for editable public-site values (contact email, phone, tagline). `store/settings.go`, `handler/settings.go`, `view/admin_settings.templ`. Public site reads via `store.GetSettings(ctx)` on each request — no caching, single-row PK lookup on WAL-mode SQLite is faster than any cache invalidation we'd write.
3. **Auth** — passwordless magic-link sign-in. `store/login_tokens.go`, `handler/auth.go`, `view/login.templ`. Sessions in `internal/session`, attached to request context by `session.Middleware`.

When you grow the project to need a content model (portfolio pieces, products, etc.), add it as a new domain following the same shape: one store file, one handler file, one or more templ views.

## Routes

Defined in `cmd/server/main.go`:

**Public**
- `GET /` — homepage
- `GET /contact` — contact form
- `POST /contact` — validates, persists an inquiry row, fires a background email to `site_settings.contact_email`

**Auth** (no session required)
- `GET /login` — sign-in form
- `POST /login` — issues a single-use magic-link token, emails it (or logs it in dev), returns "check your email"
- `GET /login/magic/{token}` — consumes the token, creates a session, redirects to `/admin`
- `POST /logout` — destroys the session, redirects to `/`

**Admin** (wrapped in `session.RequireAuth`)
- `GET /admin` — dashboard with inquiry stats
- `GET /admin/inquiries?filter=all|unresponded|responded|archived` — inbox
- `GET /admin/inquiries/{id}` — detail
- `POST /admin/inquiries/{id}/respond|unrespond|archive|unarchive` — status toggles
- `GET /admin/settings`, `POST /admin/settings` — editable site settings

The 404 catch-all is registered last (`GET /`) so explicit routes win.

## Middleware composition

`main.go` composes inside-out so the innermost middleware is closest to the handler:

```
session.Middleware(store)         ← attaches *User to request context
  → middleware.StripTrailingSlash ← canonical paths (308 redirects)
  → middleware.Logging            ← outermost; sees final status code
```

Other middleware shipped but not currently wired:
- `middleware.CSRF` — double-submit cookie pattern. Wire when admin forms need protection; uses `middleware.TemplateField(r)` to render the hidden input.
- `middleware.CORS` — for cross-origin JSON APIs. Not needed for the form-based admin.
- `middleware.Recovery` — panic catcher. Add once you start handling JSON or anything where a panic could leak partial response bytes.
- `middleware.RateLimit` — token-bucket per client IP. Useful on the public contact form once you start getting traffic.

## Database

SQLite at `DB_PATH` (defaults to `./data/app.db`). PRAGMAs applied at connection time:

```
journal_mode = WAL
synchronous = NORMAL
foreign_keys = ON
busy_timeout = 5000
```

Migrations live in `internal/migrations/*.sql` and are embedded via `//go:embed` so the deployed binary doesn't need them on disk. They run on app startup; a failed migration exits the process.

**`time.Time` storage is fragile under modernc/sqlite.** The driver serializes `time.Time` via `Stringer`, which leaves monotonic-clock junk on the value. Scanning that back into a `time.Time` then fails silently — sessions stop authenticating, etc. Always format explicitly with `store.SqliteDatetime` when inserting a `time.Time`, and parse explicitly on the way out. Columns whose value comes from `DEFAULT CURRENT_TIMESTAMP` are fine — SQLite writes a clean string in those cases.

**Nullable timestamp columns store as `sql.NullString`, not `sql.NullTime`.** SQLite's `CURRENT_TIMESTAMP` emits `YYYY-MM-DD HH:MM:SS`, which isn't RFC3339 — `sql.NullTime.Scan` silently fails on it and the whole row drops. `Inquiry.RespondedAt` and `Inquiry.ArchivedAt` are `sql.NullString`; consumers only check `.Valid`. If you add another nullable timestamp column, follow the same pattern.

## Email

`internal/email` ships two `Sender` implementations:
- `PostmarkSender` — production. Used when `POSTMARK_SERVER_TOKEN` is set.
- `LogSender` — dev fallback. Writes outgoing mail (including magic-link URLs) to slog so local testing works without API keys or DNS setup.

Selection happens once in `main.go`. Handlers depend only on the `email.Sender` interface, so swapping providers is a single-file change.

The contact form sends notification emails as a background goroutine with a fresh `context.Background()` + 15s timeout — a slow Postmark response shouldn't block the user's success render, and a client disconnect mid-render shouldn't cancel the send.

## Auth

**Magic-link only**, no passwords. `/login` collects an email; the handler creates a SHA-256-hashed single-use token in `login_tokens` and emails the raw token embedded in `{APP_BASE_URL}/login/magic/{token}`. 15-minute TTL.

Requests for unknown emails always render the same "check your email" confirmation to prevent enumeration.

Sessions are cookie-backed (`session_token`, HttpOnly, Secure, SameSite=Lax, 7-day max age). `session.Middleware` attaches `*User` to request context when a valid cookie is present. `session.FromContext(ctx)` returns `nil` for anonymous requests. Wrap protected handlers with `session.RequireAuth`.

The first admin user is created via the seed CLI:

```bash
mage seed admin@example.com "Display Name"
```

The seed CLI opens the configured `DB_PATH` directly — migrations must already have run (start the server once first, or it auto-runs on the seed open if migrations are embedded too — they're not, so start the server first).

## Adding a content model

When the project grows past brochure + CRM and needs managed content (pieces, products, posts, etc.):

1. Add a new migration in `internal/migrations/NNN_description.sql` with `-- +goose Up` / `-- +goose Down` sections. Goose runs them on startup.
2. Add a `store/<domain>.go` with the CRUD methods. Hand-rolled SQL is fine — the project doesn't standardize on sqlc unless query count grows.
3. Add a `handler/<domain>.go` with HTTP handlers. Use `session.FromContext` for the user; surface domain errors via `errors.Is(err, store.ErrXxxNotFound)`.
4. Add `view/admin_<domain>.templ` — follow the existing patterns (`adminPageHeading`, `AdminBase`).
5. Wire the routes in `main.go` under the existing admin block.

Slug invariants, soft-delete patterns, image pipelines, etc. — those decisions are project-specific. If the project needs photos, look at the schluters-metal-art repo for the image-pipeline reference (Cloudflare Images direct upload + R2 mirror); that's intentionally not in this template.

## Routes that may need adjustment per project

- The contact form has `name`, `email`, `phone`, `subject`, `message` fields. `subject` is optional; some projects may want to make it a dropdown enumerated to project-specific categories. Add the validation in `handler/contact.go`'s `validateContact` helper.
- `site_settings` has three default columns (`contact_email`, `contact_phone`, `tagline`). Add per-project columns via a new migration; extend `store.SiteSettings` and `store.SiteSettingsInput` to surface them; extend the form in `view/admin_settings.templ`.
- `view.SiteName` is set from `APP_NAME` env var at startup. The login pages and admin nav read it for branding.

## Gotchas

- **`*_templ.go` files are gitignored.** After editing a `.templ` file, run `templ generate` (or `mage build` / `mage dev`). The build fails until you do.
- **Migrations are embedded** — adding a `.sql` file to `internal/migrations/` is enough; the embed directive picks it up automatically. No separate registration.
- **`./data/` is gitignored.** The SQLite DB plus its WAL/SHM sidecar files live there. Delete the directory to reset state in dev.
- **Cookies use `Secure: true`.** They won't be set over plain HTTP except on `localhost`. If you're proxying through a non-localhost dev hostname without TLS, the session won't stick.
- **CSRF middleware is shipped but not wired** in `main.go`. Admin forms (`/admin/settings`, inquiry status toggles) currently rely on session cookies + same-origin only. Wire `middleware.CSRF` and add `@templ.Raw(middleware.TemplateField(ctx))` to forms once it's needed.
- **The cleanup worker** (`internal/cleanup`) sweeps hourly by default. It only deletes expired sessions and consumed/expired login tokens. If you add soft-delete with retention windows for content (like schluters does for projects), extend the worker accordingly.

## Project conventions

- **Hand-rolled SQL** in `internal/store`. No sqlc unless query count grows large enough that drift becomes painful.
- **Error sentinels** in the store package (`ErrInquiryNotFound`, `ErrTokenInvalid`). Handlers check via `errors.Is` and translate to HTTP status codes.
- **`InquiryInput`-style structs** for create/update operations. The Input shape carries only writable fields; the read-side struct (e.g. `Inquiry`) carries IDs, timestamps, and computed columns.
- **Time formatting/parsing** always goes through `store.SqliteDatetime`. Never let modernc/sqlite stringify a `time.Time` for you.
- **Templates use `view.SiteName`** for branding so a single env var (`APP_NAME`) reskins both the public site and admin.

## Deploy

The Dockerfile is multi-stage: builder runs `templ generate` → `tailwindcss --minify` → `go build`; runtime is `alpine:3` with the static binary plus `web/`. CGO is disabled (modernc/sqlite is pure Go), so the binary runs anywhere.

The Dockerfile sets `ENV DB_PATH=/var/lib/app/app.db` and declares `VOLUME ["/var/lib/app"]`. `docker-compose.yml` mounts a named volume `app-data` there, so the SQLite database (plus its WAL/SHM sidecars) survives `docker compose up --build` and `--force-recreate`. Don't override `DB_PATH` in `.env` to a path outside `/var/lib/app` unless you also mount that path as a volume — otherwise the DB lives inside the container and dies on rebuild.

For continuous SQLite backup (recommended in production), add a Litestream sidecar — see the schluters-metal-art repo's `docker-compose.yml` for the pattern.
