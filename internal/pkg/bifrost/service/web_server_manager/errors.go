package web_server_manager

import (
	"errors"
)

var (
	ErrDisabledService       = errors.New("config agent out of service")
	ErrDataNotParsed         = errors.New("config data not parsed")
	ErrDataSendingTimeout    = errors.New("data sending timeout")
	ErrWatchLogTimeout       = errors.New("the WatchLog operation timed out")
	ErrValidationNotExist    = errors.New("the validation process does not exist or is configured incorrectly")
	ErrEmptyConfig           = errors.New("empty configuration error")
	ErrOffstageNotExist      = errors.New("offstage is not exist")
	ErrWrongStateExpectation = errors.New("wrong state expectation")
)
