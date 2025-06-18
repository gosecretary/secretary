package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(driver, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	// For now, we'll create tables directly
	// In a production environment, you'd want to use a proper migration tool
	createTablesSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user',
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS resources (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		type TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS credentials (
		id TEXT PRIMARY KEY,
		resource_id TEXT NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS permissions (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		resource_id TEXT NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
		action TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		UNIQUE(user_id, resource_id, action)
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		resource_id TEXT NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
		start_time DATETIME NOT NULL,
		end_time DATETIME,
		status TEXT NOT NULL,
		client_ip TEXT NOT NULL,
		client_metadata TEXT,
		audit_path TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS access_requests (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		resource_id TEXT NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
		reason TEXT NOT NULL,
		status TEXT NOT NULL,
		reviewer_id TEXT REFERENCES users(id),
		review_notes TEXT,
		requested_at DATETIME NOT NULL,
		reviewed_at DATETIME,
		expires_at DATETIME,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS ephemeral_credentials (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		resource_id TEXT NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		token TEXT UNIQUE,
		expires_at DATETIME NOT NULL,
		created_at DATETIME NOT NULL,
		used_at DATETIME
	);
	`

	_, err := db.Exec(createTablesSQL)
	if err != nil {
		return err
	}

	// Run additional migrations for existing databases
	err = runAdditionalMigrations(db)
	return err
}

func runAdditionalMigrations(db *sql.DB) error {
	// Check if 'type' column exists in resources table, if not add it
	var columnExists bool
	err := db.QueryRow("PRAGMA table_info(resources)").Scan()
	if err == nil {
		// Table exists, check for type column
		rows, err := db.Query("PRAGMA table_info(resources)")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var cid int
			var name, dataType string
			var notNull, pk bool
			var defaultValue interface{}

			err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk)
			if err != nil {
				return err
			}

			if name == "type" {
				columnExists = true
				break
			}
		}

		// Add type column if it doesn't exist
		if !columnExists {
			_, err = db.Exec("ALTER TABLE resources ADD COLUMN type TEXT")
			if err != nil {
				return fmt.Errorf("failed to add type column to resources table: %w", err)
			}
		}
	}

	return nil
}
