package local

import (
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/dominikbraun/graph"
	"github.com/marmotedu/errors"
	"path/filepath"
	"strings"
)

type Config struct {
	BasicContext       `json:"config"`
	Graph              ConfigGraph `json:"-"`
	context.ConfigPath `json:"-"`
}

func (c *Config) isInGraph() bool {
	if c.Graph == nil || c.ConfigPath == nil {
		return false
	}
	_, err := c.Graph.GetConfig(configHash(c))
	return err == nil
}

func (c *Config) Father() context.Context {
	return context.NullContext()
}

func (c *Config) Clone() context.Context {
	clone := NewContext(context_type.TypeConfig, c.ContextValue)
	for i, child := range c.Children {
		clone.Insert(child.Clone(), i)
	}
	return clone
}

func (c *Config) SetValue(v string) error {
	err := c.modifyPathInGraph(v)
	if err != nil {
		return err
	}

	c.ContextValue = v

	return nil
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

func (c *Config) includeConfig(configs ...*Config) ([]*Config, error) {
	errs := make([]error, 0)
	if c.Graph == nil {
		errs = append(errs,
			errors.WithCode(code.V3ErrInvalidOperation, "this config has not been added to a certain graph"+
				" and other configs cannot be inserted into this config"))
	}
	if c.ConfigPath == nil {
		errs = append(errs,
			errors.WithCode(code.V3ErrInvalidOperation, "this config has not been banded with a `ConfigPath`"+
				" and other configs cannot be inserted into this config"))
	}
	if err := errors.NewAggregate(errs); err != nil {
		return nil, err
	}

	if configs == nil {
		return nil, errors.WithCode(code.V3ErrInvalidContext, "null config")
	}
	includedConfigs := make([]*Config, 0)
	for _, config := range configs {
		if config == nil {
			errs = append(errs, errors.WithCode(code.V3ErrInvalidContext, "config with no ConfigPath"))
			continue
		}
		if config.Graph != nil && config.Graph != c.Graph {
			config = config.Clone().(*Config)
		}
		if !config.isInGraph() {
			configpath, err := newConfigPath(c.Graph, config.Value())
			if err != nil {
				errs = append(errs, err)
				continue
			}
			config.Graph = c.Graph
			config.ConfigPath = configpath
		}

		switch config.ConfigPath.(type) {
		case *context.RelConfigPath:
			if config.BaseDir() != c.BaseDir() {
				errs = append(errs, errors.WithCode(code.V3ErrInvalidContext,
					"he relative target directory(%s) of the included configuration file does not match the directory(%s) where the main configuration file is located",
					config.BaseDir(), c.BaseDir()))
				continue
			}
		}
		err := c.Graph.AddConfig(config)
		if err != nil {
			if !errors.Is(err, graph.ErrVertexAlreadyExists) { // Temporarily unable to be covered for testing
				errs = append(errs, err)
				continue
			}
			// get cache from graph
			cache, err := c.Graph.GetConfig(configHash(config))
			if err != nil { // Temporarily unable to be covered for testing
				errs = append(errs, err)
				continue
			}

			if config != cache {
				errs = append(errs, errors.WithCode(code.V3ErrInvalidContext, "the config(%s) is already exist in graph,"+
					" and the inserted config is inconsistent with the cache in the graph", config.FullPath()))
				continue
			}
		}

		err = c.Graph.AddEdge(c, config)
		if err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
			errs = append(errs, err)
			continue
		}

		includedConfigs = append(includedConfigs, config)
	}
	return includedConfigs, errors.NewAggregate(errs)
}

func (c *Config) modifyPathInGraph(path string) error {
	if !c.isInGraph() {
		return nil
	}

	targetConfigPath, err := newConfigPath(c.Graph, path)
	if err != nil {
		return err
	}

	if c.FullPath() == targetConfigPath.FullPath() {
		return nil
	}

	targetConfig, err := c.Graph.GetConfig(targetConfigPath.FullPath())
	if err != nil {
		if !errors.Is(err, graph.ErrVertexNotFound) { // Temporarily unable to be covered for testing
			return err
		}
	} else {
		return errors.Wrapf(graph.ErrVertexAlreadyExists, "config(%s) is already exists in config graph", targetConfig.FullPath())
	}
	oldPath := c.FullPath()
	c.ConfigPath = targetConfigPath
	cache, _ := c.Graph.GetConfig(oldPath)
	if cache == c {
		return c.Graph.RenewConfigPath(oldPath)
	}
	return nil
}

