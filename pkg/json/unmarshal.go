// json 包，该包包含了用于resolve包中nginx配置对象json反序列化的相关对象和方法及函数
// 创建者: ClessLi
// 创建时间: 2020-4-13 09:37:01
package json

import (
	"encoding/json"
	"fmt"
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"regexp"
)

// unmarshaler, json反序列化内部接口对象，定义了用于nginx配置对象json反序列化所需实现的方法
type unmarshaler interface {
	UnmarshalToJSON(b []byte, configCaches *[]string) (resolv.Parser, error) // json反序列化方法
	getChildren() []*json.RawMessage                                         // 统一输出对象子对象json串内部方法
	toParser(children []resolv.Parser) (resolv.Parser, error)                // 统一解析并返回nginx配置对象的内部方法
}

// BasicContext, 用于json反序列化的上下文基础对象，定义了上下文类型的基本属性及基础方法
type BasicContext struct {
	Name     string             `json:"-"`
	Value    string             `json:"value,omitempty"`
	Children []*json.RawMessage `json:"param,omitempty"`
}

func (bc *BasicContext) getChildren() []*json.RawMessage {
	return bc.Children
}

type Config struct {
	BasicContext `json:"config"`
}

func (conf *Config) UnmarshalToJSON(b []byte, configCaches *[]string) (resolv.Parser, error) {
	return unmarshal(b, conf, configCaches)
	//return unmarshal(b)
}

func (conf *Config) toParser(children []resolv.Parser) (resolv.Parser, error) {
	return resolv.NewConf(children, conf.Value)
}

type Include struct {
	BasicContext `json:"include"`
	Key          *Key     `json:"tags"`
	Comment      *Comment `json:"comments"`
	ConfPWD      string   `json:"conf_pwd"`
}

func (i *Include) UnmarshalToJSON(b []byte, configCaches *[]string) (resolv.Parser, error) {
	return unmarshal(b, i, configCaches)
	//return unmarshal(b)
}

func (i *Include) toParser(children []resolv.Parser) (resolv.Parser, error) {
	IN, err := resolv.NewInclude(i.ConfPWD, i.Value, &[]string{})
	if err != nil {
		return nil, err
	}
	IN.Children = children
	return IN, nil
}

type Types struct {
	BasicContext `json:"types"`
}

func (t *Types) UnmarshalToJSON(b []byte, configCaches *[]string) (resolv.Parser, error) {
	return unmarshal(b, t, configCaches)
	//return unmarshal(b)
}

func (t *Types) toParser(children []resolv.Parser) (resolv.Parser, error) {
	T := resolv.NewTypes()
	T.Children = children
	return T, nil
}

type Map struct {
	BasicContext `json:"map"`
}

func (m *Map) UnmarshalToJSON(b []byte, configCaches *[]string) (resolv.Parser, error) {
	return unmarshal(b, m, configCaches)
	//return unmarshal(b)
}

func (m *Map) toParser(children []resolv.Parser) (resolv.Parser, error) {
	M := resolv.NewMap(m.Value)
	M.Children = children
	return M, nil
}

type LimitExcept struct {
	BasicContext `json:"limit_except"`
}

func (le *LimitExcept) UnmarshalToJSON(b []byte, configCaches *[]string) (resolv.Parser, error) {
	return unmarshal(b, le, configCaches)
	//return unmarshal(b)
}

func (le *LimitExcept) toParser(children []resolv.Parser) (resolv.Parser, error) {
	Le := resolv.NewLimitExcept(le.Value)
	Le.Children = children
	return Le, nil
}

type Events struct {
	BasicContext `json:"events"`
}

func (e *Events) UnmarshalToJSON(b []byte, configCaches *[]string) (resolv.Parser, error) {
	return unmarshal(b, e, configCaches)
	//return unmarshal(b)
}

func (e *Events) toParser(children []resolv.Parser) (resolv.Parser, error) {
	E := resolv.NewEvents()
	E.Children = children
	return E, nil
}

type Geo struct {
	BasicContext `json:"geo"`
}

func (g *Geo) UnmarshalToJSON(b []byte, configCaches *[]string) (resolv.Parser, error) {
	return unmarshal(b, g, configCaches)
	//return unmarshal(b)
}

