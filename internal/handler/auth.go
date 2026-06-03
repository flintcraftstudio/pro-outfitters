package handler

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/flintcraftstudio/standard-template/internal/email"
	"github.com/flintcraftstudio/standard-template/internal/session"
	"github.com/flintcraftstudio/standard-template/internal/store"
	"github.com/flintcraftstudio/standard-template/internal/view"
)

// magicLinkTTL is the lifetime of a requested sign-in link. 15 minutes
// is long enough to survive the inbox → click roundtrip on a slow email
// relay, short enough that a leaked URL has little blast radius.
const magicLinkTTL = 15 * time.Minute

// AuthDeps bundles everything the auth handlers reach for.
type AuthDeps struct {
	Store   *store.Store
	Mailer  email.Sender
	BaseURL string
	AppName string
}

// LoginPage handles GET /login and renders the sign-in form. Already-
// authenticated visitors are bounced to /admin.
func LoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if session.FromContext(r.Context()) != nil {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
		if err := view.LoginPage("", "").Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}

// RequestMagicLink handles POST /login. Issues a single-use sign-in
// token and emails a link. Always renders the confirmation page whether
// or not the email maps to a real admin user — prevents enumeration.
func RequestMagicLink(deps AuthDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		addr := strings.TrimSpace(strings.ToLower(r.FormValue("email")))
		if addr == "" {
			if err := view.LoginForm("Enter your email.", "").Render(r.Context(), w); err != nil {
				slog.Error("render error", "err", err)
			}
			return
		}

		// Do the lookup + send in the happy path. Unknown emails fall
		// through to the confirmation render without any DB write or
		// outbound request.
		if userID, _, err := deps.Store.GetUserByEmail(r.Context(), addr); err == nil {
			rawToken, err := newRawToken()
			if err != nil {
				slog.Error("token gen", "err", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			expiresAt := time.Now().Add(magicLinkTTL)
			if err := deps.Store.CreateLoginToken(r.Context(), userID, rawToken, expiresAt); err != nil {
				slog.Error("create login token", "err", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			link := fmt.Sprintf("%s/login/magic/%s", strings.TrimRight(deps.BaseURL, "/"), rawToken)
			msg := email.Message{
				To:      addr,
				Subject: fmt.Sprintf("Sign in to %s", deps.AppName),
				Text:    magicLinkText(deps.AppName, link),
				HTML:    magicLinkHTML(deps.AppName, link),
			}
			if err := deps.Mailer.Send(r.Context(), msg); err != nil {
				// Log but don't leak mail-config errors. The user sees
				// the normal confirmation — retries via the form are
				// cheap.
				slog.Error("mailer send", "err", err)
			}
		} else {
			slog.Info("magic-link request for unknown email", "email", addr)
		}

		if r.Header.Get("HX-Request") == "true" {
			if err := view.LoginCheckEmail(addr).Render(r.Context(), w); err != nil {
				slog.Error("render error", "err", err)
			}
			return
		}
		if err := view.LoginPageConfirmation(addr).Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}

// ConsumeMagicLink handles GET /login/magic/{token}. Atomically validates
// and marks the token used, then starts a session and redirects to /admin.
func ConsumeMagicLink(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.PathValue("token")

		userID, err := s.ConsumeLoginToken(r.Context(), token)
		if err != nil {
			if errors.Is(err, store.ErrTokenInvalid) {
				if err := view.LoginLinkInvalid().Render(r.Context(), w); err != nil {
					slog.Error("render error", "err", err)
				}
				return
			}
			slog.Error("consume login token", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if err := session.Create(r.Context(), w, s, userID); err != nil {
			slog.Error("session create error", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

// Logout handles POST /logout, destroys the session, and redirects.
func Logout(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := session.Destroy(r.Context(), w, r, s); err != nil {
			slog.Error("session destroy error", "err", err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func newRawToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func magicLinkText(appName, link string) string {
	return fmt.Sprintf(`Here's your sign-in link for %s:

%s

This link is good for 15 minutes. If you didn't ask to sign in, you can ignore this email.
`, appName, link)
}

func magicLinkHTML(appName, link string) string {
	return fmt.Sprintf(`<!doctype html>
<html><body style="font-family: -apple-system, Segoe UI, sans-serif; max-width: 480px; margin: 0 auto; padding: 24px; color: #18181b;">
<p>Here's your sign-in link for %s:</p>
<p><a href="%s" style="display: inline-block; padding: 10px 16px; background: #18181b; color: #ffffff; text-decoration: none; border-radius: 6px;">Sign in</a></p>
<p style="font-size: 13px; color: #71717a;">Or copy this URL: %s</p>
<p style="font-size: 13px; color: #71717a;">This link is good for 15 minutes. If you didn't ask to sign in, you can ignore this email.</p>
</body></html>
`, appName, link, link)
}
