package sitemap

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ONSdigital/dp-sitemap/config"
	"github.com/ONSdigital/log.go/v2/log"
)

type LocalSaver struct {
	fileNames map[config.Language]string
}

func NewLocalSaver(fileNames map[config.Language]string) *LocalSaver {
	return &LocalSaver{
		fileNames: fileNames,
	}
}

func (s *LocalSaver) SaveFile(lang config.Language, body io.Reader) error {
	file, err := os.OpenFile(s.fileNames[lang], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("failed to open a local file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, body)
	if err != nil {
		return fmt.Errorf("failed to copy to a local file: %w", err)
	}

	log.Info(context.Background(), fmt.Sprintf("saved file [%s], language [%s]", s.fileNames[lang], lang.String()))
	return nil
}
