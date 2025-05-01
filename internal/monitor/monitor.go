package monitor

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/deadpyxel/uptime-guard/internal/storage"
	"github.com/deadpyxel/uptime-guard/pkg/speedtest"
)

// Config holds the configuration info for the monitor
type Config struct {
	SpeedtestInterval time.Duration
	ConnCheckInterval time.Duration
	DBPath            string
}

// StartMonitoring begins the internet monitoring process
// This function will contain the main loop for scheduled tasks
func StartMonitoring(cfg Config, stProvider speedtest.SpeedtestProvider) {
	fmt.Println("Monitor started with config", cfg)

	// TODO: graceful shutdown
	stopCh := make(chan struct{})

	// Goroutine for scheduling speed tests
	go func() {
		ticker := time.NewTicker(cfg.SpeedtestInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Println("Running scheduled speed test...")
				result, err := stProvider.RunTest()
				if err != nil {
					log.Printf("error running speed test: %v\n", err)
					// TODO: How to handle this case
				}
				// Save result in the database
				err = storage.SaveSpeedtestResult(result.DownMbps, result.UpMbps, result.PingMs, result.Timestamp)
				if err != nil {
					log.Printf("Error saving speed test result: %v\n", err)
				}
			case <-stopCh:
				log.Println("Speed test scheduler stopped.")
				return
			}
		}
	}()

	// Go routine for scheduling connectivity checks and managing downtime
	go func() {
		ticker := time.NewTicker(cfg.ConnCheckInterval)
		defer ticker.Stop()

		isOnline := true            // Assume we are online initially
		var downtimeStart time.Time // Records the start of an outage
		var currDowntimeID int64    // records the database id of the current outage

		for {
			select {
			case <-ticker.C:
				log.Println("Perfoming scheduled connectivity check...")
				isConnected := CheckConnectivity()

				// Connection just came back up
				if isConnected && !isOnline {
					isOnline = true
					log.Println("Internet Connection restored...")

					// Recod the end time of the outage
					if downtimeStart.IsZero() {
						log.Println("Connection restored, but downtime start was nil")
					} else {
						now := time.Now().UTC()
						log.Printf("Outage ended at %s. Duration will be calculated on update.", now.Format(time.RFC3339))
						err := storage.UpdateDowntimeEnd(currDowntimeID, now)
						if err != nil {
							log.Printf("error updating downtime end for ID %d: %v", currDowntimeID, err)
						}
						// Reset downtime tracking
						downtimeStart = time.Time{}
						currDowntimeID = 0
					}
				} else if !isConnected && isOnline {
					// Connection just dropped
					isOnline = false
					log.Println("Internet connection lost.")

					// Record the start time of the outage
					downtimeStart = time.Now().UTC()
					log.Printf("Outage started at %s.", downtimeStart.Format(time.RFC3339))

					id, err := storage.SaveDowntimeStart(downtimeStart)
					if err != nil {
						log.Printf("error saving downtime start: %v\n", err)
					} else {
						currDowntimeID = id
						log.Printf("Downtime event saved with ID %d.", currDowntimeID)
					}
				} else if !isConnected && !isOnline {
					// Remains offline
					log.Println("Internet connection remains offline")
				} else {
					// Remains online
					log.Println("Internet connection remains online")
				}
			case <-stopCh:
				log.Println("connectivity checks stopped.")
				return
			}
		}
	}()
}

// CheckConnectivity checks if the internet connection is availabe.
// It attempts to establish a connection to a set of known reliable hosts
func CheckConnectivity() bool {
	// List of reliable hosts to check (e.g Public DNS servers)
	hosts := []string{
		"1.1.1.1:53",    // Cloudflare Public DNS
		"8.8.8.8:53",    // Google Public DNS
		"google.com:80", // Google HTTP
	}

	timeout := 5 * time.Second

	for _, host := range hosts {
		conn, err := net.DialTimeout("tcp", host, timeout)
		if err == nil {
			// Connection successful, internet is likely available
			conn.Close()
			return true
		}
	}
	// If all checks failed, assume the connection is down
	return false
}
