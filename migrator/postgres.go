package migrator

import (
	"context"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	"avito-test-task/internal/config"
)

func migratePostgres(ctx context.Context, cfg *config.Config, migrationsPath, migrationsTable string) *migrate.Migrate {
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
