package bifrost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	ngJson "github.com/ClessLi/bifrost/pkg/json/nginx"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	ngStatistics "github.com/ClessLi/bifrost/pkg/statistics/nginx"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Service, bifrost配置文件对象中web服务器信息结构体，定义管控的web服务器配置文件相关信息
type Service struct {
	ListenPort   int            `yaml:"listenPort"`
	ChunckSize   int            `yaml:"chunkSize"`
	ServiceInfos []*ServiceInfo `yaml:"infos,flow"`
	monitorChan  chan int
	waitGroup    *sync.WaitGroup
}

func (s Service) ViewConfig(req *bifrostpb.OperateRequest, qs bifrostpb.OperationService_ViewConfigServer) (err error) {
	_, err = Verify(qs.Context(), req.Token)
	if err != nil {
		return qs.Send(&bifrostpb.OperateResponse{Err: err.Error()})
	}
	info := s.getInfo(req.Location)
	if info != nil {
		//configStr := info.nginxConfig.String()

		for _, str := range info.nginxConfig.String() {
			b := []byte(str)
			for seek := 0; seek < len(b); seek += s.ChunckSize {
				if len(b) <= seek+BifrostConf.Service.ChunckSize {
					Log(DEBUG, "[%s] send config bytes from %d to end(%d)", info.Name, seek, len(b)-1)
					err = qs.Send(&bifrostpb.OperateResponse{Ret: b[seek:]})
				} else {
					Log(DEBUG, "[%s] send config bytes from %d to %d", info.Name, seek, seek+s.ChunckSize)
					err = qs.Send(&bifrostpb.OperateResponse{Ret: b[seek : seek+s.ChunckSize]})

				}
				if err != nil {
					return err
				}
			}
		}
		return err
	}
	return qs.Send(&bifrostpb.OperateResponse{Err: bifrostpb.LocationErr})
}

func (s Service) GetConfig(req *bifrostpb.OperateRequest, qs bifrostpb.OperationService_GetConfigServer) (err error) {
	_, err = Verify(qs.Context(), req.GetToken())
	if err != nil {
		return qs.Send(&bifrostpb.ConfigResponse{Err: err.Error()})
	}

	info := s.getInfo(req.Location)
	if info != nil {
		jdata, err := json.Marshal(info.nginxConfig)
		if err != nil {
			return qs.Send(&bifrostpb.ConfigResponse{Err: err.Error()})
		}

		for seek := 0; seek < len(jdata); seek += s.ChunckSize {
			if len(jdata) <= seek+BifrostConf.Service.ChunckSize {
				return qs.Send(&bifrostpb.ConfigResponse{Ret: &bifrostpb.Config{JData: jdata[seek:], Path: info.ConfPath}})
			} else {
				err = qs.Send(&bifrostpb.ConfigResponse{Ret: &bifrostpb.Config{JData: jdata[seek : seek+s.ChunckSize], Path: info.ConfPath}})
				if err != nil {
					return qs.Send(&bifrostpb.ConfigResponse{Err: err.Error()})
				}
			}
		}
		return nil
	}
	return qs.Send(&bifrostpb.ConfigResponse{Err: fmt.Sprintf("location(%s) request error", req.Location)})
}

