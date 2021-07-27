package scrape

import (
	"fmt"

	"github.com/mxschmitt/playwright-go"
)

// initializePlaywright starts playwright in a standalone function to circumvent
// some flaws in the upstream in terms of how it prints logs.
func InitializePlaywright() (pw *playwright.Playwright, err error) {
	pw, err = playwright.Run()
	if err != nil {
		return pw, fmt.Errorf("could not start playwright: %w", err)
	}

	return pw, nil
}
