package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/flintcraftstudio/standard-template/internal/session"
	"github.com/flintcraftstudio/standard-template/internal/store"
	"github.com/flintcraftstudio/standard-template/internal/view"
)

const taglineMaxLen = 200

func Settings(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := session.FromContext(r.Context())
		settings, err := s.GetSettings(r.Context())
		if err != nil {
			slog.Error("get settings", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		data := view.SettingsFormData{
			Saved:        r.URL.Query().Get("saved") == "1",
			ContactEmail: settings.ContactEmail,
			ContactPhone: settings.ContactPhone,
			Tagline:      settings.Tagline,
		}
		if err := view.AdminSettings(user.DisplayName, data).Render(r.Context(), w); err != nil {
			slog.Error("render", "err", err)
		}
	}
}

func UpdateSettings(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := session.FromContext(r.Context())

		input, formData, errs := parseSettingsForm(r)
		formData.Errors = errs

		if len(errs) > 0 {
			if err := view.AdminSettings(user.DisplayName, formData).Render(r.Context(), w); err != nil {
				slog.Error("render", "err", err)
			}
			return
		}

		if err := s.UpdateSettings(r.Context(), input); err != nil {
			slog.Error("update settings", "err", err)
			formData.Errors = map[string]string{"form": "Could not save. Please try again."}
			if err := view.AdminSettings(user.DisplayName, formData).Render(r.Context(), w); err != nil {
				slog.Error("render", "err", err)
			}
			return
		}

		http.Redirect(w, r, "/admin/settings?saved=1", http.StatusSeeOther)
	}
}

func parseSettingsForm(r *http.Request) (store.SiteSettingsInput, view.SettingsFormData, map[string]string) {
	errs := map[string]string{}
	if err := r.ParseForm(); err != nil {
		errs["form"] = "Bad form data."
		return store.SiteSettingsInput{}, view.SettingsFormData{}, errs
	}

	emailAddr := strings.TrimSpace(r.FormValue("contact_email"))
	if emailAddr != "" && !strings.Contains(emailAddr, "@") {
		errs["contact_email"] = "Enter a valid email address."
	}

	phone := strings.TrimSpace(r.FormValue("contact_phone"))
	tagline := strings.TrimSpace(r.FormValue("tagline"))
	if len(tagline) > taglineMaxLen {
		errs["tagline"] = "Keep the tagline under 200 characters."
	}

	input := store.SiteSettingsInput{
		ContactEmail: emailAddr,
		ContactPhone: phone,
		Tagline:      tagline,
	}
	formData := view.SettingsFormData{
		ContactEmail: emailAddr,
		ContactPhone: phone,
		Tagline:      tagline,
	}
	return input, formData, errs
}
