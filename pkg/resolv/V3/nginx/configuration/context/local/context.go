package local

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"

	"github.com/marmotedu/errors"
)

type BuildBasicContextConfig struct {
	ContextType    context_type.ContextType
	headStringFunc func(ctxType context_type.ContextType, value string) string
	tailStringFunc func() string
}

func OptsApplyTo(opts context.BuildOptions) (BuildBasicContextConfig, error) {
	buildConfig := BuildBasicContextConfig{}
	buildConfig.headStringFunc, buildConfig.tailStringFunc = buildHeadAndTailStringFuncs(opts)
	if buildConfig.headStringFunc == nil || buildConfig.tailStringFunc == nil {
		return buildConfig, errors.New("TODO")
	}

	buildConfig.ContextType = opts.ContextType

	return buildConfig, nil
}

func (b BuildBasicContextConfig) BasicContext() *BasicContext {
	return &BasicContext{
		Enabled:        true,
		ContextType:    b.ContextType,
		Children:       make([]context.Context, 0),
		father:         context.NullContext(),
		headStringFunc: b.headStringFunc,
		tailStringFunc: b.tailStringFunc,
	}
}

type ContextBuilderRegistrar func(func() *BasicContext) func(value string) context.Context

func hasValueBraceHeadString(ctxType context_type.ContextType, value string) string {
	contextTitle := ctxType.String()

	if value != "" {
		contextTitle += " " + value
	}

	contextTitle += " {"

	return contextTitle
}

func nonValueBraceHeadString(ctxType context_type.ContextType, _ string) string {
	return ctxType.String() + " {"
}

func nullHeadString(_ context_type.ContextType, _ string) string {
	return ""
}

func braceTailString() string {
	return "}"
}

func nullTailString() string {
	return ""
}

func directiveTailString() string {
	return ";"
}

func RegisterBuilder(opts context.BuildOptions, registrar ContextBuilderRegistrar) error {
	buildConfig, err := OptsApplyTo(opts)
	if err != nil {
		return err
	}
	builder := registrar(buildConfig.BasicContext)
	if builder == nil {
		return errors.New("the registered context builder is nil")
	}
	builderMap[opts.ContextType] = builder

	return nil
}

func RegisterParseFunc(opts parseFuncBuildOptions, parserFunc map[context_type.ContextType]parseFunc) error {
	parserFunc[opts.contextType] = func(data []byte, idx *int) context.Context {
		var value string
		if matchIndexes := opts.regex.FindIndex(data[*idx:]); matchIndexes != nil {
			if opts.valueMatchIndex >= 0 {
				value = string(opts.regex.FindSubmatch(data[*idx:])[opts.valueMatchIndex])
			}
			ctx := NewContext(opts.contextType, value)
			*idx += matchIndexes[len(matchIndexes)-1]

			return ctx
		}

		return context.NullContext()
	}

	return nil
}

func buildHeadAndTailStringFuncs(options context.BuildOptions) (func(context_type.ContextType, string) string, func() string) {
	var head func(context_type.ContextType, string) string
	var tail func() string
	switch options.ParseType {
	case context.ParseConfig:
		head = nullHeadString
		tail = nullTailString
	case context.ParseContext:
		if options.HasValue {
			head = hasValueBraceHeadString
		} else {
			head = nonValueBraceHeadString
		}
		tail = braceTailString
	case context.ParseDirective:
		head = nullHeadString
		tail = directiveTailString
	default:
		return nil, nil
	}

	return head, tail
}

func registerEventsBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeEvents,
			ParseType:   context.ParseContext,
			HasValue:    false,
		},
		func(f func() *BasicContext) func(value string) context.Context {
			return func(_ string) context.Context {
				ctx := f()
				ctx.self = ctx

				return ctx
			}
		},
	)
}

func registerEventsParseFunc() error {
	// convertor parse func
	convertorPushStackParseFuncMap[context_type.TypeEvents] = func(comment *Comment) bool {
		return !comment.Inline && RegEventsHead.MatchString(comment.Value())
	}

	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegEventsHead,
			contextType:     context_type.TypeEvents,
			valueMatchIndex: -1,
		},
		pushStackParseFuncMap,
	)
}

func registerGeoBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeGeo,
			ParseType:   context.ParseContext,
			HasValue:    true,
		},
		func(f func() *BasicContext) func(value string) context.Context {
			return func(value string) context.Context {
				ctx := f()
				ctx.ContextValue = value
				ctx.self = ctx

				return ctx
			}
		},
	)
}

