package local

import (
	"encoding/json"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/dominikbraun/graph"
	"github.com/marmotedu/errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Loader interface {
	Load() (MainContext, error)
}

type parseFunc func(data []byte, idx *int) context.Context

type jsonLoader struct {
	unmarshaler *mainUnmarshaller
	jsonBytes   []byte
}

func (j *jsonLoader) Load() (MainContext, error) {
	err := json.Unmarshal(j.jsonBytes, j.unmarshaler)
	if err != nil {
		return nil, err
	}
	return j.unmarshaler.completedMain, nil
}

func JsonLoader(data []byte) Loader {
	return &jsonLoader{
		unmarshaler: &mainUnmarshaller{},
		jsonBytes:   data,
	}
}

type fileLoader struct {
	mainConfigAbsPath string
	configGraph       ConfigGraph
	contextStack      *contextStack
}

func (f *fileLoader) Load() (MainContext, error) {
	if !filepath.IsAbs(f.mainConfigAbsPath) {
		return nil, errors.Errorf("%s is not a absolute path", f.mainConfigAbsPath)
	}

	main, err := NewMain(f.mainConfigAbsPath)
	if err != nil {
		return nil, err
	}

	f.configGraph = main.graph()

	err = f.load(main.MainConfig())
	if err != nil {
		return nil, err
	}
	return main, main.renderGraph()
}

func (f *fileLoader) load(config *Config) error {
	data, err := os.ReadFile(configHash(config))
	if err != nil {
		return errors.Wrap(err, "read file failed")
	}

	// TODO: Optimize the matching mechanism for comments with no line breaks at the end
	data = append(data, []byte("\n")...) // avoid missing line breaks at the end of comments that prevent proper matching
	idx := 0
	stackIdx := len(f.contextStack.contexts)
	err = f.contextStack.push(config)
	if err != nil {
		return errors.Wrap(err, "push config context to stack failed")
	}

	for {

		isParsed := false
		if parseBlankLine(data, &idx) {
			continue
		}

		err = parseErrLine(data, &idx, config)
		if err != nil {
			return errors.Wrap(err, "has parsed an error line")
		}

		if parseBraceEnd(data, &idx) {
			_, err = f.contextStack.pop()
			if err != nil {
				return errors.Wrap(err, "quit context from stack failed")
			}
			continue
		}

		for _, parsefunc := range pushStackParseFuncMap {
			ctx := parsefunc(data, &idx)
			if ctx != context.NullContext() {
				father, err := f.contextStack.current()
				if err != nil {
					return errors.Wrap(err, "get father context failed")
				}
				err = father.Insert(ctx, father.Len()).Error()
				if err != nil {
					return errors.Wrap(err, "insert context failed")
				}

				err = f.contextStack.push(ctx)
				if err != nil {
					return errors.Wrap(err, "push context to stack failed")
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
					return errors.Wrap(err, "get father context failed")
				}
				err = father.Insert(ctx, father.Len()).Error()
				if err != nil {
					return errors.Wrap(err, "insert context failed")
				}
				isParsed = true
				break
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

	if e := new(commentsToContextConverter).Convert(config).Error(); e != nil {
		return e
	}
	for _, pos := range config.QueryAllByKeyWords(context.NewKeyWords(context_type.TypeInclude).SetCascaded(true).SetSkipQueryFilter(func(targetCtx context.Context) bool {
		return !targetCtx.IsEnabled() || targetCtx.Type() == context_type.TypeConfig // loading without the disabled context, and skipping the child config context
	})) {
		// load include configs
		if pos.Target().Type() == context_type.TypeInclude {
			include := pos.Target().(*Include)
			_, _, isEnabledInclude, _, err := include.parsePatternPath()
			if err != nil {
				return errors.Wrap(err, "load include configs failed")
			}
			if !isEnabledInclude { // skipping the disabled context
				continue
			}
			err = f.loadInclude(pos.Target().(*Include))
			if err != nil {
				return errors.Wrap(err, "load include configs failed")
			}
		}
	}

	return nil
}

func (f *fileLoader) loadInclude(include *Include) error {
	isAbsInclude := filepath.IsAbs(include.Value())
	var paths []string
	var err error
	if isAbsInclude {
		paths, err = filepath.Glob(include.Value())
		if err != nil {
			return err
		}
	} else {
		absPaths, err := filepath.Glob(filepath.Join(filepath.Dir(f.mainConfigAbsPath), include.Value()))
		if err != nil {
			return err
		}
		for _, absPath := range absPaths {
			path, err := filepath.Rel(filepath.Dir(f.mainConfigAbsPath), absPath)
			if err != nil {
				return err
			}
			paths = append(paths, path)
		}
	}

	// no match included configs
	if len(paths) == 0 {
		return nil
	}

	// new configs
	for _, path := range paths {
		newconfig := NewContext(context_type.TypeConfig, strings.TrimSpace(path)).(*Config)
		newconfig.ConfigPath, err = newConfigPath(f.configGraph, newconfig.Value())
		if err != nil {
			return err
		}
		// adding config into configGraph
		err = f.configGraph.AddConfig(newconfig)
		if err != nil {
			if errors.Is(err, graph.ErrVertexAlreadyExists) {
				continue
			}
			return err
		}
		// loading configs from file system
		err = f.load(newconfig)
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
	if ctx == nil {
		return errors.New("input a nil")
	}
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
		*idx += matchIndexes[len(matchIndexes)-1]
		line := 1
		for i := 0; i < *idx; i++ {
			if (data)[i] == []byte("\n")[0] {
				line++
			}
		}
		return errors.WithCode(code.ErrParseFailed, "parse failed at line %d of %s", line, configHash(config))
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
