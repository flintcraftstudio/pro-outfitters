# Baseline Audit

**URL audited:** https://prooutfitters.flintcraft.studio/
**Date:** 2026-10-28T15:54:55-06:00
**Lighthouse version:** 13.3.0

> This report is the "before" snapshot for the pre-launch optimization
> pass. It is read by `/launch-checklist` at the end of the workflow to
> compute the optimization delta. Do not edit by hand — re-run
> `/audit-baseline` if you need a fresh snapshot.

---

## Lighthouse Scores

| Category        | Mobile | Desktop |
|-----------------|--------|---------|
| Performance     | 62     | 99      |
| Accessibility   | 95     | 95      |
| Best Practices  | 100    | 100     |
| SEO             | 66     | 66      |

## Core Web Vitals (Mobile)

| Metric | Value | Threshold (Good) | Status |
|--------|-------|------------------|--------|
| LCP    | 4.4s  | < 2.5s           | ✗      |
| INP    | —     | < 200ms          | field data only |
| CLS    | 0.00  | < 0.1            | ✓      |

## Supporting Metrics (Mobile)

| Metric              | Value |
|---------------------|-------|
| First Contentful Paint | 2.4s |
| Time to First Byte     | 42ms |
| Speed Index            | 5.1s |
| Total Blocking Time    | 603ms |

## Page Weight

**Total transfer size:** 843 KB across 17 requests

| Resource Type | Count | Transfer Size | Resource Size |
|---------------|-------|---------------|---------------|
| Document      | 1     | 10 KB         | 74 KB         |
| Stylesheet    | 1     | 9 KB          | 37 KB         |
| Script        | 4     | 63 KB         | 182 KB        |
| Image         | 5     | 456 KB        | 456 KB        |
| Font          | 3     | 301 KB        | 301 KB        |
| Other         | 3     | 4 KB          | 15 KB         |

## Largest Resources

1. `https://prooutfitters.flintcraft.studio/static/img/exp-fishing.avif` — 156 KB
2. `https://prooutfitters.flintcraft.studio/static/fonts/newsreader-latin-italic.woff2` — 144 KB
3. `https://prooutfitters.flintcraft.studio/static/img/exp-hunting-hun.avif` — 140 KB

---

## Notes

- Audited the staging deploy of commit f4dbed3 (Caddy, HTTP/2, gzip and
  HSTS active; HTML transferred gzip-compressed). Home page only.
- Page weight and largest resources above are from the mobile run. Desktop
  transfers 1,406 KB across 20 requests (8 images, 1,019 KB); its largest
  are `hero-fishing-1600w.avif` (269 KB), `sharptail-yurts.avif` (163 KB;
  the Sharptail lodge card) and `exp-fishing.avif` (156 KB).
- The mobile run was a full-category run rather than `--preset=perf`, so
  that mobile Accessibility, Best Practices and SEO scores exist for the
  table.
- INP is null: synthetic runs have no interaction data.
- **SEO 66 is expected on staging.** The only failing check is
  `is-crawlable`: staging runs with `ALLOW_INDEXING` unset, so every page is
  `noindex` and robots.txt disallows all. Re-check after launch config.
- **Accessibility 95, two checks failing:**
  - `color-contrast` on the home closing band ("Tell us what kind of
    week…"): axe resolved the text (#f3eee5) against the page background
    (#faf7f1), not the photo, because the band's only dark backing is the
    lazy-loaded image and an overlay behind it. Giving the band a dark
    background colour would make the contrast hold before the image
    loads and satisfy axe.
  - `label-content-name-mismatch` on the header logo link: its
    `aria-label` ("Pro Outfitters, established 1970") does not contain the
    visible text ("PRO OUTFITTERS EST. 1970"), which breaks voice control
    (WCAG 2.5.3).
- **Mobile LCP 4.4s / Speed Index 5.1s** — the obvious contributors:
  - Fonts are 301 KB, 280 KB of it Newsreader (normal 132 KB, italic
    147 KB). Subsetting would cut most of it.
  - The two home experience-card images ship single-size (156 KB and
    140 KB AVIF at 1066–1200px) to phones; Lighthouse estimates 206 KB of
    savings from responsive sizes.
  - `site.css` is render-blocking (~190 ms estimated).
- **TBT 603ms (mobile):** `site.js` shows 1.4s of main-thread time, almost
  none of it scripting; `forced-reflow-insight` also fails. Likely the
  reveal/marquee animations and the nav's `offsetHeight` reads. For
  `/audit-js`.
- Cloudflare Turnstile's `api.js` (28 KB, ~24 KB unused) loads on the home
  page, which has no form.
- htmx (17.6 KB) + Alpine (17.2 KB) transfer 35 KB combined, so both are
  minified. They, `site.css` and `site.js` are served with
  `max-age=3600`; the filenames aren't content-hashed. For
  `/go-static-assets`.
- Staging page titles end "| pro-outfitters" (lowercase): staging's
  `APP_NAME` env var should be "Pro Outfitters".
