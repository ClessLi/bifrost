package context

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
)

type ContextBuilder interface {
	SetType(contextType context_type.ContextType) ContextBuilder
	SetValue(value string) ContextBuilder
	Build() Context
}

type BuildOptions struct {
	ContextType context_type.ContextType
	ParseType   ParseType
	HasValue    bool
}

type ParseType int

const (
	ParseConfig ParseType = iota
	ParseContext
	ParseDirective
	ParseComment
)

var buildFuncMap map[context_type.ContextType]func(value string) Context = make(map[context_type.ContextType]func(value string) Context)

func RegisterContextBuilder() {

}
