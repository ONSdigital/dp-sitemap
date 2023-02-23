package sitemap

import (
	"fmt"
	"io"
	"os"
)

type LocalSaver struct {
	fileName string
}

func NewLocalSaver(fileName string) *LocalSaver {
	return &LocalSaver{
		fileName: fileName,
	}
}

func (s *LocalSaver) SaveFile(body io.Reader) error {
	file, err := os.OpenFile(s.fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("failed to open a local file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, body)
	if err != nil {
		return fmt.Errorf("failed to copy to a local file: %w", err)
	}

	return nil
}
