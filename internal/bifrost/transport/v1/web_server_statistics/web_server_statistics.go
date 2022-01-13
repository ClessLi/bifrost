package web_server_statistics

import (
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/handler"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/options"
)

var _ pbv1.WebServerStatisticsServer = &webServerStatisticsServer{}

type webServerStatisticsServer struct {
	handler handler.WebServerStatisticsHandlers
	options *options.Options
}

func NewWebServerStatisticsServer(handler handler.WebServerStatisticsHandlers, options *options.Options) pbv1.WebServerStatisticsServer {
	return &webServerStatisticsServer{
		handler: handler,
		options: options,
	}
}
