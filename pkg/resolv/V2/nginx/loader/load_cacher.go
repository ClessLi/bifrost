package loader

import (
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration/parser"
	"github.com/marmotedu/errors"
	"strings"
	"sync"
)

type LoadCacher interface {
	MainConfig() *parser.Config
	GetConfig(configName string) *parser.Config
	Keys() []string
	//CheckIncludeConfig(src, dst string) error
	SetConfig(config *parser.Config) error
}

type loadCache struct {
	mainConfig string
	cache      map[string]*parser.Config
	//loopPreventer loop_preventer.LoopPreventer
	rwLocker *sync.RWMutex
}

func (l loadCache) MainConfig() *parser.Config {
	l.rwLocker.RLock()
	defer l.rwLocker.RUnlock()
	return l.GetConfig(l.mainConfig)
}

func (l loadCache) Keys() []string {
	l.rwLocker.RLock()
	defer l.rwLocker.RUnlock()
	keys := make([]string, 0)
	for s := range l.cache {
		keys = append(keys, s)
	}
	return keys
}

func (l loadCache) GetConfig(configName string) *parser.Config {
	l.rwLocker.RLock()
	defer l.rwLocker.RUnlock()
	config, ok := l.cache[configName]
	if ok {
		return config
	}
	return nil
}

//func (l loadCache) CheckIncludeConfig(src, dst string) error {
//	return l.loopPreventer.CheckLoopPrevent(src, dst)
//}

func (l *loadCache) SetConfig(config *parser.Config) error {
	configName := config.GetValue()
	if strings.EqualFold(configName, "") {
		return errors.WithCode(code.ErrInvalidConfig, "get config name failed or null config name")
	}
	l.rwLocker.Lock()
	defer l.rwLocker.Unlock()

	l.cache[configName] = config
	return nil
}

func NewLoadCacher(configAbsPath string) LoadCacher {
	return &loadCache{
		mainConfig: configAbsPath,
		cache:      make(map[string]*parser.Config),
		//loopPreventer: loop_preventer.NewLoopPreverter(configAbsPath),
		rwLocker: new(sync.RWMutex),
	}
}
