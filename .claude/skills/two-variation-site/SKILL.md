---
name: two-variation-site
description: Build a Go/templ/Tailwind brochure website presenting two distinct brand/design directions for a client to compare side-by-side. Use when the user wants to create a two-variation comparison site, dual-brand presentation, or A/B design direction site.
user-invocable: true
argument-hint: "[business name or details]"
---

Build a brochure website that presents **two distinct brand directions** so the client can compare them and choose one. The site uses this repo's Go + templ + Tailwind stack with two complete design variations served from a single project.

## What You Are Building

- **Landing page** at `/` — a split-panel comparison hero where the client can choose either direction
- **Variation A** at `/[variation-a-slug]/` — complete site in the first design direction
- **Variation B** at `/[variation-b-slug]/` — complete site in the second design direction

Each variation is a fully navigable site with its own styling, typography, header, footer, and page layouts. They share the same business content and config.

## Architecture

This project follows the existing repo conventions: Go stdlib HTTP server, templ components, Tailwind CSS, Mage build tasks.

```
cmd/server/main.go                     # Entry point — routing, config init
internal/
  config/config.go                     # Business data + env vars
  handler/
    landing.go                         # GET / — split-panel comparison page
    variation.go                       # GET /{slug}/, GET /{slug}/about, etc.
    contact.go                         # GET/POST /contact (shared across variations)
  middleware/logging.go                # Request logging (existing)
  mail/postmark.go                     # Postmark email client (existing)
  view/
    shared.go                          # SiteName, tracking IDs, business data constants
    layout.templ                       # Base HTML wrapper (loads correct CSS per variation)
    landing.templ                      # Split-panel comparison landing page
    nav.templ                          # Navigation (accepts variation slug for scoped links)
    footer.templ                       # Footer (accepts variation slug for scoped links)
    [variation-a-slug]/
      home.templ                       # Variation A homepage
      about.templ                      # Variation A about page
      ...                              # Additional pages
    [variation-b-slug]/
      home.templ                       # Variation B homepage
      about.templ                      # Variation B about page
      ...
    contact.templ                      # Contact form (shared, styled per variation)
tailwind/
  tailwind.config.js                   # Color palettes for BOTH variations (namespaced)
  input.css                            # Font imports for both variations
web/static/
  css/site.css                         # Compiled Tailwind output (single file, both palettes)
  js/                                  # HTMX, Alpine.js, custom scripts
  images/
docs/
  brand-guide-a.md                     # Brand spec for variation A
  brand-guide-b.md                     # Brand spec for variation B
```

## Process

### Phase 0: Inputs Required

Before writing any code, you need these inputs from the user. Ask for anything missing:

1. **Business details** — name, address, phone, email, hours, social links, tagline
2. **Brand guide for each variation** — each must specify:
   - Color palette (hex values, roles, usage rules)
   - Typography stack (font families, sizes, weights for each role: headline, body, eyebrow/accent, UI/labels)
   - Voice/tone description
   - Layout personality (e.g., "airy whitespace" vs "dense color-blocking")
3. **Page copy** — approved text for each section of each page
4. **Images** — hero images, logos, team photos (or placeholders)
5. **Variation names/slugs** — short names for each direction (e.g., "warm" / "bold", "classic" / "modern")

### Phase 1: Analyze and Confirm

Audit the inputs and present a structured summary:
- Pages to build (with sections per page)
- Color/typography comparison table between variations
- Navigation structure
- Interactive features (contact form, map, etc.)

**Stop and wait for confirmation before writing code.**

### Phase 2: Tailwind Config

Extend `tailwind/tailwind.config.js` with **namespaced color palettes** for both variations. Each variation gets its own prefix (e.g., `va-*` and `vb-*`, or use the slug names like `warm-*` and `bold-*`).

