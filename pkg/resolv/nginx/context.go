// resolv.nginx 包，该包包含了bifrost.pkg.resolv.nginx项目最基础的上下文相关对象，及相关方法及函数
// 创建者： ClessLi
// 创建时间：2020-1-17 11:14:15
package nginx

import (
	"fmt"
	"regexp"
	"strings"
)

var INDENT = "    "

// Context, 上下文接口对象，定义了上下文接口需实现的增、删、改等方法.
type Context interface {
	Parser
	Insert(indexParser Parser, pType parserType, values ...string) error
	InsertByParser(indexParser Parser, contents ...Parser) error
	Add(parserType, ...string) error
	AddByParser(...Parser)
	Remove(parserType, ...string) error
	RemoveByParser(...Parser)
	Modify(indexParser Parser, pType parserType, values ...string) error
	ModifyByParser(Parser, Parser) error
	Servers() []Parser
	Server() *Server
	Params() []Parser
	// getReg() string
	// Dict() map[string]interface{}
	// UnmarshalToJSON(b []byte) error
	// BumpChildDepth(int)
	dump(string, *Caches, int) (map[string][]string, error) // TODO: 优化将dumps数据string切片换为byte切片，用于将数据直接解析给缓存
	List() (Caches, error)
}

// BasicContext, 上下文基础对象，定义了上下文类型的基本属性及基础方法.
type BasicContext struct {
	Name     parserType `json:"-"`
	Value    string     `json:"value,omitempty"`
	Children []Parser   `json:"param,omitempty"`
}

func (c *BasicContext) QueryAll(pType parserType, isRec bool, values ...string) []Parser {
	kw, err := newKW(pType, values...)
	if err != nil {
		return nil
	}
	kw.IsRec = isRec
	return c.subQueryAll([]Parser{}, *kw)
}

func (c *BasicContext) Query(pType parserType, isRec bool, values ...string) Parser {
	kw, err := newKW(pType, values...)
	if err != nil {
		return nil
	}
	kw.IsRec = isRec
	return c.subQuery(*kw)
}

func (c *BasicContext) Insert(indexParser Parser, pType parserType, values ...string) error {
	if values != nil {
		p, err := newParser(pType, values...)
		if err != nil {
			return err
		}

		return c.InsertByParser(indexParser, p)
	}
	return ParserControlNoParamError
}

// InsertByParser, BasicContext 类插入对象的方法， Context.InsertByParser(indexParser Parser, insertParsers ...Parser) error 的实现
//
// 参数:
//
//	indexParser: 基准索引子对象
//	insertParsers: 待插入子对象集
//
// 返回值:
//
//	错误
func (c *BasicContext) InsertByParser(indexParser Parser, contents ...Parser) error {
	for i, child := range c.Children {
		if indexParser != child {
			switch child.(type) {
			case *Include:
				err := child.(*Include).InsertByParser(indexParser, contents...)
				if err == ParserControlIndexNotFoundError {
					continue
				}
				return err
			}
		} else {
			// 时间、空间复杂度O(m+n)
			// tmp := append(c.Children[:i], contents...)
			// c.Children = append(tmp, c.Children[i:]...)
			return c.childInsert(i, contents...)
			// c.Children = append(append(c.Children[0:i], contents...), c.Children[i:]...)
			// return nil
		}
	}
	return ParserControlIndexNotFoundError
}

func (c *BasicContext) Add(pType parserType, values ...string) error {
	if values != nil {
		parser, err := newParser(pType, values...)
		if err != nil {
			return err
		}
		c.AddByParser(parser)
		return nil
	}
	return ParserControlNoParamError
}

// AddByParser, BasicContext 类新增子对象的方法， Context.AddByParser(...Parser) 的实现.
func (c *BasicContext) AddByParser(contents ...Parser) {
	c.Children = append(c.Children, contents...)
}

