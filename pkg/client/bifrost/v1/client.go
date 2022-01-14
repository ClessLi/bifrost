package v1

import (
	epclient "github.com/ClessLi/bifrost/pkg/client/bifrost/v1/endpoint"
	svcclient "github.com/ClessLi/bifrost/pkg/client/bifrost/v1/service"
	txpclient "github.com/ClessLi/bifrost/pkg/client/bifrost/v1/transport"
	"google.golang.org/grpc"
)

type Client struct {
	conn *grpc.ClientConn
	svcclient.Factory
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func newClient(conn *grpc.ClientConn) *Client {
	return &Client{
		conn:    conn,
		Factory: svcclient.New(epclient.New(txpclient.New(conn))),
	}
}

func New(svraddr string, opts ...grpc.DialOption) (*Client, error) {
	conn, err := grpc.Dial(svraddr, opts...)
	if err != nil {
		return nil, err
	}
	return newClient(conn), nil
}

// TODO: NewClientFromConsul
