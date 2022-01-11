package nginx

import (
	log "github.com/ClessLi/bifrost/pkg/log/v1"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/loader"
	"github.com/marmotedu/errors"
	"sync"
)

// ConfigManagerOptions defines options for nginx configuration and manager.
type ConfigManagerOptions struct {
	ServerName     string
	MainConfigPath string
	ServerBinPath  string
	BackupDir      string
	BackupCycle    int
	BackupSaveTime int
}

type ConfigsManagerOptions struct {
	Options []ConfigManagerOptions
}

func newConfigManager(options ConfigManagerOptions) (configuration.ConfigManager, error) {
	conf, err := configuration.NewConfigurationFromPath(options.MainConfigPath)
	if err != nil {
		return nil, err
	}
	return configuration.NewNginxConfigurationManager(
		loader.NewLoader(),
		conf,
		options.ServerBinPath,
		options.BackupDir,
		options.BackupCycle,
		options.BackupSaveTime,
		new(sync.RWMutex),
	), nil
}

type ConfigsManager interface {
	Start() error
	Stop() error
	GetConfigs() map[string]configuration.Configuration
}

type configsManager struct {
	cms map[string]configuration.ConfigManager
}

func (c *configsManager) Start() error {
	isStarted := make([]string, 0)
	var err error
	defer func() {
		if err != nil {
			for _, servername := range isStarted {
				err := c.cms[servername].Stop()
				if err != nil {
					log.Warnf("failed to stop %s nginx config manager, err: %w", servername, err)
				}
			}
		}
	}()
	for servername, manager := range c.cms {
		err = manager.Start()
		if err != nil {
			return err
		}
		isStarted = append(isStarted, servername)
	}
	return nil
}

func (c *configsManager) Stop() error {
	errs := make([]error, 0)
	for servername, manager := range c.cms {
		err := manager.Stop()
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "failed to stop nginx config manager %s", servername))
		}
	}
	if len(errs) > 0 {
		return errors.NewAggregate(errs)
	}
	return nil
}

func (c *configsManager) GetConfigs() map[string]configuration.Configuration {
	var configs = make(map[string]configuration.Configuration)
	for name, manager := range c.cms {
		configs[name] = manager.GetConfiguration()
	}
	return configs
}

func New(options ConfigsManagerOptions) (ConfigsManager, error) {
	cms := make(map[string]configuration.ConfigManager)
	for _, opts := range options.Options {
		cm, err := newConfigManager(opts)
		if err != nil {
			return nil, err
		}
		cms[opts.ServerName] = cm
	}
	return &configsManager{cms: cms}, nil
}
