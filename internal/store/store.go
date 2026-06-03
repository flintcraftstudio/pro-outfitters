// Package store wraps the SQLite database and exposes typed query
// methods. One Store value per process; methods are safe for concurrent
// use because *sql.DB is.
package store

import (
	"context"
	"database/sql"
	"time"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

// SqliteDatetime matches the format SQLite's CURRENT_TIMESTAMP emits, so
// values stored this way are directly comparable in WHERE clauses
// without datetime() coercion. modernc/sqlite serializes time.Time via
// fmt.Stringer (which leaves monotonic-clock junk on the value), so we
// format explicitly on every write and parse explicitly on every read.
const SqliteDatetime = "2006-01-02 15:04:05"

func (s *Store) CreateSession(ctx context.Context, token string, userID int64, expiresAt time.Time) error {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO admin_sessions (id, admin_user_id, expires_at) VALUES (?, ?, ?)",
		token, userID, expiresAt.UTC().Format(SqliteDatetime),
	)
	return err
}

func (s *Store) GetSession(ctx context.Context, token string) (int64, time.Time, error) {
	var userID int64
	var expiresAtStr string
	err := s.db.QueryRowContext(ctx,
		"SELECT admin_user_id, expires_at FROM admin_sessions WHERE id = ? AND expires_at > CURRENT_TIMESTAMP",
		token,
	).Scan(&userID, &expiresAtStr)
	if err != nil {
		return 0, time.Time{}, err
	}
	expiresAt, err := time.Parse(SqliteDatetime, expiresAtStr)
	if err != nil {
		return 0, time.Time{}, err
	}
	return userID, expiresAt, nil
}

func (s *Store) DeleteSession(ctx context.Context, token string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM admin_sessions WHERE id = ?", token)
	return err
}

func (s *Store) GetUserByID(ctx context.Context, id int64) (int64, string, string, error) {
	var userID int64
	var email, displayName string
	err := s.db.QueryRowContext(ctx,
		"SELECT id, email, display_name FROM admin_users WHERE id = ?",
		id,
	).Scan(&userID, &email, &displayName)
	return userID, email, displayName, err
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (int64, string, error) {
	var id int64
	var userEmail string
	err := s.db.QueryRowContext(ctx,
		"SELECT id, email FROM admin_users WHERE email = ?",
		email,
	).Scan(&id, &userEmail)
	return id, userEmail, err
}

// CreateUser inserts an admin user. Auth is magic-link only — no
// password is collected.
func (s *Store) CreateUser(ctx context.Context, email, displayName string) (int64, error) {
	result, err := s.db.ExecContext(ctx,
		"INSERT INTO admin_users (email, display_name) VALUES (?, ?)",
		email, displayName,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *Store) DeleteExpiredSessions(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM admin_sessions WHERE expires_at <= CURRENT_TIMESTAMP")
	return err
}

// DashboardStats is the data the admin dashboard reads on every load.
type DashboardStats struct {
	InquiriesThisMonth int
	UnrespondedCount   int
	LastInquiryAt      sql.NullTime
}

func (s *Store) GetDashboardStats(ctx context.Context) (DashboardStats, error) {
	var stats DashboardStats

	if err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM inquiries WHERE created_at >= datetime('now', 'start of month')",
	).Scan(&stats.InquiriesThisMonth); err != nil {
		return stats, err
	}

	if err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM inquiries WHERE responded_at IS NULL AND archived_at IS NULL",
	).Scan(&stats.UnrespondedCount); err != nil {
		return stats, err
	}

	// modernc/sqlite returns SQLite's CURRENT_TIMESTAMP values as the
	// raw "YYYY-MM-DD HH:MM:SS" string. sql.NullTime's Scan only
	// accepts RFC3339-shaped input and will error on this, so scan
	// into NullString and parse with SqliteDatetime explicitly.
	var raw sql.NullString
	if err := s.db.QueryRowContext(ctx,
		"SELECT MAX(created_at) FROM inquiries",
	).Scan(&raw); err != nil {
		return stats, err
	}
	if raw.Valid {
		if t, err := time.Parse(SqliteDatetime, raw.String); err == nil {
			stats.LastInquiryAt = sql.NullTime{Time: t, Valid: true}
		}
	}

	return stats, nil
}
