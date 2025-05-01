package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

var db *sql.DB

// InitializeDB sets up the SQLite database.
func InitializeDB(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	// Ping the database to ensure connection is working
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Database initialized at %s...\nCreating table structure...", dbPath)
	return createTables()
}

// CloseDB closes the database connection.
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func createTables() error {
	speedTestsCreateTableSQL := `
	CREATE TABLE IF NOT EXISTS speed_tests (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    download_mbps REAL NOT NULL,
	    upload_mbps REAL NOT NULL,
	    ping_ms REAL NOT NULL,
	    timestamp DATETIME NOT NULL
	);`

	downtimeCreateTableSQL := `
	CREATE TABLE IF NOT EXISTS downtime_events (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    start_time DATETIME NOT NULL,
	    end_time DATETIME, -- Nullable if the outage is ongoing
	    duration_seconds INTEGER -- Can calculate this, but storing might be more convenient
	);`

	_, err := db.Exec(speedTestsCreateTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create speed_tests table: %w", err)
	}

	_, err = db.Exec(downtimeCreateTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create downtime_events table: %w", err)
	}

	log.Println("Database tables created or already exist.")

	return nil
}
