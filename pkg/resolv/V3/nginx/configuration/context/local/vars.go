package local

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"regexp"
)

var (

	// regexps
	RegCommentHead           = regexp.MustCompile(`^(\s*)#+[ \r\t\f]*(.*?)\n`)
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
	RegBlankLine             = regexp.MustCompile(`^\n\s*` + LineBreak + `$`)
	RegBraceEnd              = regexp.MustCompile(`^\s*}`)
	RegErrorHeed             = regexp.MustCompile(Abnormal)

	// json unmarshal regexps
	JsonUnmarshalRegCommentHead     = regexp.MustCompile(`^\s*{[^{]*"comments"\s*:\s*"`)
	JsonUnmarshalRegEventsHead      = regexp.MustCompile(`^\s*{\s*"events"\s*:\s*{`)
	JsonUnmarshalRegGeoHead         = regexp.MustCompile(`^\s*{\s*"geo"\s*:\s*{`)
	JsonUnmarshalRegHttpHead        = regexp.MustCompile(`^\s*{\s*"http"\s*:\s*{`)
	JsonUnmarshalRegIfHead          = regexp.MustCompile(`^\s*{\s*"if"\s*:\s*{`)
	JsonUnmarshalRegIncludeHead     = regexp.MustCompile(`^\s*{\s*"include"\s*:\s*{`)
	JsonUnmarshalRegLimitExceptHead = regexp.MustCompile(`^\s*{\s*"limit_except"\s*:\s*{`)
	JsonUnmarshalRegLocationHead    = regexp.MustCompile(`^\s*{\s*"location"\s*:\s*{`)
	JsonUnmarshalRegMapHead         = regexp.MustCompile(`^\s*{\s*"map"\s*:\s*{`)
	JsonUnmarshalRegServerHead      = regexp.MustCompile(`^\s*{\s*"server"\s*:\s*{`)
	JsonUnmarshalRegStreamHead      = regexp.MustCompile(`^\s*{\s*"stream"\s*:\s*{`)
	JsonUnmarshalRegTypesHead       = regexp.MustCompile(`^\s*{\s*"types"\s*:\s*{`)
	JsonUnmarshalRegUpstreamHead    = regexp.MustCompile(`^\s*{\s*"upstream"\s*:\s*{`)

	// parse config function maps
	pushStackParseFuncMap = make(map[context_type.ContextType]parseFunc)
	inStackParseFuncMap   = make(map[context_type.ContextType]parseFunc)

	// context builder map
	builderMap = make(map[context_type.ContextType]func(value string) context.Context)

	// json unmarshal regexp matcher map
	jsonUnmarshalRegMatcherMap = make(map[context_type.ContextType]func(jsonraw []byte) bool)

	// json unmarshaler builder map
	jsonUnmarshalerBuilderMap = make(map[context_type.ContextType]func(graph ConfigGraph, father context.Context) *jsonUnmarshaler)
)
