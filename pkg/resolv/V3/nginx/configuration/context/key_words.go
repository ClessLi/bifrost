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
	IsToLeafQuery() bool
	SetIsToLeafQuery(isToLeafQuery bool) KeyWords
	SetStringMatchingValue(value string) KeyWords
	SetRegexpMatchingValue(value string) KeyWords
	SetSkipQueryFilter(filterFunc func(targetCtx Context) bool) KeyWords
	AppendMatchingFilter(filterFunc func(targetCtx Context) bool) KeyWords
}

type keywords struct {
	matchingValue        string
	isRegex              bool
	isCascaded           bool
	isToLeafQuery        bool
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

func (k *keywords) IsToLeafQuery() bool {
	return k.isToLeafQuery
}

func (k *keywords) SetIsToLeafQuery(isToLeafQuery bool) KeyWords {
	k.isToLeafQuery = isToLeafQuery

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

func (k *keywords) AppendMatchingFilter(filterFunc func(targetCtx Context) bool) KeyWords {
	k.matchingFilterFuncs = append(k.matchingFilterFuncs, filterFunc)

	return k
}

func NewKeyWords(matchingFilterFunc func(targetCtx Context) bool) KeyWords {
	kw := &keywords{
		isRegex:              false,
		isCascaded:           true,
		isToLeafQuery:        true,
		skipQueryFilterFuncs: make([]func(targetCtx Context) bool, 0),
		matchingFilterFuncs:  make([]func(targetCtx Context) bool, 0),
	}

	// set the first matching filter func
	return kw.AppendMatchingFilter(matchingFilterFunc).
		// set next matching filter func to match context's type
		AppendMatchingFilter(func(targetCtx Context) bool {
			// skip empty matching value
			if kw.matchingValue == "" {
				return true
			}

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
}

func NewKeyWordsByType(ctxtype context_type.ContextType) KeyWords {
	return NewKeyWords(func(targetCtx Context) bool {
		return targetCtx.Type() == ctxtype
	})
}
