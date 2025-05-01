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

type SpeedtestProvider interface {
	RunTest() (Result, error)
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

// RunTest performs a spedtest using a 3rd party library implementation
func (p *SpeedtestGoProvider) RunTest() (Result, error) {
	serverList, err := speedtest.FetchServers()
	if err != nil {
		return Result{}, fmt.Errorf("failed to fetch speedtest servers: %w", err)
	}

	if len(serverList) == 0 {
		return Result{}, fmt.Errorf("no speedtest servers found")
	}

	// Select bet server for the test
	targets, err := serverList.FindServer([]int{}) // []int{} is used to filter lowest latency
	if err != nil {
		return Result{}, fmt.Errorf("failed to find speedtest server: %w", err)
	}
	if len(targets) == 0 {
		return Result{}, fmt.Errorf("no suitable speedtest servers found")
	}

	server := targets[0] // Use the first (supposedely best) server
	log.Printf("Testing against server: %s (%s)\n", server.Name, server.Host)

	// Perform ping test
	err = server.PingTest(nil)
	if err != nil {
		return Result{}, fmt.Errorf("speedtest ping test failed: %w", err)
	}
	pingMs := float64(server.Latency.Milliseconds()) // Latency in Milliseconds, converted to float64

	// Perform download test
	err = server.DownloadTest()
	if err != nil {
		return Result{}, fmt.Errorf("speedtest download test failed: %w", err)
	}
	downMbps := server.DLSpeed.Mbps()

	// Perform upload test
	err = server.UploadTest()
	if err != nil {
		return Result{}, fmt.Errorf("speedtest download test failed: %w", err)
	}
	upMbps := server.ULSpeed.Mbps()

	return Result{
		DownMbps:  downMbps,
		UpMbps:    upMbps,
		PingMs:    pingMs,
		Timestamp: time.Now().UTC(),
	}, nil

}
