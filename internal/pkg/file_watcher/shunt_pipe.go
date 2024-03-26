package file_watcher

import (
	"context"
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"sync"
	"time"

	"github.com/marmotedu/errors"
)

type ShuntPipe interface {
	Output(ctx context.Context) (<-chan []byte, error)
	Input() chan<- []byte
	IsClosed() bool
	Start() error
	Close() error
}

type shuntPipe struct {
	startLocker sync.Locker
	mu          sync.RWMutex

	ctx      context.Context
	cancel   context.CancelFunc
	isClosed bool

	maxConnections int
	outputTimeout  time.Duration

	input   Pipe
	outputs []Pipe
}

func (s *shuntPipe) Output(ctx context.Context) (<-chan []byte, error) {
	err := s.checkState()
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ctx.Err() != nil {
		return nil, errors.Wrap(s.ctx.Err(), "shunt pipe is closed. %s")
	}
	if s.maxConnections <= len(s.outputs) {
		return nil, errors.Errorf(
			"the number of connections has reached the maximum (%d/%d)",
			len(s.outputs),
			s.maxConnections,
		)
	}

	outputPipe := TimeoutPipe(ctx, s.outputTimeout)

	s.outputs = append(s.outputs, outputPipe)

	return outputPipe.Output(), nil
}

func (s *shuntPipe) Input() chan<- []byte {
	return s.input.Input()
}

func (s *shuntPipe) IsClosed() bool {
	return s.checkState() != nil
}

//nolint:gocognit
func (s *shuntPipe) Start() error {
	if err := s.checkState(); err != nil {
		return err
	}

	s.startLocker.Lock()
	defer s.startLocker.Unlock()
	defer s.cleanAll()

	// clean outputs loop
	go func() {
		for {
			rmCount := 0
			s.mu.Lock()
			n := len(s.outputs)
			for i := 0; i < n; i++ {
				idx := i - rmCount
				// check child output pipe is closed or not
				if s.outputs[idx].IsClosed() { // is done
					s.outputs = append(s.outputs[:idx], s.outputs[idx+1:]...) // remove output pipe from shunt pipe
					rmCount++

					continue
				}
			}

			if len(s.outputs) == 0 {
				s.mu.Unlock()
				s.cleanAll()

				return
			}

			s.mu.Unlock()
			time.Sleep(time.Millisecond)
		}
	}()

	// run shunt pipe
	for {
		select {
		case <-s.ctx.Done():
			return nil
		case data := <-s.input.Output():
			s.mu.RLock()
			n := len(s.outputs)
			for i := 0; i < n; i++ {
				output := s.outputs[i]

				// each output pipe transmission
				go func(outputPipe Pipe) {
					defer func() {
						pInfo := recover()
						if pInfo == nil {
							return
						}
						if !outputPipe.IsClosed() {
							logV1.Warnf("panic transferring data to output pipe. %s", pInfo)
							outputPipe.Close()
						}
					}()
					select {
					case outputPipe.Input() <- data:
						return
					case <-time.After(time.Second * 30):
						outputPipe.Close()

						return
					}
				}(output)
			}

			if len(s.outputs) == 0 {
				return nil
			}

			s.mu.RUnlock()
		}
	}
}

func (s *shuntPipe) Close() error {
	if err := s.checkState(); err != nil {
		return err
	}

	go s.cleanAll()
	select {
	case <-s.ctx.Done():
		return nil
	case <-time.After(time.Minute * 2):
		return errors.New("shunt pipe close timeout")
	}
}

func (s *shuntPipe) checkState() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.isClosed {
		return errors.New("shunt pipe is closed")
	}

	return nil
}

func (s *shuntPipe) cleanAll() {
	if s.IsClosed() {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	defer func() { s.isClosed = true }() // set shunt pipe to be closed
	s.cancel()                           // close shunt pipe's context
	for _, output := range s.outputs {   // close all output pipes
		output.Close()
	}
}

func NewShuntPipe(maxConns int, outputTimeout time.Duration) (*shuntPipe, error) {
	cctx, cancel := context.WithCancel(context.Background())

	return &shuntPipe{
		startLocker:    new(sync.Mutex),
		mu:             sync.RWMutex{},
		ctx:            cctx,
		cancel:         cancel,
		maxConnections: maxConns,
		outputTimeout:  outputTimeout,
		input:          NewPipe(cctx),
		outputs:        make([]Pipe, 0),
	}, nil
}
