package configuration

import (
	"bytes"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	utilsV2 "github.com/ClessLi/bifrost/pkg/resolv/V2/utils"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"github.com/marmotedu/errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type NginxConfigManager interface {
	Start() error
	Stop(timeout time.Duration) error
	NginxConfig() NginxConfig
	ServerStatus() v1.State
	ServerVersion() string
}

type nginxConfigManager struct {
	configuration           NginxConfig
	nginxHome               string
	nginxBinFilePath        string
	regularlyTaskCycleDelay time.Duration
	regularlyTaskSignalChan chan int
	backupOpts              backupOption
	wg                      *sync.WaitGroup
	isRunning               bool
}

type backupOption struct {
	backupCycleDays     int
	backupRetentionDays int
	backupDir           string
	backupPrefix        string
	backupTimeZone      *time.Location
}

func (m *nginxConfigManager) Start() error {
	if m.isRunning {
		return errors.WithCode(code.ErrConfigManagerIsRunning, "nginx config manager is already running")
	}
	m.wg.Add(1)
	go func() {
		defer func() {
			m.isRunning = false
		}()
		defer m.wg.Done()
		err := m.regularlyRefreshAndBackup(m.regularlyTaskSignalChan)
		if err != nil {
			logV1.Errorf("regularly refresh and backup task start error. %+v", err)
		}
	}()
	return nil
}

func (m *nginxConfigManager) Stop(timeout time.Duration) error {
	if !m.isRunning {
		return errors.WithCode(code.ErrConfigManagerIsNotRunning, "nginx config manager is not running")
	}
	defer m.wg.Wait()
	select {
	case <-time.After(timeout):
		return errors.Errorf("stop nginx config manager time out for more than %d seconds", int(timeout/time.Second))
	case m.regularlyTaskSignalChan <- 9:
		m.isRunning = false
		return nil
	}
}

func (m *nginxConfigManager) NginxConfig() NginxConfig {
	return m.configuration
}

func (m *nginxConfigManager) ServerStatus() (state v1.State) {
	state = v1.UnknownState
	svrPidFilePath := filepath.Join("logs", "nginx.pid")
	pidCtx := m.configuration.Main().QueryByKeyWords(context.NewKeyWords(context_type.TypeDirective).SetRegexpMatchingValue(`pid\s+.*`)).Target()
	if pidCtx.Error() == nil {
		pidCtxKV := strings.Split(pidCtx.Value(), " ")
		if len(pidCtxKV) == 2 {
			svrPidFilePath = pidCtxKV[1]
		}
	}
	if !filepath.IsAbs(svrPidFilePath) {
		nginxHomeAbsDir, err := filepath.Abs(m.nginxHome)
		if err != nil {
			return
		}
		svrPidFilePath = filepath.Join(nginxHomeAbsDir, svrPidFilePath)
	}

	pid, err := utilsV2.GetPid(svrPidFilePath)
	if err != nil {
		return v1.Abnormal
	}
	_, err = os.FindProcess(pid)
	if err != nil {
		return v1.Abnormal
	}
	return v1.Normal
}

func (m *nginxConfigManager) ServerVersion() (version string) {
	version = "unknown"
	cmd := m.serverBinCMD("-v")
	stdoutPipe, err := cmd.StderrPipe()
	if err != nil {
		return
	}
	err = cmd.Run()
	if err != nil {
		return
	}
	buff := bytes.NewBuffer([]byte{})
	_, err = buff.ReadFrom(stdoutPipe)
	if err != nil {
		return
	}
	return strings.TrimRight(buff.String(), "\n")
}

func (m *nginxConfigManager) backup() error {
	// 开始备份
	// 归档日期初始化
	now := time.Now().In(m.backupOpts.backupTimeZone)
	backupName := utilsV2.GetBackupFileName(m.backupOpts.backupPrefix, now)
	archiveDir, err := filepath.Abs(m.NginxConfig().Main().MainConfig().BaseDir())
	if err != nil {
		return errors.Wrap(err, "failed to format archive directory")
	}
	archivePath := filepath.Join(archiveDir, backupName)

	// 确认是否为指定归档路径
	var isSpecialBackupDir bool
	backupDir := m.backupOpts.backupDir
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
	needBackup, err := utilsV2.CheckAndCleanBackups(m.backupOpts.backupPrefix, backupDir, m.backupOpts.backupRetentionDays, m.backupOpts.backupCycleDays, now)
	if err != nil {
		logV1.Warn("failed to check and clean backups, " + err.Error())
		return err
	}

	if !needBackup {
		return nil
	}

	var configPaths []string
	for _, config := range m.NginxConfig().Main().ListConfigs() {
		configPaths = append(configPaths, config.FullPath())
	}

	// 压缩归档
	logV1.Info("start backup configs")
	err = utilsV2.TarGZ(archivePath, configPaths)
	if err != nil {
		logV1.Warn("failed to backup configs, " + err.Error())
		return err
	}

	// 根据需要调整归档路径
	if isSpecialBackupDir {
		backupPath := filepath.Join(backupDir, backupName)
		err = os.Rename(archivePath, backupPath)
	}
	logV1.Info("complete configs backup")
	return err
}

