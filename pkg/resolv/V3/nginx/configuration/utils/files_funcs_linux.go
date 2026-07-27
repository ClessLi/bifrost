package utils

import (
	"syscall"
	"time"

	"github.com/ClessLi/bifrost/internal/pkg/code"
)

func FileModifyTime(filepath string) (*time.Time, error) {
	var fs syscall.Stat_t
	err := syscall.Stat(filepath, &fs)
	if err != nil && !code.Is(err, syscall.EINTR) {
		return nil, err
	}
	tt := time.Unix(fs.Mtim.Sec, fs.Mtim.Nsec)

	return &tt, nil
}
