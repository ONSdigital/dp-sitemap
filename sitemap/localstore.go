package sitemap

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ONSdigital/log.go/v2/log"
)

type LocalStore struct{}

func (s *LocalStore) SaveFile(name string, body io.Reader) error {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("failed to open a local file: %w", err)
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			log.Error(context.Background(), "failed to close a local file", closeErr)
		}
	}()

	_, err = io.Copy(file, body)
	if err != nil {
		return fmt.Errorf("failed to copy to a local file: %w", err)
	}

	return nil
}

func (s *LocalStore) GetFile(name string) (body io.ReadCloser, err error) {
	file, err := os.Open(name)
	if err != nil {
		return io.NopCloser(strings.NewReader("")), nil
	}
	return file, nil
}

func (s *LocalStore) CopyFile(src io.Reader, dest io.Writer) error {
	_, err := io.Copy(dest, src)
	if err != nil {
		return fmt.Errorf("failed to copy file : %w", err)
	}
	return nil
}

func (s *LocalStore) CreateFile(name string) (io.ReadWriteCloser, error) {
	file, err := os.Create(name)
	if err != nil {
		return nil, fmt.Errorf("failed to create file : %w", err)
	}
	return file, nil
}

func (s *LocalStore) DeleteFile(name string) error {
	err := os.Remove(name)
	if err != nil {
		return fmt.Errorf("failed to delete file : %w", err)
	}
	return nil
}
