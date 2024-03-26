package grpc_health_v1

import (
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"github.com/ClessLi/skirnir/pkg/discover"
	kitzaplog "github.com/go-kit/kit/log/zap"
	"github.com/marmotedu/errors"
	"google.golang.org/grpc"
)

type Client struct {
	*grpc.ClientConn
	*endpoints
}

func NewClientFromConsul(consulHost string, consulPort uint16, opts ...grpc.DialOption) (*Client, error) {
	discoveryClient, err := discover.NewKitConsulDiscoveryClient(consulHost, consulPort)
	if err != nil {
		return nil, err
	}
	factory := func(instance string) (interface{}, error) {
		return NewClient(instance, opts...)
	}
	relay, err := discoveryClient.DiscoverServicesClient(
		"com.github.ClessLi.api.bifrost",
		kitzaplog.NewZapSugarLogger(logV1.ZapLogger(), logV1.InfoLevel),
		factory,
	)
	if err != nil {
		return nil, err
	}
	if client, ok := relay.(*Client); ok {
		return client, nil
	}

	return nil, errors.New("failed to initialize Health Check service client")
}

func NewClient(svrAddr string, opts ...grpc.DialOption) (*Client, error) {
	conn, err := grpc.Dial(svrAddr, opts...)
	if err != nil {
		return nil, err
	}
	eps := makeEndpoints(conn)

	return newClient(conn, eps), nil
}

func newClient(conn *grpc.ClientConn, endpoints *endpoints) *Client {
	return &Client{
		ClientConn: conn,
		endpoints:  endpoints,
	}
}