func newConfigPath(configgraph ConfigGraph, newconfigpath string) (context.ConfigPath, error) {
	if configgraph == nil {
		return nil, errors.New("config graph is nil")
	}
	if strings.TrimSpace(newconfigpath) == "" {
		return nil, errors.New("new config path is null")
	}
	if !filepath.IsAbs(newconfigpath) && !filepath.IsAbs(configHash(configgraph.MainConfig())) {
		return nil, errors.Errorf("main config path(%s) is not an absolute path", configHash(configgraph.MainConfig()))
	}
	if filepath.IsAbs(newconfigpath) {
		return context.NewAbsConfigPath(newconfigpath)
	}
	return context.NewRelConfigPath(configgraph.MainConfig().BaseDir(), newconfigpath)
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
	RenewConfigPath(fullpath string) error
}

type configGraph struct {
	graph      graph.Graph[string, *Config]
	mainConfig *Config
}

func configHash(t *Config) string {
	if t == nil {
		return ""
	}
	if t.ConfigPath == nil {
		return ""
	}
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

	return c.graph.AddEdge(configHash(src), configHash(dst))
}

func (c *configGraph) RemoveEdge(src, dst *Config) error {
	err := c.graph.RemoveEdge(configHash(src), configHash(dst))
	if err != nil {
		return err
	}
	err = c.removeConfig(dst)
	if err != nil && !errors.Is(err, graph.ErrVertexHasEdges) { // Temporarily unable to be covered for testing
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
	err := c.graph.RemoveVertex(configHash(config))
	if err != nil {
		return err
	}
	config.Graph = nil
	return nil
}

func (c *configGraph) setGraphFor(config *Config) error {
	if config == nil {
		return errors.WithCode(code.V3ErrInvalidContext, "added config is nil")
	}
	if config.Graph == nil {
		config.Graph = c
	}
	if config.Graph != c {
		return errors.WithCode(code.V3ErrInvalidContext, "%s config is in another config graph", configHash(config))
	}
	return nil
}

func (c *configGraph) Topology() []*Config {
	topoHashList, err := graph.TopologicalSort(c.graph)
	if err != nil { // Temporarily unable to be covered for testing
		return nil
	}

	topo := make([]*Config, 0)
	for _, k := range topoHashList {
		if _, err := graph.ShortestPath(c.graph, configHash(c.MainConfig()), k); err != nil {
			// remove configs that are not directly or indirectly referenced by the main config
			continue
		}
		t, err := c.graph.Vertex(k)
		if err != nil { // Temporarily unable to be covered for testing
			return nil
		}
		topo = append(topo, t)
	}
	return topo
}

func (c *configGraph) MainConfig() *Config {
	return c.mainConfig
}

func (c *configGraph) RenewConfigPath(fullpath string) error {
	config, err := c.GetConfig(fullpath)
	if err != nil {
		return err
	}
	if config.FullPath() == fullpath {
		return nil
	}

	err = c.graph.AddVertex(config)
	if err != nil {
		return err
	}

	edges, err := c.graph.Edges()
	if err != nil { // Temporarily unable to be covered for testing
		return err
	}

	for _, edge := range edges {
		switch fullpath {
		case edge.Source:
			err = c.graph.AddEdge(configHash(config), edge.Target)
		case edge.Target:
			err = c.graph.AddEdge(edge.Source, configHash(config))
		default:
			continue
		}
		if err != nil { // Temporarily unable to be covered for testing
			return err
		}
		err = c.graph.RemoveEdge(edge.Source, edge.Target)
		if err != nil { // Temporarily unable to be covered for testing
			return err
		}
	}
	return c.graph.RemoveVertex(fullpath)
}

func newMainConfig(abspath string) (*Config, error) {
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
	err = setGraphForMainConfig(config)
	if err != nil { // Temporarily unable to be covered for testing
		return nil, err
	}
	return config, nil
}

func setGraphForMainConfig(mainConfig *Config) error {
	if mainConfig.isInGraph() || (mainConfig.Graph != nil && mainConfig.Graph.MainConfig() == mainConfig) {
		return errors.WithCode(code.V3ErrInvalidContext, "main config(%s) is in another config graph", configHash(mainConfig))
	}
	g := graph.New(configHash, graph.PreventCycles(), graph.Directed())
	err := g.AddVertex(mainConfig)
	if err != nil { // Temporarily unable to be covered for testing
		return err
	}
	mainConfig.Graph = &configGraph{
		graph:      g,
		mainConfig: mainConfig,
	}
	return nil
}
