package utils

import (
	"syscall"
	"time"
)

func FileModifyTime(filepath string) (*time.Time, error) {
	filename16, err := syscall.UTF16PtrFromString(filepath)
	if err != nil {
		return nil, err
	}
	h, err := syscall.CreateFile(filename16, 0, 0, nil, syscall.OPEN_EXISTING, uint32(syscall.FILE_FLAG_BACKUP_SEMANTICS), 0)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(h)
	var i syscall.ByHandleFileInformation
	if err := syscall.GetFileInformationByHandle(h, &i); err != nil {
		return nil, err
	}
	tt := time.Unix(0, i.LastWriteTime.Nanoseconds())
	return &tt, nil
}
