# Contributing to uptime-guard

Thank you for your interest in contributing to `uptime-guard`! We appreciate your help in making this project better.

This document outlines the guidelines and process for contributing to `uptime-guard`.

## How to Contribute

There are several ways you can contribute:

*   **Reporting Bugs:** If you find a bug, please open an issue on GitHub with a clear description and steps to reproduce it.
*   **Suggesting Features:** Have an idea for a new feature or improvement? Open an issue to discuss it.
*   **Writing Code:** Implement new features, fix bugs, improve documentation, or write tests.
*   **Improving Documentation:** Found a typo, unclear explanation, or missing information in the README or other documentation? Help us improve it!
*   **Testing:** Test the application on different operating systems or environments and report any issues.

## Getting Started

1.  **Fork the repository:** Click the "Fork" button on the GitHub page.
2.  **Clone your fork:**
    ```bash
    git clone https://github.com/YOUR_GITHUB_USERNAME/uptime-guard.git
    cd uptime-guard
    ```
3.  **Add the original repository as a remote:**
    ```bash
    git remote add upstream https://github.com/deadpyxel/uptime-guard.git
    ```
4.  **Install dependencies:** Go Modules will handle this.
    ```bash
    go mod download
    ```
5.  **Set up your development environment:** Ensure you have Go installed (version 1.18+ recommended).

## Workflow for Code Contributions

1.  **Sync with the latest `main` branch:**
    ```bash
    git checkout main
    git pull upstream main
    ```
2.  **Create a new branch:** Give your branch a descriptive name related to the issue or feature you're working on (e.g., `fix/database-panic`, `feat/yaml-config`).
    ```bash
    git switch -c your-branch-name
    ```
3.  **Make your changes:** Write your code, add tests, and update documentation as needed.
4.  **Test your changes:** Run all tests to ensure nothing is broken.
    ```bash
    go test ./...
    ```
5.  **Format your code:** Use `gofmt` to ensure consistent code style.
    ```bash
    gofmt -w .
    ```
6.  **Consider linting:** Running linters like `golint` or `staticcheck` can help catch potential issues.
    ```bash
    # Example using golint (install if needed: go install golang.org/x/lint/golint@latest)
    golint ./...
    ```
7.  **Commit your changes:** Write clear and concise commit messages. Follow conventional commits if possible (e.g., `feat: add yaml configuration`, `fix: resolve database connection issue`).
    ```bash
    git add .
    git commit -m "feat: describe your changes"
    ```
8.  **Push your branch to your fork:**
    ```bash
    git push origin your-branch-name
    ```
9.  **Open a Pull Request (PR):** Go to the original `uptime-guard` repository on GitHub and open a PR from your branch.
    *   Provide a clear title and description of your changes.
    *   Reference any related issues in the PR description (e.g., `Fixes #123`, `Implements #456`).

## Code Style

*   Follow standard Go idioms and best practices as outlined in the [Effective Go](https://go.dev/doc/effective_go) and [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments).
*   Use `gofmt` to format your code.
*   Aim for clear, readable, and well-commented code where necessary.

## Testing

*   Include tests for new features and bug fixes.
*   Ensure that your changes pass existing tests.
*   Write tests that are clear, concise, and easy to understand.
*   For testing database interactions, prefer using an in-memory SQLite database where possible.

## Security

If you discover a security vulnerability, please report it responsibly by contacting the project maintainer directly (e.g., via email or any linked social network) instead of opening a public issue.

## License

By contributing to `uptime-guard`, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).

Thank you for contributing to `uptime-guard`!
