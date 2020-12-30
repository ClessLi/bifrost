package web_server_manager

import (
	"errors"
)

var (
	ErrDataNotParsed       = errors.New("config data not parsed")
	ErrDataSendingTimeout  = errors.New("data sending timeout")
	ErrWatchLogTimeout     = errors.New("the WatchLog operation timed out")
	ErrValidationNotExist  = errors.New("the validation process does not exist or is configured incorrectly")
	ErrEmptyConfig         = errors.New("empty configuration error")
	ErrWrongParamPassedIn  = errors.New("wrong parameter passed in")
	ErrServiceNotAvailable = errors.New("bifrost service not available")
)
