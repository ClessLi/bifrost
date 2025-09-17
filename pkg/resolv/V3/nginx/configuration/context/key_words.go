package context

import (
	"regexp"
	"strings"

	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
)

const (
	RegexpMatchingListenPortValue = `^listen\s*(\d+)\s*\S*$`
	RegexpMatchingServerNameValue = `^server_name\s*(.+)$`
)

var SkipDisabledCtxFilterFunc = func(targetCtx Context) bool { return !targetCtx.IsEnabled() }

type KeyWords interface {
	Match(ctx Context) bool
	SkipQueryThisContext(ctx Context) bool
	Cascaded() bool
	SetCascaded(cascaded bool) KeyWords
	SetStringMatchingValue(value string) KeyWords
	SetRegexpMatchingValue(value string) KeyWords
	SetSkipQueryFilter(filterFunc func(targetCtx Context) bool) KeyWords
	SetMatchingFilter(filterFunc func(targetCtx Context) bool) KeyWords
}

type keywords struct {
	matchingType         context_type.ContextType
	matchingValue        string
	isRegex              bool
	isCascaded           bool
	skipQueryFilterFuncs []func(targetCtx Context) bool
	matchingFilterFuncs  []func(targetCtx Context) bool
}

func (k *keywords) Match(ctx Context) bool {
	for _, filterFunc := range k.matchingFilterFuncs {
		if !filterFunc(ctx) {
			return false
		}
	}

	return true
}

func (k *keywords) SkipQueryThisContext(ctx Context) bool {
	for _, filterFunc := range k.skipQueryFilterFuncs {
		if filterFunc(ctx) {
			return true
		}
	}

	return false
}

func (k *keywords) Cascaded() bool {
	return k.isCascaded
}

func (k *keywords) SetCascaded(cascaded bool) KeyWords {
	k.isCascaded = cascaded

	return k
}

func (k *keywords) SetStringMatchingValue(value string) KeyWords {
	k.isRegex = false
	k.matchingValue = value

	return k
}

func (k *keywords) SetRegexpMatchingValue(value string) KeyWords {
	k.isRegex = true
	k.matchingValue = value

	return k
}

func (k *keywords) SetSkipQueryFilter(filterFunc func(targetCtx Context) bool) KeyWords {
	k.skipQueryFilterFuncs = append(k.skipQueryFilterFuncs, filterFunc)

	return k
}

func (k *keywords) SetMatchingFilter(filterFunc func(targetCtx Context) bool) KeyWords {
	k.matchingFilterFuncs = append(k.matchingFilterFuncs, filterFunc)

	return k
}

func NewKeyWords(ctxtype context_type.ContextType) KeyWords {
	kw := &keywords{
		matchingType:         ctxtype,
		isRegex:              false,
		isCascaded:           true,
		skipQueryFilterFuncs: make([]func(targetCtx Context) bool, 0),
		matchingFilterFuncs:  make([]func(targetCtx Context) bool, 0),
	}
	// match context's type
	kw.matchingFilterFuncs = append(kw.matchingFilterFuncs, func(targetCtx Context) bool {
		return targetCtx.Type() == kw.matchingType
	})
	// match context's value
	kw.matchingFilterFuncs = append(kw.matchingFilterFuncs, func(targetCtx Context) bool {
		matched := false
		if kw.isRegex {
			var err error
			matched, err = regexp.MatchString(kw.matchingValue, targetCtx.Value())
			if err != nil {
				return false
			}
		} else {
			matched = strings.Contains(targetCtx.Value(), kw.matchingValue)
		}

		return matched
	})

	return kw
}
