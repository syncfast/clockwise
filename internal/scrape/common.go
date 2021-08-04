package scrape

import (
	"fmt"

	"github.com/mxschmitt/playwright-go"
	"github.com/syncfast/clockwise/internal/tui"
)

type MeetingImpl interface {
	VisitMeetingUrl() error
	FillBotName(botName string) error
	JoinMeeting() error
	ActivateVirtualWebcam(camName string) error
	GetParticipants(refreshInterval int, data *tui.Data) error
}

// initializePlaywright starts playwright in a standalone function to circumvent
// some flaws in the upstream in terms of how it prints logs.
func InitializePlaywright() (pw *playwright.Playwright, err error) {
	pw, err = playwright.Run()
	if err != nil {
		return pw, fmt.Errorf("could not start playwright: %w", err)
	}

	return pw, nil
}
