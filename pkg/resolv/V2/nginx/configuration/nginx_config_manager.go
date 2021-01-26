package configuration

import (
	"errors"
	"fmt"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration/parser"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/loader"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type configManager struct {
	loader         loader.Loader
	configuration  Configuration
	mainConfigPath string
	configPaths    []string
	serverBinPath  string
	rwLocker       *sync.RWMutex
}

func (c configManager) Query(keyword string) (Queryer, error) {
	return c.configuration.Query(keyword)
}

func (c configManager) QueryAll(keyword string) ([]Queryer, error) {
	return c.configuration.QueryAll(keyword)
}

func (c configManager) Self() parser.Parser {
	return nil
}

func (c configManager) fatherContext() parser.Context {
	return nil
}

func (c configManager) index() int {
	return 0
}

func (c configManager) Read() []byte {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	return c.configuration.View()
}

func (c configManager) ReadJson() []byte {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	return c.configuration.Json()
}

func (c configManager) ReadStatistics() []byte {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	return c.configuration.StatisticsByJson()
}

func (c *configManager) UpdateFromJsonBytes(data []byte) error {
	ctx, loopPreventer, err := c.loader.LoadFromJsonBytes(data)
	if err != nil {
		return err
	}
	config, ok := ctx.(*parser.Config)
	if !ok {
		return errors.New("not config json bytes")
	}
	c.rwLocker.Lock()
	defer c.rwLocker.Unlock()
	err = c.configuration.renewConfiguration(NewConfiguration(config, loopPreventer))
	if err != nil {
		return err
	}
	return c.SaveWithCheck()

}

func (c configManager) Backup(backupDir string, retentionTime, backupCycleTime int) error {
	c.rwLocker.Lock()
	defer c.rwLocker.Unlock()
	err := c.SaveWithCheck()
	if err != nil && err != ErrSameConfigFingerprint {
		return err
	}
	//system time
	TZ := time.Local
	// 归档日期初始化
	now := time.Now().In(TZ)
	dt := now.Format("20060102")
	backupName := fmt.Sprintf("nginx.%s.tgz", dt)
	archiveDir, err := filepath.Abs(filepath.Dir(c.configuration.Self().GetValue()))
	archivePath := filepath.Join(archiveDir, backupName)

	// 确认是否为指定归档路径
	var isSpecialBackupDir bool

	if backupDir != "" {
		isSpecialBackupDir = true
		backupDir, err = filepath.Abs(backupDir)
		if err != nil {
			return err
		}
	} else {
		backupDir = archiveDir
	}

	// 判断是否需要备份
	needBackup, err := utils.CheckBackups(backupName, backupDir, retentionTime, backupCycleTime, now)
	if err != nil {
		return err
	}

	if !needBackup {
		return NoBackupRequired
	}

	// 压缩归档
	err = utils.TarGZ(archivePath, c.configPaths)
	if err != nil {
		return err
	}

	// 根据需要调整归档路径
	if isSpecialBackupDir {
		backupPath := filepath.Join(backupDir, backupName)
		return os.Rename(archivePath, backupPath)
	}
	return nil

}

func (c *configManager) Reload() error {
	c.rwLocker.Lock()
	defer c.rwLocker.Unlock()
	config, configPaths, err := c.load()
	if err != nil {
		return err
	}
	if c.configuration != nil {
		err = c.configuration.renewConfiguration(config)
		if err != nil {
			if err == ErrSameConfigFingerprint {
				return NoReloadRequired
			}
			return err
		}
	} else {
		c.configuration = config
	}
	c.configPaths = configPaths
	return nil
}

func (c configManager) load() (conf Configuration, configPaths []string, err error) {
	ctx, loopPreventer, err := c.loader.LoadFromFilePath(c.mainConfigPath)
	if err != nil {
		return nil, nil, err
	}
	config, ok := ctx.(*parser.Config)
	if !ok {
		return nil, nil, errors.New("not a config file")
	}
	configPaths = c.loader.GetConfigPaths()
	return NewConfiguration(config, loopPreventer), configPaths, nil
}

func (c *configManager) SaveWithCheck() error {
	// old configs
	oldConfig, oldConfigPaths, err := c.load()
	if err != nil {
		return err
	}

	// remove old configs
	err = utils.RemoveFiles(oldConfigPaths)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = c.configuration.renewConfiguration(oldConfig)
			err = utils.RemoveFiles(oldConfigPaths)
			err = c.save()
			err = c.Check()
		}
	}()
	err = c.save()
	if err != nil {
		return err
	}

	return c.Check()
}

func (c configManager) save() error {
	dumps := c.configuration.Dump()
	configPaths := make([]string, 0)
	for s, bytes := range dumps {
		err := ioutil.WriteFile(s, bytes, 0755)
		if err != nil {
			return err
		}

		/*// debug test
		fmt.Println(s, ":")
		fmt.Println(string(bytes))
		// debug test end*/

		configPaths = append(configPaths, s)
	}
	c.configPaths = configPaths
	c.mainConfigPath = c.configuration.Self().GetValue()
	return nil
}

func (c configManager) Check() error {
	cmd := exec.Command(c.serverBinPath, "-tc", c.mainConfigPath)
	cmd.Stderr = os.Stderr
	return cmd.Run()

	/*// debug test
	return nil
	// debug test end*/
}

func NewConfigManager(serverBinPath, configAbsPath string) (*configManager, error) {
	cm := &configManager{
		loader:         loader.NewLoader(),
		mainConfigPath: configAbsPath,
		serverBinPath:  serverBinPath,
		//configuration: &configuration{
		//	rwLocker: new(sync.RWMutex),
		//},
		rwLocker: new(sync.RWMutex),
	}
	err := cm.Reload()
	if err != nil {
		return nil, err
	}
	//err = cm.Check()
	//if err != nil {
	//	return nil, err
	//}
	return cm, nil
}
