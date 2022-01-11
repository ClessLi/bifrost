package service

type Watcher interface {
	Watch(WatchRequestInfo) WatchResponseInfo
}

type watcher struct {
	offstage offstageWatcher
}

func NewWatcher(offstage offstageWatcher) Watcher {
	if offstage == nil {
		panic("offstage is nil")
	}

	return &watcher{offstage: offstage}
}

func (w watcher) Watch(req WatchRequestInfo) WatchResponseInfo {
	var dataChan <-chan []byte
	var transferErrChan <-chan error
	closeFunc := func() error {
		return ErrInvalidResponseInfo
	}
	var err error

	if req == nil {
		err = ErrNilRequestInfo
		return NewWatchResponseInfo("", closeFunc, dataChan, transferErrChan, err)
	}

	serverName := req.GetServerName()
	objectName := req.GetWatchedObjectName()
	switch req.GetRequestType() {
	case WatchLog:
		logWatcherIns, logWatcherErr := w.offstage.WatchLog(serverName, objectName)
		err = logWatcherErr
		if logWatcherIns != nil {
			dataChan = logWatcherIns.GetDataChan()
			transferErrChan = logWatcherIns.GetTransferErrorChan()
			//closeFunc = func() error {
			//	return logWatcherIns.Close()
			//}
			closeFunc = logWatcherIns.Close
		}
	default:
		err = ErrUnknownRequestType
	}

	return NewWatchResponseInfo(serverName, closeFunc, dataChan, transferErrChan, err)
}
