package robotseo

type SeoRobotModel struct {
	AllowList []string `json:"AllowList"`
	DenyList  []string `json:"DenyList"`
}
