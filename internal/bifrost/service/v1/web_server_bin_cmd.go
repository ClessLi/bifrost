package v1

import (
	"context"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
)

type WebServerBinCMDService interface {
	Exec(ctx context.Context, request *v1.ExecuteRequest) (*v1.ExecuteResponse, error)
}
