package file_watcher

import (
	"context"
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"os"
	"sync"
	"time"

	"github.com/hpcloud/tail"
	"github.com/marmotedu/errors"
)

type FileWatcher struct {
	filePath    string
	ctx         context.Context
	cancel      context.CancelFunc
	startLocker sync.Locker

	shuntPipe ShuntPipe
}

func (f *FileWatcher) Output(ctx context.Context) (<-chan []byte, error) {
	return f.shuntPipe.Output(ctx)
}

func (f *FileWatcher) Start() error {
	if f.shuntPipe.IsClosed() {
		f.cancel()

		return errors.Errorf("failed to start '%s' file watcher, shunt pipe is already closed", f.filePath)
	}

	f.startLocker.Lock()
	defer f.startLocker.Unlock()
	defer f.cancel()
	go func() {
		err := f.shuntPipe.Start()
		if err != nil {
			logV1.Warnf("file '%s' watching error. %s", f.filePath, err.Error())
		}
	}()

	logV1.Debugf("tail '%s' starting...", f.filePath)
	t, err := tail.TailFile(f.filePath, tail.Config{
		Logger: logV1.StdInfoLogger(),
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: os.SEEK_END,
		},
		ReOpen:    true,
		MustExist: true,
		Poll:      true,
		Follow:    true,
	})
	if err != nil {
		return err
	}

	// defer to stop tail
	defer func(t *tail.Tail) {
		logV1.Debugf("tail '%s' stopping...", f.filePath)
		if err = t.Stop(); err != nil {
			logV1.Warnf("tail stop error. %s", err.Error())
		}
	}(t)

	// recover panic for watching loop
	defer func() {
		pInfo := recover()
		if pInfo == nil {
			return
		}
		if !f.shuntPipe.IsClosed() {
			logV1.Warnf("panic transferring data to shunt pipe. %s", pInfo)
			if err = f.shuntPipe.Close(); err != nil {
				logV1.Warnf("failed to stop shunt pipe. %s", err.Error())
			}
		}
	}()
	// FileWatcher watching
	for !f.shuntPipe.IsClosed() {
		select {
		case line := <-t.Lines: // receive tail line
			select {
			case f.shuntPipe.Input() <- []byte(line.Text): // send to shut pipe
			case <-time.After(time.Second * 30): // send to shut pipe timeout(30s)
				err = errors.Errorf("send data to shut pipe timeout(30s), file: %s", f.filePath)

				return err
			}
		case <-f.ctx.Done(): // FileWatcher Close method has been called
			logV1.Debugf("watching file '%s' completed")

			return nil
		default: // sleep 1ms and return to loop with shut pipe closed check
			time.Sleep(time.Millisecond)
		}
	}

	return nil
}

func (f *FileWatcher) Stop() error {
	defer f.cancel()
	if f.shuntPipe.IsClosed() {
		return errors.Errorf("failed to stop '%s' file watcher, shunt pipe is already closed", f.filePath)
	}

	if err := f.closePipe(); err != nil {
		return errors.Errorf("failed to stop '%s' file watcher. %s", f.filePath, err.Error())
	}

	return nil
}

func (f *FileWatcher) IsClosed() bool {
	return f.shuntPipe.IsClosed()
}

func (f *FileWatcher) closePipe() error {
	errC := make(chan error, 1)
	select {
	case errC <- f.shuntPipe.Close():
		err := <-errC

		return err
	case <-time.After(time.Second * 30):
		return errors.New("close pipe timeout(30s)")
	}
}

func newFileWatcher(ctx context.Context, config *CompletedConfig) (*FileWatcher, error) {
	cctx, cancel := context.WithCancel(ctx)
	sp, err := NewShuntPipe(config.MaxConnections, config.OutputTimeout)
	if err != nil {
		cancel()

		return nil, err
	}

	return &FileWatcher{
		filePath:    config.filePath,
		ctx:         cctx,
		cancel:      cancel,
		startLocker: new(sync.Mutex),
		shuntPipe:   sp,
	}, nil
}
