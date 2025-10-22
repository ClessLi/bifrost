package context

import (
	"github.com/ClessLi/bifrost/internal/pkg/code"

	"github.com/marmotedu/errors"
)

type Pos interface {
	Position() (Context, int)
	Target() Context
	QueryOne(kw KeyWords) Pos
	QueryAll(kw KeyWords) PosSet
}

type pos struct {
	father     Context
	childIndex int
}

func (p pos) Position() (Context, int) {
	return p.father, p.childIndex
}

func (p pos) Target() Context {
	return p.father.Child(p.childIndex)
}

func (p pos) QueryOne(kw KeyWords) Pos {
	if kw.IsToLeafQuery() {
		return p.Target().ChildrenPosSet().QueryOne(kw)
	}

	return p.Target().FatherPosSet().QueryOne(kw)
}

func (p pos) QueryAll(kw KeyWords) PosSet {
	if kw.IsToLeafQuery() {
		return p.Target().ChildrenPosSet().QueryAll(kw)
	}

	return p.Target().FatherPosSet().QueryAll(kw)
}

func SetPos(father Context, posIdx int) Pos {
	return &pos{
		father:     father,
		childIndex: posIdx,
	}
}

func GetPos(ctx Context) Pos {
	poses := ctx.Father().ChildrenPosSet().Filter(func(pos Pos) bool {
		return pos.Target() == ctx
	})
	if len(poses.List()) != 1 {
		return ErrPos(errors.WithCode(code.ErrV3InvalidOperation, "pos of the context not found, or duplicate pos"))
	}

	return poses.List()[0]
}

type errPos struct {
	ctx *ErrorContext
}

func (e errPos) Position() (Context, int) {
	return e.ctx, -1
}

func (e errPos) Target() Context {
	return e.ctx
}

func (e errPos) QueryOne(kw KeyWords) Pos {
	return e
}

func (e errPos) QueryAll(kw KeyWords) PosSet {
	return ErrPosSet(e.ctx.Error())
}

func ErrPos(err error) Pos {
	return &errPos{ctx: ErrContext(err).(*ErrorContext)}
}

var nullPos = &errPos{nullContext}

func NullPos() Pos {
	return nullPos
}

var notFoundPos = &errPos{ErrContext(errors.WithCode(code.ErrV3ContextNotFound, "queried context not found")).(*ErrorContext)}

func NotFoundPos() Pos {
	return notFoundPos
}

type PosSet interface {
	Filter(fn func(pos Pos) bool) PosSet
	Map(fn func(pos Pos) (Pos, error)) PosSet
	MapToPosSet(fn func(pos Pos) PosSet) PosSet
	// TODO: Reduce Method
	QueryOne(kw KeyWords) Pos
	QueryAll(kw KeyWords) PosSet
	List() []Pos
	Targets() []Context
	Append(pos ...Pos) PosSet
	AppendWithPosSet(posSet PosSet) PosSet
	Error() error
}

type posSet []Pos

func (s *posSet) Filter(fn func(pos Pos) bool) (result PosSet) {
	result = NewPosSet()
	for i := range *s {
		if fn((*s)[i]) {
			result.Append((*s)[i])
		}
	}

	return
}

func (s *posSet) Map(fn func(pos Pos) (Pos, error)) (result PosSet) {
	result = NewPosSet()
	for i := range *s {
		p, err := fn((*s)[i])
		if err != nil { // return Error PosSet, if the operation function returned an error during execution.
			return ErrPosSet(err)
		}
		if p == nil { // return Null Error PosSet, if the operation function returned a nil item during execution.
			return ErrPosSet(errors.WithCode(code.ErrV3InvalidOperation, "the execution of the operation function returned a nil item"))
		}
		result.Append(p)
	}

	return
}

func (s *posSet) MapToPosSet(fn func(pos Pos) PosSet) (result PosSet) {
	result = NewPosSet()
	if errSet := s.Map(func(pos Pos) (Pos, error) { return pos, result.AppendWithPosSet(fn(pos)).Error() }); errSet.Error() != nil {
		return errSet
	}

	return
}

func (s *posSet) QueryOne(kw KeyWords) Pos {
	for _, p := range s.Filter(func(pos Pos) bool { return !kw.SkipQueryThisContext(pos.Target()) }).List() {
		if kw.Match(p.Target()) {
			return p
		}

		// if query with non-cascaded KeyWords,
		// only the next poses of the current context will be used for retrieval matching.
		if kw.Cascaded() {
			result := p.QueryOne(kw)
			if err := result.Target().Error(); err == nil {
				return result
			}
		}
	}

	return NotFoundPos()
}

func (s *posSet) QueryAll(kw KeyWords) PosSet {
	filteredSet := s.Filter(func(pos Pos) bool { return !kw.SkipQueryThisContext(pos.Target()) })

	return filteredSet.Filter(func(pos Pos) bool { return kw.Match(pos.Target()) }).
		// if query with non-cascaded KeyWords,
		// only the next poses of the current context will be used for retrieval matching.
		AppendWithPosSet(filteredSet.Filter(func(pos Pos) bool { return kw.Cascaded() }).
			MapToPosSet(func(pos Pos) PosSet { return pos.QueryAll(kw) }))
}

func (s *posSet) List() []Pos {
	return *s
}

func (s *posSet) Targets() (result []Context) {
	s.Map(
		func(pos Pos) (Pos, error) {
			result = append(result, pos.Target())

			return pos, nil
		},
	)

	return
}

func (s *posSet) Append(pos ...Pos) PosSet {
	*s = append(*s, pos...)

	return s
}

func (s *posSet) AppendWithPosSet(posSet PosSet) PosSet {
	if posSet.Error() != nil {
		return posSet
	}

	return s.Append(posSet.List()...)
}

func (s *posSet) Error() error {
	return nil
}

func NewPosSet() PosSet {
	return new(posSet)
}

type errPosSet struct {
	err error
}

func (e *errPosSet) Filter(fn func(pos Pos) bool) PosSet {
	return e
}

func (e *errPosSet) Map(fn func(pos Pos) (Pos, error)) PosSet {
	return e
}

func (e *errPosSet) MapToPosSet(fn func(pos Pos) PosSet) PosSet {
	return e
}

func (e *errPosSet) QueryOne(kw KeyWords) Pos {
	return ErrPos(errors.Wrap(e.err, "can not query with error position"))
}

func (e *errPosSet) QueryAll(kw KeyWords) PosSet {
	return e
}

func (e *errPosSet) List() []Pos {
	return nil
}

func (e *errPosSet) Targets() []Context {
	return nil
}

func (e *errPosSet) Append(pos ...Pos) PosSet {
	return e
}

func (e *errPosSet) AppendWithPosSet(posSet PosSet) PosSet {
	return ErrPosSet(errors.NewAggregate([]error{e.err, posSet.Error()}))
}

func (e *errPosSet) Error() error {
	return e.err
}

func ErrPosSet(err error) PosSet {
	return &errPosSet{err: err}
}
