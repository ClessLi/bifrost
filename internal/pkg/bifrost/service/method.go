package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ClessLi/bifrost/pkg/client/auth"
	ngLog "github.com/ClessLi/bifrost/pkg/log/nginx"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"golang.org/x/net/context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	svrPidFileKW = nginx.NewKeyWords(nginx.TypeKey, "pid", "*", false, false)
)

func (b *BifrostService) Run() {
	b.waitGroup = new(sync.WaitGroup)

	for i := 0; i < len(b.Infos); i++ {
		switch b.Infos[i].Type {
		case NGINX:
			//Log(DEBUG, "[%s] 初始化bifrost服务相关接口。。。", b.ServiceInfos[i].Name)
			loadErr := b.Infos[i].ngLoad()
			if loadErr != nil {
				fmt.Printf("[%s] load config error: %s\n", b.Infos[i].Name, loadErr)
				//Log(ERROR, "[%s] load config error: %s", b.ServiceInfos[i].Name, loadErr)
				b.Infos[i].Disable()
				break
			}

			// 检查nginx配置是否能被正常解析为json
			//Log(DEBUG, "[%s] 校验nginx配置。。。", b.ServiceInfos[i].Name)
			_, jerr := json.Marshal(b.Infos[i].nginxConfig)
			if jerr != nil {
				fmt.Printf("[%s] bifrost service failed to start. Cased by '%s'\n", b.Infos[i].Name, jerr)
				//Log(CRITICAL, "[%s] bifrost service failed to start. Cased by '%s'", b.ServiceInfos[i].Name, jerr)
				b.Infos[i].Disable()
				break
			}
			b.Infos[i].Enable()
			// DONE: 执行备份与自动加载
			b.waitGroup.Add(1)
			go b.Infos[i].Bak(b.waitGroup)
			//Log(DEBUG, "[%s] 载入备份协程", b.ServiceInfos[i].Name)
			b.waitGroup.Add(1)
			go b.Infos[i].AutoReload(b.waitGroup)
			//Log(DEBUG, "[%s] 载入自动更新配置协程", b.ServiceInfos[i].Name)
			b.Infos[i].nginxLog = ngLog.NewLog()
		case HTTPD:
			// TODO: apache httpd配置解析器
			continue
		default:
			continue
		}
	}
	// 监控系统信息
	go b.monitoring()
}

func (b *BifrostService) ConnAuthSvr() (err error) {
	b.authSvcCli, err = auth.NewClient(b.AuthServerAddr)
	if err != nil {
		return err
	}
	return nil
}

func (b *BifrostService) AuthSvrConnClose() error {
	return b.authSvcCli.Close()
}

func (b BifrostService) checkToken(ctx context.Context, token string) (pass bool, err error) {
	if b.authSvcCli == nil {
		return false, ErrConnToAuthSvr
	}
	return b.authSvcCli.Verify(ctx, token)
}

