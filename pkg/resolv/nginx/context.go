// resove 包，该包包含了项目最基础的上下文相关对象，及相关方法及函数
// 创建者： ClessLi
// 创建时间：2020-1-17 11:14:15
package nginx

import (
	"fmt"
	"regexp"
	"strings"
)

var INDENT = "    "

// Context, 上下文接口对象，定义了上下文接口需实现的增、删、改等方法
type Context interface {
	Parser
	GetIndex(Parser) int
	Insert(indexParser Parser, insertParsers ...Parser) error
	Add(...Parser)
	Remove(...Parser)
	Modify(int, Parser) error
	Servers() []Parser
	Server() *Server
	Params() []Parser
	//getReg() string
	//Dict() map[string]interface{}
	//UnmarshalToJSON(b []byte) error
	//BumpChildDepth(int)
	dump() ([]string, error)
	List() ([]string, error)
}

// BasicContext, 上下文基础对象，定义了上下文类型的基本属性及基础方法
type BasicContext struct {
	Name     string   `json:"-"`
	Value    string   `json:"value,omitempty"`
	Children []Parser `json:"param,omitempty"`
}

// GetIndex, BasicContext 类子集对象索引值查询的方法， Context.GetIndex(Parser) 的实现
//
// 参数:
//     content: Parser接口对象，各配置对象指针
// 返回值:
//     索引值，未查询到时，返回-1
func (c *BasicContext) GetIndex(content Parser) int {
	for i, child := range c.Children {
		if content == child {
			return i
		}
	}
	return -1
}

// Insert, BasicContext 类插入对象的方法， Context.Insert(indexParser Parser, insertParsers ...Parser) error 的实现
//
// 参数:
//     indexParser: 基准索引子对象
//     insertParsers: 待插入子对象集
// 返回值:
//     错误
func (c *BasicContext) Insert(indexParser Parser, insertParsers ...Parser) error {
	index := c.GetIndex(indexParser)
	if index < 0 {
		return fmt.Errorf("'%s' is not a child of (name: %s, value: %s)", indexParser.String(), c.Name, c.Value)
	} else {
		// 时间、空间复杂度O(m+n)
		c.Children = append(append(c.Children[:index], insertParsers...), c.Children[index:]...)
		return nil
	}
}

// Add, BasicContext 类新增子对象的方法， Context.Add(...Parser) 的实现
func (c *BasicContext) Add(contents ...Parser) {
	for _, content := range contents {
		/*if _, isBC := content.(Context); isBC {
			content.(Context).BumpChildDepth(c.depth+1)
		}*/
		c.Children = append(c.Children, content)
	}
}

// Remove, BasicContext 类删除子对象的方法， Context.Remove(...Parser) 的实现
func (c *BasicContext) Remove(contents ...Parser) {
	for _, content := range contents {
		for index, child := range c.Children {
			if content == child {
				c.remove(index)
			}
		}
	}
}

// Modify, BasicContext 类修改子对象的方法， Context.Modify(int, Parser) error 的实现
func (c *BasicContext) Modify(index int, content Parser) error {
	switch content.(type) {
	case Context, *Comment, *Key:
		c.Children[index] = content
	default:
		return fmt.Errorf("conf format not supported with: %T", content)
	}
	return nil
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
	//return c.Servers()[0]
}

