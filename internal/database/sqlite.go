package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

type Store struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewSQLiteStore(dbFile string) (*Store, error) {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// 1. Connection String with Pro Options:
	// _pragma=journal_mode(WAL): Concurrent reads/writes
	// _pragma=foreign_keys(1): Enforce relational integrity
	// _pragma=busy_timeout(5000): Wait 5s before failing if DB is locked
	// cache=shared: Allow multiple connections to share cache
	dsn := fmt.Sprintf("%s?_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)&cache=shared", dbFile)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		logger.Error("Failed to open database", "error", err, "path", dbFile)
		return nil, err
	}

	// Performance Tuning
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping database", "error", err, "path", dbFile)
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	store := &Store{db: db, logger: logger}

	if err = store.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	logger.Info("Database initialized successfully", "path", dbFile, "mode", "WAL", "foreign_keys", "enabled")
	return store, nil
}

func (s *Store) migrate() error {
	driver, err := sqlite.WithInstance(s.db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}
	// "file://internal/database/migrations",
	migrationManager, err := migrate.NewWithDatabaseInstance("file://internal/database/migrations", "sqlite", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration manager: %w", err)
	}

	err = migrationManager.Up()
	if err != nil && err != migrate.ErrNoChange {
		s.logger.Error("Migration Failed", "error", err)
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	version, dirty, err := migrationManager.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	s.logger.Info("Migrations applied successfully",
		"version", version,
		"dirty", dirty)

	return nil
}

// shutdown the databse
func (s *Store) Close() error {
	s.logger.Info("Closing Database Connection")
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// GetDB returns the underlying database connection (for transactions, etc.)
func (s *Store) GetDB() *sql.DB {
	return s.db
}

// Health checks if database is alive
func (s *Store) Health() error {
	return s.db.Ping()
}
