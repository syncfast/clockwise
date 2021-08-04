package scrape

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/mxschmitt/playwright-go"
	"github.com/syncfast/clockwise/internal/tui"
)

type Jitsi struct {
	url     string
	pw      *playwright.Playwright
	page    playwright.Page
	timeout float64
}

func NewJitsi(url string, pw *playwright.Playwright) *Jitsi {
	return &Jitsi{
		url:     url,
		pw:      pw,
		page:    nil,
		timeout: 5000,
	}
}

func (j *Jitsi) VisitMeetingUrl() error {
	browser, err := j.pw.Chromium.Launch()
	if err != nil {
		return fmt.Errorf("could not launch browser: %w", err)
	}

	j.page, err = browser.NewPage()
	if err != nil {
		return fmt.Errorf("could not create page: %w", err)
	}

	if _, err = j.page.Goto(j.url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
	}); err != nil {
		return fmt.Errorf("could not goto: %w", err)
	}

	return nil
}

func (j *Jitsi) FillBotName(botName string) error {
	selector := "#Prejoin-input-field-id"
	if err := j.page.Fill(selector, botName, playwright.FrameFillOptions{
		Timeout: &j.timeout,
	}); err != nil {
		return err
	}

	return nil
}

func (j *Jitsi) JoinMeeting() error {
	// Wait for and click Join button
	element, err := j.page.WaitForSelector("#lobby-screen > div.content > div.prejoin-input-area-container > div > div > div")
	if err != nil {
		return fmt.Errorf("failed to wait for join button: %w", err)
	}

	if err := element.Click(playwright.ElementHandleClickOptions{
		Timeout: &j.timeout,
	}); err != nil {
		return err
	}

	return nil
}

func (j *Jitsi) ActivateVirtualWebcam(camName string) error {
	return nil
}

// GetParticipants retrieves the total participant count from a specified
// Jitsi URL. It runs in a loop and updates the passed in `Data` struct every
// `refreshInterval` seconds.
func (j *Jitsi) GetParticipants(refreshInterval int, data *tui.Data) error {
	// Wait for and click participants sidebar
	element, err := j.page.WaitForSelector("#new-toolbox > div > div > div > div:nth-child(6)")
	if err != nil {
		return fmt.Errorf("failed to wait for participant sidebar button: %w", err)
	}

	if err := element.Click(playwright.ElementHandleClickOptions{
		Timeout: &j.timeout,
	}); err != nil {
		return err
	}

	_, err = j.page.WaitForSelector("#layout_wrapper > div.participants_pane > div")
	if err != nil {
		return fmt.Errorf("failed to wait for participant sidebar: %w", err)
	}

	for {
		res, err := j.page.QuerySelector("#layout_wrapper > div.participants_pane > div")
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
