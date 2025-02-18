// json 包，该包包含了用于resolve包中nginx配置对象json反序列化的相关对象和方法及函数
// 创建者: ClessLi
// 创建时间: 2020-4-13 09:37:01
package nginx

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
)

// unmarshaler, json反序列化内部接口对象，定义了用于nginx配置对象json反序列化所需实现的方法.
type unmarshaler interface {
	UnmarshalToJSON(b []byte, caches *nginx.Caches) (nginx.Parser, error) // json反序列化方法
	getChildren() []*json.RawMessage                                      // 统一输出对象子对象json串内部方法
	toParser(children []nginx.Parser) (nginx.Parser, error)               // 统一解析并返回nginx配置对象的内部方法
}

// BasicContext, 用于json反序列化的上下文基础对象，定义了上下文类型的基本属性及基础方法.
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

func (conf *Config) UnmarshalToJSON(b []byte, caches *nginx.Caches) (nginx.Parser, error) {
	return unmarshal(b, conf, caches)
	// return unmarshal(b)
}

func (conf *Config) toParser(children []nginx.Parser) (nginx.Parser, error) {
	return nginx.NewConf(children, conf.Value), nil
}

type Include struct {
	BasicContext `json:"include"`
	Key          *Key     `json:"tags"`
	Comment      *Comment `json:"comments"`
	ConfPWD      string   `json:"conf_pwd"`
}

func (i *Include) UnmarshalToJSON(b []byte, caches *nginx.Caches) (nginx.Parser, error) {
	return unmarshal(b, i, caches)
	// return unmarshal(b)
}

func (i *Include) toParser(children []nginx.Parser) (nginx.Parser, error) {
	IN := &nginx.Include{
		BasicContext: nginx.BasicContext{
			Name:     "include",
			Value:    i.Value,
			Children: children,
		},
		Key:     nginx.NewKey("include", i.Value),
		Comment: nginx.NewComment(fmt.Sprintf("%s %s", "include", i.Value), false),
		ConfPWD: i.ConfPWD,
	}
	//IN, err := resolv.NewInclude(i.ConfPWD, i.Value, nil, &[]string{})
	//if err != nil {
	//	return nil, err
	//}
	//IN.Children = children
	return IN, nil
}

type Types struct {
	BasicContext `json:"types"`
}

func (t *Types) UnmarshalToJSON(b []byte, caches *nginx.Caches) (nginx.Parser, error) {
	return unmarshal(b, t, caches)
	// return unmarshal(b)
}

func (t *Types) toParser(children []nginx.Parser) (nginx.Parser, error) {
	T := nginx.NewTypes()
	T.Children = children

	return T, nil
}

type Map struct {
	BasicContext `json:"map"`
}

func (m *Map) UnmarshalToJSON(b []byte, caches *nginx.Caches) (nginx.Parser, error) {
	return unmarshal(b, m, caches)
	// return unmarshal(b)
}

func (m *Map) toParser(children []nginx.Parser) (nginx.Parser, error) {
	M := nginx.NewMap(m.Value)
	M.Children = children

	return M, nil
}

type LimitExcept struct {
	BasicContext `json:"limit_except"`
}

func (le *LimitExcept) UnmarshalToJSON(b []byte, caches *nginx.Caches) (nginx.Parser, error) {
	return unmarshal(b, le, caches)
	// return unmarshal(b)
}

func (le *LimitExcept) toParser(children []nginx.Parser) (nginx.Parser, error) {
	Le := nginx.NewLimitExcept(le.Value)
	Le.Children = children

	return Le, nil
}

type Events struct {
	BasicContext `json:"events"`
}

func (e *Events) UnmarshalToJSON(b []byte, caches *nginx.Caches) (nginx.Parser, error) {
	return unmarshal(b, e, caches)
	// return unmarshal(b)
}

func (e *Events) toParser(children []nginx.Parser) (nginx.Parser, error) {
	E := nginx.NewEvents()
	E.Children = children

	return E, nil
}

type Geo struct {
	BasicContext `json:"geo"`
}

func (g *Geo) UnmarshalToJSON(b []byte, caches *nginx.Caches) (nginx.Parser, error) {
	return unmarshal(b, g, caches)
	// return unmarshal(b)
}

func (g *Geo) toParser(children []nginx.Parser) (nginx.Parser, error) {
	G := nginx.NewGeo(g.Value)
	G.Children = children

	return G, nil
}

type Http struct {
	BasicContext `json:"http"`
}

