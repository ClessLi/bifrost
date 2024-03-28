package context

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"regexp"
	"strings"
)

const (
	RegexpMatchingListenPortValue = `^listen\s*(\d+)\s*\S*$`
	RegexpMatchingServerNameValue = `^server_name\s*(.+)$`
)

type KeyWords interface {
	Match(ctx Context) bool
	Cascaded() bool
	SetCascaded(cascaded bool) KeyWords
	SetStringMatchingValue(value string) KeyWords
	SetRegexpMatchingValue(value string) KeyWords
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

func NewKeyWords(ctxtype context_type.ContextType) KeyWords {
	return &keywords{
		matchingType: ctxtype,
		isRegex:      false,
		isCascaded:   true,
	}
}
