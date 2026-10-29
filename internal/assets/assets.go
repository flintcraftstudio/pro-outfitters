// Package assets serves the embedded static files under content-hashed URLs.
//
// Every file is hashed once at startup. Templates ask for a logical path and
// get the hashed URL: URL("css/site.css") is "/static/css/site.3f9a1c2e.css".
// A hashed URL names exact bytes, so it is served with a one-year immutable
// cache; changing a file changes its name, and the next page view fetches it.
// The plain path still serves, with a short cache, for anything that can't
// go through URL: old links, the social share image, crawlers.
package assets

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/flintcraftstudio/standard-template/web"
)

// hashLen is how many hex characters of the SHA-256 go into a filename.
const hashLen = 8

// Cache-Control values for the two kinds of static URL.
const (
	CacheHashed   = "public, max-age=31536000, immutable"
	CacheUnhashed = "public, max-age=300"
)

// Types Go's mime table doesn't know; without them ServeContent sniffs the
// body and sends text/plain.
var extraTypes = map[string]string{
	".webmanifest": "application/manifest+json",
}

// Set is a static file tree indexed by logical and hashed path.
type Set struct {
	fsys      fs.FS
	byLogical map[string]*file // "css/site.css"
	byHashed  map[string]*file // "css/site.3f9a1c2e.css"
}

type file struct {
	logical string
	hashed  string
	etag    string
	data    []byte // rewritten content (CSS); nil means stream from fsys
}

// cssURL matches url(/static/...) references in a stylesheet.
var cssURL = regexp.MustCompile(`url\((['"]?)/static/([^'")?#]+)`)

// New hashes every file in fsys. Stylesheets are hashed last, after their
// url(/static/...) references are rewritten to hashed URLs, so a change to a
// font or image a stylesheet uses changes the stylesheet's name too.
func New(fsys fs.FS) (*Set, error) {
	s := &Set{fsys: fsys, byLogical: map[string]*file{}, byHashed: map[string]*file{}}
	var stylesheets []string
	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if path.Ext(p) == ".css" {
			stylesheets = append(stylesheets, p)
			return nil
		}
		f, err := fsys.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			return err
		}
		s.add(&file{logical: p}, h.Sum(nil))
		return nil
	})
	if err != nil {
		return nil, err
	}
	for _, p := range stylesheets {
		b, err := fs.ReadFile(fsys, p)
		if err != nil {
			return nil, err
		}
		b = s.rewriteCSS(b)
		sum := sha256.Sum256(b)
		s.add(&file{logical: p, data: b}, sum[:])
	}
	return s, nil
}

func (s *Set) add(f *file, sum []byte) {
	h := hex.EncodeToString(sum)[:hashLen]
	ext := path.Ext(f.logical)
	f.hashed = strings.TrimSuffix(f.logical, ext) + "." + h + ext
	f.etag = `"` + h + `"`
	s.byLogical[f.logical] = f
	s.byHashed[f.hashed] = f
}

func (s *Set) rewriteCSS(b []byte) []byte {
	return cssURL.ReplaceAllFunc(b, func(m []byte) []byte {
		sub := cssURL.FindSubmatch(m)
		if f, ok := s.byLogical[string(sub[2])]; ok {
			return []byte("url(" + string(sub[1]) + "/static/" + f.hashed)
		}
		return m
	})
}

// URL returns the hashed URL for a path under static/ ("img/hero.avif"), or
// the plain /static/ URL if no such file exists.
func (s *Set) URL(logical string) string {
	if f, ok := s.byLogical[logical]; ok {
		return "/static/" + f.hashed
	}
	return "/static/" + logical
}

// Handler serves the set. Mount it under http.StripPrefix("/static/", ...).
// Hashed paths get CacheHashed, plain paths CacheUnhashed, and anything else
// (including directories) is a 404.
func (s *Set) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, "/")
		if f, ok := s.byHashed[p]; ok {
			s.serve(w, r, f, CacheHashed)
			return
		}
		if f, ok := s.byLogical[p]; ok {
			s.serve(w, r, f, CacheUnhashed)
			return
		}
		http.NotFound(w, r)
	})
}

// ServeFile serves one file at a URL of the caller's choosing (the favicons
// browsers request from the site root) with the given Cache-Control.
func (s *Set) ServeFile(w http.ResponseWriter, r *http.Request, logical, cacheControl string) {
	f, ok := s.byLogical[logical]
	if !ok {
		http.NotFound(w, r)
		return
	}
	s.serve(w, r, f, cacheControl)
}

func (s *Set) serve(w http.ResponseWriter, r *http.Request, f *file, cacheControl string) {
	w.Header().Set("Cache-Control", cacheControl)
	// The hash doubles as a strong ETag, so a revalidating client (plain
	// paths, root favicons) gets a 304 when nothing changed.
	w.Header().Set("ETag", f.etag)
	ext := path.Ext(f.logical)
	if mime.TypeByExtension(ext) == "" {
		if t, ok := extraTypes[ext]; ok {
			w.Header().Set("Content-Type", t)
		}
	}
	if f.data != nil {
		http.ServeContent(w, r, f.logical, time.Time{}, bytes.NewReader(f.data))
		return
	}
	fh, err := s.fsys.Open(f.logical)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer fh.Close()
	rs, ok := fh.(io.ReadSeeker)
	if !ok {
		b, err := io.ReadAll(fh)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		rs = bytes.NewReader(b)
	}
	http.ServeContent(w, r, f.logical, time.Time{}, rs)
}

// embedded is the Set over the binary's embedded web/static, built on first
// use. Hashing it reads every file once; call Handler at startup so that
// happens before the first request.
var embedded = sync.OnceValue(func() *Set {
	sub, err := fs.Sub(web.Static, "static")
	if err != nil {
		panic(err)
	}
	s, err := New(sub)
	if err != nil {
		panic("assets: hashing embedded files: " + err.Error())
	}
	return s
})

// URL returns the hashed URL for an embedded path under web/static.
func URL(logical string) string { return embedded().URL(logical) }

// Handler serves the embedded files. Mount under http.StripPrefix("/static/", ...).
func Handler() http.Handler { return embedded().Handler() }

// ServeFile serves one embedded file with the given Cache-Control.
func ServeFile(w http.ResponseWriter, r *http.Request, logical, cacheControl string) {
	embedded().ServeFile(w, r, logical, cacheControl)
}