func (c *BasicContext) Remove(pType parserType, values ...string) error {
	c.RemoveByParser(c.QueryAll(pType, false, values...)...)
	return nil
}

// RemoveByParser, BasicContext 类删除子对象的方法， Context.RemoveByParser(...Parser) 的实现.
func (c *BasicContext) RemoveByParser(contents ...Parser) {
	for _, content := range contents {
		for i, child := range c.Children {
			if content == child {
				c.removeByIndex(i)
			} else {
				switch child.(type) {
				case *Include:
					child.(*Include).RemoveByParser(content)
				}
			}
		}
	}
}

func (c *BasicContext) Modify(indexParser Parser, pType parserType, values ...string) error {
	if values != nil {
		ctx, err := newParser(pType, values...)
		if err != nil {
			return err
		}

		return c.ModifyByParser(indexParser, ctx)
	}
	return ParserControlNoParamError
}

// ModifyByParser, BasicContext 类修改子对象的方法， Context.ModifyByParser(int, Parser) error 的实现.
func (c *BasicContext) ModifyByParser(indexParser, content Parser) error {
	for i, child := range c.Children {
		if child != indexParser {
			switch child.(type) {
			case *Include:
				err := child.(*Include).ModifyByParser(indexParser, content)
				if err == ParserControlIndexNotFoundError {
					continue
				} else {
					return err
				}
			}
			continue
		} else {
			c.Children[i] = content
			return nil
		}
	}
	return ParserControlIndexNotFoundError
}

func (c *BasicContext) Servers() []Parser {
	svrs := make([]Parser, 0)
	for _, child := range c.Children {
		switch child.(type) {
		case Context:
			switch child.(type) {
			case *Server:
				svrs = append(svrs, child)
			default:
				svrs = append(svrs, child.(Context).Servers()...)
			}
		}
	}
	return svrs
}

func (c *BasicContext) Server() *Server {
	for _, child := range c.Children {
		switch child.(type) {
		case Context:
			switch child.(type) {
			case *Server:
				return child.(*Server)
			default:
				if s := child.(Context).Server(); s != nil {
					return s
				}
			}
		}
	}
	return nil
}

func (c *BasicContext) Params() (parsers []Parser) {
	parsers = make([]Parser, 0)
	for _, child := range c.Children {
		switch child.(type) {
		case *Key, *Comment:
			parsers = append(parsers, child)
		case *Include:
			parsers = append(parsers, child.(*Include).Params()...)
		default:
			n := len(parsers)
			for n > 0 {
				if comment, ok := parsers[n-1].(*Comment); ok && !comment.Inline {
					parsers = parsers[:n-1]
					n--
				} else {
					break
				}
			}
		}
	}
	return
}

func (c *BasicContext) filter(kw Keywords) bool {
	var (
		selfMatch = false
		subMatch  = true
	)
	switch kw.Type {
	case TypeKey, TypeComment:
	default:
		switch kw.Type {
		case TypeEvents, TypeHttp, TypeServer, TypeStream, TypeTypes:
			kw.Value = ""
		}

		if !kw.IsReg {
			if kw.Type == c.Name && kw.Value == c.Value {
				selfMatch = true
			}
		} else {
			if regexp.MustCompile(kw.Type.ToString()).MatchString(c.Name.ToString()) && regexp.MustCompile(kw.Value).MatchString(c.Value) {
				selfMatch = true
			}
		}

		if selfMatch {
			for _, childKW := range kw.ChildKWs {
				subMatch = false
				for _, child := range c.Children {
					if child.QueryAllByKeywords(childKW) != nil {
						subMatch = true
						break
					}
				}

				if !subMatch {
					return selfMatch && subMatch
				}
			}
		}
	}

	return selfMatch && subMatch
}

func (c *BasicContext) subQueryAll(parsers []Parser, kw Keywords) []Parser {
	for _, child := range c.Children {
		if tmpParsers := child.QueryAllByKeywords(kw); tmpParsers != nil {
			parsers = append(parsers, tmpParsers...)
		}
	}
	return parsers
}

