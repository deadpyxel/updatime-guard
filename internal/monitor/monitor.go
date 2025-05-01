package monitor

import (
	"fmt"
	"log"
	"net"
	"time"
)

// StartMonitoring begins the internet monitoring process
// This function will contain the main loop for scheduled tasks
func StartMonitoring() {
	fmt.Println("Monitor started...")

	// TODO: Main scheduling loop iomplementation

	// WIP: Perform connectivity check once to verify logic
	isConnected := CheckConnectivity()
	if isConnected {
		log.Println("connectivity check: UP")
	} else {
		log.Println("connectivity check: DOWN")
	}

	// monitoring loop goes here
	// likely to run as a goroutine
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
		log.Printf("connectivity test failed for host %s: %v\n", host, err)
	}
	// If all checks failed, assume the connection is down
	return false
}
