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

func GetRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "dp-sitemap",
		Short: "CLI tool to generate and update sitemaps ",
	}
	rootCmd.AddCommand(setupGenerateCmd())
	rootCmd.AddCommand(setupUpdateCmd())
	return rootCmd
}

func setupGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate the sitemap",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Get()
			if err != nil {
				fmt.Println("Error retrieving config" + err.Error())
				os.Exit(1)
			}

			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				fmt.Println("Error binding flags: " + err.Error())
				os.Exit(1)
			}

			if !isValidURL(viper.GetString("api-url")) {
				fmt.Printf("api-url is not a valid URL: %s\n", viper.GetString("api-url"))
				os.Exit(1)
			}

			if !isValidURL(viper.GetString("zebedee-url")) {
				fmt.Printf("zebedee-url is not a valid URL: %s\n", viper.GetString("zebedee-url"))
				os.Exit(1)
			}

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

			utilities.GenerateSitemap(cfg, &flagList)
			return nil
		},
	}

	cmd.Flags().String("robots-file-path", "test_robots.txt", "path to robots file")
	cmd.Flags().String("sitemap-file-path", "test_sitemap", "path to sitemap file")
	cmd.Flags().String("api-url", "http://localhost", "elastic search api url")
	cmd.Flags().String("zebedee-url", "http://localhost:8082", "zebedee url")
	cmd.Flags().String("sitemap-index", "1", "OPENSEARCH_SITEMAP_INDEX")
	cmd.Flags().String("scroll-timeout", "2000", "OPENSEARCH_SCROLL_TIMEOUT")
	cmd.Flags().Int("scroll-size", 10, "OPENSEARCH_SCROLL_SIZE")
	cmd.Flags().Bool("fake-scroll", true, "enable fake scroll")
	return cmd
}

func setupUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the sitemap",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Get()
			if err != nil {
				fmt.Println("Error retrieving config" + err.Error())
				os.Exit(1)
			}

			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				fmt.Println("Error binding flags: " + err.Error())
				os.Exit(1)
			}

			if !isValidURL(viper.GetString("api-url")) {
				fmt.Printf("api-url is not a valid URL: %s\n", viper.GetString("api-url"))
				os.Exit(1)
			}

			if !isValidURL(viper.GetString("zebedee-url")) {
				fmt.Printf("zebedee-url is not a valid URL: %s\n", viper.GetString("zebedee-url"))
				os.Exit(1)
			}

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
			return nil
		},
	}

	cmd.Flags().String("robots-file-path", "test_robots.txt", "path to robots file")
	cmd.Flags().String("sitemap-file-path", "test_sitemap", "path to sitemap file")
	cmd.Flags().String("api-url", "http://localhost", "elastic search api url")
	cmd.Flags().String("zebedee-url", "http://localhost:8082", "zebedee url")
	cmd.Flags().String("sitemap-index", "1", "OPENSEARCH_SITEMAP_INDEX")
	cmd.Flags().String("scroll-timeout", "2000", "OPENSEARCH_SCROLL_TIMEOUT")
	cmd.Flags().Int("scroll-size", 10, "OPENSEARCH_SCROLL_SIZE")
	cmd.Flags().Bool("fake-scroll", true, "enable fake scroll")

	return cmd
}

func isValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}
