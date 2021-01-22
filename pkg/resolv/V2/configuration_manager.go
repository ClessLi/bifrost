package V2

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration"
)

type Backuper interface {
	Backup(filePath string) error
}

type Reloader interface {
	Reload() error
}

type Saver interface {
	SaveWithCheck() error
	Check() error
}

type Updater interface {
	UpdateFromJsonBytes(data []byte) error
}

type ConfigManager interface {
	Backuper
	Reloader
	Saver
	Updater
}

func NewNginxConfigurationAndManager(serverBin, configAbsPath string) (configuration.Configuration, ConfigManager, error) {
	cm, err := configuration.NewConfigManager(serverBin, configAbsPath)
	if err != nil {
		return nil, nil, err
	}
	return cm.GetConfiguration(), cm, nil
}