func (s *Service) UpdateConfig(qs bifrostpb.OperationService_UpdateConfigServer) (err error) {
	defer func() {
		if err != nil {
			_ = qs.SendAndClose(&bifrostpb.OperateResponse{Err: err.Error()})
		}
	}()
	// buffer 得初始化空字节切片，否则会存在指定容量len数的初始化字节元素在切片内，在做buffer.Write时只会在之后添加。可指定初始化容积cap
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	defer buffer.Reset()
	var info *ServiceInfo
	ip := ""
	Log(DEBUG, "Read and parse the UpdateConfig API request data")
	for {
		in, err := qs.Recv()
		if err == io.EOF {
			if info != nil {
				err = nil
				break
			} else {
				err = fmt.Errorf(bifrostpb.EmptyConfigErr)
				return err
			}
		}

		if err != nil {
			return err
		}

		ip, err = Verify(qs.Context(), in.GetToken())
		if err != nil {
			return err
		}
		if info == nil {
			info = s.getInfo(in.Location)
		}

		if info != nil {
			if in.Location != info.Name {
				err = fmt.Errorf(bifrostpb.LocationErr)
				return err
			}
			if in.Req.Path != info.ConfPath {
				err = fmt.Errorf(bifrostpb.MainConfigPathErr)
				return err
			}
			buffer.Write(in.Req.JData)
		} else {
			err = fmt.Errorf("location(%s) request error", in.Location)
			return err
		}
	}

	// update config
	//fmt.Println("获取web服务配置校验二进制文件路径")
	verifyBin, err := filepath.Abs(info.VerifyExecPath)
	if err != nil {
		Log(CRITICAL, "[%s] %s detailed error: %s", info.Name, bifrostpb.ValidationNotExist, err)
		err = fmt.Errorf(bifrostpb.ValidationNotExist)
		return err
	}

	// check config
	if buffer.Len() > 0 {
		newConfig, err := ngJson.Unmarshal(buffer.Bytes())
		if err != nil {
			Log(DEBUG, "[%s] [%s] buffer detail: %s", info.Name, ip, buffer.String())
			Log(WARN, "[%s] [%s] %s detailed error: %s", info.Name, ip, bifrostpb.ConfigUnmarshalErr, err)
			return err
		}

		// delete old config

		err = nginx.Delete(info.nginxConfig)
		message := ""
		if err != nil {
			message = fmt.Sprintf("Delete nginx ng failed. <%s>", err)
			Log(ERROR, "[%s] [%s] %s", info.Name, ip, message)
			return err
		}

		Log(INFO, "[%s] Deleted old nginx config.", info.Name)
		Log(INFO, "[%s] Verify new nginx config.", info.Name)
		newCaches, err := nginx.SaveWithCheck(newConfig, verifyBin)
		// roll back
		if err != nil {
			Log(DEBUG, "[%s] Roll back to old nginx config.", info.Name)
			message = fmt.Sprintf("Nginx ng verify failed. <%s>", err)
			Log(WARN, "[%s] %s", info.Name, message)

			Log(INFO, "[%s] Delete new nginx ng.", info.Name)
			var rollErr error
			rollErr = nginx.Delete(newConfig)
			if rollErr != nil {
				Log(ERROR, "[%s] Delete new nginx ng failed. <%s>", info.Name, err)
				message = "New nginx config verify failed. And delete new nginx config failed."
				return rollErr
			}

			Log(INFO, "[%s] Rollback nginx ng.", info.Name)
			_, rollErr = nginx.Save(info.nginxConfig)
			if rollErr != nil {
				Log(CRITICAL, "[%s] Nginx ng rollback failed. <%s>", info.Name, err)
				message = "New nginx config verify failed. And nginx config rollback failed."
				return rollErr
			}

			return err
		}
		info.confCaches = newCaches
		info.nginxConfig = newConfig
		info.ConfPath = newConfig.Value

		Log(NOTICE, "[%s] [%s] Nginx Config saved successfully", info.Name, ip)
		return qs.SendAndClose(&bifrostpb.OperateResponse{Err: bifrostpb.UpdateSuccess})
	} else {
		Log(WARN, "[%s] [%s] %s", info.Name, ip, bifrostpb.EmptyConfigErr)
		err = fmt.Errorf(bifrostpb.EmptyConfigErr)
		return err
	}

}

func (s Service) ViewStatistics(req *bifrostpb.OperateRequest, qs bifrostpb.OperationService_ViewStatisticsServer) (err error) {
	defer func() {
		if err != nil {
			_ = qs.Send(&bifrostpb.OperateResponse{Err: err.Error()})
		}
	}()

	_, err = Verify(qs.Context(), req.Token)
	if err != nil {
		return
	}
	info := s.getInfo(req.Location)
	if info != nil {
		httpServersNum, httpServers := ngStatistics.HTTPServers(info.nginxConfig)
		httpPorts := ngStatistics.HTTPPorts(info.nginxConfig)
		streamServersNum, streamPorts := ngStatistics.StreamServers(info.nginxConfig)
		jData := struct {
			HttpSvrsNum   int              `json:"http_svrs_num"`
			HttpSvrs      map[string][]int `json:"http_svrs"`
			HttpPorts     []int            `json:"http_ports"`
			StreamSvrsNum int              `json:"stream_svrs_num"`
			StreamPorts   []int            `json:"stream_ports"`
		}{HttpSvrsNum: httpServersNum, HttpSvrs: httpServers, HttpPorts: httpPorts, StreamSvrsNum: streamServersNum, StreamPorts: streamPorts}
		jBytes, err := json.Marshal(jData)
		if err != nil {
			return err
		}
		for seek := 0; seek < len(jBytes); seek += s.ChunckSize {
			if len(jBytes) <= seek+s.ChunckSize {
				err = qs.Send(&bifrostpb.OperateResponse{Ret: jBytes[seek:]})
			} else {
				err = qs.Send(&bifrostpb.OperateResponse{Ret: jBytes[seek : seek+s.ChunckSize]})

			}
			if err != nil {
				return err
			}
		}
		return err
	}
	err = fmt.Errorf(bifrostpb.LocationErr)
	return err
}

