package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

type ConfigFingerprinter interface {
	Diff(fingerprinter ConfigFingerprinter) bool
	Renew(fingerprinter ConfigFingerprinter)
}

type configFingerprinter struct {
	fingerprints map[string]string
}

func (f *configFingerprinter) Diff(fingerprinter ConfigFingerprinter) bool {
	if fp, is := fingerprinter.(*configFingerprinter); is && len(fp.fingerprints) == len(f.fingerprints) {
		for filename, fingerprint := range fp.fingerprints {
			if localFingerprint, ok := f.fingerprints[filename]; ok &&
				strings.EqualFold(localFingerprint, fingerprint) {
				continue
			}
			return true
		}
		return false
	}
	return true
}

func (f *configFingerprinter) Renew(fingerprinter ConfigFingerprinter) {
	if !f.Diff(fingerprinter) {
		return
	}
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

func NewConfigFingerprinter(dataMap map[string][]byte) ConfigFingerprinter {
	cf := &configFingerprinter{fingerprints: make(map[string]string)}
	for s, bytes := range dataMap {
		cf.setFingerprint(s, bytes)
	}
	return cf
}
