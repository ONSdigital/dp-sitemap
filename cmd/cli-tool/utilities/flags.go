package utilities

var CmdFlagFields *FlagFields

// Names of the flags used.
const (
	RobotsFilePathFlag       = "robots-file-path"
	RobotsFilePathReaderFlag = "robots-file-path-reader"
	ElasticSearchURLFlag     = "elasticsearch-url"
	ElasticSearchIndexFlag   = "elasticsearch-index"
	ScrollTimeoutFlag        = "scroll-timeout"
	ScrollSizeFlag           = "scroll-size"
	SitemapPathFlag          = "sitemap-file-path"
	SitemapPathReaderFlag    = "sitemap-file-path-reader"
	ZebedeeURLFlag           = "zebedee-url"
	FakeScrollFlag           = "fake-scroll"
)

// Config represents service configuration for dp-sitemap
type FlagFields struct {
	RobotsFilePath       string // path to the robots file that will be generated
	RobotsFilePathReader string // path to the robots file that we are reading from
	ElasticSearchURL     string // elastic search url
	ElasticSearchIndex   string // elastic search index name
	ScrollTimeout        string // elastic search scroll timeout
	ScrollSize           int    // elastic search scroll size
	SitemapPath          string // path to the sitemap file that will be generated
	SitemapPathReader    string // path to the sitemap file that we are reading from
	ZebedeeURL           string // zebedee url
	FakeScroll           bool   // toggle to use or not the fake scroll implementation that replicates elastic search
}
