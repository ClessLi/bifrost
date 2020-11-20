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

	// init bifrost svr

	// back svr run
	BifrostConf.Service.Run()
	defer BifrostConf.Service.KillCoroutines()

	svc = BifrostConf.Service
	svc = logging.LoggingMiddleware(config.KitLogger)(svc)
	endpt := endpoint.MakeBifrostEndpoint(svc)

	healthEndpt := endpoint.MakeHealthCheckEndpoint(svc)

	endpts := endpoint.BifrostEndpoints{
		BifrostEndpoint:     endpt,
		HealthCheckEndpoint: healthEndpt,
	}

	transport.ChunkSize = BifrostConf.Service.ChunckSize
	handler := transport.NewBifrostServer(ctx, endpts)

	lis, lisErr := net.Listen("tcp", fmt.Sprintf(":%d", BifrostConf.Service.Port))
	if lisErr != nil {
		return lisErr
	}

	gRPCServer := grpc.NewServer()
	bifrostpb.RegisterBifrostServiceServer(gRPCServer, handler)
	fmt.Println(logoStr)
	svrErrChan := make(chan error, 1)
	go func() {
		svrErr := gRPCServer.Serve(lis)
		svrErrChan <- svrErr
	}()

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
