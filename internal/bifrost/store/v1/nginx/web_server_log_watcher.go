package nginx

import (
	"context"
	"path/filepath"
	"regexp"
	"time"

	"github.com/marmotedu/errors"
	"github.com/marmotedu/iam/pkg/log"

	v1 "github.com/yongPhone/bifrost/api/bifrost/v1"
	"github.com/yongPhone/bifrost/internal/pkg/code"
	"github.com/yongPhone/bifrost/internal/pkg/file_watcher"
)

type webServerLogWatcherStore struct {
	watcherManager    *file_watcher.WatcherManager
	webServerLogsDirs map[string]string
}

//nolint:gocognit
func (w *webServerLogWatcherStore) Watch(
	ctx context.Context,
	request *v1.WebServerLogWatchRequest,
) (*v1.WebServerLog, error) {
	var logPath string
	if logDir, ok := w.webServerLogsDirs[request.ServerName.Name]; ok {
		logPath = filepath.Join(logDir, request.LogName)
	} else {
		return nil, errors.WithCode(code.ErrConfigurationNotFound, "web server %s is not exist", request.ServerName.Name)
	}
	outputC, err := w.watcherManager.Watch(ctx, logPath)
	if err != nil {
		return nil, err
	}
	if len(request.FilteringRegexpRule) > 0 {
		fOutputC := make(chan []byte)
		go func() {
			defer close(fOutputC)
			for {
				select {
				// case fOutputC <- filterOutput(outputC, request.FilteringRegexpRule, &needClose):
				//	if needClose {
				//		return
				//	}
				case data := <-outputC:
					ok, err := filterOutput(data, request.FilteringRegexpRule)
					if err != nil {
						data = []byte(err.Error())
					}
					if ok {
						select {
						case fOutputC <- data:
						case <-time.After(time.Second * 30):
							log.Warnf("send filtered data timeout(30s)")

							return
						}
					}
					if err != nil {
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()

		return &v1.WebServerLog{Lines: fOutputC}, nil
	}

	return &v1.WebServerLog{Lines: outputC}, nil
}

func filterOutput(data []byte, pattern string) (bool, error) {
	if data == nil {
		return true, nil
	}
	match, err := regexp.Match(pattern, data)
	if err != nil {
		log.Warnf("web server log watch error. %s", err.Error())

		return true, err
	}
	if match {
		return true, nil
	}

	return false, nil
}

func newWebServerLogWatcherStore(store *webServerStore) *webServerLogWatcherStore {
	return &webServerLogWatcherStore{
		watcherManager:    store.wm,
		webServerLogsDirs: store.logsDirs,
	}
}
