// Package migrations embeds the goose migration SQL files so the
// deployed binary doesn't need the migrations folder shipped alongside.
// main.go pulls FS and hands it to goose.SetBaseFS at startup.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
