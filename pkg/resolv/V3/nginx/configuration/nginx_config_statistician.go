package configuration

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"regexp"
	"strconv"

	"github.com/marmotedu/errors"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	utilsV2 "github.com/ClessLi/bifrost/pkg/resolv/V2/utils"
)

var (
	RegPortValue       = regexp.MustCompile(`^listen\s*(\d+)\s*\S*$`)
	RegServerNameValue = regexp.MustCompile(`^server_name\s*(.+)$`)
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
	portDirective := ctx.QueryByKeyWords(context.NewKeyWords(context_type.TypeDirective).SetRegexMatchingValue("listen .*")).Target()
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
	httpServerPoses := ctx.
		QueryByKeyWords(context.NewKeyWords(context_type.TypeHttp)).Target().
		QueryAllByKeyWords(context.NewKeyWords(context_type.TypeServer))
	httpServers := make([]context.Context, 0)
	for _, pos := range httpServerPoses {
		httpServers = append(httpServers, pos.Target())
	}

	return Ports(httpServers)
}

func HttpServers(ctx context.Context) (int, map[string][]int) {
	serverCount := 0
	httpServerPoses := ctx.
		QueryByKeyWords(context.NewKeyWords(context_type.TypeHttp)).Target().
		QueryAllByKeyWords(context.NewKeyWords(context_type.TypeServer))

	serverPortCount := make(map[string][]int)
	for _, pos := range httpServerPoses {
		serverCount++
		servernameDirective := pos.Target().QueryByKeyWords(context.NewKeyWords(context_type.TypeDirective).SetRegexMatchingValue("^server_name .*")).Target()
		if servernameDirective.Error() != nil {
			if !errors.IsCode(servernameDirective.Error(), code.ErrV3ContextNotFound) {
				return 0, nil
			}

			continue
		}
		servername := servernameDirective.(*local.Directive).Params
		serverport := Port(pos.Target().Father())
		if serverport > 0 {
			serverPortCount[servername] = append(serverPortCount[servername], serverport)
		}
	}

	return serverCount, serverPortCount
}

func StreamServers(ctx context.Context) []int {
	streamServerPoses := ctx.
		QueryByKeyWords(context.NewKeyWords(context_type.TypeStream)).Target().
		QueryAllByKeyWords(context.NewKeyWords(context_type.TypeServer))

	streamServers := make([]context.Context, 0)
	for _, pos := range streamServerPoses {
		streamServers = append(streamServers, pos.Target())
	}

	return Ports(streamServers)
}
