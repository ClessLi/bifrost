package auth

import (
	"github.com/ClessLi/bifrost/internal/pkg/auth/service"
	"google.golang.org/grpc"
	"time"
)

type Client struct {
	*grpc.ClientConn
	service.Service
}

func NewClientFromGRPCServerAddress(svrAddr string) (*Client, error) {
	conn, err := grpc.Dial(svrAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		return nil, err
	}

	eps := NewAuthClientEndpoints(
		makeLoginClientEndpoint(conn),
		makeVerifyClientEndpoint(conn),
		nil,
	)

	return NewClient(conn, eps), nil
}

func NewClient(conn *grpc.ClientConn, endpoints AuthClientEndpoints) *Client {
	return &Client{
		ClientConn: conn,
		Service:    endpoints,
	}
}
