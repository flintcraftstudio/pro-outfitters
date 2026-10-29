# Caddy Config Review

**Date:** 2026-10-29T13:40:00-06:00
**Baseline:** /docs/optimization/audit-baseline.md
**Caddyfile audited:** the `prooutfitters.flintcraft.studio` block from the
server's Caddyfile (pasted; it isn't in the repo), plus live responses
from https://prooutfitters.flintcraft.studio/

## Summary

| Check | Status |
|-------|--------|
| Compression (zstd + gzip) | ✓ |
| Cache-Control pass-through | ✓ |
| HSTS | ✓ |
| X-Content-Type-Options | ✓ |
| X-Frame-Options | ✓ |
| Referrer-Policy | ✓ |
| HTTP/3 enabled | ✓ |
| Access logging | ✓ (added 2026-10-29) |

## Audited block

```caddyfile
prooutfitters.flintcraft.studio {
    encode zstd gzip
    header {
        Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"
        X-Content-Type-Options "nosniff"
        X-Frame-Options "SAMEORIGIN"
        Referrer-Policy "strict-origin-when-cross-origin"
        -Server
    }
    reverse_proxy localhost:8106
}
```

No `tls` directive (automatic HTTPS), no `header_down` and no
`file_server`: Go serves everything, static files included.

## Findings

Verified live with curl:

- **Compression:** `/` comes back zstd when it's offered, and
  `site.css` comes back gzip when only gzip is offered.
  `Vary: Accept-Encoding` is set.
- **Cache-Control pass-through:** these values reach the client unchanged:
  - `/`: `no-cache`
  - hashed static files: `public, max-age=31536000, immutable`
  - plain static paths: `public, max-age=300`
  - favicons: `public, max-age=604800`
- **Security headers:** all four are present on 200s and on the 404 page.
  The `Server` header is stripped.
- **HTTP to HTTPS:** a request on port 80 gets a 308 redirect to https.
- **HTTP/3:** responses carry `alt-svc: h3=":443"`, and UDP 443 is open.
  The local curl has no HTTP/3 support, so the protocol itself wasn't
  exercised.

Issues:

1. **No access log** (the one ✗). Caddy only logs errors to the journal by
   default, so after launch there is no record of requests: no 404 trail and no way to trace a form submission. The snippet adds a JSON log that rotates itself.
2. **`Via: 1.1 Caddy` is still sent.** `reverse_proxy` adds it, and `-Server`
   doesn't remove it. It's cosmetic; add `-Via` next to `-Server`.
3. **HSTS `includeSubDomains; preload` needs a decision before the
   production domain moves over.**
   - On the staging subdomain it's harmless: `preload` means nothing below
     the apex.
   - On the client's own domain, `includeSubDomains` makes browsers
     refuse plain HTTP on every subdomain for a year. That would break
     anything still HTTP-only, such as old webmail or a booking tool.
   - Being on the preload list is very hard to undo.
   - Recommendation: on production, use
     `max-age=31536000; includeSubDomains` once the client's subdomains
     are confirmed HTTPS. Otherwise drop `includeSubDomains`. Don't submit
     the domain to the preload list.
4. **The app's port is published on all interfaces (this is outside Caddy).**
   - `provision.sh` writes `ports: - "${PORT}:8080"`, which binds
     `0.0.0.0:8106`.
   - Docker bypasses ufw, so the app could be reached directly. That would
     skip TLS and every header above.
   - Resolved on the server: `/opt/pro-outfitters/docker-compose.yml`
     publishes `"127.0.0.1:8106:8080"`, and `:8106` times out from
     outside.
   - Still open in the template: `provision.sh` should write
     `"127.0.0.1:${PORT}:8080"` so the next provisioned site starts out
     bound to loopback. It has uncommitted local edits, so it was left
     untouched.
5. **No Content-Security-Policy (informational).** A policy would have to
   allow:
   - Google Analytics / gtag (googletagmanager.com, google-analytics.com)
   - the Meta pixel (connect.facebook.net, facebook.com)
   - Cloudflare Turnstile (challenges.cloudflare.com, frames included)
   - the inline analytics init scripts and JSON-LD, which need nonces or
     hashes

   Start it as `Content-Security-Policy-Report-Only` once the analytics
   setup for launch is final.
6. Added `Permissions-Policy "camera=(), microphone=(), geolocation=()"`.
   This is optional: the site uses none of them.

## Recommended Snippet

Generated at:
`/tmp/claude-1000/-home-dukerupert-repos-clients-pro-outfitters/d02a1de3-d19b-41b7-ad56-6819fe8f7f64/scratchpad/Caddyfile-recommended`

It is the live block with `-Via`, `Permissions-Policy` and the access log
added.

## Migration TODO

- [x] `encode zstd gzip`
- [x] Security headers block
- [x] Cache-Control passes through from Go (verified live)
- [x] Add `log { output file /var/log/caddy/pro-outfitters.log … format json }`
      (needed `chown caddy:caddy /var/log/caddy`)
- [x] Add `-Via` and `Permissions-Policy` to the header block
- [x] Server binds the app to loopback (`127.0.0.1:8106:8080`)
- [ ] `provision.sh` template: publish `"127.0.0.1:${PORT}:8080"`
- [x] `www.prooutfitters.com` block redirecting to the apex (the
      certificate comes once DNS points here)
- [ ] `prooutfitters.com` is already in the site block: set HSTS to
      `max-age=31536000; includeSubDomains`, without `preload`, before DNS
      moves; drop `includeSubDomains` if any subdomain is HTTP-only
- [ ] (Optional) CSP in Report-Only mode once analytics are final

## Verification

After editing the Caddyfile on the server:

1. `caddy fmt --overwrite /etc/caddy/Caddyfile`
2. `caddy validate --config /etc/caddy/Caddyfile`
3. `sudo systemctl reload caddy`
4. `curl -sI https://prooutfitters.flintcraft.studio/`: expect no `via`;
   `permissions-policy` present if added
5. `sudo tail -1 /var/log/caddy/pro-outfitters.log`: expect one JSON line
   per request (the caddy user needs write access to `/var/log/caddy`,
   which the Debian package gives it)

## Notes

- The `encode` output shows both zstd and gzip, so the Caddy version is
  2.6 or later.
- `x-request-id` on responses comes from the Go logging middleware and is
  fine to keep.
