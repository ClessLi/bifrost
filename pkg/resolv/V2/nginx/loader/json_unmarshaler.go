package loader

import (
	"encoding/json"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration/parser"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/loop_preventer"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_indention"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_position"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_type"
	"regexp"
)

type UnmarshalContext interface {
	GetValue() string
	GetChildren() []*json.RawMessage
}

type unmarshalContext struct {
	//Name     string             `json:"-"`
	Value    string             `json:"value,omitempty"`
	Children []*json.RawMessage `json:"param,omitempty"`
}

func (u unmarshalContext) GetValue() string {
	return u.Value
}

func (u unmarshalContext) GetChildren() []*json.RawMessage {
	return u.Children
}

type config struct {
	unmarshalContext `json:"config"`
}

type events struct {
	unmarshalContext `json:"events"`
}

type geo struct {
	unmarshalContext `json:"geo"`
}

type http struct {
	unmarshalContext `json:"http"`
}

type _if struct {
	unmarshalContext `json:"if"`
}

type include struct {
	unmarshalContext `json:"include"`
}

type limitExcept struct {
	unmarshalContext `json:"limit_except"`
}

type location struct {
	unmarshalContext `json:"location"`
}

type _map struct {
	unmarshalContext `json:"map"`
}

type server struct {
	unmarshalContext `json:"server"`
}

type stream struct {
	unmarshalContext `json:"stream"`
}

type types struct {
	unmarshalContext `json:"types"`
}

type upstream struct {
	unmarshalContext `json:"upstream"`
}

type unmarshaler struct {
	contextType      parser_type.ParserType
	position         parser_position.ParserPosition
	indention        parser_indention.Indention
	context          parser.Context
	unmarshalContext UnmarshalContext
	LoadCacher
	loop_preventer.LoopPreventer
}

func (u *unmarshaler) UnmarshalJSON(bytes []byte) error {
	//fmt.Println(string(bytes))
	err := json.Unmarshal(bytes, u.unmarshalContext)
	if err != nil {
		return err
	}
	if u.contextType == parser_type.TypeConfig {
		if u.LoopPreventer == nil {
			u.LoopPreventer = loop_preventer.NewLoopPreverter(u.unmarshalContext.GetValue())
		} else {
			err = u.LoopPreventer.CheckLoopPrevent(u.position.Id(), u.unmarshalContext.GetValue())
			if err != nil {
				return err
			}
		}
		u.position = parser_position.NewPosition(u.unmarshalContext.GetValue())
		// 读取缓存
		if u.LoadCacher == nil {
			u.LoadCacher = NewLoadCacher(u.unmarshalContext.GetValue())
		}
		if config := u.LoadCacher.GetConfig(u.unmarshalContext.GetValue()); config != nil {
			u.context = config
			return nil
		}
		u.context = parser.NewContext(u.unmarshalContext.GetValue(), u.contextType, u.indention)
		err = u.LoadCacher.SetConfig(u.context.(*parser.Config))
		if err != nil {
			return err
		}
	} else {
		// 根据context类型创建反序列器context对象
		u.context = parser.NewContext(u.unmarshalContext.GetValue(), u.contextType, u.indention)

	}

	// parseContext, 用于解析json串归属于哪类需反序列化对象的匿名函数
	parseContext := func(b []byte, reg *regexp.Regexp) bool {
		if m := reg.Find(b); m != nil {
			return true
		} else {
			return false
		}
	}

	for _, child := range u.unmarshalContext.GetChildren() {
		var parserType parser_type.ParserType
		var p parser.Parser
		var unmarshalCtx UnmarshalContext
		indention := u.indention.NextIndention()
		switch {
		case parseContext(*child, JsonUnmarshalRegCommentHead):
			comment := parser.NewComment("", false, indention)
			err = json.Unmarshal(*child, comment)
			if err != nil {
				return err
			}
			p = comment
		case parseContext(*child, JsonUnmarshalRegIncludeHead):
			parserType = parser_type.TypeInclude
			unmarshalCtx = new(include)
		case parseContext(*child, JsonUnmarshalRegConfigHead):
			parserType = parser_type.TypeConfig
			unmarshalCtx = new(config)
		case parseContext(*child, JsonUnmarshalRegEventsHead):
			parserType = parser_type.TypeEvents
			unmarshalCtx = new(events)
		case parseContext(*child, JsonUnmarshalRegGeoHead):
			parserType = parser_type.TypeGeo
			unmarshalCtx = new(geo)
		case parseContext(*child, JsonUnmarshalRegHttpHead):
			parserType = parser_type.TypeHttp
			unmarshalCtx = new(http)
		case parseContext(*child, JsonUnmarshalRegIfHead):
			parserType = parser_type.TypeIf
			unmarshalCtx = new(_if)
		case parseContext(*child, JsonUnmarshalRegLimitExceptHead):
			parserType = parser_type.TypeLimitExcept
			unmarshalCtx = new(limitExcept)
		case parseContext(*child, JsonUnmarshalRegLocationHead):
			parserType = parser_type.TypeLocation
			unmarshalCtx = new(location)
		case parseContext(*child, JsonUnmarshalRegMapHead):
			parserType = parser_type.TypeMap
			unmarshalCtx = new(_map)
		case parseContext(*child, JsonUnmarshalRegServerHead):
			parserType = parser_type.TypeServer
			unmarshalCtx = new(server)
		case parseContext(*child, JsonUnmarshalRegStreamHead):
			parserType = parser_type.TypeStream
			unmarshalCtx = new(stream)
		case parseContext(*child, JsonUnmarshalRegTypesHead):
			parserType = parser_type.TypeTypes
			unmarshalCtx = new(types)
		case parseContext(*child, JsonUnmarshalRegUpstreamHead):
			parserType = parser_type.TypeUpstream
			unmarshalCtx = new(upstream)
		default:
			key := parser.NewKey("", "", indention)
			err = json.Unmarshal(*child, key)
			if err != nil {
				return err
			}
			p = key
		}
		if p != nil {
			err = u.context.Insert(p, u.context.Len())
			if err != nil {
				return err
			}
			continue
		}

		next := &unmarshaler{
			contextType:      parserType,
			position:         u.position,
			indention:        indention,
			context:          nil,
			unmarshalContext: unmarshalCtx,
			LoopPreventer:    u.LoopPreventer,
		}
		err = next.UnmarshalJSON(*child)
		if err != nil {
			return err
		}
		if next.context != nil {
			err = u.context.Insert(next.context, u.context.Len())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewUnmarshaler() *unmarshaler {
	return &unmarshaler{
		contextType:      parser_type.TypeConfig,
		unmarshalContext: new(config),
		indention:        parser_indention.NewIndention(),
	}
}
