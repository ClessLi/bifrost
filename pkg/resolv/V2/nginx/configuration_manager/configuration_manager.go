package configuration_manager

import "github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration"

type Backuper interface {
	Backup(backupDir string, retentionTime, backupCycleTime int) error
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

type Reader interface {
	Read() []byte
	ReadJson() []byte
	ReadStatistics() []byte
}

type ConfigManager interface {
	configuration.Queryer
	Reader
	Backuper
	Reloader
	Updater
}

//func NewNginxConfigurationAndManager(serverBin, configAbsPath string) (configuration.Configuration, ConfigManager, error) {
//	cm, err := configuration.NewConfigManager(serverBin, configAbsPath)
//	if err != nil {
//		return nil, nil, err
//	}
//	return cm.GetConfiguration(), cm, nil
//}

func NewNginxConfigurationManager(serverBin, configAbsPath string) (ConfigManager, error) {
	return configuration.NewConfigManager(serverBin, configAbsPath)
}
