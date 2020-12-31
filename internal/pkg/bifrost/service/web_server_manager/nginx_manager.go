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
	"sync"
	"time"
)

type nginxManager struct {
	name           string
	backupCycle    int
	backupSaveTime int
	backupDir      string
	confPath       string
	verifyExecPath string
	confCaches     nginx.Caches
	nginxConfig    *nginx.Config
	nginxLog       *ngLog.Log
	available      bool
	autoBackupChan chan int
	autoReloadChan chan int
	waitGroup      *sync.WaitGroup
}

func NewNginxManager(name string, backupCycle, backupSaveTime int, backupDir, confPath, verifyExecPath string) WebServerManager {
	return &nginxManager{
		name:           name,
		backupCycle:    backupCycle,
		backupSaveTime: backupSaveTime,
		backupDir:      backupDir,
		confPath:       confPath,
		verifyExecPath: verifyExecPath,
		waitGroup:      new(sync.WaitGroup),
	}

}

func (n nginxManager) DisplayConfig() (resp []byte, err error) {
	for _, s := range n.nginxConfig.String() {
		resp = append(resp, []byte(s)...)
	}
	if resp == nil {
		err = ErrDataNotParsed
		return nil, err
	}
	return resp, err
}

func (n nginxManager) GetConfig() ([]byte, error) {
	return json.Marshal(n.nginxConfig)
}

