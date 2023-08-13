package global

var CmdFlagFields *FlagFields

// Config represents service configuration for dp-sitemap
type FlagFields struct {
	RobotsFilePath       string // path to the robots file that will be generated
	RobotsFilePathReader string // path to the robots file that we are reading from
	APIURL               string // elastic search api url
	SitemapIndex         string // elastic search sitemap index
	ScrollTimeout        string // elastic search scroll timeout
	ScrollSize           int    // elastic search scroll size
	SitemapPath          string // path to the sitemap file
	ZebedeeURL           string // zebedee url
	FakeScroll           bool   // toggle to use or not the fake scroll implementation that replicates elastic search
}