func (s Service) Status(req *bifrostpb.OperateRequest, qs bifrostpb.OperationService_StatusServer) (err error) {
	defer func() {
		if err != nil {
			_ = qs.Send(&bifrostpb.OperateResponse{Err: err.Error()})
		}
	}()

	_, err = Verify(qs.Context(), req.Token)
	if err != nil {
		return
	}

	sysInfo.Time = time.Now().In(nginx.TZ).Format("2006/01/02 15:04:05")
	jBytes, err := json.Marshal(sysInfo)
	if err != nil {
		return
	}
	for seek := 0; seek < len(jBytes); seek += s.ChunckSize {
		if len(jBytes) <= seek+s.ChunckSize {
			err = qs.Send(&bifrostpb.OperateResponse{Ret: jBytes[seek:]})
		} else {
			err = qs.Send(&bifrostpb.OperateResponse{Ret: jBytes[seek : seek+s.ChunckSize]})

		}
		if err != nil {
			return
		}
	}
	return
}

func (s *Service) Run() {
	s.waitGroup = new(sync.WaitGroup)
	for i := 0; i < len(BifrostConf.Service.ServiceInfos); i++ {
		switch s.ServiceInfos[i].Type {
		case NGINX:
			Log(DEBUG, "[%s] 初始化bifrost服务相关接口。。。", s.ServiceInfos[i].Name)
			loadErr := s.ServiceInfos[i].ngLoad()
			if loadErr != nil {
				Log(ERROR, "[%s] load config error: %s", s.ServiceInfos[i].Name, loadErr)
				s.ServiceInfos[i].Disable()
				break
			}

			// 检查nginx配置是否能被正常解析为json
			Log(DEBUG, "[%s] 校验nginx配置。。。", s.ServiceInfos[i].Name)
			_, jerr := json.Marshal(s.ServiceInfos[i].nginxConfig)
			if jerr != nil {
				Log(CRITICAL, "[%s] bifrost service failed to start. Cased by '%s'", s.ServiceInfos[i].Name, jerr)
				s.ServiceInfos[i].Disable()
				break
			}
			s.ServiceInfos[i].Enable()
			// DONE: 执行备份与自动加载
			s.waitGroup.Add(1)
			go s.ServiceInfos[i].Bak(s.waitGroup)
			Log(DEBUG, "[%s] 载入备份协程", s.ServiceInfos[i].Name)
			s.waitGroup.Add(1)
			go s.ServiceInfos[i].AutoReload(s.waitGroup)
			Log(DEBUG, "[%s] 载入自动更新配置协程", s.ServiceInfos[i].Name)

		case HTTPD:
			// TODO: apache httpd配置解析器
			continue
		default:
			continue
		}
	}
	// 监控系统信息
	go s.monitoring()
}

