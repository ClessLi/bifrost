package endpoint

type ClientWatcher interface {
	GetDataChan() <-chan []byte
	GetErrChan() <-chan error
	Close() error
}

type clientLogWatcher struct {
	dataChan  <-chan []byte
	errChan   <-chan error
	closeFunc func() error
}

func (w clientLogWatcher) GetDataChan() <-chan []byte {
	return w.dataChan
}

func (w clientLogWatcher) GetErrChan() <-chan error {
	return w.errChan
}

func (w clientLogWatcher) Close() error {
	return w.closeFunc()
}

func newClientLogWatcher(dataChan <-chan []byte, errChan <-chan error, closeFunc func() error) ClientWatcher {
	return &clientLogWatcher{
		dataChan:  dataChan,
		errChan:   errChan,
		closeFunc: closeFunc,
	}
}
