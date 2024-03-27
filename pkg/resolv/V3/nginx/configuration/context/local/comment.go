package local

import (
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
	"strings"
)

type Comment struct {
	Comments string `json:"comments,omitempty"`
	Inline   bool   `json:"inline,omitempty"`

	fatherContext context.Context
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

func (c *Comment) Child(idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "comment has no context"))
}

func (c *Comment) QueryByKeyWords(kw context.KeyWords) context.Pos {
	return context.NullPos()
}

func (c *Comment) QueryAllByKeyWords(kw context.KeyWords) []context.Pos {
	return nil
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

func NewComment(comments string, isInline bool) *Comment {
	return &Comment{
		Comments:      comments,
		Inline:        isInline,
		fatherContext: context.NullContext(),
	}
}

func registerCommentParseFunc() error {
	inStackParseFuncMap[context_type.TypeComment] = func(data []byte, idx *int) context.Context {
		if subMatch := RegCommentHead.FindSubmatch(data[*idx:]); len(subMatch) == 3 { //nolint:nestif
			matchIndexes := RegCommentHead.FindIndex(data[*idx:])
			cmt := NewComment(
				string(subMatch[2]),
				!RegLineBreak.Match(subMatch[1]) && *idx != 0,
			)
			*idx += matchIndexes[len(matchIndexes)-1] - 1

			return cmt
		}
		return context.NullContext()
	}
	return nil
}
