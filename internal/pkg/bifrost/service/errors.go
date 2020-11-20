package service

import "errors"

var (
	UnknownErrCheckToken = errors.New("an unknown error occurred while verifying token")
	ErrDataNotParsed     = errors.New("config data not parsed")
	ErrUnknownSvrName    = errors.New("unknown server name")
	ErrProcessNotRunning = errors.New("process is not running")
	ErrConnToAuthSvr     = errors.New("failed to connect to authentication server")
)
