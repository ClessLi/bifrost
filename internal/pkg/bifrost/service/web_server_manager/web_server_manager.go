package web_server_manager

const (
	Unknown status = iota
	Disabled
	Abnormal
	Normal
)

type status int

type WebServerManager interface {
	DisplayConfig() ([]byte, error)
	GetConfig() ([]byte, error)
	ShowStatistics() ([]byte, error)
	UpdateConfig(data []byte, param string) error
	WatchLog(logName string) (Watcher, error)
	autoBackup()
	autoReload()
	ShowVersion() string
	DisplayStatus() status
	Start() error
	Stop() error
}
