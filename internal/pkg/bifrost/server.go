package bifrost

import (
	"errors"
	"fmt"
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/authentication"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/config"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/endpoint"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/logging"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/transport"
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

func ServerRun() error {
	if !isInit {
		return errors.New("service related configuration not initialized")
	}

	err := bifrostServiceStart()
	if err != nil {
		panic(err)
	}
	defer func() {
		err = bifrostServiceStop()
		if err != nil {
			utils.Logger.Error(err.Error())
		}
	}()
	utils.Logger.Debug("Listening system call signal")
	go ListenSignal(signalChan)
	utils.Logger.Debug("Listened system call signal")

	ctx := context.Background()

	// init service
	svc := newService()
	// init auth svc
	svc = authentication.AuthenticationMiddleware(BifrostConf.ServiceConfig.AuthServerAddr)(svc)

	// init kit logger
	svc = logging.LoggingMiddleware(config.KitLogger)(svc)

	// init kit endpoint
	endpoints := endpoint.NewBifrostEndpoints(svc)

	transport.ChunkSize = BifrostConf.ServiceConfig.ChunckSize
	handlers := transport.NewGRPCHandlers(ctx, endpoints)
	healthCheckHandler := transport.NewHealthCheckHandler(ctx, endpoints)

	lis, lisErr := net.Listen("tcp", fmt.Sprintf(":%d", BifrostConf.ServiceConfig.Port))
	if lisErr != nil {
		return lisErr
	}
	defer lis.Close()

	gRPCServer := grpc.NewServer(grpc.MaxSendMsgSize(transport.ChunkSize))
	bifrostpb.RegisterViewServiceServer(gRPCServer, handlers)
	bifrostpb.RegisterUpdateServiceServer(gRPCServer, handlers)
	bifrostpb.RegisterWatchServiceServer(gRPCServer, handlers)
	grpc_health_v1.RegisterHealthServer(gRPCServer, healthCheckHandler)
	utils.Logger.Info(logoStr)
	svrErrChan := make(chan error, 0)
	go func() {
		utils.Logger.NoticeF("bifrost service is running on %s", lis.Addr())
		svrErr := gRPCServer.Serve(lis)
		utils.Logger.NoticeF("bifrost service is stopped")
		svrErrChan <- svrErr
	}()

	go registerToRA(svrErrChan)
	defer deregisterToRA()

	var stopErr error
	select {
	case s := <-signalChan:
		if s == 9 {
			utils.Logger.Debug("bifrost service is stopping...")
			gRPCServer.Stop()
		}
		utils.Logger.Debug("stop signal error")
	case stopErr = <-svrErrChan:
		gRPCServer.Stop()
		break
	}
	return stopErr
}
