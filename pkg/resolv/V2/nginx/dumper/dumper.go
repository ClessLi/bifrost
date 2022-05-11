package dumper

import (
	"bytes"
	"errors"
	"fmt"
)

type Dumper interface {
	Len(k string) int
	Truncate(k string, n int) error
	Write(k string, b []byte)
	Done(k string) error
	Read(k string) ([]byte, error)
	ReadAll() map[string][]byte
}

var ErrCacheNotExist = errors.New("cache not exist")

type dumper struct {
	cache   map[string]*bytes.Buffer
	doneMap map[string]bool
}

func (d *dumper) Truncate(k string, n int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	if buff, ok := d.cache[k]; ok {
		if d.isDone(k) {
			return nil
		}
		buff.Truncate(n)

		return nil
	}

	return ErrCacheNotExist
}

func (d dumper) Len(k string) int {
	if buff, ok := d.cache[k]; ok {
		return buff.Len()
	}

	return 0
}

func (d *dumper) Write(k string, b []byte) {
	buff, ok := d.cache[k]
	if ok {
		if d.isDone(k) {
			return
		}
		buff.Write(b)
	} else {
		d.cache[k] = bytes.NewBuffer(b)
	}
}

func (d *dumper) Done(k string) error {
	if _, ok := d.cache[k]; !ok {
		return ErrCacheNotExist
	}
	if !d.isDone(k) {
		d.doneMap[k] = true
	}

	return nil
}

func (d *dumper) isDone(k string) bool {
	isDone, ok := d.doneMap[k]
	if ok {
		return isDone
	}
	d.doneMap[k] = false

	return false
}

func (d dumper) Read(k string) ([]byte, error) {
	buff, ok := d.cache[k]
	if !ok {
		return nil, ErrCacheNotExist
	}

	return buff.Bytes(), nil
}

func (d dumper) ReadAll() map[string][]byte {
	dumps := make(map[string][]byte)
	for k := range d.cache {
		dumps[k] = d.cache[k].Bytes()
	}

	return dumps
}

func NewDumper(k string) Dumper {
	d := &dumper{
		cache:   make(map[string]*bytes.Buffer),
		doneMap: make(map[string]bool),
	}
	d.Write(k, []byte(""))

	return d
}
