package service

import "errors"

var (
	UnknownErrCheckToken  = errors.New("an unknown error occurred while verifying token")
	ErrDataNotParsed      = errors.New("config data not parsed")
	ErrUnknownSvrName     = errors.New("unknown server name")
	ErrProcessNotRunning  = errors.New("process is not running")
	ErrConnToAuthSvr      = errors.New("failed to connect to authentication server")
	ErrChanNil            = errors.New("the channel put in is nil")
	ErrDataSendingTimeout = errors.New("data sending timeout")
	ErrWatchLogTimeout    = errors.New("the WatchLog operation timed out")
)
