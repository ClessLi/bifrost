package configuration

import (
	"errors"
)

var (
	ErrConfigurationTypeMismatch                         = errors.New("configuration type mismatch")
	ErrSameConfigFingerprint                             = errors.New("same config fingerprint")
	ErrSameConfigFingerprintBetweenFilesAndConfiguration = errors.New("same config fingerprint between files and configuration")

	ErrConfigManagerIsRunning    = errors.New("config manager is running")
	ErrConfigManagerIsNotRunning = errors.New("config manager is not running")
	//NoBackupRequired             = errors.New("no backup required")
	//NoReloadRequired             = errors.New("no reload required")
)
