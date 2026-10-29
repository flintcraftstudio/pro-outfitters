package assets

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func testSet(t *testing.T, files map[string]string) *Set {
	t.Helper()
	fsys := fstest.MapFS{}
	for name, body := range files {
		fsys[name] = &fstest.MapFile{Data: []byte(body)}
	}
	s, err := New(fsys)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func get(t *testing.T, h http.Handler, url string, header ...string) *http.Response {
	t.Helper()
	req := httptest.NewRequest("GET", url, nil)
	for i := 0; i+1 < len(header); i += 2 {
		req.Header.Set(header[i], header[i+1])
	}
	rec := httptest.NewRecorder()
	http.StripPrefix("/static/", h).ServeHTTP(rec, req)
	return rec.Result()
}

func TestURLHashesContent(t *testing.T) {
	a := testSet(t, map[string]string{"js/app.js": "one"})
	b := testSet(t, map[string]string{"js/app.js": "two"})
	ua, ub := a.URL("js/app.js"), b.URL("js/app.js")
	if !strings.HasPrefix(ua, "/static/js/app.") || !strings.HasSuffix(ua, ".js") || len(ua) != len("/static/js/app..js")+hashLen {
		t.Fatalf("URL = %q, want /static/js/app.<%d hex>.js", ua, hashLen)
	}
	if ua == ub {
		t.Errorf("different content, same URL %q", ua)
	}
	if got := a.URL("js/missing.js"); got != "/static/js/missing.js" {
		t.Errorf("unknown file URL = %q, want the plain path", got)
	}
}

func TestStylesheetPointsAtHashedFontsAndTracksThem(t *testing.T) {
	css := "@font-face{src:url(/static/fonts/a.woff2) format('woff2')} .x{background:url('/static/img/missing.png')}"
	a := testSet(t, map[string]string{"css/site.css": css, "fonts/a.woff2": "font v1"})
	b := testSet(t, map[string]string{"css/site.css": css, "fonts/a.woff2": "font v2"})

	body := readBody(t, get(t, a.Handler(), a.URL("css/site.css")))
	if want := "url(" + a.URL("fonts/a.woff2") + ")"; !strings.Contains(body, want) {
		t.Errorf("stylesheet does not reference the hashed font %s:\n%s", want, body)
	}
	if !strings.Contains(body, "url('/static/img/missing.png')") {
		t.Errorf("reference to a file that doesn't exist should be left alone:\n%s", body)
	}
	if a.URL("css/site.css") == b.URL("css/site.css") {
		t.Error("changing a font the stylesheet uses did not change the stylesheet's URL")
	}
}

func TestHandlerCachePolicy(t *testing.T) {
	s := testSet(t, map[string]string{"img/a.avif": "image", "site.webmanifest": "{}"})
	h := s.Handler()

	hashed := get(t, h, s.URL("img/a.avif"))
	if hashed.StatusCode != 200 || hashed.Header.Get("Cache-Control") != CacheHashed {
		t.Errorf("hashed: %d %q, want 200 %q", hashed.StatusCode, hashed.Header.Get("Cache-Control"), CacheHashed)
	}
	if ct := hashed.Header.Get("Content-Type"); ct != "image/avif" {
		t.Errorf("hashed Content-Type = %q", ct)
	}

	plain := get(t, h, "/static/img/a.avif")
	if plain.StatusCode != 200 || plain.Header.Get("Cache-Control") != CacheUnhashed {
		t.Errorf("plain: %d %q, want 200 %q", plain.StatusCode, plain.Header.Get("Cache-Control"), CacheUnhashed)
	}
	if again := get(t, h, "/static/img/a.avif", "If-None-Match", plain.Header.Get("ETag")); again.StatusCode != http.StatusNotModified {
		t.Errorf("revalidation with ETag = %d, want 304", again.StatusCode)
	}

	if ct := get(t, h, "/static/site.webmanifest").Header.Get("Content-Type"); ct != "application/manifest+json" {
		t.Errorf("webmanifest Content-Type = %q", ct)
	}
	for _, p := range []string{"/static/img/", "/static/img/nope.avif", "/static/img/a.deadbeef.avif"} {
		if code := get(t, h, p).StatusCode; code != http.StatusNotFound {
			t.Errorf("%s = %d, want 404", p, code)
		}
	}
}

// The real embedded tree: every URL the site can hand out must resolve.
func TestEmbeddedStylesheetFontsResolve(t *testing.T) {
	css := readBody(t, get(t, Handler(), URL("css/site.css")))
	if !strings.Contains(css, "@font-face") {
		t.Skip("site.css not built; run mage buildCSS before go test")
	}
	for _, m := range cssURL.FindAllStringSubmatch(css, -1) {
		if code := get(t, Handler(), "/static/"+m[2]).StatusCode; code != 200 {
			t.Errorf("site.css references /static/%s: %d", m[2], code)
		}
	}
}

func readBody(t *testing.T, r *http.Response) string {
	t.Helper()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
