package cmd

import (
	"strconv"

	survey "github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setCmd represents the set command.
var setCmd = &cobra.Command{
	Use:    "set",
	Short:  "Set the average annual salary of meeting participants",
	Long:   `Set the average annual salary of meeting participants. This does not need to be an exact number.`,
	PreRun: toggleDebug,
	RunE: func(cmd *cobra.Command, args []string) error {
		averageSalary := viper.GetViper().GetInt("averageSalary")

		q := &survey.Question{
			Prompt: &survey.Input{
				Message: "Set average annual salary of meeting participants:",
				Default: strconv.Itoa(averageSalary),
			},
			Validate: func(val interface{}) error {
				if _, err := strconv.Atoi(val.(string)); err != nil {
					return err
				}

				return nil
			},
		}

		survey.AskOne(q.Prompt, &averageSalary, survey.WithValidator(q.Validate))

		viper.GetViper().Set("averageSalary", averageSalary)
		if err := viper.WriteConfig(); err != nil {
			return err
		}

		log.Printf("The average annual salary of meeting participants has been updated to %v in the configuration file.", averageSalary)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
