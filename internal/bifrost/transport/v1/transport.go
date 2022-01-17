package v1

import (
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/handler"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/options"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/web_server_config"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/web_server_statistics"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/web_server_status"
)

type Factory interface {
	WebServerConfig() pbv1.WebServerConfigServer
	WebServerStatistics() pbv1.WebServerStatisticsServer
	WebServerStatus() pbv1.WebServerStatusServer
}

type transport struct {
	handlers handler.HandlersFactory
	opts     *options.Options
}

func (t *transport) WebServerConfig() pbv1.WebServerConfigServer {
	return web_server_config.NewWebServerConfigServer(t.handlers.WebServerConfig(), t.opts)
}

func (t *transport) WebServerStatistics() pbv1.WebServerStatisticsServer {
	return web_server_statistics.NewWebServerStatisticsServer(t.handlers.WebServerStatistics(), t.opts)
}

func (t *transport) WebServerStatus() pbv1.WebServerStatusServer {
	return web_server_status.NewWebServerStatusServer(t.handlers.WebServerStatus(), t.opts)
}

func New(handlers handler.HandlersFactory, opts *options.Options) Factory {
	return &transport{
		handlers: handlers,
		opts:     opts,
	}
}
