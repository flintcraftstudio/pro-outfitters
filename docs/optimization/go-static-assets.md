# Go Static Assets

**Date:** 2026-10-29T10:30:00-06:00
**Baseline:** /docs/optimization/audit-baseline.md
**Action taken:** Audit + implemented directly (at the client team's request,
in place of a reference handler file)

## Summary

| Check | Status |
|-------|--------|
| Assets served via embed.FS | ✓ |
| Content hashing | ✓ |
| Cache-Control on assets | ✓ |
| Template references via helper | ✓ |
| HTML Cache-Control set | ✓ |

## Current State

- **Static handler:** `cmd/server/main.go` — `mux.Handle("GET /static/", http.StripPrefix("/static/", assets.Handler()))`
- **Asset source:** `http.FS`-style embedded tree (`web.Static`, via `fs.Sub(web.Static, "static")`)
- **embed.FS variable:** `web.Static` (`web/static.go`, `//go:embed static`)
- **Template references:** 351 unique `/static/` URLs across 15 pages, all hashed (0 unhashed)
- **Asset helper:** `assets.URL(logical)`; in templates `assetPath("/static/...")` and the `picture` / `sizedPicture` / `bleedPicture` helpers

## Findings

State before this pass:

1. **Served from disk, not embedded.** `http.FileServer(http.Dir("web/static"))`;
   the Dockerfile copied `web/` into the runtime image.
2. **No content hashing.** Every asset kept a fixed name.
3. **Unsafe Cache-Control under fixed names.** Fonts were
   `max-age=31536000, immutable` and images `max-age=31536000`, so a changed
   font or photo never reached a returning visitor. This had already
   happened: the trimmed fonts and the 11 photos re-rendered from originals
   on 2026-10-28 kept their names. CSS and JS got `max-age=3600`, so a deploy
   reached visitors up to an hour late (seen with the reviews-marquee CSS fix).
4. **Hardcoded references.** CSS, JS and font preloads in `layout.templ` and
   `admin_layout.templ` were bare `/static/...` strings; `site.css` referenced
   the six web fonts by fixed path. Images already went through
   `imgVariant` / `imgSrcset`.
5. **No HTML Cache-Control.** Pages were left to browser heuristics.
6. **Share image wrong.** `og:image` pointed at `hero-fishing-1200w.jpg`
   and declared 1200×673; since the home hero was re-rendered from its
   original that file is 1200×800.
7. **Weight.** 113 MB in `web/static`: 9 unreferenced files (~0.7 MB) and
   136 JPG width variants (25 MB) used only as the fallback srcset for
   browsers with neither AVIF nor WebP.

## Recommended Handler

Implemented in the project rather than generated as a reference file:

- `web/static.go` — embeds `web/static`.
- `internal/assets/assets.go` — `New(fs.FS)` hashes every file (SHA-256,
  8 hex chars) into `<name>.<hash>.<ext>`. Stylesheets are hashed last,
  after their `url(/static/...)` references are rewritten to hashed URLs,
  so a font change also renames the CSS. `Handler()` serves hashed paths
  with `public, max-age=31536000, immutable` and plain paths with
  `public, max-age=300`; both carry the hash as a strong ETag (304 on
  revalidation). Directories and unknown names 404. `ServeFile` serves the
  root-path favicons (`public, max-age=604800`).
- `internal/assets/assets_test.go` — hashing, CSS rewrite and tracking,
  cache policy, ETag 304, MIME for `.webmanifest`, 404s, and that every
  font the real `site.css` references resolves. Mutation-checked: disabling
  the CSS rewrite fails two tests.

## Migration TODO

- [x] Static handler embedded and registered (`cmd/server/main.go`)
- [x] Embed path: `web/static.go` beside `web/static`
- [x] Replaced `http.FileServer` + `cacheStatic`; favicons served from the
      embedded tree
- [x] Templ files use hashed URLs: `picture.go` / `picture.templ`
      (`assetPath`, `imgVariant`, `imgSrcset`), `layout.templ` (CSS, JS,
      font preloads, `og:image`, `twitter:image`), `admin_layout.templ`
- [x] HTML responses `Cache-Control: no-cache` (`pageCacheControl`
      middleware; static and favicon handlers override it)
- [x] Tested: standalone binary run from an empty directory served all 351
      referenced URLs as 200 + immutable; plain path 300 s + 304 on
      revalidation; Docker image (125 MB, no `web/`) smoke-tested; the CSS
      hash was identical in both builds.

Also done in this pass:

- Deleted the 9 unreferenced files (`lodge-sharptail.*`, `story-flybox.*`,
  `logo-primary-boxed.png`, `logo-secondary-gold.png`,
  `trout-unlimited.jpg`).
- Dropped the 136 JPG width variants; each image keeps its base JPG as the
  single `<img>` fallback, with AVIF/WebP srcsets on the `<source>`s.
  `web/static` is now 87 MB.
- New share image `img/og-image.jpg` at 1200×630, `og:image:height` fixed.
- Dockerfile runtime stage no longer copies `web/`.
- CLAUDE.md: deploy note and a gotcha on using `assetPath`.

## Notes

- Binary is 104 MB with the assets embedded; the image is 125 MB (it carried
  the same files on disk before).
- Hashing all files takes ~0.2 s at startup; `assets.Handler()` is called
  while building the mux so it happens before the first request.
- Dev: the files are compiled in, so CSS, image or font changes need a
  restart (`mage dev` restarts; build CSS before `go build`).
- Visitors who loaded staging before this deploy hold the old fonts and
  re-rendered photos under their unhashed names for up to a year. Pages now
  request hashed names, so they fetch the new files regardless.
