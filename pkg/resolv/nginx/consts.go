package nginx

import "fmt"

type parserType string

func (pt *parserType) ToString() string {
	return fmt.Sprintf("%s", *pt)
}

const (
	DoubleQuotes               = `\s*"[^"]*"`
	SingleQuotes               = `\s*'[^']*'`
	Normal                     = `\s*[^;\s]*`
	S1                         = DoubleQuotes + `|` + SingleQuotes + `|` + Normal
	S                          = `^\s*(` + S1 + `)\s*((?:` + S1 + `)+);`
	TypeConfig      parserType = "config"
	TypeEvents      parserType = "events"
	TypeGeo         parserType = "geo"
	TypeHttp        parserType = "http"
	TypeIf          parserType = "if"
	TypeInclude     parserType = "include"
	TypeKey         parserType = "key"
	TypeLimitExcept parserType = "limit_except"
	TypeLocation    parserType = "location"
	TypeMap         parserType = "map"
	TypeServer      parserType = "server"
	TypeStream      parserType = "stream"
	TypeTypes       parserType = "types"
	TypeUpstream    parserType = "upstream"
	TypeComment     parserType = "comment"

	hashForGetList  = "ForList"
	hashForString   = "ForString"
	hashForDumpTemp = "ForDumpTemp"
)

// 整型order.
const (
	ServerPort Order = iota
)

// 字符串型order.
const (
	ServerName Order = 1000 + iota
)
