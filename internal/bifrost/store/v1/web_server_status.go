package v1

import (
	"context"

	v1 "github.com/yongPhone/bifrost/api/bifrost/v1"
)

type WebServerStatusStore interface {
	Get(ctx context.Context) (*v1.Metrics, error)
}
