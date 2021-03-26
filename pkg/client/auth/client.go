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

func NewClient(svrAddr string) (*Client, error) {
	conn, err := grpc.Dial(svrAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
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
