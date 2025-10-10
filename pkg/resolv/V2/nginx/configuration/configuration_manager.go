package configuration

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration/parser"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/loader"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/utils"

	logV1 "github.com/ClessLi/component-base/pkg/log/v1"

	"github.com/marmotedu/errors"
	"github.com/wxnacy/wgo/arrays"
)

type ConfigManager interface {
	Start() error
	Stop() error
	GetConfiguration() Configuration
	regularlyBackup(duration time.Duration, signalChan chan int) error
	regularlyReload(duration time.Duration, signalChan chan int) error
	regularlySave(duration time.Duration, signalChan chan int) error
	GetServerInfo() *v1.WebServerInfo
}

type configManager struct {
	loader                 loader.Loader
	configuration          Configuration
	configFilesFingerprint utils.ConfigFingerprinter
	mainConfigPath         string
	configPaths            []string
	backupCycle            int
	backupSaveTime         int
	backupDir              string
	serverBinPath          string
	rwLocker               *sync.RWMutex
	backupSignalChan       chan int
	reloadSignalChan       chan int
	saveSignalChan         chan int
	isRunning              bool
	waitGroup              *sync.WaitGroup
}

// GetServerInfo with web server status and version, but no server name.
func (c *configManager) GetServerInfo() *v1.WebServerInfo {
	return &v1.WebServerInfo{
		Name:    "unknown",
		Status:  c.serverStatus(),
		Version: c.serverVersion(),
	}
}

func (c configManager) serverVersion() (version string) {
	version = "unknown"
	cmd := c.serverBinCMD("-v")
	stdoutPipe, err := cmd.StderrPipe()
	if err != nil {
		return
	}
	err = cmd.Run()
	if err != nil {
		return
	}
	buf := bytes.NewBuffer([]byte{})
	_, err = buf.ReadFrom(stdoutPipe)
	if err != nil {
		return
	}

	return strings.TrimRight(buf.String(), "\n")
}

func (c configManager) serverStatus() (status v1.State) {
	status = v1.UnknownState
	svrPidFilePath := "logs/nginx.pid"
	svrPidQueryer, err := c.configuration.Query("key:sep: :reg:pid .*")
	if err == nil {
		svrPidFilePath = strings.Split(svrPidQueryer.Self().GetValue(), " ")[1]
	}

	svrPidFilePathAbs := svrPidFilePath
	if !filepath.IsAbs(svrPidFilePath) {
		svrBinAbs, absErr := filepath.Abs(c.serverBinPath)
		if absErr != nil {
			return status
		}
		svrWS, wsErr := filepath.Abs(filepath.Join(filepath.Dir(svrBinAbs), ".."))
		if wsErr != nil {
			return status
		}
		var pidErr error
		svrPidFilePathAbs, pidErr = filepath.Abs(filepath.Join(svrWS, svrPidFilePath))
		if pidErr != nil {
			return status
		}
	}

	svrPid, gPidErr := utils.GetPid(svrPidFilePathAbs)
	if gPidErr != nil {
		return v1.Abnormal
	}

	_, procErr := os.FindProcess(svrPid)
	if procErr != nil {
		return v1.Abnormal
	}

	return v1.Normal
}

func (c *configManager) GetConfiguration() Configuration {
	return c.configuration
}

