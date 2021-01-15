package loader

import (
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/context"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/context/parser"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/parser_position"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/parser_type"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Loader interface {
	LoadFromFilePath(path string) (context.Context, error)
	LoadFromJsonBytes(data []byte) (context.Context, error)
}

type loader struct {
	workDir string
	cacher  LoadCacher
}

func (l *loader) LoadFromFilePath(path string) (context.Context, error) {
	configAbsPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	if l.cacher == nil {
		l.cacher = NewLoadCacher(configAbsPath)
	} else {
		if config := l.cacher.GetConfig(configAbsPath); config != nil {
			return config, nil
		}
	}
	file, err := os.Open(configAbsPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	configData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var parserErr error
	index := 0
	parsers := make([]parser.Parser, 0)
	configDeep := 0

	position := parser_position.NewParserPosition(configAbsPath, configDeep, 0)
	config := context.NewContext("", parser_type.TypeConfig, position)

	parseContext := func(reg *regexp.Regexp) bool {
		var (
			contextValue    string
			contextType     parser_type.ParserType
			contextPosition = parser_position.NewParserPosition(configAbsPath, configDeep, 0)
		)

		if matchIndexes := reg.FindIndex(configData[index:]); matchIndexes != nil {
			configDeep++
			switch reg {
			case RegEventsHead:
				contextType = parser_type.TypeEvents
			case RegHttpHead:
				contextType = parser_type.TypeHttp
			case RegStreamHead:
				contextType = parser_type.TypeStream
			case RegServerHead:
				contextType = parser_type.TypeServer
			case RegLocationHead:
				contextType = parser_type.TypeLocation
				contextValue = string(reg.FindSubmatch(configData[index:])[1])
			case RegIfHead:
				contextType = parser_type.TypeIf
				contextValue = string(reg.FindSubmatch(configData[index:])[1])
			case RegUpstreamHead:
				contextType = parser_type.TypeUpstream
				contextValue = string(reg.FindSubmatch(configData[index:])[1])
			case RegGeoHead:
				contextType = parser_type.TypeGeo
				contextValue = string(reg.FindSubmatch(configData[index:])[1])
			case RegMapHead:
				contextType = parser_type.TypeMap
				contextValue = string(reg.FindSubmatch(configData[index:])[1])
			case RegLimitExceptHead:
				contextType = parser_type.TypeLimitExcept
				contextValue = string(reg.FindSubmatch(configData[index:])[1])
			case RegTypesHead:
				contextType = parser_type.TypeTypes
			default:
				configDeep--
				return false
			}
			ctx := context.NewContext(contextValue, contextType, contextPosition)
			if ctx != nil {
				parsers = append([]parser.Parser{ctx}, parsers...)
				index += matchIndexes[len(matchIndexes)-1]
				return true
			}
		}
		return false
	}

	parseContextEnd := func() bool {
		if matchIndexes := RegContextEnd.FindIndex(configData[index:]); matchIndexes != nil {
			index += matchIndexes[len(matchIndexes)-1]
			if lowerContext, ok := checkContext(parsers); ok {
				parsers = append(parsers[:0], parsers[1:]...)
				if upperContext, ok := checkContext(parsers); ok {
					parserErr = upperContext.Insert(lowerContext, upperContext.Len())
					if parserErr != nil {
						return false
					}
				} else {
					parserErr = config.Insert(lowerContext, config.Len())
					if parserErr != nil {
						return false
					}
				}
			}
			configDeep--
			return true
		}
		return false
	}

	parseComment := func() bool {
		if subMatch := RegCommentHead.FindSubmatch(configData[index:]); len(subMatch) == 3 {
			matchIndexes := RegCommentHead.FindIndex(configData[index:])
			index += matchIndexes[len(matchIndexes)-1]
			cmt := parser.NewComment(string(subMatch[2]), !strings.Contains(string(subMatch[1]), "\n") && index != 0, position)
			if ctx, ok := checkContext(parsers); ok {
				parserErr = ctx.Insert(cmt, ctx.Len())
				if parserErr != nil {
					return false
				}
			} else {
				parserErr = config.Insert(cmt, config.Len())
				if parserErr != nil {
					return false
				}
			}
			return true
		}
		return false
	}

	parseKeyOrInclude := func(reg *regexp.Regexp) bool {
		var (
			key   string
			value string
			p     parser.Parser
		)

		if matchIndexes := reg.FindIndex(configData[index:]); matchIndexes != nil {
			subMatch := reg.FindSubmatch(configData[index:])
			index += matchIndexes[len(matchIndexes)-1]
			switch reg {
			case RegKey:
				key = string(subMatch[1])
			case RegKeyValue:
				key = string(subMatch[1])
				value = string(subMatch[2])
			default:
				return false
			}

			if strings.EqualFold(string(subMatch[1]), parser_type.TypeInclude.String()) {
				p = context.NewContext(value, parser_type.TypeInclude, position)
				parserErr = l.loadIncludeConfigs(p.(*context.Include))
				if parserErr != nil {
					return false
				}
			} else {
				p = parser.NewKey(key, value, position)
			}

			if ctx, ok := checkContext(parsers); ok {
				parserErr = ctx.Insert(p, ctx.Len())
				if parserErr != nil {
					return false
				}
			} else {
				parserErr = config.Insert(p, config.Len())
				if parserErr != nil {
					return false
				}
			}
			return true
		}
		return false
	}

	for {
		switch {
		case parseContext(RegEventsHead),
			parseContext(RegHttpHead),
			parseContext(RegStreamHead),
			parseContext(RegServerHead),
			parseContext(RegLocationHead),
			parseContext(RegIfHead),
			parseContext(RegUpstreamHead),
			parseContext(RegGeoHead),
			parseContext(RegMapHead),
			parseContext(RegLimitExceptHead),
			parseContext(RegTypesHead),
			parseComment(),
			parseContextEnd(),
			parseKeyOrInclude(RegKeyValue),
			parseKeyOrInclude(RegKey):
			if parserErr != nil {
				return nil, parserErr
			}
			continue
		}
		break
	}
	return config, nil

}

func (l *loader) LoadFromJsonBytes(data []byte) (context.Context, error) {
	panic("implement me")
}

func (l *loader) loadIncludeConfigs(include *context.Include) error {
	configAbsPaths, err := filepath.Glob(filepath.Join(l.workDir, include.Value))
	if err != nil {
		return err
	}

	for _, path := range configAbsPaths {

		config, err := l.LoadFromFilePath(path)
		if err != nil {
			return err
		}

		err = include.Insert(config, include.Len())
		if err != nil {
			return err
		}

	}
	return nil
}

func NewLoader(workPath string) (Loader, error) {
	absPath, err := filepath.Abs(workPath)
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(absPath)
	return &loader{
		workDir: dir,
	}, nil
}

func checkContext(parsers []parser.Parser) (context.Context, bool) {
	if parsers != nil && len(parsers) > 0 {
		ctx, isContext := parsers[0].(context.Context)
		return ctx, isContext
	} else {
		return nil, false
	}
}
