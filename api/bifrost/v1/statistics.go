package v1

type Statistics struct {
	Http   HttpAggregate   `json:"http"`
	Stream StreamAggregate `json:"stream"`
}

type HttpAggregate struct {
	Ports   []uint16                       `json:"ports"`
	Servers map[string]HttpServerAggregate `json:"servers"`
}

type HttpServerAggregate struct {
	ServerName string   `json:"server-name"`
	Ports      []uint16 `json:"ports"`
}

type StreamAggregate struct {
	Ports []uint16 `json:"ports"`
}
