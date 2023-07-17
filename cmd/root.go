package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dp-sitemap",
	Short: "dp-sitemap is a cli tool to generate sitemaps easily & locally",
	Long: `dp-sitemap is a comprehensive tool for generating and managing sitemaps of your website. 
	The tool provides functionality for creating sitemaps from different data sources, updating existing sitemaps, and integrating with external services

Sitemap-cli is a CLI tool for Go that empowers applications.
This application is a tool to generate the needed files
to quickly.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default s $HOME/.dp-sitemap.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
