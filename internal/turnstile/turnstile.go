// Package turnstile verifies Cloudflare Turnstile challenge tokens.
// A Verifier is safe for concurrent use; one per process is enough.
//
// When the secret key is empty the Verifier is in "disabled" mode:
// Verify returns nil without hitting the network. That's the dev
// fallback — the contact form skips the widget when no site key is
// configured, so the server-side check would be spurious anyway.
package turnstile

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// ErrFailed is returned when Cloudflare says the token is bad. Missing
// tokens, expired tokens, and "widget-level" failures all collapse to
// this single error.
var ErrFailed = errors.New("turnstile verification failed")

// Verifier holds the secret key + HTTP client. Secret="" disables the
// verifier (Verify is a no-op) so dev runs without Cloudflare creds.
type Verifier struct {
	Secret string
	HTTP   *http.Client
}

func New(secret string) *Verifier {
	return &Verifier{Secret: secret, HTTP: &http.Client{Timeout: 10 * time.Second}}
}

// Enabled reports whether the verifier will actually check tokens.
func (v *Verifier) Enabled() bool { return v != nil && v.Secret != "" }

type siteverifyResponse struct {
	Success     bool     `json:"success"`
	ErrorCodes  []string `json:"error-codes"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	Action      string   `json:"action"`
}

// Verify exchanges a token (from the widget) for a yes/no. remoteIP is
// optional; when supplied Cloudflare uses it as an additional signal.
func (v *Verifier) Verify(ctx context.Context, token, remoteIP string) error {
	if !v.Enabled() {
		return nil
	}
	if token == "" {
		return ErrFailed
	}

	form := url.Values{}
	form.Set("secret", v.Secret)
	form.Set("response", token)
	if remoteIP != "" {
		form.Set("remoteip", remoteIP)
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://challenges.cloudflare.com/turnstile/v0/siteverify",
		bytes.NewBufferString(form.Encode()),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := v.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("turnstile: %w", err)
	}
	defer resp.Body.Close()

	var out siteverifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return fmt.Errorf("turnstile decode: %w", err)
	}
	if !out.Success {
		return ErrFailed
	}
	return nil
}
