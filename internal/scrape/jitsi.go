package scrape

import (
	"fmt"
	"strconv"
	"time"
	"regexp"

	"github.com/mxschmitt/playwright-go"
	"github.com/syncfast/clockwise/internal/tui"
)

// getParticipants retrieves the total participant count from a specified Jitsi
// URL. It runs in a loop and updates the passed in `Data` struct every
// `refreshInterval` seconds.
func GetParticipantsJitsi(url string, refreshInterval int, data *tui.Data, pw *playwright.Playwright) error {
	var timeout float64 = 5000

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

	selector := "#Prejoin-input-field-id"
	if err := page.Fill(selector, "clockwise-bot", playwright.FrameFillOptions{
		Timeout: &timeout,
	}); err != nil {
		return err
	}

	// Wait for and click Join button
	page.WaitForSelector("#lobby-screen > div.content > div.prejoin-input-area-container > div > div > div")

	if err := page.Click("#lobby-screen > div.content > div.prejoin-input-area-container > div > div > div", playwright.PageClickOptions{
		Timeout: &timeout,
	}); err != nil {
		return err
	}

	// Wait for and click participants sidebar
	page.WaitForSelector("#new-toolbox > div > div > div > div:nth-child(6)")
	if err := page.Click("#new-toolbox > div > div > div > div:nth-child(6)", playwright.PageClickOptions{
		Timeout: &timeout,
	}); err != nil {
		return err
	}

	page.WaitForSelector("#layout_wrapper > div.participants_pane > div")

	for {
		res, err := page.QuerySelector("#layout_wrapper > div.participants_pane > div")
		if err != nil {
			return err
		}

		span, err := res.InnerHTML()
		if err != nil {
			return err
		}

		re := regexp.MustCompile(`Meeting participants \(([0-9]+)\)`)
		match_str := re.FindStringSubmatch(span)

		count, err := strconv.Atoi(match_str[1])
		if err != nil {
			return err
		}

		data.SetCount(count - 1)

		time.Sleep(time.Second * time.Duration(refreshInterval))
	}
}
