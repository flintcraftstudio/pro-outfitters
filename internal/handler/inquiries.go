package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/flintcraftstudio/standard-template/internal/session"
	"github.com/flintcraftstudio/standard-template/internal/store"
	"github.com/flintcraftstudio/standard-template/internal/view"
)

func Inquiries(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := session.FromContext(r.Context())
		filter := parseInquiryFilter(r.URL.Query().Get("filter"))
		items, err := s.ListInquiries(r.Context(), filter)
		if err != nil {
			slog.Error("list inquiries", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		if err := view.AdminInquiries(user.DisplayName, filter, items).Render(r.Context(), w); err != nil {
			slog.Error("render", "err", err)
		}
	}
}

func InquiryDetail(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := session.FromContext(r.Context())
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		inq, err := s.GetInquiry(r.Context(), id)
		if errors.Is(err, store.ErrInquiryNotFound) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			slog.Error("get inquiry", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		if err := view.AdminInquiryDetail(user.DisplayName, *inq).Render(r.Context(), w); err != nil {
			slog.Error("render", "err", err)
		}
	}
}

// ToggleInquiry flips one of the two status columns (responded_at or
// archived_at) on or off.
func ToggleInquiry(s *store.Store, field string, set bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		switch field {
		case "responded":
			err = s.MarkInquiryResponded(r.Context(), id, set)
		case "archived":
			err = s.MarkInquiryArchived(r.Context(), id, set)
		default:
			http.Error(w, "bad field", http.StatusBadRequest)
			return
		}
		if err != nil {
			if errors.Is(err, store.ErrInquiryNotFound) {
				http.NotFound(w, r)
				return
			}
			slog.Error("toggle inquiry", "field", field, "set", set, "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// Archiving drops the user back to the inbox list; "responded"
		// just refreshes the detail view.
		if field == "archived" && set {
			http.Redirect(w, r, "/admin/inquiries", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/admin/inquiries/"+strconv.FormatInt(id, 10), http.StatusSeeOther)
	}
}

func parseInquiryFilter(s string) store.InquiryFilter {
	switch store.InquiryFilter(s) {
	case store.FilterUnresponded, store.FilterResponded, store.FilterArchived, store.FilterAll:
		return store.InquiryFilter(s)
	default:
		return store.FilterAll
	}
}
