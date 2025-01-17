package web_server_bin_cmd

import (
	"context"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
)

func (w *webServerBinCMDService) Exec(ctx context.Context, req *v1.ExecuteRequest) (*v1.ExecuteResponse, error) {
	return w.store.WebServerBinCMD().Exec(ctx, req)
}
