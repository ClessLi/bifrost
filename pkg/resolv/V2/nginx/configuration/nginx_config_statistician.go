package configuration

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/utils"
	"regexp"
	"strconv"
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
	ServerCount int
	PortCount   []int
}

type Statistics struct {
	HttpSvrsNum   int              `json:"http_svrs_num"`
	HttpSvrs      map[string][]int `json:"http_svrs"`
	HttpPorts     []int            `json:"http_ports"`
	StreamSvrsNum int              `json:"stream_svrs_num"`
	StreamPorts   []int            `json:"stream_ports"`
}

type Statistician struct {
	configuration *configuration
}

func (s Statistician) HttpInfo() HttpInfo {
	serverCount, serverPortCount := HttpServers(s.configuration)
	return HttpInfo{
		ServerCount:     serverCount,
		ServerPortCount: serverPortCount,
		PortCount:       HttpPorts(s.configuration),
	}
}

func (s Statistician) StreamInfo() StreamInfo {
	serverCount, portCount := StreamServers(s.configuration)
	return StreamInfo{
		ServerCount: serverCount,
		PortCount:   portCount,
	}
}

func (s Statistician) Statistics() Statistics {
	httpInfo := s.HttpInfo()
	streamInfo := s.StreamInfo()
	return Statistics{
		HttpSvrsNum:   httpInfo.ServerCount,
		HttpSvrs:      httpInfo.ServerPortCount,
		HttpPorts:     httpInfo.PortCount,
		StreamSvrsNum: streamInfo.ServerCount,
		StreamPorts:   streamInfo.PortCount,
	}
}

func NewStatistician(c *configuration) *Statistician {
	return &Statistician{configuration: c}
}

func Port(q Queryer) int {
	keyQueryer, err := q.Query("key:sep: :reg: listen .*")
	if err != nil {
		return -1
	}
	portValue := keyQueryer.Self().GetValue()
	if RegPortValue.MatchString(portValue) {
		portStr := RegPortValue.FindStringSubmatch(portValue)[1]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return -1
		}
		return port
	}
	return -1
}

func Ports(qs []Queryer) []int {
	ports := make([]int, 0)
	for _, q := range qs {
		port := Port(q)
		if port > 0 {
			ports = utils.SortInsertUniqInt(ports, port)
		}
	}
	return ports
}

func HttpPorts(q Queryer) []int {
	httpQueryer, err := q.Query("http")
	if err != nil {
		return nil
	}
	serverQueryers, err := httpQueryer.QueryAll("server")
	if err != nil {
		return nil
	}
	return Ports(serverQueryers)
}

func HttpServers(q Queryer) (int, map[string][]int) {
	serverCount := 0
	httpQueryer, err := q.Query("http")
	if err != nil {
		return serverCount, nil
	}
	serverQueryers, err := httpQueryer.QueryAll("server")
	if err != nil {
		return serverCount, nil
	}

	serverPortCount := make(map[string][]int)
	for _, serverQueryer := range serverQueryers {
		serverCount++
		serverNameKeyQueryer, err := serverQueryer.Query("key:sep: :reg: server_name .*")
		if err != nil {
			if err != nginx.ErrNotFound {
				return 0, nil
			}
			continue
		}
		serverNameValue := serverNameKeyQueryer.Self().GetValue()
		if RegServerNameValue.MatchString(serverNameValue) {
			serverName := RegServerNameValue.FindStringSubmatch(serverNameValue)[1]
			port := Port(serverQueryer)
			if port > 0 {
				// 只去server里的第一个侦听端口
				serverPortCount[serverName] = append(serverPortCount[serverName], port)
			}
		}
	}
	return serverCount, serverPortCount
}

func StreamServers(q Queryer) (int, []int) {
	serverCount := 0
	streamQueryer, err := q.Query("stream")
	if err != nil {
		return serverCount, nil
	}
	serverQueryers, err := streamQueryer.QueryAll("server")
	ports := Ports(serverQueryers)
	return serverCount, ports
}
