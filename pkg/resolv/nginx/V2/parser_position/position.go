package parser_position

import (
	"strings"
	"sync"
)

const INDENT = "    "

type ParserPosition interface {
	ConfigAbsPath() string
	ConfigDeep() int
	GlobalDeep() int
	//SetConfigDeep(int)
	SetGlobalDeep(int)
	ConfigIndents() string
	GlobalIndents() string
	NextPosition() ParserPosition
}

type parserPosition struct {
	configAbsPath string
	configDeep    int
	globalDeep    int
	indent        string
	rwLock        *sync.RWMutex
}

func (p parserPosition) ConfigAbsPath() string {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.configAbsPath
}

func (p parserPosition) ConfigDeep() int {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.configDeep
}

func (p parserPosition) GlobalDeep() int {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.globalDeep
}

//func (p *parserPosition) SetConfigDeep(deep int) {
//	p.rwLock.Lock()
//	defer p.rwLock.Unlock()
//	p.configDeep = deep
//}

func (p *parserPosition) SetGlobalDeep(deep int) {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	p.globalDeep = deep
}

func (p parserPosition) ConfigIndents() string {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.indents(p.configDeep)
}

func (p parserPosition) GlobalIndents() string {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.indents(p.globalDeep)
}

func (p parserPosition) NextPosition() ParserPosition {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return NewParserPosition(p.configAbsPath, p.configDeep+1, p.globalDeep+1)
}

func (p parserPosition) indents(deep int) string {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return strings.Repeat(p.indent, deep)
}

func NewParserPosition(path string, configDeep, globalDeep int) ParserPosition {
	return &parserPosition{
		configAbsPath: path,
		configDeep:    configDeep,
		globalDeep:    globalDeep,
		indent:        INDENT,
		rwLock:        new(sync.RWMutex),
	}
}
