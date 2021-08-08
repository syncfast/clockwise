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
	Short:  "Set the average annual salary of meeting participants and currency representation",
	Long:   `Set the average annual salary of meeting participants. This does not need to be an exact number.`,
	PreRun: toggleDebug,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Fetch currently set values from config or default values
		averageSalaryPrev := viper.GetViper().GetInt("averageSalary")
		currencySymbolPrev := viper.GetViper().GetString("currencySymbol")

		q := []*survey.Question{
			{
				Name: "averageSalary",
				Prompt: &survey.Input{
					Message: "Set average annual salary of meeting participants:",
					Default: strconv.Itoa(averageSalaryPrev),
				},
				Validate: func(val interface{}) error {
					if _, err := strconv.Atoi(val.(string)); err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name: "currencySymbol",
				Prompt: &survey.Input{
					Message: "Set symbol or abbreviation of your local currency:",
					Default: currencySymbolPrev,
				},
			},
		}

		answers := struct {
			AverageSalary  int
			CurrencySymbol string
		}{}

		err := survey.Ask(q, &answers)
		if err != nil {
			return err
		}

		viper.GetViper().Set("averageSalary", answers.AverageSalary)
		viper.GetViper().Set("currencySymbol", answers.CurrencySymbol)
		if err := viper.WriteConfig(); err != nil {
			return err
		}

		log.Printf(
			"The average annual salary of meeting participants has been updated to %s %v in the configuration file.",
			answers.CurrencySymbol,
			answers.AverageSalary,
		)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
