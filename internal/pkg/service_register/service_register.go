package service_register

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
)

type ServiceRegister func(server *grpc.Server, healthzSvr *health.Server)

type ServiceRegisterGenerator interface {
	Generate() map[string]ServiceRegister
}
