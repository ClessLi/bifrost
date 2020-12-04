package bifrost

import (
	"context"
	"errors"
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/config"
	bifrostendpoint "github.com/ClessLi/bifrost/internal/pkg/bifrost/endpoint"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service"
	"github.com/ClessLi/skirnir/pkg/discover"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"time"
)

type Client struct {
	*grpc.ClientConn
	service.Service
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
	eps := bifrostendpoint.BifrostEndpoints{
		ViewConfigEndpoint: grpctransport.NewClient(
			conn,
			"bifrostpb.BifrostService",
			"ViewConfig",
			decodeRequest,
			encodeResponse,
			&bifrostpb.OperateResponse{},
		).Endpoint(),
		GetConfigEndpoint: grpctransport.NewClient(
			conn,
			"bifrostpb.BifrostService",
			"GetConfig",
			decodeRequest,
			encodeResponse,
			&bifrostpb.ConfigResponse{},
		).Endpoint(),
		UpdateConfigEndpoint: grpctransport.NewClient(
			conn,
			"bifrostpb.BifrostService",
			"UpdateConfig",
			decodeRequest,
			encodeResponse,
			&bifrostpb.OperateResponse{},
		).Endpoint(),
		ViewStatisticsEndpoint: grpctransport.NewClient(
			conn,
			"bifrostpb.BifrostService",
			"ViewStatistics",
			decodeRequest,
			encodeResponse,
			&bifrostpb.OperateResponse{},
		).Endpoint(),
		StatusEndpoint: grpctransport.NewClient(
			conn,
			"bifrostpb.BifrostService",
			"Status",
			decodeRequest,
			encodeResponse,
			&bifrostpb.OperateResponse{},
		).Endpoint(),
		WatchLogEndpoint: bifrostendpoint.MakeWatchLogClientEndpoint(conn),
	}
	return &Client{
		ClientConn: conn,
		Service:    eps,
	}, nil
}

func decodeRequest(ctx context.Context, r interface{}) (request interface{}, err error) {
	return r, nil
}

func encodeResponse(ctx context.Context, r interface{}) (response interface{}, err error) {
	return r, nil
}
