package file_watcher

import (
	"context"
	log "github.com/ClessLi/bifrost/pkg/log/v1"
	"github.com/hpcloud/tail"
	"github.com/marmotedu/errors"
	"os"
	"sync"
	"time"
)

type FileWatcher struct {
	filePath    string
	ctx         context.Context
	cancel      context.CancelFunc
	startLocker sync.Locker

	shuntPipe ShuntPipe
}

func (f *FileWatcher) Output() (<-chan []byte, error) {
	return f.shuntPipe.Output()
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
			log.Warnf("file '%s' watching error. %s", f.filePath, err.Error())
		}
	}()

	t, err := tail.TailFile(f.filePath, tail.Config{
		Logger: log.StdInfoLogger(),
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

	defer func(t *tail.Tail) {
		err := t.Stop()
		if err != nil {
			log.Warnf("tail stop error. %s", err.Error())
		}
	}(t)

	defer func() {
		pInfo := recover()
		if pInfo == nil {
			return
		}
		if !f.shuntPipe.IsClosed() {
			log.Warnf("panic transferring data to shunt pipe. %s", pInfo)
			err := f.shuntPipe.Close()
			if err != nil {
				log.Warnf("failed to stop shunt pipe. %s", err.Error())
			}
		}
	}()
	for {
		select {
		case f.shuntPipe.Input() <- []byte((<-t.Lines).Text):
		case <-f.ctx.Done():
			return nil
		}
	}
}

func (f *FileWatcher) Stop() error {
	defer f.cancel()
	if f.shuntPipe.IsClosed() {
		return errors.Errorf("failed to stop '%s' file watcher, shunt pipe is already closed", f.filePath)
	}

	err := f.closePipe()
	if err != nil {
		return errors.Errorf("failed to stop '%s' file watcher. %s", f.filePath, err.Error())
	}
	return nil
}

func (f *FileWatcher) IsClosed() bool {
	return f.shuntPipe.IsClosed()
}

func (f *FileWatcher) closePipe() error {
	errC := make(chan error)
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