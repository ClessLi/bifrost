package web_server_manager

const (
	Unknown State = iota
	Disabled
	Initializing
	Abnormal
	Normal
)

type State int

type WebServerConfigViewer interface {
	DisplayConfig() ([]byte, error)
	GetConfig() ([]byte, error)
	ShowStatistics() ([]byte, error)
	DisplayWebServerStatus() State
	DisplayWebServerVersion() string
}

type WebServerConfigUpdater interface {
	UpdateConfig(data []byte) error
}

type WebServerConfigWatcher interface {
	WatchLog(logName string) (LogWatcher, error)
}

type WebServerConfigService interface {
	serverName() string
	//checkConfigsHash() (bool, error)
	configReload() error
	configBackup() error
	WebServerConfigViewer
	WebServerConfigUpdater
	WebServerConfigWatcher
}

type WebServerConfigServicesHandler struct {
	services map[string]WebServerConfigService
}

func (m WebServerConfigServicesHandler) ServerNames() []string {
	names := make([]string, 0)
	for s, service := range m.services {
		if service != nil {
			names = append(names, s)
		}
	}
	return names
}

func (m WebServerConfigServicesHandler) DisplayConfig(serverName string) ([]byte, error) {
	offstage, err := m.getService(serverName)
	if err != nil {
		return nil, err
	}
	return offstage.DisplayConfig()
}

func (m WebServerConfigServicesHandler) GetConfig(serverName string) ([]byte, error) {
	offstage, err := m.getService(serverName)
	if err != nil {
		return nil, err
	}
	return offstage.GetConfig()
}

func (m WebServerConfigServicesHandler) ShowStatistics(serverName string) ([]byte, error) {
	offstage, err := m.getService(serverName)
	if err != nil {
		return nil, err
	}
	return offstage.ShowStatistics()
}

func (m WebServerConfigServicesHandler) DisplayStatus(serverName string) (State, error) {
	offstage, err := m.getService(serverName)
	if err != nil {
		return Unknown, err
	}
	return offstage.DisplayWebServerStatus(), nil
}

func (m WebServerConfigServicesHandler) DisplayVersion(serverName string) (string, error) {
	offstage, err := m.getService(serverName)
	if err != nil {
		return "", err
	}
	return offstage.DisplayWebServerVersion(), nil
}

func (m *WebServerConfigServicesHandler) UpdateConfig(serverName string, data []byte) error {
	offstage, err := m.getService(serverName)
	if err != nil {
		return err
	}
	return offstage.UpdateConfig(data)
}

func (m WebServerConfigServicesHandler) WatchLog(serverName, logName string) (LogWatcher, error) {
	offstage, err := m.getService(serverName)
	if err != nil {
		return nil, err
	}
	return offstage.WatchLog(logName)
}

func (m WebServerConfigServicesHandler) getService(serverName string) (WebServerConfigService, error) {
	if backstage, ok := m.services[serverName]; ok {
		return backstage, nil
	}
	return nil, ErrOffstageNotExist
}

func newWebServerConfigServiceHandler(offstages ...WebServerConfigService) *WebServerConfigServicesHandler {
	if offstages == nil || len(offstages) < 1 {
		return nil
	}
	handler := &WebServerConfigServicesHandler{services: make(map[string]WebServerConfigService)}
	for _, offstage := range offstages {
		handler.services[offstage.serverName()] = offstage
	}
	return handler
}
