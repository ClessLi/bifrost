package v1

import "time"

type Options struct {
	ChunkSize   int
	RecvTimeout time.Duration
}
