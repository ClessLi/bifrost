package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	ngLog "github.com/ClessLi/bifrost/pkg/log/nginx"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type offstageViewer interface {
	DisplayConfig(serverName string) ([]byte, error)
	GetConfig(serverName string) ([]byte, error)
	ShowStatistics(serverName string) ([]byte, error)
	DisplayStatus() ([]byte, error)
}

type offstageUpdater interface {
	UpdateConfig(serverName string, data []byte) error
}

type offstageWatcher interface {
	WatchLog(serverName, logName string) (LogWatcher, error)
}

type OffstageManager interface {
	Start() error
	Stop() error
}

type Metrics interface {
	Start() error
	Stop() error
	Report() ([]byte, error)
}

type metrics struct {
	// DONE: 添加web服务版本信息、web服务状态信息(README.md需调整相关接口文档)
	OS                string          `json:"system"`
	Time              string          `json:"time"`
	Cpu               string          `json:"cpu"`
	Mem               string          `json:"mem"`
	Disk              string          `json:"disk"`
	StatusList        []WebServerInfo `json:"status_list"`
	BifrostVersion    string          `json:"bifrost_version"`
	isStoped          bool
	monitorErrChan    chan error
	webServerInfoFunc func() []WebServerInfo
	locker            *sync.RWMutex
}

func (m *metrics) Start() error {
	m.locker.Lock()
	defer m.locker.Unlock()
	if !m.isStoped {
		return ErrMonitoringStarted
	}

	m.isStoped = false
	go func() {
		var sysErr error
		defer func() {
			m.locker.Lock()
			m.isStoped = true
			m.locker.Unlock()
			utils.Logger.DebugF("system info monitor is stopping...")
			m.monitorErrChan <- sysErr
		}()
		platform, _, release, OSErr := host.PlatformInformation()
		if OSErr != nil {
			utils.Logger.FatalF("Failed to initialize monitor, cased by '%s'", OSErr)
			m.isStoped = true
		} else {
			m.OS = fmt.Sprintf("%s %s", platform, release)
		}
		for !m.isStoped {
			// 监控数据获取

			// cpu监控
			cpupct, sysErr := cpu.Percent(time.Second*5, false)
			if sysErr != nil {
				return
			}

			// 内存监控
			vmem, sysErr := mem.VirtualMemory()
			if sysErr != nil {
				return
			}

			// 磁盘监控
			diskInfo, sysErr := disk.Usage("/")
			if sysErr != nil {
				return
			}

			// 监控数据写入
			m.locker.Lock()
			// web server 监控
			m.StatusList = m.webServerInfoFunc()
			m.Cpu = fmt.Sprintf("%.2f", cpupct[0])
			m.Mem = fmt.Sprintf("%.2f", vmem.UsedPercent)
			m.Disk = fmt.Sprintf("%.2f", diskInfo.UsedPercent)

			// 监控时间戳
			m.Time = fmt.Sprintf("%s", time.Now().In(nginx.TZ).Format("2006/01/02 15:04:05"))
			//time.Sleep(time.Second * 5)

			m.locker.Unlock()
		}
	}()
	return nil
}

func (m *metrics) Stop() error {
	if m.isStoped {
		return ErrMonitoringServiceSuspension
	}
	m.locker.Lock()
	m.isStoped = true
	m.locker.Unlock()
	select {
	case err := <-m.monitorErrChan:
		return err
	case <-time.After(time.Second * 30):
		return ErrStopMonitoringTimeout
	}
}

func (m *metrics) Report() ([]byte, error) {
	m.locker.RLock()
	defer m.locker.RUnlock()
	return json.Marshal(m)
}

func NewMetrics(webServerInfoFunc func() []WebServerInfo, errChan chan error) Metrics {
	return &metrics{
		StatusList:        make([]WebServerInfo, 0),
		BifrostVersion:    utils.Version(),
		isStoped:          true,
		monitorErrChan:    errChan,
		webServerInfoFunc: webServerInfoFunc,
		locker:            new(sync.RWMutex),
	}
}

type WebServerInfo struct {
	Name    string `json:"name"`
	Status  State  `json:"status"`
	Version string `json:"version"`
}

func NewWebServerInfo(name string) WebServerInfo {
	return WebServerInfo{
		Name:    name,
		Status:  UnknownState,
		Version: "unknown",
	}
}

type Offstage struct {
	webServerConfigServices map[string]WebServerConfigService
	webServerConfigManagers map[string]configuration.ConfigManager
	metrics                 Metrics
}

func (o Offstage) Range(rangeFunc func(serverName string, configService WebServerConfigService) bool) {
	for s, configService := range o.webServerConfigServices {
		if !rangeFunc(s, configService) {
			return
		}
	}
}

func (o Offstage) DisplayConfig(serverName string) ([]byte, error) {
	service, has := o.webServerConfigServices[serverName]
	if has {
		return service.configuration.View(), nil
	}
	return nil, ErrWebServerConfigServiceNotExist
}

func (o Offstage) GetConfig(serverName string) ([]byte, error) {
	service, has := o.webServerConfigServices[serverName]
	if has {
		return service.configuration.Json(), nil
	}
	return nil, ErrWebServerConfigServiceNotExist
}

func (o Offstage) ShowStatistics(serverName string) ([]byte, error) {
	service, has := o.webServerConfigServices[serverName]
	if has {
		return service.configuration.StatisticsByJson(), nil
	}
	return nil, ErrWebServerConfigServiceNotExist
}

