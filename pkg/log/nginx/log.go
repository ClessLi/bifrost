// logger.nginx 包，该包包含了bifrost.pkg.logger.nginx项目用于读取nginx日志的相关对象，及相关方法及函数
// 创建者： ClessLi
// 创建时间：2020-10-27 17:02:30
package nginx

import (
	"bytes"
	"github.com/hpcloud/tail"
	"os"
	"path/filepath"
)

// TODO: nginx日志读取包
type Log struct {
	log   map[string]*tail.Tail
	locks map[string]bool
}

func (l *Log) StartWatch(logName, workspace string) (err error) {
	isLocked, ok := l.locks[logName]
	if ok && isLocked {
		return ErrLogIsLocked
	}
	dir, err := os.Stat(workspace)
	if err != nil {
		return err
	}
	if !dir.IsDir() {
		return ErrLogsDirPath
	}

	if _, ok := l.log[logName]; ok {
		return ErrLogBufferIsExist
	}

	logPath, err := filepath.Abs(filepath.Join(workspace, logName))
	if err != nil {
		return err
	}

	tails, err := tail.TailFile(logPath, tail.Config{
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: 2,
		},
		ReOpen:    true,
		MustExist: true,
		Poll:      true,
		Follow:    true,
	})
	defer func() {
		if err != nil {
			l.StopWatch(logName)
		}
	}()
	if err != nil {
		return err
	}
	l.log[logName] = tails
	l.locks[logName] = true
	return nil
}

func (l *Log) Watch(logName string) (logData []byte, err error) {
	if isLocked, ok := l.locks[logName]; !ok || !isLocked {
		_ = l.StopWatch(logName)
		return nil, ErrUnknownLockError
	}
	tails, ok := l.log[logName]
	if !ok {
		return nil, ErrLogBufferIsNotExist
	}
	// TODO: 优化缓存回收机制，类似rpc存根回收新增机制
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	select {
	case line, ok := <-tails.Lines:
		if ok {
			buf.Write([]byte(line.Text))
			buf.Write([]byte("\n"))
		} else {
			break
		}
	}
	logData = buf.Bytes()
	return logData, err
}

func (l *Log) StopWatch(logName string) error {
	tails, ok := l.log[logName]
	if !ok {
		return ErrLogBufferIsNotExist
	}
	//fmt.Printf("stop tail %s\n", logName)
	err := tails.Stop()
	delete(l.log, logName)
	delete(l.locks, logName)
	//fmt.Printf("stop tail %s compelet\n", logName)
	return err
}

func NewLog() *Log {
	return &Log{
		log:   make(map[string]*tail.Tail),
		locks: make(map[string]bool),
	}
}
