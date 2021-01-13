package bifrost

import (
	"errors"
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/config"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service"
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	"github.com/ClessLi/skirnir/pkg/discover"
	uuid "github.com/satori/go.uuid"
	"net"
	"os"
	"strings"
)

func getProc(path string) (*os.Process, error) {
	pid, pidErr := utils.GetPid(path)
	if pidErr != nil {
		return nil, pidErr
	}
	return os.FindProcess(pid)
}

func rmPidFile(path string) {
	rmPidFileErr := os.Remove(path)
	if rmPidFileErr != nil {
		utils.Logger.Error(rmPidFileErr.Error())
	}
	utils.Logger.Notice("bifrost.pid has been removed.")
}

// configCheck, 检查bifrost配置项是否完整
// 返回值:
//     错误
func configCheck() error {
	if BifrostConf == nil {
		return fmt.Errorf("bifrost config load error")
	}
	if len(BifrostConf.ServiceConfig.WebServerConfigInfos) == 0 {
		return fmt.Errorf("bifrost services config load error")
	}
	if BifrostConf.LogDir == "" {
		return fmt.Errorf("bifrost log config load error")
	}
	// 初始化服务信息配置
	if BifrostConf.ServiceConfig.Port == 0 {
		BifrostConf.ServiceConfig.Port = 12321
	}
	if BifrostConf.ServiceConfig.ChunckSize == 0 {
		BifrostConf.ServiceConfig.ChunckSize = 4194304
	}

	if BifrostConf.RAConfig != nil {
		if BifrostConf.RAConfig.Host == "" || BifrostConf.RAConfig.Port == 0 {
			BifrostConf.RAConfig = nil
		}
	}
	return nil
}

func registerToRA(errChan chan<- error) {
	if BifrostConf.RAConfig == nil {
		return
	}

	var err error
	discoveryClient, err = discover.NewKitConsulRegistryClient(BifrostConf.RAConfig.Host, BifrostConf.RAConfig.Port)
	if err != nil {
		utils.Logger.WarningF("Get Consul Client failed. Cased by: %s", err)
		errChan <- err
		return
	}

	svcName := "com.github.ClessLi.api.bifrost"
	//svcName := "bifrostpb.BifrostService"
	//svcName := "BifrostService"
	//svcName := "Health"

	instanceId = svcName + "-" + uuid.NewV4().String()

	instanceIP, err := externalIP()
	if err != nil {
		utils.Logger.WarningF("Failed to initialize service instance IP. Cased by: %s", err)
		errChan <- err
		return
	}
	instanceHost := instanceIP.String()

	if !discoveryClient.Register(svcName, instanceId, instanceHost, BifrostConf.ServiceConfig.Port, nil, config.KitLogger) {
		err = fmt.Errorf("register service %s failed", svcName)
		utils.Logger.Warning(err.Error())
		errChan <- err
		instanceId = ""
		return
	}
}

func deregisterToRA() {
	if discoveryClient != nil && !strings.EqualFold(instanceId, "") {
		if discoveryClient.DeRegister(instanceId, config.KitLogger) {
			utils.Logger.InfoF("bifrost service (instance ID is '%s') has been unregistered from RA '%s:%d'", instanceId, BifrostConf.RAConfig.Host, BifrostConf.RAConfig.Port)
		} else {
			utils.Logger.WarningF("bifrost service (instance ID is '%s') failed to deregister from RA '%s:%d'", instanceId, BifrostConf.RAConfig.Host, BifrostConf.RAConfig.Port)
		}
	}
}

func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("external ip failed")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func newService() service.Service {
	bifrostService := BifrostConf.ServiceConfig.BifrostServiceController.GetService()
	return service.NewService(bifrostService)
}

func bifrostServiceStart() error {
	return BifrostConf.ServiceConfig.BifrostServiceController.Start()
}

func bifrostServiceStop() error {
	return BifrostConf.ServiceConfig.BifrostServiceController.Stop()
}
