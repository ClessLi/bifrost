package nginx

import "errors"

var (
	ErrLogsDirPath         = errors.New("logs dir is not exist or is not a directory")
	ErrLogBufferIsNotExist = errors.New("log buffer is not exist")
	ErrLogBufferIsExist    = errors.New("log buffer is exist")
	ErrLogIsLocked         = errors.New("log is locked")
	ErrUnknownLockError    = errors.New("unknown lock error")
)
