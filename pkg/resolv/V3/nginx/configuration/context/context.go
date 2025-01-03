package context

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
)

type Context interface {
	// write `Context` methods
	Insert(ctx Context, idx int) Context
	Remove(idx int) Context
	Modify(ctx Context, idx int) Context
	// read `Context` methods
	Father() Context
	Child(idx int) Context

	QueryByKeyWords(kw KeyWords) Pos
	QueryAllByKeyWords(kw KeyWords) []Pos

	Clone() Context

	// write methods
	SetValue(v string) error
	SetFather(ctx Context) error
	// read methods
	HasChild() bool
	Len() int
	Value() string
	Type() context_type.ContextType
	Error() error

	// dump
	ConfigLines(isDumping bool) ([]string, error)

	// Enable/Disable conversion methods
	IsEnabled() bool
	Enable() Context
	Disable() Context
}
