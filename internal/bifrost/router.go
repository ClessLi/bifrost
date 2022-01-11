package bifrost

import (
	storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"
	txpv1 "github.com/ClessLi/bifrost/internal/bifrost/transport/v1"
	genericgrpcserver "github.com/ClessLi/bifrost/internal/pkg/server"
)

func initRouter(server *genericgrpcserver.GenericGRPCServer) {
	initMiddleware(server)
	initController(server)
}

func initMiddleware(server *genericgrpcserver.GenericGRPCServer) {

}

func initController(server *genericgrpcserver.GenericGRPCServer) {
	// v1 handlers
	storeIns := storev1.Client()
	bifrostTransport := txpv1.New(storeIns, server)
	{
		// register bifrost services
		registers := txpv1.NewBifrostServiceRegister(bifrostTransport)
		server.RegisterServices(registers.Generate())
	}
}
