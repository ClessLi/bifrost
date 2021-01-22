package loader

import (
	"encoding/json"
	"fmt"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration/parser"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/loop_preventer"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_indention"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_position"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_type"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type Loader interface {
	LoadFromFilePath(path string) (parser.Context, loop_preventer.LoopPreventer, error)
	LoadFromJsonBytes(data []byte) (parser.Context, loop_preventer.LoopPreventer, error)
	GetConfigPaths() []string
}

type loader struct {
	workDir string
	cacher  LoadCacher
	locker  *sync.RWMutex
	loop_preventer.LoopPreventer
}

func (l *loader) LoadFromFilePath(path string) (parser.Context, loop_preventer.LoopPreventer, error) {
	l.locker.Lock()
	defer l.locker.Unlock()
	configAbsPath, err := filepath.Abs(path)
	if err != nil {
		return nil, nil, err
	}
	l.workDir = filepath.Dir(configAbsPath)
	l.cacher = NewLoadCacher(configAbsPath)
	l.LoopPreventer = loop_preventer.NewLoopPreverter(configAbsPath)
	//config := parser.NewContext(configAbsPath, parser_type.TypeConfig, parser_indention.NewIndention())
	ctx, err := l.loadFromConfigPosition(configAbsPath)
	return ctx, l.LoopPreventer, err
}

func (l *loader) loadFromConfigPosition(configAbsPath string) (parser.Context, error) {

	if config := l.cacher.GetConfig(configAbsPath); config != nil {
		return config, nil
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

	var parseErr error
	index := 0
	parsers := make([]parser.Parser, 0)
	configDeep := 0
	//positions := make([]parser_position.ParserPosition, 0)
	//positions = append(positions, configPos)
	indentions := make([]parser_indention.Indention, 0)
	indentions = append(indentions, parser_indention.NewIndention())

	config := parser.NewContext(configAbsPath, parser_type.TypeConfig, indentions[0])

	parseBlankLine := func() bool {
		if matchIndexes := RegBlankLine.FindIndex(configData[index:]); matchIndexes != nil {
			index += matchIndexes[len(matchIndexes)-1]
			return true
		}
		return false
	}

	parseErrorHead := func() bool {
		if matchIndexes := RegErrorHeed.FindIndex(configData[index:]); matchIndexes != nil {
			index += matchIndexes[0]
			return true
		}
		return false
	}

	parseContext := func(reg *regexp.Regexp, indention parser_indention.Indention) bool {
		var (
			contextValue string
			contextType  parser_type.ParserType
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
			ctx := parser.NewContext(contextValue, contextType, indention)
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
					parseErr = upperContext.Insert(lowerContext, upperContext.Len())
					if parseErr != nil {
						return false
					}

				} else {
					parseErr = config.Insert(lowerContext, config.Len())
					if parseErr != nil {
						return false
					}
				}
			}
			configDeep--
			return true
		}
		return false
	}

	parseComment := func(indention parser_indention.Indention) bool {
		if subMatch := RegCommentHead.FindSubmatch(configData[index:]); len(subMatch) == 3 {
			matchIndexes := RegCommentHead.FindIndex(configData[index:])
			index += matchIndexes[len(matchIndexes)-1]
			cmt := parser.NewComment(string(subMatch[2]), !strings.Contains(string(subMatch[1]), "\n") && index != 0, indention)
			if ctx, ok := checkContext(parsers); ok {
				parseErr = ctx.Insert(cmt, ctx.Len())
				if parseErr != nil {
					return false
				}
			} else {
				parseErr = config.Insert(cmt, config.Len())
				if parseErr != nil {
					return false
				}
			}
			return true
		}
		return false
	}

	parseKeyOrInclude := func(reg *regexp.Regexp, indention parser_indention.Indention) bool {
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
				p = parser.NewContext(value, parser_type.TypeInclude, indention)

			} else {
				p = parser.NewKey(key, value, indention)
			}

			if ctx, ok := checkContext(parsers); ok {
				parseErr = ctx.Insert(p, ctx.Len())
				if parseErr != nil {
					return false
				}
			} else {
				parseErr = config.Insert(p, config.Len())
				if parseErr != nil {
					return false
				}
			}

			if inc, ok := p.(*parser.Include); ok {
				inc.Position = parser_position.NewPosition(config.GetPosition())
				parseErr = l.loadIncludeConfigs(inc)
				if parseErr != nil {
					return false
				}
			}

			return true
		}
		return false
	}

	for {
		// initial indention
		for i := len(indentions); i <= configDeep; i++ {
			indentions = append(indentions, indentions[i-1].NextIndention())
		}

		indention := indentions[configDeep]

		switch {
		case parseErrorHead():
			parseErr = parseErrLine(configAbsPath, &configData, index)
		case parseBlankLine(),
			parseContext(RegEventsHead, indention),
			parseContext(RegHttpHead, indention),
			parseContext(RegStreamHead, indention),
			parseContext(RegServerHead, indention),
			parseContext(RegLocationHead, indention),
			parseContext(RegIfHead, indention),
			parseContext(RegUpstreamHead, indention),
			parseContext(RegGeoHead, indention),
			parseContext(RegMapHead, indention),
			parseContext(RegLimitExceptHead, indention),
			parseContext(RegTypesHead, indention),
			parseComment(indention),
			parseContextEnd(),
			parseKeyOrInclude(RegKeyValue, indention),
			parseKeyOrInclude(RegKey, indention):
			if parseErr != nil {
				return nil, parseErr
			}
			continue
		}
		break
	}
	err = l.cacher.SetConfig(config.(*parser.Config))
	if err != nil {
		return nil, err
	}
	return config, parseErr

}

