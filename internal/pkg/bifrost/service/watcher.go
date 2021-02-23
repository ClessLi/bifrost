package service

type Watcher interface {
	Watch(info WatchRequestInfo) WatchResponseInfo
}

type watcher struct {
	offstage offstageWatcher
}

func NewWatcher(offstage offstageWatcher) Watcher {
	return &watcher{offstage: offstage}
}

func (w watcher) Watch(req WatchRequestInfo) WatchResponseInfo {
	serverName := req.GetServerName()
	objectName := req.GetObjectName()
	var dataChan <-chan []byte
	var transferErrChan <-chan error
	closeFunc := func() error { return nil }
	var err error
	switch req.GetRequestType() {
	case WatchLog:
		var logWatcher LogWatcher
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
	return NewWatchResponseInfo(serverName, closeFunc, dataChan, transferErrChan, err)
}
