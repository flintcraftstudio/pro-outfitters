package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type csrfContextKey struct{}
type csrfFieldKey struct{}

// CSRFConfig holds all configuration for the CSRF middleware.
type CSRFConfig struct {
	// Secret is the HMAC signing key. Must be at least 32 bytes.
	// Load from environment — never hardcode. A short or missing
	// secret causes a panic at middleware init.
	Secret []byte

	// CookieName is the name of the CSRF cookie. Defaults to "__csrf".
	CookieName string

	// FieldName is the hidden form field name checked on mutation requests.
	// Defaults to "csrf_token".
	FieldName string

	// HeaderName is checked before the form field on mutation requests.
	// Pair with htmx via hx-headers='{"X-CSRF-Token": "..."}' or a global
	// htmx config. Defaults to "X-CSRF-Token".
	HeaderName string

	// CookiePath scopes the cookie. Defaults to "/".
	CookiePath string

	// SameSite controls the SameSite cookie attribute. Defaults to
	// http.SameSiteLaxMode.
	SameSite http.SameSite

	// InsecureDev disables the Secure cookie flag. Must be explicitly
	// set to true for local HTTP development.
	InsecureDev bool

	// ErrorHandler writes the full response on CSRF validation failure.
	// Defaults to plain-text 403 Forbidden if nil.
	ErrorHandler func(w http.ResponseWriter, r *http.Request)
}

func newCSRFConfig(c CSRFConfig) CSRFConfig {
	if len(c.Secret) < 32 {
		panic(fmt.Sprintf(
			"csrf: Secret must be at least 32 bytes, got %d — "+
				"load from environment with os.Getenv and never hardcode",
			len(c.Secret),
		))
	}
	if c.CookieName == "" {
		c.CookieName = "__csrf"
	}
	if c.FieldName == "" {
		c.FieldName = "csrf_token"
	}
	if c.HeaderName == "" {
		c.HeaderName = "X-CSRF-Token"
	}
	if c.CookiePath == "" {
		c.CookiePath = "/"
	}
	if c.SameSite == http.SameSiteDefaultMode {
		c.SameSite = http.SameSiteLaxMode
	}
	if c.ErrorHandler == nil {
		c.ErrorHandler = func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		}
	}
	return c
}

var safeMethods = map[string]bool{
	http.MethodGet:     true,
	http.MethodHead:    true,
	http.MethodOptions: true,
	http.MethodTrace:   true,
}

// CSRF returns middleware that issues and validates CSRF tokens using
// the signed double-submit cookie pattern.
func CSRF(config CSRFConfig) func(http.Handler) http.Handler {
	c := newCSRFConfig(config)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := ""
			cookieVal := ""
			if cookie, err := r.Cookie(c.CookieName); err == nil {
				parts := strings.SplitN(cookie.Value, ".", 2)
				if len(parts) == 2 {
					token = parts[0]
					cookieVal = cookie.Value
				}
			}
			if token == "" {
				token = generateCSRFToken()
				cookieVal = cookieValue(c.Secret, token)
			}

			http.SetCookie(w, &http.Cookie{
				Name:     c.CookieName,
				Value:    cookieVal,
				Path:     c.CookiePath,
				HttpOnly: true,
				Secure:   !c.InsecureDev,
				SameSite: c.SameSite,
			})

			ctx := context.WithValue(r.Context(), csrfContextKey{}, token)
			ctx = context.WithValue(ctx, csrfFieldKey{}, c.FieldName)
			r = r.WithContext(ctx)

			if safeMethods[r.Method] {
				next.ServeHTTP(w, r)
				return
			}

			submitted := r.Header.Get(c.HeaderName)
			if submitted == "" {
				if err := r.ParseForm(); err != nil {
					c.ErrorHandler(w, r)
					return
				}
				submitted = r.FormValue(c.FieldName)
			}

			if submitted == "" {
				c.ErrorHandler(w, r)
				return
			}

			if !verifyCSRFToken(c.Secret, cookieVal, submitted) {
				c.ErrorHandler(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Token returns the raw CSRF token from the request context.
// Returns an empty string if the middleware has not been applied.
func Token(r *http.Request) string {
	token, _ := r.Context().Value(csrfContextKey{}).(string)
	return token
}

// TemplateField returns an html/template.HTML snippet containing a
// hidden input field populated with the CSRF token.
func TemplateField(r *http.Request) template.HTML {
	token := Token(r)
	if token == "" {
		return ""
	}
	fieldName, _ := r.Context().Value(csrfFieldKey{}).(string)
	if fieldName == "" {
		fieldName = "csrf_token"
	}
	return template.HTML(fmt.Sprintf(
		`<input type="hidden" name="%s" value="%s">`,
		template.HTMLEscapeString(fieldName),
		template.HTMLEscapeString(token),
	))
}

func generateCSRFToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("csrf: failed to read from crypto/rand: %v", err))
	}
	return base64.URLEncoding.EncodeToString(b)
}

func signToken(secret []byte, token string) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(token))
	return base64.URLEncoding.EncodeToString(mac.Sum(nil))
}

func cookieValue(secret []byte, token string) string {
	return token + "." + signToken(secret, token)
}

func verifyCSRFToken(secret []byte, cookieVal, submitted string) bool {
	expected := cookieValue(secret, submitted)
	return hmac.Equal([]byte(expected), []byte(cookieVal))
}
