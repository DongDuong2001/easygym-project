// Package db holds the SQL migrations and the first-run seed, embedded into the binaries.
package db

import "embed"

//go:embed migrations/*.sql
var Migrations embed.FS

//go:embed seed/seed.json
var Seed []byte
