package nginx

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"

	"github.com/marmotedu/errors"
)

type Config struct {
	ManagersConfig       map[string]*configuration.ManagerConfig
	DomainNameServerIPv4 string
}

func (c *Config) Complete() (*CompletedConfig, error) {
	var errs []error
	var err error
	cc := &CompletedConfig{ManagersConfig: make(map[string]*configuration.CompletedManagerConfig)}

	for svrname, config := range c.ManagersConfig {
		cc.ManagersConfig[svrname], err = config.Complete()
		errs = append(errs, err)
	}

	cc.DomainNameServerIPv4 = c.DomainNameServerIPv4

	return cc, errors.NewAggregate(errs)
}

type CompletedConfig struct {
	ManagersConfig       map[string]*configuration.CompletedManagerConfig
	DomainNameServerIPv4 string
}

func (cc *CompletedConfig) NewConfigsManager() (ConfigsManager, error) {
	var errs []error
	var err error
	m := &configsManager{make(map[string]configuration.NginxConfigManager)}
	for svrname, config := range cc.ManagersConfig {
		m.managerMap[svrname], err = config.NewNginxConfigManager()
		errs = append(errs, err)
	}

	if aerr := errors.NewAggregate(errs); aerr != nil {
		return nil, aerr
	}

	utilsV3.SetDomainNameResolver(utilsV3.NewDNSClient(cc.DomainNameServerIPv4))

	return m, nil
}
