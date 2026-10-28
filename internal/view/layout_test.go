package view

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b bytes.Buffer
	if err := c.Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

// The analytics init scripts must carry the configured IDs. templ does not
// interpolate single-brace expressions inside <script>, so a regression
// would ship the literal "{ GtagID }" and record nothing.
func TestAnalyticsIDsAreInterpolated(t *testing.T) {
	defer func(g, p string) { GtagID, PixelID = g, p }(GtagID, PixelID)
	GtagID, PixelID = "G-TEST123", "1234567890"

	html := render(t, HomePage())
	for _, want := range []string{`gtag('config', "G-TEST123")`, `fbq('init', "1234567890")`} {
		if !strings.Contains(html, want) {
			t.Errorf("home page missing %s", want)
		}
	}
	if strings.Contains(html, "{ GtagID }") || strings.Contains(html, "{ PixelID }") {
		t.Error("home page contains an uninterpolated analytics placeholder")
	}
}

// htmx and Turnstile load only where a form uses them, and only at page
// level: the form fragment htmx swaps back in must not reload either.
func TestFormScriptsLoadOnlyWhereNeeded(t *testing.T) {
	defer func(k string) { TurnstileSiteKey = k }(TurnstileSiteKey)
	TurnstileSiteKey = "test-site-key"

	const htmx, turnstile = "/static/js/htmx.min.js", "challenges.cloudflare.com/turnstile"
	cases := []struct {
		name          string
		html          string
		htmx, captcha bool
	}{
		{"home", render(t, HomePage()), false, false},
		{"contact page", render(t, ContactPage(nil, nil, false)), true, true},
		{"contact form fragment", render(t, ContactForm(nil, nil, false)), false, false},
		{"login page", render(t, LoginPage("", "")), true, false},
		{"login form fragment", render(t, LoginForm("", "")), false, false},
	}
	for _, c := range cases {
		if got := strings.Count(c.html, htmx); got != map[bool]int{true: 1, false: 0}[c.htmx] {
			t.Errorf("%s: htmx script tags = %d, want htmx=%v", c.name, got, c.htmx)
		}
		if got := strings.Count(c.html, turnstile); got != map[bool]int{true: 1, false: 0}[c.captcha] {
			t.Errorf("%s: Turnstile script tags = %d, want turnstile=%v", c.name, got, c.captcha)
		}
	}
}
