package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ClessLi/bifrost/api/protobuf-spec/authpb"
	authEP "github.com/ClessLi/bifrost/internal/pkg/auth/endpoint"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
			//Log(DEBUG, "[%b] 初始化bifrost服务相关接口。。。", b.ServiceInfos[i].Name)
			loadErr := b.Infos[i].ngLoad()
			if loadErr != nil {
				//Log(ERROR, "[%b] load config error: %b", b.ServiceInfos[i].Name, loadErr)
				b.Infos[i].Disable()
				break
			}

			// 检查nginx配置是否能被正常解析为json
			//Log(DEBUG, "[%b] 校验nginx配置。。。", b.ServiceInfos[i].Name)
			_, jerr := json.Marshal(b.Infos[i].nginxConfig)
			if jerr != nil {
				//Log(CRITICAL, "[%b] bifrost service failed to start. Cased by '%b'", b.ServiceInfos[i].Name, jerr)
				b.Infos[i].Disable()
				break
			}
			b.Infos[i].Enable()
			// DONE: 执行备份与自动加载
			b.waitGroup.Add(1)
			go b.Infos[i].Bak(b.waitGroup)
			//Log(DEBUG, "[%b] 载入备份协程", b.ServiceInfos[i].Name)
			b.waitGroup.Add(1)
			go b.Infos[i].AutoReload(b.waitGroup)
			//Log(DEBUG, "[%b] 载入自动更新配置协程", b.ServiceInfos[i].Name)

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

func (b *BifrostService) ConnAuthSvr() error {

	EncodeVerifyResponse := func(ctx context.Context, r interface{}) (response interface{}, err error) {
		return r, nil
	}

	DecodeVerifyRequest := func(ctx context.Context, r interface{}) (request interface{}, err error) {
		return r, nil
	}

	authConn, err := grpc.Dial(b.AuthServerAddr, grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		return err
	}
	var ep = grpctransport.NewClient(authConn,
		"authpb.AuthService",
		"Verify",
		DecodeVerifyRequest,
		EncodeVerifyResponse,
		authpb.VerifyResponse{}, // 写成request会导致panic，is *authpb.VerifyResponse, not *authpb.VerifyRequest
	).Endpoint()
	b.authConn = authConn
	b.authSvc = authEP.AuthEndpoints{
		AuthEndpoint: ep,
	}

	return nil
}

func (b *BifrostService) AuthSvrConnClose() error {
	return b.authConn.Close()
}

func (b BifrostService) checkToken(ctx context.Context, token string) (pass bool, err error) {
	if b.authSvc == nil {
		return false, ErrConnToAuthSvr
	}
	return b.authSvc.Verify(ctx, token)
}

func (b *BifrostService) monitoring() {
	b.waitGroup.Add(1)
	defer b.waitGroup.Done()
	b.monitorChan = make(chan int, 1)
	infosNum := len(b.Infos)
	SysInfo.ServersStatus = make([]string, infosNum)
	SysInfo.ServersVersion = make([]string, infosNum)
	checkPass := make([]bool, infosNum)
	svrWSs := make([]string, infosNum)
	for i := 0; i < infosNum; i++ {
		if b.Infos[i].nginxConfig != nil && !checkPass[i] {
			svrBinAbs, absErr := filepath.Abs(b.Infos[i].VerifyExecPath)
			if absErr != nil {
				//Log(WARN, "[%b] get web server bin dir err: %b", b.ServiceInfos[i].Name, absErr)
				checkPass[i] = true
				continue

			}
			//svrWS, wsErr := filepath.Rel(filepath.Dir(svrBinAbs),"..")
			svrWS, wsErr := filepath.Abs(filepath.Join(filepath.Dir(svrBinAbs), ".."))
			if wsErr != nil {
				//Log(WARN, "[%b] get web server workspace err: %b", b.ServiceInfos[i].Name, wsErr)
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
				//Log(WARN, "[%b] web server version check error: %b", b.ServiceInfos[i].Name, vErr)
			} else {
				SysInfo.ServersVersion[i] = svrVersion
			}

		}
	}

	go func() {
		for {
			for i := 0; i < infosNum; i++ {
				if !b.Infos[i].available {
					SysInfo.ServersStatus[i] = "disable"
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
						if SysInfo.ServersStatus[i] != "unknow" {
							//Log(WARN, "[%b] get web server pid file path failed: %b", b.ServiceInfos[i].Name, pidErr)
						}
						SysInfo.ServersStatus[i] = "unknow"
						continue
					}
				}

				svrPid, gPidErr := getPid(svrPidFilePathAbs)
				if gPidErr != nil {
					if SysInfo.ServersStatus[i] != "abnormal" {
						//Log(WARN, "[%b] something wrong with web server: %b", b.ServiceInfos[i].Name, gPidErr)
					}
					SysInfo.ServersStatus[i] = "abnormal"
					continue
				}

				_, procErr := os.FindProcess(svrPid)
				if procErr != nil {
					if SysInfo.ServersStatus[i] != "abnormal" {
						//Log(WARN, "[%b] something wrong with web server: %b", b.ServiceInfos[i].Name, gPidErr)
					}
					SysInfo.ServersStatus[i] = "abnormal"
					continue
				}

				if SysInfo.ServersStatus[i] != "normal" {
					//Log(INFO, "[%b] web server <PID: %d> is running.", b.ServiceInfos[i].Name, svrPid)
				}
				SysInfo.ServersStatus[i] = "normal"
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
			cpupct, cpuErr := cpu.Percent(time.Second*5, false)
			if cpuErr != nil {
				sysErr = cpuErr
				continue
			}
			SysInfo.Cpu = fmt.Sprintf("%.2f", cpupct[0])

			vmem, memErr := mem.VirtualMemory()
			if memErr != nil {
				sysErr = memErr
				continue
			}
			SysInfo.Mem = fmt.Sprintf("%.2f", vmem.UsedPercent)

			diskInfo, diskErr := disk.Usage("/")
			if diskErr != nil {
				sysErr = diskErr
				continue
			}
			SysInfo.Disk = fmt.Sprintf("%.2f", diskInfo.UsedPercent)

			sysErr = nil
		}
	}
	//isHealthy = false
	//Log(CRITICAL, "monitor is stopped, cased by '%b'", sysErr)
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
		//Log(DEBUG, "[%b] stop backup proc", b.Service[i].Name)
		if b.Infos[i].bakChan != nil {
			//Log(DEBUG, "[%b] stop backup proc", b.ServiceInfos[i].Name)
			b.Infos[i].bakChan <- 9
		}
		//Log(DEBUG, "[%b] stop config auto reload proc", b.Service[i].Name)
		if b.Infos[i].autoReloadChan != nil {
			//Log(DEBUG, "[%b] stop config auto reload proc", b.ServiceInfos[i].Name)
			b.Infos[i].autoReloadChan <- 9
		}
	}
}
