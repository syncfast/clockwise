package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const version = "0.0.1"

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:    "version",
	Short:  "Print version",
	Long:   "Print version",
	PreRun: toggleDebug,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
