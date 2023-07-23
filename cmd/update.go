package cmd

import (
	"fmt"
	"os"

	"github.com/ONSdigital/dp-sitemap/cmd/utilities"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Generate sitemap files using default parameters",
	Long:  `A tool to update generate and update sitemap files`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Get()
		if err != nil {
			fmt.Println("Error retrieving config"+err.Error(), cfg)
			os.Exit(1)
		}
		viper.BindPFlags(cmd.Flags()) // Bind Flags with Viper

		// Validate APIURL
		if !isValidURL(viper.GetString("api-url")) {
			fmt.Printf("api-url is not a valid URL: %s\n", viper.GetString("api-url"))
			os.Exit(1)
		}

		// Validate ZebedeeUrl
		if !isValidURL(viper.GetString("zebedee-url")) {
			fmt.Printf("zebedee-url is not a valid URL: %s\n", viper.GetString("zebedee-url"))
			os.Exit(1)
		}

		// Create FlagFields structure
		flagList := utilities.FlagFields{
			RobotsFilePath: viper.GetString("robots-file-path"),
			APIURL:         viper.GetString("api-url"),
			SitemapIndex:   viper.GetString("sitemap-index"),
			ScrollTimeout:  viper.GetString("scroll-timeout"),
			ScrollSize:     viper.GetInt("scroll-size"),
			SitemapPath:    viper.GetString("sitemap-file-path"),
			ZebedeeURL:     viper.GetString("zebedee-url"),
			FakeScroll:     viper.GetBool("fake-scroll"),
		}

		utilities.UpdateSitemap(cfg, &flagList)
	},
}

func init() {

	rootCmd.AddCommand(updateCmd)

	updateCmd.PersistentFlags().String("robots-file-path", "test_robots.txt", "path to robots file")
	updateCmd.PersistentFlags().String("sitemap-file-path", "test_sitemap", "path to sitemap file")
	updateCmd.PersistentFlags().String("api-url", "http://localhost", "elastic search api url")
	updateCmd.PersistentFlags().String("zebedee-url", "http://localhost:8082", "zebedee url")
	updateCmd.PersistentFlags().String("sitemap-index", "1", "OPENSEARCH_SITEMAP_INDEX")
	updateCmd.PersistentFlags().String("scroll-timeout", "2000", "OPENSEARCH_SCROLL_TIMEOUT")
	updateCmd.PersistentFlags().Int("scroll-size", 10, "OPENSEARCH_SCROLL_SIZE")
	updateCmd.PersistentFlags().Bool("fake-scroll", true, "enable fake scroll")
	updateCmd.PersistentFlags().Bool("update-sitemap", false, "update the sitemap")
}
