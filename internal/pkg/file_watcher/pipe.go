package file_watcher

import (
	"context"
	"sync"
	"time"
)

type Pipe interface {
	Close()
	IsClosed() bool
	Input() chan<- []byte
	Output() <-chan []byte
}

type pipe struct {
	mu       sync.RWMutex
	cancel   context.CancelFunc
	isClosed bool
	c        chan []byte
}

func (p *pipe) Close() {
	p.cancel()
}

func (p *pipe) IsClosed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.isClosed
}

func (p *pipe) Input() chan<- []byte {
	return p.c
}

func (p *pipe) Output() <-chan []byte {
	return p.c
}

func newPipe(ctx context.Context, cancel context.CancelFunc) Pipe {
	channel := make(chan []byte)
	p := &pipe{
		mu:       sync.RWMutex{},
		cancel:   cancel,
		isClosed: false,
		c:        channel,
	}
	go func() {
		<-ctx.Done()
		p.mu.Lock()
		defer p.mu.Unlock()
		p.isClosed = true
		close(channel)
	}()

	return p
}

func NewPipe(ctx context.Context) Pipe {
	cctx, cancel := context.WithCancel(ctx)

	return newPipe(cctx, cancel)
}

func TimeoutPipe(ctx context.Context, timeout time.Duration) Pipe {
	cctx, cancel := context.WithTimeout(ctx, timeout)

	return newPipe(cctx, cancel)
}
