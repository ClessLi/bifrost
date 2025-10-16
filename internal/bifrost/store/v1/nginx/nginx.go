package nginx

import (
	"sync"
	"time"

	storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"
	"github.com/ClessLi/bifrost/internal/pkg/file_watcher"
	"github.com/ClessLi/bifrost/internal/pkg/monitor"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx"

	"github.com/marmotedu/errors"
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

func (w *webServerStore) WebServerLogWatcher() storev1.WebServerLogWatcherStore {
	return newWebServerLogWatcherStore(w)
}

func (w *webServerStore) WebServerBinCMD() storev1.WebServerBinCMDStore {
	return newNginxBinCMDStore(w)
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
	webSvrConfManager nginx.ConfigsManager,
	svrLogsDirs map[string]string,
	m monitor.Monitor,
	logWatcherManager *file_watcher.WatcherManager,
) (storev1.StoreFactory, error) {
	if webSvrConfManager == nil && nginxStoreFactory == nil {
		return nil, errors.New("failed to get web server store factory")
	}

	once.Do(func() {
		// build nginx store factory
		nginxStoreFactory = &webServerStore{
			configsManger:  webSvrConfManager,
			monitor:        m,
			watcherManager: logWatcherManager,
			logsDirs:       svrLogsDirs,
		}
	})

	if nginxStoreFactory == nil {
		return nil, errors.New("failed to setup nginx store factory")
	}

	return nginxStoreFactory, nil
}
