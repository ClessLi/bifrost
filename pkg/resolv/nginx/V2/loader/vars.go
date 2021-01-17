package loader

import (
	"regexp"
	"time"
)

var (
	//system vars
	TZ = time.Local

	// regexp
	RegEndWithCR       = regexp.MustCompile("}\n+$")
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
	RegKey             = regexp.MustCompile(`^\s*(\S+);`)
	RegPortValue       = regexp.MustCompile(`^(\d+)\s*\S*$`)
	RegBlankLine       = regexp.MustCompile(`^\s*` + LineBreak)
	RegErrorHeed       = regexp.MustCompile(Abnormal)
	RegLine            = regexp.MustCompile(LineBreak)

	//KeywordHTTP    = NewKeyWords(TypeHttp, "", "", false, true)
	//KeywordStream  = NewKeyWords(TypeStream, "", "", false, true)
	//KeywordSvrName = NewKeyWords(TypeKey, `server_name`, `*`, false, true)
	//KeywordPort    = NewKeyWords(TypeKey, `^listen$`, `.*`, true, true)
	//KeywordLocations = NewKeyWords(TypeLocation, "", `.*`, true, true)
)
