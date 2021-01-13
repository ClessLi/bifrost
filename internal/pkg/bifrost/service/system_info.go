package service

import (
	"encoding/json"
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service/web_server_manager"
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"sync"
	"time"
)

type systemInfo struct {
	// DONE: 添加web服务版本信息、web服务状态信息(README.md需调整相关接口文档)
	OS                          string           `json:"system"`
	Time                        string           `json:"time"`
	Cpu                         string           `json:"cpu"`
	Mem                         string           `json:"mem"`
	Disk                        string           `json:"disk"`
	StatusList                  []*webServerInfo `json:"status_list"`
	BifrostVersion              string           `json:"bifrost_version"`
	isStoped                    bool
	monitorErrChan              chan error
	displayWebServerVersionFunc func(string) string
	displayWebServerStatusFunc  func(string) web_server_manager.State
	locker                      *sync.RWMutex
}

func (s systemInfo) DisplayStatus() ([]byte, error) {
	s.locker.RLock()
	defer s.locker.RUnlock()
	if !s.isStoped {
		return json.Marshal(s)
	}
	return nil, ErrMonitoringServiceSuspension
}

func (s *systemInfo) Start() error {
	s.locker.Lock()
	defer s.locker.Unlock()
	if !s.isStoped {
		return ErrMonitoringStarted
	}
	for _, serverMonitor := range s.StatusList {
		serverMonitor.Version = s.displayWebServerVersionFunc(serverMonitor.Name)
	}

	s.isStoped = false
	go func() {
		var sysErr error
		defer func() {
			s.locker.Lock()
			s.isStoped = true
			s.locker.Unlock()
			//utils.Logger.DebugF("system info monitor is stopping...")
			s.monitorErrChan <- sysErr
		}()
		for !s.isStoped {
			// 监控数据获取
			// web server 监控
			stateList := make([]web_server_manager.State, 0)
			for _, serverMonitor := range s.StatusList {
				stateList = append(stateList, s.displayWebServerStatusFunc(serverMonitor.Name))
			}

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
			s.locker.Lock()
			for i := 0; i < len(s.StatusList); i++ {
				s.StatusList[i].Status = stateList[i]
			}
			s.Cpu = fmt.Sprintf("%.2f", cpupct[0])
			s.Mem = fmt.Sprintf("%.2f", vmem.UsedPercent)
			s.Disk = fmt.Sprintf("%.2f", diskInfo.UsedPercent)

			// 监控时间戳
			s.Time = fmt.Sprintf("%s", time.Now().In(nginx.TZ).Format("2006/01/02 15:04:05"))
			//time.Sleep(time.Second * 5)

			s.locker.Unlock()
		}
	}()
	return nil

}

func (s *systemInfo) Stop() error {
	if s.isStoped {
		return ErrMonitoringServiceSuspension
	}
	s.locker.Lock()
	s.isStoped = true
	s.locker.Unlock()
	select {
	case err := <-s.monitorErrChan:
		return err
	case <-time.After(time.Second * 30):
		return ErrStopMonitoringTimeout
	}
}

func NewSysInfo(servicesHandler *web_server_manager.WebServerConfigServicesHandler) *systemInfo {
	if len(servicesHandler.ServerNames()) < 1 {
		return nil
	}
	platform, _, release, OSErr := host.PlatformInformation()
	if OSErr != nil {
		utils.Logger.FatalF("Failed to initialize monitor, cased by '%s'", OSErr)
		return nil
	}
	systemInfo := &systemInfo{
		OS:             fmt.Sprintf("%s %s", platform, release),
		StatusList:     make([]*webServerInfo, 0),
		BifrostVersion: utils.Version(),
		monitorErrChan: make(chan error),
		locker:         new(sync.RWMutex),
		isStoped:       true,
		displayWebServerVersionFunc: func(name string) string {
			v, _ := servicesHandler.DisplayVersion(name)
			return v
		},
		displayWebServerStatusFunc: func(name string) web_server_manager.State {
			s, _ := servicesHandler.DisplayStatus(name)
			return s
		},
	}
	for _, name := range servicesHandler.ServerNames() {
		systemInfo.StatusList = append(systemInfo.StatusList, newWebServerMonitor(name))
	}

	return systemInfo
}

type webServerInfo struct {
	Name    string                   `json:"name"`
	Status  web_server_manager.State `json:"status"`
	Version string                   `json:"version"`
}

func newWebServerMonitor(name string) *webServerInfo {
	return &webServerInfo{
		Name:    name,
		Status:  web_server_manager.Unknown,
		Version: "unknown",
	}
}
