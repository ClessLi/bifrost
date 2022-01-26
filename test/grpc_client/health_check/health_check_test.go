package health_check

import (
	"context"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/fake"
	clientv1 "github.com/ClessLi/bifrost/pkg/client/grpc_health_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"testing"
)

//func TestHealthCheck(t *testing.T) {
//	conn, _ := grpc.Dial(bifrostSvrAddr, grpc.WithInsecure())
//	defer conn.Close()
//	c := grpc_health_v1.NewHealthClient(conn)
//	reply, _ := c.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
//	t.Log(reply.Status)
//}
func testGRPCServer() (*grpc.Server, *health.Server) {
	server := grpc.NewServer()
	pbv1.RegisterWebServerConfigServer(server, fake.New().WebServerConfig())
	pbv1.RegisterWebServerStatisticsServer(server, fake.New().WebServerStatistics())
	pbv1.RegisterWebServerStatusServer(server, fake.New().WebServerStatus())
	pbv1.RegisterWebServerLogWatcherServer(server, fake.New().WebServerLogWatcher())
	healthSvr := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthSvr)
	return server, healthSvr
}

func TestNewClient(t *testing.T) {
	server, healthSvr := testGRPCServer()
	address := "192.168.220.1:8888"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		t.Fatalf(err.Error())
	}

	go func() {
		err = server.Serve(lis)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()
	defer server.Stop()

	type args struct {
		svrAddr string
		opts    []grpc.DialOption
	}
	tests := []struct {
		name        string
		args        args
		servicename string
		wantState   clientv1.HealthStatus
	}{
		{
			name: "test bifrost web server config",
			args: args{
				svrAddr: address,
				opts:    []grpc.DialOption{grpc.WithInsecure()},
			},
			servicename: "bifrostpb.WebServerConfig",
			wantState:   clientv1.SERVING,
		},
		{
			name: "test bifrost web server statistics",
			args: args{
				svrAddr: address,
				opts:    []grpc.DialOption{grpc.WithInsecure()},
			},
			servicename: "bifrostpb.WebServerStatistics",
			wantState:   clientv1.SERVING,
		},
		{
			name: "test bifrost web server status",
			args: args{
				svrAddr: address,
				opts:    []grpc.DialOption{grpc.WithInsecure()},
			},
			servicename: "bifrostpb.WebServerStatus",
			wantState:   clientv1.SERVING,
		},
		{
			name: "test bifrost web server log watcher",
			args: args{
				svrAddr: address,
				opts:    []grpc.DialOption{grpc.WithInsecure()},
			},
			servicename: "bifrostpb.WebServerLogWatcher",
			wantState:   clientv1.SERVING,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			healthSvr.SetServingStatus(tt.servicename, grpc_health_v1.HealthCheckResponse_ServingStatus(tt.wantState))
			got, err := clientv1.NewClient(tt.args.svrAddr, tt.args.opts...)
			if err != nil {
				t.Fatalf(err.Error())
			}
			state, _ := got.Check(context.Background(), tt.servicename)
			if state != tt.wantState {
				t.Fatalf("got state %s, want %s", clientv1.StatusString(state), clientv1.StatusString(tt.wantState))
			}
			t.Logf("service %s is %s", tt.servicename, clientv1.StatusString(state))
		})
	}
}
