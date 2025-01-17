package v1

import (
	"context"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
)

type WebServerLogWatcherStore interface {
	Watch(ctx context.Context, request *v1.WebServerLogWatchRequest) (*v1.WebServerLog, error)
}
