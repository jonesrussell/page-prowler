/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
	rootCmd.AddCommand(resultsCmd)

	// Set the default value of the siteid flag from the viper configuration
	if err := rootCmd.PersistentFlags().Lookup("siteid").Value.Set(viper.GetString("CRAWLSITEID")); err != nil {
		log.Fatalf("Failed to set siteid flag: %v", err)
	}

	// Define the siteid flag as a persistent flag on the rootCmd
	rootCmd.PersistentFlags().StringVarP(&Siteid, "siteid", "s", "", "Site ID")
}
