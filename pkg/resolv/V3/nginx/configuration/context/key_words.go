package context

type KeyWords interface {
	Match(ctx Context) bool
	Cascaded() bool
	SetCascaded(cascaded bool)
}
