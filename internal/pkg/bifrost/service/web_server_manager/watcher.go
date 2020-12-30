package web_server_manager

import (
	"errors"
	"time"
)

type Watcher interface {
	GetDataChan() <-chan []byte
	GetErrChan() <-chan error
	Close() error
	getSignalChan() chan int
	inputDataChan() chan<- []byte
	inputErrChan() chan<- error
}

type logWatcher struct {
	DataC    chan []byte
	ErrC     chan error
	SignalC  chan int
	isClosed bool
}

func (l logWatcher) GetDataChan() <-chan []byte {
	return l.DataC
}

func (l logWatcher) GetErrChan() <-chan error {
	return l.ErrC
}

func (l logWatcher) Close() error {
	if l.isClosed {
		return errors.New("LogWatcher was closed")
	}
	select {
	case l.SignalC <- 9:
		l.isClosed = true
		return nil
	case <-time.After(time.Second * 30):
		return errors.New("time out to close LogWatcher")
	}
}

func (l logWatcher) getSignalChan() chan int {
	return l.SignalC
}

func (l logWatcher) inputDataChan() chan<- []byte {
	return l.DataC
}

func (l logWatcher) inputErrChan() chan<- error {
	return l.ErrC
}

func (l logWatcher) getSignal() int {
	return <-l.SignalC
}

func NewLogWatcher() Watcher {
	return &logWatcher{
		DataC:    make(chan []byte, 0),
		ErrC:     make(chan error, 0),
		SignalC:  make(chan int, 0),
		isClosed: false,
	}
}
