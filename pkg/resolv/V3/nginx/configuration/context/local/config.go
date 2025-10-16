package local

import (
	"path/filepath"
	"strings"

	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"

	"github.com/dominikbraun/graph"
	"github.com/marmotedu/errors"
)

type Config struct {
	BasicContext
	context.ConfigPath `json:"-"`
}

func (c *Config) isInGraph() bool {
	if c.ConfigPath == nil {
		return false
	}
	mainCtx, err := c.mainContext()
	if err != nil {
		return false
	}
	if mainCtx.graph() == nil {
		return false
	}
	_, err = mainCtx.GetConfig(configHash(c))

	return err == nil
}

func (c *Config) FatherPosSet() context.PosSet {
	if c.isInGraph() {
		// return ErrPosSet, which is result of the MainContext's FatherPosSet(), if config itself is main config
		mainCtx, err := c.mainContext()
		if err == nil && mainCtx.MainConfig() == c {
			return mainCtx.FatherPosSet()
		}

		// return father Include pos list
		fatherIncludePosSet := context.NewPosSet()
		edges, err := mainCtx.graph().(*configGraph).graph.Edges()
		if err != nil {
			return context.ErrPosSet(err)
		}

		for _, edge := range edges {
			if edge.Target == configHash(c) {
				fc, err := mainCtx.GetConfig(edge.Source)
				if err != nil {
					return context.ErrPosSet(err)
				}
				fatherIncludePosSet.AppendWithPosSet(
					fc.ChildrenPosSet().QueryAll(
						context.NewKeyWordsByType(context_type.TypeInclude),
					).Filter(func(pos context.Pos) bool {
						t := pos.Target()
						i, ok := t.(*Include)
						if !ok {
							return false
						}

						i.loadLocker.Lock()
						defer i.loadLocker.Unlock()
						for _, ic := range i.snapshot.includedConfigs {
							if ic == c {
								return true
							}
						}

						return false
					}),
				)
			}
		}

		return fatherIncludePosSet
	}

	return context.NewPosSet()
}

func (c *Config) Clone() context.Context {
	clone := NewContext(context_type.TypeConfig, c.ContextValue)
	if !c.Enabled {
		clone.Disable()
	}
	for i, child := range c.Children {
		clone.Insert(child.Clone(), i)
	}

	return clone
}

func (c *Config) SetValue(v string) error {
	mainCtx, err := c.mainContext()
	if err != nil {
		return err
	}

	return mainCtx.RenameConfig(configHash(c), v)
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

			lines = append(lines, clines...)
		}
	}

	if !c.IsEnabled() && len(lines) > 0 {
		for i := range lines {
			lines[i] = "# " + lines[i]
		}
	}

	return lines, nil
}

func (c *Config) mainContext() (MainContext, error) {
	if c.ConfigPath == nil {
		return nil, errors.WithCode(code.ErrV3InvalidOperation, "this config has not been banded with a `ConfigPath`"+
			" and other configs cannot be inserted into this config")
	}

	fatherMain, ok := c.Father().(MainContext)
	if !ok || fatherMain.graph() == nil {
		return nil, errors.WithCode(code.ErrV3InvalidOperation, "this config has not been added to a certain graph"+
			" and other configs cannot be inserted into this config")
	}

	return fatherMain, nil
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
	checkOperatedVertex(v *Config) error
	addVertex(v *Config) error
	removeVertex(v *Config) error
	cleanupGraph() error
	renderGraph() error
	rerenderGraph() error

	AddEdge(src, dst *Config) error
	RemoveEdge(src, dst *Config, keepDst bool) error
	Topology() []*Config
	ListConfigs() []*Config
	MainConfig() *Config
	AddConfig(config *Config) error
	RemoveConfig(config *Config) error
	GetConfig(fullpath string) (*Config, error)
	RenameConfig(oldFullPath, newPath string) error
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

	err = c.graph.AddEdge(configHash(src), configHash(dst))
	if err != nil && !errors.Is(err, graph.ErrEdgeAlreadyExists) { // allow repeated addition of edges
		return err
	}

	return nil
}

func (c *configGraph) RemoveEdge(src, dst *Config, keepDst bool) error {
	err := c.graph.RemoveEdge(configHash(src), configHash(dst))
	if err != nil && !errors.Is(err, graph.ErrEdgeNotFound) { // allow repeated deletion of edges
		return err
	}
	if keepDst {
		return nil
	}
	err = c.removeVertex(dst)
	if err != nil && !errors.Is(err, graph.ErrVertexHasEdges) { // Temporarily unable to be covered for testing
		return err
	}

	return nil
}

func (c *configGraph) GetConfig(fullpath string) (*Config, error) {
	return c.graph.Vertex(fullpath)
}

func (c *configGraph) checkOperatedVertex(v *Config) error {
	if c == nil {
		return errors.WithCode(code.ErrV3InvalidOperation, "this is a nil ConfigGraph")
	}
	if v == nil {
		return errors.WithCode(code.ErrV3InvalidContext, "added config is nil")
	}
	if v.isInGraph() && v.Father() != c.MainConfig().Father() {
		return errors.WithCode(code.ErrV3InvalidContext, "%s config is in another config graph", configHash(v))
	}
	if v.ConfigPath != nil {
		if _, ok := v.ConfigPath.(*context.RelConfigPath); ok && v.BaseDir() != c.MainConfig().BaseDir() {
			return errors.WithCode(code.ErrV3InvalidContext,
				"he relative target directory(%s) of the included configuration file does not match the directory(%s) where the main configuration file is located",
				v.BaseDir(), c.MainConfig().BaseDir())
		}
	}

	return nil
}

