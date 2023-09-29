package main

import (
	"context"

	"github.com/ONSdigital/log.go/v2/log"
)

func main() {
	ctx := context.Background()
	cmdErr := GetRootCommand().Execute()
	if cmdErr != nil {
		log.Error(ctx, "error initialising the CLI tool", cmdErr)
		return
	}
}
