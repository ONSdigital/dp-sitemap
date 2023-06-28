package sitemap

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ONSdigital/log.go/v2/log"
)

type DefaultAdder struct{}

func (a *DefaultAdder) Add(oldSitemap io.Reader, url *URL) (fileName string, size int, err error) {
	// create a temporary file
	file, err := os.CreateTemp("", "sitemap-incr")
	if err != nil {
		return "", 0, fmt.Errorf("failed to create publishing sitemap file: %w", err)
	}
	fileName = file.Name()
	log.Info(context.Background(), "created publishing sitemap file "+fileName)
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			log.Error(context.Background(), "failed to close publishing sitemap file", closeErr)
		}
		// clean up the temporary file if we're returning with an error
		if err != nil {
			removeErr := os.Remove(fileName)
			if removeErr != nil {
				log.Error(context.Background(), "failed to remove publishing sitemap file", removeErr)
				return
			}
			log.Info(context.Background(), "removed publishing sitemap file "+fileName)
		}
	}()

	// get the old sitemap
	sitemapReader := UrlsetReader{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Xhtml: "http://www.w3.org/1999/xhtml",
	}
	decoder := xml.NewDecoder(oldSitemap)
	err = decoder.Decode(&sitemapReader)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return fileName, 0, fmt.Errorf("failed to decode old sitemap: %w", err)
		}
	}

	// move old sitemap urls to new sitemap
	sitemap := Urlset{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Xhtml: "http://www.w3.org/1999/xhtml",
	}

	for _, oldURL := range sitemapReader.URL {
		var newURL URL
		newURL.Loc = oldURL.Loc
		newURL.Lastmod = oldURL.Lastmod
		newURL.Alternate = &AlternateURL{}
		if oldURL.Alternate != nil {
			newURL.Alternate.Rel = oldURL.Alternate.Rel
			newURL.Alternate.Link = oldURL.Alternate.Link
			newURL.Alternate.Lang = oldURL.Alternate.Lang
		}
		sitemap.URL = append(sitemap.URL, newURL)
	}

	// add new URL
	if url != nil {
		sitemap.URL = append(sitemap.URL, *url)
	}

	// output result into the file
	_, err = file.WriteString(xml.Header)
	if err != nil {
		return fileName, 0, fmt.Errorf("failed to write xml doctype: %w", err)
	}
	enc := xml.NewEncoder(file)
	enc.Indent("", "  ")
	err = enc.Encode(sitemap)
	if err != nil {
		return fileName, 0, fmt.Errorf("failed to encode sitemap: %w", err)
	}

	return fileName, len(sitemap.URL), nil
}
