package fake

import (
	"context"
	"io"

	"github.com/marmotedu/errors"

	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	log "github.com/ClessLi/bifrost/pkg/log/v1"
)

type webServerConfig struct{}

func (w webServerConfig) GetServerNames(ctx context.Context, null *pbv1.Null) (*pbv1.ServerNames, error) {
	log.Info("get web server names")

	return &pbv1.ServerNames{Names: []*pbv1.ServerName{{Name: "test1"}, {Name: "test2"}}}, nil
}

func (w webServerConfig) Get(servername *pbv1.ServerName, stream pbv1.WebServerConfig_GetServer) error {
	log.Infof("get web server config %s", servername.Name)

	return nil
}

func (w webServerConfig) Update(stream pbv1.WebServerConfig_UpdateServer) error {
	conf, err := stream.Recv()
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	log.Infof("update web server config %s", conf.GetServerName())

	return nil
}

var _ pbv1.WebServerConfigServer = webServerConfig{}
