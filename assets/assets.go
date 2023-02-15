package assets

import (
	"context"
	"embed"
	"fmt"
	"io"
)

var (
	//go:embed robot/*
	static embed.FS
)

//go:generate moq -out mock/filesysteminterface.go -pkg mock . FileSystemInterface
type FileSystemInterface interface {
	Get(_ context.Context, path string) ([]byte, error)
}
type Embeddedfs struct{}

func NewFromEmbeddedFilesystem() Embeddedfs {
	return Embeddedfs{}
}

func (s Embeddedfs) Get(_ context.Context, path string) ([]byte, error) {
	file, err := static.Open(fmt.Sprintf("%s/%s", "robot", path))
	if err != nil {
		return nil, fmt.Errorf("cannot open file %w", err)
	}
	defer file.Close()
	b, err := io.ReadAll(file)
	return b, err
}