func (c *configManager) regularlyBackup(duration time.Duration, signalChan chan int) error {
	// regularly backup is disabled, when c.backupCycle or c.backupSaveTime is less equal zero.
	if c.backupCycle <= 0 || c.backupSaveTime <= 0 {
		return nil
	}

	ticker := time.NewTicker(duration)
	var backupErr error
	for backupErr == nil {
		// 等待触发
		select {
		case <-ticker.C:
		case signal := <-signalChan:
			if signal == 9 {
				return backupErr
			}
		}

		// 1) save with check
		err := c.SaveWithCheck()
		// 非指纹相同的报错则退出备份
		if err != nil && !errors.IsCode(err, code.ErrSameConfigFingerprint) &&
			!errors.IsCode(err, code.ErrSameConfigFingerprints) {
			backupErr = err

			continue
		}
		// 2) 开始备份
		// system time
		TZ := time.Local
		// 归档日期初始化
		now := time.Now().In(TZ)
		backupPrefix := "nginx.conf"
		backupName := utils.GetBackupFileName(backupPrefix, now)
		archiveDir, err := filepath.Abs(filepath.Dir(c.configuration.Self().GetValue()))
		if err != nil {
			backupErr = errors.Wrap(err, "failed to format archive directory")

			continue
		}
		archivePath := filepath.Join(archiveDir, backupName)

		// 确认是否为指定归档路径
		var isSpecialBackupDir bool

		if c.backupDir != "" {
			isSpecialBackupDir = true
			c.backupDir, err = filepath.Abs(c.backupDir)
			if err != nil {
				return err
			}
		} else {
			c.backupDir = archiveDir
		}

		// 初始化时，完成该操作↑

		// 判断是否需要备份
		needBackup, err := utils.CheckAndCleanBackups(backupPrefix, c.backupDir, c.backupSaveTime, c.backupCycle, now)
		if err != nil {
			logV1.Warn("failed to check and clean backups, " + err.Error())
			backupErr = err

			continue
		}

		if !needBackup {
			continue
		}

		// 压缩归档
		c.rwLocker.RLock()
		logV1.Info("start backup configs")
		err = utils.TarGZ(archivePath, c.configPaths)
		c.rwLocker.RUnlock()
		if err != nil {
			logV1.Warn("failed to backup configs, " + err.Error())
			backupErr = err

			continue
		}

		// 根据需要调整归档路径
		if isSpecialBackupDir {
			backupPath := filepath.Join(c.backupDir, backupName)
			backupErr = os.Rename(archivePath, backupPath)
		}
		logV1.Info("complete configs backup")
	}

	return backupErr
}

func (c *configManager) regularlyReload(duration time.Duration, signalChan chan int) error {
	ticker := time.NewTicker(duration)
	var reloadErr error
	for reloadErr == nil {
		// 等待触发
		select {
		case <-ticker.C:
		case signal := <-signalChan:
			if signal == 9 {
				return reloadErr
			}
		}

		// 1) load
		config, configPaths, err := c.load()
		if err != nil {
			reloadErr = err

			continue
		}

		// 2) 判断manager配置指纹与文件指纹是否一致
		if !c.configFilesFingerprint.Diff(config.getConfigFingerprinter()) {
			continue
		}
		// 2) 不一致则重载文件配置
		err = c.configuration.renewConfiguration(config)
		if err != nil {
			if !errors.IsCode(err, code.ErrSameConfigFingerprint) {
				reloadErr = err
			}

			continue
		}
		c.rwLocker.Lock()
		c.configPaths = configPaths
		c.configFilesFingerprint.Renew(config.getConfigFingerprinter())
		c.rwLocker.Unlock()
	}

	return reloadErr
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
	newConfiguration := NewConfiguration(config, loopPreventer, new(sync.RWMutex))

	return newConfiguration, configPaths, nil
}

func (c *configManager) regularlySave(duration time.Duration, signalChan chan int) error {
	ticker := time.NewTicker(duration)
	var saveErr error
	for saveErr == nil {
		// 等待触发
		select {
		case <-ticker.C:
		case signal := <-signalChan:
			if signal == 9 {
				return saveErr
			}
		}
		// 1) 判断manager配置指纹与内存配置指纹是否一致
		if !c.configuration.getConfigFingerprinter().Diff(c.configFilesFingerprint) {
			continue
		}
		// 1) 不一致则save with check
		saveErr = c.SaveWithCheck()
	}

	return saveErr
}

func (c *configManager) SaveWithCheck() error {
	// 1) load old configs
	oldConfig, oldConfigPaths, err := c.load()
	if err != nil {
		return err
	}

	// 2) 比较内存配置指纹与old配置指纹是否一致
	if !c.configuration.getConfigFingerprinter().Diff(oldConfig.getConfigFingerprinter()) {
		return errors.WithCode(
			code.ErrSameConfigFingerprints,
			"same config fingerprint between files and configuration",
		)
	}

	// 2) 不一致则save内存配置
	// remove old configs
	err = utils.RemoveFiles(oldConfigPaths)
	if err != nil {
		return err
	}

	configPaths, err := c.save()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			// 3) check失败则将old配置写入内存和写入本地文件，更新manager配置指纹为old配置指纹
			err = c.configuration.renewConfiguration(oldConfig)
			c.configFilesFingerprint.Renew(oldConfig.getConfigFingerprinter())
			err = utils.RemoveFiles(configPaths)
			configPaths, err = c.save()
			err = c.Check()
		}
		c.rwLocker.Lock()
		defer c.rwLocker.Unlock()
		// 3) check成功则更新manager配置指纹为内存配置指纹
		c.configFilesFingerprint.Renew(c.configuration.getConfigFingerprinter())
		c.configPaths = configPaths
		c.mainConfigPath = c.configuration.getMainConfigPath()
	}()

	// 3) check
	return c.Check()
}

