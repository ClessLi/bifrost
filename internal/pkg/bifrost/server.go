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
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

func ServerRun() error {

	if !isInit {
		return errors.New("service related configuration not initialized")
	}

	Log(DEBUG, "Listening system call signal")
	go ListenSignal(signalChan)
	Log(DEBUG, "Listened system call signal")

	ctx := context.Background()

	var svc service.Service
	// init auth svc
	connAuthErr := BifrostConf.Service.ConnAuthSvr()
	if connAuthErr != nil {
		return connAuthErr
	}
	defer BifrostConf.Service.AuthSvrConnClose()

	// init bifrost svr

	// back svr run
	BifrostConf.Service.Run()
	defer BifrostConf.Service.KillCoroutines()

	svc = BifrostConf.Service
	svc = logging.LoggingMiddleware(config.KitLogger)(svc)

	endpts := endpoint.BifrostEndpoints{
		ViewConfigEndpoint:     endpoint.MakeViewConfigEndpoint(svc),
		GetConfigEndpoint:      endpoint.MakeGetConfigEndpoint(svc),
		UpdateConfigEndpoint:   endpoint.MakeUpdateConfigEndpoint(svc),
		ViewStatisticsEndpoint: endpoint.MakeViewStatisticsEndpoint(svc),
		StatusEndpoint:         endpoint.MakeStatusEndpoint(svc),
		WatchLogEndpoint:       endpoint.MakeWatchLogEndpoint(svc),
		HealthCheckEndpoint:    endpoint.MakeHealthCheckEndpoint(svc),
	}

	transport.ChunkSize = BifrostConf.Service.ChunckSize
	handler := transport.NewBifrostServer(ctx, endpts)
	healthCheckHandler := transport.NewHealthCheck(ctx, endpts)

	lis, lisErr := net.Listen("tcp", fmt.Sprintf(":%d", BifrostConf.Service.Port))
	if lisErr != nil {
		return lisErr
	}
	defer lis.Close()

	gRPCServer := grpc.NewServer(grpc.MaxSendMsgSize(transport.ChunkSize))
	bifrostpb.RegisterBifrostServiceServer(gRPCServer, handler)
	grpc_health_v1.RegisterHealthServer(gRPCServer, healthCheckHandler)
	fmt.Println(logoStr)
	svrErrChan := make(chan error, 1)
	go func() {
		svrErr := gRPCServer.Serve(lis)
		Log(NOTICE, "bifrost service is running on %s", lis.Addr())
		svrErrChan <- svrErr
	}()

	go registerToRA(svrErrChan)
	defer deregisterToRA()

	var stopErr error
	select {
	case s := <-signalChan:
		if s == 9 {
			Log(DEBUG, "stopping...")
			gRPCServer.Stop()
		}
		Log(DEBUG, "stop signal error")
	case stopErr = <-svrErrChan:
		gRPCServer.Stop()
		break
	}
	return stopErr
}
