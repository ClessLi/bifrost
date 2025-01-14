package local

import (
	"encoding/json"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
	"path/filepath"
	"strings"
	"sync"
)

type includeSnapshot struct {
	isEnabled          bool
	isInTopology       bool
	mainContext        MainContext
	fatherConfig       *Config
	includePatternPath string
	includedConfigs    []*Config
}

func (i includeSnapshot) matchConfig(config *Config) error {
	if config == nil {
		return errors.WithCode(code.ErrV3InvalidContext, "pattern(%s) cannot match nil", i.includePatternPath)
	}
	isMatch, err := filepath.Match(i.includePatternPath, config.FullPath())
	if err != nil {
		return errors.WithCode(code.ErrV3InvalidContext, "pattern(%s) match included config(%s) failed, cased by: %v", i.includePatternPath, config.FullPath(), err)
	}
	if !isMatch {
		return errors.WithCode(code.ErrV3InvalidContext, "pattern(%s) cannot match included config(%s)", i.includePatternPath, config.FullPath())
	}
	return nil
}

type Include struct {
	enabled      bool
	ContextValue string
	loadLocker   *sync.RWMutex

	fatherContext context.Context

	snapshot *includeSnapshot
}

func (i *Include) MarshalJSON() ([]byte, error) {
	i.loadLocker.RLock()
	defer i.loadLocker.RUnlock()
	marshalCtx := struct {
		Enabled     bool                     `json:"enabled,omitempty"`
		ContextType context_type.ContextType `json:"context-type"`
		Value       string                   `json:"value,omitempty"`
	}(struct {
		Enabled     bool
		ContextType context_type.ContextType
		Value       string
	}{
		Enabled:     i.enabled,
		ContextType: context_type.TypeInclude,
		Value:       i.ContextValue,
	})

	return json.Marshal(marshalCtx)
}

func (i *Include) parsePatternPath() (main MainContext, fatherConfig *Config, isEnabled bool, pattern string, err error) {
	isEnabled, fatherConfig, err = i.fatherConfig()
	if !i.enabled {
		isEnabled = false
	}
	if err != nil {
		return nil, fatherConfig, isEnabled, "", err
	}
	mainCtx, ok := fatherConfig.Father().(MainContext)
	if !ok {
		return nil, fatherConfig, isEnabled, "", errors.WithCode(code.ErrV3InvalidContext, "this include context is not bound to a main context")
	}
	if filepath.IsAbs(i.Value()) {
		return mainCtx, fatherConfig, isEnabled, i.Value(), nil
	}
	return mainCtx, fatherConfig, isEnabled, filepath.Join(mainCtx.MainConfig().BaseDir(), i.Value()), nil
}

func (i *Include) PatternPath() string {
	i.loadLocker.RLock()
	defer i.loadLocker.RUnlock()
	_, _, _, pattern, _ := i.parsePatternPath()
	return pattern
}

func (i *Include) parseSnapshot() (*includeSnapshot, error) {
	if i.snapshot != nil {
		return i.snapshot, nil
	}
	mainCtx, fatherConfig, isEnabled, pattern, err := i.parsePatternPath()
	if err != nil {
		return nil, err
	}
	var isInTopology = false
	for _, config := range mainCtx.Topology() {
		if fatherConfig == config {
			isInTopology = true
			break
		}
	}

	includes := make([]*Config, 0)
	if len(strings.TrimSpace(pattern)) > 0 {
		for _, config := range mainCtx.ListConfigs() {
			isMatch, err := filepath.Match(pattern, config.FullPath())
			if err == nil && isMatch {
				includes = append(includes, config)
			}
		}
	}
	i.snapshot = &includeSnapshot{
		isEnabled:          isEnabled,
		isInTopology:       isInTopology,
		mainContext:        mainCtx,
		fatherConfig:       fatherConfig,
		includePatternPath: pattern,
		includedConfigs:    includes,
	}
	return i.snapshot, nil
}

func (i *Include) Configs() []*Config {
	i.loadLocker.RLock()
	defer i.loadLocker.RUnlock()
	snapshot, err := i.parseSnapshot()
	if err != nil {
		return nil
	}
	return snapshot.includedConfigs
}

