package decoder

import "context"

// Decoder defines the decoder interface for grpc request.
type Decoder interface {
	DecodeRequest(ctx context.Context, r interface{}) (interface{}, error)
}