func (h *Http) UnmarshalToJSON(b []byte, caches *nginx.Caches) (nginx.Parser, error) {
	return unmarshal(b, h, caches)
	// return unmarshal(b)
}

func (h *Http) toParser(children []nginx.Parser) (nginx.Parser, error) {
	H := nginx.NewHttp()
	H.Children = children

	return H, nil
}

type Stream struct {
	BasicContext `json:"stream"`
}

func (st *Stream) UnmarshalToJSON(b []byte, caches *nginx.Caches) (nginx.Parser, error) {
	return unmarshal(b, st, caches)
	// return unmarshal(b)
}

func (st *Stream) toParser(children []nginx.Parser) (nginx.Parser, error) {
	St := nginx.NewStream()
	St.Children = children

	return St, nil
}

type Upstream struct {
	BasicContext `json:"upstream"`
}

func (u *Upstream) UnmarshalToJSON(b []byte, caches *nginx.Caches) (nginx.Parser, error) {
	return unmarshal(b, u, caches)
	// return unmarshal(b)
}

func (u *Upstream) toParser(children []nginx.Parser) (nginx.Parser, error) {
	U := nginx.NewUpstream(u.Value)
	U.Children = children

	return U, nil
}

type Server struct {
	BasicContext `json:"server"`
}

func (s *Server) UnmarshalToJSON(b []byte, caches *nginx.Caches) (nginx.Parser, error) {
	return unmarshal(b, s, caches)
	// return unmarshal(b)
}

func (s *Server) toParser(children []nginx.Parser) (nginx.Parser, error) {
	S := nginx.NewServer()
	S.Children = children

	return S, nil
}

type Location struct {
	BasicContext `json:"location"`
}

func (l *Location) UnmarshalToJSON(b []byte, caches *nginx.Caches) (nginx.Parser, error) {
	return unmarshal(b, l, caches)
	// return unmarshal(b)
}

func (l *Location) toParser(children []nginx.Parser) (nginx.Parser, error) {
	L := nginx.NewLocation(l.Value)
	L.Children = children

	return L, nil
}

type If struct {
	BasicContext `json:"if"`
}

func (i *If) UnmarshalToJSON(b []byte, caches *nginx.Caches) (nginx.Parser, error) {
	return unmarshal(b, i, caches)
	// return unmarshal(b)
}

func (i *If) toParser(children []nginx.Parser) (nginx.Parser, error) {
	I := nginx.NewIf(i.Value)
	I.Children = children

	return I, nil
}

type Key struct {
	nginx.Key
}

func (k *Key) UnmarshalToJSON(b []byte, _ *nginx.Caches) (nginx.Parser, error) {
	err := json.Unmarshal(b, k)

	return &k.Key, err
}

func (k *Key) getChildren() []*json.RawMessage {
	return nil
}

func (k *Key) toParser(_ []nginx.Parser) (nginx.Parser, error) {
	return &k.Key, nil
}

type Comment struct {
	nginx.Comment
}

func (c *Comment) UnmarshalToJSON(b []byte, _ *nginx.Caches) (nginx.Parser, error) {
	err := json.Unmarshal(b, c)

	return &c.Comment, err
}

func (c *Comment) getChildren() []*json.RawMessage {
	return nil
}

func (c *Comment) toParser(_ []nginx.Parser) (nginx.Parser, error) {
	return &c.Comment, nil
}

// Unmarshal, json反序列化并返回nginx配置对象的函数.
func Unmarshal(b []byte) (*nginx.Config, error) {
	caches := nginx.NewCaches()
	// allConfigs := make(map[string]*nginx.Config, 0)
	parser, err := unmarshal(b, &Config{}, &caches)
	if err != nil {
		return nil, err
	}
	conf, ok := parser.(*nginx.Config)
	if !ok {
		return nil, fmt.Errorf("unmarshal error")
	}

	return conf, nil
	/* 测试解析器Config对象指针映射
	conf, err := unmarshal(b, &Config{}, &[]string{}, &allConfigs)
	fmt.Println(allConfigs)
	return conf, err */
}

