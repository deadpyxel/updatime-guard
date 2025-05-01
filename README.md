# uptime-guard

![Go](https://github.com/deadpyxel/uptime-guard/actions/workflows/go.yml/badge.svg)
[![GitHub license](https://img.shields.io/github/license/deadpyxel/uptime-guard.svg)](https://github.com/deadpyxel/uptime-guard/blob/main/LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/deadpyxel/uptime-guard.svg?style=social&label=Star)](https://github.com/deadpyxel/uptime-guard)
[![Go Report Card](https://goreportcard.com/badge/github.com/deadpyxel/uptime-guard)](https://goreportcard.com/report/github.com/deadpyxel/uptime-guard)

Internet monitoring tool with speed test and downtime tracking.

## Overview

`uptime-guard` is a self-contained command-line tool written in Go to monitor your internet connection's speed and uptime. 
It periodically performs speed tests and connectivity checks, recording the results and downtime events in an SQLite database. 
This data can be used for analysis, helping you understand your internet service performance over time and potentially identify issues with your ISP.

> [!IMPORTANT]
> The project is currently under active development.

## Features

*   Periodic internet speed testing (Download, Upload, Ping).
*   Utilizes the `showwin/speedtest-go` library for native Go speed tests against Speedtest.net servers.
*   Periodic internet connectivity checks.
*   Detection and recording of internet downtime events with start and end times.
*   Persistence of speed test results and downtime events in a local SQLite database.
*   Configurable intervals for speed tests and connectivity checks.
*   Basic graceful shutdown on receiving termination signals (like Ctrl+C).

## Installation

Currently, `uptime-guard` is distributed using . To use it, you'll need to have Go installed (version 1.18 or higher recommended).

1.  **Install the tool:**

    ```bash
    go install github.com/deadpyxel/uptime-guard@latest
    ```

2.  **Run the application:**

    ```bash
    uptime-guard
    ```

    The tool will start monitoring your internet connection. A `uptimeguard.db` file will be created in the same directory to store the data.

    To stop the application, press `Ctrl+C`.

## Notes

`uptime-guard` is primarily a long-running background process.

*   **Running:** Execute `go run main.go` or build the executable and run it.
*   **Configuration:** Currently, the intervals for speed tests and connectivity checks are hardcoded in `main.go`. Future versions will include external configuration options.
*   **Data:** The collected data is stored in the `uptimeguard.db` SQLite file. You can use any SQLite browser tool (like [DB Browser for SQLite](https://sqlitebrowser.org/)) to view and analyze the data directly.

    *   `speed_tests` table: Contains records of each speed test performed.
    *   `downtime_events` table: Contains records of detected internet outages.

## Development

We welcome contributions! If you'd like to contribute to `uptime-guard`, here's how to get started:

1.  **Prerequisites:**
    *   Go (version 1.18+)
    *   Git

2.  **Clone the repository:**

    ```bash
    git clone https://github.com/deadpyxel/uptime-guard.git
    cd uptime-guard
    ```

3.  **Install dependencies:** Go Modules will automatically handle this when you build or run the project.

    ```bash
    go build
    ```

4.  **Run the code:**

    ```bash
    go run main.go
    ```

5.  **Run tests:**

    ```bash
    go test ./...
    ```

    This will run all tests in the project.

6.  **Code Style:** We generally follow standard Go idioms and best practices.

## Future Enhancements

*   External configuration options (e.g., using YAML or command-line flags) for intervals, database path, etc.
*   More robust error handling and logging.
*   Implementation of data retrieval functions to calculate and display statistics within the application.
*   A simple web frontend using HTMX and Templ for visualizing data.
*   Packaging the application for easier distribution (e.g., using GoReleaser).
*   Adding more comprehensive tests.

## Contributing

We welcome contributions from the community! Please see the [CONTRIBUTING.md](CONTRIBUTING.md) file (to be created) for details on how to contribute.

## Acknowledgements

This project utilizes the following excellent open-source libraries:

*   [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite driver for Go.
*   [github.com/showwin/speedtest-go](https://github.com/showwin/speedtest-go) - Go native Speedtest.net client library.

We are grateful to the maintainers and contributors of these projects.

## License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT). See the [LICENSE](LICENSE) file for details.
