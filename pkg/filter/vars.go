package filter

import "github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"

var (
	KeywordHTTP      = resolv.NewKeyWords(resolv.TypeHttp, "", "", false)
	KeywordStream    = resolv.NewKeyWords(resolv.TypeStream, "", "", false)
	KeywordSvrName   = resolv.NewKeyWords(resolv.TypeKey, `^server_name$`, `.*`, true)
	KeywordPort      = resolv.NewKeyWords(resolv.TypeKey, `^listen$`, `.*`, true)
	keywordLocations = resolv.NewKeyWords(resolv.TypeLocation, "", `.*`, true)
)
