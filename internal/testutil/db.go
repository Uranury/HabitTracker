package testutil

import (
	"database/sql"
	"testing"
	"time"

	"github.com/Uranury/HabitTracker/migrations"
	migratesqlite3 "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/mattn/go-sqlite3"
)

// NewTestDB opens an in-memory SQLite database and runs all migrations.
// The database is closed automatically when the test finishes.
func NewTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	sqlDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("testutil: open db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)

	if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatalf("testutil: enable foreign keys: %v", err)
	}

	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		t.Fatalf("testutil: create migration source: %v", err)
	}

	driver, err := migratesqlite3.WithInstance(sqlDB, &migratesqlite3.Config{})
	if err != nil {
		t.Fatalf("testutil: create migration driver: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "sqlite3", driver)
	if err != nil {
		t.Fatalf("testutil: create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("testutil: run migrations: %v", err)
	}

	db := sqlx.NewDb(sqlDB, "sqlite3")
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// SeedUser inserts a minimal user row and returns its UUID.
// Use this to satisfy the foreign key constraint when testing habits or check-ins.
func SeedUser(t *testing.T, db *sqlx.DB) uuid.UUID {
	t.Helper()
	id := uuid.New()
	now := time.Now().Format(time.RFC3339)
	_, err := db.Exec(
		`INSERT INTO users (id, username, password, time_zone, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		id.String(), "testuser_"+id.String()[:8], "hashed", "UTC", now, now,
	)
	if err != nil {
		t.Fatalf("testutil: seed user: %v", err)
	}
	return id
}
