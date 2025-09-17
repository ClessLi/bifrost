package nginx

import (
	"context"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"
)

type webServerBinCMDStore struct {
	store *webServerStore
}

func (w webServerBinCMDStore) Exec(ctx context.Context, request *v1.ExecuteRequest) (*v1.ExecuteResponse, error) {
	cmd, err := w.store.configsManger.GenServerBinCMD(request.ServerName, request.Args...)
	if err != nil {
		return nil, err
	}

	if output, err := cmd.CombinedOutput(); err != nil {
		return &v1.ExecuteResponse{
			Successful:     false,
			StandardOutput: []byte{},
			StandardError:  output,
		}, nil
	} else {
		return &v1.ExecuteResponse{
			Successful:     true,
			StandardOutput: output,
			StandardError:  []byte{},
		}, nil
	}
}

var _ storev1.WebServerBinCMDStore = &webServerBinCMDStore{}

func newNginxBinCMDStore(store *webServerStore) storev1.WebServerBinCMDStore {
	return &webServerBinCMDStore{store: store}
}
