// Package encoder encode request
package encoder

import "context"

type Encoder interface {
	EncodeRequest(ctx context.Context, req interface{}) (interface{}, error)
}

type Factory interface {
	WebServerConfig() Encoder
	WebServerStatistics() Encoder
	WebServerStatus() Encoder
	WebServerLogWatcher() Encoder
	WebServerBinCMD() Encoder
}

type factory struct{}

func (f factory) WebServerConfig() Encoder {
	return new(webServerConfig)
}

func (f factory) WebServerStatistics() Encoder {
	return new(webServerStatistics)
}

func (f factory) WebServerStatus() Encoder {
	return new(webServerStatus)
}

func (f factory) WebServerLogWatcher() Encoder {
	return new(webServerLogWatcher)
}

func (f factory) WebServerBinCMD() Encoder {
	return new(webServerBinCMD)
}

var _ Factory = factory{}

func New() Factory {
	return new(factory)
}
