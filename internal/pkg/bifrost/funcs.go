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
	"github.com/ClessLi/bifrost/pkg/client/auth"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration/parser"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/loader"
	"github.com/ClessLi/bifrost/pkg/server_log/nginx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"gopkg.in/yaml.v2"
	"net"
	"os"
	"path/filepath"
	"sync"
)

var (
	onceForAuthSvcClient       sync.Once
	singletonAuthServiceClient *auth.Client
	onceForBifrostConf         sync.Once
	singletonBifrostConf       *Config
)

func getBifrostConfInstance() *Config {
	onceForBifrostConf.Do(func() {
		if singletonBifrostConf != nil {
			return
		}
		// 初始化bifrost配置
		confData, err := utils.ReadFile(*confPath)
		if err != nil {
			panic(err)
		}
		singletonBifrostConf = new(Config)
		// 加载bifrost配置
		err = yaml.Unmarshal(confData, singletonBifrostConf)
		if err != nil {
			panic(err)
		}

		// 配置必填项检查
		err = singletonBifrostConf.check()
		if err != nil {
			panic(err)
		}

		// 初始化日志
		logDir, err := filepath.Abs(singletonBifrostConf.LogDir)
		if err != nil {
			panic(err)
		}

		logPath := filepath.Join(logDir, "bifrost.log")
		utils.Logf, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}

		utils.InitLogger(utils.Logf, singletonBifrostConf.LogConfig.Level)

		// 初始化应用运行日志输出
		stdoutPath := filepath.Join(logDir, "bifrost.out")
		utils.Stdoutf, err = os.OpenFile(stdoutPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}

		os.Stdout = utils.Stdoutf
		os.Stderr = utils.Stdoutf
	})
	return singletonBifrostConf
}

func getAuthServiceClientInstance() *auth.Client {
	onceForAuthSvcClient.Do(func() {
		if singletonAuthServiceClient != nil {
			return
		}
		var err error
		singletonAuthServiceClient, err = auth.NewClientFromGRPCServerAddress(getBifrostConfInstance().ServiceConfig.AuthServerAddr)
		if err != nil {
			panic(err)
		}
	})
	return singletonAuthServiceClient
}

func getProc(path string) (*os.Process, error) {
	pid, pidErr := utils.GetPid(path)
	if pidErr != nil {
		return nil, pidErr
	}
	return os.FindProcess(pid)
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

func newService(errChan chan error) (service.Service, service.OffstageManager) {
	webServerConfigServices := make(map[string]service.WebServerConfigService)
	configManagers := make(map[string]configuration.ConfigManager)
	for _, info := range getBifrostConfInstance().ServiceConfig.WebServerConfigInfos {
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
	svc = authentication.AuthenticationMiddleware(getAuthServiceClientInstance())(svc)
	// init kit logger
	svc = logging.LoggingMiddleware(config.KitLogger(utils.Stdoutf))(svc)
	return svc, offstage
}

func newGRPCServer(chunkSize int, svc service.Service) *grpc.Server {
	//ctx := context.Background()
	// init kit endpoint
	endpoints := endpoint.NewBifrostEndpoints(svc)

	// init kit transport
	transport.ChunkSize = chunkSize
	handlers := transport.NewGRPCHandlers(endpoints)
	healthCheckHandler := transport.NewHealthCheckHandler(endpoints)

	// init gRPC server
	gRPCServer := grpc.NewServer(grpc.MaxSendMsgSize(transport.ChunkSize))
	bifrostpb.RegisterViewServiceServer(gRPCServer, handlers)
	bifrostpb.RegisterUpdateServiceServer(gRPCServer, handlers)
	bifrostpb.RegisterWatchServiceServer(gRPCServer, handlers)
	grpc_health_v1.RegisterHealthServer(gRPCServer, healthCheckHandler)
	return gRPCServer
}
