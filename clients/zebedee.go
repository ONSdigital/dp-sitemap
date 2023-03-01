package clients

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

//go:generate moq -out mock/zebedee.go -pkg mock . ZebedeeClient

// ZebedeeClient defines the zebedee client
type ZebedeeClient interface {
	Checker(context.Context, *healthcheck.CheckState) error
	GetFileSize(ctx context.Context, userAccessToken, collectionID, lang, uri string) (zebedee.FileSize, error)
}
