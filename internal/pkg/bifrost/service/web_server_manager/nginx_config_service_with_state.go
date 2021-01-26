package web_server_manager

type nginxConfigServiceWithState struct {
	service *nginxConfigServiceWithLocker
	status  State
}

func (n *nginxConfigServiceWithState) Status() State {
	return n.status
}

func (n *nginxConfigServiceWithState) SetState(state State) {
	n.status = state
}

func (n *nginxConfigServiceWithState) Lock() {
	n.service.Lock()
}

func (n *nginxConfigServiceWithState) Unlock() {
	n.service.Unlock()
}

//
//func (n *nginxConfigServiceWithState) checkConfigsHash() (bool, error) {
//	if n.status > Disabled {
//		return n.service.checkConfigsHash()
//	}
//	return false, ErrDisabledService
//}

func (n *nginxConfigServiceWithState) configReload() error {
	if n.status > Disabled {
		return n.service.configReload()
	}
	return ErrDisabledService
}

func (n *nginxConfigServiceWithState) configBackup() error {
	if n.status > Disabled {
		return n.service.configBackup()
	}
	return ErrDisabledService
}

func (n *nginxConfigServiceWithState) serverName() string {
	return n.service.serverName()
}

func (n *nginxConfigServiceWithState) DisplayConfig() ([]byte, error) {
	if n.status > Disabled {
		return n.service.DisplayConfig()
	}
	return nil, ErrDisabledService
}

func (n *nginxConfigServiceWithState) GetConfig() ([]byte, error) {
	if n.status > Disabled {
		return n.service.GetConfig()
	}
	return nil, ErrDisabledService
}

func (n *nginxConfigServiceWithState) ShowStatistics() ([]byte, error) {
	if n.status > Disabled {
		return n.service.ShowStatistics()
	}
	return nil, ErrDisabledService
}

func (n *nginxConfigServiceWithState) DisplayWebServerStatus() State {
	if n.status > Disabled {
		return n.service.DisplayWebServerStatus()
	}
	return Unknown
}

func (n *nginxConfigServiceWithState) DisplayWebServerVersion() string {
	if n.status > Disabled {
		return n.service.DisplayWebServerVersion()
	}
	return "unknown"
}

func (n *nginxConfigServiceWithState) UpdateConfig(data []byte) error {
	if n.status > Disabled {
		return n.service.UpdateConfig(data)
	}
	return ErrDisabledService
}

func (n *nginxConfigServiceWithState) WatchLog(logName string) (LogWatcher, error) {
	if n.status > Disabled {
		return n.service.WatchLog(logName)
	}
	return nil, ErrDisabledService
}

func newNginxConfigServiceWithState(info WebServerConfigInfo) *nginxConfigServiceWithState {
	service := newNginxConfigServiceWithLocker(info)
	state := Unknown
	if service == nil {
		state = Abnormal
	}
	return &nginxConfigServiceWithState{
		service: service,
		status:  state,
	}
}