func (c *configGraph) addVertex(v *Config) error {
	err := c.setFatherFor(v)
	if err != nil {
		return err
	}

	return c.graph.AddVertex(v)
}

func (c *configGraph) cleanupGraph() error {
	var errs []error
	edges, err := c.graph.Edges()
	if err != nil {
		return err
	}
	for _, edge := range edges {
		errs = append(errs, c.graph.RemoveEdge(edge.Source, edge.Target))
	}

	return errors.NewAggregate(errs)
}

func (c *configGraph) renderGraph() error {
	return c.MainConfig().
		ChildrenPosSet().
		QueryAll(context.NewKeyWordsByType(context_type.TypeInclude).SetCascaded(true)).
		Filter(
			func(pos context.Pos) bool {
				_, ok := pos.Target().(*Include)

				return ok
			},
		).
		Map( // TODO: reduce for aggregate errors
			func(pos context.Pos) (context.Pos, error) {
				return pos, pos.Target().(*Include).load()
			},
		).
		Error()
}

func (c *configGraph) rerenderGraph() error {
	err := c.cleanupGraph()
	if err != nil {
		return err
	}

	return c.renderGraph()
}

func (c *configGraph) AddConfig(config *Config) error {
	err := c.addVertex(config)
	if err != nil {
		return err
	}
	// call include context reload
	return c.rerenderGraph()
}

func (c *configGraph) removeVertex(v *Config) error {
	err := c.graph.RemoveVertex(configHash(v))
	if err != nil {
		return err
	}
	v.father = context.NullContext()

	return nil
}

func (c *configGraph) RemoveConfig(config *Config) error {
	_, err := c.GetConfig(configHash(config))
	if err != nil {
		return err
	}

	if configHash(c.mainConfig) == configHash(config) {
		return errors.WithCode(code.ErrV3InvalidContext, "cannot remove the main config(%s)", configHash(config))
	}

	err = c.cleanupGraph()
	if err != nil {
		return err
	}

	return errors.NewAggregate([]error{c.removeVertex(config), c.renderGraph()})
}

func (c *configGraph) setFatherFor(config *Config) error {
	err := c.checkOperatedVertex(config)
	if err != nil {
		return err
	}
	config.ConfigPath, err = newConfigPath(c, config.ContextValue)
	if err != nil {
		return err
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

func (c *configGraph) ListConfigs() []*Config {
	edgeMap, err := c.graph.AdjacencyMap()
	if err != nil {
		return nil
	}
	list := make([]*Config, 0)
	for k := range edgeMap {
		conf, err := c.graph.Vertex(k)
		if err != nil {
			return nil
		}
		list = append(list, conf)
	}

	return list
}

func (c *configGraph) MainConfig() *Config {
	return c.mainConfig
}

func (c *configGraph) RenameConfig(oldFullPath, newPath string) error {
	// check if the old config exists
	targetConfig, err := c.GetConfig(oldFullPath)
	if err != nil {
		return err
	}
	// check if the new path is valid
	targetConfigPath, err := newConfigPath(c, newPath)
	if err != nil {
		return err
	}
	// check if the old and new paths are the same
	if oldFullPath == strings.TrimSpace(targetConfigPath.FullPath()) {
		return nil
	}
	// check if the config of the new path exists
	existedConfig, err := c.GetConfig(strings.TrimSpace(targetConfigPath.FullPath()))
	if err != nil {
		if !errors.Is(err, graph.ErrVertexNotFound) { // Temporarily unable to be covered for testing
			return err
		}
	} else {
		return errors.Wrapf(graph.ErrVertexAlreadyExists, "the config(%s) already exists in config graph", configHash(existedConfig))
	}

	// clean up the config graph
	err = c.cleanupGraph()
	if err != nil {
		return err
	}

	// rename the target config
	// if the target config is the main config, set the config path with a full path
	if targetConfig == c.mainConfig {
		// check if the directory of the new path of the main config is the same as the original one
		if c.mainConfig.BaseDir() != filepath.Dir(strings.TrimSpace(targetConfigPath.FullPath())) {
			return errors.WithCode(code.ErrV3InvalidOperation, "the main config path cannot be modified to '%s'."+
				" when modifying the path of the main config, the directory where the main config is located cannot be changed", newPath)
		}
		mainConfigFather := c.mainConfig.father
		err = c.removeVertex(targetConfig)
		if err != nil {
			return err
		}
		targetConfig.father = mainConfigFather
		targetConfig.ContextValue = strings.TrimSpace(targetConfigPath.FullPath())
	} else {
		err = c.removeVertex(targetConfig)
		if err != nil {
			return err
		}
		targetConfig.ContextValue = newPath
	}

	// add config and rerender the config graph
	return errors.NewAggregate([]error{c.addVertex(targetConfig), c.renderGraph()})
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
