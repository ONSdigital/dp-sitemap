package utilities

import (
	"testing"

	"github.com/ONSdigital/dp-sitemap/config"
	. "github.com/smartystreets/goconvey/convey"
)

// // Mock dependencies

// type MockRoundTripper struct {
// 	signed bool
// }

// func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
// 	return &http.Response{}, nil
// }

// func MockNewAWSSignerRoundTripper(...interface{}) (http.RoundTripper, error) {
// 	return &MockRoundTripper{signed: true}, nil
// }

// func MockNewClient(config es710.Config) (*es710.Client, error) {
// 	if _, ok := config.Transport.(*MockRoundTripper); !ok {
// 		return nil, errors.New("transport is not signed")
// 	}
// 	return &es710.Client{}, nil
// }

func TestGenerateSitemap(t *testing.T) {
	Convey("Given configuration and commandline flags", t, func() {
		cfg := &config.Config{
			OpenSearchConfig: config.OpenSearchConfig{
				Signer: true,
			},
		}
		cmd := &FlagFields{}

		Convey("When GenerateSitemap is called", func() {

			GenerateSitemap(cfg, cmd)

			Convey("It should configure transport with AWS signer", func() {
				// Here, you can check the aspects of the function that are observable from the outside,
				// such as effects on global state or function outputs.
				// In this mock, we haven't actually given a way to observe the effect,
				// so this is just a placeholder.
			})

			Convey("Elasticsearch client should be initialized with the correct transport", func() {
				// Similar to the above, this would require you to be able to observe
				// the configuration of the Elasticsearch client.
			})
		})
	})
}
