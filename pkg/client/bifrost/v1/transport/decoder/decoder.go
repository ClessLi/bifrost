// Package decoder decode response
package decoder

import "context"

type Decoder interface {
	DecodeResponse(ctx context.Context, resp interface{}) (interface{}, error)
}

type Factory interface {
	WebServerConfig() Decoder
	WebServerStatistics() Decoder
	WebServerStatus() Decoder
}

type factory struct{}

func (f factory) WebServerConfig() Decoder {
	return new(webConfigServer)
}

func (f factory) WebServerStatistics() Decoder {
	return new(webServerStatistics)
}

func (f factory) WebServerStatus() Decoder {
	return new(webServerStatus)
}

var _ Factory = factory{}

func New() Factory {
	return new(factory)
}
