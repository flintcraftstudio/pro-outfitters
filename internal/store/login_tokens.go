package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"
)

// ErrTokenInvalid covers all the ways a magic-link token can fail
// validation — unknown hash, expired, or already consumed. Callers
// surface a single generic "invalid or expired" message to the user
// regardless of which case applied, so we don't narrow the reason.
var ErrTokenInvalid = errors.New("login token invalid or expired")

// CreateLoginToken persists a single-use magic-link token hashed at
// rest. The caller keeps the raw token to embed in the email URL;
// the database only sees SHA-256(raw), so a compromised DB dump
// can't be replayed to mint valid sessions.
func (s *Store) CreateLoginToken(ctx context.Context, userID int64, rawToken string, expiresAt time.Time) error {
	hash := hashLoginToken(rawToken)
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO login_tokens (admin_user_id, token_hash, expires_at) VALUES (?, ?, ?)",
		userID, hash, expiresAt.UTC().Format(SqliteDatetime),
	)
	return err
}

// ConsumeLoginToken validates the raw token and atomically marks it
// used, returning the owning admin_user_id. Returns ErrTokenInvalid
// when the hash is unknown, expired, or already consumed.
func (s *Store) ConsumeLoginToken(ctx context.Context, rawToken string) (int64, error) {
	hash := hashLoginToken(rawToken)

	var userID int64
	err := s.db.QueryRowContext(ctx, `
		UPDATE login_tokens
		   SET used_at = CURRENT_TIMESTAMP
		 WHERE token_hash = ?
		   AND used_at IS NULL
		   AND expires_at > CURRENT_TIMESTAMP
		RETURNING admin_user_id
	`, hash).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrTokenInvalid
	}
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// DeleteExpiredLoginTokens purges expired and already-consumed rows.
// Safe to call from the cleanup worker; no user-facing effect.
func (s *Store) DeleteExpiredLoginTokens(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx,
		"DELETE FROM login_tokens WHERE expires_at <= CURRENT_TIMESTAMP OR used_at IS NOT NULL",
	)
	return err
}

func hashLoginToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