func (g *Geo) toParser(children []resolv.Parser) (resolv.Parser, error) {
	G := resolv.NewGeo(g.Value)
	G.Children = children
	return G, nil
}

type Http struct {
	BasicContext `json:"http"`
}

func (h *Http) UnmarshalToJSON(b []byte, configCaches *[]string) (resolv.Parser, error) {
	return unmarshal(b, h, configCaches)
	//return unmarshal(b)
}

func (h *Http) toParser(children []resolv.Parser) (resolv.Parser, error) {
	H := resolv.NewHttp()
	H.Children = children
	return H, nil
}

type Stream struct {
	BasicContext `json:"stream"`
}

func (st *Stream) UnmarshalToJSON(b []byte, configCaches *[]string) (resolv.Parser, error) {
	return unmarshal(b, st, configCaches)
	//return unmarshal(b)
}

func (st *Stream) toParser(children []resolv.Parser) (resolv.Parser, error) {
	St := resolv.NewStream()
	St.Children = children
	return St, nil
}

type Upstream struct {
	BasicContext `json:"upstream"`
}

func (u *Upstream) UnmarshalToJSON(b []byte, configCaches *[]string) (resolv.Parser, error) {
	return unmarshal(b, u, configCaches)
	//return unmarshal(b)
}

func (u *Upstream) toParser(children []resolv.Parser) (resolv.Parser, error) {
	U := resolv.NewUpstream(u.Value)
	U.Children = children
	return U, nil
}

type Server struct {
	BasicContext `json:"server"`
}

func (s *Server) UnmarshalToJSON(b []byte, configCaches *[]string) (resolv.Parser, error) {
	return unmarshal(b, s, configCaches)
	//return unmarshal(b)
}

func (s *Server) toParser(children []resolv.Parser) (resolv.Parser, error) {
	S := resolv.NewServer()
	S.Children = children
	return S, nil
}

type Location struct {
	BasicContext `json:"location"`
}

func (l *Location) UnmarshalToJSON(b []byte, configCaches *[]string) (resolv.Parser, error) {
	return unmarshal(b, l, configCaches)
	//return unmarshal(b)
}

func (l *Location) toParser(children []resolv.Parser) (resolv.Parser, error) {
	L := resolv.NewLocation(l.Value)
	L.Children = children
	return L, nil
}

type If struct {
	BasicContext `json:"if"`
}

func (i *If) UnmarshalToJSON(b []byte, configCaches *[]string) (resolv.Parser, error) {
	return unmarshal(b, i, configCaches)
	//return unmarshal(b)
}

func (i *If) toParser(children []resolv.Parser) (resolv.Parser, error) {
	I := resolv.NewIf(i.Value)
	I.Children = children
	return I, nil
}

type Key struct {
	resolv.Key
}

func (k *Key) UnmarshalToJSON(b []byte, configCaches *[]string) (resolv.Parser, error) {
	err := json.Unmarshal(b, k)
	return &k.Key, err
}

func (k *Key) getChildren() []*json.RawMessage {
	return nil
}

func (k *Key) toParser(_ []resolv.Parser) (resolv.Parser, error) {
	return &k.Key, nil
}

type Comment struct {
	resolv.Comment
}

func (c *Comment) UnmarshalToJSON(b []byte, configCaches *[]string) (resolv.Parser, error) {
	err := json.Unmarshal(b, c)
	return &c.Comment, err
}

func (c *Comment) getChildren() []*json.RawMessage {
	return nil
}

func (c *Comment) toParser(_ []resolv.Parser) (resolv.Parser, error) {
	return &c.Comment, nil
}

// Unmarshal, json反序列化并返回nginx配置对象的函数
func Unmarshal(b []byte) (resolv.Parser, error) {
	return unmarshal(b, &Config{}, &[]string{})
}