func (l *loader) LoadFromJsonBytes(data []byte) (parser.Context, loop_preventer.LoopPreventer, error) {
	l.locker.Lock()
	defer l.locker.Unlock()
	unmarshaler := NewUnmarshaler()
	err := json.Unmarshal(data, unmarshaler)
	if err != nil {
		return nil, nil, err
	}
	l.workDir = filepath.Dir(unmarshaler.position.Id())
	l.cacher = unmarshaler.LoadCacher
	l.LoopPreventer = unmarshaler.LoopPreventer
	return unmarshaler.context, l.LoopPreventer, nil
}

func (l loader) GetConfigPaths() []string {
	l.locker.RLock()
	defer l.locker.RUnlock()
	return l.cacher.Keys()
}

func (l *loader) loadIncludeConfigs(include *parser.Include) error {
	configAbsPaths, err := filepath.Glob(filepath.Join(l.workDir, include.GetValue()))
	if err != nil {
		return err
	}

	for _, path := range configAbsPaths {

		// 校验引入的Config是否形成环路
		err := l.LoopPreventer.CheckLoopPrevent(include.GetPosition(), path)
		if err != nil {
			return err
		}

		// 加载引入的Config
		config, err := l.loadFromConfigPosition(path)
		if err != nil {
			return err
		}

		// Include引入Config
		err = include.Insert(config, include.Len())
		if err != nil {
			return err
		}

	}
	return nil
}

func NewLoader() Loader {
	return &loader{
		locker: new(sync.RWMutex),
	}
}

func checkContext(parsers []parser.Parser) (parser.Context, bool) {
	if parsers != nil && len(parsers) > 0 {
		ctx, isContext := parsers[0].(parser.Context)
		return ctx, isContext
	} else {
		return nil, false
	}
}

func parseErrLine(path string, configData *[]byte, index int) error {
	line := 1
	for i := 0; i < index; i++ {
		if (*configData)[i] == []byte("\n")[0] {
			line++
		}
	}
	return fmt.Errorf("%s at line %d of %s", ConfigParseError, line, path)
}
