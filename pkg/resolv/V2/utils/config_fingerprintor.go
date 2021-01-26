package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

type ConfigFingerprinter interface {
	Diff(fingerprinter ConfigFingerprinter) bool
}

type configFingerprinter struct {
	fingerprints map[string]string
}

func (f *configFingerprinter) Diff(fingerprinter ConfigFingerprinter) bool {
	if fp, ok := fingerprinter.(*configFingerprinter); ok && len(fp.fingerprints) == len(f.fingerprints) {
		for filename, fingerprint := range fp.fingerprints {
			if localFingerprint, ok := f.fingerprints[filename]; ok && strings.EqualFold(localFingerprint, fingerprint) {
				continue
			}
			return true
		}
		return false
	}
	return true
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
