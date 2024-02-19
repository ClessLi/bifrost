package local

import (
	"encoding/json"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/graph"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
	"os"
	"path/filepath"
	"regexp"
)

type Loader interface {
	Load() (*Main, error)
}

type parseFunc func(data []byte, idx *int) context.Context

type jsonLoader struct {
	unmarshaler *mainUnmarshaler
	jsonBytes   []byte
}

func (j *jsonLoader) Load() (*Main, error) {
	err := json.Unmarshal(j.jsonBytes, j.unmarshaler)
	if err != nil {
		return nil, err
	}
	return j.unmarshaler.completedMain, nil
}

func JsonLoader(data []byte) Loader {
	return &jsonLoader{
		unmarshaler: &mainUnmarshaler{unmarshalContext: new(_main)},
		jsonBytes:   data,
	}
}

type fileLoader struct {
	mainConfigAbsPath string
	configGraph       ConfigGraph
	contextStack      *contextStack
}

func (f *fileLoader) Load() (*Main, error) {
	if !filepath.IsAbs(f.mainConfigAbsPath) {
		return nil, errors.Errorf("%s is not a absolute path", f.mainConfigAbsPath)
	}

	mainCtx := NewContext(context_type.TypeMain, f.mainConfigAbsPath)
	err := mainCtx.Error()
	if err != nil {
		return nil, err
	}

	main, ok := mainCtx.(*Main)
	if !ok {
		return nil, errors.New("failed to build main context")
	}

	f.configGraph = main.ConfigGraph

	err = f.load(main.Config)

	return main, err
}

func (f *fileLoader) load(config *Config) error {
	data, err := os.ReadFile(config.FullPath())
	if err != nil {
		return err
	}

	idx := 0
	stackIdx := len(f.contextStack.contexts)
	err = f.contextStack.push(config)
	if err != nil {
		return err
	}

	for {

		isParsed := false
		if parseBlankLine(data, &idx) {
			continue
		}

		err = parseErrLine(data, &idx, config)
		if err != nil {
			return err
		}

		if parseBraceEnd(data, &idx) {
			_, err = f.contextStack.pop()
			if err != nil {
				return err
			}
			continue
		}

		for _, parsefunc := range pushStackParseFuncMap {
			ctx := parsefunc(data, &idx)
			if ctx != context.NullContext() {
				father, err := f.contextStack.current()
				if err != nil {
					return err
				}
				err = father.Insert(ctx, father.Len()).Error()
				if err != nil {
					return err
				}

				err = f.contextStack.push(ctx)
				if err != nil {
					return err
				}
				isParsed = true
				break
			}
		}
		if isParsed {
			continue
		}

		for _, parsefunc := range inStackParseFuncMap {
			ctx := parsefunc(data, &idx)
			if ctx != context.NullContext() {
				father, err := f.contextStack.current()
				if err != nil {
					return err
				}
				err = father.Insert(ctx, father.Len()).Error()
				if err != nil {
					return err
				}
				isParsed = true
				break
			}
			// load include configs
			if ctx.Type() == context_type.TypeInclude {
				err = f.loadInclude(ctx.(*Include))
				if err != nil {
					return err
				}
			}
		}
		if isParsed {
			continue
		}

		break
	}
	_, err = f.contextStack.pop()
	if err != nil {
		return err
	}
	if stackIdx != len(f.contextStack.contexts) {
		return errors.WithCode(code.ErrParseFailed, "context stack is not empty")
	}

	return nil
}

func (f *fileLoader) loadInclude(include *Include) error {
	isAbsInclude := filepath.IsAbs(include.Value())
	var paths []string
	var err error
	if isAbsInclude {
		paths, err = filepath.Glob(include.Value())
	} else {
		paths, err = filepath.Glob(filepath.Join(filepath.Dir(f.mainConfigAbsPath), include.Value()))
	}
	if err != nil {
		return err
	}

	// new configs
	newconfigs := make([]*Config, 0)
	includedconfigs := make([]*Config, 0)
	for _, path := range paths {
		var configpath context.ConfigPath
		if isAbsInclude {
			configpath, err = context.NewAbsConfigPath(path)
		} else {
			configpath, err = context.NewRelConfigPath(filepath.Dir(f.mainConfigAbsPath), path)
		}
		if err != nil {
			return err
		}

		// get config cache
		cache, err := f.configGraph.GetConfig(configpath.FullPath())
		if err == nil { // has cache
			includedconfigs = append(includedconfigs, cache)
			continue
		} else if !errors.Is(err, graph.ErrVertexNotExist) {
			return err
		}

		// build new config
		config, ok := NewContext(context_type.TypeConfig, path).(*Config)
		if !ok {
			return errors.Errorf("failed to build included config %s", path)
		}
		config.ConfigPath = configpath
		newconfigs = append(newconfigs, config)
	}
	includedconfigs = append(includedconfigs, newconfigs...)

	// include configs
	err = include.InsertConfig(includedconfigs...)
	if err != nil {
		return err
	}
	// load new configs
	for _, config := range newconfigs {
		err = f.load(config)
		if err != nil {
			return err
		}
	}
	return nil
}

func FileLoader(configpath string) Loader {
	return &fileLoader{
		mainConfigAbsPath: filepath.Clean(configpath),
		contextStack:      newContextStack(),
	}
}

type contextStack struct {
	contexts []context.Context
}

func (s *contextStack) current() (context.Context, error) {
	if len(s.contexts) == 0 {
		return context.NullContext(), errors.New("empty context stack")
	}
	return s.contexts[len(s.contexts)-1], nil
}

func (s *contextStack) pop() (context.Context, error) {
	ctx, err := s.current()
	if err != nil {
		return ctx, err
	}
	s.contexts = s.contexts[:len(s.contexts)-1]
	return ctx, nil
}

func (s *contextStack) push(ctx context.Context) error {
	err := ctx.Error()
	if err != nil {
		return err
	}
	s.contexts = append(s.contexts, ctx)
	return nil
}

func newContextStack() *contextStack {
	return &contextStack{contexts: make([]context.Context, 0)}
}

type parseFuncBuildOptions struct {
	regex           *regexp.Regexp
	contextType     context_type.ContextType
	valueMatchIndex int
}

func parseBlankLine(data []byte, idx *int) bool {
	if matchIndexes := RegBlankLine.FindIndex(data[*idx:]); matchIndexes != nil {
		*idx += matchIndexes[len(matchIndexes)-1]

		return true
	}

	return false
}

func parseErrLine(data []byte, idx *int, config *Config) error {
	if matchIndexes := RegErrorHeed.FindIndex(data[*idx:]); matchIndexes != nil {
		*idx += matchIndexes[0]
		line := 1
		for i := 0; i < *idx; i++ {
			if (data)[i] == []byte("\n")[0] {
				line++
			}
		}
		return errors.WithCode(code.ErrParseFailed, "parse failed at line %d of %s", line, config.FullPath())
	}
	return nil
}

func parseBraceEnd(data []byte, idx *int) bool {
	if matchIndexes := RegBraceEnd.FindIndex(data[*idx:]); matchIndexes != nil { //nolint:nestif
		*idx += matchIndexes[len(matchIndexes)-1]

		return true
	}

	return false
}
