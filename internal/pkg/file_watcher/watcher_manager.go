package file_watcher

import (
	"context"
	"sync"

	"github.com/ClessLi/bifrost/internal/pkg/code"

	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
)

type WatcherManager struct {
	config *Config
	mu     sync.RWMutex

	watchers map[string]*FileWatcher
}

func (wm *WatcherManager) Watch(ctx context.Context, file string) (<-chan []byte, error) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	cconf, err := wm.config.Complete(file)
	if err != nil {
		return nil, err
	}
	if watcher, has := wm.watchers[cconf.filePath]; has && !watcher.IsClosed() {
		return watcher.Output(ctx)
	}
	watcher, output, err := cconf.NewFileWatcher(ctx)
	if err != nil {
		return nil, err
	}
	wm.watchers[watcher.filePath] = watcher

	return output, nil
}

func (wm *WatcherManager) StopAll() error {
	var errs []error
	wm.mu.Lock()
	defer wm.mu.Unlock()
	for filePath, watcher := range wm.watchers {
		logV1.Infof("stopping watching file '%s'", filePath)
		err := watcher.Stop()
		if err != nil {
			errs = append(errs, code.Errorf("failed to stop watching file '%s'. %s", filePath, err.Error()))
		} else {
			delete(wm.watchers, filePath)
		}
	}

	return code.NewAggregate(errs)
}

func NewWatcherManager(config *Config) *WatcherManager {
	return &WatcherManager{
		config:   config,
		mu:       sync.RWMutex{},
		watchers: make(map[string]*FileWatcher),
	}
}