func (b *BifrostService) monitoring() {
	b.waitGroup.Add(1)
	defer b.waitGroup.Done()
	b.monitorChan = make(chan int, 1)
	infosNum := len(b.Infos)
	//SysInfo.ServersStatus = make([]string, infosNum)
	//SysInfo.ServersStatus = make(map[string]string, infosNum)
	SysInfo.StatusList = make([]ServerStatus, infosNum)
	//SysInfo.ServersVersion = make([]string, infosNum)
	//SysInfo.ServersVersion = make(map[string]string, infosNum)
	checkPass := make([]bool, infosNum)
	svrWSs := make([]string, infosNum)
	for i := 0; i < infosNum; i++ {
		SysInfo.StatusList[i] = newStatus(b.Infos[i].Name)
		if b.Infos[i].nginxConfig != nil && !checkPass[i] {
			svrBinAbs, absErr := filepath.Abs(b.Infos[i].VerifyExecPath)
			if absErr != nil {
				//Log(WARN, "[%s] get web server bin dir err: %s", b.ServiceInfos[i].Name, absErr)
				checkPass[i] = true
				continue

			}
			//svrWS, wsErr := filepath.Rel(filepath.Dir(svrBinAbs),"..")
			svrWS, wsErr := filepath.Abs(filepath.Join(filepath.Dir(svrBinAbs), ".."))
			if wsErr != nil {
				//Log(WARN, "[%s] get web server workspace err: %s", b.ServiceInfos[i].Name, wsErr)
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
				//Log(WARN, "[%s] web server version check error: %s", b.ServiceInfos[i].Name, vErr)
			} else {
				SysInfo.StatusList[i].setVersion(svrVersion)
			}

		}
	}

	go func() {
		for {
			for i := 0; i < infosNum; i++ {
				if !b.Infos[i].available {
					SysInfo.StatusList[i].setStatus(disable)
					continue
				}
				svrPidFilePath := "logs/nginx.pid"
				svrPidFileKey, ok := b.Infos[i].nginxConfig.QueryByKeywords(svrPidFileKW).(*nginx.Key)
				if ok && svrPidFileKey != nil {
					svrPidFilePath = svrPidFileKey.Value
				}

				svrPidFilePathAbs := svrPidFilePath
				if !filepath.IsAbs(svrPidFilePath) {
					var pidErr error
					svrPidFilePathAbs, pidErr = filepath.Abs(filepath.Join(svrWSs[i], svrPidFilePath))
					if pidErr != nil {
						//if SysInfo.StatusList[i].ServerStatus != "unknow" {
						//Log(WARN, "[%s] get web server pid file path failed: %s", b.ServiceInfos[i].Name, pidErr)
						//}
						SysInfo.StatusList[i].setStatus(unknown)
						continue
					}
				}

				svrPid, gPidErr := getPid(svrPidFilePathAbs)
				if gPidErr != nil {
					//if SysInfo.StatusList[i].ServerStatus != "abnormal" {
					//Log(WARN, "[%s] something wrong with web server: %s", b.ServiceInfos[i].Name, gPidErr)
					//}
					SysInfo.StatusList[i].setStatus(abnormal)
					continue
				}

				_, procErr := os.FindProcess(svrPid)
				if procErr != nil {
					//if SysInfo.StatusList[i].ServerStatus != "abnormal" {
					//Log(WARN, "[%s] something wrong with web server: %s", b.ServiceInfos[i].Name, gPidErr)
					//}
					SysInfo.StatusList[i].setStatus(abnormal)
					continue
				}

				//if SysInfo.StatusList[i].ServerStatus != "normal" {
				//Log(INFO, "[%s] web server <PID: %d> is running.", b.ServiceInfos[i].Name, svrPid)
				//}
				SysInfo.StatusList[i].setStatus(normal)
			}

			time.Sleep(1 * time.Minute)
		}
	}()

	var sysErr error
	for sysErr == nil {

		select {
		case s := <-b.monitorChan: // 获取管道传入信号
			if s == 9 { // 为9时，停止监控
				//Log(NOTICE, "monitor has been finished.")
				return
			}
		default:
			// cpu监控
			cpupct, cpuErr := cpu.Percent(time.Second*5, false)
			if cpuErr != nil {
				sysErr = cpuErr
				continue
			}
			SysInfo.Cpu = fmt.Sprintf("%.2f", cpupct[0])

			// 内存监控
			vmem, memErr := mem.VirtualMemory()
			if memErr != nil {
				sysErr = memErr
				continue
			}
			SysInfo.Mem = fmt.Sprintf("%.2f", vmem.UsedPercent)

			// 磁盘监控
			diskInfo, diskErr := disk.Usage("/")
			if diskErr != nil {
				sysErr = diskErr
				continue
			}
			SysInfo.Disk = fmt.Sprintf("%.2f", diskInfo.UsedPercent)

			// 监控时间戳
			SysInfo.Time = fmt.Sprintf("%s", time.Now().In(nginx.TZ).Format("2006/01/02 15:04:05"))
			sysErr = nil
		}
	}
	//isHealthy = false
	//Log(CRITICAL, "monitor is stopped, cased by '%s'", sysErr)
	return
}

// getInfo, ServiceInfo获取Service指针对象的方法
// 参数:
//     name: bifrost服务名
func (b BifrostService) getInfo(svrName string) (info *Info, err error) {
	for idx := range b.Infos {
		if b.Infos[idx].Name == svrName {
			return b.Infos[idx], nil
		}
	}
	return nil, ErrUnknownSvrName
}

// KillCoroutines, ServiceInfo关闭协程任务的方法
func (b *BifrostService) KillCoroutines() {
	defer b.waitGroup.Wait()
	if b.monitorChan != nil {
		//Log(DEBUG, "stop monitor proc")
		b.monitorChan <- 9
	}
	for i := 0; i < len(b.Infos); i++ {
		//Log(DEBUG, "[%s] stop backup proc", b.Service[i].Name)
		if b.Infos[i].bakChan != nil {
			//Log(DEBUG, "[%s] stop backup proc", b.ServiceInfos[i].Name)
			b.Infos[i].bakChan <- 9
		}
		//Log(DEBUG, "[%s] stop config auto reload proc", b.Service[i].Name)
		if b.Infos[i].autoReloadChan != nil {
			//Log(DEBUG, "[%s] stop config auto reload proc", b.ServiceInfos[i].Name)
			b.Infos[i].autoReloadChan <- 9
		}
	}
}
