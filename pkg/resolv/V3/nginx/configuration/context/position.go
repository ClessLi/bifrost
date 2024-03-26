package context

import (
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/marmotedu/errors"
)

type Pos interface {
	Position() (Context, int)
	Target() Context
}

type pos struct {
	father     Context
	childIndex int
}

func (p pos) Position() (Context, int) {
	return p.father, p.childIndex
}

func (p pos) Target() Context {
	return p.father.Child(p.childIndex)
}

func SetPos(father Context, posIdx int) Pos {
	return &pos{
		father:     father,
		childIndex: posIdx,
	}
}

type errPos struct {
	ctx *ErrorContext
}

func (e errPos) Position() (Context, int) {
	return e.ctx, -1
}

func (e errPos) Target() Context {
	return e.ctx
}

var nullPos = &errPos{nullContext}

func NullPos() Pos {
	return nullPos
}

var notFoundPos = &errPos{ErrContext(errors.WithCode(code.V3ErrContextNotFound, "queried context not found")).(*ErrorContext)}

func NotFoundPos() Pos {
	return notFoundPos
}
