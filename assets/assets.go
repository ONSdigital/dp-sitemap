package assets

import (
	"context"
	"embed"
	"fmt"
	"io"
)

var (
	//go:embed robot/*
	//go:embed sitemap/*
	static embed.FS
)

type EmbeddedFile string

const (
	Robots  EmbeddedFile = "robots"
	Sitemap EmbeddedFile = "sitemap"
)

func (l EmbeddedFile) String() string {
	switch l {
	case Robots:
		return "robot"
	default:
		return "sitemap"
	}
}

//go:generate moq -out mock/filesysteminterface.go -pkg mock . FileSystemInterface
type FileSystemInterface interface {
	Get(
		context.Context, EmbeddedFile, string) ([]byte, error)
}
type Embeddedfs struct{}

func NewFromEmbeddedFilesystem() Embeddedfs {
	return Embeddedfs{}
}

func (s Embeddedfs) Get(_ context.Context, embeddedFile EmbeddedFile, path string) ([]byte, error) {
	file, err := static.Open(embeddedFile.String() + "/" + path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %w", err)
	}
	defer file.Close()
	b, err := io.ReadAll(file)
	return b, err
}
