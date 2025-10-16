package bifrost

import (
	"github.com/ClessLi/bifrost/internal/bifrost/config"
	storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"
	storev1nginx "github.com/ClessLi/bifrost/internal/bifrost/store/v1/nginx"
	"github.com/ClessLi/bifrost/internal/pkg/file_watcher"
	"github.com/ClessLi/bifrost/internal/pkg/monitor"
	genericgrpcserver "github.com/ClessLi/bifrost/internal/pkg/server"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"

	logV1 "github.com/ClessLi/component-base/pkg/log/v1"

	"github.com/marmotedu/iam/pkg/shutdown"
	"github.com/marmotedu/iam/pkg/shutdown/shutdownmanagers/posixsignal"
)

type bifrostServer struct {
	gs                      *shutdown.GracefulShutdown
	genericGRPCServer       *genericgrpcserver.GenericGRPCServer
	webSvrConfigsManager    nginx.ConfigsManager
	webSvrMonitor           monitor.Monitor
	webSvrLogWatcherManager *file_watcher.WatcherManager
	webSvrLogDirs           map[string]string
}

type preparedBifrostServer struct {
	*bifrostServer
}

func createBifrostServer(cfg *config.Config) (*bifrostServer, error) {
	logV1.Debug("create bifrost server...")
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	genericConfig, err := buildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	webSvrConfManagerConfig, svrLogDirs, err := buildWebServerConfigsManagerConfig(cfg)
	if err != nil {
		return nil, err
	}

	logWatcherConfig, err := buildLogWatcherConfig(cfg)
	if err != nil {
		return nil, err
	}

	monitorConfig, err := buildMonitorConfig(cfg)
	if err != nil {
		return nil, err
	}

	genericServer, err := genericConfig.Complete().NewGRPCServer()
	if err != nil {
		return nil, err
	}

	configsManagerCC, err := webSvrConfManagerConfig.Complete()
	if err != nil {
		return nil, err
	}
	configsManager, err := configsManagerCC.NewConfigsManager()
	if err != nil {
		return nil, err
	}

	m, err := monitorConfig.Complete().NewMonitor()
	if err != nil {
		return nil, err
	}

	server := &bifrostServer{
		gs:                      gs,
		genericGRPCServer:       genericServer,
		webSvrConfigsManager:    configsManager,
		webSvrLogDirs:           svrLogDirs,
		webSvrMonitor:           m,
		webSvrLogWatcherManager: file_watcher.NewWatcherManager(logWatcherConfig),
	}

	return server, nil
}

func (b *bifrostServer) PrepareRun() preparedBifrostServer {
	logV1.Debug("prepare run...")
	b.initStore()
	initRouter(b.genericGRPCServer)

	b.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		var err error

		b.genericGRPCServer.Close()

		storeIns := storev1.Client()
		if storeIns != nil {
			err = storeIns.Close()
		}

		return err
	}))

	return preparedBifrostServer{b}
}

func (p preparedBifrostServer) Run() error {
	logV1.Debug("preparedBifrostServer run...")
	if err := p.gs.Start(); err != nil {
		logV1.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	logV1.Infof("the web server configs manager is going to run...")
	err := p.webSvrConfigsManager.Start()
	if err != nil {
		logV1.Fatalf("start web server configs manager failed: %s", err.Error())
	}

	logV1.Infof("the web server monitor is going to run...")
	err = p.webSvrMonitor.Start()
	if err != nil {
		logV1.Fatalf("start web server monitor failed: %s", err.Error())
	}

	logV1.Infof("the generic gRPC server is going to run...")

	return p.genericGRPCServer.Run()
}

func (b *bifrostServer) initStore() {
	logV1.Debug("bifrost server init store...")
	storeIns, err := storev1nginx.GetNginxStoreFactory(b.webSvrConfigsManager, b.webSvrLogDirs, b.webSvrMonitor, b.webSvrLogWatcherManager)
	if err != nil {
		logV1.Fatalf("init nginx store failed: %+v", err)
	}
	storev1.SetClient(storeIns)
}

func buildGenericConfig(cfg *config.Config) (genericConfig *genericgrpcserver.Config, lastErr error) {
	logV1.Debug("build generic config, apply all options to generic config")
	genericConfig = genericgrpcserver.NewConfig()

	if lastErr = cfg.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.SecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.InsecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.GRPCServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.RAOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	return
}

func buildWebServerConfigsManagerConfig(cfg *config.Config) (*nginx.Config, map[string]string, error) {
	webSvrConfigsManagerConfig := &nginx.Config{ManagersConfig: make(map[string]*configuration.ManagerConfig)}
	// apply options to configuration for ConfigsManager
	cfg.WebServerConfigsOptions.ApplyGenericOptsTo(webSvrConfigsManagerConfig)

	svrLogDirs := make(map[string]string)
	for _, itemOpts := range cfg.WebServerConfigsOptions.WebServerConfigs {
		if itemOpts.ServerType == "nginx" {
			itemOpts.ApplyToNginx(webSvrConfigsManagerConfig)
		}
		svrLogDirs[itemOpts.ServerName] = itemOpts.LogsDirPath
	}

	return webSvrConfigsManagerConfig, svrLogDirs, nil
}

func buildLogWatcherConfig(cfg *config.Config) (*file_watcher.Config, error) {
	watcherConfig := file_watcher.NewConfig()

	return watcherConfig, cfg.WebServerLogWatcherOptions.ApplyTo(watcherConfig)
}

func buildMonitorConfig(cfg *config.Config) (*monitor.Config, error) {
	return &monitor.Config{
		MonitoringSyncInterval:      cfg.MonitorOptions.SyncInterval,
		MonitoringCycle:             cfg.MonitorOptions.CycleTime,
		MonitoringFrequencyPerCycle: cfg.MonitorOptions.FrequencyPerCycle,
	}, nil
}
