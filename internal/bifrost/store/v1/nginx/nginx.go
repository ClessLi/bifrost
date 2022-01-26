package nginx

import (
	storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"
	"github.com/ClessLi/bifrost/internal/pkg/file_watcher"
	"github.com/ClessLi/bifrost/internal/pkg/monitor"
	genericoptions "github.com/ClessLi/bifrost/internal/pkg/options"
	log "github.com/ClessLi/bifrost/pkg/log/v1"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx"
	"github.com/marmotedu/errors"
	"sync"
)

const (
	nginxServer = "nginx"
)

type webServerStore struct {
	cms      nginx.ConfigsManager
	m        monitor.Monitor
	wm       *file_watcher.WatcherManager
	logsDirs map[string]string
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
		w.cms.Stop(),
		w.m.Stop(),
		w.wm.StopAll(),
	})
}

var _ storev1.StoreFactory = &webServerStore{}

var (
	nginxStoreFactory storev1.StoreFactory
	once              sync.Once
)

func GetNginxStoreFactory(webSvrConfOpts *genericoptions.WebServerConfigsOptions, monitorOpts *genericoptions.MonitorOptions, webSvrLogWatcherOpts *genericoptions.WebServerLogWatcherOptions) (storev1.StoreFactory, error) {
	if webSvrConfOpts == nil && nginxStoreFactory == nil {
		return nil, errors.New("failed to get nginx store factory")
	}

	var err error
	var cms nginx.ConfigsManager
	var m monitor.Monitor
	once.Do(func() {
		// init and start config managers and log watcher manager
		cmsOpts := nginx.ConfigsManagerOptions{Options: make([]nginx.ConfigManagerOptions, 0)}
		svrLogsDirs := make(map[string]string)
		for _, itemOpts := range webSvrConfOpts.WebServerConfigs {
			if itemOpts.ServerType == nginxServer {
				cmsOpts.Options = append(cmsOpts.Options, nginx.ConfigManagerOptions{
					ServerName:     itemOpts.ServerName,
					MainConfigPath: itemOpts.ConfigPath,
					ServerBinPath:  itemOpts.VerifyExecPath,
					BackupDir:      itemOpts.BackupDir,
					BackupCycle:    itemOpts.BackupCycle,
					BackupSaveTime: itemOpts.BackupSaveTime,
				})
			}
			svrLogsDirs[itemOpts.ServerName] = itemOpts.LogsDirPath
		}
		cms, err = nginx.New(cmsOpts)
		if err != nil {
			return
		}

		err = cms.Start()
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
			err := m.Start()
			if err != nil {
				log.Fatal(err.Error())

				return
			}
		}()

		// build nginx store factory
		nginxStoreFactory = &webServerStore{
			cms:      cms,
			m:        m,
			wm:       wm,
			logsDirs: svrLogsDirs,
		}
	})

	if nginxStoreFactory == nil || err != nil {
		return nil, errors.Errorf("failed to get nginx store factory, nginx store factory: %+v, err: %w", nginxStoreFactory, err)
	}

	return nginxStoreFactory, nil
}
