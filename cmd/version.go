package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Long:  "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		version := "0.0.1"
		fmt.Println(version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
