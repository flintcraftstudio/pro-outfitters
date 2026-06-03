# standard-template

Brochure-site starter for FlintCraft client projects. Ships with a small CRM (contact-form inbox) and CMS plumbing (admin auth, site settings) so most projects can grow into managed content without a rebuild.

## Tech Stack

- **Go** (stdlib `net/http`) — server, routing, handlers
- **templ** — type-safe HTML templating
- **Tailwind CSS** (standalone CLI) — utility-first styling
- **HTMX** — form interactions without page reloads
- **Alpine.js** — lightweight client-side interactivity
- **Mage** — build task runner
- **SQLite** (modernc/sqlite, pure Go) + **goose** migrations
- **Postmark** — transactional email (contact-form notifications + magic-link auth)
- **Cloudflare Turnstile** — bot protection

## Getting Started

```bash
mage installTailwind             # download Tailwind CLI (once)
mage dev                         # generate templ, build CSS, run server on :8080
mage seed admin@example.com      # create the first admin user
```

Or build + run the production binary in one step:

```bash
mage Start                       # build CSS + templ + Go, then run ./bin/server
```

Sign in at `/login` — with no Postmark token configured, the magic-link URL is logged to the server output.

Production build:

```bash
mage build                       # CSS + templ generate + go build
./bin/server
```

Docker:

```bash
docker compose up
```

## Routes

**Public**
- `GET /` — homepage
- `GET /contact`, `POST /contact` — contact form (persists an inquiry, then emails the configured contact address)

**Auth**
- `GET /login`, `POST /login`, `GET /login/magic/{token}`, `POST /logout` — passwordless magic-link flow

**Admin** (all behind `session.RequireAuth`)
- `GET /admin` — dashboard
- `GET /admin/inquiries?filter=all|unresponded|responded|archived` — inbox
- `GET /admin/inquiries/{id}` — detail; `POST .../respond|unrespond|archive|unarchive` — status toggles
- `GET /admin/settings`, `POST /admin/settings` — editable contact info + tagline

## Project Structure

```
cmd/
  server/main.go                  # wires DB, migrations, routes, workers
  seed/main.go                    # CLI: create an admin user
internal/
  config/                         # env var loading
  migrations/                     # goose-managed SQL, embedded into the binary
  store/                          # SQLite persistence
  session/                        # cookie-backed sessions, RequireAuth middleware
  email/                          # Postmark + log-only sender
  turnstile/                      # CF Turnstile verifier
  middleware/                     # logging, csrf, cors, ratelimit, recovery, trailing-slash
  cleanup/                        # background sweep: expired sessions + login tokens
  handler/                        # HTTP handlers (public + admin)
  view/                           # templ templates (public + admin)
tailwind/
  tailwind.config.js              # color palette, fonts, content paths
  input.css                       # font imports, Tailwind directives
web/static/                       # compiled CSS + JS (htmx, alpine)
```

## Environment Variables

All optional with graceful degradation. See `.env.example`.

| Variable | Purpose |
|---|---|
| `PORT` | Server port (default: 8080) |
| `APP_NAME` | Display name in templates and magic-link emails |
| `APP_BASE_URL` | Used to build absolute URLs for magic-link emails |
| `DB_PATH` | SQLite database path (default: `./data/app.db`) |
| `POSTMARK_SERVER_TOKEN` | Postmark API key. Empty → log sender (writes mail to stdout) |
| `POSTMARK_FROM` | Sender address (must match a Postmark Sender Signature in production) |
| `CLEANUP_INTERVAL` | Cleanup-worker tick (default: `1h`) |
| `GTAG_ID` | Google Analytics tag |
| `PIXEL_ID` | Facebook Pixel ID |
| `TURNSTILE_SITE_KEY` | Cloudflare Turnstile site key |
| `TURNSTILE_SECRET_KEY` | Cloudflare Turnstile secret key |

## Deployment

This template deploys via **GitHub Actions** to a VPS running Docker + Caddy.

### Flow

1. Push to `main` triggers `.github/workflows/deploy.yml`
2. Docker image is built and pushed to `ghcr.io/flintcraftstudio/{project}:latest`
3. SSH into VPS triggers `docker compose pull && docker compose up -d` in `/opt/{project}/`
4. Caddy reverse-proxies `{domain}` to the container's allocated port

