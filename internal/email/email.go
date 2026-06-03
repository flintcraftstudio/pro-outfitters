// Package email sends transactional messages. Two implementations ship:
//
//   - PostmarkSender — production, talks to Postmark's REST API.
//   - LogSender      — dev fallback, writes the message to slog instead of sending.
//
// Selection happens in main.go: when POSTMARK_SERVER_TOKEN is set the
// app picks PostmarkSender, otherwise LogSender. Handlers depend only
// on the Sender interface so swapping providers later is a single-file
// change.
package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// Sender is the minimal shape every implementation honors.
type Sender interface {
	Send(ctx context.Context, msg Message) error
}

// Message is the transactional envelope. Text is required; HTML is
// optional and clients that accept it will prefer it.
type Message struct {
	To      string
	Subject string
	Text    string
	HTML    string
	ReplyTo string
}

// LogSender prints the outgoing message to the app logger. Used in
// local dev so magic-link URLs surface in the server log without an
// API key or DNS setup.
type LogSender struct{}

func (LogSender) Send(_ context.Context, m Message) error {
	slog.Info("email (dev log sender)",
		"to", m.To,
		"subject", m.Subject,
		"text", m.Text,
	)
	return nil
}

// PostmarkSender is the production path. Uses the REST API directly
// instead of an SDK to avoid adding a dependency for what's effectively
// a single POST.
// https://postmarkapp.com/developer/api/email-api
type PostmarkSender struct {
	ServerToken string
	From        string
	HTTP        *http.Client
}

// NewPostmarkSender returns a ready client. From accepts either a bare
// address ("admin@example.com") or a display-name form
// ("Name <admin@example.com>"). The address must correspond to a
// Postmark Sender Signature or a verified domain on the server.
func NewPostmarkSender(serverToken, from string) *PostmarkSender {
	return &PostmarkSender{
		ServerToken: serverToken,
		From:        from,
		HTTP:        &http.Client{Timeout: 15 * time.Second},
	}
}

type postmarkRequest struct {
	From          string `json:"From"`
	To            string `json:"To"`
	Subject       string `json:"Subject"`
	TextBody      string `json:"TextBody,omitempty"`
	HtmlBody      string `json:"HtmlBody,omitempty"`
	ReplyTo       string `json:"ReplyTo,omitempty"`
	MessageStream string `json:"MessageStream"`
}

func (s *PostmarkSender) Send(ctx context.Context, m Message) error {
	body, err := json.Marshal(postmarkRequest{
		From:          s.From,
		To:            m.To,
		Subject:       m.Subject,
		TextBody:      m.Text,
		HtmlBody:      m.HTML,
		ReplyTo:       m.ReplyTo,
		MessageStream: "outbound",
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.postmarkapp.com/email", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("X-Postmark-Server-Token", s.ServerToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("postmark: %s: %s", resp.Status, string(raw))
	}
	return nil
}
