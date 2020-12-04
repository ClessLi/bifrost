package service

const (
	unknown status = iota
	disable
	abnormal
	normal
)

type status int

type systemInfo struct {
	// DONE: 添加web服务版本信息、web服务状态信息(README.md需调整相关接口文档)
	OS             string         `json:"system"`
	Time           string         `json:"time"`
	Cpu            string         `json:"cpu"`
	Mem            string         `json:"mem"`
	Disk           string         `json:"disk"`
	StatusList     []ServerStatus `json:"status_list"`
	BifrostVersion string         `json:"bifrost_version"`
}

type ServerStatus struct {
	Name    string `json:"name"`
	Status  status `json:"status"`
	Version string `json:"version"`
}

func newStatus(name string) ServerStatus {
	return ServerStatus{
		Name:    name,
		Status:  unknown,
		Version: "unknown",
	}
}

func (s *ServerStatus) setVersion(v string) {
	s.Version = v
}

func (s *ServerStatus) setStatus(status status) {
	s.Status = status
}

var SysInfo = new(systemInfo)