### VPS Provisioning

Run `provision.sh` on the VPS to set up a new project:

```bash
sudo ./provision.sh <project> <domain>
```

This creates:
- `/opt/{project}/docker-compose.yml` — pulls from GHCR, maps an allocated port to `8080`
- `/opt/{project}/.env` — app secrets (Postmark, Turnstile, tracking pixels, etc.)
- `/etc/caddy/sites/{project}.caddy` — HTTPS reverse proxy with security headers
- SSH deploy key restricted to `docker compose pull && up -d`

### GitHub Secrets

Set these in **Settings > Secrets and variables > Actions**:

| Secret | Value |
|---|---|
| `VPS_HOST` | VPS IP or hostname |
| `VPS_USER` | `deploy` |
| `VPS_SSH_KEY` | Private key from `provision.sh` output |

The `GITHUB_TOKEN` (automatic) handles GHCR authentication.

### Container Contract

- Image listens on port **8080**
- Config via environment variables (see `.env` on VPS)
- Service name is `app` in both local and production `docker-compose.yml`
- A named volume `app-data` is mounted at `/var/lib/app` so the SQLite database survives container rebuilds — when you adapt the production `docker-compose.yml`, mirror this mount

## Gotchas

- `*_templ.go` files are gitignored — run `templ generate` (or `mage build`) after editing a `.templ` file or the build fails.
- Migrations are embedded via `//go:embed` and run on app startup; a failed migration exits the process.
- `time.Time` storage in SQLite is fragile under modernc/sqlite — always format/parse explicitly via `store.SqliteDatetime` rather than letting the driver serialize.

---

## Claude Code Skills

This repo includes Claude Code skills invoked with `/skill-name` in conversation. Each skill is a structured prompt that guides Claude through a specific workflow.

### `/two-variation-site`

Build a brochure website presenting two distinct brand/design directions for a client to compare side-by-side.

**When to use:** Starting a new client project where you want to present two visual directions (e.g., "warm" vs "bold") from a single codebase.

**Inputs required before code is written:**

1. **Business details** — name, address, phone, email, hours, social links, tagline
2. **Brand guide for each variation** — color palette (hex values + roles), typography stack (families, sizes, weights for headline/body/accent/UI), voice/tone, layout personality
3. **Page copy** — approved text for each section of each page
4. **Images** — hero images, logos, team photos (or placeholders to use)
5. **Variation names/slugs** — evocative short names for each direction (e.g., "warm" / "bold", "classic" / "modern")

**What it produces:**

- Split-panel landing page at `/` comparing both directions
- Complete variation A site at `/[slug-a]/`
- Complete variation B site at `/[slug-b]/`
- Namespaced Tailwind color palettes (`va-*`, `vb-*`)
- Separate templ templates per variation with shared business data
- Scoped navigation (links stay within each variation's URL space)

**Example:**

```
/two-variation-site Henderson Bakery
```

Claude will ask for any missing inputs before writing code, present a structured summary for confirmation, then build the full site.

---

### `/qc`

Run a pre-deploy quality control check against the standard-tier site checklist.

**When to use:** The site is nearing completion and needs a final review before deployment. This is a read-only audit — it reports issues but does not fix them.

**Inputs required:**

1. **Client name** — used in the report header

**What it checks (7 sections):**

| Section | What it verifies |
|---|---|
| Project Structure | Required files exist (`main.go`, handlers, templates, Dockerfile, etc.) |
| Build Verification | `mage build` succeeds, `go vet` and `golangci-lint` pass |
| Functionality | Routing, contact form validation, Postmark integration, config safety |
| SEO | Unique titles, meta descriptions, heading hierarchy, Open Graph tags, robots.txt, sitemap |
| Accessibility | Semantic HTML, image alt text, form labels, keyboard navigation, skip link |
| Security | Security headers, CSRF protection, no hardcoded secrets, no localhost references |
| Deployment Readiness | `.env.example` complete, Docker config clean, all placeholder copy replaced |

**Output:** A structured report with tiered findings:

- **FAIL** — must fix before deploy
- **WARN** — should fix, does not block deploy
- **PASS** — requirement met

Final status: `READY`, `READY WITH WARNINGS`, or `NOT READY`.

**Example:**

```
/qc Henderson Bakery
```
