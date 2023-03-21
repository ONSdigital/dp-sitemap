package sitemap

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//go:generate moq -out mock/s3uploader.go -pkg mock . S3Uploader

type S3Uploader interface {
	Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

type S3Saver struct {
	uploader S3Uploader
	bucket   string
	fileKey  map[config.Language]string
}

func NewS3Saver(uploader S3Uploader, bucket string, fileKey map[config.Language]string) *S3Saver {
	return &S3Saver{
		uploader: uploader,
		bucket:   bucket,
		fileKey:  fileKey,
	}
}

func (s *S3Saver) SaveFile(lang config.Language, body io.Reader) error {
	k := lang.String()
	_, err := s.uploader.Upload(&s3manager.UploadInput{
		Body:   body,
		Bucket: &s.bucket,
		Key:    &k,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to s3: %w", err)
	}
	log.Info(context.Background(), fmt.Sprintf("saved key [%s], bucket [%s]", k, s.bucket))
	return nil
}

func (s *S3Saver) UploadFiles(paths []string) error {
	for _, path := range paths {
		body, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open robots file: %w", err)
		}
		fileName := filepath.Base(path)
		_, err = s.uploader.Upload(&s3manager.UploadInput{
			Body:   body,
			Bucket: &s.bucket,
			Key:    &fileName,
		})
		if err != nil {
			return fmt.Errorf("failed to upload file to s3: %w", err)
		}
		log.Info(context.Background(), fmt.Sprintf("saved key [%s], bucket [%s]", fileName, s.bucket))
	}
	return nil
}
