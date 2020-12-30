package service

import "errors"

var (
	UnknownErrCheckToken = errors.New("an unknown error occurred while verifying token")
	UnknownRequestType   = errors.New("an unknown request type")
	ErrUnknownSvrName    = errors.New("unknown server name")
	ErrConnToAuthSvr     = errors.New("failed to connect to authentication server")
	// Responder Error
	ErrNotWatcherResponse = errors.New("it's not a response from watcher")
	ErrNotBytesResponse   = errors.New("it's not a response with Bytes")
	ErrParamNotPassedIn   = errors.New("parameter not passed in")
	// Monitor Error
	ErrStopMonitoringTimeout = errors.New("stop monitoring timeout")
)
