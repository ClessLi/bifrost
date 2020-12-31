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

const (
	unknown status = iota
	disabled
	abnormal
	normal
)

type status int

type systemInfo struct {
	// DONE: 添加web服务版本信息、web服务状态信息(README.md需调整相关接口文档)
	OS             string              `json:"system"`
	Time           string              `json:"time"`
	Cpu            string              `json:"cpu"`
	Mem            string              `json:"mem"`
	Disk           string              `json:"disk"`
	StatusList     []*webServerMonitor `json:"status_list"`
	BifrostVersion string              `json:"bifrost_version"`
	isStoped       bool
	monitorErrChan chan error
	locker         *sync.Mutex
}

func NewSysInfo(webServerManagers map[string]web_server_manager.WebServerManager) Monitor {
	if webServerManagers == nil {
		return nil
	}
	platform, _, release, OSErr := host.PlatformInformation()
	if OSErr != nil {
		utils.Logger.FatalF("Failed to initialize systemInfo, cased by '%s'", OSErr)
		return nil
	}
	sysInfo := &systemInfo{
		OS:             fmt.Sprintf("%s %s", platform, release),
		StatusList:     make([]*webServerMonitor, 0),
		BifrostVersion: utils.Version(),
		monitorErrChan: make(chan error),
		locker:         new(sync.Mutex),
	}
	for name, manager := range webServerManagers {
		sysInfo.StatusList = append(sysInfo.StatusList, newWebServerMonitor(name, manager))
	}
	return sysInfo
}

func (s systemInfo) DisplayStatus() ([]byte, error) {
	return json.Marshal(s)
}

func (s *systemInfo) Start() error {
	s.locker.Lock()
	s.isStoped = false
	s.locker.Unlock()
	for _, serverMonitor := range s.StatusList {
		serverMonitor.ShowVersion()
	}

	go func() {
		var sysErr error
		defer func() {
			utils.Logger.DebugF("system info monitor is stopping...")
			s.monitorErrChan <- sysErr
		}()
		for !s.isStoped {
			// web server 监控
			for _, monitor := range s.StatusList {
				monitor.DisplayStatus()
			}

			// cpu监控
			cpupct, sysErr := cpu.Percent(time.Second*5, false)
			if sysErr != nil {
				return
			}
			s.Cpu = fmt.Sprintf("%.2f", cpupct[0])

			// 内存监控
			vmem, sysErr := mem.VirtualMemory()
			if sysErr != nil {
				return
			}
			s.Mem = fmt.Sprintf("%.2f", vmem.UsedPercent)

			// 磁盘监控
			diskInfo, sysErr := disk.Usage("/")
			if sysErr != nil {
				return
			}
			s.Disk = fmt.Sprintf("%.2f", diskInfo.UsedPercent)

			// 监控时间戳
			s.Time = fmt.Sprintf("%s", time.Now().In(nginx.TZ).Format("2006/01/02 15:04:05"))
			//time.Sleep(time.Second * 5)
		}
	}()
	return nil

}

func (s *systemInfo) Stop() error {
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

type webServerMonitor struct {
	Name              string `json:"name"`
	Status            status `json:"status"`
	Version           string `json:"version"`
	showVersionFunc   func() string
	displayStatusFunc func() status
}

func newWebServerMonitor(name string, manager web_server_manager.WebServerManager) *webServerMonitor {
	return &webServerMonitor{
		Name:            name,
		Status:          unknown,
		Version:         "unknown",
		showVersionFunc: manager.ShowVersion,
		displayStatusFunc: func() status {
			var s = unknown
			switch manager.DisplayStatus() {
			case web_server_manager.Disabled:
				s = disabled
			case web_server_manager.Abnormal:
				s = abnormal
			case web_server_manager.Normal:
				s = normal
			}
			return s
		},
	}
}

func (s *webServerMonitor) ShowVersion() string {
	s.Version = s.showVersionFunc()
	return s.Version
}

func (s *webServerMonitor) DisplayStatus() status {
	s.Status = s.displayStatusFunc()
	return s.Status
}

func (s *webServerMonitor) setVersion(v string) {
	s.Version = v
}

func (s *webServerMonitor) setStatus(status status) {
	s.Status = status
}
