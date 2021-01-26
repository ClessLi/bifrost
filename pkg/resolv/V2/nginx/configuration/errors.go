package configuration

import (
	"errors"
)

var (
	ErrConfigurationTypeMismatch = errors.New("configuration type mismatch")
	ErrSameConfigFingerprint     = errors.New("same config fingerprint")

	NoBackupRequired = errors.New("no backup required")
	NoReloadRequired = errors.New("no reload required")
)
