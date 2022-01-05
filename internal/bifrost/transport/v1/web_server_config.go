package v1

import "github.com/go-kit/kit/transport/grpc"

type WebServerConfigHandler interface {
	HandlerGet() grpc.Handler
	HandlerUpdate() grpc.Handler
}