func (n nginxManager) ShowStatistics() ([]byte, error) {
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

func (n *nginxManager) UpdateConfig(data []byte, param string) error {
	//utils.Logger.Debug("nginx manager update config")
	switch param {
	case "full":
		//utils.Logger.Debug("nginx manager full update config")
		if len(data) > 0 {

			// check config
			newConfig, err := ngJson.Unmarshal(data)
			if err != nil {
				//Log(DEBUG, "[%s] [%s] buffer detail: %s", info.Name, ip, buffer.String())
				//Log(WARN, "[%s] [%s] %s detailed error: %s", info.Name, ip, bifrostpb.ConfigUnmarshalErr, err)
				return err
			}

			return n.updateFullConfig(newConfig)

		}
		return ErrEmptyConfig
	default:
		return ErrWrongParamPassedIn
	}
}

func (n nginxManager) WatchLog(logName string) (Watcher, error) {
	// 开始监控日志
	ticker := time.Tick(time.Second)
	timeout := time.After(time.Minute * 30)
	err := n.startWatchLog(logName)
	//err = n.nginxLog.Watch(logName)
	if err != nil {
		return nil, err
	}

	watcher := NewLogWatcher()
	// 监听终止信号和每秒读取日志并发送
	//fmt.Println("监听终止信号及准备发送日志")
	go func() {
		defer func() {
			if err != nil {
				watcher.inputErrChan() <- err
			}
		}()
		for {
			select {
			case s := <-watcher.getSignalChan():
				//fmt.Println("sig is:", s)
				if s == 9 {
					//fmt.Println("watch log stopping...")
					err = n.stopWatchLog(logName)
					//fmt.Println("watch log is stopped")
					return
				}
			case <-ticker:
				//fmt.Println("读取日志")
				data, err := n.nginxLog.Watch(logName)
				if err != nil {
					return
				}
				if len(data) > 0 {
					select {
					case watcher.inputDataChan() <- data:
						// 日志推送后，客户端已经终止，handler日志推送阻断且发送了终止信号，由于日志推送阻断，接收终止信息被积压
						//fmt.Println("svc发送日志成功")
					case <-time.After(time.Second * 30):
						_ = n.stopWatchLog(logName)
						err = ErrDataSendingTimeout
						return
					}
				}
			case <-timeout:
				err = ErrWatchLogTimeout
				return
			}
		}
	}()
	return watcher, nil
}

func (n *nginxManager) updateFullConfig(newConfig *nginx.Config) error {
	//fmt.Println("获取web服务配置校验二进制文件路径")
	verifyBin, err := filepath.Abs(n.verifyExecPath)
	if err != nil {
		//Log(CRITICAL, "[%s] %s detailed error: %s", info.Name, bifrostpb.ValidationNotExist, err)
		err = ErrValidationNotExist
		return err
	}

	// delete old config
	err = nginx.Delete(n.nginxConfig)
	//message := ""
	if err != nil {
		//message = fmt.Sprintf("Delete nginx ng failed. <%s>", err)
		//Log(ERROR, "[%s] [%s] %s", info.Name, ip, message)
		return err
	}

	//Log(INFO, "[%s] Deleted old nginx config.", info.Name)
	//Log(INFO, "[%s] Verify new nginx config.", info.Name)
	newCaches, err := nginx.SaveWithCheck(newConfig, verifyBin)
	// roll back
	if err != nil {
		//Log(DEBUG, "[%s] Roll back to old nginx config.", info.Name)
		//message = fmt.Sprintf("Nginx ng verify failed. <%s>", err)
		//Log(WARN, "[%s] %s", info.Name, message)

		//Log(INFO, "[%s] Delete new nginx ng.", info.Name)
		var rollErr error
		rollErr = nginx.Delete(newConfig)
		if rollErr != nil {
			//Log(ERROR, "[%s] Delete new nginx ng failed. <%s>", info.Name, err)
			//message = "New nginx config verify failed. And delete new nginx config failed."
			return rollErr
		}

		//Log(INFO, "[%s] Rollback nginx ng.", info.Name)
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

	//Log(NOTICE, "[%s] [%s] Nginx Config saved successfully", info.Name, ip)
	return nil
}

func (n *nginxManager) ConfigLoad() error {
	if n.available {
		return n.confLoad()
	}
	return ErrServiceNotAvailable
}

// confLoad, ServerInfo的nginx配置加载方法，根据nginx配置文件信息加载nginx配置并记录文件基准信息
// 返回值:
//     nginx配置对象指针
//     错误
func (n *nginxManager) confLoad() error {
	// 加载nginx配置并获取缓存
	//Log(DEBUG, "[%s] load config...", n.name)
	path, caches, err := nginx.Load(n.confPath)
	if err != nil {
		//Log(DEBUG, "[%s] load config failed: %s", n.Name, err.Error())
		return err
	}

	// 记录缓存
	n.confCaches = caches
	n.confPath = path
	n.nginxConfig, err = n.confCaches.GetConfig(n.confPath)
	if err != nil {
		//Log(DEBUG, "[%s] load config failed: %s", n.Name, err.Error())
		return err
	}
	//Log(DEBUG, "[%s] load config success", n.name)

	return nil

}

// autoBackup, ServerInfo的nginx配置文件备份方法
// 参数:
//     c: 整型管道，用于停止备份
func (n *nginxManager) autoBackup() {
	n.waitGroup.Add(1)
	//defer func() {
	//	if !n.available {
	//		//Log(INFO, "[%s] %s, Nginx Config backup is stop.", n.Name, ErrServiceNotAvailable)
	//	} else {
	//		//Log(NOTICE, "[%s] Nginx Config backup is stop.", n.Name)
	//	}
	//}()
	n.autoBackupChan = make(chan int)
	go func() {
		defer n.waitGroup.Done()
		//defer close(n.autoBackupChan)
		for n.available {
			select {
			case <-time.NewTicker(5 * time.Minute).C: // 每5分钟定时执行备份操作
				utils.Logger.DebugF("[%s] Nginx Config check and backup", n.name)
				n.bak()
			case sig := <-n.autoBackupChan: // 获取管道传入信号
				if sig == 9 { // 为9时，停止备份
					utils.Logger.InfoF("[%s] Nginx Config autoBackup method stopping...", n.name)
					goto stopHere
				}
			}
		}
	stopHere:
		return
	}()
}

// bak, ServerInfo的nginx配置文件备份子方法
// 参数:
func (n *nginxManager) bak() {
	if n.available {
		bakPath, bErr := nginx.Backup(n.nginxConfig, "nginx.conf", n.backupSaveTime, n.backupCycle, n.backupDir)

		if bErr != nil && (!os.IsExist(bErr) && bErr != nginx.NoBackupRequired) { // 备份失败
			//Log(CRITICAL, "[%s] Nginx Config backup to %s, but failed. <%s>", n.Name, bakPath, bErr)
			fmt.Printf("[%s] Nginx Config backup to %s, but failed. <%s>\n", n.name, bakPath, bErr)
			//Log(NOTICE, "[%s] Nginx Config backup is stop.", n.Name)
		} else if bErr == nil { // 备份成功
			//Log(INFO, "[%s] Nginx Config backup to %s", n.Name, bakPath)
			fmt.Printf("[%s] Nginx Config backup to %s\n", n.name, bakPath)
		}
	}
}

// autoReload, ServerInfo的web服务器配置文件自动热加载方法
// 参数:
//     c: 整型管道，用于停止备份
func (n *nginxManager) autoReload() {
	n.waitGroup.Add(1)
	//defer func() {
	//	if !n.available {
	//		Log(INFO, "[%s] %s, Nginx Config auto reload is stop.", ErrServiceNotAvailable, n.Name)
	//	} else {
	//		Log(NOTICE, "[%s] Nginx Config auto reload is stop.", n.Name)
	//	}
	//}()
	n.autoReloadChan = make(chan int)
	go func() {
		defer n.waitGroup.Done()
		//defer close(n.autoReloadChan)
		for n.available {
			select {
			case <-time.NewTicker(30 * time.Second).C: // 每30秒检查一次nginx配置文件是否已在后台更新
				utils.Logger.DebugF("[%s] Nginx Config check and reloading", n.name)
				reloadErr := n.confReload()
				if reloadErr != nil && reloadErr != nginx.NoReloadRequired {
					utils.Logger.WarningF("[%s] Nginx Config reload failed, cased by '%s'", n.name, reloadErr)
				} else if reloadErr == nil {
					utils.Logger.InfoF("[%s] Nginx Config reload successfully", n.name)
				}
			case sig := <-n.autoReloadChan: // 获取管道传入信号
				if sig == 9 { // 为9时，停止备份
					utils.Logger.InfoF("[%s] Nginx Config autoReload method stopping...", n.name)
					goto stopHere
				}
			}
		}
	stopHere:
		return
	}()
}

// autoReload, ServerInfo的web服务器配置文件自动热加载子方法
func (n *nginxManager) confReload() error {
	if n.available {
		// 校验配置文件是否更新
		isSame, checkErr := n.checkHash()
		if checkErr != nil {
			return checkErr
		}

		// 如果有差别，则重新读取配置
		if !isSame {
			//Log(DEBUG, "[%s] reloading nginx config", n.Name)
			return n.confLoad()
		}
		return nginx.NoReloadRequired
	}
	return ErrServiceNotAvailable
}

func (n *nginxManager) startWatchLog(logName string) error {
	logDir, err := filepath.Abs(filepath.Join(filepath.Dir(filepath.Dir(n.verifyExecPath)), "logs"))
	if err != nil {
		return err
	}
	return n.nginxLog.StartWatch(logName, logDir)
}

//func (n *nginxManager) WatchLog(logName string) (data []byte, err error) {
//	return n.nginxLog.Watch(logName)
//}

func (n *nginxManager) stopWatchLog(logName string) error {
	return n.nginxLog.StopWatch(logName)
}

// checkHash, ServerInfo的web服务器配置文件是否已更改校验方法
func (n nginxManager) checkHash() (isSame bool, err error) {
	isSame = true
	for path := range n.confCaches {
		if isSame, err = n.confCaches.CheckHash(path); !isSame {
			return
		}
	}
	return
}

func (n *nginxManager) Enable() {
	n.available = true
	//Log(INFO, "[%s] bifrost service enabled", n.Name)
}

func (n *nginxManager) Disable() {
	n.available = false
	//Log(INFO, "[%s] bifrost service disabled", n.Name)
}

func (n nginxManager) ShowVersion() string {
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

func (n nginxManager) DisplayStatus() status {
	if n.available {
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
			//if SysInfo.StatusList[i].ServerStatus != "abnormal" {
			//Log(WARN, "[%s] something wrong with web server: %s", b.ServiceInfos[i].Name, gPidErr)
			//}
			return Abnormal
		}

		_, procErr := os.FindProcess(svrPid)
		if procErr != nil {
			return Abnormal
		}

		//if SysInfo.StatusList[i].ServerStatus != "normal" {
		//Log(INFO, "[%s] web server <PID: %d> is running.", b.ServiceInfos[i].Name, svrPid)
		//}
		return Normal
	}
	return Disabled
}

func (n *nginxManager) ManagementStart() (err error) {
	//Log(DEBUG, "[%s] 初始化bifrost服务相关接口。。。", b.ServiceInfos[i].Name)
	n.Enable()
	defer func() {
		if err != nil {
			n.Disable()
		}
	}()
	err = n.ConfigLoad()
	if err != nil {
		fmt.Printf("[%s] load config error: %s\n", n.name, err)
		//Log(ERROR, "[%s] load config error: %s", b.ServiceInfos[i].Name, loadErr)
		return err
	}

	// 检查nginx配置是否能被正常解析为json
	//Log(DEBUG, "[%s] 校验nginx配置。。。", b.ServiceInfos[i].Name)
	_, err = json.Marshal(n.nginxConfig)
	if err != nil {
		fmt.Printf("[%s] bifrost service failed to start. Cased by '%s'\n", n.name, err)
		//Log(CRITICAL, "[%s] bifrost service failed to start. Cased by '%s'", b.ServiceInfos[i].Name, jerr)
		return err
	}

	// DONE: 执行备份与自动加载
	n.autoBackup()
	//Log(DEBUG, "[%s] 载入备份协程", b.ServiceInfos[i].Name)
	n.autoReload()
	//Log(DEBUG, "[%s] 载入自动更新配置协程", b.ServiceInfos[i].Name)
	n.nginxLog = ngLog.NewLog()
	return nil
}

func (n *nginxManager) ManagementStop() error {
	defer n.Disable()
	defer n.waitGroup.Wait()
	//utils.Logger.DebugF("[%s] stop backup proc", n.name)
	if n.autoBackupChan != nil {
		utils.Logger.DebugF("[%s] stop backup proc", n.name)
		n.autoBackupChan <- 9
	}
	//utils.Logger.DebugF("[%s] stop config auto reload proc", n.name)
	if n.autoReloadChan != nil {
		utils.Logger.DebugF("[%s] stop config auto reload proc", n.name)
		n.autoReloadChan <- 9
	}
	return nil
}
