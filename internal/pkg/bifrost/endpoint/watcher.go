package endpoint

type Watcher struct {
	dataChan  <-chan []byte
	errChan   <-chan error
	closeFunc func() error
}

func (w Watcher) GetDataChan() <-chan []byte {
	return w.dataChan
}

func (w Watcher) GetErrChan() <-chan error {
	return w.errChan
}

func (w Watcher) Close() error {
	return w.closeFunc()
}

func NewWatcher(dataChan <-chan []byte, errChan <-chan error, closeFunc func() error) *Watcher {
	return &Watcher{
		dataChan:  dataChan,
		errChan:   errChan,
		closeFunc: closeFunc,
	}
}
