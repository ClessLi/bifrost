package local

import (
	"encoding/json"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
	"regexp"
	"strings"
)

type JsonUnmarshalContext interface {
	Type() context_type.ContextType
	GetValue() string
	GetChildren() []*json.RawMessage
}

type jsonUnmarshalContext struct {
	Value       string             `json:"value,omitempty"`
	Children    []*json.RawMessage `json:"params,omitempty"`
	contextType context_type.ContextType
}

func (u jsonUnmarshalContext) Type() context_type.ContextType {
	return u.contextType
}

func (u jsonUnmarshalContext) GetValue() string {
	return u.Value
}

func (u jsonUnmarshalContext) GetChildren() []*json.RawMessage {
	return u.Children
}

type jsonUnmarshalMain struct {
	MainConfig string                        `json:"main-config,omitempty"`
	Configs    map[string][]*json.RawMessage `json:"configs,omitempty"`
}

type jsonUnmarshalComment struct {
	Comments string `json:"comments,omitempty"`
	Inline   bool   `json:"inline,omitempty"`
}

func (c jsonUnmarshalComment) Type() context_type.ContextType {
	if c.Inline {
		return context_type.TypeInlineComment
	}
	return context_type.TypeComment
}

func (c jsonUnmarshalComment) GetValue() string {
	return c.Comments
}

func (c jsonUnmarshalComment) GetChildren() []*json.RawMessage {
	return []*json.RawMessage{}
}

func registerCommentJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeComment, JsonUnmarshalRegCommentHead)
}

func registerCommentJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeComment, func() JsonUnmarshalContext {
		return new(jsonUnmarshalComment)
	})
}

type jsonUnmarshalConfig struct {
	jsonUnmarshalContext `json:"config"`
}

type jsonUnmarshalDirective struct {
	Name   string `json:"directive,omitempty"`
	Params string `json:"params,omitempty"`
}

func (d jsonUnmarshalDirective) Type() context_type.ContextType {
	return context_type.TypeDirective
}

func (d jsonUnmarshalDirective) GetValue() string {
	v := strings.TrimSpace(d.Name)
	if params := strings.TrimSpace(d.Params); len(params) > 0 {
		v += " " + params
	}
	return v
}

func (d jsonUnmarshalDirective) GetChildren() []*json.RawMessage {
	return []*json.RawMessage{}
}

func registerDirectiveJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeDirective, func() JsonUnmarshalContext {
		return new(jsonUnmarshalDirective)
	})
}

type jsonUnmarshalEvents struct {
	jsonUnmarshalContext `json:"events"`
}

func registerEventsJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeEvents, JsonUnmarshalRegEventsHead)
}

func registerEventsJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeEvents, func() JsonUnmarshalContext {
		u := new(jsonUnmarshalEvents)
		u.contextType = context_type.TypeEvents
		return u
	})
}

type jsonUnmarshalGeo struct {
	jsonUnmarshalContext `json:"geo"`
}

func registerGEOJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeGeo, JsonUnmarshalRegGeoHead)
}

func registerGEOJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeGeo, func() JsonUnmarshalContext {
		u := new(jsonUnmarshalGeo)
		u.contextType = context_type.TypeGeo
		return u
	})
}

type jsonUnmarshalHttp struct {
	jsonUnmarshalContext `json:"http"`
}

func registerHTTPJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeHttp, JsonUnmarshalRegHttpHead)
}

func registerHTTPJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeHttp, func() JsonUnmarshalContext {
		u := new(jsonUnmarshalHttp)
		u.contextType = context_type.TypeHttp
		return u
	})
}

type jsonUnmarshalIf struct {
	jsonUnmarshalContext `json:"if"`
}

func registerIFJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeIf, JsonUnmarshalRegIfHead)
}

func registerIFJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeIf, func() JsonUnmarshalContext {
		u := new(jsonUnmarshalIf)
		u.contextType = context_type.TypeIf
		return u
	})
}

type jsonUnmarshalInclude struct {
	jsonUnmarshalContext `json:"include"`
}

func registerIncludeJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeInclude, JsonUnmarshalRegIncludeHead)
}

func registerIncludeJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeInclude, func() JsonUnmarshalContext {
		u := new(jsonUnmarshalInclude)
		u.contextType = context_type.TypeInclude
		return u
	})
}

type jsonUnmarshalLimitExcept struct {
	jsonUnmarshalContext `json:"limit-except"`
}

func registerLimitExceptJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeLimitExcept, JsonUnmarshalRegLimitExceptHead)
}

func registerLimitExceptJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeLimitExcept, func() JsonUnmarshalContext {
		u := new(jsonUnmarshalLimitExcept)
		u.contextType = context_type.TypeLimitExcept
		return u
	})
}

type jsonUnmarshalLocation struct {
	jsonUnmarshalContext `json:"location"`
}

func registerLocationJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeLocation, JsonUnmarshalRegLocationHead)
}

func registerLocationJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeLocation, func() JsonUnmarshalContext {
		u := new(jsonUnmarshalLocation)
		u.contextType = context_type.TypeLocation
		return u
	})
}

type jsonUnmarshalMap struct {
	jsonUnmarshalContext `json:"map"`
}

func registerMapJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeMap, JsonUnmarshalRegMapHead)
}

func registerMapJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeMap, func() JsonUnmarshalContext {
		u := new(jsonUnmarshalMap)
		u.contextType = context_type.TypeMap
		return u
	})
}

type jsonUnmarshalServer struct {
	jsonUnmarshalContext `json:"server"`
}

func registerServerJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeServer, JsonUnmarshalRegServerHead)
}

func registerServerJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeServer, func() JsonUnmarshalContext {
		u := new(jsonUnmarshalServer)
		u.contextType = context_type.TypeServer
		return u
	})
}

type jsonUnmarshalStream struct {
	jsonUnmarshalContext `json:"stream"`
}

func registerStreamJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeStream, JsonUnmarshalRegStreamHead)
}

func registerStreamJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeStream, func() JsonUnmarshalContext {
		u := new(jsonUnmarshalStream)
		u.contextType = context_type.TypeStream
		return u
	})
}

type jsonUnmarshalTypes struct {
	jsonUnmarshalContext `json:"types"`
}

func registerTypesJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeTypes, JsonUnmarshalRegTypesHead)
}

func registerTypesJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeTypes, func() JsonUnmarshalContext {
		u := new(jsonUnmarshalTypes)
		u.contextType = context_type.TypeTypes
		return u
	})
}

type jsonUnmarshalUpstream struct {
	jsonUnmarshalContext `json:"upstream"`
}

func registerUpstreamJsonRegMatcher() error {
	return RegisterJsonRegMatcher(context_type.TypeUpstream, JsonUnmarshalRegUpstreamHead)
}

func registerUpstreamJsonUnmarshalerBuilder() error {
	return RegisterJsonUnmarshalerBuilder(context_type.TypeUpstream, func() JsonUnmarshalContext {
		u := new(jsonUnmarshalUpstream)
		u.contextType = context_type.TypeUpstream
		return u
	})
}

type mainUnmarshaler struct {
	unmarshalContext *jsonUnmarshalMain
	completedMain    MainContext
}

