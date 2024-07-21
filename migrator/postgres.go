package migrator

import (
	"context"
	"errors"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"avito-test-task/internal/config"
)

const (
	maxRetries = 30
	delay      = time.Second
)

// migratePostgres applies migrations to the postgres database.
// It first waits postgres for 30 seconds in case if it's not ready yet and then
// applies migrations stored in a file migrationsPath.
func migratePostgres(ctx context.Context, cfg *config.Config, migrationsPath, migrationsTable string) *migrate.Migrate {
	err := waitForPostgres(cfg.DB.ConnectionString(), maxRetries, delay)
	if err != nil {
		panic(err)
	}

	ensurePgsDBExists(ctx, &cfg.DB)
	m, err := migrate.New(
		"file://"+migrationsPath,
		cfg.DB.ConnectionString()+"&x-migrations-table="+migrationsTable,
	)
	if err != nil {
		panic(err)
	}

	return m
}

// waitForPostgres waits maxRetries * delay until postgres is ready to accept connections
// or returns error in case if postgres wasn't up til that point.
func waitForPostgres(url string, maxRetries int, delay time.Duration) error {
	for range maxRetries {
		conn, err := pgx.Connect(context.Background(), url)
		if err == nil {
			_ = conn.Close(context.Background())
			return nil
		}
		time.Sleep(delay)
	}
	return errors.New("postgres didn't become available within the specified time")
}

// ensurePgsDBExists connects to postgres 'postgres' database,
// creates a new database that's specified in the db config.DB and
// connects to the created database.
func ensurePgsDBExists(ctx context.Context, db *config.DB) {
	dbName := db.DBName
	db.DBName = "postgres"
	connStr := db.ConnectionString()
	db.DBName = dbName
	dbPool, e := pgxpool.New(ctx, connStr)
	if e != nil {
		panic(e)
	}
	defer dbPool.Close()

	_, _ = dbPool.Exec(ctx, "CREATE DATABASE \""+dbName+"\";")
}