func (i *Include) fatherConfig() (isEnabled bool, fatherConfig *Config, err error) {
	fatherCtx := i.fatherContext
	isEnabled = fatherCtx.IsEnabled()
	fatherConfig, ok := fatherCtx.(*Config)
	for !ok {
		switch fatherCtx.Type() {
		case context_type.TypeErrContext:
			return isEnabled, fatherConfig, fatherCtx.(*context.ErrorContext).AppendError(errors.WithCode(code.ErrV3InvalidContext, "found an error config context")).Error()
		case context_type.TypeMain:
			return isEnabled, fatherConfig, errors.WithCode(code.ErrV3InvalidContext, "found an Main context")
		}

		fatherCtx = fatherCtx.Father()
		if !fatherCtx.IsEnabled() {
			isEnabled = false
		}
		fatherConfig, ok = fatherCtx.(*Config)
	}
	return isEnabled, fatherConfig, nil
}

func (i *Include) FatherConfig() (*Config, error) {
	i.loadLocker.RLock()
	defer i.loadLocker.RUnlock()
	_, fc, err := i.fatherConfig()
	return fc, err
}

func (i *Include) Insert(ctx context.Context, idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "include cannot insert by index"))
}

func (i *Include) Remove(idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "include cannot remove by index"))
}

func (i *Include) Modify(ctx context.Context, idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "include cannot modify by index"))
}

func (i *Include) Father() context.Context {
	i.loadLocker.RLock()
	defer i.loadLocker.RUnlock()
	return i.fatherContext
}

func (i *Include) Child(idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.ErrV3InvalidOperation, "include cannot get child config by index"))
}

func (i *Include) QueryByKeyWords(kw context.KeyWords) context.Pos {
	i.loadLocker.RLock()
	defer i.loadLocker.RUnlock()
	// Avoid loop queries
	snapshot, err := i.parseSnapshot()
	if err != nil || !snapshot.isEnabled || !snapshot.isInTopology { // If the snapshot is not resolved, disabled or not in the topology, there is no need to load includes
		return context.NotFoundPos()
	}
	if !kw.Cascaded() {
		return context.NotFoundPos()
	}
	for _, config := range i.Configs() {
		if kw.SkipQueryThisContext(config) {
			continue
		}
		pos := config.QueryByKeyWords(kw)
		if pos != context.NotFoundPos() {
			return pos
		}
	}
	return context.NotFoundPos()
}

func (i *Include) QueryAllByKeyWords(kw context.KeyWords) []context.Pos {
	i.loadLocker.RLock()
	defer i.loadLocker.RUnlock()
	poses := make([]context.Pos, 0)
	// Avoid loop queries
	snapshot, err := i.parseSnapshot()
	if err != nil || !snapshot.isEnabled || !snapshot.isInTopology { // If the snapshot is not resolved, disabled or not in the topology, there is no need to load includes
		return poses
	}
	if !kw.Cascaded() {
		return poses
	}
	for _, config := range i.Configs() {
		if kw.SkipQueryThisContext(config) {
			continue
		}
		poses = append(poses, config.QueryAllByKeyWords(kw)...)
	}
	return poses
}

func (i *Include) Clone() context.Context {
	clone := NewContext(context_type.TypeInclude, i.ContextValue)
	if !i.enabled {
		return clone.Disable()
	}
	return clone
}

func (i *Include) SetValue(v string) error {
	return errors.WithCode(code.ErrV3InvalidOperation, "setting the value of include context is unsafe")
}

func (i *Include) load() error {
	i.snapshot = nil // fresh snapshot
	snapshot, err := i.parseSnapshot()
	if err != nil || !snapshot.isEnabled || !snapshot.isInTopology { // If the snapshot is not resolved, disabled or not in the topology, there is no need to load includes
		return nil
	}
	var errs []error
	for _, includedConfig := range snapshot.includedConfigs {
		err = snapshot.mainContext.AddEdge(snapshot.fatherConfig, includedConfig)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		// call includes context loading, which is located in included configs
		includeKeyWords := context.NewKeyWords(context_type.TypeInclude).
			SetCascaded(true).
			SetSkipQueryFilter(func(targetCtx context.Context) bool {
				// skip the include context of indirect connections
				targetInclude, ok := targetCtx.(*Include)
				if !ok {
					return false
				}
				_, subFatherConfig, _ := targetInclude.fatherConfig()
				return subFatherConfig != includedConfig
			}).SetSkipQueryFilter(func(targetCtx context.Context) bool {
			// skip self
			return targetCtx == i
		})
		for _, pos := range includedConfig.QueryAllByKeyWords(includeKeyWords) {
			include, ok := pos.Target().(*Include)
			if !ok {
				errs = append(errs, errors.WithCode(code.ErrV3InvalidContext, "[%v] is not an Include context", pos.Target()))
				continue
			}

			errs = append(errs, include.load())
		}
	}
	return errors.NewAggregate(errs)
}

