package web_server_manager

import "sync"

type nginxConfigServiceWithLocker struct {
	service WebServerConfigService
	locker  *sync.RWMutex
}

func (n nginxConfigServiceWithLocker) Lock() {
	n.locker.Lock()
}

func (n nginxConfigServiceWithLocker) Unlock() {
	n.locker.Unlock()
}

func (n nginxConfigServiceWithLocker) RLock() {
	n.locker.RLock()
}

func (n nginxConfigServiceWithLocker) RUnlock() {
	n.locker.RUnlock()
}

func (n nginxConfigServiceWithLocker) checkConfigsHash() (bool, error) {
	n.RLock()
	defer n.RUnlock()
	return n.service.checkConfigsHash()
}

func (n nginxConfigServiceWithLocker) configLoad() error {
	n.RLock()
	defer n.RUnlock()
	return n.service.configLoad()
}

func (n nginxConfigServiceWithLocker) configBackup() error {
	n.RLock()
	defer n.RUnlock()
	return n.service.configBackup()
}

func (n nginxConfigServiceWithLocker) serverName() string {
	n.RLock()
	defer n.RUnlock()
	return n.service.serverName()
}

func (n nginxConfigServiceWithLocker) DisplayConfig() ([]byte, error) {
	n.RLock()
	defer n.RUnlock()
	return n.service.DisplayConfig()
}

func (n nginxConfigServiceWithLocker) GetConfig() ([]byte, error) {
	n.RLock()
	defer n.RUnlock()
	return n.service.GetConfig()
}

func (n nginxConfigServiceWithLocker) ShowStatistics() ([]byte, error) {
	n.RLock()
	defer n.RUnlock()
	return n.service.ShowStatistics()
}

func (n nginxConfigServiceWithLocker) DisplayWebServerStatus() State {
	n.RLock()
	defer n.RUnlock()
	return n.service.DisplayWebServerStatus()
}

func (n nginxConfigServiceWithLocker) DisplayWebServerVersion() string {
	n.RLock()
	defer n.RUnlock()
	return n.service.DisplayWebServerVersion()
}

func (n nginxConfigServiceWithLocker) UpdateConfig(data []byte) error {
	n.Lock()
	defer n.Unlock()
	return n.service.UpdateConfig(data)
}

func (n nginxConfigServiceWithLocker) WatchLog(logName string) (LogWatcher, error) {
	n.RLock()
	defer n.RUnlock()
	return n.service.WatchLog(logName)
}

func newNginxConfigServiceWithLocker(info WebServerConfigInfo) *nginxConfigServiceWithLocker {
	return &nginxConfigServiceWithLocker{
		service: newNginxConfigService(info),
		locker:  new(sync.RWMutex),
	}
}
