package parser_type

type ParserType string

func (t ParserType) String() string {
	return string(t)
}

const (
	TypeConfig      ParserType = "config"
	TypeEvents      ParserType = "events"
	TypeGeo         ParserType = "geo"
	TypeHttp        ParserType = "http"
	TypeIf          ParserType = "if"
	TypeInclude     ParserType = "include"
	TypeKey         ParserType = "key"
	TypeLimitExcept ParserType = "limit_except"
	TypeLocation    ParserType = "location"
	TypeMap         ParserType = "map"
	TypeServer      ParserType = "server"
	TypeStream      ParserType = "stream"
	TypeTypes       ParserType = "types"
	TypeUpstream    ParserType = "upstream"
	TypeComment     ParserType = "comment"
)
