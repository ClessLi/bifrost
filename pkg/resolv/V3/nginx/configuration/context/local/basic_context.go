package local

import (
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"

	"github.com/marmotedu/errors"
)

const INDENT = "    "

type BasicContext struct {
	Enabled      bool                     `json:"enabled,omitempty"`
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

	return b.self.Remove(idx).Insert(ctx, idx)
}

func (b *BasicContext) Father() context.Context {
	return b.father
}

func (b *BasicContext) FatherPosSet() context.PosSet {
	if b.father.Type() == context_type.TypeConfig || b.father.Type() == context_type.TypeMain {
		return b.father.FatherPosSet()
	}
	fatherPoses := context.NewPosSet()
	matched := false
	b.father.Father().ChildrenPosSet().Map(func(pos context.Pos) (context.Pos, error) {
		if !matched && pos.Target() == b.father {
			matched = true
			fatherPoses.Append(pos)
		}

		return pos, nil
	})

	return fatherPoses
}

func (b *BasicContext) Child(idx int) context.Context {
	if idx >= b.Len() || idx < 0 {
		return context.ErrContext(errors.WithCode(code.ErrV3ContextIndexOutOfRange, "index(%d) out of range", idx))
	}

	return b.Children[idx]
}

func (b *BasicContext) ChildrenPosSet() context.PosSet {
	childrenPoses := context.NewPosSet()
	for i := 0; i < b.Len(); i++ {
		childrenPoses.Append(context.SetPos(b, i))
	}

	return childrenPoses
}

func (b *BasicContext) Clone() context.Context {
	clone := NewContext(b.Type(), b.Value())
	if !b.Enabled {
		clone.Disable()
	}
	for i, child := range b.Children {
		clone.Insert(child.Clone(), i)
	}

	return clone
}

func (b *BasicContext) SetValue(v string) error {
	b.ContextValue = v

	return nil
}

func (b *BasicContext) operateIncludes(handle func(include *Include) error) error {
	if b == nil || b.self == nil {
		return nil
	}
	// call include context handle some task
	return b.self.ChildrenPosSet().
		QueryAll(context.NewKeyWordsByType(context_type.TypeInclude).SetCascaded(true)).
		Filter(func(pos context.Pos) bool {
			_, ok := pos.Target().(*Include)

			return ok
		}).
		Map(func(pos context.Pos) (context.Pos, error) { return pos, handle(pos.Target().(*Include)) }).
		Error()
}

func (b *BasicContext) unloadIncludes() error {
	return b.operateIncludes(func(include *Include) error {
		return include.unload()
	})
}

func (b *BasicContext) loadIncludes() error {
	return b.operateIncludes(func(include *Include) error {
		return include.load()
	})
}

func (b *BasicContext) reloadIncludes() error {
	return b.operateIncludes(func(include *Include) error {
		return include.reload()
	})
}

func (b *BasicContext) SetFather(ctx context.Context) error {
	err := b.unloadIncludes()
	if err != nil {
		return err
	}
	b.father = ctx

	return b.loadIncludes()
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

	if !b.IsEnabled() && len(lines) > 0 {
		for i := range lines {
			lines[i] = "# " + lines[i]
		}
	}

	return lines, nil
}

func (b *BasicContext) IsEnabled() bool {
	return b.Enabled
}

func (b *BasicContext) Enable() context.Context {
	if b.Enabled {
		return b.self
	}
	b.Enabled = true
	err := b.reloadIncludes()
	if err != nil {
		return context.ErrContext(err)
	}

	return b.self
}

func (b *BasicContext) Disable() context.Context {
	if !b.Enabled {
		return b.self
	}
	b.Enabled = false
	err := b.reloadIncludes()
	if err != nil {
		return context.ErrContext(err)
	}

	return b.self
}

func newBasicContext(ctxType context_type.ContextType, head func(context_type.ContextType, string) string, tail func() string) BasicContext {
	return BasicContext{
		Enabled:        true,
		ContextType:    ctxType,
		Children:       make([]context.Context, 0),
		father:         context.NullContext(),
		headStringFunc: head,
		tailStringFunc: tail,
	}
}

func getFatherContextByType(ctx context.Context, contextType context_type.ContextType) context.Context {
	if ctx.Type() == context_type.TypeMain {
		return errCtxGetFatherCtxFromMainByType
	}

	father := ctx.Father()
	if father.Error() != nil || father.Type() == contextType {
		return father
	}

	return getFatherContextByType(father, contextType)
}
