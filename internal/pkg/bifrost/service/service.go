package service

import (
	"encoding/json"
	"github.com/ClessLi/bifrost/pkg/client/auth"
	ngStatistics "github.com/ClessLi/bifrost/pkg/statistics/nginx"
	"golang.org/x/net/context"
	"sync"
	"time"
)

type Service interface {
	ViewConfig(ctx context.Context, token, svrName string) (data []byte, err error)
	GetConfig(ctx context.Context, token, srvName string) (jsonData []byte, err error)
	UpdateConfig(ctx context.Context, token, svrName string, jsonData []byte) (data []byte, err error)
	ViewStatistics(ctx context.Context, token, svrName string) (jsonData []byte, err error)
	Status(ctx context.Context, token string) (jsonData []byte, err error)
	// TODO: WatchLog暂用锁机制，一个日志文件仅允许一个终端访问
	//WatchLog(ctx context.Context, token, svrName, logName string, dataChan chan<- []byte, signal <-chan int) error
	WatchLog(ctx context.Context, token, svrName, logName string) (logWatcher *LogWatcher, err error)
}

// BifrostService, bifrost配置文件对象中web服务器信息结构体，定义管控的web服务器配置文件相关信息
type BifrostService struct {
	Port           uint16  `yaml:"Port"`
	ChunckSize     int     `yaml:"ChunkSize"`
	AuthServerAddr string  `yaml:"AuthServerAddr"`
	Infos          []*Info `yaml:"Infos,flow"`
	monitorChan    chan int
	authSvcCli     *auth.Client
	waitGroup      *sync.WaitGroup
}

func (b *BifrostService) ViewConfig(ctx context.Context, token, svrName string) (data []byte, err error) {
	var pass bool
	pass, err = b.checkToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if !pass {
		return nil, UnknownErrCheckToken
	}

	var info *Info
	info, err = b.getInfo(svrName)
	if err != nil {
		return nil, err
	}
	for _, s := range info.nginxConfig.String() {
		data = append(data, []byte(s)...)
	}
	if data == nil {
		err = ErrDataNotParsed
		return nil, err
	}
	return data, err
}

func (b *BifrostService) GetConfig(ctx context.Context, token, svrName string) (jsonData []byte, err error) {
	var pass bool
	pass, err = b.checkToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if !pass {
		return nil, UnknownErrCheckToken
	}

	var info *Info
	info, err = b.getInfo(svrName)
	if err != nil {
		return nil, err
	}

	return json.Marshal(info.nginxConfig)
}

func (b *BifrostService) UpdateConfig(ctx context.Context, token, svrName string, jsonData []byte) (data []byte, err error) {
	var pass bool
	pass, err = b.checkToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if !pass {
		return nil, UnknownErrCheckToken
	}

	var info *Info
	info, err = b.getInfo(svrName)
	if err != nil {
		return nil, err
	}

	return info.update(jsonData)
}

func (b *BifrostService) ViewStatistics(ctx context.Context, token, svrName string) (jsonData []byte, err error) {
	var pass bool
	pass, err = b.checkToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if !pass {
		return nil, UnknownErrCheckToken
	}

	var info *Info
	info, err = b.getInfo(svrName)
	if err != nil {
		return nil, err
	}

	httpServersNum, httpServers := ngStatistics.HTTPServers(info.nginxConfig)
	httpPorts := ngStatistics.HTTPPorts(info.nginxConfig)
	streamServersNum, streamPorts := ngStatistics.StreamServers(info.nginxConfig)
	statistics := struct {
		HttpSvrsNum   int              `json:"http_svrs_num"`
		HttpSvrs      map[string][]int `json:"http_svrs"`
		HttpPorts     []int            `json:"http_ports"`
		StreamSvrsNum int              `json:"stream_svrs_num"`
		StreamPorts   []int            `json:"stream_ports"`
	}{HttpSvrsNum: httpServersNum, HttpSvrs: httpServers, HttpPorts: httpPorts, StreamSvrsNum: streamServersNum, StreamPorts: streamPorts}
	return json.Marshal(statistics)
}

func (b *BifrostService) Status(ctx context.Context, token string) (jsonData []byte, err error) {
	var pass bool
	pass, err = b.checkToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if !pass {
		return nil, UnknownErrCheckToken
	}

	// TODO: SysInfo lock mechanism
	return json.Marshal(SysInfo)
}

