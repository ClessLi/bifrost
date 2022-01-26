package v1

import (
	"context"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
)

type WebServerStatisticsStore interface {
	Get(ctx context.Context, servername *v1.ServerName) (*v1.Statistics, error)
}
