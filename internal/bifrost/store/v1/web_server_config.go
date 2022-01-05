package v1

import (
	"context"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
)

type WebServerConfigStore interface {
	Get(ctx context.Context, name *v1.ServerName) (*v1.WebServerConfig, error)
	Update(ctx context.Context, config *v1.WebServerConfig) error
}