func (c *BasicContext) subQuery(kw Keywords) Parser {
	for _, child := range c.Children {
		if tmpParser := child.QueryByKeywords(kw); tmpParser != nil {
			return tmpParser
		}
	}
	return nil
}

func (c BasicContext) String() []string {
	caches := NewCaches()
	return c.string(&caches, 0)
}

func (c *BasicContext) string(caches *Caches, deep int) []string {
	ret := make([]string, 0)
	ind := strings.Repeat(INDENT, deep)

	contextTitle := ind + c.getTitle()

	ret = append(ret, contextTitle)

	for _, child := range c.Children {
		switch child.(type) {
		case *Key:
			ret = append(ret, child.string(caches, deep+1)[0])
		case *Comment:
			if child.(*Comment).Inline && len(ret) >= 1 {
				ret[len(ret)-1] = strings.TrimRight(ret[len(ret)-1], "\n") + "  " + child.String()[0]
			} else {
				ret = append(ret, child.string(caches, deep+1)[0])
			}
		case Context:
			strs := child.string(caches, deep+1)
			//for _, str := range strs {
			//	ret = append(ret, INDENT+str)
			//}
			ret = append(ret, strs...)
		default:
			str := child.string(caches, deep+1)
			if str != nil {
				ret = append(ret, str...)
			}
		}
	}
	ret[len(ret)-1] = RegEndWithCR.ReplaceAllString(ret[len(ret)-1], "}\n")
	ret = append(ret, ind+"}\n")

	return ret
}

func (c *BasicContext) dump(configPath string, caches *Caches, deep int) (map[string][]string, error) {
	dumps := make(map[string][]string)
	dmp := make([]string, 0)
	ind := strings.Repeat(INDENT, deep)
	contextTitle := ind + c.getTitle()
	dmp = append(dmp, contextTitle)

	for _, child := range c.Children {
		switch child.(type) {
		case *Key:
			dmp = append(dmp, child.string(caches, deep+1)[0])
		case *Comment:
			if child.(*Comment).Inline && len(dmp) >= 1 {
				dmp[len(dmp)-1] = strings.TrimRight(dmp[len(dmp)-1], "\n") + "  " + child.String()[0]
			} else {
				dmp = append(dmp, child.string(caches, deep+1)[0])
			}
		case Context:
			newDumps, err := child.(Context).dump(configPath, caches, deep+1)
			if err != nil {
				return nil, err
			}

			for dmpPath, data := range newDumps {
				if _, ok := dumps[dmpPath]; ok || dmpPath == configPath {
					continue
				}
				dumps[dmpPath] = data
			}

			if d, ok := newDumps[configPath]; ok {
				//for _, line := range d {
				//	dmp = append(dmp, INDENT+line)
				//}
				dmp = append(dmp, d...)
			}
		default:
			lines := child.string(caches, deep+1)
			if lines != nil {
				dmp = append(dmp, lines...)
			}
		}
	}
	dmp[len(dmp)-1] = RegEndWithCR.ReplaceAllString(dmp[len(dmp)-1], "}\n")
	dmp = append(dmp, ind+"}\n")

	dumps[configPath] = dmp

	return dumps, nil
}

func (c *BasicContext) List() (caches Caches, err error) {
	caches = NewCaches()
	for _, child := range c.Children {
		switch child.(type) {
		case Context:
			subCaches, err := child.(Context).List()
			if err != nil {
				return nil, err
			} else if subCaches != nil {
				for _, cache := range subCaches {
					err = caches.SetCache(cache.config, hashForGetList)

					if err != nil && err != IsInCaches {
						return nil, err
					}
				}
			}
		}
	}
	return
}

func (c *BasicContext) removeByIndex(index int) {
	c.Children = append(c.Children[0:index], c.Children[index+1:]...)
}

