package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func InitDB(ctx context.Context, driverName, dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return db, nil
}

func RunMigrations(migrationsPath, driverName, dsn string) error {
	migrationsURL := migrationsPath
	if !strings.HasPrefix(migrationsURL, "file://") {
		migrationsURL = "file://" + migrationsPath
	}

	migrationsDB, err := sql.Open(driverName, dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	defer func() {
		_ = migrationsDB.Close()
	}()

	driver, err := sqlite3.WithInstance(migrationsDB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsURL, driverName, driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
