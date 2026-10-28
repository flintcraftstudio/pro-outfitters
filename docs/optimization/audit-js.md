# JS Audit

**Date:** 2026-10-28T16:25:00-06:00
**Baseline:** /docs/optimization/audit-baseline.md
**Action taken:** Report only (nothing to vendor: htmx and Alpine are already
local and minified; Turnstile and analytics must load from their own origins)

## Summary

| Metric                       | Before  | After  | Delta |
|------------------------------|---------|--------|-------|
| External library payload     | 28 KB   | 28 KB  | 0%    |
| Third-party CDN requests     | 1       | 1      | 0     |
| Unminified libraries         | 0       | 0      | 0     |
| Duplicate library loads      | 0       | 0      | 0     |
| Unused extensions/plugins    | 0       | 0*     | -     |

*No extensions or plugins are loaded. Two libraries load site-wide but are
used on few pages; see Findings.

Sizes are gzip transfer sizes from the baseline run (home page, mobile).

## Discovered Inventory

### External Libraries

| Library | Version | URL | Size | Minified | Defer | Async |
|---------|---------|-----|------|----------|-------|-------|
| Cloudflare Turnstile | v0 | https://challenges.cloudflare.com/turnstile/v0/api.js | 28 KB | yes | yes | yes |
| Google tag (gtag.js) | n/a | https://www.googletagmanager.com/gtag/js?id={GtagID} | ~90 KB (typical) | yes | no | yes |
| Meta pixel (fbevents.js, injected) | 2.0 | https://connect.facebook.net/en_US/fbevents.js | ~40 KB (typical) | yes | n/a | yes |

Turnstile loads whenever `TURNSTILE_SITE_KEY` is set (it is on staging).
gtag.js and the Meta pixel load only when `GTAG_ID` / `PIXEL_ID` are set;
neither is set on staging, so they were not in the baseline.

### Local Scripts

| File | Size | Minified | Defer |
|------|------|----------|-------|
| /static/js/htmx.min.js (htmx 2.0.4) | 50 KB (17.6 KB gz) | yes | yes |
| /static/js/alpine.min.js (Alpine 3.14.8) | 44 KB (17.2 KB gz) | yes | yes |
| /static/js/site.js | 3.6 KB (1.6 KB gz) | no | yes |

All three load from `layout.templ:65-74` on every public page. The admin
layout (`admin_layout.templ:12`) loads htmx separately; that is a different
page shell, not a duplicate.

### Inline Scripts ≥ 1 KB

| File:Line | Size | Classification |
|-----------|------|----------------|
| (none)    | -    | -              |

The inline blocks are all under 1 KB: the `.js` class flag
(`layout.templ:22`, 52 B), the gtag init (`layout.templ:80`, ~170 B) and the
Meta pixel init (`layout.templ:89`, ~550 B).

## Findings

### Unminified Loads

None. `site.js` is unminified but 1.6 KB gzipped; minifying it is not worth a
build step.

### CDN Dependencies

- Cloudflare Turnstile: must load from challenges.cloudflare.com; do not
  vendor. It is, however, loaded on every page (`layout.templ:56-58`) while
  the only widget is on `/contact` (`contact.templ:107`). The home page pays
  28 KB, ~24 KB of it unused (baseline `unused-javascript`), plus a third-
  party connection, for nothing. Load it on `/contact` only.

### Unused Extensions

None loaded.

### Libraries loaded where unused

- **htmx** (17.6 KB gz): used only by the contact form (`contact.templ:55`)
  and the sign-in form (`login.templ:22`), loaded on every public page.
  Loading it on those two pages only saves 17.6 KB and a request on the
  rest.
- **Alpine** (17.2 KB gz): used only by the mobile nav toggle
  (`nav.templ:4-44`: `x-data`, `@click`, `x-show`, `x-transition`). It is
  needed on every page, but the toggle is ~20 lines of vanilla JS in
  `site.js`; replacing it would drop Alpine's 17.2 KB, a request, and its
  ~80 ms long task on mobile. Optional: weigh against keeping Alpine for
  future interactivity.

### Duplicate Loads

None.

### Large Inline Scripts

None ≥ 1 KB.

### Missing defer/async

None. Every script is `defer` or `async`.

### Tailwind Play CDN

Not present.

### Bug: analytics IDs are not interpolated

`gtagInit` and `pixelInit` (`layout.templ:79-101`) write `'{ GtagID }'`
and `'{ PixelID }'` inside `<script>` blocks. templ does not interpolate
single-brace expressions in script content, so the generated page emits the
literal text: `gtag('config', '{ GtagID }')` and
`fbq('init', '{ PixelID }')` (confirmed in `layout_templ.go`). With
`GTAG_ID` / `PIXEL_ID` set, both tags load but register under a
non-existent ID and record nothing. Fix with templ's script interpolation,
which JSON-encodes the value: `gtag('config', {{ GtagID }});` and
`fbq('init', {{ PixelID }});`. Not live yet (neither ID is set), but it
would break analytics at launch.

### Main-thread work attributed to site.js

The baseline attributes 1.4 s of mobile main-thread time and a 593 ms long
task to `site.js`, but only 10 ms of it is script execution; the rest is
style and layout it forces. `forced-reflow-insight` pins it to the nav
script reading `hero.offsetHeight` (`site.js:19`, 291 ms of reflow at
startup; `site.js:21`, 52 ms on load). An IntersectionObserver on the hero
(with a -90px top root margin) gives the same transparent-to-paper switch
without ever reading layout, and drops the scroll and resize listeners.

## Vendored Files

None (see Action taken).

## Templates Requiring Migration

### Remove these tags

- [ ] `internal/view/layout.templ:56-58` — remove the site-wide Turnstile
      `<script>`; load it from the contact page instead.
- [ ] `internal/view/layout.templ:65` — remove the site-wide htmx
      `<script>`; load it on the contact and sign-in pages.

### Add these tags

- [ ] `internal/view/contact.templ` — add
      `<script src="https://challenges.cloudflare.com/turnstile/v0/api.js" async defer></script>`
      (when `TurnstileSiteKey != ""`) and
      `<script src="/static/js/htmx.min.js" defer></script>`.
- [ ] `internal/view/login.templ` — add
      `<script src="/static/js/htmx.min.js" defer></script>`.

### Remove unused

- [ ] Optional: replace Alpine in `nav.templ` with a vanilla toggle in
      `site.js`, then remove `layout.templ:66`.

### Extract inline

None.

### Other

- [ ] Fix analytics interpolation in `layout.templ:84` and `:98`
      (`{{ GtagID }}`, `{{ PixelID }}`).
- [ ] Rewrite the nav state in `site.js:6-31` to use an
      IntersectionObserver instead of `offsetHeight`.

## Notes

- Static files are embedded? No: `cmd/server/main.go` serves
  `web/static` from disk (`http.FileServer(http.Dir("web/static"))`), and
  the Dockerfile copies `web/` into the image. JS and CSS are served with
  `max-age=3600` and unhashed names; content hashing is for
  `/go-static-assets`.
- htmx must stay `defer`red after it moves into page bodies: the contact
  form's `hx-post` is processed on DOMContentLoaded either way.
- The `htmx:afterSwap` Turnstile re-render in `site.js:64-72` keeps working
  when both scripts load only on `/contact`.

## Applied

2026-10-28: the analytics interpolation fix, htmx and Turnstile scoped to the
pages with forms (`htmxScript` / `turnstileScript` in `layout.templ`, called
from `contactContent` and `loginContent`), and the nav state moved to an
IntersectionObserver. `internal/view/layout_test.go` guards the first two.
Alpine (optional item) left in place.
