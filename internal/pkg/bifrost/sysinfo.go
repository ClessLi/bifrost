package bifrost

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"net/http"
	"time"
)

type systemInfo struct {
	// TODO: 添加web服务版本信息、web服务状态信息(README.md需调整相关接口文档)
	OS             string `json:"system"`
	Time           string `json:"time"`
	Cpu            string `json:"cpu"`
	Mem            string `json:"mem"`
	Disk           string `json:"disk"`
	BifrostVersion string `json:"bifrost_version"`
}

func monitoring(signal chan int) {
	var sysErr error
	for sysErr == nil {

		select {
		case s := <-signal: // 获取管道传入信号
			if s == 9 { // 为9时，停止备份
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
