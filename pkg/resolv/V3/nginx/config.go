package nginx

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	"github.com/marmotedu/errors"
)

type Config struct {
	ManagersConfig map[string]*configuration.ManagerConfig
}

func (c *Config) Complete() (*CompletedConfig, error) {
	var errs []error
	var err error
	cc := &CompletedConfig{make(map[string]*configuration.CompletedManagerConfig)}

	for svrname, config := range c.ManagersConfig {
		cc.ManagersConfig[svrname], err = config.Complete()
		errs = append(errs, err)
	}
	return cc, errors.NewAggregate(errs)
}

type CompletedConfig struct {
	ManagersConfig map[string]*configuration.CompletedManagerConfig
}

func (cc *CompletedConfig) NewConfigsManager() (ConfigsManager, error) {
	var errs []error
	var err error
	m := &configsManager{make(map[string]configuration.NginxConfigManager)}
	for svrname, config := range cc.ManagersConfig {
		m.managerMap[svrname], err = config.NewNginxConfigManager()
		errs = append(errs, err)
	}
	return m, errors.NewAggregate(errs)
}