func (c *BasicContext) getTitle() string {
	contextTitle := ""
	/*for i := 0; i < c.depth; i++ {
		contextTitle += INDENT
	}*/
	contextTitle += c.Name.ToString()

	if c.Value != "" {
		contextTitle += " " + c.Value
	}

	contextTitle += " {\n"
	return contextTitle
}

func (c *BasicContext) childInsert(i int, contents ...Parser) error {
	if contents != nil {
		max := len(c.Children)
		if max > 0 && max > i {
			//cLen := len(contents)
			//mvLen := max-i
			//// 扩容切片
			//for a := 0; a < cLen; a++ {
			//	c.Children = append(c.Children, c.Children[max-1])
			//	max++
			//}
			//// 移动元素
			//for m := 0; m < mvLen; m++ {
			//	c.Children[max-m-2] = c.Children[max-cLen-2-m]
			//}
			//// 插入元素
			//for j := 0; j < cLen; j++ {
			//	c.Children[i+j] = contents[j]
			//}
			tmp := append([]Parser{}, c.Children[i:]...)
			c.Children = append(append(c.Children[:i], contents...), tmp...)

			return nil
		} else if max == 0 && i == 0 {
			c.Children = append(c.Children, contents...)
			return nil
		}
	}
	return ParserControlParamsError
}

func newParser(pType parserType, values ...string) (Parser, error) {
	var parser Parser
	if values != nil {
		isMatch := false
		switch pType {
		case TypeComment:
			if ms := regexp.MustCompile(`^#+[ \r\t\f]*(.*)$`).FindStringSubmatch(values[0]); len(ms) == 2 {
				return NewComment(ms[1], false), nil
			} else {
				return nil, ParserControlParamsError
			}
		case TypeKey:
			keyValue := ""
			kv := strings.Split(values[0], ":")
			if len(kv) > 1 {
				keyValue = strings.Join(kv[1:], ":")
			}
			keyName := kv[0]
			return NewKey(keyName, keyValue), nil
		case TypeGeo:
			parser = NewGeo(values[0])
			isMatch = true
		case TypeIf:
			parser = NewIf(values[0])
			isMatch = true
		case TypeLimitExcept:
			parser = NewLimitExcept(values[0])
			isMatch = true
		case TypeLocation:
			parser = NewLocation(values[0])
			isMatch = true
		case TypeMap:
			parser = NewMap(values[0])
			isMatch = true
		case TypeUpstream:
			parser = NewUpstream(values[0])
			isMatch = true
		case TypeEvents:
			parser = NewEvents()
		case TypeHttp:
			parser = NewHttp()
		case TypeServer:
			parser = NewServer()
		case TypeStream:
			parser = NewStream()
		case TypeTypes:
			parser = NewTypes()
		default:
			return nil, fmt.Errorf("unknown nginx context type: %s", pType)
		}

		if isMatch {
			if len(values) > 1 {
				values = values[1:]
			} else {
				values = nil
			}
		}
	} else {
		switch pType {
		case TypeEvents:
			parser = NewEvents()
		case TypeHttp:
			parser = NewHttp()
		case TypeServer:
			parser = NewServer()
		case TypeStream:
			parser = NewStream()
		case TypeTypes:
			parser = NewTypes()
		default:
			return nil, fmt.Errorf("unknown nginx context type: %s", pType)
		}
	}

	if ctx, ok := parser.(Context); ok && values != nil {
		for _, value := range values {
			if ms := regexp.MustCompile(`#+[ \r\t\f]*(.*?)`).FindStringSubmatch(value); len(ms) == 2 {
				err := ctx.Add(TypeComment, ms[1])
				if err != nil {
					return nil, err
				}
			} else {
				err := ctx.Add(TypeKey, value)
				if err != nil {
					return nil, err
				}
			}
		}
		return ctx, nil
	}
	return parser, nil
}
