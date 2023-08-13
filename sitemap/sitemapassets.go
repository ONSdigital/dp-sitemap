package sitemap

import "embed"

//go:embed static/sitemap_en.json
//go:embed static/sitemap_cy.json

var folder embed.FS

func GetStaticSitemap(filename string) ([]byte, error) {
	file, err := folder.ReadFile("static/" + filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}
