package service

import "github.com/ClessLi/bifrost/internal/pkg/bifrost/service/web_server_manager"

type Watcher interface {
	Watch(WatchRequester) WatchResponder
}

type watcher struct {
	offstage offstageWatcher
}

func NewWatcher(offstage offstageWatcher) Watcher {
	return &watcher{offstage: offstage}
}

func (w watcher) Watch(req WatchRequester) WatchResponder {
	serverName := req.GetServerName()
	objectName := req.GetObjectName()
	var dataChan <-chan []byte
	var transferErrChan <-chan error
	closeFunc := func() error { return nil }
	var err error
	switch req.GetRequestType() {
	case WatchLog:
		var logWatcher web_server_manager.LogWatcher
		logWatcher, err = w.offstage.WatchLog(serverName, objectName)
		if logWatcher != nil {
			dataChan = logWatcher.GetDataChan()
			transferErrChan = logWatcher.GetTransferErrorChan()
			closeFunc = func() error {
				return logWatcher.Close()
			}
		}
	default:
		err = UnknownRequestType
	}
	return NewWatchResponder(serverName, closeFunc, dataChan, transferErrChan, err)
}
