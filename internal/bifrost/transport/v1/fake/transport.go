package fake

import (
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	txpv1 "github.com/ClessLi/bifrost/internal/bifrost/transport/v1"
)

type transport struct {
}

func (t transport) WebServerStatistics() pbv1.WebServerStatisticsServer {
	return webServerStatistics{}
}

func (t transport) WebServerConfig() pbv1.WebServerConfigServer {
	return webServerConfig{}
}

func New() txpv1.Factory {
	return transport{}
}
