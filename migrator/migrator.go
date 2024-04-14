package migrator

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/golang-migrate/migrate/v4"

	"avito-test-task/internal/config"
)

const (
	up   = "up"
	down = "down"
)

func Migrate(
	ctx context.Context,
	out io.Writer,
	args []string,
	getenv func(string) (string, bool),
	db string,
) {
	flagSet := flag.NewFlagSet("migrate", flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)
	migrationsPath := flagSet.String("migrations-path", "", "path to migrations")
	migrationsTable := flagSet.String("migrations-table", "migrations", "name of migrations table")
	direction := flagSet.String("direction", up, "migration direction (up/down)")

	_ = flagSet.Parse(args)
	if migrationsPath == nil || *migrationsPath == "" {
		panic("migrations-path is required")
	}
	if direction == nil || *direction != up && *direction != down {
		panic("direction must be either 'up' or 'down'")
	}

	m := new(migrate.Migrate)

	switch db {
	case "postgres":
		cfg := config.MustLoad(args, getenv)
		m = migratePostgres(ctx, cfg, *migrationsPath, *migrationsTable)
	default:
		panic("unsupported db")
	}

	switch *direction {
	case up:
		migrateUp(m, out)
	case down:
		migrateDown(m, out)
	}

	_, _ = fmt.Fprintln(out, "migrations applied")
}

func migrateUp(m *migrate.Migrate, out io.Writer) {
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			_, _ = fmt.Fprintln(out, "no migrations to apply")
		} else {
			panic(err)
		}
	}
}

func migrateDown(m *migrate.Migrate, out io.Writer) {
	if err := m.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			_, _ = fmt.Fprintln(out, "no migrations to rollback")
		} else {
			panic(err)
		}
	}
}
