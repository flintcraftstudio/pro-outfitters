package handler

import (
	"log/slog"
	"net/http"

	"github.com/flintcraftstudio/standard-template/internal/view"
)

// Lodges handles GET /lodges — the lodges overview page.
func Lodges() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := view.LodgesPage().Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}

// Lodge handles GET /lodges/{slug} — an individual lodge detail page. An
// unknown slug renders the 404 page so explicit-but-missing lodges don't
// 200 with an empty shell.
func Lodge() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		l, ok := view.LodgeBySlug(slug)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			if err := view.NotFoundPage().Render(r.Context(), w); err != nil {
				slog.Error("render error", "err", err)
			}
			return
		}
		if err := view.LodgePage(l).Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}

// SmithRiver handles GET /smith-river — the featured trip page.
func SmithRiver() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := view.SmithRiverPage().Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}

// About handles GET /about — the company story + team page.
func About() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := view.AboutPage().Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}
