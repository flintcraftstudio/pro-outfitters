package store

import (
	"context"
	"strings"
	"time"
)

// SiteSettings is the editable public-site configuration — one row,
// id=1. String fields default to "" in the schema so the public
// renderer can treat empty as "not set, hide this element" without
// nil-checking. Per-project columns can be added via additional
// migrations; extend this struct to surface them.
type SiteSettings struct {
	ContactEmail string
	ContactPhone string
	Tagline      string
	UpdatedAt    time.Time
}

type SiteSettingsInput struct {
	ContactEmail string
	ContactPhone string
	Tagline      string
}

func (s *Store) GetSettings(ctx context.Context) (SiteSettings, error) {
	var out SiteSettings
	var updatedStr string
	err := s.db.QueryRowContext(ctx, `
		SELECT contact_email, contact_phone, tagline, updated_at
		FROM site_settings WHERE id = 1
	`).Scan(&out.ContactEmail, &out.ContactPhone, &out.Tagline, &updatedStr)
	if err != nil {
		return out, err
	}
	out.UpdatedAt, _ = time.Parse(SqliteDatetime, updatedStr)
	return out, nil
}

func (s *Store) UpdateSettings(ctx context.Context, in SiteSettingsInput) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE site_settings
		SET contact_email = ?, contact_phone = ?, tagline = ?, updated_at = ?
		WHERE id = 1
	`,
		strings.TrimSpace(in.ContactEmail),
		strings.TrimSpace(in.ContactPhone),
		strings.TrimSpace(in.Tagline),
		time.Now().UTC().Format(SqliteDatetime),
	)
	return err
}
