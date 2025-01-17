package nginx

import (
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	"github.com/marmotedu/errors"
	"os/exec"
	"time"
)

type ConfigsManager interface {
	Start() error
	Stop(timeout time.Duration) error
	GetConfigs() map[string]configuration.NginxConfig
	GetServerInfos() []*v1.WebServerInfo
	GetServersBinCMD() map[string]func(arg ...string) *exec.Cmd
}

type configsManager struct {
	managerMap map[string]configuration.NginxConfigManager
}

var _ ConfigsManager = &configsManager{}

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
	configs = make(map[string]configuration.NginxConfig)
	for svrname, manager := range m.managerMap {
		configs[svrname] = manager.NginxConfig()
	}
	return
}

func (m *configsManager) GetServerInfos() (infos []*v1.WebServerInfo) {
	infos = make([]*v1.WebServerInfo, 0)
	for svrname, manager := range m.managerMap {
		infos = append(infos, &v1.WebServerInfo{
			Name:    svrname,
			Status:  manager.ServerStatus(),
			Version: manager.ServerVersion(),
		})
	}
	return
}

func (m *configsManager) GetServersBinCMD() map[string]func(arg ...string) *exec.Cmd {
	cmds := make(map[string]func(arg ...string) *exec.Cmd)
	for svrname, manager := range m.managerMap {
		cmds[svrname] = func(arg ...string) *exec.Cmd {
			return manager.ServerBinCMD(arg...)
		}
	}
	return cmds
}
