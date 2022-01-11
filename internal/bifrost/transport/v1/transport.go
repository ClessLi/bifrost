package v1

import (
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
	svcv1 "github.com/ClessLi/bifrost/internal/bifrost/service/v1"
	storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/handler"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/options"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/web_server_config"
	"github.com/ClessLi/bifrost/internal/pkg/server"
	"time"
)

type Factory interface {
	WebServerConfig() pbv1.WebServerConfigServer
}

type transport struct {
	eps  epv1.EndpointsFactory
	opts *options.Options
}

func (t *transport) WebServerConfig() pbv1.WebServerConfigServer {
	return web_server_config.NewWebServerConfigServer(handler.NewWebServerConfigHandler(t.eps), t.opts)
}

func New(store storev1.StoreFactory, server *server.GenericGRPCServer) Factory {
	svc := svcv1.NewServiceFactory(store)
	eps := epv1.NewEndpoints(svc)
	opts := &options.Options{
		ChunkSize:          server.ChunkSize,
		RecvTimeoutMinutes: server.ReceiveTimeout / time.Minute,
	}
	return &transport{
		eps:  eps,
		opts: opts,
	}
}
