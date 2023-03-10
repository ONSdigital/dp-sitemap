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

func (a *DefaultAdder) Add(oldSitemap io.Reader, url URL) (fileName string, err error) {
	// create a temporary file
	file, err := os.CreateTemp("", "sitemap-incr")
	if err != nil {
		return "", fmt.Errorf("failed to create incremental sitemap file: %w", err)
	}
	fileName = file.Name()
	log.Info(context.Background(), "created incremental sitemap file "+fileName)
	defer func() {
		file.Close()
		// clean up the temporary file if we're returning with an error
		if err != nil {
			removeErr := os.Remove(fileName)
			if removeErr != nil {
				log.Error(context.Background(), "failed to remove incremental sitemap file", err)
				return
			}
			log.Info(context.Background(), "removed incremental sitemap file "+fileName)
		}
	}()

	// get the old sitemap
	sitemap := Urlset{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}
	decoder := xml.NewDecoder(oldSitemap)
	err = decoder.Decode(&sitemap)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return fileName, fmt.Errorf("failed to decode old sitemap: %w", err)
		}
	}

	// add new URL
	sitemap.URL = append(sitemap.URL, url)

	// output result into the file
	_, err = file.WriteString(xml.Header)
	if err != nil {
		return fileName, fmt.Errorf("failed to write xml doctype: %w", err)
	}
	enc := xml.NewEncoder(file)
	enc.Indent("", "  ")
	err = enc.Encode(sitemap)
	if err != nil {
		return fileName, fmt.Errorf("failed to encode sitemap: %w", err)
	}

	return fileName, nil
}