func (i *Include) unload() error {
	defer func() { // release snapshot
		i.snapshot = nil
	}()
	snapshot, err := i.parseSnapshot()
	if err != nil { // If the snapshot is not resolved, there is no need to unload includes
		return nil
	}
	var errs []error
	for _, includedConfig := range snapshot.includedConfigs {
		err = snapshot.mainContext.RemoveEdge(snapshot.fatherConfig, includedConfig, true)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		// call includes context unloading, which is located in included configs
		includeKeyWords := context.NewKeyWords(context_type.TypeInclude).
			SetCascaded(true).
			SetSkipQueryFilter(func(targetCtx context.Context) bool {
				// skip the include context of indirect connections
				targetInclude, ok := targetCtx.(*Include)
				if !ok {
					return false
				}
				_, subFatherConfig, _ := targetInclude.fatherConfig()
				return subFatherConfig != includedConfig
			}).SetSkipQueryFilter(func(targetCtx context.Context) bool {
			// skip self
			return targetCtx == i
		})
		for _, pos := range includedConfig.QueryAllByKeyWords(includeKeyWords) {
			include, ok := pos.Target().(*Include)
			if !ok {
				errs = append(errs, errors.WithCode(code.ErrV3InvalidContext, "[%v] is not an Include context", pos.Target()))
				continue
			}

			errs = append(errs, include.unload())
		}
	}
	return errors.NewAggregate(errs)
}

func (i *Include) reload() error {
	_ = i.unload()
	return i.load()
}

func (i *Include) SetFather(ctx context.Context) error {
	i.loadLocker.Lock()
	defer i.loadLocker.Unlock()
	// exclude old configs, pass unload error
	_ = i.unload()
	// include new configs
	i.fatherContext = ctx
	if i.enabled { // reload when include is enabled
		return i.load()
	}
	return nil
}

func (i *Include) HasChild() bool {
	return i.Len() > 0
}

func (i *Include) Len() int {
	return len(i.Configs())
}

func (i *Include) Value() string {
	return filepath.Clean(i.ContextValue)
}

func (i *Include) Type() context_type.ContextType {
	return context_type.TypeInclude
}

func (i *Include) Error() error {
	return nil
}

func (i *Include) ConfigLines(isDumping bool) ([]string, error) {
	i.loadLocker.RLock()
	defer i.loadLocker.RUnlock()
	lines := make([]string, 0)
	if !i.enabled {
		defer func() {
			if len(lines) > 0 {
				for idx := range lines {
					lines[idx] = "# " + lines[idx]
				}
			}
		}()
	}
	if isDumping {
		lines = append(lines, "include "+i.ContextValue+";")
		return lines, nil
	}

	lines = append(lines, "# include <== "+i.ContextValue)
	if i.enabled { // Avoid loop rendering
		for _, config := range i.Configs() {
			configlines, err := config.ConfigLines(isDumping)
			if err != nil {
				return nil, err
			}
			lines = append(lines, configlines...)
		}
	}
	return lines, nil
}

func (i *Include) IsEnabled() bool {
	i.loadLocker.RLock()
	defer i.loadLocker.RUnlock()
	return i.enabled
}

func (i *Include) Enable() context.Context {
	i.loadLocker.Lock()
	defer i.loadLocker.Unlock()
	if i.enabled {
		return i
	}

	i.enabled = true
	err := i.load()
	if err != nil {
		i.enabled = false
		return context.ErrContext(err, i.unload())
	}
	return i
}

func (i *Include) Disable() context.Context {
	i.loadLocker.Lock()
	defer i.loadLocker.Unlock()
	if !i.enabled {
		return i
	}

	err := i.unload()
	if err != nil {
		return context.ErrContext(err, i.load())
	}
	i.enabled = false
	return i
}

func registerIncludeBuilder() error {
	builderMap[context_type.TypeInclude] = func(value string) context.Context {
		return &Include{
			enabled:       true,
			ContextValue:  value,
			loadLocker:    new(sync.RWMutex),
			fatherContext: context.NullContext(),
		}
	}
	return nil
}
