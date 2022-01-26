package options

import "time"

type Options struct {
	ChunkSize          int
	RecvTimeoutMinutes time.Duration
}
