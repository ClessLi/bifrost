package configuration

import (
	"errors"
	"fmt"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration/parser"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/loader"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/utils"
	"os"
	"path/filepath"
	"sync"
)

type configManager struct {
	loader         loader.Loader
	configuration  Configuration
	mainConfigPath string
	configPaths    []string
	serverBinPath  string
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
	c.configuration = NewConfiguration(config, loopPreventer)
	return nil
}

func (c configManager) Backup(filePath string) error {
	err := c.SaveWithCheck()
	if err != nil {
		return err
	}
	archivePath, err := filepath.Abs(filepath.Join(filepath.Dir(c.configuration.Self().GetValue()), "nginx.conf.tgz"))
	if err != nil {
		return err
	}

	err = utils.TarGZ(archivePath, c.configPaths)
	if err != nil {
		return err
	}
	if filePath != "" {
		filePath, err := filepath.Abs(filePath)
		if err != nil {
			return err
		}
		return os.Rename(archivePath, filePath)
	}
	return nil

}

func (c *configManager) Reload() error {
	config, configPaths, err := c.load()
	if err != nil {
		return err
	}
	c.configuration = config
	//c.mainConfigPath = config.GetValue()
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
			c.configuration = oldConfig
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
		/*err := ioutil.WriteFile(s, bytes, 0755)
		if err != nil {
			return err
		}*/

		// debug test
		fmt.Println(bytes)
		// debug test end

		configPaths = append(configPaths, s)
	}
	c.configPaths = configPaths
	c.mainConfigPath = c.configuration.Self().GetValue()
	return nil
}

func (c configManager) Check() error {
	/*cmd := exec.Command(c.serverBinPath, "-tc", c.mainConfigPath)
	cmd.Stderr = os.Stderr
	return cmd.Run()*/

	// debug test
	return nil
	// debug test end
}

func (c configManager) GetConfiguration() Configuration {
	return c.configuration
}

func NewConfigManager(serverBinPath, configAbsPath string) (*configManager, error) {
	cm := &configManager{
		loader:         loader.NewLoader(),
		mainConfigPath: configAbsPath,
		serverBinPath:  serverBinPath,
		configuration: &configuration{
			rwLocker: new(sync.RWMutex),
		},
	}
	err := cm.Reload()
	if err != nil {
		return nil, err
	}
	err = cm.Check()
	if err != nil {
		return nil, err
	}
	return cm, nil
}
