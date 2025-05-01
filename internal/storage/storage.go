package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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

// SaveSpeedtestResult creates a new entry on speed_tests with the result of the speedtest
func SaveSpeedtestResult(downMbps, upMbps, pingMs float64, timestamp time.Time) error {
	q := `
	INSERT INTO speed_tests (download_mbps, upload_mbps, ping_ms, timestamp)
	VALUES (?, ?, ?, ?);`

	_, err := db.Exec(q, downMbps, upMbps, pingMs, timestamp)
	if err != nil {
		return fmt.Errorf("failed to save speed test result: %w", err)
	}
	log.Printf("Saved speedtest result...")
	return nil
}

// SaveDowntimeStart creates a new entry in downtime_events setting the start_time
func SaveDowntimeStart(startTime time.Time) (int64, error) {
	q := `
	INSERT INTO downtime_events (start_time)
	VALUES (?);`

	result, err := db.Exec(q, startTime)
	if err != nil {
		return 0, fmt.Errorf("failed to save downtime start: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last inserted ID after saving downtime start: %w", err)
	}

	log.Printf("Saved downtime start at %s with ID %d\n", startTime.Format(time.RFC3339), id)
	return id, nil
}

// UpdateDowntimeEnd updates a downtime_event with matching id, setting the endtime and updating the duration
func UpdateDowntimeEnd(id int64, endTime time.Time) error {
	// Retrieve start time to calculate duration
	var startTime time.Time
	err := db.QueryRow("SELECT start_time FROM downtime_events WHERE id = ?", id).Scan(&startTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("downtime event with ID %d not found", id)
		}
		return fmt.Errorf("failed to get start time for downtime event with ID %d: %w", id, err)
	}

	duration := endTime.Sub(startTime).Seconds()

	q := `
	UPDATE downtime_events
	SET end_time = ?, duration_seconds = ?
	WHERE id = ?;`

	_, err = db.Exec(q, endTime, int(duration), id)
	if err != nil {
		return fmt.Errorf("failed to update downtime end for ID %d: %w", id, err)
	}

	log.Printf("Updated downtime event ID %d with end time %s and duration %.2f",
		id, endTime.Format(time.RFC3339), duration)
	return nil
}

// GetAvgSpeed returns the average download and upload speed in a given timeframe
func GetAvgSpeed(startTime, endTime time.Time) (float64, float64, error) {
	q := `
	SELECT AVG(download_mbps), AVG(upload_mbps)
	FROM speed_tests
	WHERE timestamp BETWEEN ? AND ?;`

	var avgDown, avgUp sql.NullFloat64 // Use NullFloat64 for handling no results
	err := db.QueryRow(q, startTime, endTime).Scan(&avgDown, &avgUp)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get average speed: %w", err)
	}

	// Handle cases where there are no records in the time range
	downSpd := 0.0
	if avgDown.Valid {
		downSpd = avgDown.Float64
	}
	upSpd := 0.0
	if avgUp.Valid {
		upSpd = avgUp.Float64
	}

	return downSpd, upSpd, nil
}

// GetTotalDowntime retrieves the sum in seconds of all downtime in a timeframe
func GetTotalDowntime(startTime, endTime time.Time) (int, error) {
	q := `
	SELECT SUM(duration_seconds)
	FROM downtime_events
	WHERE start_time BETWEEN ? AND ? AND end_time IS NOT NULL;`

	var totalDuration sql.NullInt64 // Use NullInt64 for handling no results

	err := db.QueryRow(q, startTime, endTime).Scan(&totalDuration)
	if err != nil {
		return 0, fmt.Errorf("failed to get total downtime: %w", err)
	}

	duration := 0
	if totalDuration.Valid {
		duration = int(totalDuration.Int64)
	}

	return duration, nil
}
