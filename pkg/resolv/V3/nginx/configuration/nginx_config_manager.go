package configuration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	utilsV2 "github.com/ClessLi/bifrost/pkg/resolv/V2/utils"
	nginxContext "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"

	logV1 "github.com/ClessLi/component-base/pkg/log/v1"

	"github.com/marmotedu/errors"
	"golang.org/x/sync/errgroup"
)

type NginxConfigManager interface {
	Start() error
	Stop(timeout time.Duration) error
	NginxConfig() NginxConfig
	ServerStatus() v1.State
	ServerVersion() string
	ServerBinCMD(arg ...string) *exec.Cmd
}

type nginxConfigManager struct {
	configuration           NginxConfig
	nginxHome               string
	nginxBinFilePath        string
	regularlyTaskCycleDelay time.Duration
	ctx                     context.Context
	cancel                  context.CancelFunc
	backupOpts              backupOption
	wg                      *sync.WaitGroup
	eg                      *errgroup.Group
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
	// TODO: Server Port listen state, proxy server connection state
	if m.isRunning || (m.ctx != nil && m.ctx.Err() == nil) {
		return errors.WithCode(code.ErrConfigManagerIsRunning, "the nginx config manager is already running")
	}

	if m.ctx == nil || m.cancel == nil {
		m.ctx, m.cancel = context.WithCancel(context.Background())
	}
	var work context.Context
	m.eg, work = errgroup.WithContext(m.ctx)
	// cron jobs
	m.eg.Go(func() error {
		defer func() {
			m.isRunning = false
		}()

		return m.regularlyTask(work)
	})

	return m.ctx.Err()
}

func (m *nginxConfigManager) Stop(timeout time.Duration) error {
	if !m.isRunning {
		return errors.WithCode(code.ErrConfigManagerIsNotRunning, "the nginx config manager is not running")
	}
	stopWork, done := context.WithTimeout(context.TODO(), timeout)
	defer done()
	stopeg, _ := errgroup.WithContext(stopWork)
	stopeg.Go(func() error {
		defer done()
		m.cancel()

		return m.eg.Wait()
	})

	if err := stopeg.Wait(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errors.Errorf("stop the nginx config manager time out for more than %d seconds", int(timeout/time.Second))
		} else if !errors.Is(err, context.Canceled) {
			return errors.Errorf("stop the nginx config manager error. %+v", err)
		}
		m.isRunning = false

		return nil
	}

	return errors.WithCode(code.ErrConfigManagerIsRunning, "failed to stop the nginx config manager")
}

func (m *nginxConfigManager) NginxConfig() NginxConfig {
	return m.configuration
}

func (m *nginxConfigManager) ServerStatus() (state v1.State) {
	state = v1.UnknownState
	svrPidFilePath := filepath.Join("logs", "nginx.pid")
	pidCtx := m.configuration.Main().
		ChildrenPosSet().
		QueryOne(nginxContext.NewKeyWordsByType(context_type.TypeDirective).
			SetSkipQueryFilter(nginxContext.SkipDisabledCtxFilterFunc).
			SetRegexpMatchingValue(`pid\s+.*`)).
		Target()
	if pidCtx.Error() == nil {
		pidCtxKV := strings.Split(pidCtx.Value(), " ")
		if len(pidCtxKV) == 2 {
			svrPidFilePath = pidCtxKV[1]
		}
	}
	if !filepath.IsAbs(svrPidFilePath) {
		nginxHomeAbsDir, err := filepath.Abs(m.nginxHome)
		if err != nil {
			return state
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
	cmd := m.ServerBinCMD("-v")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}

	return strings.TrimRight(string(output), "\n")
}

func (m *nginxConfigManager) ServerBinCMD(arg ...string) *exec.Cmd {
	arg = append(arg, "-c", m.configuration.Main().MainConfig().FullPath())

	return exec.Command(m.nginxBinFilePath, arg...)
}

func (m *nginxConfigManager) backup() error {
	logV1.Debugf("start to backup nginx configs")
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
	logV1.Debug("start to refresh nginx configs")
	fsMain, fsFingerprinter, err := m.load()
	if err != nil {
		return err
	}

	if !fsFingerprinter.Diff(utilsV3.NewConfigFingerprinter(m.configuration.Dump()).Fingerprints()) {
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

func (m *nginxConfigManager) reparseProxyInfo() error {
	logV1.Debugf("start to reparse proxy info")

	return m.configuration.Main().ChildrenPosSet().QueryAll(
		nginxContext.NewKeyWords(func(targetCtx nginxContext.Context) bool {
			return targetCtx.Type() == context_type.TypeDirHTTPProxyPass ||
				targetCtx.Type() == context_type.TypeDirStreamProxyPass ||
				targetCtx.Type() == context_type.TypeDirUninitializedProxyPass
		}).
			SetSkipQueryFilter(nginxContext.SkipDisabledCtxFilterFunc), // skip reparsing the disabled `ProxyPass`.
	).
		Map(func(pos nginxContext.Pos) (nginxContext.Pos, error) {
			if err := pos.Target().Error(); err != nil {
				return pos, err
			}
			pp, ok := pos.Target().(local.ProxyPass)
			if !ok {
				return pos, fmt.Errorf("invalid ProxyPass: %s", pos.Target().Value())
			}

			return pos, pp.ReparseParams()
		}).Error()
}

func (m *nginxConfigManager) regularlyTask(ctx context.Context) error {
	// regularly backup is disabled, when m.backupOpts.backupCycleDays or m.backupOpts.backupRetentionDays is less equal zero.
	backupIsDisabled := m.backupOpts.backupCycleDays <= 0 || m.backupOpts.backupRetentionDays <= 0

	ticker := time.NewTicker(m.regularlyTaskCycleDelay)
	for {
		// 等待触发
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return ctx.Err()
		}

		taskeg, _ := errgroup.WithContext(ctx)
		// refresh configs
		err := m.refresh()
		if err != nil {
			logV1.Errorf("refresh nginx config failed, cased by: %+v", err)
		}

		// backup configs after refreshing
		if !backupIsDisabled {
			taskeg.Go(func() error {
				bakerr := m.backup()
				if bakerr != nil {
					logV1.Errorf("backup nginx config failed, cased by: %+v", bakerr)
				}

				return nil
			})
		}

		// proxy information parsing task after refreshing
		taskeg.Go(func() error {
			parseerr := m.reparseProxyInfo()
			if parseerr != nil {
				logV1.Errorf("reparse proxy info failed, cased by: %+v", parseerr)
			}

			return nil
		})

		err = taskeg.Wait()
		if err != nil && !errors.Is(err, context.Canceled) {
			logV1.Errorf("regularly task executed failed, cased by: %+v", err)
		}
	}
}

func (m *nginxConfigManager) save() error {
	dumps := m.configuration.Dump()
	for file, dumpbuff := range dumps {
		err := os.WriteFile(file, dumpbuff.Bytes(), 0o600)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *nginxConfigManager) check() error {
	cmd := m.ServerBinCMD("-t")
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
