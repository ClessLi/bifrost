package local

import (
	"encoding/json"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
	"strings"
)

type jsonUnmarshalContext struct {
	Enabled     bool                     `json:"enabled,omitempty"`
	ContextType context_type.ContextType `json:"context-type"`
	Value       string                   `json:"value,omitempty"`
	Children    []*json.RawMessage       `json:"params,omitempty"`
}

type jsonUnmarshalMain struct {
	MainConfig string                      `json:"main-config,omitempty"`
	Configs    map[string]*json.RawMessage `json:"configs,omitempty"`
}

type mainUnmarshaler struct {
	completedMain MainContext
}

func (m *mainUnmarshaler) UnmarshalJSON(bytes []byte) error {
	var unmarshalMain jsonUnmarshalMain
	err := json.Unmarshal(bytes, &unmarshalMain)
	if err != nil {
		return err
	}
	main, err := NewMain(unmarshalMain.MainConfig)
	if err != nil {
		return err
	}

	m.completedMain = main

	// add configs into graph, and unmarshal configs
	for configpath, configRaw := range unmarshalMain.Configs {
		//var configHashString string
		if configpath != unmarshalMain.MainConfig {
			config, ok := NewContext(context_type.TypeConfig, configpath).(*Config)
			if !ok {
				return errors.Errorf("failed to build config '%v'", configpath)
			}
			config.ConfigPath, err = newConfigPath(m.completedMain, configpath)
			if err != nil {
				return err
			}
			err = m.completedMain.addVertex(config)
			if err != nil {
				return err
			}
		}
		configMarshaler := jsonUnmarshaler{
			configGraph:      m.completedMain.graph(),
			completedContext: context.NullContext(),
			fatherContext:    m.completedMain,
		}
		err = configMarshaler.UnmarshalJSON(*configRaw)
		if err != nil {
			return err
		}
	}

	return main.rerenderGraph()
}

type jsonUnmarshaler struct {
	configGraph      ConfigGraph
	completedContext context.Context
	fatherContext    context.Context
}

func (u *jsonUnmarshaler) UnmarshalJSON(bytes []byte) error {
	var unmarshalCtx jsonUnmarshalContext
	// unmarshal context, it's self
	err := json.Unmarshal(bytes, &unmarshalCtx)
	if err != nil {
		return err
	}

	switch unmarshalCtx.ContextType {
	case context_type.TypeConfig:
		err = u.unmarshalConfig(&unmarshalCtx)
		if err != nil {
			return err
		}
	//case context_type.TypeInclude:
	//	// insert the include context to be unmarshalled into its father, and unmarshal itself
	//	return u.unmarshalInclude(&unmarshalCtx)
	default:
		u.completedContext = NewContext(unmarshalCtx.ContextType, unmarshalCtx.Value)
		if err = u.completedContext.Error(); err != nil {
			return err
		}

		// insert the context to be unmarshalled into its father
		if err = u.fatherContext.Insert(u.completedContext, u.fatherContext.Len()).Error(); err != nil {
			return err
		}
	}

	// unmarshal context's children
	for _, childRaw := range unmarshalCtx.Children {
		err = u.nextUnmarshaler().UnmarshalJSON(*childRaw)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *jsonUnmarshaler) nextUnmarshaler() *jsonUnmarshaler {
	return &jsonUnmarshaler{
		configGraph:      u.configGraph,
		completedContext: context.NullContext(),
		fatherContext:    u.completedContext,
	}
}

func (u *jsonUnmarshaler) unmarshalConfig(unmashalctx *jsonUnmarshalContext) error {
	configPath, err := newConfigPath(u.configGraph, unmashalctx.Value)
	if err != nil {
		return err
	}
	cache, err := u.configGraph.GetConfig(strings.TrimSpace(configPath.FullPath()))
	if err != nil {
		return err
	}
	u.completedContext = cache
	return nil
}

//func (u *jsonUnmarshaler) unmarshalInclude(unmarshalctx *jsonUnmarshalContext) error {
//	u.completedContext = NewContext(unmarshalctx.ContextType, unmarshalctx.Value)
//	err := u.completedContext.Error()
//	if err != nil {
//		return err
//	}
//	if unmarshalctx.Enabled {
//		u.completedContext.Enable()
//	} else {
//		u.completedContext.Disable()
//	}
//
//	// insert the Include context to be unmarshalled into its father
//	if err = u.fatherContext.Insert(u.completedContext, u.fatherContext.Len()).Error(); err != nil {
//		return err
//	}
//
//	// unmarshal included configs
//	configs := make([]*Config, 0)
//	for _, childRaw := range unmarshalctx.Children {
//		var path string
//		err := json.Unmarshal(*childRaw, &path)
//		if err != nil {
//			return err
//		}
//		configPath, err := newConfigPath(u.configGraph, path)
//		if err != nil {
//			return err
//		}
//
//		// get config cache
//		cache, err := u.configGraph.GetConfig(strings.TrimSpace(configPath.FullPath()))
//		if err != nil {
//			return err
//		}
//		configs = append(configs, cache)
//
//	}
//	// include configs
//	return u.completedContext.(*Include).InsertConfig(configs...)
//}
