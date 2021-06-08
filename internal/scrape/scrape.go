// Package scrape is a web scraper that relies on playwright-go to continuously
// retrieve the participant count of a specified zoom meeting.
package scrape

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mxschmitt/playwright-go"
	"github.com/syncfast/clockwise/internal/tui"
)

// initializePlaywright initializes playwright to circumvent some flaws in the
// upstream in terms of how it prints logs.
func InitializePlaywright() (pw *playwright.Playwright, err error) {
	pw, err = playwright.Run()
	if err != nil {
		return pw, fmt.Errorf("could not start playwright: %w", err)
	}

	return pw, nil
}

// getParticipants retrieves the total participant count from a specified zoom
// URL. It runs in a loop and updates the passed in `Data` struct every
// `refreshInterval` seconds.
func GetParticipants(url string, refreshInterval int, data *tui.Data, pw *playwright.Playwright) error {
	var timeout float64 = 5000

	url = mutateURL(url)

	browser, err := pw.Chromium.Launch()
	if err != nil {
		return fmt.Errorf("could not launch browser: %w", err)
	}

	page, err := browser.NewPage()
	if err != nil {
		return fmt.Errorf("could not create page: %w", err)
	}

	if _, err = page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
	}); err != nil {
		return fmt.Errorf("could not goto: %w", err)
	}

	selector := "text=Your Name"
	if err := page.Fill(selector, "clockwise", playwright.FrameFillOptions{
		Timeout: &timeout,
	}); err != nil {
		return err
	}

	page.WaitForSelector("button#joinBtn")

	if err := page.Click("button#joinBtn", playwright.PageClickOptions{
		Timeout: &timeout,
	}); err != nil {
		return err
	}

	page.WaitForSelector(".footer-button__number-counter")

	for {
		res, err := page.QuerySelector(".footer-button__number-counter")
		if err != nil {
			return err
		}

		span, err := res.InnerHTML()
		if err != nil {
			return err
		}

		stringCount := span[6 : len(span)-7]
		count, err := strconv.Atoi(stringCount)
		if err != nil {
			return err
		}

		data.SetCount(count)

		time.Sleep(time.Second * time.Duration(refreshInterval))
	}

	// TODO: Revisit this. I might be able to just remove it altogether.
	if err = browser.Close(); err != nil {
		return fmt.Errorf("could not close browser: %w", err)
	}

	if err = pw.Stop(); err != nil {
		return fmt.Errorf("could not stop Playwright: %w", err)
	}

	return nil
}

// mutateURL converts the generic meeting URL into the browser-specific URL.
func mutateURL(url string) string {
	return strings.Replace(url, "/j/", "/wc/join/", 1)
}
