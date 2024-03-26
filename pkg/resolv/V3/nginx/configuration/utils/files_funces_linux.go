package utils

import (
	"github.com/marmotedu/errors"
	"syscall"
	"time"
)

func FileModifyTime(filepath string) (*time.Time, error) {
	var fs syscall.Stat_t
	err := syscall.Stat(filepath, &fs)
	if err != nil && !errors.Is(err, syscall.EINTR) {
		return nil, err
	}
	tt := time.Unix(fs.Mtim.Sec, fs.Mtim.Nsec)
	return &tt, nil
}
