package parser_indention

import (
	"strings"
	"sync"
)

const INDENT = "    "

type Indention interface {
	SetGlobalDeep(deep int)
	GlobalDeep() int
	GlobalIndents() string
	ConfigIndents() string
	NextIndention() Indention
}

type indention struct {
	configDeep int
	globalDeep int
	indent     string
	locker     *sync.RWMutex
}

func (i *indention) SetGlobalDeep(deep int) {
	i.locker.Lock()
	defer i.locker.Unlock()
	i.globalDeep = deep
}

func (i indention) GlobalDeep() int {
	i.locker.RLock()
	defer i.locker.RUnlock()

	return i.globalDeep
}

func (i indention) GlobalIndents() string {
	i.locker.RLock()
	defer i.locker.RUnlock()

	return strings.Repeat(i.indent, i.globalDeep)
}

func (i indention) ConfigIndents() string {
	return strings.Repeat(i.indent, i.configDeep)
}

func (i indention) NextIndention() Indention {
	i.locker.RLock()
	defer i.locker.RUnlock()

	return &indention{
		configDeep: i.configDeep + 1,
		globalDeep: i.globalDeep + 1,
		indent:     i.indent,
		locker:     new(sync.RWMutex),
	}
}

func NewIndention() Indention {
	return &indention{
		configDeep: 0,
		globalDeep: 0,
		indent:     INDENT,
		locker:     new(sync.RWMutex),
	}
}
