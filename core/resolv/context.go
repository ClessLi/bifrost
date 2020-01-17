// resove 包，该包包含了项目最基础的上下文相关对象，及相关方法及函数
// 创建者： ClessLi
// 创建时间：2020-1-17 11:14:15
package resolv

import (
	"fmt"
)

// Context, 上下文接口对象，定义了上下文接口需实现的增、删、改等方法
type Context interface {
	Add(...interface{})
	Remove(...interface{})
	Modify(int, interface{}) error
	//Filter(string, string) *Context
	//getReg() string
}

// BasicContext, 上下文基础对象，定义了上下文类型的基本属性及基础方法
type BasicContext struct {
	Name     string
	Value    string
	depth    int
	Children []interface{}
}

// Add, BasicContext 类新增子对象的方法， Context.Add(...interface{}) 的实现
func (c *BasicContext) Add(contents ...interface{}) {
	for _, content := range contents {
		c.Children = append(c.Children, content)
	}
}

// Remove, BasicContext 类删除子对象的方法， Context.Remove(...interface{}) 的实现
func (c *BasicContext) Remove(contents ...interface{}) {
	for _, content := range contents {
		for index, child := range c.Children {
			if content == child {
				c.remove(index)
			}
		}
	}
}

// Modify, BasicContext 类修改子对象的方法， Context.Modify(int, interface{}) error 的实现
func (c *BasicContext) Modify(index int, content interface{}) error {
	switch content.(type) {
	case Context:
		c.Children[index] = content
	case Comment:
		c.Children[index] = content
	case Key:
		c.Children[index] = content
	default:
		return fmt.Errorf("conf format not supported with: %T", content)
	}
	return nil
}

// Filter, BasicContext 类生成过滤对象的方法， Context.Filter(string, string) []*Context 的实现
//TODO: 过滤器
//func (c *BasicContext) Filter(btype, name string) *Context {
//}

//func (c *BasicContext) getReg() string {
//	return c.Value
//}

func (c *BasicContext) remove(index int) {
	c.Children = append(c.Children[:index], c.Children[index+1:])
}

type Comment struct {
	Comments string
	inline   bool
}

type Key struct {
	Name  string
	Value string
}
