package configuration

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/loader"
	"sync"
)

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
	Queryer
	Reader
	Backuper
	Reloader
	Updater
}

func NewNginxConfigurationManager(serverBinPath, configAbsPath string) (ConfigManager, error) {
	cm := &configManager{
		loader:         loader.NewLoader(),
		mainConfigPath: configAbsPath,
		serverBinPath:  serverBinPath,
		rwLocker:       new(sync.RWMutex),
	}
	err := cm.Reload()
	if err != nil {
		return nil, err
	}
	return cm, nil
}
