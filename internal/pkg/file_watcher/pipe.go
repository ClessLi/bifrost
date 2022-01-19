package file_watcher

import (
	"context"
	"github.com/marmotedu/errors"
	"sync"
	"time"
)

type ShuntPipe struct {
	mu             sync.RWMutex
	ctx            context.Context
	cancel         context.CancelFunc
	maxConnections int
	outputTimeout  time.Duration
	input          *InputPipe
	outputs        []*OutputPipe
}

func (s *ShuntPipe) AddOutput(outputC chan []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ctx.Err() != nil {
		return errors.Wrap(s.ctx.Err(), "shunt pipe is closed. %s")
	}
	if s.maxConnections <= len(s.outputs) {
		return errors.Errorf("the number of connections has reached the maximum (%d/%d)", len(s.outputs), s.maxConnections)
	}

	output := newOutputPipe(s.ctx, s.outputTimeout, outputC)

	go func() {
		defer close(output.c)
		defer output.cancel()
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-output.ctx.Done():
				return
			}
		}
	}()

	s.outputs = append(s.outputs, output)
	return nil

}

func (s *ShuntPipe) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer s.cancel()
	defer s.input.cancel()
	for {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		case data := <-s.input.c:
			rmCount := 0
			n := len(s.outputs)
			for i := 0; i < n; i++ {
				idx := i - rmCount
				if s.outputs[idx].ctx.Err() != nil {
					s.outputs = append(s.outputs[:idx], s.outputs[idx+1:]...)
					rmCount++
					continue
				}
				output := s.outputs[idx]
				go func(pipe *OutputPipe) {
					select {
					case pipe.c <- data:
						return
					case <-pipe.ctx.Done():
						return
					case <-time.After(time.Second * 30):
						pipe.cancel()
						return
					}
				}(output)
			}

			if len(s.outputs) == 0 {
				return nil
			}
		}
	}
}

func (s *ShuntPipe) Close() error {
	timeoutClose, cancel := context.WithTimeout(s.input.ctx, time.Minute*2)
	defer cancel()
	go func() {
		s.input.cancel()
	}()
	select {
	case <-s.ctx.Done():
		return nil
	case <-timeoutClose.Done():
		return errors.New("shunt pipe close timeout")
	}
}

func NewShuntPipe(maxConns int, outputTimeout time.Duration, input *InputPipe) (*ShuntPipe, error) {
	if input.ctx.Err() != nil {
		return nil, errors.Errorf("input pipe is closed. %s", input.ctx.Err().Error())
	}
	ctx, cancel := context.WithCancel(input.ctx)
	return &ShuntPipe{
		mu:             sync.RWMutex{},
		ctx:            ctx,
		cancel:         cancel,
		maxConnections: maxConns,
		outputTimeout:  outputTimeout,
		input:          input,
		outputs:        make([]*OutputPipe, 0),
	}, nil
}

type InputPipe struct {
	ctx    context.Context
	cancel context.CancelFunc
	c      <-chan []byte
}

func NewInputPipe(ctx context.Context, c <-chan []byte) *InputPipe {
	cctx, cancel := context.WithCancel(ctx)
	return &InputPipe{
		ctx:    cctx,
		cancel: cancel,
		c:      c,
	}
}

type OutputPipe struct {
	ctx    context.Context
	cancel context.CancelFunc
	c      chan []byte
}

func newOutputPipe(ctx context.Context, timeout time.Duration, c chan []byte) *OutputPipe {
	cctx, cancel := context.WithTimeout(ctx, timeout)
	return &OutputPipe{
		ctx:    cctx,
		cancel: cancel,
		c:      c,
	}
}
