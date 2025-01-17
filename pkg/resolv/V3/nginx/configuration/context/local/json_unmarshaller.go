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

type mainUnmarshaller struct {
	completedMain MainContext
}

func (m *mainUnmarshaller) UnmarshalJSON(bytes []byte) error {
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
		configUnmarshaller := jsonUnmarshaller{
			configGraph:      m.completedMain.graph(),
			completedContext: context.NullContext(),
			fatherContext:    m.completedMain,
		}
		err = configUnmarshaller.UnmarshalJSON(*configRaw)
		if err != nil {
			return err
		}
	}

	return main.rerenderGraph()
}

type jsonUnmarshaller struct {
	configGraph      ConfigGraph
	completedContext context.Context
	fatherContext    context.Context
}

func (u *jsonUnmarshaller) UnmarshalJSON(bytes []byte) error {
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

	// set enabled/disabled to completed context
	if unmarshalCtx.Enabled {
		err = u.completedContext.Enable().Error()
	} else {
		err = u.completedContext.Disable().Error()
	}
	if err != nil {
		return err
	}

	// unmarshal context's children
	for _, childRaw := range unmarshalCtx.Children {
		err = u.nextUnmarshaller().UnmarshalJSON(*childRaw)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *jsonUnmarshaller) nextUnmarshaller() *jsonUnmarshaller {
	return &jsonUnmarshaller{
		configGraph:      u.configGraph,
		completedContext: context.NullContext(),
		fatherContext:    u.completedContext,
	}
}

func (u *jsonUnmarshaller) unmarshalConfig(unmashalctx *jsonUnmarshalContext) error {
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