func (b *BifrostService) WatchLog(ctx context.Context, token, svrName, logName string) (logWatcher *LogWatcher, err error) {
	//fmt.Println("svc接收到请求")
	//if dataChan == nil || signal == nil {
	//	return ErrChanNil
	//}
	// 认证
	//fmt.Println("svc认证请求")
	var pass bool
	pass, err = b.checkToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if !pass {
		return nil, UnknownErrCheckToken
	}

	// 获取info
	var info *Info
	info, err = b.getInfo(svrName)
	if err != nil {
		return nil, err
	}

	// 开始监控日志
	ticker := time.Tick(time.Second)
	timeout := time.After(time.Minute * 30)
	err = info.StartWatchLog(logName)
	if err != nil {
		return nil, err
	}

	logWatcher = NewLogWatcher()
	// 监听终止信号和每秒读取日志并发送
	//fmt.Println("监听终止信号及准备发送日志")
	go func() {
		for {
			select {
			case s := <-logWatcher.SignalC:
				if s == 9 {
					//fmt.Println("watch log stopping...")
					err = info.StopWatchLog(logName)
					//fmt.Println("watch log is stopped")
					if err != nil {
						logWatcher.ErrC <- err
					}
					return
				}
			case <-ticker:
				//fmt.Println("读取日志")
				data, err := info.WatchLog(logName)
				if err != nil {
					stopErr := info.StopWatchLog(logName)
					if stopErr != nil {
						err = stopErr
					}
					logWatcher.ErrC <- err
					return
				}
				if len(data) == 0 {
					break
				}
				select {
				case logWatcher.DataC <- data:
					// 日志推送后，客户端已经终止，handler日志推送阻断且发送了终止信号，由于日志推送阻断，接收终止信息被积压
					//fmt.Println("svc发送日志成功")
				case <-time.After(time.Second * 30):
					_ = info.StopWatchLog(logName)
					logWatcher.ErrC <- ErrDataSendingTimeout
					return
				}
			case <-timeout:
				logWatcher.ErrC <- ErrWatchLogTimeout
				return
			}
		}
	}()
	return logWatcher, nil
}

//func (b *BifrostService) WatchLog(ctx context.Context, token, svrName, logName string, dataChan chan<- []byte, signal <-chan int) (err error) {
//	//fmt.Println("svc接收到请求")
//	if dataChan == nil || signal == nil {
//		return ErrChanNil
//	}
//	// 认证
//	//fmt.Println("svc认证请求")
//	var pass bool
//	pass, err = b.checkToken(ctx, token)
//	if err != nil {
//		return err
//	}
//	if !pass {
//		return UnknownErrCheckToken
//	}
//
//	// 获取info
//	var info *Info
//	info, err = b.getInfo(svrName)
//	if err != nil {
//		return err
//	}
//
//	// 开始监控日志
//	ticker := time.Tick(time.Second)
//	timeout := time.After(time.Minute * 30)
//	err = info.StartWatchLog(logName)
//	if err != nil {
//		return err
//	}
//
//	// 监听终止信号和每秒读取日志并发送
//	//fmt.Println("监听终止信号及准备发送日志")
//	for {
//		select {
//		case s := <-signal:
//			if s == 9 {
//				//fmt.Println("watch log stopping...")
//				err = info.StopWatchLog(logName)
//				//fmt.Println("watch log is stopped")
//				return err
//			}
//		case <-ticker:
//			//fmt.Println("读取日志")
//			data, err := info.WatchLog(logName)
//			if err != nil {
//				stopErr := info.StopWatchLog(logName)
//				if stopErr != nil {
//					err = stopErr
//				}
//				return err
//			}
//			if len(data) == 0 {
//				break
//			}
//			select {
//			case dataChan <- data:
//				// 日志推送后，客户端已经终止，handler日志推送阻断且发送了终止信号，由于日志推送阻断，接收终止信息被积压
//				//fmt.Println("svc发送日志成功")
//			case <-time.After(time.Second * 30):
//				_ = info.StopWatchLog(logName)
//				return ErrDataSendingTimeout
//			}
//		case <-timeout:
//			return ErrWatchLogTimeout
//		}
//	}
//}

func (b BifrostService) GetPort() uint16 {
	return b.Port
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
