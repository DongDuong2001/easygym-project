// Command migrate applies the embedded SQL migrations: go run ./cmd/migrate [up|down|status].
package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"

	"easygym/db"
	"easygym/internal/platform"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "migrate:", err)
		os.Exit(1)
	}
}

func run() error {
	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	cfg, err := platform.LoadConfig()
	if err != nil {
		return err
	}
	conn, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx := context.Background()
	// The version table lives outside "public" so Supabase's REST API never exposes it.
	if _, err := conn.ExecContext(ctx, `create schema if not exists platform`); err != nil {
		return err
	}
	store, err := database.NewStore(database.DialectPostgres, "platform.goose_db_version")
	if err != nil {
		return err
	}
	migrations, _ := fs.Sub(db.Migrations, "migrations")
	p, err := goose.NewProvider("", conn, migrations, goose.WithStore(store))
	if err != nil {
		return err
	}

	switch cmd {
	case "up":
		res, err := p.Up(ctx)
		for _, r := range res {
			fmt.Println(r)
		}
		return err
	case "down":
		r, err := p.Down(ctx)
		if r != nil {
			fmt.Println(r)
		}
		return err
	case "status":
		st, err := p.Status(ctx)
		for _, s := range st {
			fmt.Printf("%-10s %s\n", s.State, s.Source.Path)
		}
		return err
	}
	return fmt.Errorf("unknown command %q (use up, down or status)", cmd)
}
