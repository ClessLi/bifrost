package bifrost

import (
	"errors"
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/authentication"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/config"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/endpoint"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/logging"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/service"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/transport"
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	"github.com/ClessLi/bifrost/pkg/log/nginx"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration/parser"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/loader"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"path/filepath"
	"sync"
)

func getProc(path string) (*os.Process, error) {
	pid, pidErr := utils.GetPid(path)
	if pidErr != nil {
		return nil, pidErr
	}
	return os.FindProcess(pid)
}

//func rmPidFile(path string) {
//	rmPidFileErr := os.Remove(path)
//	if rmPidFileErr != nil {
//		utils.Logger.Error(rmPidFileErr.Error())
//	}
//	utils.Logger.Notice("bifrost.pid has been removed.")
//}

//// configCheck, 检查bifrost配置项是否完整
//// 返回值:
////     错误
//func configCheck() error {
//	if BifrostConf == nil {
//		return fmt.Errorf("bifrost config load error")
//	}
//	if len(BifrostConf.ServiceConfig.WebServerConfigInfos) == 0 {
//		return fmt.Errorf("bifrost services config load error")
//	}
//	if BifrostConf.LogDir == "" {
//		return fmt.Errorf("bifrost log config load error")
//	}
//	// 初始化服务信息配置
//	if BifrostConf.ServiceConfig.Port == 0 {
//		BifrostConf.ServiceConfig.Port = 12321
//	}
//	if BifrostConf.ServiceConfig.ChunckSize == 0 {
//		BifrostConf.ServiceConfig.ChunckSize = 4194304
//	}
//
//	if BifrostConf.RAConfig != nil {
//		if BifrostConf.RAConfig.Host == "" || BifrostConf.RAConfig.Port == 0 {
//			BifrostConf.RAConfig = nil
//		}
//	}
//	return nil
//}

//func registerToRA(errChan chan<- error) {
//	if BifrostConf.RAConfig == nil {
//		return
//	}
//
//	var err error
//	discoveryClient, err = discover.NewKitConsulRegistryClient(BifrostConf.RAConfig.Host, BifrostConf.RAConfig.Port)
//	if err != nil {
//		utils.Logger.WarningF("Get Consul Client failed. Cased by: %s", err)
//		errChan <- err
//		return
//	}
//
//	svcName := "com.github.ClessLi.api.bifrost"
//	//svcName := "bifrostpb.BifrostService"
//	//svcName := "BifrostService"
//	//svcName := "Health"
//
//	instanceId = svcName + "-" + uuid.NewV4().String()
//
//	instanceIP, err := externalIP()
//	if err != nil {
//		utils.Logger.WarningF("Failed to initialize service instance IP. Cased by: %s", err)
//		errChan <- err
//		return
//	}
//	instanceHost := instanceIP.String()
//
//	if !discoveryClient.Register(svcName, instanceId, instanceHost, BifrostConf.ServiceConfig.Port, nil, config.KitLogger) {
//		err = fmt.Errorf("register service %s failed", svcName)
//		utils.Logger.Warning(err.Error())
//		errChan <- err
//		instanceId = ""
//		return
//	}
//}

//func deregisterToRA() {
//	if discoveryClient != nil && !strings.EqualFold(instanceId, "") {
//		if discoveryClient.DeRegister(instanceId, config.KitLogger) {
//			utils.Logger.InfoF("bifrost service (instance ID is '%s') has been unregistered from RA '%s:%d'", instanceId, BifrostConf.RAConfig.Host, BifrostConf.RAConfig.Port)
//		} else {
//			utils.Logger.WarningF("bifrost service (instance ID is '%s') failed to deregister from RA '%s:%d'", instanceId, BifrostConf.RAConfig.Host, BifrostConf.RAConfig.Port)
//		}
//	}
//}

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

func newService(bifrostConf *Config, errChan chan error) (service.Service, service.OffstageManager) {
	webServerConfigServices := make(map[string]service.WebServerConfigService)
	configManagers := make(map[string]configuration.ConfigManager)
	for _, info := range bifrostConf.ServiceConfig.WebServerConfigInfos {
		// 加载配置文件对象
		l := loader.NewLoader()
		ctx, loopPreventer, err := l.LoadFromFilePath(info.ConfPath)
		if err != nil {
			utils.Logger.FatalF("[%s] Load error: %s", info.Name, err)
		}
		c := configuration.NewConfiguration(ctx.(*parser.Config), loopPreventer, new(sync.RWMutex))
		// 初始化日志目录
		logsDir, err := filepath.Abs(filepath.Join(filepath.Dir(filepath.Dir(info.VerifyExecPath)), "logs"))
		if err != nil {
			utils.Logger.FatalF("[%s] Initialization logs directory error: %s", info.Name, err)
		}
		// 添加服务到web服务器配置服务表
		webServerConfigServices[info.Name] = service.NewWebServerConfigService(c, info.VerifyExecPath, logsDir, nginx.NewLog())
		// 添加管理服务到配置服务表
		configManagers[info.Name] = configuration.NewNginxConfigurationManager(l, c, info.VerifyExecPath, info.BackupDir, info.BackupCycle, info.BackupSaveTime, new(sync.RWMutex))
	}
	// 初始化后台服务对象
	var offstage *service.Offstage
	offstage = service.NewOffstage(webServerConfigServices, configManagers, service.NewMetrics(func() []service.WebServerInfo {
		webServerInfos := make([]service.WebServerInfo, 0)
		offstage.Range(func(serverName string, configService service.WebServerConfigService) bool {
			info := service.NewWebServerInfo(serverName)
			info.Version = configService.ServerVersion()
			info.Status = configService.ServerStatus()
			webServerInfos = append(webServerInfos, info)
			return true
		})
		return webServerInfos
	}, errChan))

	// init service
	svc := service.NewService(service.NewViewer(offstage), service.NewUpdater(offstage), service.NewWatcher(offstage))
	// init auth svc
	svc = authentication.AuthenticationMiddleware(bifrostConf.ServiceConfig.AuthServerAddr)(svc)
	// init kit logger
	svc = logging.LoggingMiddleware(config.KitLogger(utils.Stdoutf))(svc)
	return svc, offstage
}

func newGRPCServer(chunkSize int, svc service.Service) *grpc.Server {
	ctx := context.Background()
	// init kit endpoint
	endpoints := endpoint.NewBifrostEndpoints(svc)

	// init kit transport
	transport.ChunkSize = chunkSize
	handlers := transport.NewGRPCHandlers(ctx, endpoints)
	healthCheckHandler := transport.NewHealthCheckHandler(ctx, endpoints)

	// init gRPC server
	gRPCServer := grpc.NewServer(grpc.MaxSendMsgSize(transport.ChunkSize))
	bifrostpb.RegisterViewServiceServer(gRPCServer, handlers)
	bifrostpb.RegisterUpdateServiceServer(gRPCServer, handlers)
	bifrostpb.RegisterWatchServiceServer(gRPCServer, handlers)
	grpc_health_v1.RegisterHealthServer(gRPCServer, healthCheckHandler)
	return gRPCServer
}