func (s *Service) monitoring() {
	s.waitGroup.Add(1)
	defer s.waitGroup.Done()
	s.monitorChan = make(chan int, 1)
	infosNum := len(s.ServiceInfos)
	sysInfo.ServersStatus = make([]string, infosNum)
	sysInfo.ServersVersion = make([]string, infosNum)
	checkPass := make([]bool, infosNum)
	svrWSs := make([]string, infosNum)
	for i := 0; i < infosNum; i++ {
		if s.ServiceInfos[i].nginxConfig != nil && !checkPass[i] {
			svrBinAbs, absErr := filepath.Abs(s.ServiceInfos[i].VerifyExecPath)
			if absErr != nil {
				Log(WARN, "[%s] get web server bin dir err: %s", s.ServiceInfos[i].Name, absErr)
				checkPass[i] = true
				continue

			}
			//svrWS, wsErr := filepath.Rel(filepath.Dir(svrBinAbs),"..")
			svrWS, wsErr := filepath.Abs(filepath.Join(filepath.Dir(svrBinAbs), ".."))
			if wsErr != nil {
				Log(WARN, "[%s] get web server workspace err: %s", s.ServiceInfos[i].Name, wsErr)
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
				Log(WARN, "[%s] web server version check error: %s", s.ServiceInfos[i].Name, vErr)
			} else {
				sysInfo.ServersVersion[i] = svrVersion
			}

		}
	}

	go func() {
		for {
			for i := 0; i < infosNum; i++ {
				if !s.ServiceInfos[i].available {
					sysInfo.ServersStatus[i] = "disable"
					continue
				}
				svrPidFilePath := "logs/nginx.pid"
				svrPidFileKey, ok := s.ServiceInfos[i].nginxConfig.QueryByKeywords(svrPidFileKW).(*nginx.Key)
				if ok && svrPidFileKey != nil {
					svrPidFilePath = svrPidFileKey.Value
				}

				svrPidFilePathAbs := svrPidFilePath
				if !filepath.IsAbs(svrPidFilePath) {
					var pidErr error
					svrPidFilePathAbs, pidErr = filepath.Abs(filepath.Join(svrWSs[i], svrPidFilePath))
					if pidErr != nil {
						if sysInfo.ServersStatus[i] != "unknow" {
							Log(WARN, "[%s] get web server pid file path failed: %s", s.ServiceInfos[i].Name, pidErr)
						}
						sysInfo.ServersStatus[i] = "unknow"
						continue
					}
				}

				svrPid, gPidErr := getPid(svrPidFilePathAbs)
				if gPidErr != nil {
					if sysInfo.ServersStatus[i] != "abnormal" {
						Log(WARN, "[%s] something wrong with web server: %s", s.ServiceInfos[i].Name, gPidErr)
					}
					sysInfo.ServersStatus[i] = "abnormal"
					continue
				}

				_, procErr := os.FindProcess(svrPid)
				if procErr != nil {
					if sysInfo.ServersStatus[i] != "abnormal" {
						Log(WARN, "[%s] something wrong with web server: %s", s.ServiceInfos[i].Name, gPidErr)
					}
					sysInfo.ServersStatus[i] = "abnormal"
					continue
				}

				if sysInfo.ServersStatus[i] != "normal" {
					Log(INFO, "[%s] web server <PID: %d> is running.", s.ServiceInfos[i].Name, svrPid)
				}
				sysInfo.ServersStatus[i] = "normal"
			}

			time.Sleep(1 * time.Minute)
		}
	}()

	var sysErr error
	for sysErr == nil {

		select {
		case s := <-s.monitorChan: // 获取管道传入信号
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
			sysInfo.Cpu = fmt.Sprintf("%.2f", cpupct[0])

			vmem, memErr := mem.VirtualMemory()
			if memErr != nil {
				sysErr = memErr
				continue
			}
			sysInfo.Mem = fmt.Sprintf("%.2f", vmem.UsedPercent)

			diskInfo, diskErr := disk.Usage("/")
			if diskErr != nil {
				sysErr = diskErr
				continue
			}
			sysInfo.Disk = fmt.Sprintf("%.2f", diskInfo.UsedPercent)

			sysErr = nil
		}
	}
	//isHealthy = false
	Log(CRITICAL, "monitor is stopped, cased by '%s'", sysErr)
	return
}

// getInfo, ServiceInfo获取Service指针对象的方法
// 参数:
//     name: bifrost服务名
func (s *Service) getInfo(name string) *ServiceInfo {
	for idx := range s.ServiceInfos {
		if s.ServiceInfos[idx].Name == name {
			return s.ServiceInfos[idx]
		}
	}
	return nil
}

// killCoroutines, ServiceInfo关闭协程任务的方法
func (s *Service) killCoroutines() {
	defer s.waitGroup.Wait()
	if s.monitorChan != nil {
		Log(DEBUG, "stop monitor proc")
		s.monitorChan <- 9
	}
	for i := 0; i < len(s.ServiceInfos); i++ {
		//Log(DEBUG, "[%s] stop backup proc", s.Service[i].Name)
		if s.ServiceInfos[i].bakChan != nil {
			Log(DEBUG, "[%s] stop backup proc", s.ServiceInfos[i].Name)
			s.ServiceInfos[i].bakChan <- 9
		}
		//Log(DEBUG, "[%s] stop config auto reload proc", s.Service[i].Name)
		if s.ServiceInfos[i].autoReloadChan != nil {
			Log(DEBUG, "[%s] stop config auto reload proc", s.ServiceInfos[i].Name)
			s.ServiceInfos[i].autoReloadChan <- 9
		}
	}
}
