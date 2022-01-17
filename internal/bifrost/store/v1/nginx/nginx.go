package nginx

import (
	storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"
	"github.com/ClessLi/bifrost/internal/pkg/monitor"
	genericoptions "github.com/ClessLi/bifrost/internal/pkg/options"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx"
	"github.com/marmotedu/errors"
	"sync"
)

const (
	nginxServer = "nginx"
)

type webServerStore struct {
	cms nginx.ConfigsManager
	m   monitor.Monitor
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

func (w *webServerStore) Close() error {
	return w.cms.Stop()
}

var _ storev1.StoreFactory = &webServerStore{}

var (
	nginxStoreFactory storev1.StoreFactory
	once              sync.Once
)

func GetNginxStoreFactory(webSvrConfOpts *genericoptions.WebServerConfigsOptions, monitorOpts *genericoptions.MonitorOptions) (storev1.StoreFactory, error) {
	if webSvrConfOpts == nil && nginxStoreFactory == nil {
		return nil, errors.New("failed to get nginx store factory")
	}

	var err error
	var cms nginx.ConfigsManager
	var m monitor.Monitor
	once.Do(func() {
		cmsOpts := nginx.ConfigsManagerOptions{Options: make([]nginx.ConfigManagerOptions, 0)}
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
		}
		cms, err = nginx.New(cmsOpts)
		if err != nil {
			return
		}

		err = cms.Start()
		if err != nil {
			return
		}

		mconf := &monitor.Config{
			MonitoringSyncInterval:      monitorOpts.SyncInterval,
			MonitoringCycle:             monitorOpts.CycleTime,
			MonitoringFrequencyPerCycle: monitorOpts.FrequencyPerCycle,
		}
		m, err = mconf.Complete().NewMonitor()
		if err != nil {
			return
		}

		nginxStoreFactory = &webServerStore{
			cms: cms,
			m:   m,
		}
	})

	if nginxStoreFactory == nil || err != nil {
		return nil, errors.Errorf("failed to get nginx store factory, nginx store factory: %+v, err: %w", nginxStoreFactory, err)
	}

	return nginxStoreFactory, nil
}
