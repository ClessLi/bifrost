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
	config := BuildBasicContextConfig{}
	config.headStringFunc, config.tailStringFunc = buildHeadAndTailStringFuncs(opts)
	if config.headStringFunc == nil || config.tailStringFunc == nil {
		return config, errors.New("TODO")
	}

	config.ContextType = opts.ContextType
	return config, nil
}

func (b BuildBasicContextConfig) BasicContext() BasicContext {
	return BasicContext{
		ContextType:    b.ContextType,
		Children:       make([]context.Context, 0),
		father:         context.NullContext(),
		headStringFunc: b.headStringFunc,
		tailStringFunc: b.tailStringFunc,
	}
}

type ContextBuilderRegistrar func(func() BasicContext) func(value string) context.Context

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
	config, err := OptsApplyTo(opts)
	if err != nil {
		return err
	}
	builder := registrar(config.BasicContext)
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

type Main struct {
	*Config `json:"main"`
}

func registerMainBuilder() error {
	builderMap[context_type.TypeMain] = func(value string) context.Context {
		mainConfig, err := newMainConfig(value)
		if err != nil {
			return context.ErrContext(err)
		}
		main := &Main{mainConfig}
		return main
	}
	return nil
}

type Events struct {
	BasicContext `json:"events"`
}

func registerEventsBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeEvents,
			ParseType:   context.ParseContext,
			HasValue:    false,
		},
		func(f func() BasicContext) func(value string) context.Context {
			return func(_ string) context.Context {
				ctx := &Events{f()}
				ctx.self = ctx
				return ctx
			}
		},
	)
}

func registerEventsParseFunc() error {
	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegEventsHead,
			contextType:     context_type.TypeEvents,
			valueMatchIndex: -1,
		},
		pushStackParseFuncMap,
	)
}

type Geo struct {
	BasicContext `json:"geo"`
}

func registerGeoBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeGeo,
			ParseType:   context.ParseContext,
			HasValue:    true,
		},
		func(f func() BasicContext) func(value string) context.Context {
			return func(value string) context.Context {
				ctx := &Geo{f()}
				ctx.ContextValue = value
				ctx.self = ctx
				return ctx
			}
		},
	)
}

func registerGeoParseFunc() error {
	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegGeoHead,
			contextType:     context_type.TypeGeo,
			valueMatchIndex: 1,
		},
		pushStackParseFuncMap,
	)
}

type Http struct {
	BasicContext `json:"http"`
}

func registerHttpBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeHttp,
			ParseType:   context.ParseContext,
			HasValue:    false,
		},
		func(f func() BasicContext) func(value string) context.Context {
			return func(_ string) context.Context {
				ctx := &Http{f()}
				ctx.self = ctx
				return ctx
			}
		},
	)
}

func registerHttpParseFunc() error {
	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegHttpHead,
			contextType:     context_type.TypeHttp,
			valueMatchIndex: -1,
		},
		pushStackParseFuncMap,
	)
}

type If struct {
	BasicContext `json:"if"`
}

func registerIfBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeIf,
			ParseType:   context.ParseContext,
			HasValue:    true,
		},
		func(f func() BasicContext) func(value string) context.Context {
			return func(value string) context.Context {
				ctx := &If{f()}
				ctx.ContextValue = value
				ctx.self = ctx
				return ctx
			}
		},
	)
}

func registerIfParseFunc() error {
	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegIfHead,
			contextType:     context_type.TypeIf,
			valueMatchIndex: 1,
		},
		pushStackParseFuncMap,
	)
}

type LimitExcept struct {
	BasicContext `json:"limit_except"`
}

func registerLimitExceptBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeLimitExcept,
			ParseType:   context.ParseContext,
			HasValue:    true,
		},
		func(f func() BasicContext) func(value string) context.Context {
			return func(value string) context.Context {
				ctx := &LimitExcept{f()}
				ctx.ContextValue = value
				ctx.self = ctx
				return ctx
			}
		},
	)
}