func (c configManager) save() ([]string, error) {
	dumps := c.configuration.Dump()
	configPaths := make([]string, 0)
	for s, bytes := range dumps {
		// 判断是否为单元测试
		if len(os.Args) > 3 && os.Args[1] == "-test.v" && os.Args[2] == "-test.run" {
			logV1.Infof("%s: %s", s, string(bytes))

			continue
		}
		err := ioutil.WriteFile(s, bytes, 0o600)
		if err != nil {
			return nil, err
		}

		/*// debug test
		fmt.Println(s, ":")
		fmt.Println(string(bytes))
		// debug test end*/

		configPaths = append(configPaths, s)
	}

	return configPaths, nil
}

func (c configManager) Check() error {
	// 判断是否为单元测试
	if arrays.ContainsString(os.Args, "-test.v") >= 0 && arrays.ContainsString(os.Args, "-test.run") >= 0 {
		return nil
	}
	// cmd := exec.Command(c.serverBinPath, "-tc", c.mainConfigPath)
	cmd := c.serverBinCMD("-t")
	cmd.Stderr = os.Stderr

	return cmd.Run()

	/*// debug test
	return nil
	// debug test end*/
}

func (c configManager) serverBinCMD(arg ...string) *exec.Cmd {
	arg = append(arg, "-c", c.mainConfigPath)

	return exec.Command(c.serverBinPath, arg...) //nolint:gosec
}

func (c *configManager) Start() error {
	if c.isRunning {
		return errors.WithCode(code.ErrConfigManagerIsRunning, "config manager is already running")
	}
	c.waitGroup.Add(3)
	go func() {
		err := c.regularlyBackup(time.Minute*5, c.backupSignalChan)
		if err != nil {
			logV1.Errorf("regularly backup error. %+v", err)
		}
	}()
	go func() {
		err := c.regularlyReload(time.Second*30, c.reloadSignalChan)
		if err != nil {
			logV1.Errorf("regularly reload error. %+v", err)
		}
	}()
	go func() {
		err := c.regularlySave(time.Second*10, c.saveSignalChan)
		if err != nil {
			logV1.Errorf("regularly save error. %+v", err)
		}
	}()
	c.isRunning = true

	return nil
}

func (c *configManager) Stop() error {
	if !c.isRunning {
		return errors.WithCode(code.ErrConfigManagerIsNotRunning, "config manager is not running")
	}
	errorStr := ""
	waittime := time.Second * 2
	stopGoroutineFunc := func(goroutineName string, signalChan chan int, timeout <-chan time.Time) {
		select {
		case <-timeout:
			if errorStr != "" {
				errorStr += ", "
			}
			errorStr += fmt.Sprintf(
				"stop goroutine %s timed out for more than %d",
				goroutineName,
				int(waittime/time.Second),
			)

			break
		case signalChan <- 9:
			break
		}
		c.waitGroup.Done()
	}
	go stopGoroutineFunc("backup", c.backupSignalChan, time.After(waittime))
	go stopGoroutineFunc("reload", c.reloadSignalChan, time.After(waittime))
	go stopGoroutineFunc("save", c.saveSignalChan, time.After(waittime))
	c.waitGroup.Wait()
	if len(errorStr) > 0 {
		return errors.New(errorStr)
	}
	c.isRunning = false

	return nil
}

func NewNginxConfigurationManager(
	loader loader.Loader,
	configuration Configuration,
	serverBinPath, backupDir string,
	backupCycle, backupSaveTime int,
	rwLocker *sync.RWMutex,
) ConfigManager {
	fingerprinter := utils.NewConfigFingerprinter(make(map[string][]byte))
	fingerprinter.Renew(configuration.getConfigFingerprinter())
	cm := &configManager{
		loader:                 loader,
		configuration:          configuration,
		configFilesFingerprint: fingerprinter,
		mainConfigPath:         configuration.getMainConfigPath(),
		configPaths:            make([]string, 0),
		serverBinPath:          serverBinPath,
		backupDir:              backupDir,
		backupCycle:            backupCycle,
		backupSaveTime:         backupSaveTime,
		rwLocker:               rwLocker,
		backupSignalChan:       make(chan int),
		reloadSignalChan:       make(chan int),
		saveSignalChan:         make(chan int),
		waitGroup:              new(sync.WaitGroup),
	}
	for s := range cm.configuration.Dump() {
		cm.configPaths = append(cm.configPaths, s)
	}

	return cm
}