```js
/** @type {import('tailwindcss').Config} */
const defaultTheme = require("tailwindcss/defaultTheme");

module.exports = {
  content: ["./internal/view/**/*.templ"],
  theme: {
    extend: {
      colors: {
        // Variation A palette
        va: {
          dark:    "#...",
          panel:   "#...",
          primary: "#...",
          // ... full palette
        },
        // Variation B palette
        vb: {
          dark:    "#...",
          panel:   "#...",
          primary: "#...",
          // ... full palette
        },
      },
      fontFamily: {
        // Variation A fonts
        "va-display": ['"Font A Display"', ...defaultTheme.fontFamily.serif],
        "va-body":    ['"Font A Body"', ...defaultTheme.fontFamily.sans],
        // Variation B fonts
        "vb-display": ['"Font B Display"', ...defaultTheme.fontFamily.sans],
        "vb-body":    ['"Font B Body"', ...defaultTheme.fontFamily.sans],
      },
    },
  },
  plugins: [],
};
```

Update `tailwind/input.css` to import fonts for both variations:

```css
@import url("https://fonts.googleapis.com/css2?family=...");

@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  ::selection {
    background: rgba(...);
    color: #...;
  }
}
```

### Phase 3: Config and Business Data

Add business data to `internal/view/shared.go` as package-level constants or variables. Follow the existing pattern where tracking IDs are set from config at startup.

```go
package view

const SiteName = "Business Name"

// Business info — set from config at startup or as constants
var (
    BusinessPhone   string
    BusinessEmail   string
    BusinessAddress string
    // ... hours, social links, etc.
)
```

For variation metadata, define a struct:

```go
type Variation struct {
    Slug        string
    Name        string
    Description string
    // Any variation-specific data needed by templates
}

var Variations = []Variation{
    {Slug: "warm", Name: "Warm", Description: "..."},
    {Slug: "bold", Name: "Bold", Description: "..."},
}
```

### Phase 4: Routing

Add routes in `cmd/server/main.go` following the existing stdlib `http.ServeMux` pattern:

```go
// Landing page
mux.Handle("GET /", handler.Landing())

// Variation pages — use path values (Go 1.22+ routing)
mux.Handle("GET /{slug}/", handler.VariationHome())
mux.Handle("GET /{slug}/about", handler.VariationPage("about"))
mux.Handle("GET /{slug}/services", handler.VariationPage("services"))
// ... additional pages

// Contact (shared)
mux.Handle("GET /{slug}/contact", handler.Contact())
mux.Handle("POST /{slug}/contact", handler.ContactSubmit(mailer, cfg.TurnstileSecretKey))

// Static files
mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
```

### Phase 5: Handlers

Create handlers following the existing pattern (function returning `http.HandlerFunc`):

```go
// internal/handler/landing.go
func Landing() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := view.LandingPage().Render(r.Context(), w); err != nil {
            slog.Error("render error", "err", err)
        }
    }
}

// internal/handler/variation.go
func VariationHome() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        slug := r.PathValue("slug")
        // Validate slug is a known variation
        if err := view.VariationHomePage(slug).Render(r.Context(), w); err != nil {
            slog.Error("render error", "err", err)
        }
    }
}
```

### Phase 6: Templates (templ)

Follow the existing templ patterns: `@Base(title, description, content)` composition, reusable components as templ functions.

**Layout** — modify `layout.templ` to accept an optional variation slug so it can conditionally set the body class or load variation-specific styles:

```go
templ Base(title, description, slug string, body templ.Component) {
    <!DOCTYPE html>
    <html lang="en" class={ variationRootClass(slug) }>
        <head>
            // ... existing head content
            <link rel="stylesheet" href="/static/css/site.css"/>
        </head>
        <body class={ variationBodyClass(slug) }>
            @Nav(slug)
            <main>
                @body
            </main>
            @Footer(slug)
        </body>
    </html>
}
```

**Landing page** — `landing.templ`:

```go
templ LandingPage() {
    @Base("Choose Your Direction", "Compare two design directions.", "", landingContent())
}

templ landingContent() {
    <section class="min-h-screen flex flex-col sm:flex-row">
        // Left panel — Variation A preview
        <a href="/warm/" class="flex-1 bg-va-dark p-12 ...">
            <h2 class="font-va-display ...">Warm</h2>
            // Color swatches, personality description, CTA
        </a>
        // Right panel — Variation B preview
        <a href="/bold/" class="flex-1 bg-vb-dark p-12 ...">
            <h2 class="font-vb-display ...">Bold</h2>
            // Color swatches, personality description, CTA
        </a>
    </section>
}
```

