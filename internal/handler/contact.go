package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/flintcraftstudio/standard-template/internal/email"
	"github.com/flintcraftstudio/standard-template/internal/store"
	"github.com/flintcraftstudio/standard-template/internal/turnstile"
	"github.com/flintcraftstudio/standard-template/internal/view"
)

// ContactDeps bundles everything the public contact form handler reaches for.
type ContactDeps struct {
	Store     *store.Store
	Mailer    email.Sender
	Turnstile *turnstile.Verifier
	AppName   string
}

// Contact handles GET /contact and renders the contact form.
func Contact() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := view.ContactPage(nil, nil, false).Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}

// ContactSubmit handles POST /contact: validates input, persists an
// inquiry to the database, and fires off a notification email to the
// site contact (best-effort, in a background goroutine).
func ContactSubmit(deps ContactDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		values := map[string]string{
			"name":    strings.TrimSpace(r.FormValue("name")),
			"email":   strings.TrimSpace(r.FormValue("email")),
			"phone":   strings.TrimSpace(r.FormValue("phone")),
			"subject": strings.TrimSpace(r.FormValue("subject")),
			"message": strings.TrimSpace(r.FormValue("message")),
		}

		errs := validateContact(values)
		if len(errs) > 0 {
			if err := view.ContactForm(errs, values, false).Render(r.Context(), w); err != nil {
				slog.Error("render error", "err", err)
			}
			return
		}

		if err := deps.Turnstile.Verify(r.Context(), r.FormValue("cf-turnstile-response"), r.RemoteAddr); err != nil {
			errs := map[string]string{"form": "Verification failed. Please try again."}
			if err := view.ContactForm(errs, values, false).Render(r.Context(), w); err != nil {
				slog.Error("render error", "err", err)
			}
			return
		}

		input := store.InquiryInput{
			Name:    values["name"],
			Email:   values["email"],
			Phone:   values["phone"],
			Subject: values["subject"],
			Message: values["message"],
		}
		id, err := deps.Store.CreateInquiry(r.Context(), input)
		if err != nil {
			slog.Error("create inquiry", "err", err)
			errs := map[string]string{"form": "Could not save your message. Please try again."}
			if err := view.ContactForm(errs, values, false).Render(r.Context(), w); err != nil {
				slog.Error("render error", "err", err)
			}
			return
		}

		// Fire-and-forget notification email so a slow Postmark blip
		// doesn't block the visitor's success render.
		settings, _ := deps.Store.GetSettings(r.Context())
		go sendInquiryNotification(deps.Mailer, deps.AppName, settings.ContactEmail, id, input)

		if err := view.ContactForm(nil, nil, true).Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}

func validateContact(values map[string]string) map[string]string {
	errs := make(map[string]string)
	if values["name"] == "" {
		errs["name"] = "Name is required."
	}
	if values["email"] == "" {
		errs["email"] = "Email is required."
	} else if !strings.Contains(values["email"], "@") {
		errs["email"] = "Enter a valid email address."
	}
	if values["message"] == "" {
		errs["message"] = "Message is required."
	}
	return errs
}

// sendInquiryNotification emails the configured site contact when a new
// inquiry lands. Runs with a fresh background context so a client
// disconnect mid-render doesn't cancel the send.
func sendInquiryNotification(mailer email.Sender, appName, notifyEmail string, id int64, in store.InquiryInput) {
	if notifyEmail == "" {
		slog.Warn("contact notification skipped: site_settings.contact_email is empty",
			"inquiry_id", id,
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	subject := fmt.Sprintf("New inquiry from %s", in.Name)
	if in.Subject != "" {
		subject = fmt.Sprintf("New inquiry: %s", in.Subject)
	}
	body := fmt.Sprintf(
		"New inquiry #%d on %s.\n\n"+
			"From:    %s <%s>\n"+
			"Phone:   %s\n"+
			"Subject: %s\n\n"+
			"-- Message --\n\n%s\n",
		id, appName, in.Name, in.Email, defaultStr(in.Phone, "—"), defaultStr(in.Subject, "—"), in.Message,
	)

	msg := email.Message{
		To:      notifyEmail,
		Subject: subject,
		Text:    body,
		ReplyTo: in.Email,
	}
	if err := mailer.Send(ctx, msg); err != nil {
		slog.Error("inquiry notification send", "err", err, "inquiry_id", id)
	}
}

func defaultStr(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}
