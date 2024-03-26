package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

type ConfigFingerprinter interface {
	Diff(fingerprinter ConfigFingerprinter) bool
	NewerThan(timestamp time.Time) bool
	Renew(fingerprinter ConfigFingerprinter)
}

type configFingerprinter struct {
	fingerprints map[string]string
	timestamp    time.Time
	locker       *sync.Mutex
}

func (f *configFingerprinter) Diff(fingerprinter ConfigFingerprinter) bool {
	if fp, is := fingerprinter.(*configFingerprinter); is && len(fp.fingerprints) == len(f.fingerprints) {
		for filename, fingerprint := range fp.fingerprints {
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

func (f *configFingerprinter) Renew(fingerprinter ConfigFingerprinter) {
	if !f.Diff(fingerprinter) {
		return
	}
	f.locker.Lock()
	defer f.locker.Unlock()
	defer func() {
		f.timestamp = time.Now()
	}()
	if fp, is := fingerprinter.(*configFingerprinter); is {
		f.fingerprints = make(map[string]string)

		for name, fingerprint := range fp.fingerprints {
			f.fingerprints[name] = fingerprint
		}
	}
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
		fingerprints: make(map[string]string),
		locker:       new(sync.Mutex),
		timestamp:    timestamp,
	}
	for s, buff := range buffMap {
		cf.setFingerprint(s, buff.Bytes())
	}
	return cf
}
