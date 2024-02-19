package local

import (
	"encoding/json"
	"github.com/ClessLi/bifrost/pkg/graph"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
	"path/filepath"
	"regexp"
	"strings"
)

type UnmarshalContext interface {
	Type() context_type.ContextType
	GetValue() string
	GetChildren() []*json.RawMessage
}

type unmarshalIncludeContext struct {
	Value    string                      `json:"value,omitempty"`
	Children map[string]*json.RawMessage `json:"param,omitempty"`
}

func (u unmarshalIncludeContext) Type() context_type.ContextType {
	return context_type.TypeInclude
}

func (u unmarshalIncludeContext) GetValue() string {
	return u.Value
}

func (u unmarshalIncludeContext) GetChildren() []*json.RawMessage {
	return nil
}

type unmarshalContext struct {
	Value       string             `json:"value,omitempty"`
	Children    []*json.RawMessage `json:"param,omitempty"`
	contextType context_type.ContextType
}

func (u unmarshalContext) Type() context_type.ContextType {
	return u.contextType
}

func (u unmarshalContext) GetValue() string {
	return u.Value
}

func (u unmarshalContext) GetChildren() []*json.RawMessage {
	return u.Children
}

type _main struct {
	config `json:"main"`
}

type comment struct {
	Comments string `json:"comments,omitempty"`
	Inline   bool   `json:"inline,omitempty"`
}

func (c comment) Type() context_type.ContextType {
	if c.Inline {
		return context_type.TypeInlineComment
	}
	return context_type.TypeComment
}

func (c comment) GetValue() string {
	return c.Comments
}

func (c comment) GetChildren() []*json.RawMessage {
	return []*json.RawMessage{}
}

func registerCommentJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeComment, JsonUnmarshalRegCommentHead)
}

func registerCommentJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeComment, func() UnmarshalContext {
		return new(comment)
	})
}

type config struct {
	unmarshalContext `json:"config"`
}

type directive struct {
	Name   string `json:"name,omitempty"`
	Params string `json:"params,omitempty"`
}

func (d directive) Type() context_type.ContextType {
	return context_type.TypeDirective
}

func (d directive) GetValue() string {
	return strings.Join([]string{d.Name, d.Params}, " ")
}

func (d directive) GetChildren() []*json.RawMessage {
	return []*json.RawMessage{}
}

func registerDirectiveJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeDirective, func() UnmarshalContext {
		return new(directive)
	})
}

type events struct {
	unmarshalContext `json:"events"`
}

func registerEventsJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeEvents, JsonUnmarshalRegEventsHead)
}

func registerEventsJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeEvents, func() UnmarshalContext {
		return new(events)
	})
}

type geo struct {
	unmarshalContext `json:"geo"`
}

func registerGEOJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeGeo, JsonUnmarshalRegGeoHead)
}

func registerGEOJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeGeo, func() UnmarshalContext {
		return new(geo)
	})
}

type http struct {
	unmarshalContext `json:"http"`
}

func registerHTTPJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeHttp, JsonUnmarshalRegHttpHead)
}

func registerHTTPJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeHttp, func() UnmarshalContext {
		return new(http)
	})
}

type _if struct {
	unmarshalContext `json:"if"`
}

func registerIFJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeIf, JsonUnmarshalRegIfHead)
}

func registerIFJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeIf, func() UnmarshalContext {
		return new(_if)
	})
}

type include struct {
	unmarshalIncludeContext `json:"include"`
}

func registerIncludeJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeInclude, JsonUnmarshalRegIncludeHead)
}

func registerIncludeJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeInclude, func() UnmarshalContext {
		return new(include)
	})
}

type limitExcept struct {
	unmarshalContext `json:"limit_except"`
}

func registerLimitExceptJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeLimitExcept, JsonUnmarshalRegLimitExceptHead)
}

func registerLimitExceptJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeLimitExcept, func() UnmarshalContext {
		return new(limitExcept)
	})
}

type location struct {
	unmarshalContext `json:"location"`
}

func registerLocationJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeLocation, JsonUnmarshalRegLocationHead)
}

func registerLocationJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeLocation, func() UnmarshalContext {
		return new(location)
	})
}

type _map struct {
	unmarshalContext `json:"map"`
}

func registerMapJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeMap, JsonUnmarshalRegMapHead)
}

func registerMapJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeMap, func() UnmarshalContext {
		return new(_map)
	})
}

type server struct {
	unmarshalContext `json:"server"`
}

func registerServerJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeServer, JsonUnmarshalRegServerHead)
}

func registerServerJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeServer, func() UnmarshalContext {
		return new(server)
	})
}

type stream struct {
	unmarshalContext `json:"stream"`
}

func registerStreamJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeStream, JsonUnmarshalRegStreamHead)
}

func registerStreamJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeStream, func() UnmarshalContext {
		return new(stream)
	})
}

type types struct {
	unmarshalContext `json:"types"`
}

func registerTypesJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeTypes, JsonUnmarshalRegTypesHead)
}

func registerTypesJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeTypes, func() UnmarshalContext {
		return new(types)
	})
}

type upstream struct {
	unmarshalContext `json:"upstream"`
}

func registerUpstreamJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeUpstream, JsonUnmarshalRegUpstreamHead)
}

func registerUpstreamJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeUpstream, func() UnmarshalContext {
		return new(upstream)
	})
}

type mainUnmarshaler struct {
	unmarshalContext *_main
	completedMain    *Main
}

func (m *mainUnmarshaler) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, m.unmarshalContext)
	if err != nil {
		return err
	}
	mainCtx := NewContext(context_type.TypeMain, m.unmarshalContext.GetValue())
	if err = mainCtx.Error(); err != nil {
		return err
	}
	main, ok := mainCtx.(*Main)
	if !ok {
		return errors.New("failed to build main context")
	}

	m.completedMain = main

	mainConfigUnmarshaler := &unmarshaler{
		unmarshalContext: m.unmarshalContext.config,
		configGraph:      m.completedMain.ConfigGraph,
		completedContext: m.completedMain.Config,
		fatherContext:    context.NullContext(),
	}

	for _, childRaw := range mainConfigUnmarshaler.unmarshalContext.GetChildren() {
		err = mainConfigUnmarshaler.nextUnmarshaler(childRaw).UnmarshalJSON(*childRaw)
		if err != nil {
			return err
		}
	}
	return nil
}

type unmarshaler struct {
	unmarshalContext UnmarshalContext
	configGraph      ConfigGraph
	completedContext context.Context
	fatherContext    context.Context
}

