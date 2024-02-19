package local

import (
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/dominikbraun/graph"
	"github.com/marmotedu/errors"
	"strings"
)

type Config struct {
	BasicContext       `json:"config"`
	ConfigGraph        `json:"-"`
	context.ConfigPath `json:"-"`
}

func (c *Config) Father() context.Context {
	return context.NullContext()
}

func (c *Config) SetFather(_ context.Context) error {
	return errors.WithCode(code.V3ErrInvalidOperation, "cannot set father for Config Context")
}

func (c *Config) ConfigLines(isDumping bool) ([]string, error) {
	lines := make([]string, 0)
	for _, child := range c.Children {
		clines, err := child.ConfigLines(isDumping)
		if err != nil {
			return nil, err
		}
		if clines != nil {
			for _, cline := range clines {
				lines = append(lines, cline)
			}
		}
	}

	return lines, nil
}

func (c *Config) IncludeConfig(configs ...*Config) error {
	if configs == nil {
		return errors.WithCode(code.V3ErrInvalidContext, "null config")
	}
	errs := make([]error, 0)
	for _, config := range configs {

		switch config.ConfigPath.(type) {
		case *context.RelConfigPath:
			if config.BaseDir() != c.BaseDir() {
				errs = append(errs, errors.WithCode(code.V3ErrInvalidContext, "he relative target directory(%s) of the included configuration file does not match the directory(%s) where the main configuration file is located", config.BaseDir(), c.BaseDir()))
				continue
			}
		case nil:
			errs = append(errs, errors.WithCode(code.V3ErrInvalidContext, "config with no ConfigPath"))
		}

		errs = append(errs, c.AddEdge(c, config))
	}
	return errors.NewAggregate(errs)
}

func registerConfigBuilder() error {
	builderMap[context_type.TypeConfig] = func(value string) context.Context {
		ctx := &Config{BasicContext: newBasicContext(context_type.TypeConfig, nullHeadString, nullTailString)}
		ctx.ContextValue = value
		ctx.self = ctx
		return ctx
	}
	return nil
}

type ConfigGraph interface {
	AddEdge(src, dst *Config) error
	RemoveEdge(src, dst *Config) error
	Topology() []*Config
	MainConfig() *Config
	AddConfig(config *Config) error
	GetConfig(fullpath string) (*Config, error)
}

type configGraph struct {
	graph      graph.Graph[string, *Config]
	mainConfig *Config
}

func configHash(t *Config) string {
	return strings.TrimSpace(t.FullPath())
}

func (c *configGraph) AddEdge(src, dst *Config) error {
	if src == nil || dst == nil || configHash(src) == "" || configHash(dst) == "" {
		return errors.WithCode(code.V3ErrInvalidContext, "source or destination config is nil")
	}
	err := c.setGraphFor(src)
	if err != nil {
		return err
	}
	err = c.setGraphFor(dst)
	if err != nil {
		return err
	}

	return c.graph.AddEdge(src.FullPath(), dst.FullPath())
}

func (c *configGraph) RemoveEdge(src, dst *Config) error {
	err := c.graph.RemoveEdge(src.FullPath(), dst.FullPath())
	if err != nil {
		return err
	}
	err = c.removeConfig(dst)
	if err != nil && !errors.Is(err, graph.ErrVertexHasEdges) {
		return err
	}
	return nil
}

func (c *configGraph) GetConfig(fullpath string) (*Config, error) {
	return c.graph.Vertex(fullpath)
}

func (c *configGraph) AddConfig(config *Config) error {
	err := c.setGraphFor(config)
	if err != nil {
		return err
	}
	return c.graph.AddVertex(config)
}

func (c *configGraph) removeConfig(config *Config) error {
	err := c.graph.RemoveVertex(config.FullPath())
	if err != nil {
		return err
	}
	config.ConfigGraph = nil
	return nil
}

func (c *configGraph) setGraphFor(config *Config) error {
	if config == nil {
		return errors.WithCode(code.V3ErrInvalidContext, "added config is nil")
	}
	if config.ConfigGraph == nil {
		config.ConfigGraph = c
	}
	if config.ConfigGraph != c {
		return errors.WithCode(code.V3ErrInvalidContext, "%s config is in another config graph", config.FullPath())
	}
	return nil
}

func (c *configGraph) Topology() []*Config {
	topoHashList, err := graph.TopologicalSort(c.graph)
	if err != nil {
		return nil
	}
	topo := make([]*Config, 0)
	for _, k := range topoHashList {
		t, err := c.graph.Vertex(k)
		if err != nil {
			return nil
		}
		topo = append(topo, t)
	}
	return topo
}

func (c *configGraph) MainConfig() *Config {
	return c.mainConfig
}

func newMainConfig(abspath string) (*Config, error) {
	confpath, err := context.NewAbsConfigPath(abspath)
	if err != nil {
		return nil, err
	}
	g := graph.New(configHash, graph.PreventCycles(), graph.Directed())
	config := &Config{
		BasicContext: newBasicContext(context_type.TypeConfig, nullHeadString, nullTailString),
		ConfigPath:   confpath,
	}
	err = g.AddVertex(config)
	if err != nil {
		return nil, err
	}
	config.self = config
	config.ContextValue = config.FullPath()
	config.ConfigGraph = &configGraph{
		graph:      g,
		mainConfig: config,
	}
	return config, nil
}
