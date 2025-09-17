package configuration

import (
	"strconv"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	utilsV2 "github.com/ClessLi/bifrost/pkg/resolv/V2/utils"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"

	"github.com/marmotedu/errors"
)

type HttpInfo struct {
	ServerCount     int
	ServerPortCount map[string][]int
	PortCount       []int
}

type StreamInfo struct {
	PortCount []int
}

type Statistician interface {
	HttpInfo() HttpInfo
	StreamInfo() StreamInfo
	Statistics() *v1.Statistics
}

type statistician struct {
	configuration NginxConfig
}

func (s *statistician) HttpInfo() HttpInfo {
	serverCount, serverPortCount := HttpServers(s.configuration.Main())

	return HttpInfo{
		ServerCount:     serverCount,
		ServerPortCount: serverPortCount,
		PortCount:       HttpPorts(s.configuration.Main()),
	}
}

func (s *statistician) StreamInfo() StreamInfo {
	return StreamInfo{
		PortCount: StreamServers(s.configuration.Main()),
	}
}

func (s *statistician) Statistics() *v1.Statistics {
	httpInfo := s.HttpInfo()
	streamInfo := s.StreamInfo()

	return &v1.Statistics{
		HttpSvrsNum:   httpInfo.ServerCount,
		HttpSvrs:      httpInfo.ServerPortCount,
		HttpPorts:     httpInfo.PortCount,
		StreamSvrsNum: len(streamInfo.PortCount),
		StreamPorts:   streamInfo.PortCount,
	}
}

func NewStatistician(c NginxConfig) Statistician {
	return &statistician{configuration: c}
}

func Port(ctx context.Context) int {
	portDirective := ctx.ChildrenPosSet().
		QueryOne(context.NewKeyWords(context_type.TypeDirective).
			SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc).
			SetRegexpMatchingValue(context.RegexpMatchingListenPortValue)).
		Target()
	if portDirective.Error() != nil {
		return -1
	}
	port, err := strconv.Atoi(portDirective.(*local.Directive).Params)
	if err != nil {
		return -1
	}

	return port
}

func Ports(contexts []context.Context) []int {
	ports := make([]int, 0)
	for _, ctx := range contexts {
		port := Port(ctx)
		if port > 0 {
			ports = utilsV2.SortInsertUniqInt(ports, port)
		}
	}

	return ports
}

func HttpPorts(ctx context.Context) []int {
	return Ports(
		ctx.ChildrenPosSet().
			QueryOne(context.NewKeyWords(context_type.TypeHttp).SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).
			QueryAll(context.NewKeyWords(context_type.TypeServer).SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).
			Targets(),
	)
}

func HttpServers(ctx context.Context) (int, map[string][]int) {
	serverCount := 0
	serverPortCount := make(map[string][]int)
	err := ctx.ChildrenPosSet().
		QueryOne(context.NewKeyWords(context_type.TypeHttp).SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).
		QueryAll(context.NewKeyWords(context_type.TypeServer).SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).
		Map(
			func(pos context.Pos) (context.Pos, error) {
				serverCount++
				servernameDirective := pos.QueryOne(
					context.NewKeyWords(context_type.TypeDirective).
						SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc).
						SetRegexpMatchingValue(context.RegexpMatchingServerNameValue),
				).Target()
				if servernameDirective.Error() != nil {
					if !errors.IsCode(servernameDirective.Error(), code.ErrV3ContextNotFound) {
						return pos, servernameDirective.Error()
					}

					return pos, nil
				}
				servername := servernameDirective.(*local.Directive).Params
				serverport := Port(pos.Target())
				if serverport > 0 {
					serverPortCount[servername] = append(serverPortCount[servername], serverport)
				}

				return pos, nil
			},
		).
		Error()
	if err != nil {
		return 0, nil
	}

	return serverCount, serverPortCount
}

func StreamServers(ctx context.Context) []int {
	return Ports(
		ctx.ChildrenPosSet().
			QueryOne(context.NewKeyWords(context_type.TypeStream).SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).
			QueryAll(context.NewKeyWords(context_type.TypeServer).SetSkipQueryFilter(context.SkipDisabledCtxFilterFunc)).
			Targets(),
	)
}
