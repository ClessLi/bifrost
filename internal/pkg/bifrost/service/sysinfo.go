package service

type systemInfo struct {
	// DONE: 添加web服务版本信息、web服务状态信息(README.md需调整相关接口文档)
	OS   string `json:"system"`
	Time string `json:"time"`
	Cpu  string `json:"cpu"`
	Mem  string `json:"mem"`
	Disk string `json:"disk"`
	//ServersStatus  []string `json:"servers_status"`
	ServersStatus map[string]string `json:"servers_status"`
	//ServersVersion []string `json:"servers_version"`
	ServersVersion map[string]string `json:"servers_version"`
	BifrostVersion string            `json:"bifrost_version"`
}

var SysInfo = new(systemInfo)
