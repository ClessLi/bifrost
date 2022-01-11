package web_server_config

import (
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/handler"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/options"
)

var _ pbv1.WebServerConfigServer = &webServerConfigServer{}

type webServerConfigServer struct {
	handler handler.WebServerConfigHandler
	options *options.Options
}

func NewWebServerConfigServer(handler handler.WebServerConfigHandler, options *options.Options) pbv1.WebServerConfigServer {
	return &webServerConfigServer{
		handler: handler,
		options: options,
	}
}
