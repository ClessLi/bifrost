package local

import (
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
)

const INDENT = "    "

type BasicContext struct {
	ContextType  context_type.ContextType `json:"-"`
	ContextValue string                   `json:"value,omitempty"`
	Children     []context.Context        `json:"param,omitempty"`

	father context.Context
	self   context.Context

	headStringFunc func(ctxType context_type.ContextType, value string) string
	tailStringFunc func() string
}

func (b *BasicContext) Insert(ctx context.Context, idx int) context.Context {
	if idx >= b.Len() {
		idx = b.Len()
	}
	n := b.Len()
	b.Children = append(b.Children, nil)
	for i := n; i > idx; i-- {
		b.Children[i] = b.Children[i-1]
	}
	b.Children[idx] = ctx

	// set father for inserted ctx
	err := ctx.SetFather(b.self)
	if err != nil {
		return context.ErrContext(errors.WithCode(code.V3ErrSetFatherContextFailed, err.Error()))
	}

	return b.self
}

func (b *BasicContext) Remove(idx int) context.Context {
	if idx < b.Len() {
		// release father ctx
		err := b.Children[idx].SetFather(context.NullContext())
		if err != nil {
			return context.ErrContext(errors.WithCode(code.V3ErrSetFatherContextFailed, err.Error()))
		}

		b.Children = append(b.Children[:idx], b.Children[idx+1:]...)
	}
	return b.self
}

func (b *BasicContext) Modify(ctx context.Context, idx int) context.Context {
	if idx < b.Len() {
		b.Remove(idx).Insert(ctx, idx)
	}
	return b.self
}

func (b *BasicContext) Father() context.Context {
	return b.father
}

func (b *BasicContext) Child(idx int) context.Context {
	if idx >= b.Len() {
		return context.ErrContext(errors.WithCode(code.V3ErrContextIndexOutOfRange, "index(%d) out of range", idx))
	}
	return b.Children[idx]
}

func (b *BasicContext) QueryByKeyWords(kw context.KeyWords) context.Pos {
	for idx, child := range b.Children {
		if kw.Match(child) {
			return context.SetPos(b.self, idx)
		}

		// if query with non-cascaded KeyWords,
		// only the children of the current context will be used for retrieval matching.
		if kw.Cascaded() {
			pos := child.QueryByKeyWords(kw)
			c, _ := pos.Position()
			if c.Error() == nil {
				return pos
			}
		}
	}
	return context.NullPos()
}

func (b *BasicContext) QueryAllByKeyWords(kw context.KeyWords) []context.Pos {
	poses := make([]context.Pos, 0)
	for idx, child := range b.Children {
		if kw.Match(child) {
			poses = append(poses, context.SetPos(b.self, idx))
		}

		// if query with non-cascaded KeyWords,
		// only the children of the current context will be used for retrieval matching.
		if kw.Cascaded() {
			childPoses := child.QueryAllByKeyWords(kw)
			if len(childPoses) > 0 {
				poses = append(poses, childPoses...)
			}
		}
	}
	return poses
}

func (b *BasicContext) Clone() context.Context {
	clone := NewContext(b.Type(), b.Value())
	for i, child := range b.Children {
		clone.Insert(child.Clone(), i)
	}
	return clone
}

func (b *BasicContext) SetValue(v string) error {
	b.ContextValue = v
	return nil
}

func (b *BasicContext) SetFather(ctx context.Context) error {
	b.father = ctx
	return nil
}

func (b *BasicContext) HasChild() bool {
	return b.Len() > 0
}

func (b *BasicContext) Len() int {
	return len(b.Children)
}

func (b *BasicContext) Value() string {
	return b.ContextValue
}

func (b *BasicContext) Type() context_type.ContextType {
	return b.ContextType
}

func (b *BasicContext) Error() error {
	return nil
}

func (b *BasicContext) ConfigLines(isDumping bool) ([]string, error) {
	lines := make([]string, 0)
	title := b.headStringFunc(b.ContextType, b.Value())
	tail := b.tailStringFunc()
	if len(tail) == 0 {
		return nil, errors.WithCode(code.ErrParseFailed, "parse tail string failed")
	}
	if len(title) > 0 {
		lines = append(lines, title)
	}
	// TODO: watch out for nil
	for _, child := range b.Children {
		clines, err := child.ConfigLines(isDumping)
		if err != nil {
			return nil, err
		}
		if clines != nil {
			if child.Type() == context_type.TypeInlineComment && len(lines) > 0 {
				lines[len(lines)-1] += INDENT + clines[0]
				continue
			}

			for _, cline := range clines {
				lines = append(lines, INDENT+cline)
			}
		}
	}

	lines = append(lines, tail)

	return lines, nil
}

func newBasicContext(ctxType context_type.ContextType, head func(context_type.ContextType, string) string, tail func() string) BasicContext {
	return BasicContext{
		ContextType:    ctxType,
		Children:       make([]context.Context, 0),
		father:         context.NullContext(),
		headStringFunc: head,
		tailStringFunc: tail,
	}
}
