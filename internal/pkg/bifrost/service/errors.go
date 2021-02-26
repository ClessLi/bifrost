package service

import "errors"

var (
	UnknownErrCheckToken              = errors.New("an unknown error occurred while verifying token")
	UnknownRequestType                = errors.New("an unknown request type")
	ErrConnToAuthSvr                  = errors.New("failed to connect to authentication server")
	ErrWebServerConfigServiceNotExist = errors.New("web server config service is not exist")
	// Monitor Error
	ErrStopMonitoringTimeout       = errors.New("stop monitoring timeout")
	ErrMonitoringServiceSuspension = errors.New("monitoring service suspension")
	ErrMonitoringStarted           = errors.New("monitoring started")
	// offstage Error
	ErrDataSendingTimeout     = errors.New("data sending timeout")
	ErrWatchLogTimeout        = errors.New("the WatchLog operation timed out")
	ErrLogWatcherCloseTimeout = errors.New("the LogWatcher Close timed out")
	// responseInfo Error
	ErrInvalidResponseInfo = errors.New("invalid response info object")
)