// unmarshal, json反序列化并返回nginx配置对象的内部函数.
func unmarshal(b []byte, p unmarshaler, caches *nginx.Caches) (nginx.Parser, error) {
	switch p.(type) {
	case *Key, *Comment:
		return p.UnmarshalToJSON(b, caches)
	default:
		err := json.Unmarshal(b, p)
		if err != nil {
			return nil, err
		}
		conf, ok := p.(*Config)
		if !ok {
			children, cerr := unmarshalChildren(p.getChildren(), caches)
			if cerr != nil {
				return nil, cerr
			}

			return p.toParser(children)
		}

		if caches.IsCached(conf.Value) {
			conf, _ := caches.GetConfig(conf.Value)

			return conf, nil
		}

		setErr := caches.SetCache(&nginx.Config{BasicContext: nginx.BasicContext{
			Name:     nginx.TypeConfig,
			Value:    conf.Value,
			Children: nil,
		}}, hashForJsonUnmarshalTemp)
		if setErr != nil {
			return nil, setErr
		}

		children, cerr := unmarshalChildren(p.getChildren(), caches)
		if cerr != nil {
			return nil, cerr
		}

		// 获取解析器Config对象，并添加到解析器Config对象指针映射中
		config, tpErr := p.toParser(children)
		if tpErr != nil {
			return nil, tpErr
		}

		setErr = caches.SetCache(config.(*nginx.Config), hashForJsonUnmarshal)
		if setErr != nil {
			return nil, setErr
		}

		return config, nil
	}
}

// unmarshalChildren, 解析并反序列化子json串切片对象的内部函数.
//
//nolint:funlen,gocognit,gocyclo
func unmarshalChildren(bytes []*json.RawMessage, caches *nginx.Caches) (children []nginx.Parser, err error) {
	// parseContext, 用于解析json串归属于哪类需反序列化对象的匿名函数
	parseContext := func(b []byte, reg *regexp.Regexp) bool {
		m := reg.Find(b)

		return m != nil
	}

	for _, b := range bytes {
		switch {
		case parseContext(*b, RegCommentHead):
			comment, err := unmarshal(*b, &Comment{}, caches)
			if err != nil {
				return nil, err
			}
			children = append(children, comment)
		case parseContext(*b, RegIncludeHead):
			include, err := unmarshal(*b, &Include{}, caches)
			if err != nil {
				return nil, err
			}
			children = append(children, include)
		case parseContext(*b, RegConfigHead):
			config, err := unmarshal(*b, &Config{}, caches)
			if err != nil {
				return nil, err
			}
			children = append(children, config)
		case parseContext(*b, RegEventsHead):
			events, err := unmarshal(*b, &Events{}, caches)
			if err != nil {
				return nil, err
			}
			children = append(children, events)
		case parseContext(*b, RegGeoHead):
			geo, err := unmarshal(*b, &Geo{}, caches)
			if err != nil {
				return nil, err
			}
			children = append(children, geo)
		case parseContext(*b, RegHttpHead):
			http, err := unmarshal(*b, &Http{}, caches)
			if err != nil {
				return nil, err
				// return
			}
			children = append(children, http)
		case parseContext(*b, RegIfHead):
			i, err := unmarshal(*b, &If{}, caches)
			if err != nil {
				return nil, err
				// return
			}
			children = append(children, i)
		case parseContext(*b, RegLimitExceptHead):
			le, err := unmarshal(*b, &LimitExcept{}, caches)
			if err != nil {
				return nil, err
				// return
			}
			children = append(children, le)
		case parseContext(*b, RegLocationHead):
			l, err := unmarshal(*b, &Location{}, caches)
			if err != nil {
				return nil, err
				// return
			}
			children = append(children, l)
		case parseContext(*b, RegMapHead):
			m, err := unmarshal(*b, &Map{}, caches)
			if err != nil {
				return nil, err
				// return
			}
			children = append(children, m)
		case parseContext(*b, RegServerHead):
			svr, err := unmarshal(*b, &Server{}, caches)
			if err != nil {
				return nil, err
				// return
			}
			children = append(children, svr)
		case parseContext(*b, RegStreamHead):
			st, err := unmarshal(*b, &Stream{}, caches)
			if err != nil {
				return nil, err
				// return
			}
			children = append(children, st)
		case parseContext(*b, RegTypesHead):
			t, err := unmarshal(*b, &Types{}, caches)
			if err != nil {
				return nil, err
				// return
			}
			children = append(children, t)
		case parseContext(*b, RegUpstreamHead):
			u, err := unmarshal(*b, &Upstream{}, caches)
			if err != nil {
				return nil, err
				// return
			}
			children = append(children, u)
		default:
			k, err := unmarshal(*b, &Key{}, caches)
			if err != nil {
				return nil, err
				// return
			}
			children = append(children, k)
		}
	}

	return children, nil
}
