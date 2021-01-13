package web_server_manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	ngJson "github.com/ClessLi/bifrost/pkg/json/nginx"
	ngLog "github.com/ClessLi/bifrost/pkg/log/nginx"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	ngStatistics "github.com/ClessLi/bifrost/pkg/statistics/nginx"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type nginxConfigService struct {
	name           string
	backupCycle    int
	backupSaveTime int
	backupDir      string
	confPath       string
	verifyExecPath string
	confCaches     nginx.Caches
	nginxConfig    *nginx.Config
	nginxLog       *ngLog.Log
}

func (n nginxConfigService) serverName() string {
	return n.name
}

func (n nginxConfigService) DisplayWebServerStatus() State {
	svrPidFileKW := nginx.NewKeyWords(nginx.TypeKey, "pid", "*", false, false)
	svrPidFilePath := "logs/nginx.pid"
	svrPidFileKey, ok := n.nginxConfig.QueryByKeywords(svrPidFileKW).(*nginx.Key)
	if ok && svrPidFileKey != nil {
		svrPidFilePath = svrPidFileKey.Value
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
		//cmd.Stderr = Stdoutf
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
	for _, s := range n.nginxConfig.String() {
		resp = append(resp, []byte(s)...)
	}
	if resp == nil {
		err = ErrDataNotParsed
		return nil, err
	}
	return resp, err
}

func (n nginxConfigService) GetConfig() ([]byte, error) {
	return json.Marshal(n.nginxConfig)
}

func (n nginxConfigService) ShowStatistics() ([]byte, error) {
	httpServersNum, httpServers := ngStatistics.HTTPServers(n.nginxConfig)
	httpPorts := ngStatistics.HTTPPorts(n.nginxConfig)
	streamServersNum, streamPorts := ngStatistics.StreamServers(n.nginxConfig)
	statistics := struct {
		HttpSvrsNum   int              `json:"http_svrs_num"`
		HttpSvrs      map[string][]int `json:"http_svrs"`
		HttpPorts     []int            `json:"http_ports"`
		StreamSvrsNum int              `json:"stream_svrs_num"`
		StreamPorts   []int            `json:"stream_ports"`
	}{HttpSvrsNum: httpServersNum, HttpSvrs: httpServers, HttpPorts: httpPorts, StreamSvrsNum: streamServersNum, StreamPorts: streamPorts}
	return json.Marshal(statistics)
}

func (n *nginxConfigService) UpdateConfig(data []byte) error {
	if len(data) > 0 {

		// check config
		newConfig, err := ngJson.Unmarshal(data)
		if err != nil {
			//Log(DEBUG, "[%s] [%s] buffer detail: %s", info.Name, ip, buffer.String())
			//Log(WARN, "[%s] [%s] %s detailed error: %s", info.Name, ip, bifrostpb.ConfigUnmarshalErr, err)
			return err
		}

		//fmt.Println("获取web服务配置校验二进制文件路径")
		verifyBin, err := filepath.Abs(n.verifyExecPath)
		if err != nil {
			err = ErrValidationNotExist
			return err
		}

		// delete old config
		err = nginx.Delete(n.nginxConfig)
		if err != nil {
			//message = fmt.Sprintf("Delete nginx ng failed. <%s>", err)
			//Log(ERROR, "[%s] [%s] %s", info.Name, ip, message)
			return err
		}

		newCaches, err := nginx.SaveWithCheck(newConfig, verifyBin)
		// roll back
		if err != nil {
			var rollErr error
			rollErr = nginx.Delete(newConfig)
			if rollErr != nil {
				//Log(ERROR, "[%s] Delete new nginx ng failed. <%s>", info.Name, err)
				//message = "New nginx config verify failed. And delete new nginx config failed."
				return rollErr
			}

			_, rollErr = nginx.Save(n.nginxConfig)
			if rollErr != nil {
				//Log(CRITICAL, "[%s] Nginx ng rollback failed. <%s>", info.Name, err)
				//message = "New nginx config verify failed. And nginx config rollback failed."
				return rollErr
			}

			return err
		}
		n.confCaches = newCaches
		n.nginxConfig = newConfig
		n.confPath = newConfig.Value

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

// checkHash, ServerInfo的web服务器配置文件是否已更改校验方法
func (n nginxConfigService) checkConfigsHash() (isSame bool, err error) {
	isSame = true
	for path := range n.confCaches {
		if isSame, err = n.confCaches.CheckHash(path); !isSame {
			return
		}
	}
	return
}

// confLoad, ServerInfo的nginx配置加载方法，根据nginx配置文件信息加载nginx配置并记录文件基准信息
// 返回值:
//     nginx配置对象指针
//     错误
func (n *nginxConfigService) configLoad() error {
	// 加载nginx配置并获取缓存
	//Log(DEBUG, "[%s] load config...", n.name)
	path, caches, err := nginx.Load(n.confPath)
	if err != nil {
		//Log(DEBUG, "[%s] load config failed: %s", n.Name, err.Error())
		return err
	}

	config, err := caches.GetConfig(n.confPath)
	if err != nil {
		//Log(DEBUG, "[%s] load config failed: %s", n.Name, err.Error())
		return err
	}
	// 检查nginx配置是否能被正常解析为json
	//Log(DEBUG, "[%s] 校验nginx配置。。。", b.ServiceInfos[i].Name)
	_, err = json.Marshal(config)
	if err != nil {
		//fmt.Printf("[%s] bifrost service failed to start. Cased by '%s'\n", n.name, err)
		utils.Logger.CriticalF("[%s] bifrost service failed to start. Cased by '%s'", n.name, err)
		return err
	}
	// 记录缓存
	n.confCaches = caches
	n.confPath = path
	n.nginxConfig = config

	//Log(DEBUG, "[%s] load config success", n.name)
	return nil
}

// bak, ServerInfo的nginx配置文件备份子方法
// 参数:
func (n *nginxConfigService) configBackup() error {
	bakPath, bErr := nginx.Backup(n.nginxConfig, "nginx.conf", n.backupSaveTime, n.backupCycle, n.backupDir)

	if bErr != nil && (!os.IsExist(bErr) && bErr != nginx.NoBackupRequired) { // 备份失败
		//Log(CRITICAL, "[%s] Nginx Config backup to %s, but failed. <%s>", n.Name, bakPath, bErr)
		utils.Logger.CriticalF("[%s] Nginx Config backup to %s, but failed. <%s>", n.name, bakPath, bErr)
		//Log(NOTICE, "[%s] Nginx Config backup is stop.", n.Name)
		return bErr
	} else if bErr == nil { // 备份成功
		//Log(INFO, "[%s] Nginx Config backup to %s", n.Name, bakPath)
		utils.Logger.InfoF("[%s] Nginx Config backup to %s", n.name, bakPath)
	}
	return nil
}

func newNginxConfigService(info WebServerConfigInfo) WebServerConfigService {
	return &nginxConfigService{
		name:           info.Name,
		backupCycle:    info.BackupCycle,
		backupSaveTime: info.BackupSaveTime,
		backupDir:      info.BackupDir,
		confPath:       info.ConfPath,
		verifyExecPath: info.VerifyExecPath,
		nginxLog:       ngLog.NewLog(),
	}
}
