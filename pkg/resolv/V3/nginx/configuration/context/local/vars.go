package local

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"regexp"
)

var (

	// regexps
	RegCommentHead           = regexp.MustCompile(`^(\s*)#+[ \t\f]*([^\r\n]*?)` + LineBreak + `+`)
	RegDirectiveWithValue    = regexp.MustCompile(S)
	RegDirectiveWithoutValue = regexp.MustCompile(`^\s*(` + Normal + `)\s*;`)
	RegEventsHead            = regexp.MustCompile(`^\s*events\s*{`)
	RegGeoHead               = regexp.MustCompile(`^\s*geo\s*([^;]*?)\s*{`)
	RegHttpHead              = regexp.MustCompile(`^\s*http\s*{`)
	RegIfHead                = regexp.MustCompile(`^\s*if\s*([^;]*?)\s*{`)
	RegLimitExceptHead       = regexp.MustCompile(`^\s*limit_except\s*([^;]*?)\s*{`)
	RegLocationHead          = regexp.MustCompile(`^\s*location\s*([^;]*?)\s*{`)
	RegMapHead               = regexp.MustCompile(`^\s*map\s*([^;]*?)\s*{`)
	RegServerHead            = regexp.MustCompile(`^\s*server\s*{`)
	RegStreamHead            = regexp.MustCompile(`^\s*stream\s*{`)
	RegTypesHead             = regexp.MustCompile(`^\s*types\s*{`)
	RegUpstreamHead          = regexp.MustCompile(`^\s*upstream\s*([^;]*?)\s*{`)
	RegBlankLine             = regexp.MustCompile(`^` + LineBreak + `\s*` + LineBreak + `$`)
	RegBraceEnd              = regexp.MustCompile(`^\s*}`)
	RegErrorHeed             = regexp.MustCompile(Abnormal)
	RegLineBreak             = regexp.MustCompile(LineBreak)

	// parse config function maps
	pushStackParseFuncMap = make(map[context_type.ContextType]parseFunc)
	inStackParseFuncMap   = make(map[context_type.ContextType]parseFunc)

	// commentsToContexts convertor parsing function map
	convertorPushStackParseFuncMap = make(map[context_type.ContextType]func(comment *Comment) bool)

	// context builder map
	builderMap = make(map[context_type.ContextType]func(value string) context.Context)
)
