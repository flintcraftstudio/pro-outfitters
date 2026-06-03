package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/flintcraftstudio/standard-template/internal/session"
	"github.com/flintcraftstudio/standard-template/internal/store"
	"github.com/flintcraftstudio/standard-template/internal/view"
)

// Dashboard handles GET /admin and renders the admin home.
func Dashboard(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := session.FromContext(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		stats, err := s.GetDashboardStats(r.Context())
		if err != nil {
			slog.Error("dashboard stats error", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		var lastInquiry *time.Time
		if stats.LastInquiryAt.Valid {
			t := stats.LastInquiryAt.Time
			lastInquiry = &t
		}

		if err := view.AdminDashboard(user.DisplayName, stats.InquiriesThisMonth, stats.UnrespondedCount, lastInquiry).Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}
