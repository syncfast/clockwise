package scrape

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mxschmitt/playwright-go"
	"github.com/syncfast/clockwise/internal/tui"
)

type Zoom struct {
	url     string
	pw      *playwright.Playwright
	page    playwright.Page
	timeout float64
}

func NewZoom(url string, pw *playwright.Playwright) *Zoom {
	return &Zoom{
		url:     url,
		pw:      pw,
		page:    nil,
		timeout: 5000,
	}
}

func (z *Zoom) VisitMeetingUrl() error {
	if strings.Contains(z.url, "zoom.us/my/") {
		return fmt.Errorf(`Error: clockwise is not compatible with Zoom Personal Meeting IDs at the moment.
			Disabling your PMI is as as simple as clicking a checkbox.
			Please visit https://support.zoom.us/hc/en-us/articles/203276937-Using-Personal-Meeting-ID-PMI- for more info.`)
	}

	z.url = mutateURL(z.url)

	browser, err := z.pw.Chromium.Launch()
	if err != nil {
		return fmt.Errorf("could not launch browser: %w", err)
	}

	page, err := browser.NewPage()
	if err != nil {
		return fmt.Errorf("could not create page: %w", err)
	}

	if _, err = page.Goto(z.url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
	}); err != nil {
		return fmt.Errorf("could not goto: %w", err)
	}

	return nil
}

func (z *Zoom) FillBotName(botName string) error {
	selector := "text=Your Name"
	if err := z.page.Fill(selector, "clockwise-bot", playwright.FrameFillOptions{
		Timeout: &z.timeout,
	}); err != nil {
		return err
	}

	return nil
}

func (z *Zoom) JoinMeeting() error {
	element, err := z.page.WaitForSelector("button#joinBtn")
	if err != nil {
		return fmt.Errorf("failed to wait for join button: %w", err)
	}

	if err := element.Click(playwright.ElementHandleClickOptions{
		Timeout: &z.timeout,
	}); err != nil {
		return err
	}

	return nil
}

func (z *Zoom) ActivateVirtualWebcam(camName string) error {
	return nil
}

// GetParticipants retrieves the total participant count from a specified
// zoom URL. It runs in a loop and updates the passed in `Data` struct every
// `refreshInterval` seconds.
func (z *Zoom) GetParticipants(refreshInterval int, data *tui.Data) error {
	_, err := z.page.WaitForSelector(".footer-button__number-counter")
	if err != nil {
		return fmt.Errorf("failed to wait for join button: %w", err)
	}

	for {
		res, err := z.page.QuerySelector(".footer-button__number-counter")
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

		data.SetCount(count - 1)

		time.Sleep(time.Second * time.Duration(refreshInterval))
	}
}

// mutateURL converts the generic meeting URL into the browser-specific URL.
func mutateURL(url string) string {
	return strings.Replace(url, "/j/", "/wc/join/", 1)
}