// unmarshal, json反序列化并返回nginx配置对象的内部函数
func unmarshal(b []byte, p unmarshaler, configCaches *[]string) (resolv.Parser, error) {
	switch p.(type) {
	case *Key, *Comment:
		return p.UnmarshalToJSON(b, configCaches)
	default:
		err := json.Unmarshal(b, p)
		if err != nil {
			return nil, err
		}
		if conf, ok := p.(*Config); ok {
			if inCaches(conf.Value, configCaches) {
				return nil, fmt.Errorf("config '%s' is already loaded", conf.Value)
			}
			newCaches := *configCaches
			newCaches = append(newCaches, conf.Value)
			children, cerr := unmarshalChildren(p.getChildren(), &newCaches)
			if cerr != nil {
				return nil, cerr
			}
			return p.toParser(children)
		} else {
			children, cerr := unmarshalChildren(p.getChildren(), configCaches)
			if cerr != nil {
				return nil, cerr
			}
			return p.toParser(children)
		}
	}
}

func inCaches(path string, caches *[]string) bool {
	for _, cache := range *caches {
		if path == cache {
			return true
		}
	}
	return false
}

// unmarshalChildren, 解析并反序列化子json串切片对象的内部函数
func unmarshalChildren(bytes []*json.RawMessage, configCaches *[]string) (children []resolv.Parser, err error) {
	// parseContext, 用于解析json串归属于哪类需反序列化对象的匿名函数
	parseContext := func(b []byte, reg *regexp.Regexp) bool {
		if m := reg.Find(b); m != nil {
			return true
		} else {
			return false
		}
	}

	for _, b := range bytes {
		switch {
		case parseContext(*b, RegCommentHead):
			comment, err := unmarshal(*b, &Comment{}, configCaches)
			if err != nil {
				return nil, err
			}
			children = append(children, comment)
		case parseContext(*b, RegIncludeHead):
			include, err := unmarshal(*b, &Include{}, configCaches)
			if err != nil {
				return nil, err
			}
			children = append(children, include)
		case parseContext(*b, RegConfigHead):
			config, err := unmarshal(*b, &Config{}, configCaches)
			if err != nil {
				return nil, err
			}
			children = append(children, config)
		case parseContext(*b, RegEventsHead):
			events, err := unmarshal(*b, &Events{}, configCaches)
			if err != nil {
				return nil, err
			}
			children = append(children, events)
		case parseContext(*b, RegGeoHead):
			geo, err := unmarshal(*b, &Geo{}, configCaches)
			if err != nil {
				return nil, err
			}
			children = append(children, geo)
		case parseContext(*b, RegHttpHead):
			http, err := unmarshal(*b, &Http{}, configCaches)
			if err != nil {
				return nil, err
				//return
			}
			children = append(children, http)
		case parseContext(*b, RegIfHead):
			i, err := unmarshal(*b, &If{}, configCaches)
			if err != nil {
				return nil, err
				//return
			}
			children = append(children, i)
		case parseContext(*b, RegLimitExceptHead):
			le, err := unmarshal(*b, &LimitExcept{}, configCaches)
			if err != nil {
				return nil, err
				//return
			}
			children = append(children, le)
		case parseContext(*b, RegLocationHead):
			l, err := unmarshal(*b, &Location{}, configCaches)
			if err != nil {
				return nil, err
				//return
			}
			children = append(children, l)
		case parseContext(*b, RegMapHead):
			m, err := unmarshal(*b, &Map{}, configCaches)
			if err != nil {
				return nil, err
				//return
			}
			children = append(children, m)
		case parseContext(*b, RegServerHead):
			svr, err := unmarshal(*b, &Server{}, configCaches)
			if err != nil {
				return nil, err
				//return
			}
			children = append(children, svr)
		case parseContext(*b, RegStreamHead):
			st, err := unmarshal(*b, &Stream{}, configCaches)
			if err != nil {
				return nil, err
				//return
			}
			children = append(children, st)
		case parseContext(*b, RegTypesHead):
			t, err := unmarshal(*b, &Types{}, configCaches)
			if err != nil {
				return nil, err
				//return
			}
			children = append(children, t)
		case parseContext(*b, RegUpstreamHead):
			u, err := unmarshal(*b, &Upstream{}, configCaches)
			if err != nil {
				return nil, err
				//return
			}
			children = append(children, u)
		default:
			k, err := unmarshal(*b, &Key{}, configCaches)
			if err != nil {
				return nil, err
				//return
			}
			children = append(children, k)
		}
	}

	return children, nil

}
