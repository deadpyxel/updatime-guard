package storage

import (
	"database/sql"
	"os"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Open an in-memory database
	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	// Assign the test DB to the package level db variable for tests to use
	// this is a pragmatic approach for testing this specific structure.
	// For more complex scenarios dependency injection myght be a better option
	db = testDB

	// Create tables in the in-memory database
	err = createTables()
	if err != nil {
		t.Fatalf("failed to create tables in in-memory database: %v", err)
	}

	return testDB
}

func teardownTestDB(testDB *sql.DB) {
	if testDB != nil {
		testDB.Close()
		// Reset the package-level db variable
		db = nil
	}
}

// Helper function for flaot comparison
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func TestInitializeDB(t *testing.T) {
	// Test with a temporary file path
	dbPath := t.TempDir() + "/test.db"
	err := InitializeDB(dbPath)
	if err != nil {
		t.Fatalf("InitializeDB failed: %v", err)
	}

	defer func() {
		// Close the database initialized by InitializeDB
		closeErr := CloseDB()
		if closeErr != nil {
			t.Errorf("Error closing DB in TestInitializeDB: %v", err)
		}
	}()

	// Check if the file was created (basic check)
	_, err = os.Stat(dbPath)
	if os.IsNotExist(err) {
		t.Errorf("database file was not created at %s", dbPath)
	}

	// TODO: Add tests for table structure
}

func TestSaveSpeedtestResult(t *testing.T) {
	testDB := setupTestDB(t)
	defer teardownTestDB(testDB)

	download := 150.5
	upload := 30.2
	ping := 10.1
	timestamp := time.Now().UTC() // Use UTC for consitentcy

	err := SaveSpeedtestResult(download, upload, ping, timestamp)
	if err != nil {
		t.Errorf("SaveSpeedtestResult failed with error: %v", err)
	}

	// Verify the data was inserted
	var count int
	err = testDB.QueryRow("SELECT COUNT(*) FROM speed_tests;").Scan(&count)
	if err != nil {
		t.Errorf("Failed to count saved speed test results: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 speed test result, got %d", count)
	}

	// TODO: Check inserted values for consistency (are we testing the database layer?)
}

func TestSaveDowntime(t *testing.T) {
	testDB := setupTestDB(t)
	defer teardownTestDB(testDB)

	startTime := time.Now().UTC().Add(-10 * time.Minute) // set the start time for 10 minutes prior
	var savedID int64

	t.Run("SaveDowntimeStart creates new entry with correct start_time", func(t *testing.T) {
		id, err := SaveDowntimeStart(startTime)
		if err != nil {
			t.Errorf("SaveDowntimeStart failed with error: %v", err)
		}
		savedID = id

		// Verify the start time was inserted
		var savedStartTime time.Time
		var endTime sql.NullTime
		err = testDB.QueryRow("SELECT start_time, end_time FROM downtime_events WHERE id = ?;", id).Scan(
			&savedStartTime, &endTime)
		if err != nil {
			t.Fatalf("failed to get saved downtime start: %v", err)
		}
		if !savedStartTime.Equal(startTime) {
			t.Errorf("Saved start time %s does not match expected %s", savedStartTime, startTime)
		}
		if endTime.Valid {
			t.Errorf("Expected end_time to be NULL, got %v instead", endTime)
		}
	})

	t.Run("UpdateDowntimeEnd updated end_time and duration without changing start_time", func(t *testing.T) {
		endTime := time.Now().UTC().Add(-2 * time.Minute) // 2 minutes ago
		err := UpdateDowntimeEnd(savedID, endTime)
		if err != nil {
			t.Fatalf("UpdateDowntimeEnd failed with error: %v", err)
		}

		// Verify the end time and duration were updated
		var updatedEndTime sql.NullTime
		var duration int
		err = testDB.QueryRow("SELECT end_time, duration_seconds FROM downtime_events WHERE id = ?", savedID).Scan(
			&updatedEndTime, &duration)
		if err != nil {
			t.Fatalf("failed to get updated downtime event: %v", err)
		}
		if !updatedEndTime.Valid || !updatedEndTime.Time.Equal(endTime) {
			t.Errorf("Updated end time %s does not match expected %s",
				updatedEndTime.Time.Format(time.RFC3339Nano), endTime.Format(time.RFC3339Nano))
		}

		expectedDuration := int(endTime.Sub(startTime).Seconds())
		if duration != expectedDuration {
			t.Errorf("Expected duration %d, got %d", expectedDuration, duration)
		}
	})

	// TODO: add test to verify sequential calls of UpdateDowntimeEnd update the end_time every time
}

