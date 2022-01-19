package file_watcher

import (
	"context"
	log "github.com/ClessLi/bifrost/pkg/log/v1"
	"github.com/hpcloud/tail"
	"github.com/marmotedu/errors"
	"os"
	"time"
)

type FileWatcher struct {
	filePath string
	isClosed bool
	inputC   chan []byte
	pipe     *ShuntPipe
}

func (f *FileWatcher) start() error {
	defer func() { f.isClosed = true }()
	defer close(f.inputC)
	go func() {
		err := f.pipe.Start()
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

	for {
		select {
		case f.inputC <- []byte((<-t.Lines).Text):
		case <-f.pipe.ctx.Done():
			return f.pipe.ctx.Err()
		}
	}
}

func (f *FileWatcher) Stop() error {
	go func() {
		err := f.pipe.Close()
		if err != nil {
			log.Warnf("file '%s' stop watching error. %s", f.filePath, err.Error())
		}
	}()
	select {
	case <-f.pipe.ctx.Done():
		return f.pipe.ctx.Err()
	case <-time.After(time.Second * 30):
		return errors.Errorf("stop watching file '%s' timeout(30s)", f.filePath)
	}
}

func (f *FileWatcher) AddOutput(outputC chan []byte) error {
	return f.pipe.AddOutput(outputC)
}

func newFileWatcher(config *CompletedConfig) (*FileWatcher, error) {
	inputC := make(chan []byte)
	pipe, err := NewShuntPipe(config.MaxConnections, config.OutputTimeout, NewInputPipe(context.Background(), inputC))
	if err != nil {
		return nil, err
	}
	return &FileWatcher{
		filePath: config.filePath,
		isClosed: false,
		inputC:   inputC,
		pipe:     pipe,
	}, nil
}
