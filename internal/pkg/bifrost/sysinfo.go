package bifrost

import (
	"bytes"
	"fmt"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type systemInfo struct {
	// DONE: 添加web服务版本信息、web服务状态信息(README.md需调整相关接口文档)
	OS             string   `json:"system"`
	Time           string   `json:"time"`
	Cpu            string   `json:"cpu"`
	Mem            string   `json:"mem"`
	Disk           string   `json:"disk"`
	ServersStatus  []string `json:"servers_status"`
	ServersVersion []string `json:"servers_version"`
	BifrostVersion string   `json:"bifrost_version"`
}

func monitoring(signal chan int) {
	svrsNum := len(BifrostConf.WebServerInfo.Servers)
	si.ServersStatus = make([]string, svrsNum)
	si.ServersVersion = make([]string, svrsNum)
	checkPass := make([]bool, svrsNum)
	svrWSs := make([]string, svrsNum)
	for i := 0; i < svrsNum; i++ {
		if BifrostConf.WebServerInfo.Servers[i].nginxConfig != nil && !checkPass[i] {
			svrBinAbs, absErr := filepath.Abs(BifrostConf.WebServerInfo.Servers[i].VerifyExecPath)
			if absErr != nil {
				Log(WARN, "[%s] get web server bin dir err: %s", BifrostConf.WebServerInfo.Servers[i].Name, absErr)
				checkPass[i] = true
				continue

			}
			//svrWS, wsErr := filepath.Rel(filepath.Dir(svrBinAbs),"..")
			svrWS, wsErr := filepath.Abs(filepath.Join(filepath.Dir(svrBinAbs), ".."))
			if wsErr != nil {
				Log(WARN, "[%s] get web server workspace err: %s", BifrostConf.WebServerInfo.Servers[i].Name, wsErr)
				checkPass[i] = true
				continue
			}
			svrWSs[i] = svrWS
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
				Log(WARN, "[%s] web server version check error: %s", BifrostConf.WebServerInfo.Servers[i].Name, vErr)
			} else {
				si.ServersVersion[i] = svrVersion
			}

		}
	}

	go func() {
		for {
			for i := 0; i < svrsNum; i++ {
				svrPidFilePath := "logs/nginx.pid"
				svrPidFileKey, ok := BifrostConf.WebServerInfo.Servers[i].nginxConfig.QueryByKeywords(svrPidFileKW).(*nginx.Key)
				if ok && svrPidFileKey != nil {
					svrPidFilePath = svrPidFileKey.Value
				}

				svrPidFilePathAbs := svrPidFilePath
				if !filepath.IsAbs(svrPidFilePath) {
					var pidErr error
					svrPidFilePathAbs, pidErr = filepath.Abs(filepath.Join(svrWSs[i], svrPidFilePath))
					if pidErr != nil {
						if si.ServersStatus[i] != "unknow" {
							Log(WARN, "[%s] get web server pid file path failed: %s", BifrostConf.WebServerInfo.Servers[i].Name, pidErr)
						}
						si.ServersStatus[i] = "unknow"
						continue
					}
				}

				svrPid, gPidErr := getPid(svrPidFilePathAbs)
				if gPidErr != nil {
					if si.ServersStatus[i] != "abnormal" {
						Log(WARN, "[%s] something wrong with web server: %s", BifrostConf.WebServerInfo.Servers[i].Name, gPidErr)
					}
					si.ServersStatus[i] = "abnormal"
					continue
				}

				_, procErr := os.FindProcess(svrPid)
				if procErr != nil {
					if si.ServersStatus[i] != "abnormal" {
						Log(WARN, "[%s] something wrong with web server: %s", BifrostConf.WebServerInfo.Servers[i].Name, gPidErr)
					}
					si.ServersStatus[i] = "abnormal"
					continue
				}

				if si.ServersStatus[i] != "normal" {
					Log(INFO, "[%s] web server <PID: %d> is running.", BifrostConf.WebServerInfo.Servers[i].Name, svrPid)
				}
				si.ServersStatus[i] = "normal"
			}

			time.Sleep(1 * time.Minute)
		}
	}()

	var sysErr error
	for sysErr == nil {

		select {
		case s := <-signal: // 获取管道传入信号
			if s == 9 { // 为9时，停止监控
				Log(NOTICE, "monitor has been finished.")
				return
			}
		default:
			cpupct, cpuErr := cpu.Percent(time.Second*5, false)
			if cpuErr != nil {
				sysErr = cpuErr
				continue
			}
			si.Cpu = fmt.Sprintf("%.2f", cpupct[0])

			vmem, memErr := mem.VirtualMemory()
			if memErr != nil {
				sysErr = memErr
				continue
			}
			si.Mem = fmt.Sprintf("%.2f", vmem.UsedPercent)

			diskInfo, diskErr := disk.Usage("/")
			if diskErr != nil {
				sysErr = diskErr
				continue
			}
			si.Disk = fmt.Sprintf("%.2f", diskInfo.UsedPercent)

			sysErr = nil
		}
	}
	isHealthy = false
	Log(CRITICAL, "monitor is stopped, cased by '%s'", sysErr)
	return
}

func sysInfo(c *gin.Context) {
	status := "unkown"
	var message interface{} = "null"
	h := gin.H{
		"status":  &status,
		"message": &message,
	}

	// 获取接口传参
	strToken, hasToken := c.GetQuery("token")
	if !hasToken {
		status = "failed"
		message = "Token cannot be empty"
		Log(NOTICE, "[%s] token verify failed, message is: '%s'", c.ClientIP(), message)
		c.JSON(http.StatusBadRequest, &h)
		return
	}

	// 校验token
	_, err := verifyAction(strToken)
	if err != nil {
		//c.String(http.StatusNotFound, err.Error())
		status = "failed"
		message = err.Error()
		Log(NOTICE, "[%s] Verified failed", c.ClientIP())
		c.JSON(http.StatusNotFound, &h)
		return
	}

	status = "success"
	si.Time = time.Now().Format("2006/01/02 15:04:05")
	message = si
	c.JSON(http.StatusOK, &h)

}