**Variation pages** — each variation gets its own directory under `internal/view/`:

```go
// internal/view/warm/home.templ
package warm

templ HomePage() {
    @view.Base("Home", "...", "warm", homeContent())
}

templ homeContent() {
    <section class="bg-va-dark text-va-cream ...">
        // Hero, services, CTA — all using va-* Tailwind classes
    </section>
}
```

**Nav and Footer** — accept slug to scope links:

```go
templ Nav(slug string) {
    <header>
        <nav>
            if slug != "" {
                <a href={ templ.SafeURL("/" + slug + "/") }>Home</a>
                <a href={ templ.SafeURL("/" + slug + "/about") }>About</a>
                <a href={ templ.SafeURL("/" + slug + "/contact") }>Contact</a>
            }
        </nav>
    </header>
}
```

Key rules:
- Navigation links within a variation MUST stay within that variation's URL space
- No cross-linking between variations (except back to `/` landing page)
- Use the variation's namespaced Tailwind classes (`va-*` vs `vb-*`)
- Content is the same across variations — only presentation differs
- Reuse existing component patterns (`eyebrow()`, `stat()`, `serviceCard()`)

### Phase 7: Build Pipeline

The existing Mage tasks handle everything. No changes needed to `magefile.go` or `Dockerfile`:

- `mage BuildCSS` compiles both variation palettes into a single `site.css`
- `mage BuildGo` runs `templ generate` then `go build`
- `mage Build` does both
- Docker multi-stage build works unchanged

Run `templ generate` after creating new `.templ` files.

### Phase 8: Accessibility Checklist

Every page must have:
- Skip-to-content link
- Semantic HTML (`<header>`, `<main>`, `<nav>`, `<section>`, `<footer>`)
- ARIA labels on interactive elements
- `focus-visible` outlines
- `prefers-reduced-motion` support
- Alt text on all images
- Sufficient color contrast (WCAG AA minimum)
- Keyboard-navigable menus

## Design Principles

1. **Config over code** — business data in `shared.go` or config, never hardcoded in templates
2. **Tailwind utility-first** — use Tailwind classes, no separate CSS files per variation
3. **Namespaced palettes** — each variation's colors prefixed in Tailwind config (e.g., `va-*`, `vb-*`)
4. **templ composition** — `@Base()` wraps everything, reusable components as templ functions
5. **Stdlib HTTP** — no router frameworks, use Go 1.22+ `http.ServeMux` path values
6. **Mage build tasks** — use existing `mage Build`, `mage Dev` workflow
7. **Accessibility first** — not an afterthought
8. **Security** — honeypot + Turnstile on forms, CORS validation

## What "Done" Looks Like

- `mage Build && ./bin/server` runs without errors
- `/` shows the split-panel comparison landing page
- `/[variation-a-slug]/` renders a complete, navigable site in variation A's style
- `/[variation-b-slug]/` renders a complete, navigable site in variation B's style
- Both variations use the same business data from `shared.go`/config
- Navigation within each variation stays within its URL space
- Mobile-responsive using Tailwind's responsive prefixes
- All accessibility requirements met
- `templ generate` produces no errors
- Docker build succeeds (`docker compose up`)

## Customization Guidance

When building the two variations:

1. **Write two brand guides first** — each should fully specify colors, typography, voice, and layout personality. The more specific and opinionated, the more distinct the variations will feel.

2. **Name the variations evocatively** — "warm" vs "bold" works better than "option-1" vs "option-2". The names should hint at the personality of each direction.

3. **Vary the right things**:
   - Color palette (not just accent colors — different background strategies)
   - Typography family (serif vs sans-serif, or two different sans-serifs)
   - Layout density (airy vs compact)
   - Visual texture (flat vs textured, curved vs angular)
   - Hero treatment (single image vs split panel, overlay text vs adjacent text)

4. **Keep the same things the same**:
   - All business data (shared.go / config)
   - Page copy / content
   - URL structure (just different root prefix)
   - Information architecture
   - Accessibility standards

5. **The landing page sells the comparison** — it should make each direction's personality immediately obvious through color, type, and layout differences visible in the panels.
