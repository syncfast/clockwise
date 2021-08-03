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

		force_jitsi, err := cmd.Flags().GetBool("jitsi")
		if err != nil {
			return err
		}

		manual := false
		if url == "" {
			manual = true
		}

		var scrape_func = scrape.GetParticipantsZoom
		if !manual {
			// Checking optional force_jitsi flag first
			if strings.Contains(url, "meet.jit.si") || force_jitsi {
				scrape_func = scrape.GetParticipantsJitsi
			} else if strings.Contains(url, "zoom") {
				scrape_func = scrape.GetParticipantsZoom
			} else {
				return fmt.Errorf("Provided url does not contain known domain")
			}
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

			log.Info("Initializing TUI.")
			url, err := cmd.Flags().GetString("url")
			go func() {
				err = scrape_func(url, 1, &data, pw)
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
