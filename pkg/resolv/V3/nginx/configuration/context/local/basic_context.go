package local

import (
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
)

const INDENT = "    "

type BasicContext struct {
	ContextType  context_type.ContextType `json:"context-type"`
	ContextValue string                   `json:"value,omitempty"`
	Children     []context.Context        `json:"params,omitempty"`

	father context.Context
	self   context.Context

	headStringFunc func(ctxType context_type.ContextType, value string) string
	tailStringFunc func() string
}

func (b *BasicContext) Insert(ctx context.Context, idx int) context.Context {
	// negative index
	if idx < 0 {
		return context.ErrContext(errors.WithCode(code.ErrV3ContextIndexOutOfRange, "index(%d) out of range", idx))
	}

	// refuse to insert nil
	if ctx == nil {
		return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "refuse to insert nil"))
	}

	// refuse to insert error context and config context
	switch ctx.Type() {
	case context_type.TypeErrContext:
		errctx, ok := ctx.(*context.ErrorContext)
		if ok {
			return errctx.AppendError(errors.WithCode(code.ErrV3InvalidOperation, "refuse to insert error context"))
		}
		return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "refuse to insert invalid context"))
	case context_type.TypeConfig:
		_, ok := ctx.(*Config)
		if ok {
			return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "refuse to insert config context"))
		}
		return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "refuse to insert invalid context"))
	}

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
		return context.ErrContext(errors.WithCode(code.ErrV3SetFatherContextFailed, err.Error()))
	}

	return b.self
}

func (b *BasicContext) Remove(idx int) context.Context {
	// negative index
	if idx < 0 {
		return context.ErrContext(errors.WithCode(code.ErrV3ContextIndexOutOfRange, "index(%d) out of range", idx))
	}

	if idx < b.Len() {
		// release father ctx
		err := b.Children[idx].SetFather(context.NullContext())
		if err != nil {
			return context.ErrContext(errors.WithCode(code.ErrV3SetFatherContextFailed, err.Error()))
		}

		b.Children = append(b.Children[:idx], b.Children[idx+1:]...)
	}
	return b.self
}

func (b *BasicContext) Modify(ctx context.Context, idx int) context.Context {
	// refuse to modify to nil
	if ctx == nil {
		return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "refuse to modify to nil"))
	}

	// refuse to modify to error context
	if ctx.Type() == context_type.TypeErrContext {
		errctx, ok := ctx.(*context.ErrorContext)
		if ok {
			return errctx.AppendError(errors.WithCode(code.ErrV3InvalidOperation, "refuse to modify to error context"))
		}
		return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "refuse to modify to invalid context"))
	}

	// if the context before and after modification is the same, no modification will be made
	if ctx == b.Child(idx) {
		return b.self
	}

	return b.Remove(idx).Insert(ctx, idx)
}

func (b *BasicContext) Father() context.Context {
	return b.father
}

func (b *BasicContext) Child(idx int) context.Context {
	if idx >= b.Len() || idx < 0 {
		return context.ErrContext(errors.WithCode(code.ErrV3ContextIndexOutOfRange, "index(%d) out of range", idx))
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
	return context.NotFoundPos()
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
	if len(title) > 0 {
		lines = append(lines, title)
	}
	for idx, child := range b.Children {
		if child == nil {
			return nil, errors.WithCode(code.ErrV3InvalidOperation, "child(index:%d) is nil", idx)
		}
		clines, err := child.ConfigLines(isDumping)
		if err != nil {
			return nil, err
		}
		if clines != nil {
			if child.Type() == context_type.TypeInlineComment && len(lines) > 0 &&
				b.Child(idx-1).Type() != context_type.TypeComment &&
				b.Child(idx-1).Type() != context_type.TypeInlineComment {
				lines[len(lines)-1] += INDENT + clines[0]
				continue
			}

			for _, cline := range clines {
				lines = append(lines, INDENT+cline)
			}
		}
	}

	if len(tail) > 0 {
		lines = append(lines, tail)
	}

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
