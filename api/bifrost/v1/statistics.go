package v1

type Statistics struct {
	HttpSvrsNum   int              `json:"http_svrs_num"`
	HttpSvrs      map[string][]int `json:"http_svrs"`
	HttpPorts     []int            `json:"http_ports"`
	StreamSvrsNum int              `json:"stream_svrs_num"`
	StreamPorts   []int            `json:"stream_ports"`
}
