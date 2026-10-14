# Legacy redirects: reference

Internal reference for the old-site redirects. The client-facing summary is
[old-site-addresses.md](old-site-addresses.md).

- Code: `internal/handler/redirects.go` (the `legacyRedirects` map and the
  `LegacyRedirects` middleware), wired in `cmd/server/main.go`.
- Tests: `internal/handler/redirects_test.go`.
- Added in commit 9dd29e3, 2026-10-14.

## Source of the list

Taken from the old WordPress site on 2026-10-14:

- Its Yoast sitemap index, `https://prooutfitters.com/sitemap_index.xml`, and
  each sub-sitemap (post, page, portfolio_page, category, portfolio_category,
  testimonials_category, slides_category, author).
- The internal links on its homepage, which turned up nothing beyond the
  sitemaps except the feeds.

To recheck before cutover:

```bash
for s in post page portfolio_page category portfolio_category; do
  curl -sL "https://prooutfitters.com/$s-sitemap.xml" \
    | grep -o '<loc>[^<]*</loc>' | grep -v wp-content | sed 's/<\/\?loc>//g'
done
```

## How it behaves

- **GET and HEAD only.** Anything else passes through.
- **Matching:** the path is lowercased and the trailing slash is trimmed. The old site used a slash on every URL.
- **Response:** a 301 to the new path, with the query string dropped.
- **Order:** it runs ahead of `StripTrailingSlash`, so an old URL takes a single 301 instead of a 308 followed by a 301.
- **Not covered:** paths missing from the map fall through to the normal routes. For example, `/about/` gets a 308 to `/about`.

## The map

| Old path | New path |
|---|---|
| `/lodges-rivers/` | `/lodges` |
| `/contact-pro-outfitters/` | `/contact` |
| `/news-videos/` | `/` (pending decision 1) |
| `/trip-details-nforkcrossinglodge/` | `/lodges/north-fork-crossing` |
| `/smith-river-packing-list/` | `/packing-lists/smith-river` |
| `/kids-camp-what-to-bring/` | `/packing-lists/kids-camp` |
| `/bird-hunting-equipment-list/` | `/packing-lists/bird-hunting` |
| `/portfolio_page/` | `/lodges` |
| `/portfolio_page/smith-river-montana/` | `/smith-river` |
| `/portfolio_page/north-fork-crossing-lodge/` | `/lodges/north-fork-crossing` |
| `/portfolio_page/sharptail-lodge/` | `/lodges/sharptail` |
| `/portfolio_page/yurt-craig/` | `/lodges/yurt-at-craig` |
| `/portfolio_page/kids-camp-prooutfitters/` | `/kids-camp` |
| `/portfolio_page/blackfoot-river/` | `/blackfoot-river` (pending decision 2) |
| `/2020/08/08/sharptail-lodge-wingshooting/` | `/lodges/sharptail` |
| `/2020/08/08/community-escape-smith-river-float-trip/` | `/smith-river` |
| `/2020/07/16/sharptail-grouse-lodge/` | `/lodges/sharptail` |
| `/2021/12/02/yeti-sharptail-taupe-collection/` | `/lodges/sharptail` |
| `/category/bird-hunting/` | `/bird-hunting` |
| `/category/trout-fishing/` | `/fly-fishing` |
| `/category/smith-river-montana/` | `/smith-river` |
| `/category/sharptail-lodge/` | `/lodges/sharptail` |
| `/portfolio-category/featured/` | `/` |
| `/portfolio-category/lodges/` | `/lodges` |
| `/portfolio-category/rivers/` | `/fly-fishing` |

These paths match a new route already and need no entry: `/`, `/about/`.

These are deliberately left to 404: `/feed/`, `/comments/feed/`,
`/author/devon/`, `/testimonials-category/verified-guests/`,
`/slides-category/original-slider/`, and old `wp-content` image URLs.

## Pending decisions (client)

1. **News & Videos.** If Brandon wants something kept, give it a page and repoint `/news-videos`. Otherwise leave it going to `/`. This is request list item 21.
2. **Blackfoot River.** If it folds into North Fork, change three things:
   - point `/portfolio_page/blackfoot-river` at `/lodges/north-fork-crossing`;
   - delete the `/blackfoot-river` route and its page;
   - add a `/blackfoot-river` → `/lodges/north-fork-crossing` entry.

   This is request list item 20.

## Changing the map

1. Edit `legacyRedirects` in `internal/handler/redirects.go`. Keys are lowercase with no trailing slash.
2. If you add a new target route, add it to the `routes` set in `TestLegacyRedirectTargetsExist`. That list is kept by hand to match `cmd/server/main.go`.
3. Run `go test ./internal/handler/`.

## Launch-day checklist

- [ ] Recheck the old sitemaps (command above) for anything added since 2026-10-14.
- [ ] Set `APP_BASE_URL=https://prooutfitters.com` and `ALLOW_INDEXING=true` in the server's env file. Restart.
- [ ] Caddy: redirect `www.prooutfitters.com` → `https://prooutfitters.com{uri}` (301). Caddy already upgrades `http://` to `https://`.
- [ ] Point DNS at the VPS. Change only the web records; leave the MX and email records alone (request list items 27–28).
- [ ] Spot-check after DNS propagates:

  ```bash
  for u in /lodges-rivers/ /portfolio_page/sharptail-lodge/ /smith-river-packing-list/ \
           /2020/08/08/community-escape-smith-river-float-trip/ /contact-pro-outfitters/; do
    curl -s -o /dev/null -w "%{http_code} $u -> %{redirect_url}\n" "https://prooutfitters.com$u"
  done
  curl -sI http://www.prooutfitters.com/portfolio_page/yurt-craig/ | grep -i '^location'
  ```

- [ ] Submit `https://prooutfitters.com/sitemap.xml` in Google Search Console. Watch the coverage report for 404s over the first couple of weeks.
