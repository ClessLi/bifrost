package service

import (
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service/web_server_manager"
)

// BifrostService, bifrost配置文件对象中web服务器信息结构体，定义管控的web服务器配置文件相关信息
type BifrostService struct {
	webServerConfigServicesHandler *web_server_manager.WebServerConfigServicesHandler
	monitor                        Monitor
}

func (b BifrostService) DisplayConfig(serverName string) ([]byte, error) {
	return b.webServerConfigServicesHandler.DisplayConfig(serverName)
}

func (b BifrostService) GetConfig(serverName string) ([]byte, error) {
	return b.webServerConfigServicesHandler.GetConfig(serverName)
}

func (b BifrostService) ShowStatistics(serverName string) ([]byte, error) {
	return b.webServerConfigServicesHandler.ShowStatistics(serverName)
}

func (b BifrostService) DisplayStatus() ([]byte, error) {
	return b.monitor.DisplayStatus()
}

func (b *BifrostService) UpdateConfig(serverName string, data []byte) error {
	return b.webServerConfigServicesHandler.UpdateConfig(serverName, data)
}

func (b BifrostService) WatchLog(serverName, logName string) (web_server_manager.LogWatcher, error) {
	return b.webServerConfigServicesHandler.WatchLog(serverName, logName)
}
