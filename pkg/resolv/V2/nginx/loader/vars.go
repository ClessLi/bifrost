package loader

import (
	"regexp"
)

var (

	// regexp
	// RegEndWithCR       = regexp.MustCompile("}\n+$").
	RegEventsHead      = regexp.MustCompile(`^\s*events\s*{`)
	RegHttpHead        = regexp.MustCompile(`^\s*http\s*{`)
	RegStreamHead      = regexp.MustCompile(`^\s*stream\s*{`)
	RegServerHead      = regexp.MustCompile(`^\s*server\s*{`)
	RegLocationHead    = regexp.MustCompile(`^\s*location\s*([^;]*?)\s*{`)
	RegIfHead          = regexp.MustCompile(`^\s*if\s*([^;]*?)\s*{`)
	RegUpstreamHead    = regexp.MustCompile(`^\s*upstream\s*([^;]*?)\s*{`)
	RegGeoHead         = regexp.MustCompile(`^\s*geo\s*([^;]*?)\s*{`)
	RegMapHead         = regexp.MustCompile(`^\s*map\s*([^;]*?)\s*{`)
	RegLimitExceptHead = regexp.MustCompile(`^\s*limit_except\s*([^;]*?)\s*{`)
	RegTypesHead       = regexp.MustCompile(`^\s*types\s*{`)
	RegContextEnd      = regexp.MustCompile(`^\s*}`)
	RegCommentHead     = regexp.MustCompile(`^(\s*)#+[ \r\t\f]*(.*?)\n`)
	RegKeyValue        = regexp.MustCompile(S)
	RegKey             = regexp.MustCompile(`^\s*(` + Normal + `)\s*;`)
	RegBlankLine       = regexp.MustCompile(`^\n\s*` + LineBreak + `$`)
	RegErrorHeed       = regexp.MustCompile(Abnormal)
	// RegLine            = regexp.MustCompile(LineBreak).

	// json unmarshal.

	JsonUnmarshalRegEventsHead      = regexp.MustCompile(`^\s*{\s*"events"\s*:\s*{`)
	JsonUnmarshalRegHttpHead        = regexp.MustCompile(`^\s*{\s*"http"\s*:\s*{`)
	JsonUnmarshalRegStreamHead      = regexp.MustCompile(`^\s*{\s*"stream"\s*:\s*{`)
	JsonUnmarshalRegServerHead      = regexp.MustCompile(`^\s*{\s*"server"\s*:\s*{`)
	JsonUnmarshalRegLocationHead    = regexp.MustCompile(`^\s*{\s*"location"\s*:\s*{`)
	JsonUnmarshalRegIfHead          = regexp.MustCompile(`^\s*{\s*"if"\s*:\s*{`)
	JsonUnmarshalRegUpstreamHead    = regexp.MustCompile(`^\s*{\s*"upstream"\s*:\s*{`)
	JsonUnmarshalRegGeoHead         = regexp.MustCompile(`^\s*{\s*"geo"\s*:\s*{`)
	JsonUnmarshalRegMapHead         = regexp.MustCompile(`^\s*{\s*"map"\s*:\s*{`)
	JsonUnmarshalRegLimitExceptHead = regexp.MustCompile(`^\s*{\s*"limit_except"\s*:\s*{`)
	JsonUnmarshalRegTypesHead       = regexp.MustCompile(`^\s*{\s*"types"\s*:\s*{`)
	JsonUnmarshalRegIncludeHead     = regexp.MustCompile(`^\s*{\s*"include"\s*:\s*{`)
	JsonUnmarshalRegConfigHead      = regexp.MustCompile(`^\s*{\s*"config"\s*:\s*{`)
	JsonUnmarshalRegCommentHead     = regexp.MustCompile(`^\s*{\s*"comments"\s*:\s*"`)
	// KeywordHTTP    = NewKeyWords(TypeHttp, "", "", false, true)
	// KeywordStream  = NewKeyWords(TypeStream, "", "", false, true)
	// KeywordSvrName = NewKeyWords(TypeKey, `server_name`, `*`, false, true)
	// KeywordPort    = NewKeyWords(TypeKey, `^listen$`, `.*`, true, true)
	// KeywordLocations = NewKeyWords(TypeLocation, "", `.*`, true, true).
)
