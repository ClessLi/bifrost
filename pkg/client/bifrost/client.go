package bifrost

import (
	"errors"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/config"
	"github.com/ClessLi/skirnir/pkg/discover"
	"google.golang.org/grpc"
	"os"
	"time"
)

type Client struct {
	*grpc.ClientConn
	*bifrostClientEndpoints
}

func NewClientFromConsul(consulHost string, consulPort uint16) (*Client, error) {
	discoveryClient, err := discover.NewKitConsulDiscoveryClient(consulHost, consulPort)
	if err != nil {
		return nil, err
	}
	factory := func(instance string) (interface{}, error) {
		return NewClientFromServerAddress(instance)
	}
	relay, err := discoveryClient.DiscoverServicesClient("com.github.ClessLi.api.bifrost", config.KitLogger(os.Stdout), factory)
	if err != nil {
		return nil, err
	}
	if client, ok := relay.(*Client); ok {
		return client, nil
	}
	return nil, errors.New("failed to initialize Bifrost service client")
}

func NewClientFromServerAddress(svrAddr string) (*Client, error) {
	conn, err := grpc.Dial(svrAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		return nil, err
	}
	eps := NewBifrostClientEndpoints(conn)
	return NewClient(conn, eps), nil
}

func NewClient(conn *grpc.ClientConn, endpoints *bifrostClientEndpoints) *Client {
	return &Client{
		ClientConn:             conn,
		bifrostClientEndpoints: endpoints,
	}
}
