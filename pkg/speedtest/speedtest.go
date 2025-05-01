package speedtest

import (
	"fmt"
	"log"
	"time"

	"github.com/showwin/speedtest-go/speedtest"
)

// Result holds the result of a speedtest
type Result struct {
	DownMbps  float64
	UpMbps    float64
	PingMs    float64
	Timestamp time.Time
}

func RunSpeedtest() (Result, error) {
	fmt.Println("Running simulated speed test...")

	return Result{
		DownMbps:  100.0,
		UpMbps:    20.0,
		PingMs:    15.0,
		Timestamp: time.Now(),
	}, nil
}

// TODO: Add other providers for cases we want to use the CLI directly
// or some tool that is able to interface with ours given properly parsable output
// The idea of providers is that we would be able to further extend the supported speed test solutiosn we have

// SpeedtestGoProvider is a speed test provider that uses the showwin/speedtest-go library
type SpeedtestGoProvider struct {
	user *speedtest.User
}

// NewSpeedtestGoProvider creats a new SpeedtestGoProvider
func NewSpeedtestGoProvider() (*SpeedtestGoProvider, error) {
	user, err := speedtest.FetchUserInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch speedtest user info: %w", err)
	}

	log.Printf("Speedtest user info: %+v\n", user)
	return &SpeedtestGoProvider{user: user}, nil
}