func TestGetAVGSpeed(t *testing.T) {
	testDB := setupTestDB(t)
	defer teardownTestDB(testDB)

	now := time.Now().UTC()
	t.Run("GetAvgSpeed correctly computes resutls for period with entries", func(t *testing.T) {
		// Insert some test data
		_ = SaveSpeedtestResult(100, 20, 15, now.Add(-time.Hour))
		_ = SaveSpeedtestResult(120, 25, 15, now.Add(-30*time.Minute))
		_ = SaveSpeedtestResult(80, 15, 20, now.Add(-15*time.Minute))
		_ = SaveSpeedtestResult(100, 20, 10, now.Add(-5*time.Minute))

		avgDL, avgUP, err := GetAvgSpeed(now.Add(-2*time.Hour), now)
		if err != nil {
			t.Fatalf("GetAvgSpeed failed with error: %v", err)
		}

		// Expected averages: 100.0 DL, 20.0 UP, 15 Ping
		// Use a tolerance for float comparison
		if abs(avgDL-100.0) > 0.001 {
			t.Errorf("Expected average download speed 100.0, got %.2f", avgDL)
		}
		if abs(avgUP-20.0) > 0.001 {
			t.Errorf("Expected average upload speed 20.0, got %.2f", avgUP)
		}
	})

	t.Run("GetAvgSpeed returns 0 for both speeds in period with no data", func(t *testing.T) {
		avgDL, avgUP, err := GetAvgSpeed(now.Add(-3*time.Hour), now.Add(-2*time.Hour-time.Minute))
		if err != nil {
			t.Fatalf("GetAvgSpeed failed with error: %v", err)
		}

		if avgDL != 0 || avgUP != 0 {
			t.Errorf("Expected average speed 0 for time range with no data, got %.2f DL and %.2f UP instead", avgDL, avgUP)
		}
	})
}

func TestGetTotalDowntime(t *testing.T) {
	testDB := setupTestDB(t)
	defer teardownTestDB(testDB)

	// Insert some test downtime data
	now := time.Now().UTC()

	// Completed outage 1 (5 minutes)
	id1, _ := SaveDowntimeStart(now.Add(-30 * time.Minute))
	_ = UpdateDowntimeEnd(id1, now.Add(-25*time.Minute)) // Duration ~300 seconds

	// Completed outage 2 (10 minutes)
	id2, _ := SaveDowntimeStart(now.Add(-20 * time.Minute))
	_ = UpdateDowntimeEnd(id2, now.Add(-10*time.Minute)) // Duration ~600 seconds

	// Ongoing outage (should not be counted by default)
	_, _ = SaveDowntimeStart(now.Add(-5 * time.Minute))

	// Total downtime in the last hour
	totalDT, err := GetTotalDowntime(now.Add(-1*time.Hour), now)
	if err != nil {
		t.Fatalf("GetTotalDowntime failed: %v", err)
	}

	// Expected total duration = 300 + 600 = 900 seconds
	expectedTotalDuration := 900
	if expectedTotalDuration != totalDT {
		t.Errorf("Expected total downtime %d seconds, got %d", expectedTotalDuration, totalDT)
	}

	// Test with a time range with no completed outages
	totalDTNoData, err := GetTotalDowntime(now.Add(-40*time.Minute), now.Add(-35*time.Minute))
	if err != nil {
		t.Fatalf("GetTotalDowntime with no data failed: %v", err)
	}
	if totalDTNoData != 0 {
		t.Errorf("Expected total downtime 0 for time range with no data, got %d", totalDTNoData)
	}
}
