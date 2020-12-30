package service

type RequestType int

// WebServerType, web服务器类型对象，定义web服务器所属类型
type WebServerType string

const (
	Unknown RequestType = iota
	DisplayConfig
	GetConfig
	ShowStatistics
	DisplayStatus

	UpdateConfig

	WatchLog

	// Web服务类型
	NGINX WebServerType = "nginx"
	HTTPD WebServerType = "httpd"
)
