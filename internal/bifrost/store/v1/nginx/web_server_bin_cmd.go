package nginx

import (
	"bytes"
	"context"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/marmotedu/errors"
	"os/exec"
)

type webServerBinCMDStore struct {
	serversBinCMD map[string]func(arg ...string) *exec.Cmd
}

func (w webServerBinCMDStore) Exec(ctx context.Context, request *v1.ExecuteRequest) (*v1.ExecuteResponse, error) {
	if f, has := w.serversBinCMD[request.ServerName]; has {
		cmd := f(request.Args...)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		var isSuccessful = cmd.Run() == nil

		return &v1.ExecuteResponse{
			Successful:     isSuccessful,
			StandardOutput: stdout.Bytes(),
			StandardError:  stderr.Bytes(),
		}, nil
	}

	return nil, errors.WithCode(code.ErrWebServerNotFound, "nginx server '%s' not found", request.ServerName)
}

var _ storev1.WebServerBinCMDStore = &webServerBinCMDStore{}

func newNginxBinCMDStore(store *webServerStore) storev1.WebServerBinCMDStore {
	return &webServerBinCMDStore{serversBinCMD: store.configsManger.GetServersBinCMD()}
}
