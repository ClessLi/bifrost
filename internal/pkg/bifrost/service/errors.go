package service

import "errors"

var (
	UnknownErrCheckToken = errors.New("an unknown error occurred while verifying token")
	UnknownRequestType   = errors.New("an unknown request type")
	ErrConnToAuthSvr     = errors.New("failed to connect to authentication server")
	// Monitor Error
	ErrStopMonitoringTimeout       = errors.New("stop monitoring timeout")
	ErrMonitoringServiceSuspension = errors.New("monitoring service suspension")
	ErrMonitoringStarted           = errors.New("monitoring started")
)
