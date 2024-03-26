package nginx

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"sync"
	"time"

	"github.com/marmotedu/errors"

	storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"
	"github.com/ClessLi/bifrost/internal/pkg/file_watcher"
	"github.com/ClessLi/bifrost/internal/pkg/monitor"
	genericoptions "github.com/ClessLi/bifrost/internal/pkg/options"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx"
)

const (
	nginxServer = "nginx"
)

type webServerStore struct {
	configsManger  nginx.ConfigsManager
	monitor        monitor.Monitor
	watcherManager *file_watcher.WatcherManager
	logsDirs       map[string]string
}

func (w *webServerStore) WebServerStatus() storev1.WebServerStatusStore {
	return newWebServerStatusStore(w)
}

func (w *webServerStore) WebServerConfig() storev1.WebServerConfigStore {
	return newNginxConfigStore(w)
}

func (w *webServerStore) WebServerStatistics() storev1.WebServerStatisticsStore {
	return newNginxStatisticsStore(w)
}

func (w *webServerStore) WebServerLogWatcher() storev1.WebServerLogWatcher {
	return newWebServerLogWatcherStore(w)
}

func (w *webServerStore) Close() error {
	return errors.NewAggregate([]error{
		w.configsManger.Stop(5 * time.Minute),
		w.monitor.Stop(),
		w.watcherManager.StopAll(),
	})
}

var _ storev1.StoreFactory = &webServerStore{}

var (
	nginxStoreFactory storev1.StoreFactory
	once              sync.Once
)

func GetNginxStoreFactory(
	webSvrConfOpts *genericoptions.WebServerConfigsOptions,
	monitorOpts *genericoptions.MonitorOptions,
	webSvrLogWatcherOpts *genericoptions.WebServerLogWatcherOptions,
) (storev1.StoreFactory, error) {
	if webSvrConfOpts == nil && nginxStoreFactory == nil {
		return nil, errors.New("failed to get nginx store factory")
	}

	var err error
	var configsManager nginx.ConfigsManager
	var m monitor.Monitor
	once.Do(func() {
		// init and start config managers and log watcher manager
		cmConf := &nginx.Config{ManagersConfig: make(map[string]*configuration.ManagerConfig)}
		svrLogsDirs := make(map[string]string)
		for _, itemOpts := range webSvrConfOpts.WebServerConfigs {
			if itemOpts.ServerType == nginxServer {
				itemOpts.ApplyToNginx(cmConf)
			}
			svrLogsDirs[itemOpts.ServerName] = itemOpts.LogsDirPath
		}
		cmCompletedConf, err := cmConf.Complete()
		if err != nil {
			return
		}
		configsManager, err = cmCompletedConf.NewConfigsManager()
		if err != nil {
			return
		}

		err = configsManager.Start()
		if err != nil {
			return
		}

		wmconf := file_watcher.NewConfig()
		err = webSvrLogWatcherOpts.ApplyTo(wmconf)
		if err != nil {
			return
		}
		wm := file_watcher.NewWatcherManager(wmconf)

		// init and start monitor
		mconf := &monitor.Config{
			MonitoringSyncInterval:      monitorOpts.SyncInterval,
			MonitoringCycle:             monitorOpts.CycleTime,
			MonitoringFrequencyPerCycle: monitorOpts.FrequencyPerCycle,
		}
		m, err = mconf.Complete().NewMonitor()
		if err != nil {
			return
		}

		go func() {
			if err := m.Start(); err != nil { //nolint:govet
				logV1.Fatal(err.Error())

				return
			}
		}()

		// build nginx store factory
		nginxStoreFactory = &webServerStore{
			configsManger:  configsManager,
			monitor:        m,
			watcherManager: wm,
			logsDirs:       svrLogsDirs,
		}
	})

	if nginxStoreFactory == nil || err != nil {
		return nil, errors.Errorf( //nolint:govet
			"failed to get nginx store factory, nginx store factory: %+v, err: %w",
			nginxStoreFactory,
			err,
		)
	}

	return nginxStoreFactory, nil
}
