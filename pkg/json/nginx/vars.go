package nginx

import "regexp"

var (
	RegEventsHead      = regexp.MustCompile(`^\s*{\s*"events"\s*:\s*{`)
	RegHttpHead        = regexp.MustCompile(`^\s*{\s*"http"\s*:\s*{`)
	RegStreamHead      = regexp.MustCompile(`^\s*{\s*"stream"\s*:\s*{`)
	RegServerHead      = regexp.MustCompile(`^\s*{\s*"server"\s*:\s*{`)
	RegLocationHead    = regexp.MustCompile(`^\s*{\s*"location"\s*:\s*{`)
	RegIfHead          = regexp.MustCompile(`^\s*{\s*"if"\s*:\s*{`)
	RegUpstreamHead    = regexp.MustCompile(`^\s*{\s*"upstream"\s*:\s*{`)
	RegGeoHead         = regexp.MustCompile(`^\s*{\s*"geo"\s*:\s*{`)
	RegMapHead         = regexp.MustCompile(`^\s*{\s*"map"\s*:\s*{`)
	RegLimitExceptHead = regexp.MustCompile(`^\s*{\s*"limit_except"\s*:\s*{`)
	RegTypesHead       = regexp.MustCompile(`^\s*{\s*"types"\s*:\s*{`)
	RegIncludeHead     = regexp.MustCompile(`^\s*{\s*"include"\s*:\s*{`)
	RegConfigHead      = regexp.MustCompile(`^\s*{\s*"config"\s*:\s*{`)
	RegCommentHead     = regexp.MustCompile(`^\s*{\s*"comments"\s*:\s*"`)
	// RegKeyValue        = regexp.MustCompile(S)
	// RegKey             = regexp.MustCompile(`^\s*(\S+);`).
)
