package encoder

import "context"

type Encoder interface {
	EncodeResponse(ctx context.Context, r interface{}) (interface{}, error)
}