func (u *unmarshaler) UnmarshalJSON(bytes []byte) error {
	// unmarshal context, it's self
	err := json.Unmarshal(bytes, u.unmarshalContext)
	if err != nil {
		return err
	}

	switch u.unmarshalContext.Type() {
	case context_type.TypeConfig:
		// pre unmarshal config context from include unmarshaler
		if err = u.completedContext.Error(); err != nil {
			return err
		}
		// check config path
		if u.completedContext.(*Config).FullPath() != u.unmarshalContext.GetValue() {
			return errors.Errorf("config paths are different between completed config(%s) and unmarshalled config(%s)", u.completedContext.(*Config).FullPath(), u.unmarshalContext.GetValue())
		}
	case context_type.TypeInclude:
		// insert the include context to be unmarshalled into its father, and unmarshal itself
		return u.unmarshalInclude()
	case context_type.TypeComment:
		return u.fatherContext.Insert(NewComment(u.unmarshalContext.GetValue(), false), u.fatherContext.Len()).Error()
	case context_type.TypeInlineComment:
		return u.fatherContext.Insert(NewComment(u.unmarshalContext.GetValue(), true), u.fatherContext.Len()).Error()
	case context_type.TypeDirective:
		return u.fatherContext.Insert(NewDirective(u.unmarshalContext.(*directive).Name, u.unmarshalContext.(*directive).Params), u.fatherContext.Len()).Error()
	default:
		u.completedContext = NewContext(u.unmarshalContext.Type(), u.unmarshalContext.GetValue())
		if err = u.completedContext.Error(); err != nil {
			return err
		}

		// insert the context to be unmarshalled into its father
		if err = u.fatherContext.Insert(u.completedContext, u.fatherContext.Len()).Error(); err != nil {
			return err
		}
	}

	// unmarshal context's children
	for _, childRaw := range u.unmarshalContext.GetChildren() {
		err = u.nextUnmarshaler(childRaw).UnmarshalJSON(*childRaw)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *unmarshaler) nextUnmarshaler(message *json.RawMessage) *unmarshaler {
	matchedType := context_type.TypeDirective
	for contextType, matcher := range jsonUnmarshalRegMatcherMap {
		if matcher(*message) {
			matchedType = contextType
			break
		}
	}

	return jsonUnmarshalerBuilderMap[matchedType](u.configGraph, u.completedContext)
}

func (u *unmarshaler) unmarshalInclude() error {
	unmarshalIncludeCtx, ok := u.unmarshalContext.(*unmarshalIncludeContext)
	if !ok {
		return errors.New("unmarshal context is not unmarshalIncludeContext")
	}
	u.completedContext = NewContext(context_type.TypeInclude, unmarshalIncludeCtx.GetValue())
	err := u.completedContext.Error()
	if err != nil {
		return err
	}

	isAbsInclude := filepath.IsAbs(u.completedContext.Value())
	// unmarshal included configs
	// new configs
	newconfigs := make([]*Config, 0)
	includedconfigs := make([]*Config, 0)
	for path := range unmarshalIncludeCtx.Children {
		var configpath context.ConfigPath
		if isAbsInclude {
			configpath, err = context.NewAbsConfigPath(path)
		} else {
			configpath, err = context.NewRelConfigPath(u.configGraph.MainConfig().BaseDir(), path)
		}
		if err != nil {
			return err
		}

		// get config cache
		cache, err := u.configGraph.GetConfig(configpath.FullPath())
		if err == nil { // has cache
			includedconfigs = append(includedconfigs, cache)
			continue
		} else if !errors.Is(err, graph.ErrVertexNotExist) {
			return err
		}

		// build new config
		config, ok := NewContext(context_type.TypeConfig, path).(*Config)
		if !ok {
			return errors.Errorf("failed to build included config %s", path)
		}
		config.ConfigPath = configpath
		newconfigs = append(newconfigs, config)
	}
	includedconfigs = append(includedconfigs, newconfigs...)
	// include configs
	err = u.completedContext.(*Include).InsertConfig(includedconfigs...)
	if err != nil {
		return err
	}
	// unmarshal new configs
	for _, c := range newconfigs {
		configUnmarshaler := &unmarshaler{
			unmarshalContext: new(config),
			configGraph:      u.configGraph,
			completedContext: c,
			fatherContext:    context.NullContext(),
		}
		err = configUnmarshaler.UnmarshalJSON(*unmarshalIncludeCtx.Children[c.FullPath()])
		if err != nil {
			return err
		}
	}
	return nil
}

func RegisterJsonRegMatcher(contextType context_type.ContextType, regexp *regexp.Regexp) error {
	jsonUnmarshalRegMatcherMap[contextType] = func(jsonraw []byte) bool {
		return regexp.Find(jsonraw) != nil
	}
	return nil
}

func RegisterJsonUnmarshalerBuilder(contextType context_type.ContextType, newFunc func() UnmarshalContext) error {
	jsonUnmarshalerBuilderMap[contextType] = func(graph ConfigGraph, father context.Context) *unmarshaler {
		return &unmarshaler{
			unmarshalContext: newFunc(),
			configGraph:      graph,
			completedContext: context.NullContext(),
			fatherContext:    father,
		}
	}
	return nil
}

func registerJsonRegMatchers() error {
	errs := make([]error, 0)
	errs = append(errs,
		registerCommentJsonRegMatcher(),
		registerEventsJsonRegMatcher(),
		registerGEOJsonRegMatcher(),
		registerHTTPJsonRegMatcher(),
		registerIFJsonRegMatcher(),
		registerIncludeJsonRegMatcher(),
		registerLimitExceptJsonRegMatcher(),
		registerLocationJsonRegMatcher(),
		registerMapJsonRegMatcher(),
		registerServerJsonRegMatcher(),
		registerStreamJsonRegMatcher(),
		registerTypesJsonRegMatcher(),
		registerUpstreamJsonRegMatcher(),
	)
	return errors.NewAggregate(errs)
}

func registerJsonUnmarshalerBuilders() error {
	errs := make([]error, 0)
	errs = append(errs,
		registerCommentJsonUnmarshalerBuilder(),
		registerDirectiveJsonUnmarshalerBuilder(),
		registerEventsJsonUnmarshalerBuilder(),
		registerGEOJsonUnmarshalerBuilder(),
		registerHTTPJsonUnmarshalerBuilder(),
		registerIFJsonUnmarshalerBuilder(),
		registerIncludeJsonUnmarshalerBuilder(),
		registerLimitExceptJsonUnmarshalerBuilder(),
		registerLocationJsonUnmarshalerBuilder(),
		registerMapJsonUnmarshalerBuilder(),
		registerServerJsonUnmarshalerBuilder(),
		registerStreamJsonUnmarshalerBuilder(),
		registerTypesJsonUnmarshalerBuilder(),
		registerUpstreamJsonUnmarshalerBuilder(),
	)
	return errors.NewAggregate(errs)
}
