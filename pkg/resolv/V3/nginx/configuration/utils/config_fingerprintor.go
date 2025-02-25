package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

type ConfigFingerprints map[string]string

type ConfigFingerprinter interface {
	Diff(fingerprints ConfigFingerprints) bool
	NewerThan(timestamp time.Time) bool
	Renew(fingerprints ConfigFingerprints)
	Fingerprints() ConfigFingerprints
}

type configFingerprinter struct {
	fingerprints ConfigFingerprints
	timestamp    time.Time
	locker       *sync.Mutex
}

func (f *configFingerprinter) Diff(fingerprints ConfigFingerprints) bool {
	if len(fingerprints) == len(f.fingerprints) {
		for filename, fingerprint := range fingerprints {
			if localFingerprint, ok := f.fingerprints[filename]; ok &&
				localFingerprint == fingerprint {
				continue
			}
			return true
		}
		return false
	}
	return true
}

func (f *configFingerprinter) NewerThan(timestamp time.Time) bool {
	return f.timestamp.After(timestamp)
}

func (f *configFingerprinter) Renew(fingerprints ConfigFingerprints) {
	if !f.Diff(fingerprints) {
		return
	}
	f.locker.Lock()
	defer f.locker.Unlock()
	defer func() {
		f.timestamp = time.Now()
	}()
	f.fingerprints = make(ConfigFingerprints)
	for name, fingerprint := range fingerprints {
		f.fingerprints[name] = fingerprint
	}
}

func (f *configFingerprinter) Fingerprints() ConfigFingerprints {
	fp := make(ConfigFingerprints)
	for s := range f.fingerprints {
		fp[s] = f.fingerprints[s]
	}
	return fp
}

func (f *configFingerprinter) setFingerprint(filename string, data []byte) {
	hash := sha256.New()
	hash.Write(data)
	f.fingerprints[filename] = hex.EncodeToString(hash.Sum(nil))
}

func NewConfigFingerprinter(buffMap map[string]*bytes.Buffer) ConfigFingerprinter {
	return NewConfigFingerprinterWithTimestamp(buffMap, time.Now())
}

func NewConfigFingerprinterWithTimestamp(buffMap map[string]*bytes.Buffer, timestamp time.Time) ConfigFingerprinter {
	cf := &configFingerprinter{
		fingerprints: make(ConfigFingerprints),
		locker:       new(sync.Mutex),
		timestamp:    timestamp,
	}
	for s, buff := range buffMap {
		cf.setFingerprint(s, buff.Bytes())
	}
	return cf
}

func SimpleConfigFingerprinter(fingerprints ConfigFingerprints) ConfigFingerprinter {
	return &configFingerprinter{
		fingerprints: fingerprints,
		timestamp:    time.Now(),
		locker:       new(sync.Mutex),
	}
}
