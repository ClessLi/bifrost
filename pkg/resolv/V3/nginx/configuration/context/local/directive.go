package local

import (
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
	"strings"
)

type Directive struct {
	Name   string `json:"directive,omitempty"`
	Params string `json:"params,omitempty"`

	fatherContext context.Context
}

func (d *Directive) Insert(ctx context.Context, idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "directive cannot insert context"))
}

func (d *Directive) Remove(idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "directive cannot remove context"))
}

func (d *Directive) Modify(ctx context.Context, idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "directive cannot modify context"))
}

func (d *Directive) Father() context.Context {
	return d.fatherContext
}

func (d *Directive) Child(idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "directive has no child context"))
}

func (d *Directive) QueryByKeyWords(kw context.KeyWords) context.Pos {
	return context.NullPos()
}

func (d *Directive) QueryAllByKeyWords(kw context.KeyWords) []context.Pos {
	return nil
}

func (d *Directive) Clone() context.Context {
	return &Directive{
		Name:          d.Name,
		Params:        d.Params,
		fatherContext: d.fatherContext,
	}
}

func (d *Directive) SetValue(v string) error {
	kv := strings.SplitN(strings.TrimSpace(v), " ", 2)
	if len(strings.TrimSpace(kv[0])) == 0 {
		return errors.WithCode(code.ErrV3InvalidOperation, "set value for directive failed, cased by: split null value")
	}
	d.Name = strings.TrimSpace(kv[0])
	if len(kv) == 2 {
		d.Params = strings.TrimSpace(kv[1])
	} else {
		d.Params = ""
	}
	return nil
}

func (d *Directive) SetFather(ctx context.Context) error {
	d.fatherContext = ctx
	return nil
}

func (d *Directive) HasChild() bool {
	return false
}

func (d *Directive) Len() int {
	return 0
}

func (d *Directive) Value() string {
	v := strings.TrimSpace(d.Name)
	if params := strings.TrimSpace(d.Params); len(params) > 0 {
		v += " " + params
	}
	return v
}

func (d *Directive) Type() context_type.ContextType {
	return context_type.TypeDirective
}

func (d *Directive) Error() error {
	return nil
}

func (d *Directive) ConfigLines(isDumping bool) ([]string, error) {
	return []string{d.Value() + ";"}, nil
}

func NewDirective(name, params string) *Directive {
	return &Directive{
		Name:          strings.TrimSpace(name),
		Params:        strings.TrimSpace(params),
		fatherContext: context.NullContext(),
	}
}

func registerDirectiveParseFunc() error {
	inStackParseFuncMap[context_type.TypeDirective] = func(data []byte, idx *int) context.Context {
		if matchIndexes := RegDirectiveWithoutValue.FindIndex(data[*idx:]); matchIndexes != nil { //nolint:nestif
			subMatch := RegDirectiveWithoutValue.FindSubmatch(data[*idx:])
			*idx += matchIndexes[len(matchIndexes)-1]
			key := string(subMatch[1])
			return NewDirective(key, "")
		}

		if matchIndexes := RegDirectiveWithValue.FindIndex(data[*idx:]); matchIndexes != nil { //nolint:nestif
			subMatch := RegDirectiveWithValue.FindSubmatch(data[*idx:])
			*idx += matchIndexes[len(matchIndexes)-1]
			name := string(subMatch[1])
			value := string(subMatch[2])
			if name == string(context_type.TypeInclude) {
				return context.NullContext()
			}
			return NewDirective(name, value)
		}
		return context.NullContext()
	}
	return nil
}
