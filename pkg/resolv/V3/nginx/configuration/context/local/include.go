package local

import (
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/dominikbraun/graph"
	"github.com/marmotedu/errors"
	"path/filepath"
)

type Include struct {
	ContextValue string             `json:"value,omitempty"`
	Configs      map[string]*Config `json:"param,omitempty"`

	fatherContext context.Context
}

func (i *Include) FatherConfig() (*Config, error) {
	fatherCtx := i.Father()
	fatherConfig, ok := fatherCtx.(*Config)
	for !ok {
		if fatherCtx.Type() == context_type.TypeErrContext {
			return nil, fatherCtx.(*context.ErrorContext).AppendError(errors.WithCode(code.V3ErrInvalidContext, "found a error config context")).Error()
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

	for _, config := range configs {
		// match config path
		if config.ConfigPath == nil {
			return errors.WithCode(code.V3ErrInvalidContext, "config with no ConfigPath")
		}

		err = i.matchConfigPath(config.RelativePath())
		if err != nil {
			err = i.matchConfigPath(config.FullPath())
			if err != nil {
				return err
			}
		}

		// clone inserted config
		//clone := config.Clone().(*Config)

		err = fatherConfig.AddConfig(config)
		if err != nil {
			if !errors.Is(err, graph.ErrVertexAlreadyExists) {
				return err
			}
			config, err = fatherConfig.GetConfig(config.FullPath())
			if err != nil {
				return err
			}
		}

		// config graph add edge
		err = fatherConfig.IncludeConfig(config)
		if err != nil {
			return err
		}
		i.Configs[config.FullPath()] = config

	}
	return nil
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
	for _, config := range configs {
		_, has := i.Configs[config.FullPath()]
		if !has {
			continue
		}
		delete(i.Configs, config.FullPath())
		err = fatherConfig.RemoveEdge(fatherConfig, config)
		if err != nil {
			return err
		}
		// TODO:删除本地文件
	}
	return nil
}

func (i *Include) Modify(ctx context.Context, idx int) context.Context {
	return context.ErrContext(errors.WithCode(code.V3ErrInvalidOperation, "include cannot modify by index"))
}

func (i *Include) ModifyConfig(configs ...*Config) error {
	if configs == nil {
		return errors.WithCode(code.V3ErrInvalidContext, "null config")
	}
	err := i.RemoveConfig(configs...)
	if err != nil {
		return err
	}
	err = i.InsertConfig(configs...)
	return err
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
		if pos != context.NullPos() {
			return pos
		}
	}
	return context.NullPos()
}

func (i *Include) QueryAllByKeyWords(kw context.KeyWords) []context.Pos {
	poses := make([]context.Pos, 0)
	for _, config := range i.Configs {
		poses = append(poses, config.QueryAllByKeyWords(kw)...)
	}
	return poses
}

func (i *Include) Clone() context.Context {
	configs := make(map[string]*Config)
	for path, config := range i.Configs {
		configs[path] = config.Clone().(*Config)
	}
	return &Include{
		ContextValue:  i.ContextValue,
		Configs:       configs,
		fatherContext: i.Father(),
	}
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
	return i.ContextValue
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
	isMatch, err := filepath.Match(i.ContextValue, path)
	if err != nil {
		return errors.WithCode(code.V3ErrInvalidContext, "pattern(%s) match included config(%s) failed, cased by: %v", i.ContextValue, path, err)
	}
	if !isMatch {
		return errors.WithCode(code.V3ErrInvalidContext, "pattern(%s) cannot match included config(%s)", i.ContextValue, path, err)
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
			*idx += matchIndexes[len(matchIndexes)-1]
			key := string(subMatch[1])
			value := string(subMatch[2])
			if key == string(context_type.TypeInclude) {
				return NewContext(context_type.TypeInclude, value)
			}
		}
		return context.NullContext()
	}
	return nil
}
