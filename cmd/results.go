package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Siteid string

// resultsCmd represents the results command
var resultsCmd = &cobra.Command{
	Use:   "results",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Could not read config file")
	}

	// Bind the environment variable to the flag
	err = viper.BindEnv("siteid")
	if err != nil {
		log.Fatalf("Failed to bind env var: %v", err)
	}

	// Define the siteid flag and set its default value from the environment variable
	resultsCmd.PersistentFlags().StringVarP(&Siteid, "siteid", "s", viper.GetString("siteid"), "Set siteid for redis set key")

	RootCmd.AddCommand(resultsCmd)
}
