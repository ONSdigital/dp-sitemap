package sitemap

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//go:generate moq -out mock/s3uploader.go -pkg mock . S3Uploader

type S3Uploader interface {
	Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

type S3Saver struct {
	uploader S3Uploader
	bucket   string
	fileKey  string
}

func NewS3Saver(uploader S3Uploader, bucket, fileKey string) *S3Saver {
	return &S3Saver{
		uploader: uploader,
		bucket:   bucket,
		fileKey:  fileKey,
	}
}

func (s *S3Saver) SaveFile(body io.Reader) error {
	_, err := s.uploader.Upload(&s3manager.UploadInput{
		Body:   body,
		Bucket: &s.bucket,
		Key:    &s.fileKey,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to s3: %w", err)
	}
	return nil
}
