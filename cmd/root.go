package cmd

import (
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var rootCmd = &cobra.Command{
	Use:    "clockwise",
	Short:  "Clockwise is a meeting cost calculator designed to encourage more efficient meetings",
	Long:   "Clockwise is a meeting cost calculator designed to encourage more efficient meetings.",
	PreRun: toggleDebug,
	RunE:   rootFunc,
}

func rootFunc(cmd *cobra.Command, args []string) error {
	cmd.Help()
	os.Exit(1)

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "verbose logging")
}

// initConfig reads and/or initializes the configuration file.
func initConfig() {
	viper.SetDefault("averageSalary", 150000)

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	configFolder := home + "/.config/clockwise/"

	// Create config directory. Ignore errors (if it already exists).
	_ = os.MkdirAll(configFolder, os.ModePerm)

	viper.AddConfigPath(configFolder)
	viper.SetConfigName("clockwise") // Implicitly assumes .yaml extension.

	err = viper.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		// If the configuration file isn't found, create a new one.
		if ok {
			if err := viper.WriteConfigAs(configFolder + "clockwise.yaml"); err != nil {
				log.Println(err)
			}
		} else {
			log.Fatal(err)
		}
	}

	viper.AutomaticEnv() // Read in environment variables that match.
}