func (o Offstage) DisplayStatus() ([]byte, error) {
	return o.metrics.Report()
}

func (o *Offstage) UpdateConfig(serverName string, data []byte) error {
	svc, has := o.webServerConfigServices[serverName]
	if has {
		return svc.configuration.UpdateFromJsonBytes(data)
	}
	return ErrWebServerConfigServiceNotExist
}

func (o Offstage) WatchLog(serverName, logName string) (LogWatcher, error) {
	svc, has := o.webServerConfigServices[serverName]
	if !has {
		return nil, ErrWebServerConfigServiceNotExist
	}
	// 开始监控日志
	ticker := time.Tick(time.Second)
	timeout := time.After(time.Minute * 30)
	err := svc.log.StartWatch(logName, svc.logsDir)
	//err = n.nginxLog.Watch(logName)
	if err != nil {
		return nil, err
	}

	dataChan := make(chan []byte)
	transferErrChan := make(chan error)
	signalChan := make(chan int)

	lw := NewLogWatcher(dataChan, transferErrChan, func() error {
		signalChan <- 9
		return nil
	})
	// 监听终止信号和每秒读取日志并发送
	//fmt.Println("监听终止信号及准备发送日志")
	go func(sigChan chan int) {
		var transferErr error
		defer func() {
			if transferErr != nil {
				utils.Logger.WarningF("[%s] watch log (log file: %s) goroutine is stop with error: %s", serverName, logName, transferErr)
				transferErrChan <- transferErr
			}
			_ = svc.log.StopWatch(logName)
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
				data, transferErr := svc.log.Watch(logName)
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
	return lw, nil
}

func (o *Offstage) Start() error {
	// 启动web服务配置文件管理器
	for mngName, manager := range o.webServerConfigManagers {
		utils.Logger.DebugF("[%s] config manager starting...", mngName)
		err := manager.Start()
		if err != nil {
			panic(fmt.Errorf("[%s] config manager start error: %s", mngName, err))
		}
		utils.Logger.DebugF("[%s] config manager is started", mngName)
	}
	return o.metrics.Start()

}

func (o *Offstage) Stop() error {
	errorStr := ""
	for mngName, manager := range o.webServerConfigManagers {
		err := manager.Stop()
		if err != nil {
			if errorStr != "" {
				errorStr += "; "
			}
			errorStr += fmt.Sprintf("[%s] manager stop error: %s", mngName, err)
		}
	}
	err := o.metrics.Stop()
	if err != nil {
		if errorStr != "" {
			errorStr += "; "
		}
		errorStr += fmt.Sprintf("metrics stop error: %s", err)
	}
	if errorStr != "" {
		return fmt.Errorf(errorStr)
	}
	return nil
}

func NewOffstage(services map[string]WebServerConfigService, configManagers map[string]configuration.ConfigManager, metrics Metrics) *Offstage {
	return &Offstage{
		webServerConfigServices: services,
		webServerConfigManagers: configManagers,
		metrics:                 metrics,
	}
}

type WebServerConfigService struct {
	configuration configuration.Configuration
	serverBinPath string
	logsDir       string
	log           *ngLog.Log
}

func (w WebServerConfigService) ServerVersion() string {
	svrBinAbs, absErr := filepath.Abs(w.serverBinPath)
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

func (w WebServerConfigService) ServerStatus() State {
	svrPidFilePath := "logs/nginx.pid"
	svrPidQueryer, err := w.configuration.Query("key:sep: :reg:pid .*")
	if err == nil {
		svrPidFilePath = strings.Split(svrPidQueryer.Self().GetValue(), " ")[1]
	}

	svrPidFilePathAbs := svrPidFilePath
	if !filepath.IsAbs(svrPidFilePath) {
		svrBinAbs, absErr := filepath.Abs(w.serverBinPath)
		if absErr != nil {
			return UnknownState
		}
		svrWS, wsErr := filepath.Abs(filepath.Join(filepath.Dir(svrBinAbs), ".."))
		if wsErr != nil {
			return UnknownState
		}
		var pidErr error
		svrPidFilePathAbs, pidErr = filepath.Abs(filepath.Join(svrWS, svrPidFilePath))
		if pidErr != nil {
			return UnknownState
		}
	}

	svrPid, gPidErr := utils.GetPid(svrPidFilePathAbs)
	if gPidErr != nil {
		return Abnormal
	}

	_, procErr := os.FindProcess(svrPid)
	if procErr != nil {
		return Abnormal
	}
	return Normal
}

func NewWebServerConfigService(configuration configuration.Configuration, serverBinPath, logsDir string, log *ngLog.Log) WebServerConfigService {
	return WebServerConfigService{
		configuration: configuration,
		serverBinPath: serverBinPath,
		logsDir:       logsDir,
		log:           log,
	}
}

type LogWatcher interface {
	GetDataChan() <-chan []byte
	GetTransferErrorChan() <-chan error
	Close() error
}

type logWatcher struct {
	dataChan          chan []byte
	transferErrorChan chan error
	closeFunc         func() error
}

func NewLogWatcher(dataChan chan []byte, errChan chan error, closeFunc func() error) LogWatcher {
	return &logWatcher{
		dataChan:          dataChan,
		transferErrorChan: errChan,
		closeFunc:         closeFunc,
	}
}

func (l logWatcher) GetDataChan() <-chan []byte {
	return l.dataChan
}

func (l logWatcher) GetTransferErrorChan() <-chan error {
	return l.transferErrorChan
}

func (l logWatcher) Close() error {
	return l.closeFunc()
}