func registerGeoParseFunc() error {
	// convertor parse func
	convertorPushStackParseFuncMap[context_type.TypeGeo] = func(comment *Comment) bool {
		return !comment.Inline && RegGeoHead.MatchString(comment.Value())
	}

	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegGeoHead,
			contextType:     context_type.TypeGeo,
			valueMatchIndex: 1,
		},
		pushStackParseFuncMap,
	)
}

func registerHttpBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeHttp,
			ParseType:   context.ParseContext,
			HasValue:    false,
		},
		func(f func() *BasicContext) func(value string) context.Context {
			return func(_ string) context.Context {
				ctx := f()
				ctx.self = ctx

				return ctx
			}
		},
	)
}

func registerHttpParseFunc() error {
	// convertor parse func
	convertorPushStackParseFuncMap[context_type.TypeHttp] = func(comment *Comment) bool {
		return !comment.Inline && RegHttpHead.MatchString(comment.Value())
	}

	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegHttpHead,
			contextType:     context_type.TypeHttp,
			valueMatchIndex: -1,
		},
		pushStackParseFuncMap,
	)
}

func registerIfBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeIf,
			ParseType:   context.ParseContext,
			HasValue:    true,
		},
		func(f func() *BasicContext) func(value string) context.Context {
			return func(value string) context.Context {
				ctx := f()
				ctx.ContextValue = value
				ctx.self = ctx

				return ctx
			}
		},
	)
}

func registerIfParseFunc() error {
	// convertor parse func
	convertorPushStackParseFuncMap[context_type.TypeIf] = func(comment *Comment) bool {
		return !comment.Inline && RegIfHead.MatchString(comment.Value())
	}

	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegIfHead,
			contextType:     context_type.TypeIf,
			valueMatchIndex: 1,
		},
		pushStackParseFuncMap,
	)
}

func registerLimitExceptBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeLimitExcept,
			ParseType:   context.ParseContext,
			HasValue:    true,
		},
		func(f func() *BasicContext) func(value string) context.Context {
			return func(value string) context.Context {
				ctx := f()
				ctx.ContextValue = value
				ctx.self = ctx

				return ctx
			}
		},
	)
}

func registerLimitExceptParseFunc() error {
	// convertor parse func
	convertorPushStackParseFuncMap[context_type.TypeLimitExcept] = func(comment *Comment) bool {
		return !comment.Inline && RegLimitExceptHead.MatchString(comment.Value())
	}

	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegLimitExceptHead,
			contextType:     context_type.TypeLimitExcept,
			valueMatchIndex: 1,
		},
		pushStackParseFuncMap,
	)
}

func registerLocationBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeLocation,
			ParseType:   context.ParseContext,
			HasValue:    true,
		},
		func(f func() *BasicContext) func(value string) context.Context {
			return func(value string) context.Context {
				ctx := f()
				ctx.ContextValue = value
				ctx.self = ctx

				return ctx
			}
		},
	)
}

func registerLocationParseFunc() error {
	// convertor parse func
	convertorPushStackParseFuncMap[context_type.TypeLocation] = func(comment *Comment) bool {
		return !comment.Inline && RegLocationHead.MatchString(comment.Value())
	}

	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegLocationHead,
			contextType:     context_type.TypeLocation,
			valueMatchIndex: 1,
		},
		pushStackParseFuncMap,
	)
}

func registerMapBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeMap,
			ParseType:   context.ParseContext,
			HasValue:    true,
		},
		func(f func() *BasicContext) func(value string) context.Context {
			return func(value string) context.Context {
				ctx := f()
				ctx.ContextValue = value
				ctx.self = ctx

				return ctx
			}
		},
	)
}

func registerMapParseFunc() error {
	// convertor parse func
	convertorPushStackParseFuncMap[context_type.TypeMap] = func(comment *Comment) bool {
		return !comment.Inline && RegMapHead.MatchString(comment.Value())
	}

	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegMapHead,
			contextType:     context_type.TypeMap,
			valueMatchIndex: 1,
		},
		pushStackParseFuncMap,
	)
}

func registerServerBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeServer,
			ParseType:   context.ParseContext,
			HasValue:    false,
		},
		func(f func() *BasicContext) func(value string) context.Context {
			return func(_ string) context.Context {
				ctx := f()
				ctx.self = ctx

				return ctx
			}
		},
	)
}

