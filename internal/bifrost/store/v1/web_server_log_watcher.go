package v1

import (
	"context"

	v1 "github.com/yongPhone/bifrost/api/bifrost/v1"
)

type WebServerLogWatcher interface {
	Watch(ctx context.Context, request *v1.WebServerLogWatchRequest) (*v1.WebServerLog, error)
}
