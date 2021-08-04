package scrape

import (
	"fmt"

	"github.com/syncfast/clockwise/internal/tui"
	"github.com/mxschmitt/playwright-go"
)

// Function prototype for per-platform participant count scraping
type Scraper func(url string, refreshInterval int, data *tui.Data, pw *playwright.Playwright) error

// initializePlaywright starts playwright in a standalone function to circumvent
// some flaws in the upstream in terms of how it prints logs.
func InitializePlaywright() (pw *playwright.Playwright, err error) {
	pw, err = playwright.Run()
	if err != nil {
		return pw, fmt.Errorf("could not start playwright: %w", err)
	}

	return pw, nil
}
