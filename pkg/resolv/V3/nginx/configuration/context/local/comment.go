package local

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"

	"github.com/marmotedu/errors"
)

type Comment struct {
	Comments string
	Inline   bool

	fatherContext context.Context
}

func (c *Comment) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ContextType context_type.ContextType `json:"context-type"`
		Value       string                   `json:"value"`
	}{
		ContextType: c.Type(),
		Value:       c.Value(),
	})
}

func (c *Comment) Insert(ctx context.Context, idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "comment cannot insert context"))
}

func (c *Comment) Remove(idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "comment cannot remove context"))
}

func (c *Comment) Modify(ctx context.Context, idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "comment cannot modify context"))
}

func (c *Comment) Father() context.Context {
	return c.fatherContext
}

func (c *Comment) FatherPosSet() context.PosSet {
	if c.fatherContext.Type() == context_type.TypeConfig || c.fatherContext.Type() == context_type.TypeMain {
		return c.fatherContext.FatherPosSet()
	}
	fatherPoses := context.NewPosSet()
	matched := false
	c.fatherContext.Father().ChildrenPosSet().Map(func(pos context.Pos) (context.Pos, error) {
		if !matched && pos.Target() == c.fatherContext {
			matched = true
			fatherPoses.Append(pos)
		}

		return pos, nil
	})

	return fatherPoses
}

func (c *Comment) Child(idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "comment has no context"))
}

func (c *Comment) ChildrenPosSet() context.PosSet {
	return context.NewPosSet()
}

func (c *Comment) Clone() context.Context {
	return &Comment{
		Comments:      c.Comments,
		Inline:        c.Inline,
		fatherContext: c.fatherContext,
	}
}

func (c *Comment) SetValue(v string) error {
	c.Comments = v

	return nil
}

func (c *Comment) SetFather(ctx context.Context) error {
	c.fatherContext = ctx

	return nil
}

func (c *Comment) HasChild() bool {
	return false
}

func (c *Comment) Len() int {
	return 0
}

func (c *Comment) Value() string {
	return c.Comments
}

func (c *Comment) Type() context_type.ContextType {
	if c.Inline {
		return context_type.TypeInlineComment
	}

	return context_type.TypeComment
}

func (c *Comment) Error() error {
	return nil
}

func (c *Comment) ConfigLines(isDumping bool) ([]string, error) {
	if len(strings.TrimSpace(c.Value())) == 0 {
		return []string{"#"}, nil
	}

	return []string{"# " + c.Value()}, nil
}

func (c *Comment) IsEnabled() bool {
	return true
}

func (c *Comment) Enable() context.Context {
	return c
}

func (c *Comment) Disable() context.Context {
	return c
}

func registerCommentBuilder() error {
	builderMap[context_type.TypeComment] = func(value string) context.Context {
		return &Comment{
			Comments:      value,
			Inline:        false,
			fatherContext: context.NullContext(),
		}
	}
	builderMap[context_type.TypeInlineComment] = func(value string) context.Context {
		return &Comment{
			Comments:      value,
			Inline:        true,
			fatherContext: context.NullContext(),
		}
	}

	return nil
}

func registerCommentParseFunc() error {
	inStackParseFuncMap[context_type.TypeComment] = func(data []byte, idx *int) context.Context {
		if subMatch := RegCommentHead.FindSubmatch(data[*idx:]); len(subMatch) == 3 { //nolint:nestif
			matchIndexes := RegCommentHead.FindIndex(data[*idx:])
			ct := context_type.TypeComment
			if !RegLineBreak.Match(subMatch[1]) && *idx != 0 {
				ct = context_type.TypeInlineComment
			}
			cmt := NewContext(ct, string(subMatch[2]))
			*idx += matchIndexes[len(matchIndexes)-1] - 1

			return cmt
		}

		return context.NullContext()
	}

	return nil
}

type CommentsToContextConverter interface {
	Convert(ctx context.Context) context.Context
}

