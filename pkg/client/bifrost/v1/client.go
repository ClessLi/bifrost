package v1

import (
	txpclient "github.com/ClessLi/bifrost/pkg/client/bifrost/v1/transport"
	"google.golang.org/grpc"
)

type Client struct {
	*grpc.ClientConn
	txpclient.Factory
}
