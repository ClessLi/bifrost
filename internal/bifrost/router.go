package bifrost

import (
	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
	svcv1 "github.com/ClessLi/bifrost/internal/bifrost/service/v1"
	storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"
	txpv1 "github.com/ClessLi/bifrost/internal/bifrost/transport/v1"
	handlerv1 "github.com/ClessLi/bifrost/internal/bifrost/transport/v1/handler"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/options"
	genericgrpcserver "github.com/ClessLi/bifrost/internal/pkg/server"
	"time"
)

func initRouter(server *genericgrpcserver.GenericGRPCServer) {
	initMiddleware(server)
	initController(server)
}

func initMiddleware(server *genericgrpcserver.GenericGRPCServer) {

}

func initController(server *genericgrpcserver.GenericGRPCServer) {
	// v1 transport
	storeIns := storev1.Client()
	svc := svcv1.NewServiceFactory(storeIns)
	eps := epv1.NewEndpoints(svc)
	hs := handlerv1.NewHandlersFactory(eps)
	opts := &options.Options{
		ChunkSize:          server.ChunkSize,
		RecvTimeoutMinutes: server.ReceiveTimeout / time.Minute,
	}

	txp := txpv1.New(hs, opts)
	{
		// register bifrost services
		registers := txpv1.NewBifrostServiceRegister(txp)
		server.RegisterServices(registers.Generate())
	}
}