type commentsToContextConverter struct{}

func (c commentsToContextConverter) queryComments(ctx context.Context) (map[context.Context]sort.IntSlice, error) {
	store := make(map[context.Context]sort.IntSlice)

	return store, ctx.ChildrenPosSet().QueryAll(context.NewKeyWords(context_type.TypeComment).SetCascaded(true)).Map(
		func(pos context.Pos) (context.Pos, error) {
			if err := pos.Target().Error(); err != nil {
				return pos, err
			}
			father, idx := pos.Position()
			store[father] = append(store[father], idx)

			return pos, nil
		},
	).Error()
}

func (c commentsToContextConverter) convertSubPoses(fatherCtx context.Context, subIdxes []int) error {
	insertIdx := subIdxes[0]
	// parse comments
	lines := make([]string, 0)
	for _, idx := range subIdxes {
		cmt, ok := fatherCtx.Child(idx).(*Comment)
		if !ok {
			return errors.WithCode(code.ErrV3ConversionToContextFailed, "cannot convert non Comment context")
		}
		line := cmt.Value()

		// context line
		if RegDirectiveWithValue.MatchString(line) ||
			RegDirectiveWithoutValue.MatchString(line) ||
			RegEventsHead.MatchString(line) ||
			RegGeoHead.MatchString(line) ||
			RegHttpHead.MatchString(line) ||
			RegIfHead.MatchString(line) ||
			RegLimitExceptHead.MatchString(line) ||
			RegLocationHead.MatchString(line) ||
			RegMapHead.MatchString(line) ||
			RegServerHead.MatchString(line) ||
			RegStreamHead.MatchString(line) ||
			RegTypesHead.MatchString(line) ||
			RegUpstreamHead.MatchString(line) ||
			RegBraceEnd.MatchString(line) ||
			RegCommentHead.MatchString(line+"\n") {
			lines = append(lines, line)
		} else {
			lines = append(lines, "# "+line)
		}
	}

	// parse bytes
	stack := newContextStack()
	// TODO: Optimize the matching mechanism for comments with no line breaks at the end
	data := []byte(strings.Join(lines, "\n") + "\n") // avoid missing line breaks at the end of comments that prevent proper matching
	idx := 0
	tmpCtx := &BasicContext{
		Enabled:     false,
		ContextType: "TempContext",
		Children:    make([]context.Context, 0),
		father:      context.NullContext(),
		headStringFunc: func(ctxType context_type.ContextType, value string) string {
			return ""
		},
		tailStringFunc: func() string {
			return ""
		},
	}
	tmpCtx.self = tmpCtx
	_ = stack.push(tmpCtx)
	for {
		isParsed := false
		if parseBlankLine(data, &idx) {
			continue
		}

		if matchIndexes := RegErrorHeed.FindIndex(data[idx:]); matchIndexes != nil {
			return errors.WithCode(code.ErrV3ConversionToContextFailed, "has parsed an error line")
		}

		if parseBraceEnd(data, &idx) {
			ctx, err := stack.pop()
			if err != nil {
				return errors.WithCode(code.ErrV3ConversionToContextFailed, "quit context from stack failed")
			}
			if ctx == tmpCtx {
				// Will split duplicate '}' into independent '# }' comments
				tmpCtx.Insert(NewContext(context_type.TypeComment, "}"), tmpCtx.Len())
				_ = stack.push(ctx)
			}

			continue
		}

		for _, parsefunc := range pushStackParseFuncMap {
			ctx := parsefunc(data, &idx)
			if ctx != context.NullContext() {
				father, err := stack.current()
				if err != nil {
					return errors.WithCode(code.ErrV3ConversionToContextFailed, "get father context failed")
				}
				err = father.Insert(ctx, father.Len()).Error()
				if err != nil {
					return errors.WithCode(code.ErrV3ConversionToContextFailed, "insert context failed")
				}

				err = stack.push(ctx)
				if err != nil {
					return errors.WithCode(code.ErrV3ConversionToContextFailed, "push context to stack failed")
				}
				isParsed = true

				break
			}
		}
		if isParsed {
			continue
		}

		for _, parsefunc := range inStackParseFuncMap {
			ctx := parsefunc(data, &idx)
			if ctx != context.NullContext() {
				father, err := stack.current()
				if err != nil {
					return errors.WithCode(code.ErrV3ConversionToContextFailed, "get father context failed")
				}
				err = father.Insert(ctx, father.Len()).Error()
				if err != nil {
					return errors.WithCode(code.ErrV3ConversionToContextFailed, "insert context failed")
				}
				isParsed = true

				break
			}
		}
		if isParsed {
			continue
		}

		if idx == len(data)-1 {
			break
		}

		return errors.WithCode(code.ErrV3ConversionToContextFailed, "the comments content was not successfully parsed")
	}
	finalCtx, err := stack.pop()
	if err != nil {
		return errors.WithCode(code.ErrV3ConversionToContextFailed, "stack exception for conversion and parsing")
	}
	if finalCtx != tmpCtx {
		return errors.WithCode(code.ErrV3ConversionToContextFailed, "lack of sufficient `context` ending symbols")
	}

	// rewrite converted ctx to father ctx
	for i := len(subIdxes) - 1; i >= 0; i-- {
		if e := fatherCtx.Remove(subIdxes[i]).Error(); e != nil {
			return e
		}
	}

	for i := finalCtx.Len() - 1; i >= 0; i-- {
		if finalCtx.Child(i).Type() != context_type.TypeComment && finalCtx.Child(i).Type() != context_type.TypeInlineComment {
			if e := c.Convert(finalCtx.Child(i)).Disable().Error(); e != nil {
				return e
			}
		}
		if e := fatherCtx.Insert(finalCtx.Child(i), insertIdx).Error(); e != nil {
			return e
		}
	}

	return nil
}

