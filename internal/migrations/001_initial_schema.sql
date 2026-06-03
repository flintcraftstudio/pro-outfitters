-- +goose Up

CREATE TABLE admin_users (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    email         TEXT    NOT NULL UNIQUE,
    display_name  TEXT    NOT NULL,
    last_login_at TEXT,
    created_at    TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE admin_sessions (
    id            TEXT    PRIMARY KEY,
    admin_user_id INTEGER NOT NULL REFERENCES admin_users(id) ON DELETE CASCADE,
    created_at    TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at    TEXT    NOT NULL,
    user_agent    TEXT,
    ip_address    TEXT
);

CREATE INDEX idx_admin_sessions_user    ON admin_sessions(admin_user_id);
CREATE INDEX idx_admin_sessions_expires ON admin_sessions(expires_at);

CREATE TABLE login_tokens (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    admin_user_id INTEGER NOT NULL REFERENCES admin_users(id) ON DELETE CASCADE,
    token_hash    TEXT    NOT NULL UNIQUE,
    expires_at    TEXT    NOT NULL,
    used_at       TEXT,
    created_at    TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_login_tokens_hash    ON login_tokens(token_hash);
CREATE INDEX idx_login_tokens_expires ON login_tokens(expires_at);

CREATE TABLE site_settings (
    id            INTEGER PRIMARY KEY CHECK (id = 1),
    contact_email TEXT    NOT NULL DEFAULT '',
    contact_phone TEXT    NOT NULL DEFAULT '',
    tagline       TEXT    NOT NULL DEFAULT '',
    updated_at    TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO site_settings (id) VALUES (1);

CREATE TABLE inquiries (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    name         TEXT    NOT NULL,
    email        TEXT    NOT NULL,
    phone        TEXT,
    subject      TEXT    NOT NULL DEFAULT '',
    message      TEXT    NOT NULL,
    responded_at TEXT,
    archived_at  TEXT,
    created_at   TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_inquiries_created     ON inquiries(created_at DESC);
CREATE INDEX idx_inquiries_unresponded ON inquiries(responded_at) WHERE responded_at IS NULL;

-- +goose Down

DROP TABLE inquiries;
DROP TABLE site_settings;
DROP INDEX IF EXISTS idx_login_tokens_hash;
DROP INDEX IF EXISTS idx_login_tokens_expires;
DROP TABLE login_tokens;
DROP TABLE admin_sessions;
DROP TABLE admin_users;
