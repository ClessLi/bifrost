package v1

import (
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/ClessLi/bifrost/internal/pkg/service_register"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	// GRPCInstancePrefix defines the grpc server instance prefix used by all bifrost service.
	GRPCInstancePrefix = "com.github.ClessLi.api"
)

type bifrostServiceRegister struct {
	instancePrefixName string
	factory            Factory
}

func (b *bifrostServiceRegister) Generate() map[string]service_register.ServiceRegister {
	return map[string]service_register.ServiceRegister{
		b.instancePrefixName + ".bifrostpb.WebServerConfig": func(server *grpc.Server, healthzSvr *health.Server) {
			if healthzSvr != nil {
				healthzSvr.SetServingStatus(b.instancePrefixName+".bifrostpb.WebServerConfig", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
			}
			pbv1.RegisterWebServerConfigServer(server, b.factory.WebServerConfig())
		},
	}
}

func NewBifrostServiceRegister(factory Factory) service_register.ServiceRegisterGenerator {
	return &bifrostServiceRegister{
		instancePrefixName: GRPCInstancePrefix,
		factory:            factory,
	}
}