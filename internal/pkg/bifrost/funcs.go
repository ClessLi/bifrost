package bifrost

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/config"
	"github.com/ClessLi/skirnir/pkg/discover"
	"github.com/apsdehal/go-logger"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

// readFile, 读取文件函数
// 参数:
//     path: 文件路径字符串
// 返回值:
//     文件数据
//     错误
func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fd, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return fd, nil
}

// PathExists, 判断文件路径是否存在函数
// 参数:
//     path: 待判断的文件路径字符串
// 返回值:
//     true: 存在; false: 不存在
//     错误
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil || os.IsExist(err) {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, err
	} else {
		return false, nil
	}
}

// Log, 日志记录函数
// 参数:
//     level: 日志级别对象
//     message: 需记录的日志信息字符串
func Log(level logger.LogLevel, message string, a ...interface{}) {
	myLogger.Log(level, fmt.Sprintf(message, a...))
}

func getProc(path string) (*os.Process, error) {
	pid, pidErr := getPid(path)
	if pidErr != nil {
		return nil, pidErr
	}
	return os.FindProcess(pid)
}

func rmPidFile(path string) {
	rmPidFileErr := os.Remove(path)
	if rmPidFileErr != nil {
		Log(ERROR, rmPidFileErr.Error())
	}
	Log(NOTICE, "bifrost.pid has been removed.")
}

// getPid, 查询pid文件并返回pid
// 返回值:
//     pid
//     错误
func getPid(path string) (int, error) {
	// 判断pid文件是否存在
	if _, err := os.Stat(path); err == nil || os.IsExist(err) { // 存在
		// 读取pid文件
		pidBytes, readPidErr := readFile(path)
		if readPidErr != nil {
			Log(ERROR, readPidErr.Error())
			return -1, readPidErr
		}

		// 去除pid后边的换行符
		pidBytes = bytes.TrimRight(pidBytes, "\n")

		// 转码pid
		pid, toIntErr := strconv.Atoi(string(pidBytes))
		if toIntErr != nil {
			Log(ERROR, toIntErr.Error())
			return -1, toIntErr
		}

		return pid, nil
	} else { // 不存在
		return -1, procStatusNotRunning
	}
}

// configCheck, 检查bifrost配置项是否完整
// 返回值:
//     错误
func configCheck() error {
	if BifrostConf == nil {
		return fmt.Errorf("bifrost config load error")
	}
	if len(BifrostConf.Service.Infos) == 0 {
		return fmt.Errorf("bifrost services config load error")
	}
	if BifrostConf.LogDir == "" {
		return fmt.Errorf("bifrost log config load error")
	}
	// 初始化服务信息配置
	if BifrostConf.Service.Port == 0 {
		BifrostConf.Service.Port = 12321
	}
	if BifrostConf.Service.ChunckSize == 0 {
		BifrostConf.Service.ChunckSize = 4194304
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
		Log(WARN, "Get Consul Client failed. Cased by: %s", err)
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
		Log(WARN, "Failed to initialize service instance IP. Cased by: %s", err)
		errChan <- err
		return
	}
	instanceHost := instanceIP.String()

	if !discoveryClient.Register(svcName, instanceId, instanceHost, BifrostConf.Service.Port, nil, config.KitLogger) {
		err = fmt.Errorf("register service %s failed", svcName)
		Log(WARN, err.Error())
		errChan <- err
		instanceId = ""
		return
	}
}

func deregisterToRA() {
	if discoveryClient != nil && !strings.EqualFold(instanceId, "") {
		if discoveryClient.DeRegister(instanceId, config.KitLogger) {
			Log(INFO, "bifrost service (instance ID is '%s') has been unregistered from RA '%s:%d'", instanceId, BifrostConf.RAConfig.Host, BifrostConf.RAConfig.Port)
		} else {
			Log(WARN, "bifrost service (instance ID is '%s') failed to deregister from RA '%s:%d'", instanceId, BifrostConf.RAConfig.Host, BifrostConf.RAConfig.Port)
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
