package event

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-sitemap/clients"
	mock2 "github.com/ONSdigital/dp-sitemap/clients/mock"
	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/dp-sitemap/sitemap"
	"github.com/ONSdigital/dp-sitemap/sitemap/mock"
	"io"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHandle(t *testing.T) {
	Convey("When having a content published event", t, func() {
		store := &mock.FileStoreMock{}
		fetcher := &mock.FetcherMock{}
		zebedeeClient := &mock2.ZebedeeClientMock{}
		handler := NewContentPublishedHandler(store, fetcher)
		cfg, _ := config.Get()
		content := &ContentPublished{
			URI:          "economy/environmentalaccounts/articles/testarticle3",
			DataType:     "theDateType",
			CollectionID: "theCollectionId",
			JobID:        "theJobId",
			SearchIndex:  "theSearchIndex",
			TraceID:      "theTraceId",
		}

		fetcher.URLVersionsFunc = func(ctx context.Context, path, lastmod string) (*sitemap.URL, *sitemap.URL) {
			return &sitemap.URL{Loc: path, Lastmod: "2006-01-02T15:04:05Z"}, &sitemap.URL{Loc: path, Lastmod: "2006-01-02T15:04:05Z"}
		}

		fetcher.URLVersionFunc = func(ctx context.Context, path string, lastmod string, lang string) *sitemap.URL {
			return &sitemap.URL{Loc: path, Lastmod: lastmod}
		}

		fetcher.GetZebedeeClientFunc = func() clients.ZebedeeClient {
			return zebedeeClient
		}

		zebedeeClient.GetPageDescriptionFunc = func(ctx context.Context, userAccessToken string, collectionID string, lang string, uri string) (zebedee.PageDescription, error) {
			return zebedee.PageDescription{URI: uri, Description: zebedee.Description{
				ReleaseDate: "2006-01-02T15:04:05Z",
			}}, nil
		}

		store.GetFileFunc = func(name string) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader("")), nil
		}

		store.CreateFileFunc = func(name string) (io.ReadWriteCloser, error) {
			file, _ := os.CreateTemp("", name)
			return file, nil
		}

		store.CopyFileFunc = func(src io.Reader, dest io.Writer) error {
			return nil
		}

		err := handler.Handle(context.Background(), cfg, content)
		Convey("There should be no error", func() {
			So(err, ShouldBeNil)
		})
	})
}
