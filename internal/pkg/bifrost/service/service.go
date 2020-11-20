package service

import (
	"encoding/json"
	authSvc "github.com/ClessLi/bifrost/internal/pkg/auth/service"
	ngStatistics "github.com/ClessLi/bifrost/pkg/statistics/nginx"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"sync"
)

type Service interface {
	ViewConfig(ctx context.Context, token, svrName string) (data []byte, err error)
	GetConfig(ctx context.Context, token, srvName string) (jsonData []byte, err error)
	UpdateConfig(ctx context.Context, token, svrName string, jsonData []byte) (data []byte, err error)
	ViewStatistics(ctx context.Context, token, svrName string) (jsonData []byte, err error)
	Status(ctx context.Context, token string) (jsonData []byte, err error)
}

// BifrostService, bifrost配置文件对象中web服务器信息结构体，定义管控的web服务器配置文件相关信息
type BifrostService struct {
	Port           int     `yaml:"Port"`
	ChunckSize     int     `yaml:"ChunkSize"`
	AuthServerAddr string  `yaml:"AuthServerAddr"`
	Infos          []*Info `yaml:"Infos,flow"`
	monitorChan    chan int
	authConn       *grpc.ClientConn
	authSvc        authSvc.Service
	waitGroup      *sync.WaitGroup
}

func (b *BifrostService) ViewConfig(ctx context.Context, token, svrName string) (data []byte, err error) {
	var pass bool
	pass, err = b.checkToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if !pass {
		return nil, UnknownErrCheckToken
	}

	var info *Info
	info, err = b.getInfo(svrName)
	if err != nil {
		return nil, err
	}
	for _, s := range info.nginxConfig.String() {
		data = append(data, []byte(s)...)
	}
	if data == nil {
		err = ErrDataNotParsed
		return nil, err
	}
	return data, err
}

func (b *BifrostService) GetConfig(ctx context.Context, token, svrName string) (jsonData []byte, err error) {
	var pass bool
	pass, err = b.checkToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if !pass {
		return nil, UnknownErrCheckToken
	}

	var info *Info
	info, err = b.getInfo(svrName)
	if err != nil {
		return nil, err
	}

	return json.Marshal(info.nginxConfig)
}

func (b *BifrostService) UpdateConfig(ctx context.Context, token, svrName string, jsonData []byte) (data []byte, err error) {
	var pass bool
	pass, err = b.checkToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if !pass {
		return nil, UnknownErrCheckToken
	}

	var info *Info
	info, err = b.getInfo(svrName)
	if err != nil {
		return nil, err
	}

	return info.update(jsonData)
}

func (b *BifrostService) ViewStatistics(ctx context.Context, token, svrName string) (jsonData []byte, err error) {
	var pass bool
	pass, err = b.checkToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if !pass {
		return nil, UnknownErrCheckToken
	}

	var info *Info
	info, err = b.getInfo(svrName)
	if err != nil {
		return nil, err
	}

	httpServersNum, httpServers := ngStatistics.HTTPServers(info.nginxConfig)
	httpPorts := ngStatistics.HTTPPorts(info.nginxConfig)
	streamServersNum, streamPorts := ngStatistics.StreamServers(info.nginxConfig)
	statistics := struct {
		HttpSvrsNum   int              `json:"http_svrs_num"`
		HttpSvrs      map[string][]int `json:"http_svrs"`
		HttpPorts     []int            `json:"http_ports"`
		StreamSvrsNum int              `json:"stream_svrs_num"`
		StreamPorts   []int            `json:"stream_ports"`
	}{HttpSvrsNum: httpServersNum, HttpSvrs: httpServers, HttpPorts: httpPorts, StreamSvrsNum: streamServersNum, StreamPorts: streamPorts}
	return json.Marshal(statistics)
}

func (b *BifrostService) Status(ctx context.Context, token string) (jsonData []byte, err error) {
	var pass bool
	pass, err = b.checkToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if !pass {
		return nil, UnknownErrCheckToken
	}

	// TODO: SysInfo lock mechanism
	return json.Marshal(SysInfo)
}

func (b BifrostService) GetPort() int {
	return b.Port
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