func (c commentsToContextConverter) Convert(ctx context.Context) context.Context {
	store, err := c.queryComments(ctx)
	if err != nil {
		return context.ErrContext(err)
	}
	for father, slice := range store {
		slice.Sort()
		conIdxes := c.sliceContinuousIndexes(slice)
		configLike := father.Type() == context_type.TypeConfig && len(conIdxes) == 1 && father.Len() == len(conIdxes[0])
		// reverse order conversion, automatic rewriting, and convergence of sub contexts
		for i := len(conIdxes) - 1; i >= 0; i-- {
			err = c.convertSubPoses(father, conIdxes[i])
			if err != nil {
				if errors.IsCode(err, code.ErrV3ConversionToContextFailed) { // skip conversion when conversion fails
					continue
				}

				return context.ErrContext(err)
			}
		}

		// determine and handle `configLike`
		for i := 0; i < father.Len() && configLike; i++ {
			configLike = father.Child(i).Type() == context_type.TypeComment ||
				father.Child(i).Type() == context_type.TypeInlineComment ||
				!father.Child(i).IsEnabled()
		}
		if configLike && father.IsEnabled() {
			father.Disable()
			for i := 0; i < father.Len(); i++ {
				father.Child(i).Enable()
			}
			c.Convert(father) // resolve and convert disabled children contexts in the disabled config context
		}
	}

	return ctx
}

func (c commentsToContextConverter) sliceContinuousIndexes(indexes []int) [][]int {
	conIdxes := make([][]int, 0)
	s := 0
	idxo := indexes[s]
	for e := 1; e <= len(indexes); e++ {
		if e == len(indexes) {
			conIdxes = append(conIdxes, indexes[s:e])

			continue
		}
		if indexes[e] != idxo+1 {
			if s == e-1 {
				conIdxes = append(conIdxes, []int{indexes[s]})
			} else {
				conIdxes = append(conIdxes, indexes[s:e])
			}
			s = e
		}
		idxo = indexes[e]
	}

	return conIdxes
}
