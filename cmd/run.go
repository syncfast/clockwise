package cmd

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/syncfast/clockwise/internal/scrape"
	"github.com/syncfast/clockwise/internal/tui"
)

// runCmd represents the run command.
var runCmd = &cobra.Command{
	Use:    "run",
	Short:  "Run clockwise",
	Long:   ``,
	PreRun: toggleDebug,
	RunE: func(cmd *cobra.Command, args []string) error {
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			return err
		}

		forceJitsi, err := cmd.Flags().GetBool("jitsi")
		if err != nil {
			return err
		}

		manual := false
		if url == "" {
			manual = true
		}

		// We declare data here because it's consumed by both the `tui` and
		// `scrape` packages.
		var data tui.Data

		if !manual {
			log.Println("Initializing playwright to scrape participant count.")
			pw, err := scrape.InitializePlaywright()
			if err != nil {
				return err
			}

			// Checking optional force_jitsi flag first
			var meetingImpl scrape.MeetingImpl
			switch {
			case forceJitsi || strings.Contains(url, "meet.jit.si"):
				meetingImpl = scrape.NewJitsi(url, pw)
			case strings.Contains(url, "zoom"):
				meetingImpl = scrape.NewZoom(url, pw)
			default:
				return fmt.Errorf("Provided url does not contain known domain")
			}

			log.Info("Initializing TUI.")
			go func() {
				meetingImpl.VisitMeetingUrl()
				meetingImpl.FillBotName("clockwise-bot")
				meetingImpl.JoinMeeting()
				// FIXME: Deactivated until ffmpeg vcam gets implemented
				// meetingImpl.ActivateVirtualWebcam("")

				err = meetingImpl.GetParticipants(1, &data)
				if err != nil {
					log.Fatal(err)
				}
			}()
		}

		tui.Start(manual, &data)
		log.Info("Clockwise has been stopped.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("url", "u", "", "Meeting URL")
	runCmd.Flags().Bool("jitsi", false, "Force Jitsi URL scraping")
}
