package fake

import (
	pbv1 "github.com/yongPhone/bifrost/api/protobuf-spec/bifrostpb/v1"
	log "github.com/yongPhone/bifrost/pkg/log/v1"
)

type webServerStatistics struct{}

func (w webServerStatistics) Get(servername *pbv1.ServerName, stream pbv1.WebServerStatistics_GetServer) error {
	log.Infof("get %s web server statistics", servername.Name)

	return nil
}
