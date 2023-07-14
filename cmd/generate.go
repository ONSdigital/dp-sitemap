/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/ONSdigital/dp-sitemap/cmd/utilities"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate sitemap files using default parametssers",
	Long:  `A tool to generate the sitemap`,
	Run: func(cmd *cobra.Command, args []string) {
		//rf, _ := cmd.Flags().GetString("robots-file-path")
		sm, _ := cmd.Flags().GetString("sitemap-index")

		fmt.Print(cmd)
		fmt.Print("sm is : ", sm)

		cfg, err := config.Get()
		if err != nil {
			fmt.Println("Error retrieving config" + err.Error())
			os.Exit(1)
		}

		_, commandLine := utilities.ValidateCommandLines()
		// if !commandLine.valid {
		// 	os.Exit(1)
		// }

		// if commandLine.generate_sitemap {
		utilities.GenerateSitemap(cfg, commandLine)
		//		}

		// if commandLine.update_sitemap {
		// 	// Your update sitemap code here
		// }

		// GenerateRobotFile(cfgxs, commandLine)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.PersistentFlags().String("robots-file-path", "test_robots.txt", "path to robots file")
	generateCmd.PersistentFlags().String("sitemap-file-path", "test_sitemap", "path to sitemap file")
	generateCmd.PersistentFlags().String("api-url", "http://localhost", "elastic search api url")
	generateCmd.PersistentFlags().String("zebedee-url", "http://localhost:8082", "zebedee url")
	generateCmd.PersistentFlags().String("sitemap-index", "1", "OPENSEARCH_SITEMAP_INDEX")
	generateCmd.PersistentFlags().String("scroll-timeout", "2000", "OPENSEARCH_SCROLL_TIMEOUT")
	generateCmd.PersistentFlags().Int("scroll-size", 10, "OPENSEARCH_SCROLL_SIZE")
	generateCmd.PersistentFlags().Bool("enable-fake-scroll", true, "enable fake scroll")
	generateCmd.PersistentFlags().Bool("generate-sitemap", false, "generate the sitemap")
	generateCmd.PersistentFlags().Bool("update-sitemap", false, "update the sitemap")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
