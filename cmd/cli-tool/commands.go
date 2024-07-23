package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/ONSdigital/dp-sitemap/cmd/cli-tool/utilities"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "dp-sitemap",
		Short: "CLI tool to generate and update sitemaps ",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
	}
	rootCmd.PersistentFlags().String(utilities.RobotsFilePathFlag, "test_robots.txt", "path to robots file that will be generated")
	rootCmd.PersistentFlags().String(utilities.RobotsFilePathReaderFlag, "", "path to external robot file (optional)")
	rootCmd.PersistentFlags().String(utilities.SitemapPathFlag, "test_sitemap", "path to sitemap file")
	rootCmd.PersistentFlags().String(utilities.SitemapPathReaderFlag, "./sitemap/static/", "path to external sitemap files (optional)")
	rootCmd.PersistentFlags().String(utilities.ElasticSearchURLFlag, "http://localhost", "elastic search url")
	rootCmd.PersistentFlags().String(utilities.ZebedeeURLFlag, "http://localhost:8082", "zebedee url")
	rootCmd.PersistentFlags().String(utilities.ElasticSearchIndexFlag, "1", "OPENSEARCH_SITEMAP_INDEX")
	rootCmd.PersistentFlags().String(utilities.ScrollTimeoutFlag, "2000", "OPENSEARCH_SCROLL_TIMEOUT")
	rootCmd.PersistentFlags().Int(utilities.ScrollSizeFlag, 10, "OPENSEARCH_SCROLL_SIZE")
	rootCmd.PersistentFlags().Bool(utilities.FakeScrollFlag, true, "enable fake scroll")
	rootCmd.AddCommand(setupGenerateCmd())
	rootCmd.AddCommand(setupUpdateCmd())
	rootCmd.AddCommand(setupLoadStaticSitemapCmd())
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

			if !isValidURL(viper.GetString(utilities.ElasticSearchURLFlag)) {
				fmt.Printf("%s is not a valid URL: %s\n", utilities.ElasticSearchURLFlag, viper.GetString(utilities.ElasticSearchURLFlag))
				os.Exit(1)
			}

			if !isValidURL(viper.GetString(utilities.ZebedeeURLFlag)) {
				fmt.Printf("%s is not a valid URL: %s\n", utilities.ZebedeeURLFlag, viper.GetString(utilities.ZebedeeURLFlag))
				os.Exit(1)
			}

			flagList := utilities.FlagFields{
				RobotsFilePath:       viper.GetString(utilities.RobotsFilePathFlag),
				RobotsFilePathReader: viper.GetString(utilities.RobotsFilePathReaderFlag),
				ElasticSearchURL:     viper.GetString(utilities.ElasticSearchURLFlag),
				ElasticSearchIndex:   viper.GetString(utilities.ElasticSearchIndexFlag),
				ScrollTimeout:        viper.GetString(utilities.ScrollTimeoutFlag),
				ScrollSize:           viper.GetInt(utilities.ScrollSizeFlag),
				SitemapPath:          viper.GetString(utilities.SitemapPathFlag),
				SitemapPathReader:    viper.GetString(utilities.SitemapPathReaderFlag),
				ZebedeeURL:           viper.GetString(utilities.ZebedeeURLFlag),
				FakeScroll:           viper.GetBool(utilities.FakeScrollFlag),
			}
			utilities.CmdFlagFields = &flagList
			utilities.GenerateSitemap(cfg, &flagList)
			utilities.GenerateRobotFile(cfg, &flagList)
			return nil
		},
	}

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

			if !isValidURL(viper.GetString(utilities.ElasticSearchURLFlag)) {
				fmt.Printf("%s is not a valid URL: %s\n", utilities.ElasticSearchURLFlag, viper.GetString(utilities.ElasticSearchURLFlag))
				os.Exit(1)
			}

			if !isValidURL(viper.GetString(utilities.ZebedeeURLFlag)) {
				fmt.Printf("%s is not a valid URL: %s\n", utilities.ZebedeeURLFlag, viper.GetString(utilities.ZebedeeURLFlag))
				os.Exit(1)
			}

			flagList := utilities.FlagFields{
				RobotsFilePath:       viper.GetString(utilities.RobotsFilePathFlag),
				RobotsFilePathReader: viper.GetString(utilities.RobotsFilePathReaderFlag),
				ElasticSearchURL:     viper.GetString(utilities.ElasticSearchURLFlag),
				ElasticSearchIndex:   viper.GetString(utilities.ElasticSearchIndexFlag),
				ScrollTimeout:        viper.GetString(utilities.ScrollTimeoutFlag),
				ScrollSize:           viper.GetInt(utilities.ScrollSizeFlag),
				SitemapPath:          viper.GetString(utilities.SitemapPathFlag),
				SitemapPathReader:    viper.GetString(utilities.SitemapPathReaderFlag),
				ZebedeeURL:           viper.GetString(utilities.ZebedeeURLFlag),
				FakeScroll:           viper.GetBool(utilities.FakeScrollFlag),
			}
			utilities.CmdFlagFields = &flagList
			err = utilities.UpdateSitemap(cfg, &flagList)
			if err != nil {
				return err
			}
			utilities.GenerateRobotFile(cfg, &flagList)
			return nil
		},
	}
	return cmd
}

func setupLoadStaticSitemapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load",
		Short: "Load the static sitemap",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Get()
			if err != nil {
				fmt.Println("Error retrieving config" + err.Error())
				os.Exit(1)
			}

			flagList := utilities.FlagFields{
				RobotsFilePathReader: viper.GetString(utilities.RobotsFilePathReaderFlag),
				SitemapPath:          viper.GetString(utilities.SitemapPathFlag),
				SitemapPathReader:    viper.GetString(utilities.SitemapPathReaderFlag),
			}
			utilities.CmdFlagFields = &flagList
			utilities.LoadStaticSitemap(cfg, &flagList)
			return nil
		},
	}
	return cmd
}

func isValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}
