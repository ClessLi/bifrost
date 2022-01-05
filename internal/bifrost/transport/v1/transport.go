package v1

import (
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/web_server_config"
)

type TransportFactory interface {
	WebServerConfig() pbv1.WebServerConfigServer
}

var _ TransportFactory = &transportFactory{}

type transportFactory struct {
	handler HandlersFactory
	options *Options
}

func (t *transportFactory) WebServerConfig() pbv1.WebServerConfigServer {
	return web_server_config.NewWebServerConfigServer(t.handler, t.options)
}
