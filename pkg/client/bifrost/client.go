package bifrost

import (
	"errors"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/config"
	bifrostendpoint "github.com/ClessLi/bifrost/internal/pkg/bifrost/endpoint"
	"github.com/ClessLi/skirnir/pkg/discover"
	"google.golang.org/grpc"
	"time"
)

type Client struct {
	*grpc.ClientConn
	*bifrostendpoint.BifrostClientEndpoints
}

func NewClientFromConsul(consulHost string, consulPort uint16) (*Client, error) {
	discoveryClient, err := discover.NewKitConsulDiscoveryClient(consulHost, consulPort)
	if err != nil {
		return nil, err
	}
	factory := func(instance string) (interface{}, error) {
		return NewClient(instance)
	}
	relay, err := discoveryClient.DiscoverServicesClient("com.github.ClessLi.api.bifrost", config.KitLogger, factory)
	if err != nil {
		return nil, err
	}
	if client, ok := relay.(*Client); ok {
		return client, nil
	}
	return nil, errors.New("failed to initialize Bifrost service client")
}

func NewClient(svrAddr string) (*Client, error) {
	conn, err := grpc.Dial(svrAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		return nil, err
	}
	eps := bifrostendpoint.NewBifrostClient(conn)
	return &Client{
		ClientConn:             conn,
		BifrostClientEndpoints: eps,
	}, nil
}
