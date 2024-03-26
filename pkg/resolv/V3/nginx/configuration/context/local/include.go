package local

import (
	"encoding/json"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
	"path/filepath"
)

type Include struct {
	ContextValue string             `json:"value,omitempty"`
	Configs      map[string]*Config `json:"param,omitempty"`

	fatherContext context.Context
}

func (i *Include) MarshalJSON() ([]byte, error) {
	marshalCtx := struct {
		Include struct {
			Value  string   `json:"value,omitempty"`
			Params []string `json:"params,omitempty"`
		} `json:"include,omitempty"`
	}{Include: struct {
		Value  string   `json:"value,omitempty"`
		Params []string `json:"params,omitempty"`
	}(struct {
		Value  string
		Params []string
	}{Value: i.Value(), Params: make([]string, 0)})}

	for _, config := range i.Configs {
		marshalCtx.Include.Params = append(marshalCtx.Include.Params, config.Value())
	}
	return json.Marshal(marshalCtx)
}

func (i *Include) FatherConfig() (*Config, error) {
	fatherCtx := i.Father()
	fatherConfig, ok := fatherCtx.(*Config)
	for !ok {
		switch fatherCtx.Type() {
		case context_type.TypeErrContext:
			return nil, fatherCtx.(*context.ErrorContext).AppendError(errors.WithCode(code.V3ErrInvalidContext, "found an error config context")).Error()
		case context_type.TypeMain:
			return nil, errors.WithCode(code.V3ErrInvalidContext, "found an Main context")
		}

		if fatherCtx.Type() == context_type.TypeErrContext {
		}
		fatherCtx = fatherCtx.Father()
		fatherConfig, ok = fatherCtx.(*Config)
	}
	return fatherConfig, nil
}

func (i *Include) Insert(ctx context.Context, idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.V3ErrInvalidOperation, "include cannot insert by index"))
}

func (i *Include) InsertConfig(configs ...*Config) error {
	if configs == nil {
		return errors.WithCode(code.V3ErrInvalidContext, "null config")
	}

	// find father config
	fatherConfig, err := i.FatherConfig()
	if err != nil {
		return err
	}
	fatherMain, ok := fatherConfig.Father().(*Main)
	if !ok {
		return errors.WithCode(code.V3ErrInvalidContext, "this include context is not bound to a main context")
	}

	includingConfigs := make([]*Config, 0)
	for _, config := range configs {
		if config == nil {
			return errors.WithCode(code.V3ErrInvalidContext, "nil config")
		}
		// match config path
		if config.ConfigPath == nil {
			config.ConfigPath, err = newConfigPath(fatherMain.ConfigGraph, config.Value())
			if err != nil {
				return err
			}
		}

		err = i.matchConfigPath(config.RelativePath())
		if err != nil {
			err = i.matchConfigPath(config.FullPath())
			if err != nil {
				return err
			}
		}

		includingConfigs = append(includingConfigs, config)

	}
	includedConfigs, err := fatherConfig.includeConfig(includingConfigs...)
	for _, config := range includedConfigs {
		i.Configs[configHash(config)] = config
	}
	return err
}

func (i *Include) Remove(idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.V3ErrInvalidOperation, "include cannot remove by index"))
}

func (i *Include) RemoveConfig(configs ...*Config) error {
	if configs == nil {
		return errors.WithCode(code.V3ErrInvalidContext, "null config")
	}
	// find father config
	fatherConfig, err := i.FatherConfig()
	if err != nil {
		return err
	}

	removingConfigs := make([]*Config, 0)
	for _, config := range configs {
		_, has := i.Configs[configHash(config)]
		if has {
			removingConfigs = append(removingConfigs, config)
		}
	}

	removedConfigs, err := fatherConfig.removeIncludedConfig(removingConfigs...)
	for _, config := range removedConfigs {
		delete(i.Configs, configHash(config))
	}
	return err
}

func (i *Include) Modify(ctx context.Context, idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.V3ErrInvalidOperation, "include cannot modify by index"))
}

func (i *Include) ModifyConfig(configs ...*Config) error {
	if configs == nil {
		return errors.WithCode(code.V3ErrInvalidContext, "null config")
	}
	return errors.NewAggregate([]error{i.RemoveConfig(configs...), i.InsertConfig(configs...)})
}

func (i *Include) Father() context.Context {
	return i.fatherContext
}

func (i *Include) Child(idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.V3ErrInvalidOperation, "include cannot get child config by index"))
}

func (i *Include) ChildConfig(fullpath string) (*Config, error) {
	config, has := i.Configs[fullpath]
	if !has {
		return nil, errors.Errorf("%s config has not been included", fullpath)
	}
	return config, nil
}

func (i *Include) QueryByKeyWords(kw context.KeyWords) context.Pos {
	for _, child := range i.Configs {
		pos := child.QueryByKeyWords(kw)
		if pos != context.NotFoundPos() {
			return pos
		}
	}
	return context.NotFoundPos()
}

func (i *Include) QueryAllByKeyWords(kw context.KeyWords) []context.Pos {
	poses := make([]context.Pos, 0)
	for _, config := range i.Configs {
		poses = append(poses, config.QueryAllByKeyWords(kw)...)
	}
	return poses
}

func (i *Include) Clone() context.Context {
	//configs := make([]*Config, 0)
	//for _, config := range i.Configs {
	//	configs = append(configs, config) // clone config's pointer
	//}
	return NewContext(context_type.TypeInclude, i.ContextValue)
}

func (i *Include) SetValue(v string) error {
	return errors.WithCode(code.V3ErrInvalidOperation, "cannot set include context's value")
}

func (i *Include) SetFather(ctx context.Context) error {
	i.fatherContext = ctx
	return nil
}

func (i *Include) HasChild() bool {
	return i.Len() > 0
}

func (i *Include) Len() int {
	return len(i.Configs)
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
	lines := make([]string, 0)
	if isDumping {
		lines = append(lines, "include "+i.ContextValue+";")
		return lines, nil
	}

	lines = append(lines, "# include <== "+i.ContextValue)
	for _, config := range i.Configs {
		configlines, err := config.ConfigLines(isDumping)
		if err != nil {
			return nil, err
		}
		lines = append(lines, configlines...)
	}
	return lines, nil
}

func (i *Include) matchConfigPath(path string) error {
	isMatch, err := filepath.Match(i.Value(), path)
	if err != nil {
		return errors.WithCode(code.V3ErrInvalidContext, "pattern(%s) match included config(%s) failed, cased by: %v", i.ContextValue, path, err)
	}
	if !isMatch {
		return errors.WithCode(code.V3ErrInvalidContext, "pattern(%s) cannot match included config(%s)", i.ContextValue, path)
	}
	return nil
}

func registerIncludeBuild() error {
	builderMap[context_type.TypeInclude] = func(value string) context.Context {
		return &Include{
			ContextValue:  value,
			Configs:       make(map[string]*Config),
			fatherContext: context.NullContext(),
		}
	}
	return nil
}

func registerIncludeParseFunc() error {
	pushStackParseFuncMap[context_type.TypeInclude] = func(data []byte, idx *int) context.Context {
		if matchIndexes := RegDirectiveWithValue.FindIndex(data[*idx:]); matchIndexes != nil { //nolint:nestif
			subMatch := RegDirectiveWithValue.FindSubmatch(data[*idx:])
			key := string(subMatch[1])
			value := string(subMatch[2])
			if key == string(context_type.TypeInclude) {
				*idx += matchIndexes[len(matchIndexes)-1]
				return NewContext(context_type.TypeInclude, value)
			}
		}
		return context.NullContext()
	}
	return nil
}
