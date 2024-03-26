package nginx

import (
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	"github.com/marmotedu/errors"
	"time"
)

type ConfigsManager interface {
	Start() error
	Stop(timeout time.Duration) error
	GetConfigs() map[string]configuration.NginxConfig
	GetServerInfos() []*v1.WebServerInfo
}

type configsManager struct {
	managerMap map[string]configuration.NginxConfigManager
}

func (m *configsManager) Start() error {
	var errs []error
	for _, manager := range m.managerMap {
		errs = append(errs, manager.Start())
	}
	return errors.NewAggregate(errs)
}

func (m *configsManager) Stop(timeout time.Duration) error {
	var errs []error
	for _, manager := range m.managerMap {
		errs = append(errs, manager.Stop(timeout))
	}
	return errors.NewAggregate(errs)
}

func (m *configsManager) GetConfigs() (configs map[string]configuration.NginxConfig) {
	for svrname, manager := range m.managerMap {
		configs[svrname] = manager.NginxConfig()
	}
	return
}

func (m *configsManager) GetServerInfos() (infos []*v1.WebServerInfo) {
	for svrname, manager := range m.managerMap {
		infos = append(infos, &v1.WebServerInfo{
			Name:    svrname,
			Status:  manager.ServerStatus(),
			Version: manager.ServerVersion(),
		})
	}
	return
}
