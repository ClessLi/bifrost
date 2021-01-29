package web_server_manager

import (
	"bytes"
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	ngLog "github.com/ClessLi/bifrost/pkg/log/nginx"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type nginxConfigService struct {
	name               string
	backupCycle        int
	backupSaveTime     int
	backupDir          string
	confPath           string
	verifyExecPath     string
	nginxConfigManager configuration.ConfigManager
	nginxLog           *ngLog.Log
}

func (n nginxConfigService) serverName() string {
	return n.name
}

func (n nginxConfigService) DisplayWebServerStatus() State {
	svrPidFilePath := "logs/nginx.pid"
	svrPidQueryer, err := n.nginxConfigManager.Query("key:sep: :reg:pid .*")
	if err == nil {
		svrPidFilePath = strings.Split(svrPidQueryer.Self().GetValue(), " ")[1]
	}

	svrPidFilePathAbs := svrPidFilePath
	if !filepath.IsAbs(svrPidFilePath) {
		svrBinAbs, absErr := filepath.Abs(n.verifyExecPath)
		if absErr != nil {
			//Log(WARN, "[%s] get web server bin dir err: %s", b.ServiceInfos[i].Name, absErr)
			return Unknown
		}
		svrWS, wsErr := filepath.Abs(filepath.Join(filepath.Dir(svrBinAbs), ".."))
		if wsErr != nil {
			//Log(WARN, "[%s] get web server workspace err: %s", b.ServiceInfos[i].Name, wsErr)
			return Unknown
		}
		var pidErr error
		svrPidFilePathAbs, pidErr = filepath.Abs(filepath.Join(svrWS, svrPidFilePath))
		if pidErr != nil {
			//if SysInfo.StatusList[i].ServerStatus != "unknow" {
			//Log(WARN, "[%s] get web server pid file path failed: %s", b.ServiceInfos[i].Name, pidErr)
			//}
			return Unknown
		}
	}

	svrPid, gPidErr := utils.GetPid(svrPidFilePathAbs)
	if gPidErr != nil {
		utils.Logger.InfoF("[%s] %s", n.serverName(), gPidErr)
		return Abnormal
	}

	_, procErr := os.FindProcess(svrPid)
	if procErr != nil {
		utils.Logger.InfoF("[%s] %s", n.serverName(), procErr)
		return Abnormal
	}
	return Normal
}

func (n nginxConfigService) DisplayWebServerVersion() string {
	svrBinAbs, absErr := filepath.Abs(n.verifyExecPath)
	if absErr != nil {
		return "The absolute path of nginx process binary file is abnormal"
	}
	svrVersion, vErr := func() (string, error) {
		cmd := exec.Command(svrBinAbs, "-v")
		stdoutPipe, pipeErr := cmd.StderrPipe()
		if pipeErr != nil {
			return "", pipeErr
		}

		startErr := cmd.Start()
		if startErr != nil {
			return "", startErr
		}

		buff := bytes.NewBuffer([]byte{})
		_, rbErr := buff.ReadFrom(stdoutPipe)
		if rbErr != nil {
			return "", rbErr
		}

		return strings.TrimRight(buff.String(), "\n"), cmd.Wait()
	}()

	if vErr != nil {
		return fmt.Sprintf("nginx server version check error: %s\n", vErr)
	}
	return svrVersion
}

func (n nginxConfigService) DisplayConfig() (resp []byte, err error) {
	resp = n.nginxConfigManager.Read()
	if resp == nil {
		err = ErrDataNotParsed
		return nil, err
	}
	return resp, err
}

func (n nginxConfigService) GetConfig() (resp []byte, err error) {
	resp = n.nginxConfigManager.ReadJson()
	if resp == nil {
		err = ErrDataNotParsed
		return nil, err
	}
	return resp, err
}

func (n nginxConfigService) ShowStatistics() (resp []byte, err error) {
	resp = n.nginxConfigManager.ReadStatistics()
	if resp == nil {
		err = ErrDataNotParsed
		return nil, err
	}
	return resp, err
}

func (n *nginxConfigService) UpdateConfig(data []byte) error {
	if len(data) > 0 {
		err := n.nginxConfigManager.UpdateFromJsonBytes(data)
		if err != nil {
			return err
		}
		return nil
	}
	return ErrEmptyConfig
}

func (n nginxConfigService) WatchLog(logName string) (LogWatcher, error) {
	// 开始监控日志
	ticker := time.Tick(time.Second)
	timeout := time.After(time.Minute * 30)
	err := n.startWatchLog(logName)
	//err = n.nginxLog.Watch(logName)
	if err != nil {
		return nil, err
	}

	dataChan := make(chan []byte)
	transferErrChan := make(chan error)
	signalChan := make(chan int)

	watcher := NewLogWatcher(dataChan, transferErrChan, func() error {
		signalChan <- 9
		return nil
	})
	// 监听终止信号和每秒读取日志并发送
	//fmt.Println("监听终止信号及准备发送日志")
	go func(sigChan chan int) {
		var transferErr error
		defer func() {
			if transferErr != nil {
				utils.Logger.WarningF("[%s] watch log (log file: %s) goroutine is stop with error: %s", n.name, logName, transferErr)
				transferErrChan <- transferErr
			}
			_ = n.stopWatchLog(logName)
		}()
		for {
			select {
			case s := <-sigChan:
				if s == 9 {
					//fmt.Println("watch log stopping...")
					//fmt.Println("watch log is stopped")
					return
				}
			case <-ticker:
				//fmt.Println("读取日志")
				data, transferErr := n.nginxLog.Watch(logName)
				if transferErr != nil {
					return
				}
				if len(data) > 0 {
					select {
					case dataChan <- data:
						// 日志推送后，客户端已经终止，handler日志推送阻断且发送了终止信号，由于日志推送阻断，接收终止信息被积压
						//fmt.Println("svc发送日志成功")
					case <-time.After(time.Second * 30):
						transferErr = ErrDataSendingTimeout
						return
					}
				}
			case <-timeout:
				transferErr = ErrWatchLogTimeout
				return
			}
		}
	}(signalChan)
	return watcher, nil
}

func (n *nginxConfigService) startWatchLog(logName string) error {
	logDir, err := filepath.Abs(filepath.Join(filepath.Dir(filepath.Dir(n.verifyExecPath)), "logs"))
	if err != nil {
		return err
	}
	return n.nginxLog.StartWatch(logName, logDir)
}

func (n *nginxConfigService) stopWatchLog(logName string) error {
	return n.nginxLog.StopWatch(logName)
}

// configReload, nginxConfigService的nginx配置重载方法，根据nginx配置文件信息加载nginx配置并记录文件基准信息
// 返回值:
//     nginx配置对象指针
//     错误
func (n *nginxConfigService) configReload() error {
	return n.nginxConfigManager.Reload()
}

// bak, ServerInfo的nginx配置文件备份子方法
// 参数:
func (n *nginxConfigService) configBackup() error {
	return n.nginxConfigManager.Backup(n.backupDir, n.backupSaveTime, n.backupCycle)
}

func newNginxConfigService(info WebServerConfigInfo) WebServerConfigService {
	configAbsPath, err := filepath.Abs(info.ConfPath)
	if err != nil {
		utils.Logger.FatalF("[%s] config path error: %s", info.Name, err)
		return nil
	}
	nginxConfigManager, err := configuration.NewNginxConfigurationManager(info.VerifyExecPath, configAbsPath)
	if err != nil {
		utils.Logger.FatalF("[%s] config manager initial error: %s", info.Name, err)
		return nil
	}
	return &nginxConfigService{
		name:               info.Name,
		backupCycle:        info.BackupCycle,
		backupSaveTime:     info.BackupSaveTime,
		backupDir:          info.BackupDir,
		confPath:           info.ConfPath,
		verifyExecPath:     info.VerifyExecPath,
		nginxConfigManager: nginxConfigManager,
		nginxLog:           ngLog.NewLog(),
	}
}
