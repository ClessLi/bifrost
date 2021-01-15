package dump_cacher

import (
	"bytes"
	"errors"
)

type DumpCacher interface {
	Write(k string, b []byte)
	Read(k string) ([]byte, error)
	ReadAll() map[string][]byte
}

var (
	ErrCacheNotExist = errors.New("cache not exist")
)

type dumpCache map[string]*bytes.Buffer

func (d *dumpCache) Write(k string, b []byte) {
	buff, ok := (*d)[k]
	if ok {
		buff.Write(b)
	} else {
		(*d)[k] = bytes.NewBuffer(b)
	}
}

func (d dumpCache) Read(k string) ([]byte, error) {
	buff, ok := d[k]
	if !ok {
		return nil, ErrCacheNotExist
	}
	return buff.Bytes(), nil
}

func (d dumpCache) ReadAll() map[string][]byte {
	dumps := make(map[string][]byte)
	for k := range d {
		dumps[k] = d[k].Bytes()
	}
	return dumps
}

func NewDumpCacher(k string, b []byte) DumpCacher {
	cacher := make(dumpCache)
	cacher.Write(k, b)
	return &cacher
}
