package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// RespondedAt and ArchivedAt are sql.NullString rather than sql.NullTime
// because modernc/sqlite stores timestamps as SQLite's native
// "YYYY-MM-DD HH:MM:SS" text format (the output of CURRENT_TIMESTAMP),
// which sql.NullTime.Scan can't parse — NullTime only accepts RFC3339.
// The rest of the app only checks .Valid on these fields (to show
// status badges and filter tabs); the exact moment is never read.
type Inquiry struct {
	ID          int64
	Name        string
	Email       string
	Phone       sql.NullString
	Subject     string
	Message     string
	RespondedAt sql.NullString
	ArchivedAt  sql.NullString
	CreatedAt   time.Time
}

type InquiryInput struct {
	Name    string
	Email   string
	Phone   string
	Subject string
	Message string
}

// InquiryFilter picks which subset of inquiries the inbox returns.
//
//	FilterAll          — everything except archived
//	FilterUnresponded  — responded_at IS NULL AND archived_at IS NULL
//	FilterResponded    — responded_at IS NOT NULL AND archived_at IS NULL
//	FilterArchived     — archived_at IS NOT NULL
type InquiryFilter string

const (
	FilterAll         InquiryFilter = "all"
	FilterUnresponded InquiryFilter = "unresponded"
	FilterResponded   InquiryFilter = "responded"
	FilterArchived    InquiryFilter = "archived"
)

var ErrInquiryNotFound = errors.New("inquiry not found")

func (s *Store) CreateInquiry(ctx context.Context, in InquiryInput) (int64, error) {
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO inquiries (name, email, phone, subject, message)
		VALUES (?, ?, ?, ?, ?)
	`,
		strings.TrimSpace(in.Name),
		strings.TrimSpace(in.Email),
		nullableText(in.Phone),
		strings.TrimSpace(in.Subject),
		in.Message,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// InquiryListItem is the inbox-row projection — stripped of message
// body (the list only needs the first line; full message is on the
// detail view).
type InquiryListItem struct {
	ID        int64
	Name      string
	Email     string
	Subject   string
	Preview   string // first line of message, trimmed
	CreatedAt time.Time
	Responded bool
	Archived  bool
}

func (s *Store) ListInquiries(ctx context.Context, filter InquiryFilter) ([]InquiryListItem, error) {
	where := "WHERE archived_at IS NULL"
	switch filter {
	case FilterUnresponded:
		where = "WHERE archived_at IS NULL AND responded_at IS NULL"
	case FilterResponded:
		where = "WHERE archived_at IS NULL AND responded_at IS NOT NULL"
	case FilterArchived:
		where = "WHERE archived_at IS NOT NULL"
	case FilterAll, "":
		// default: unarchived only
	}

	q := fmt.Sprintf(`
		SELECT id, name, email, subject, message, responded_at, archived_at, created_at
		FROM inquiries
		%s
		ORDER BY created_at DESC
		LIMIT 100
	`, where)

	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []InquiryListItem
	for rows.Next() {
		var it InquiryListItem
		var message, createdStr string
		var responded, archived sql.NullString
		if err := rows.Scan(
			&it.ID, &it.Name, &it.Email, &it.Subject, &message,
			&responded, &archived, &createdStr,
		); err != nil {
			return nil, err
		}
		it.Preview = firstLine(message, 140)
		it.CreatedAt, _ = time.Parse(SqliteDatetime, createdStr)
		it.Responded = responded.Valid
		it.Archived = archived.Valid
		out = append(out, it)
	}
	return out, rows.Err()
}

func (s *Store) GetInquiry(ctx context.Context, id int64) (*Inquiry, error) {
	var it Inquiry
	var createdStr string
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, email, phone, subject, message,
		       responded_at, archived_at, created_at
		FROM inquiries WHERE id = ?
	`, id).Scan(
		&it.ID, &it.Name, &it.Email, &it.Phone, &it.Subject, &it.Message,
		&it.RespondedAt, &it.ArchivedAt, &createdStr,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInquiryNotFound
	}
	if err != nil {
		return nil, err
	}
	it.CreatedAt, _ = time.Parse(SqliteDatetime, createdStr)
	return &it, nil
}

func (s *Store) MarkInquiryResponded(ctx context.Context, id int64, responded bool) error {
	var q string
	if responded {
		q = "UPDATE inquiries SET responded_at = CURRENT_TIMESTAMP WHERE id = ?"
	} else {
		q = "UPDATE inquiries SET responded_at = NULL WHERE id = ?"
	}
	res, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrInquiryNotFound
	}
	return nil
}

func (s *Store) MarkInquiryArchived(ctx context.Context, id int64, archived bool) error {
	var q string
	if archived {
		q = "UPDATE inquiries SET archived_at = CURRENT_TIMESTAMP WHERE id = ?"
	} else {
		q = "UPDATE inquiries SET archived_at = NULL WHERE id = ?"
	}
	res, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrInquiryNotFound
	}
	return nil
}

// CountUnrespondedInquiries powers the dashboard quick-action badge.
// Excludes archived so the badge doesn't light up for old things the
// user already moved on from.
func (s *Store) CountUnrespondedInquiries(ctx context.Context) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM inquiries
		WHERE responded_at IS NULL AND archived_at IS NULL
	`).Scan(&n)
	return n, err
}

func firstLine(s string, max int) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	s = strings.TrimSpace(s)
	if len(s) > max {
		s = s[:max] + "…"
	}
	return s
}

func nullableText(s string) any {
	if s = strings.TrimSpace(s); s == "" {
		return nil
	}
	return s
}
