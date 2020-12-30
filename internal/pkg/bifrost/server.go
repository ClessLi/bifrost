package bifrost

import (
	"errors"
	"fmt"
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/config"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/endpoint"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/logging"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service"
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

	utils.Logger.Debug("Listening system call signal")
	go ListenSignal(signalChan)
	utils.Logger.Debug("Listened system call signal")

	ctx := context.Background()

	// init service
	svc := newService()
	// init auth svc
	svc = service.AuthenticationMiddleware(BifrostConf.Service.AuthServerAddr)(svc)

	// init kit logger
	svc = logging.LoggingMiddleware(config.KitLogger)(svc)
	defer svc.Stop()

	// init kit endpoint
	endpoints := endpoint.NewBifrostEndpoints(svc)

	transport.ChunkSize = BifrostConf.Service.ChunckSize
	handler := transport.NewBifrostServer(ctx, endpoints)
	healthCheckHandler := transport.NewHealthCheck(ctx, endpoints)

	lis, lisErr := net.Listen("tcp", fmt.Sprintf(":%d", BifrostConf.Service.Port))
	if lisErr != nil {
		return lisErr
	}
	defer lis.Close()

	gRPCServer := grpc.NewServer(grpc.MaxSendMsgSize(transport.ChunkSize))
	bifrostpb.RegisterBifrostServiceServer(gRPCServer, handler)
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
