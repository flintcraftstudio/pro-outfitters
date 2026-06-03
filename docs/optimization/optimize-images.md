# Image Optimization

**Date:** 2026-06-03
**Baseline:** none captured (`/audit-baseline` was not run first)

## Summary

| Metric                         | Before  | After   | Delta |
|--------------------------------|---------|---------|-------|
| Referenced image bytes (JPG)   | 6867 KB | —       | —     |
| Same images, AVIF (base size)  | —       | 3312 KB | −52%  |
| Same images, WebP (base size)  | —       | 4725 KB | −32%  |
| Referenced source images       | 22      | 22      | —     |
| Full treatment (responsive)    | —       | 7       | —     |
| Content images (single size)   | —       | 15      | —     |
| Skipped (unused)               | —       | 3       | —     |

Modern browsers fetch AVIF; the figures above are the like-for-like base-size
comparison. The seven full-bleed images additionally ship responsive widths, so
a phone fetches far less than even the AVIF base size (e.g. the home hero serves
a 63 KB 800w AVIF instead of the 525 KB original JPG).

**Homepage** (the page PageSpeed measured): the eight content/hero JPGs totalled
~2005 KB and now serve ~964 KB as base AVIF, and less again where a responsive
width is picked.

## What changed

Unlike the default skill flow, the `<img>` → helper migration was **completed**,
not left as a TODO. All ten content `<img>` tags across the site now render
through the helper.

### Helper

Written to `internal/view/picture.templ` (+ `internal/view/picture.go`), in the
existing `view` package — this project keeps all templates in one package rather
than a `components/` dir.

- `picture(src, alt, class)` — content images. AVIF + WebP `<source>` siblings
  with the original JPG as `<img>` fallback. `loading="lazy"`, `decoding="async"`.
- `bleedPicture(src, alt, class, eager)` — full-bleed images with responsive
  AVIF/WebP/JPG srcsets at `sizes="100vw"`. `eager=true` sets `fetchpriority="high"`
  + `loading="eager"` for a page's LCP hero; `eager=false` lazy-loads
  below-the-fold bands.

Both wrap `<picture class="contents">` (display:contents) so the layout classes
on the inner `<img>` resolve against the real parent (grid cell / positioned
section) exactly as the bare `<img>` did — no layout change.

`src` is the full `/static/img/foo.jpg` path; the AVIF/WebP/width siblings are
derived in `picture.go` (`imgVariant`, `imgSrcset`).

### Conversions

- AVIF: `avifenc --speed 6 --min 20 --max 30`
- WebP: `cwebp -q 82 -metadata none`
- Originals: EXIF stripped losslessly with `jpegtran -copy none -optimize`
  (originals kept, bytes unchanged as fallback).
- Responsive widths 400/800/1200/1600 generated for the seven full-bleed images
  (every source is ≥1752px wide, so no upscaling and all four widths exist).

## Classification

**Full treatment — responsive srcset (7):** `hero-fishing`, `fishing-alt1`,
`fishing-alt2`, `nfork-hero`, `sharptail-hero`, `yurt-hero`, `smith-hero`.
These are the `100vw` full-bleed heroes/bands (home hero, the two dark CTA bands,
and the four interior-page heroes flowing through `detailHero`).

**Content — single-size AVIF/WebP (15):** `story-flybox`, `exp-hunting`,
`steward-smith`, `about-brandon`, `lodge-northfork`, `lodge-sharptail`,
`lodge-yurt`, `nfork-dining`, `nfork-pond`, `sharptail-dogs`, `sharptail-yurts`,
`smith-trout`, `smith-wading`, `yurt-bedroom`, `yurt-wildlife`.

## Templates migrated

- [x] `internal/view/home.templ` (6: hero, story-flybox, 2 experience cards,
      lodge cards loop, steward-smith, fishing-alt2 band)
- [x] `internal/view/subpages.templ` (3: `detailHero`, `gallery`,
      `closingInquiry` — shared by every interior page)
- [x] `internal/view/about.templ` (1: about-brandon)

`detailHero`/`gallery`/`closingInquiry` are shared components, so the lodge
detail, Smith River, and About pages were all covered by the subpages edits.

## Per-image results (referenced originals)

| Source            | JPG    | AVIF  | WebP  |
|-------------------|--------|-------|-------|
| hero-fishing      | 525 KB | 270 KB| 407 KB|
| smith-hero        | 645 KB | —     | —     |
| fishing-alt1      | 535 KB | —     | —     |
| smith-wading      | 483 KB | —     | —     |
| exp-hunting       | 502 KB | 216 KB| 353 KB|
| yurt-wildlife     | 457 KB | —     | —     |
| sharptail-hero    | 446 KB | —     | —     |
| fishing-alt2      | 411 KB | 167 KB| 248 KB|
| (… 14 more)       |        |       |       |

Totals: 6867 KB JPG → 3312 KB AVIF base (−52%), plus responsive widths for the
seven full-bleed images.

## Flagged for manual review

- `logo-primary-boxed.png` (30 KB), `logo-secondary-gold.png` (18 KB),
  `trout-unlimited.jpg` (11 KB) — **not referenced** anywhere in the codebase.
  Left untouched; delete if confirmed dead, or wire up + optimize if intended.

## Broken references

- (none)

## Notes

- Go's static `FileServer` serves `.avif`/`.webp` with correct MIME types
  (verified: `Content-Type: image/avif`). Cache-Control on `/static/` is set by
  `cacheStatic` in `cmd/server/main.go` (images: 1 year).
- New image files are served from `web/static/img` on disk and copied into the
  runtime image by the Dockerfile (`COPY --from=build /src/web /web`).
