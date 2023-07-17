package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/ONSdigital/dp-sitemap/cmd/utilities"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate sitemap files using default parameters",
	Long:  `A tool to generate the sitemap`,
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.Get()
		if err != nil {
			fmt.Println("Error retrieving config" + err.Error())
			os.Exit(1)
		}

		viper.BindPFlags(cmd.Flags()) // Bind Flags with Viper

		// Validate ApiUrl
		if !isValidUrl(viper.GetString("api-url")) {
			fmt.Printf("api-url is not a valid URL: %s\n", viper.GetString("api-url"))
			os.Exit(1)
		}

		// Validate ZebedeeUrl
		if !isValidUrl(viper.GetString("zebedee-url")) {
			fmt.Printf("zebedee-url is not a valid URL: %s\n", viper.GetString("zebedee-url"))
			os.Exit(1)
		}

		// Create FlagFields structure
		flagList := utilities.FlagFields{
			RobotsFilePath:  viper.GetString("robots-file-path"),
			ApiUrl:          viper.GetString("api-url"),
			SitemapIndex:    viper.GetString("sitemap-index"),
			ScrollTimeout:   viper.GetString("scroll-timeout"),
			ScrollSize:      viper.GetInt("scroll-size"),
			SitemapPath:     viper.GetString("sitemap-file-path"),
			ZebedeeUrl:      viper.GetString("zebedee-url"),
			FakeScroll:      viper.GetBool("fake-scroll"),
			GenerateSitemap: viper.GetBool("generate-sitemap"),
			UpdateSitemap:   viper.GetBool("update-sitemap"),
		}

		utilities.GenerateSitemap(cfg, &flagList)
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
	generateCmd.PersistentFlags().Bool("fake-scroll", true, "enable fake scroll")
	generateCmd.PersistentFlags().Bool("generate-sitemap", false, "generate the sitemap")
	generateCmd.PersistentFlags().Bool("update-sitemap", false, "update the sitemap")

}

func isValidUrl(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}
