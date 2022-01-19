package file_watcher

import (
	"github.com/marmotedu/errors"
	"sync"
)

type WatcherManager struct {
	config   *Config
	watchers map[string]*FileWatcher
	mu       sync.RWMutex
}

func (wm *WatcherManager) Watch(file string, outputC chan []byte) error {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	cconf, err := wm.config.Complete(file)
	if err != nil {
		return err
	}
	if watcher, has := wm.watchers[cconf.filePath]; has && !watcher.isClosed {
		return watcher.AddOutput(outputC)
	}
	watcher, err := cconf.NewFileWatcher(outputC)
	if err != nil {
		return err
	}
	wm.watchers[watcher.filePath] = watcher
	return nil
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
		watchers: make(map[string]*FileWatcher),
		mu:       sync.RWMutex{},
	}
}
