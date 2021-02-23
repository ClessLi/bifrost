package bifrost

type WatchClient interface {
	GetDataChan() <-chan []byte
	GetErrChan() <-chan error
	Close() error
}

type logWatchClient struct {
	dataChan  <-chan []byte
	errChan   <-chan error
	closeFunc func() error
}

func (w logWatchClient) GetDataChan() <-chan []byte {
	return w.dataChan
}

func (w logWatchClient) GetErrChan() <-chan error {
	return w.errChan
}

func (w logWatchClient) Close() error {
	return w.closeFunc()
}

func newLogWatcherClient(dataChan <-chan []byte, errChan <-chan error, closeFunc func() error) WatchClient {
	return &logWatchClient{
		dataChan:  dataChan,
		errChan:   errChan,
		closeFunc: closeFunc,
	}
}
