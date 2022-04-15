package file_watcher

import (
	"context"
	"sync"

	"github.com/marmotedu/errors"
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
		err := watcher.Stop()
		if err != nil {
			errs = append(errs, errors.Errorf("failed to stop watching file '%s'. %s", filePath, err.Error()))
		} else {
			delete(wm.watchers, filePath)
		}
	}

	return errors.NewAggregate(errs)
}

func NewWatcherManager(config *Config) *WatcherManager {
	return &WatcherManager{
		config:   config,
		mu:       sync.RWMutex{},
		watchers: make(map[string]*FileWatcher),
	}
}
