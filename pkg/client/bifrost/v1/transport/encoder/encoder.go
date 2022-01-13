// Package encoder encode request
package encoder

import "context"

type Encoder interface {
	EncodeRequest(ctx context.Context, req interface{}) (interface{}, error)
}

type Factory interface {
	WebServerConfig() Encoder
	WebServerStatistics() Encoder
}

type factory struct{}

func (f factory) WebServerConfig() Encoder {
	return new(webServerConfig)
}

func (f factory) WebServerStatistics() Encoder {
	return new(webServerStatistics)
}

var _ Factory = factory{}

func New() Factory {
	return new(factory)
}
