package bifrost

import (
	"github.com/ClessLi/bifrost/internal/bifrost/config"
	storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"
	storev1nginx "github.com/ClessLi/bifrost/internal/bifrost/store/v1/nginx"
	genericoptions "github.com/ClessLi/bifrost/internal/pkg/options"
	genericgrpcserver "github.com/ClessLi/bifrost/internal/pkg/server"

	logV1 "github.com/ClessLi/component-base/pkg/log/v1"

	"github.com/marmotedu/iam/pkg/shutdown"
	"github.com/marmotedu/iam/pkg/shutdown/shutdownmanagers/posixsignal"
)

type bifrostServer struct {
	gs                   *shutdown.GracefulShutdown
	genericGRPCServer    *genericgrpcserver.GenericGRPCServer
	webSvrConfigsOpts    *genericoptions.WebServerConfigsOptions
	monitorOpts          *genericoptions.MonitorOptions
	webSvrLogWatcherOpts *genericoptions.WebServerLogWatcherOptions
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

	genericServer, err := genericConfig.Complete().NewGRPCServer()
	if err != nil {
		return nil, err
	}

	server := &bifrostServer{
		gs:                   gs,
		genericGRPCServer:    genericServer,
		webSvrConfigsOpts:    cfg.WebServerConfigsOptions,
		monitorOpts:          cfg.MonitorOptions,
		webSvrLogWatcherOpts: cfg.WebServerLogWatcherOptions,
	}

	return server, nil
}

func (b *bifrostServer) PrepareRun() preparedBifrostServer {
	logV1.Debug("prepare run...")
	b.initStore()
	initRouter(b.genericGRPCServer)
	b.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		var err error
		storeIns := storev1.Client()
		if storeIns != nil {
			err = storeIns.Close()
		}
		b.genericGRPCServer.Close()

		return err
	}))

	return preparedBifrostServer{b}
}

func (p preparedBifrostServer) Run() error {
	logV1.Debug("prepareBifrostServer run...")
	if err := p.gs.Start(); err != nil {
		logV1.Fatalf("start shutdown manager failed: %s", err.Error())
	}
	logV1.Infof("the generic gRPC server is going to run...")

	return p.genericGRPCServer.Run()
}

func (b *bifrostServer) initStore() {
	logV1.Debug("bifrost server init store...")
	storeIns, err := storev1nginx.GetNginxStoreFactory(b.webSvrConfigsOpts, b.monitorOpts, b.webSvrLogWatcherOpts)
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
