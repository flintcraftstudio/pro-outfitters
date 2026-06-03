package handler

import (
	"log/slog"
	"net/http"

	"github.com/flintcraftstudio/standard-template/internal/view"
)

// NotFound handles requests to unknown routes and renders a 404 page.
func NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if err := view.NotFoundPage().Render(r.Context(), w); err != nil {
			slog.Error("render error", "err", err)
		}
	}
}
