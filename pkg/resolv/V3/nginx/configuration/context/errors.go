package context

import (
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
)

const ()

var (
	ErrIndexOutOfRange             = errors.New("index out of range")
	ErrInsertParserTypeError       = errors.New("insert parser type error")
	ErrInsertIntoErrorContext      = errors.WithCode(code.ErrV3OperationOnErrorContext, "insert context into an error context")
	ErrRemoveFromErrorContext      = errors.WithCode(code.ErrV3OperationOnErrorContext, "remove context from an error context")
	ErrModifyFromErrorContext      = errors.WithCode(code.ErrV3OperationOnErrorContext, "modify context from an error context")
	ErrSetValueToErrorContext      = errors.WithCode(code.ErrV3OperationOnErrorContext, "set value to an error context")
	ErrSetFatherToErrorContext     = errors.WithCode(code.ErrV3OperationOnErrorContext, "set father to an error context")
	ErrGetFatherFromErrorContext   = errors.WithCode(code.ErrV3OperationOnErrorContext, "get father from an error context")
	ErrGetChildFromErrorContext    = errors.WithCode(code.ErrV3OperationOnErrorContext, "get child from an error context")
	ErrParseConfigFromErrorContext = errors.WithCode(code.ErrV3OperationOnErrorContext, "parse config from an error context")
	ErrNullPosition                = errors.WithCode(code.ErrV3NullContextPosition, "null position")

	ErrNullContext = errors.New("null context")
)

type ErrorContext struct {
	errors []error
}

func (e *ErrorContext) IsEnabled() bool {
	return true
}

func (e *ErrorContext) Enable() Context {
	return e
}

func (e *ErrorContext) Disable() Context {
	return e
}

func ErrContext(errs ...error) Context {
	return errContext(errs)
}

func errContext(errs []error) *ErrorContext {
	if errs == nil {
		return nullContext
	}

	cErrs := make([]error, 0)
	for _, err := range errs {
		if err != nil {
			cErrs = append(cErrs, err)
		}
	}

	if len(cErrs) == 0 {
		return nullContext
	}

	return &ErrorContext{errors: cErrs}
}

func (e *ErrorContext) Insert(ctx Context, idx int) Context {
	return e.AppendError(ErrInsertIntoErrorContext)
}

func (e *ErrorContext) Remove(idx int) Context {
	return e.AppendError(ErrRemoveFromErrorContext)
}

func (e *ErrorContext) Modify(ctx Context, idx int) Context {
	return e.AppendError(ErrModifyFromErrorContext)
}

func (e *ErrorContext) Father() Context {
	return e.AppendError(ErrGetFatherFromErrorContext)
}

func (e *ErrorContext) Child(idx int) Context {
	return e.AppendError(ErrGetChildFromErrorContext)
}

func (e *ErrorContext) QueryByKeyWords(kw KeyWords) Pos {
	return nullPos
}

func (e *ErrorContext) QueryAllByKeyWords(kw KeyWords) []Pos {
	return nil
}

func (e *ErrorContext) Clone() Context {
	return e.clone()
}

func (e *ErrorContext) clone() *ErrorContext {
	return errContext(e.errors)
}

func (e *ErrorContext) SetValue(v string) error {
	return e.AppendError(ErrSetValueToErrorContext).Error()
}

func (e *ErrorContext) SetFather(ctx Context) error {
	return e.AppendError(ErrSetFatherToErrorContext).Error()
}

func (e *ErrorContext) HasChild() bool {
	return false
}

func (e *ErrorContext) Len() int {
	return 0
}

func (e *ErrorContext) Value() string {
	return ""
}

func (e *ErrorContext) Type() context_type.ContextType {
	return context_type.TypeErrContext
}

func (e *ErrorContext) Error() error {
	return errors.NewAggregate(e.errors)
}

func (e *ErrorContext) ConfigLines(isDumping bool) ([]string, error) {
	return nil, e.AppendError(ErrParseConfigFromErrorContext).Error()
}

func (e *ErrorContext) AppendError(err error) Context {
	if err != nil {
		if e == nullContext {
			clone := e.clone()
			clone.errors = append(clone.errors, err)
			return clone
		}
		e.errors = append(e.errors, err)
	}
	return e
}

var nullContext = &ErrorContext{errors: []error{ErrNullContext}}

func NullContext() Context {
	return nullContext
}