func (c *BasicContext) Params() (parsers []Parser) {
	parsers = make([]Parser, 0)
	for _, child := range c.Children {
		switch child.(type) {
		case *Key, *Comment:
			parsers = append(parsers, child)
		case *Include:
			for _, incChild := range child.(*Include).Children {
				if subConf, ok := incChild.(*Config); ok {
					parsers = append(parsers, subConf.Params()...)
				}
			}
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

func (c *BasicContext) filter(kw KeyWords) bool {
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
			if regexp.MustCompile(kw.Type).MatchString(c.Name) && regexp.MustCompile(kw.Value).MatchString(c.Value) {
				selfMatch = true
			}
		}

		if selfMatch {

			for _, childKW := range kw.ChildKWs {
				subMatch = false
				for _, child := range c.Children {
					if child.QueryAll(childKW) != nil {
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

func (c *BasicContext) subQueryAll(parsers []Parser, kw KeyWords) []Parser {
	for _, child := range c.Children {
		//parsers = append(parsers, child.QueryAll(kw)...)
		if tmpParsers := child.QueryAll(kw); tmpParsers != nil {
			parsers = append(parsers, tmpParsers...)
		}
	}
	return parsers
}

func (c *BasicContext) subQuery(kw KeyWords) Parser {
	for _, child := range c.Children {
		//parsers = append(parsers, child.QueryAll(kw)...)
		if tmpParser := child.Query(kw); tmpParser != nil {
			return tmpParser
		}
	}
	return nil
}

//func (c *BasicContext) getReg() string {
//	return c.Value
//}

//func (c *BasicContext) Dict() map[string]interface{} {
//}

func (c *BasicContext) String() []string {
	ret := make([]string, 0)

	contextTitle := c.getTitle()

	ret = append(ret, contextTitle)

	for _, child := range c.Children {
		switch child.(type) {
		case *Key:
			ret = append(ret, INDENT+child.String()[0])
		case *Comment:
			if child.(*Comment).Inline && len(ret) >= 1 {
				ret[len(ret)-1] = strings.TrimRight(ret[len(ret)-1], "\n") + "  " + child.String()[0]
			} else {
				ret = append(ret, INDENT+child.String()[0])
			}
		case Context:
			strs := child.String()
			//ret = append(ret, INDENT+strs[0])
			//for _, str := range strs[1:] {
			for _, str := range strs {
				ret = append(ret, INDENT+str)
			}
		default:
			str := child.String()
			if str != nil {
				ret = append(ret, str...)
			}
		}
	}
	ret[len(ret)-1] = RegEndWithCR.ReplaceAllString(ret[len(ret)-1], "}\n")
	ret = append(ret, "}\n\n")

	return ret
}

func (c *BasicContext) dump() ([]string, error) {
	ret := make([]string, 0)
	contextTitle := c.getTitle()
	ret = append(ret, contextTitle)

	for _, child := range c.Children {
		switch child.(type) {
		case *Key:
			ret = append(ret, INDENT+child.String()[0])
		case *Comment:
			if child.(*Comment).Inline && len(ret) >= 1 {
				ret[len(ret)-1] = strings.TrimRight(ret[len(ret)-1], "\n") + "  " + child.String()[0]
			} else {
				ret = append(ret, INDENT+child.String()[0])
			}
		case Context:
			strs, err := child.(Context).dump()
			if err != nil {
				return ret, err
			}

			for _, str := range strs {
				ret = append(ret, INDENT+str)
			}
		default:
			str := child.String()
			if str != nil {
				ret = append(ret, str...)
			}
		}
	}
	ret[len(ret)-1] = RegEndWithCR.ReplaceAllString(ret[len(ret)-1], "}\n")
	ret = append(ret, "}\n\n")

	return ret, nil
}

func (c *BasicContext) List() (ret []string, err error) {
	for _, child := range c.Children {
		switch child.(type) {
		case Context:
			//fmt.Println("farther", c.Name, "child", child)
			l, err := child.(Context).List()
			if err != nil {
				return nil, err
			}

			ret = append(ret, l...)
		}
	}
	return ret, nil
}

func (c *BasicContext) remove(index int) {
	c.Children = append(c.Children[:index], c.Children[index+1:]...)
}

func (c *BasicContext) getTitle() string {
	contextTitle := ""
	/*for i := 0; i < c.depth; i++ {
		contextTitle += INDENT
	}*/
	contextTitle += c.Name

	if c.Value != "" {
		contextTitle += " " + c.Value
	}

	contextTitle += " {\n"
	return contextTitle
}
