package web_server_manager

type LogWatcher interface {
	GetDataChan() <-chan []byte
	GetTransferErrorChan() <-chan error
	Close() error
}

type logWatcher struct {
	dataChan          chan []byte
	transferErrorChan chan error
	closeFunc         func() error
}

func NewLogWatcher(dataChan chan []byte, errChan chan error, closeFunc func() error) LogWatcher {
	return &logWatcher{
		dataChan:          dataChan,
		transferErrorChan: errChan,
		closeFunc:         closeFunc,
	}
}

func (l logWatcher) GetDataChan() <-chan []byte {
	return l.dataChan
}

func (l logWatcher) GetTransferErrorChan() <-chan error {
	return l.transferErrorChan
}

func (l logWatcher) Close() error {
	return l.closeFunc()
}