func (m *mainUnmarshaler) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, m.unmarshalContext)
	if err != nil {
		return err
	}
	main, err := NewMain(m.unmarshalContext.MainConfig)
	if err != nil {
		return err
	}

	m.completedMain = main

	toBeUnmarshalledConfigs := make(map[string]*jsonUnmarshalConfig)
	for value, rawMessages := range m.unmarshalContext.Configs {
		var configHashString string
		if value != m.unmarshalContext.MainConfig {
			config, ok := NewContext(context_type.TypeConfig, value).(*Config)
			if !ok {
				return errors.Errorf("failed to build config '%v'", value)
			}
			config.ConfigPath, err = newConfigPath(m.completedMain, value)
			if err != nil {
				return err
			}
			err = m.completedMain.AddConfig(config)
			if err != nil {
				return err
			}
			configHashString = configHash(config)
		} else {
			configHashString = configHash(m.completedMain.MainConfig())
		}
		toBeUnmarshalledConfigs[configHashString] = &jsonUnmarshalConfig{jsonUnmarshalContext{
			Value:       value,
			Children:    rawMessages,
			contextType: context_type.TypeConfig,
		}}
	}

	// unmarshal configs
	for configHashString, unmarshalConfig := range toBeUnmarshalledConfigs {
		cache, err := m.completedMain.GetConfig(configHashString)
		if err != nil {
			return err
		}
		unmarshaler := &jsonUnmarshaler{
			unmarshalContext: unmarshalConfig,
			configGraph:      m.completedMain.graph(),
			completedContext: cache,
		}
		for _, childRaw := range unmarshaler.unmarshalContext.GetChildren() {
			err = unmarshaler.nextUnmarshaler(childRaw).UnmarshalJSON(*childRaw)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type jsonUnmarshaler struct {
	unmarshalContext JsonUnmarshalContext
	configGraph      ConfigGraph
	completedContext context.Context
	fatherContext    context.Context
}

func (u *jsonUnmarshaler) UnmarshalJSON(bytes []byte) error {
	// unmarshal context, it's self
	err := json.Unmarshal(bytes, u.unmarshalContext)
	if err != nil {
		return err
	}

	switch u.unmarshalContext.Type() {
	case context_type.TypeConfig:
		return errors.New("invalid JSON data: the unmarshal operation for config should be completed in the unmarshal operation of the main unmarshaler")
	case context_type.TypeInclude:
		// insert the include context to be unmarshalled into its father, and unmarshal itself
		return u.unmarshalInclude()
	case context_type.TypeComment:
		return u.fatherContext.Insert(NewComment(u.unmarshalContext.GetValue(), false), u.fatherContext.Len()).Error()
	case context_type.TypeInlineComment:
		return u.fatherContext.Insert(NewComment(u.unmarshalContext.GetValue(), true), u.fatherContext.Len()).Error()
	case context_type.TypeDirective:
		return u.fatherContext.Insert(NewDirective(u.unmarshalContext.(*jsonUnmarshalDirective).Name, u.unmarshalContext.(*jsonUnmarshalDirective).Params), u.fatherContext.Len()).Error()
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

func (u *jsonUnmarshaler) nextUnmarshaler(message *json.RawMessage) *jsonUnmarshaler {
	matchedType := context_type.TypeDirective
	for contextType, matcher := range jsonUnmarshalRegMatcherMap {
		if matcher(*message) {
			matchedType = contextType
			break
		}
	}

	return jsonUnmarshalerBuilderMap[matchedType](u.configGraph, u.completedContext)
}

func (u *jsonUnmarshaler) unmarshalInclude() error {
	unmarshalInclude, ok := u.unmarshalContext.(*jsonUnmarshalInclude)
	if !ok {
		return errors.New("unmarshal context is not jsonUnmarshalInclude")
	}
	u.completedContext = NewContext(context_type.TypeInclude, unmarshalInclude.GetValue())
	err := u.completedContext.Error()
	if err != nil {
		return err
	}

	// insert the Include context to be unmarshalled into its father
	if err = u.fatherContext.Insert(u.completedContext, u.fatherContext.Len()).Error(); err != nil {
		return err
	}

	// unmarshal included configs
	configs := make([]*Config, 0)
	for _, childRaw := range unmarshalInclude.GetChildren() {
		var path string
		err := json.Unmarshal(*childRaw, &path)
		if err != nil {
			return err
		}
		configPath, err := newConfigPath(u.configGraph, path)
		if err != nil {
			return err
		}

		// get config cache
		cache, err := u.configGraph.GetConfig(strings.TrimSpace(configPath.FullPath()))
		if err != nil {
			return err
		}
		configs = append(configs, cache)

	}
	// include configs
	return u.completedContext.(*Include).InsertConfig(configs...)
}

func RegisterJsonRegMatcher(contextType context_type.ContextType, regexp *regexp.Regexp) error {
	jsonUnmarshalRegMatcherMap[contextType] = func(jsonraw []byte) bool {
		return regexp.Find(jsonraw) != nil
	}
	return nil
}

func RegisterJsonUnmarshalerBuilder(contextType context_type.ContextType, newFunc func() JsonUnmarshalContext) error {
	jsonUnmarshalerBuilderMap[contextType] = func(graph ConfigGraph, father context.Context) *jsonUnmarshaler {
		return &jsonUnmarshaler{
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
