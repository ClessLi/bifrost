package nginx

import (
	storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"
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
}

func (w *webServerStore) Close() error {
	return w.cms.Stop()
}

func (w *webServerStore) WebServerConfig() storev1.WebServerConfigStore {
	return newNginxConfigStore(w)
}

var _ storev1.StoreFactory = &webServerStore{}

var (
	nginxStoreFactory storev1.StoreFactory
	once              sync.Once
)

func GetNginxStoreFactory(opts []*genericoptions.WebServerConfigOptions) (storev1.StoreFactory, error) {
	if opts == nil && nginxStoreFactory == nil {
		return nil, errors.New("failed to get nginx store factory")
	}

	var err error
	var cms nginx.ConfigsManager
	once.Do(func() {
		options := nginx.ConfigsManagerOptions{Options: make([]nginx.ConfigManagerOptions, 0)}
		for _, itemOpts := range opts {
			//if len(strings.TrimSpace(itemOpts.ServerName)) == 0 {
			//	err := errors.Errorf("the server name of the %dth options in WebServerConfigsOptions is empty", i)
			//	panic(err)
			//}
			if itemOpts.ServerType == nginxServer {
				options.Options = append(options.Options, nginx.ConfigManagerOptions{
					MainConfigPath: itemOpts.ConfigPath,
					ServerBinPath:  itemOpts.VerifyExecPath,
					BackupDir:      itemOpts.BackupDir,
					BackupCycle:    itemOpts.BackupCycle,
					BackupSaveTime: itemOpts.BackupSaveTime,
				})
			}
		}
		cms, err = nginx.New(options)
		if err != nil {
			return
		}

		err = cms.Start()
		if err != nil {
			return
		}

		nginxStoreFactory = &webServerStore{cms: cms}
	})

	if nginxStoreFactory == nil || err != nil {
		return nil, errors.Errorf("failed to get nginx store factory, nginx store factory: %+v, err: %w", nginxStoreFactory, err)
	}

	return nginxStoreFactory, nil
}
