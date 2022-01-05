package v1

import "context"

type Encoder interface {
	EncodeResponse(ctx context.Context, r interface{}) (interface{}, error)
}