func registerLimitExceptParseFunc() error {
	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegLimitExceptHead,
			contextType:     context_type.TypeLimitExcept,
			valueMatchIndex: 1,
		},
		pushStackParseFuncMap,
	)
}

type Location struct {
	BasicContext `json:"location"`
}

func registerLocationBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeLocation,
			ParseType:   context.ParseContext,
			HasValue:    true,
		},
		func(f func() BasicContext) func(value string) context.Context {
			return func(value string) context.Context {
				ctx := &Location{f()}
				ctx.ContextValue = value
				ctx.self = ctx
				return ctx
			}
		},
	)
}

func registerLocationParseFunc() error {
	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegLocationHead,
			contextType:     context_type.TypeLocation,
			valueMatchIndex: 1,
		},
		pushStackParseFuncMap,
	)
}

type Map struct {
	BasicContext `json:"map"`
}

func registerMapBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeMap,
			ParseType:   context.ParseContext,
			HasValue:    true,
		},
		func(f func() BasicContext) func(value string) context.Context {
			return func(value string) context.Context {
				ctx := &Map{f()}
				ctx.ContextValue = value
				ctx.self = ctx
				return ctx
			}
		},
	)
}

func registerMapParseFunc() error {
	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegMapHead,
			contextType:     context_type.TypeMap,
			valueMatchIndex: 1,
		},
		pushStackParseFuncMap,
	)
}

type Server struct {
	BasicContext `json:"server"`
}

func registerServerBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeServer,
			ParseType:   context.ParseContext,
			HasValue:    false,
		},
		func(f func() BasicContext) func(value string) context.Context {
			return func(_ string) context.Context {
				ctx := &Server{f()}
				ctx.self = ctx
				return ctx
			}
		},
	)
}

func registerServerParseFunc() error {
	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegServerHead,
			contextType:     context_type.TypeServer,
			valueMatchIndex: -1,
		},
		pushStackParseFuncMap,
	)
}

type Stream struct {
	BasicContext `json:"stream"`
}

func registerStreamBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeStream,
			ParseType:   context.ParseContext,
			HasValue:    false,
		},
		func(f func() BasicContext) func(value string) context.Context {
			return func(_ string) context.Context {
				ctx := &Stream{f()}
				ctx.self = ctx
				return ctx
			}
		},
	)
}

func registerStreamParseFunc() error {
	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegStreamHead,
			contextType:     context_type.TypeStream,
			valueMatchIndex: -1,
		},
		pushStackParseFuncMap,
	)
}

type Types struct {
	BasicContext `json:"types"`
}

func registerTypesBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeTypes,
			ParseType:   context.ParseContext,
			HasValue:    false,
		},
		func(f func() BasicContext) func(value string) context.Context {
			return func(_ string) context.Context {
				ctx := &Types{f()}
				ctx.self = ctx
				return ctx
			}
		},
	)
}

func registerTypesParseFunc() error {
	return RegisterParseFunc(
		parseFuncBuildOptions{
			regex:           RegTypesHead,
			contextType:     context_type.TypeTypes,
			valueMatchIndex: -1,
		},
		pushStackParseFuncMap,
	)
}

type Upstream struct {
	BasicContext `json:"upstream"`
}

func registerUpstreamBuilder() error {
	return RegisterBuilder(
		context.BuildOptions{
			ContextType: context_type.TypeUpstream,
			ParseType:   context.ParseContext,
			HasValue:    true,
		},
		func(f func() BasicContext) func(value string) context.Context {
			return func(value string) context.Context {
				ctx := &Upstream{f()}
				ctx.ContextValue = value
				ctx.self = ctx
				return ctx
			}
		},
	)
}

func registerUpstreamParseFunc() error {
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
		registerConfigBuilder(),
		registerMainBuilder(),
		registerIncludeBuild(),
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
	)
	return errors.NewAggregate(errs)
}

func registerContextParseFuncs() error {
	errs := make([]error, 0)
	errs = append(errs,
		registerIncludeParseFunc(),
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
	)
	return errors.NewAggregate(errs)
}
