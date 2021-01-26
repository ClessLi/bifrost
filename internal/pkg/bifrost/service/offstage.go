package service

import "github.com/ClessLi/bifrost/internal/pkg/bifrost/service/web_server_manager"

type offstageViewer interface {
	DisplayConfig(serverName string) ([]byte, error)
	GetConfig(serverName string) ([]byte, error)
	ShowStatistics(serverName string) ([]byte, error)
	DisplayStatus() ([]byte, error)
}

type offstageUpdater interface {
	UpdateConfig(serverName string, data []byte) error
}

type offstageWatcher interface {
	WatchLog(serverName, logName string) (web_server_manager.LogWatcher, error)
}
