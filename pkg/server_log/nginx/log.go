// Package nginx 包含了bifrost.pkg.logger.nginx项目用于读取nginx日志的相关对象，及相关方法及函数
// 创建者： ClessLi
// 创建时间：2020-10-27 17:02:30
package nginx

import (
	"bytes"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/hpcloud/tail"
	"github.com/marmotedu/errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// TODO: nginx日志读取包
type Log struct {
	logs  map[string]*logBuffer
	locks map[string]bool
}

type logBuffer struct {
	tail   *tail.Tail
	buffer *bytes.Buffer
	rwLock *sync.Mutex
}

func newLogBuffer(t *tail.Tail, buf *bytes.Buffer) *logBuffer {
	return &logBuffer{
		tail:   t,
		buffer: buf,
		rwLock: new(sync.Mutex),
	}
}

func (l Log) getLog(logName string) (*logBuffer, error) {
	if isLocked, ok := l.locks[logName]; !ok || !isLocked {
		//_ = l.StopWatch(logName)
		return nil, errors.WithCode(code.ErrUnknownLockError, "the log lock exception")
	}
	log, ok := l.logs[logName]
	if !ok {
		return nil, errors.WithCode(code.ErrLogBufferIsNotExist, "get log buffer failed.")
	}
	return log, nil
}

func (l *Log) StartWatch(logName, workspace string) (err error) {
	isLocked, ok := l.locks[logName]
	if ok && isLocked {
		return errors.WithCode(code.ErrLogIsLocked, "the log lock is locked")
	}

	defer func() {
		if err != nil {
			_ = l.StopWatch(logName)
		}
	}()
	dir, err := os.Stat(workspace)
	if err != nil {
		return err
	}
	if !dir.IsDir() {
		return errors.WithCode(code.ErrLogsDirPath, "logs dir is not a directory")
	}

	log, _ := l.getLog(logName)
	if log != nil {
		return errors.WithCode(code.ErrLogBufferIsExist, "log buffer is already exist.")
	}

	logPath, err := filepath.Abs(filepath.Join(workspace, logName))
	if err != nil {
		return err
	}

	t, err := tail.TailFile(logPath, tail.Config{
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: 2,
		},
		ReOpen:    true,
		MustExist: true,
		Poll:      true,
		Follow:    true,
	})

	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	l.logs[logName] = newLogBuffer(t, buf)
	l.locks[logName] = true
	go func() {
		var (
			log *logBuffer
			err error
		)

		for {
			log, err = l.getLog(logName)
			if err != nil || log == nil {
				return
			}
			select {
			case line, ok := <-log.tail.Lines:
				if ok {
					//fmt.Println("read lock")
					log.rwLock.Lock()
					log.buffer.Write([]byte(line.Text))
					log.buffer.Write([]byte("\n"))
					log.rwLock.Unlock()
					//fmt.Println("read unlock")
				} else {
					time.Sleep(time.Millisecond)
					break
				}
			}
		}
	}()
	return nil
}

func (l *Log) Watch(logName string) (logData []byte, err error) {
	defer func() {
		if err != nil {
			_ = l.StopWatch(logName)
		}
	}()

	// TODO: 优化缓存回收机制，类似rpc存根回收新增机制
	log, err := l.getLog(logName)
	if err != nil {
		return nil, err
	}
	//fmt.Println("write lock")
	log.rwLock.Lock()
	logData = log.buffer.Bytes()
	log.buffer.Reset()
	log.rwLock.Unlock()
	//fmt.Println("write unlock")
	return logData, nil
}

func (l *Log) StopWatch(logName string) error {
	log, ok := l.logs[logName]
	if !ok {
		return errors.WithCode(code.ErrLogBufferIsNotExist, "get log buffer failed")
	}
	//fmt.Printf("stop tail %s\n", logName)
	err := log.tail.Stop()
	delete(l.logs, logName)
	_, ok = l.locks[logName]
	if !ok {
		return errors.WithCode(code.ErrLogIsUnlocked, "the log lock is unlocked in advance")
	}
	delete(l.locks, logName)
	//fmt.Printf("stop tail %s compelet\n", logName)
	return err
}

func NewLog() *Log {
	return &Log{
		logs:  make(map[string]*logBuffer),
		locks: make(map[string]bool),
	}
}
