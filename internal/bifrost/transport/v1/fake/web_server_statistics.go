package fake

import (
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"

	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
)

type webServerStatistics struct{}

func (w webServerStatistics) Get(servername *pbv1.ServerName, stream pbv1.WebServerStatistics_GetServer) error {
	logV1.Infof("get %s web server statistics", servername.Name)

	return nil
}
