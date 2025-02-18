package local

import (
	"encoding/json"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
)

type MainContext interface {
	context.Context
	ConfigGraph

	graph() ConfigGraph
}

type Main struct {
	ConfigGraph
}

func (m *Main) graph() ConfigGraph {
	return m.ConfigGraph
}

func (m *Main) MarshalJSON() ([]byte, error) {
	if m.ConfigGraph == nil || m.ConfigGraph.MainConfig() == nil || !m.ConfigGraph.MainConfig().isInGraph() {
		return nil, errors.New("MainContext is not completed")
	}
	marshalCtx := struct {
		MainConfig string             `json:"main-config"`
		Configs    map[string]*Config `json:"configs"`
	}{
		MainConfig: m.MainConfig().Value(),
		Configs:    make(map[string]*Config),
	}

	for _, config := range m.ConfigGraph.ListConfigs() {
		marshalCtx.Configs[config.Value()] = config
	}

	return json.Marshal(marshalCtx)
}

func (m *Main) Father() context.Context {
	return m
}

func (m *Main) Child(idx int) context.Context {
	return m.MainConfig().Child(idx)
}

func (m *Main) ChildrenPosSet() context.PosSet {
	return m.MainConfig().ChildrenPosSet().Map(
		func(pos context.Pos) (context.Pos, error) {
			if f, idx := pos.Position(); f == m.MainConfig().self {
				return context.SetPos(m, idx), nil
			}
			return pos, nil
		},
	)
}

func (m *Main) SetValue(v string) error {
	if len(m.ListConfigs()) > 1 {
		return errors.New("cannot set value for MainContext with non empty graph")
	}
	return m.MainConfig().SetValue(v)
}

func (m *Main) SetFather(ctx context.Context) error {
	return errors.New("cannot set father for MainContext")
}

func (m *Main) HasChild() bool {
	return m.MainConfig().HasChild()
}

func (m *Main) Len() int {
	return m.MainConfig().Len()
}

func (m *Main) Value() string {
	return m.MainConfig().Value()
}

func (m *Main) Error() error {
	return nil
}

func (m *Main) ConfigLines(isDumping bool) ([]string, error) {
	return m.MainConfig().ConfigLines(isDumping)
}

func (m *Main) Insert(ctx context.Context, idx int) context.Context {
	if got := m.MainConfig().Insert(ctx, idx); got == m.MainConfig().self {
		return m
	} else {
		return got
	}
}

func (m *Main) Remove(idx int) context.Context {
	if got := m.MainConfig().Remove(idx); got == m.MainConfig().self {
		return m
	} else {
		return got
	}
}

func (m *Main) Modify(ctx context.Context, idx int) context.Context {
	if got := m.MainConfig().Modify(ctx, idx); got == m.MainConfig().self {
		return m
	} else {
		return got
	}
}

func (m *Main) Clone() context.Context {
	cloneConfigPath, err := context.NewAbsConfigPath(m.Value())
	if err != nil {
		return context.ErrContext(err)
	}
	cloneConfig := m.MainConfig().Clone().(*Config)
	cloneConfig.ConfigPath = cloneConfigPath
	g, err := newConfigGraph(cloneConfig)
	if err != nil {
		return context.ErrContext(err)
	}
	clone := &Main{ConfigGraph: g}
	err = cloneConfig.SetFather(clone)
	if err != nil {
		return context.ErrContext(err)
	}
	return clone
}

func NewMain(abspath string) (MainContext, error) {
	confpath, err := context.NewAbsConfigPath(abspath)
	if err != nil {
		return nil, err
	}
	config := &Config{
		BasicContext: newBasicContext(context_type.TypeConfig, nullHeadString, nullTailString),
		ConfigPath:   confpath,
	}
	config.self = config
	config.ContextValue = abspath
	g, err := newConfigGraph(config)
	if err != nil { // Temporarily unable to be covered for testing
		return nil, err
	}
	m := &Main{ConfigGraph: g}
	m.MainConfig().father = m
	return m, nil
}

func (m *Main) Type() context_type.ContextType {
	return context_type.TypeMain
}

func (m *Main) IsEnabled() bool {
	return m.MainConfig().IsEnabled()
}

func (m *Main) Enable() context.Context {
	m.MainConfig().Enable()
	return m
}

func (m *Main) Disable() context.Context {
	m.MainConfig().Disable()
	return m
}
