package context_type

type ContextType string

func (t ContextType) String() string {
	return string(t)
}

const (
	TypeMain          ContextType = "main"
	TypeConfig        ContextType = "config"
	TypeEvents        ContextType = "events"
	TypeGeo           ContextType = "geo"
	TypeHttp          ContextType = "http"
	TypeIf            ContextType = "if"
	TypeInclude       ContextType = "include"
	TypeDirective     ContextType = "directive"
	TypeLimitExcept   ContextType = "limit_except"
	TypeLocation      ContextType = "location"
	TypeMap           ContextType = "map"
	TypeServer        ContextType = "server"
	TypeStream        ContextType = "stream"
	TypeTypes         ContextType = "types"
	TypeUpstream      ContextType = "upstream"
	TypeComment       ContextType = "comment"
	TypeInlineComment ContextType = "inline_comment"
	TypeErrContext    ContextType = "error_context"
)
