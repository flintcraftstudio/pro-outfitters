// seed creates an admin user. Used during initial setup to bootstrap
// the magic-link login flow. Usage:
//
//	go run ./cmd/seed admin@example.com
//	go run ./cmd/seed admin@example.com "Display Name"
//
// Reads DB_PATH from the environment (defaults to ./data/app.db).
// Migrations must already have run — point this at a database the
// server has opened at least once.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/flintcraftstudio/standard-template/internal/store"

	_ "modernc.org/sqlite"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: seed <email> [display_name]\n")
		os.Exit(1)
	}

	email := os.Args[1]
	displayName := defaultDisplayName(email)
	if len(os.Args) >= 3 {
		displayName = os.Args[2]
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/app.db"
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	st := store.New(db)
	id, err := st.CreateUser(context.Background(), email, displayName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating user: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created admin user %s (id=%d, display_name=%q)\n", email, id, displayName)
}

func defaultDisplayName(email string) string {
	if i := strings.IndexByte(email, '@'); i > 0 {
		return email[:i]
	}
	return email
}
