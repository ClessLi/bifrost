package resolv

const (
	DoubleQuotes    = `\s*"[^"]*"`
	SingleQuotes    = `\s*'[^']*'`
	Normal          = `\s*[^;\s]*`
	S1              = DoubleQuotes + `|` + SingleQuotes + `|` + Normal
	S               = `^\s*(` + S1 + `)\s*((?:` + S1 + `)+);`
	TypeConfig      = "config"
	TypeEvents      = "events"
	TypeGeo         = "geo"
	TypeHttp        = "http"
	TypeIf          = "if"
	TypeInclude     = "include"
	TypeKey         = "key"
	TypeLimitExcept = "limit_except"
	TypeLocation    = "location"
	TypeMap         = "map"
	TypeServer      = "server"
	TypeStream      = "stream"
	TypeTypes       = "types"
	TypeUpstream    = "upstream"
	TypeComment     = "comment"
)

// 整型order
const (
	ServerPort Order = iota
)

// 字符串型order
const (
	ServerName Order = 1000 + iota
)
