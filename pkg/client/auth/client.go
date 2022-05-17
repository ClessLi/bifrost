package auth

import (
	"time"

	"google.golang.org/grpc"

	"github.com/yongPhone/bifrost/internal/pkg/auth/service"
)

type Client struct {
	*grpc.ClientConn
	service.Service
}

func NewClient(svrAddr string) (*Client, error) {
	conn, err := grpc.Dial(svrAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second)) //nolint:staticcheck
	if err != nil {
		return nil, err
	}

	eps := newAuthClientEndpoints(
		makeLoginClientEndpoint(conn),
		makeVerifyClientEndpoint(conn),
	)

	return newClient(conn, eps), nil
}

func newClient(conn *grpc.ClientConn, endpoints *authClientEndpoints) *Client {
	return &Client{
		ClientConn: conn,
		Service:    endpoints,
	}
}
