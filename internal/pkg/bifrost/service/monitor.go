package service

type Monitor interface {
	DisplayStatus() ([]byte, error)
	Start() error
	Stop() error
}

type WebServerMonitor interface {
	ShowVersion() string
	DisplayStatus() status
}
