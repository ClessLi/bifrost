package service

import (
	"errors"
	"time"
)

type LogWatcher struct {
	DataC    chan []byte
	ErrC     chan error
	SignalC  chan int
	isClosed bool
}

func (l LogWatcher) Close() error {
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

func NewLogWatcher() *LogWatcher {
	return &LogWatcher{
		DataC:    make(chan []byte, 1),
		ErrC:     make(chan error, 1),
		SignalC:  make(chan int, 1),
		isClosed: false,
	}
}
