package faker

import (
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	log "github.com/ClessLi/bifrost/pkg/log/v1"
	"io"
)

type webServerConfig struct {
}

func (w webServerConfig) Get(name *pbv1.ServerName, stream pbv1.WebServerConfig_GetServer) error {
	log.Infof("get web server config %s", name)
	return nil
}

func (w webServerConfig) Update(stream pbv1.WebServerConfig_UpdateServer) error {
	conf, err := stream.Recv()
	if err != nil && err != io.EOF {
		return err
	}

	log.Infof("update web server config %s", conf.GetServerName())
	return nil
}

var _ pbv1.WebServerConfigServer = webServerConfig{}
