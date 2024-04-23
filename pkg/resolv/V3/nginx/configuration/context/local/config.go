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
	BasicContext
	context.ConfigPath `json:"-"`
}

func (c *Config) isInGraph() bool {
	if c.ConfigPath == nil {
		return false
	}
	fatherMain, ok := c.Father().(MainContext)
	if !ok || fatherMain.graph() == nil {
		return false
	}
	_, err := fatherMain.GetConfig(configHash(c))
	return err == nil
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

func (c *Config) SetFather(ctx context.Context) error {
	if _, ok := ctx.(MainContext); !ok {
		return errors.WithCode(code.ErrV3InvalidOperation, "the set father is not Main Context")
	}
	c.father = ctx
	return nil
}

func (c *Config) ConfigLines(isDumping bool) ([]string, error) {
	lines := make([]string, 0)
	for idx, child := range c.Children {
		clines, err := child.ConfigLines(isDumping)
		if err != nil {
			return nil, err
		}
		if clines != nil {
			if child.Type() == context_type.TypeInlineComment && len(lines) > 0 &&
				c.Child(idx-1).Type() != context_type.TypeComment &&
				c.Child(idx-1).Type() != context_type.TypeInlineComment {
				lines[len(lines)-1] += INDENT + clines[0]
				continue
			}

			for _, cline := range clines {
				lines = append(lines, cline)
			}
		}
	}

	return lines, nil
}

func (c *Config) checkIncludedConfigs(configs []*Config) error {
	if c.ConfigPath == nil {
		return errors.WithCode(code.ErrV3InvalidOperation, "this config has not been banded with a `ConfigPath`"+
			" and other configs cannot be inserted into this config")
	}

	fatherMain, ok := c.Father().(MainContext)
	if !ok || fatherMain.graph() == nil {
		return errors.WithCode(code.ErrV3InvalidOperation, "this config has not been added to a certain graph"+
			" and other configs cannot be inserted into this config")
	}

	if configs == nil {
		return errors.WithCode(code.ErrV3InvalidContext, "null config")
	}
	for _, config := range configs {
		if config == nil {
			return errors.WithCode(code.ErrV3InvalidContext, "config with no ConfigPath")
		}
		if config.ConfigPath != nil {
			if _, ok := config.ConfigPath.(*context.RelConfigPath); ok && config.ConfigPath.BaseDir() != fatherMain.MainConfig().BaseDir() {
				return errors.WithCode(code.ErrV3InvalidContext,
					"he relative target directory(%s) of the included configuration file does not match the directory(%s) where the main configuration file is located",
					config.BaseDir(), fatherMain.MainConfig().BaseDir())
			}
		}
	}
	return nil
}

func (c *Config) includeConfig(configs ...*Config) ([]*Config, error) {
	includedConfigs := make([]*Config, 0)
	err := c.checkIncludedConfigs(configs)
	if err != nil {
		return includedConfigs, err
	}

	var errs []error
	fatherMain := c.Father().(MainContext)

	for _, config := range configs {
		if config.isInGraph() && config.Father() != fatherMain {
			config = config.Clone().(*Config)
		}

		configpath, err := newConfigPath(fatherMain, config.Value())
		if err != nil {
			errs = append(errs, err)
			continue
		}
		config.ConfigPath = configpath

		err = fatherMain.AddConfig(config)
		if err != nil {
			if !errors.Is(err, graph.ErrVertexAlreadyExists) { // Temporarily unable to be covered for testing
				errs = append(errs, err)
				continue
			}
			// get cache from graph
			cache, err := fatherMain.GetConfig(configHash(config))
			if err != nil { // Temporarily unable to be covered for testing
				errs = append(errs, err)
				continue
			}

			config = cache
		}

		err = fatherMain.AddEdge(c, config)
		if err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) {
			errs = append(errs, err)
			continue
		}

		includedConfigs = append(includedConfigs, config)
	}
	return includedConfigs, errors.NewAggregate(errs)
}

func (c *Config) removeIncludedConfig(configs ...*Config) ([]*Config, error) {
	removedConfigs := make([]*Config, 0)
	err := c.checkIncludedConfigs(configs)
	if err != nil {
		return removedConfigs, err
	}

	var errs []error
	fatherMain := c.Father().(MainContext)
	for _, config := range configs {
		// get cache from graph
		cache, err := fatherMain.GetConfig(configHash(config))
		if err != nil {
			errs = append(errs, err)
			removedConfigs = append(removedConfigs, config)
			continue
		}
		removedConfigs = append(removedConfigs, cache)
		errs = append(errs, fatherMain.RemoveEdge(c, cache))
	}
	return removedConfigs, errors.NewAggregate(errs)
}

func (c *Config) modifyPathInGraph(path string) error {
	if !c.isInGraph() {
		return nil
	}
	fatherMain := c.Father().(MainContext)

	targetConfigPath, err := newConfigPath(fatherMain, path)
	if err != nil {
		return err
	}

	if configHash(c) == strings.TrimSpace(targetConfigPath.FullPath()) {
		return nil
	}

	targetConfig, err := fatherMain.GetConfig(strings.TrimSpace(targetConfigPath.FullPath()))
	if err != nil {
		if !errors.Is(err, graph.ErrVertexNotFound) { // Temporarily unable to be covered for testing
			return err
		}
	} else {
		return errors.Wrapf(graph.ErrVertexAlreadyExists, "config(%s) is already exists in config graph", configHash(targetConfig))
	}
	oldPath := configHash(c)
	c.ConfigPath = targetConfigPath
	cache, _ := fatherMain.GetConfig(oldPath)
	if cache == c {
		return fatherMain.RenewConfigPath(oldPath)
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
		return errors.WithCode(code.ErrV3InvalidContext, "source or destination config is nil")
	}
	err := c.setFatherFor(src)
	if err != nil {
		return err
	}
	err = c.setFatherFor(dst)
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
	err := c.setFatherFor(config)
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
	config.father = context.NullContext()
	return nil
}

func (c *configGraph) setFatherFor(config *Config) error {
	if config == nil {
		return errors.WithCode(code.ErrV3InvalidContext, "added config is nil")
	}
	if config.isInGraph() && config.Father() != c.MainConfig().Father() {
		return errors.WithCode(code.ErrV3InvalidContext, "%s config is in another config graph", configHash(config))
	}
	return config.SetFather(c.MainConfig().Father())
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
	if configHash(config) == fullpath {
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

func newConfigGraph(mainConfig *Config) (ConfigGraph, error) {
	if mainConfig.ConfigPath == nil {
		return nil, errors.WithCode(code.ErrV3InvalidContext, "main config's ConfigPath is nil")
	}
	if mainConfig.isInGraph() {
		return nil, errors.WithCode(code.ErrV3InvalidContext, "main config(%s) is in another config graph", configHash(mainConfig))
	}
	g := graph.New(configHash, graph.PreventCycles(), graph.Directed())
	err := g.AddVertex(mainConfig)
	if err != nil { // Temporarily unable to be covered for testing
		return nil, err
	}
	return &configGraph{
		graph:      g,
		mainConfig: mainConfig,
	}, nil
}
