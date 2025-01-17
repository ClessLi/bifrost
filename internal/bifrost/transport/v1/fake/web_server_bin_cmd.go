package fake

import (
	"context"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
)

type webServerBinCMD struct{}

func (w webServerBinCMD) Exec(ctx context.Context, request *pbv1.ExecuteRequest) (*pbv1.ExecuteResponse, error) {
	logV1.Infof("web server binary command excuting...")
	return &pbv1.ExecuteResponse{
		Successful: true,
		Msg:        []byte("success\n"),
	}, nil
}