func (m *nginxConfigManager) refresh() error {
	fsMain, fsFingerprinter, err := m.load()
	if err != nil {
		return err
	}

	if !fsFingerprinter.Diff(utilsV3.NewConfigFingerprinter(m.configuration.Dump())) {
		// 指纹一致不做刷新
		return nil
	}

	if fsFingerprinter.NewerThan(m.configuration.UpdatedTimestamp()) {
		err := m.configuration.(*nginxConfig).renewMainContext(fsMain)
		if err != nil && !errors.IsCode(err, code.ErrSameConfigFingerprint) {
			return err
		}
	} else {
		// TODO: 配置文件删除功能独立出来
		// // 仅清理本地磁盘上拓扑图中的配置
		// var fsConfigPaths []string
		// for _, config := range fsMain.Topology() {
		// 	fsConfigPaths = append(fsConfigPaths, config.FullPath())
		// }
		// err = utilsV2.RemoveFiles(fsConfigPaths)
		// if err != nil {
		// 	return err
		// }

		// 回写内存配置到本地磁盘
		err = m.saveWithCheck()
		if err != nil { // 保存异常，则回退
			// TODO: 配置文件删除功能独立出来
			// // 仅清理内存中在拓扑图中的配置
			// var memConfigPaths []string
			// for _, config := range m.configuration.Main().Topology() {
			// 	memConfigPaths = append(memConfigPaths, config.FullPath())
			// }

			// err = utilsV2.RemoveFiles(memConfigPaths)
			// if err != nil {
			// 	logV1.Errorf("remove failure update configs failed. %+v", err)
			// }

			// 回退为原本地磁盘上配置
			err = m.configuration.(*nginxConfig).renewMainContext(fsMain)
			if err != nil && !errors.IsCode(err, code.ErrSameConfigFingerprint) {
				logV1.Errorf("roll back nginx config failed. %+v", err)
				return err
			}
			return m.saveWithCheck()
		}
	}
	return nil
}

func (m *nginxConfigManager) regularlyRefreshAndBackup(signalChan chan int) error {
	// regularly backup is disabled, when c.backupCycleDays or c.backupRetentionDays is less equal zero.
	backupIsDisabled := m.backupOpts.backupCycleDays <= 0 || m.backupOpts.backupRetentionDays <= 0

	ticker := time.NewTicker(m.regularlyTaskCycleDelay)
	for {
		// 等待触发
		select {
		case <-ticker.C:
		case signal := <-signalChan:
			if signal == 9 {
				return nil
			}
		}
		err := m.refresh()
		if err != nil {
			logV1.Errorf("refresh nginx config failed, cased by: %+v", err)
		}

		if !backupIsDisabled {
			err = m.backup()
			if err != nil {
				logV1.Errorf("backup nginx config failed, cased by: %+v", err)
			}
		}
	}
}

func (m *nginxConfigManager) serverBinCMD(arg ...string) *exec.Cmd {
	arg = append(arg, "-c", m.configuration.Main().MainConfig().FullPath())
	return exec.Command(m.nginxBinFilePath, arg...)
}

func (m *nginxConfigManager) save() error {
	dumps := m.configuration.Dump()
	for file, dumpbuff := range dumps {
		err := os.WriteFile(file, dumpbuff.Bytes(), 0600)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *nginxConfigManager) check() error {
	cmd := m.serverBinCMD("-t")
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (m *nginxConfigManager) saveWithCheck() error {
	err := m.save()
	if err != nil {
		return err
	}
	return m.check()
}

func (m *nginxConfigManager) load() (local.MainContext, utilsV3.ConfigFingerprinter, error) {
	localMain, timestamp, err := loadMainContextFromFS(m.configuration.Main().Value())
	if err != nil {
		return nil, nil, err
	}
	fingerprinter := utilsV3.NewConfigFingerprinterWithTimestamp(dumpMainContext(localMain), timestamp)
	return localMain, fingerprinter, nil
}
