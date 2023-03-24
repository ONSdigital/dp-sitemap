package sitemap

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//go:generate moq -out mock/s3client.go -pkg mock . S3Client

type S3Client interface {
	Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
	Get(key string) (io.ReadCloser, *int64, error)
	BucketName() string
}

type S3Store struct {
	client S3Client
}

func NewS3Store(client S3Client) *S3Store {
	return &S3Store{
		client: client,
	}
}

func (s *S3Store) SaveFile(name string, body io.Reader) error {
	key := filepath.Base(name) // trim the path
	bucket := s.client.BucketName()
	_, err := s.client.Upload(&s3manager.UploadInput{
		Body:   body,
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to s3: %w", err)
	}
	return nil
}

func (s *S3Store) GetFile(name string) (body io.ReadCloser, err error) {
	file, _, err := s.client.Get(name)
	if err != nil {
		return io.NopCloser(strings.NewReader("")), nil
	}
	return file, nil
}