func registerServerParseFunc() error {
	// convertor parse func
	convertorPushStackParseFuncMap[context_type.TypeServer] = func(comment *Comment) bool {
		return !comment.Inline && RegServerHead.MatchString(comment.Value())
	}

	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegServerHead,
			contextType:     context_type.TypeServer,
			valueMatchIndex: -1,
		},
		pushStackParseFuncMap,
	)
}

func registerStreamBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeStream,
			ParseType:   context.ParseContext,
			HasValue:    false,
		},
		func(f func() *BasicContext) func(value string) context.Context {
			return func(_ string) context.Context {
				ctx := f()
				ctx.self = ctx

				return ctx
			}
		},
	)
}

func registerStreamParseFunc() error {
	// convertor parse func
	convertorPushStackParseFuncMap[context_type.TypeStream] = func(comment *Comment) bool {
		return !comment.Inline && RegStreamHead.MatchString(comment.Value())
	}

	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegStreamHead,
			contextType:     context_type.TypeStream,
			valueMatchIndex: -1,
		},
		pushStackParseFuncMap,
	)
}

func registerTypesBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeTypes,
			ParseType:   context.ParseContext,
			HasValue:    false,
		},
		func(f func() *BasicContext) func(value string) context.Context {
			return func(_ string) context.Context {
				ctx := f()
				ctx.self = ctx

				return ctx
			}
		},
	)
}

func registerTypesParseFunc() error {
	// convertor parse func
	convertorPushStackParseFuncMap[context_type.TypeTypes] = func(comment *Comment) bool {
		return !comment.Inline && RegTypesHead.MatchString(comment.Value())
	}

	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegTypesHead,
			contextType:     context_type.TypeTypes,
			valueMatchIndex: -1,
		},
		pushStackParseFuncMap,
	)
}

func registerUpstreamBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeUpstream,
			ParseType:   context.ParseContext,
			HasValue:    true,
		},
		func(f func() *BasicContext) func(value string) context.Context {
			return func(value string) context.Context {
				ctx := f()
				ctx.ContextValue = value
				ctx.self = ctx

				return ctx
			}
		},
	)
}

func registerUpstreamParseFunc() error {
	// convertor parse func
	convertorPushStackParseFuncMap[context_type.TypeUpstream] = func(comment *Comment) bool {
		return !comment.Inline && RegUpstreamHead.MatchString(comment.Value())
	}

	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegUpstreamHead,
			contextType:     context_type.TypeUpstream,
			valueMatchIndex: 1,
		},
		pushStackParseFuncMap,
	)
}

func NewContext(contextType context_type.ContextType, value string) context.Context {
	builder, has := builderMap[contextType]
	if !has {
		return context.ErrContext(errors.Errorf("not found context builder for %s", contextType))
	}

	return builder(value)
}

func registerContextBuilders() error {
	errs := make([]error, 0)
	errs = append(errs,
		registerCommentBuilder(),
		registerConfigBuilder(),
		registerDirectiveBuilder(),
		registerIncludeBuilder(),
		registerEventsBuilder(),
		registerGeoBuilder(),
		registerHttpBuilder(),
		registerIfBuilder(),
		registerLimitExceptBuilder(),
		registerLocationBuilder(),
		registerMapBuilder(),
		registerServerBuilder(),
		registerStreamBuilder(),
		registerTypesBuilder(),
		registerUpstreamBuilder(),
		registerTmpProxyPassBuilder(),
		registerHTTPProxyPassBuilder(),
		registerStreamProxyPassBuilder(),
	)

	return errors.NewAggregate(errs)
}

func registerContextParseFuncs() error {
	errs := make([]error, 0)
	errs = append(errs,
		registerEventsParseFunc(),
		registerGeoParseFunc(),
		registerHttpParseFunc(),
		registerIfParseFunc(),
		registerLimitExceptParseFunc(),
		registerLocationParseFunc(),
		registerMapParseFunc(),
		registerServerParseFunc(),
		registerStreamParseFunc(),
		registerTypesParseFunc(),
		registerUpstreamParseFunc(),
		registerDirectiveParseFunc(),
		registerCommentParseFunc(),
		registerProxyPassParseFunc(),
	)

	return errors.NewAggregate(errs)
}

func FatherPosSetWithoutInclude(ctx context.Context) context.PosSet {
	ps := ctx.FatherPosSet()
	if ps.Error() != nil {
		return ps
	}
	return ps.MapToPosSet(func(pos context.Pos) context.PosSet {
		if pos.Target().Type() == context_type.TypeInclude {
			return FatherPosSetWithoutInclude(pos.Target())
		}
		return context.NewPosSet().Append(pos)
	})
}
