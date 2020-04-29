package filter

import "github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"

var (
	KeywordHTTP      = resolv.NewKeyWords("http", "", "", false)
	KeywordStream    = resolv.NewKeyWords("stream", "", "", false)
	KeywordSvrName   = resolv.NewKeyWords("key", `^server_name$`, `.*`, true)
	KeywordPorts     = resolv.NewKeyWords("key", `^listen$`, `.*`, true)
	keywordLocations = resolv.NewKeyWords("location", "", `.*`, true)
)
