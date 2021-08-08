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
		averageSalary := viper.GetViper().GetInt("averageSalary")
		currencySymbol := viper.GetViper().GetString("currencySymbol")

		qSalary := &survey.Question{
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

		qCurrencySymbol := &survey.Question{
			Prompt: &survey.Input{
				Message: "Set symbol or abbreviation of your local currency:",
				Default: "$",
			},
		}

		err := survey.AskOne(qSalary.Prompt, &averageSalary, survey.WithValidator(qSalary.Validate))
		if err != nil {
			return err
		}

		// No validation required for currencySymbol
		err = survey.AskOne(qCurrencySymbol.Prompt, &currencySymbol)
		if err != nil {
			return err
		}

		viper.GetViper().Set("averageSalary", averageSalary)
		viper.GetViper().Set("currencySymbol", currencySymbol)
		if err := viper.WriteConfig(); err != nil {
			return err
		}

		log.Printf(
			"The average annual salary of meeting participants has been updated to %s %v in the configuration file.",
			currencySymbol,
			averageSalary,
		)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
