package web_server_bin_cmd

import (
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/handler"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/options"
)

type webServerBinCMDServer struct {
	handler handler.WebServerBinCMDHandlers
	options *options.Options
}

func NewWebServerBinCMDServer(
	handler handler.WebServerBinCMDHandlers,
	options *options.Options,
) *webServerBinCMDServer {
	return &webServerBinCMDServer{
		handler: handler,
		options: options,
	}
}
