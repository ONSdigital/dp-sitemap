package sitemap

import "context"

type Scroll interface {
	StartScroll(ctx context.Context, result interface{}) error
	GetScroll(ctx context.Context, id string, result interface{}) error
}
