package context

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"regexp"
	"strings"
)

type KeyWords interface {
	Match(ctx Context) bool
	Cascaded() bool
	SetCascaded(cascaded bool) KeyWords
}

type keywords struct {
	matchingType  context_type.ContextType
	matchingValue string
	isRegex       bool
	isCascaded    bool
}

func (k *keywords) Match(ctx Context) bool {
	// match context's type
	matched := ctx.Type() == k.matchingType
	// match context's value
	if matched {
		matched = false //nolint:wastedassign,ineffassign
		if k.isRegex {
			var err error
			matched, err = regexp.MatchString(k.matchingValue, ctx.Value())
			if err != nil {
				return false
			}
		} else {
			matched = strings.Contains(ctx.Value(), k.matchingValue)
		}
	}
	return matched
}

func (k *keywords) Cascaded() bool {
	return k.isCascaded
}

func (k *keywords) SetCascaded(cascaded bool) KeyWords {
	k.isCascaded = cascaded
	return k
}

func NewKeyWords(ctxtype context_type.ContextType, matching string, isregex, iscascaded bool) KeyWords {
	return &keywords{
		matchingType:  ctxtype,
		matchingValue: matching,
		isRegex:       isregex,
		isCascaded:    iscascaded,
	}
}
