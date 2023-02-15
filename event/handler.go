package event

import (
	"context"

	"github.com/ONSdigital/dp-sitemap/config"
)

// TODO: remove hello called example handler
// HelloCalledHandler ...
type HelloCalledHandler struct {
}

// Handle takes a single event.
func (h *HelloCalledHandler) Handle(ctx context.Context, cfg *config.Config, event *HelloCalled) (err error) {
	return nil
}
