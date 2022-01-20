package nginx

import (
	"context"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/internal/pkg/file_watcher"
	"github.com/marmotedu/errors"
	"github.com/marmotedu/iam/pkg/log"
	"path/filepath"
	"regexp"
	"time"
)

type webServerLogWatcherStore struct {
	watcherManager    *file_watcher.WatcherManager
	webServerLogsDirs map[string]string
}

func (w *webServerLogWatcherStore) Watch(ctx context.Context, request *v1.WebServerLogWatchRequest) (*v1.WebServerLog, error) {
	var logPath string
	if logDir, ok := w.webServerLogsDirs[request.ServerName.Name]; ok {
		logPath = filepath.Join(logDir, request.LogName)
	} else {
		return nil, errors.WithCode(code.ErrConfigurationNotFound, "web server %s is not exist", request.ServerName.Name)
	}
	outputC, err := w.watcherManager.Watch(logPath)
	if err != nil {
		return nil, err
	}
	if len(request.FilteringRegexpRule) > 0 {
		fOutputC := make(chan []byte)
		go func() {
			defer close(fOutputC)
			needClose := false
			for {
				select {
				case fOutputC <- filterOutput(outputC, request.FilteringRegexpRule, &needClose):
					if needClose {
						return
					}
				case <-ctx.Done():
					needClose = true
					return
				}
			}
		}()
		return &v1.WebServerLog{Lines: fOutputC}, nil
	}
	return &v1.WebServerLog{Lines: outputC}, nil
}

func filterOutput(output <-chan []byte, pattern string, needClose *bool) []byte {
	for !*needClose {
		data := <-output
		if data == nil {
			*needClose = true
			return nil
		}
		match, err := regexp.Match(pattern, data)
		if err != nil {
			*needClose = true
			log.Warnf("web server log watch error. %s", err.Error())
			return []byte(err.Error())
		}
		if match {
			return data
		}

		time.Sleep(time.Second)
	}
	return nil
}

func newWebServerLogWatcherStore(store *webServerStore) *webServerLogWatcherStore {
	return &webServerLogWatcherStore{
		watcherManager:    store.wm,
		webServerLogsDirs: store.logsDirs,
	}
}
