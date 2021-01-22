package configuration

import (
	"encoding/json"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration/parser"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/dumper"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/loop_preventer"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_type"
	"sync"
)

// keyword string: <parser type>[':sep: <value string>', ':sep: :reg: <value regexp>']
//
// e.g. for Nginx Config keyword string:
//     1) server
//     2) location:sep: :reg: \^\~\s+\/
//     3) key:sep: server_name test1\.com
//     4) comment:sep: :reg: .*
type Configuration interface {
	Queryer
	InsertByKeyword(insertParser parser.Parser, keyword string) error
	InsertByQueryer(insertParser parser.Parser, queryer Queryer) error
	//InsertByIndex(insertParser parser.Parser, targetContext parser.Context, index int) error
	RemoveByKeyword(keyword string) error
	RemoveByQueryer(queryer Queryer) error
	//RemoveByIndex(targetContext parser.Context, index int) error
	ModifyByKeyword(modifyParser parser.Parser, keyword string) error
	ModifyByQueryer(modifyParser parser.Parser, queryer Queryer) error
	//ModifyByIndex(modifyParser parser.Parser, targetContext parser.Context, index int) error
	View() []byte
	Json() []byte
	StatisticsByJson() []byte
	Dump() map[string][]byte
	//setConfig(config *parser.Config)
}

type configuration struct {
	config        *parser.Config
	rwLocker      *sync.RWMutex
	loopPreventer loop_preventer.LoopPreventer
}

func (c *configuration) InsertByKeyword(insertParser parser.Parser, keyword string) error {
	c.rwLocker.Lock()
	defer c.rwLocker.Unlock()
	parserKeyword, err := parseKeyword(keyword)
	if err != nil {
		return err
	}
	target, idx := c.config.Query(parserKeyword)
	//return target.Insert(insertParser, idx)
	return c.insertByIndex(insertParser, target, idx)
}

func (c *configuration) InsertByQueryer(insertParser parser.Parser, queryer Queryer) error {
	c.rwLocker.Lock()
	defer c.rwLocker.Unlock()
	//return queryer.fatherContext().Insert(insertParser, queryer.index())
	return c.insertByIndex(insertParser, queryer.fatherContext(), queryer.index())
}

func (c *configuration) insertByIndex(insertParser parser.Parser, targetContext parser.Context, index int) error {
	if insertParser.GetType() == parser_type.TypeConfig {
		err := c.loopPreventer.CheckLoopPrevent(targetContext.GetPosition(), insertParser.GetPosition())
		if err != nil {
			return err
		}
	}
	return targetContext.Insert(insertParser, index)
}

//func (c *configuration) InsertByIndex(insertParser parser.Parser, targetContext parser.Context, index int) error {
//	c.rwLocker.Lock()
//	defer c.rwLocker.Unlock()
//	return targetContext.Insert(insertParser, index)
//}

func (c *configuration) RemoveByKeyword(keyword string) error {
	c.rwLocker.Lock()
	defer c.rwLocker.Unlock()
	parserKeyword, err := parseKeyword(keyword)
	if err != nil {
		return err
	}
	target, idx := c.config.Query(parserKeyword)
	return target.Remove(idx)
}

func (c *configuration) RemoveByQueryer(queryer Queryer) error {
	c.rwLocker.Lock()
	defer c.rwLocker.Unlock()
	return queryer.fatherContext().Remove(queryer.index())
}

func (c *configuration) removeByIndex(targetContext parser.Context, index int) error {
	if targetContext.GetType() == parser_type.TypeInclude {
		removeParser, err := targetContext.GetChild(index)
		if err != nil {
			return err
		}
		err = c.loopPreventer.RemoveRoute(targetContext.GetPosition(), removeParser.GetPosition())
		if err != nil {
			return err
		}
	}
	return targetContext.Remove(index)
}

//func (c *configuration) RemoveByIndex(targetContext parser.Context, index int) error {
//	c.rwLocker.Lock()
//	defer c.rwLocker.Unlock()
//	return targetContext.Remove(index)
//}

func (c *configuration) ModifyByKeyword(modifyParser parser.Parser, keyword string) error {
	c.rwLocker.Lock()
	defer c.rwLocker.Unlock()
	parserKeyword, err := parseKeyword(keyword)
	if err != nil {
		return err
	}
	target, idx := c.config.Query(parserKeyword)
	return target.Modify(modifyParser, idx)
}

func (c *configuration) ModifyByQueryer(modifyParser parser.Parser, queryer Queryer) error {
	c.rwLocker.Lock()
	defer c.rwLocker.Unlock()
	return queryer.fatherContext().Modify(modifyParser, queryer.index())
}

//func (c *configuration) ModifyByIndex(modifyParser parser.Parser, targetContext parser.Context, index int) error {
//	c.rwLocker.Lock()
//	defer c.rwLocker.Unlock()
//	return targetContext.Modify(modifyParser, index)
//}

func (c *configuration) Query(keyword string) (Queryer, error) {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	parserKeyword, err := parseKeyword(keyword)
	if err != nil {
		return nil, err
	}
	foundCtx, index := c.config.Query(parserKeyword)
	if foundCtx == nil {
		return nil, nginx.ErrNotFound
	}
	foundParser, err := foundCtx.GetChild(index)
	if err != nil {
		return nil, err
	}
	return &queryer{
		Parser:    foundParser,
		fatherCtx: foundCtx,
		selfIndex: index,
	}, nil
}

func (c *configuration) QueryAll(keyword string) ([]Queryer, error) {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	parserKeyword, err := parseKeyword(keyword)
	if err != nil {
		return nil, err
	}
	queryers := make([]Queryer, 0)
	for foundCtx, indexes := range c.config.QueryAll(parserKeyword) {
		if indexes == nil {
			continue
		}
		for _, index := range indexes {
			foundParser, err := foundCtx.GetChild(index)
			if err != nil {
				return nil, err
			}
			queryers = append(queryers, &queryer{
				Parser:    foundParser,
				fatherCtx: foundCtx,
				selfIndex: index,
			})
		}
	}
	return queryers, nil
}

func (c configuration) Self() parser.Parser {
	return c.config
}

func (c configuration) fatherContext() parser.Context {
	return nil
}

func (c configuration) index() int {
	return 0
}

func (c *configuration) View() []byte {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	return c.config.Bytes()
}

func (c *configuration) Json() []byte {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	data, err := json.Marshal(c.config)
	if err != nil {
		return nil
	}
	return data
}

func (c *configuration) StatisticsByJson() []byte {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	statistician := NewStatistician(c)
	statistics := statistician.Statistics()
	data, err := json.Marshal(statistics)
	if err != nil {
		return nil
	}
	return data
}

func (c *configuration) Dump() map[string][]byte {
	c.rwLocker.RLock()
	defer c.rwLocker.RUnlock()
	d := dumper.NewDumper(c.config.GetValue())
	_ = c.config.Dump(d)
	return d.ReadAll()
}

func (c *configuration) setConfig(config *parser.Config) {
	c.rwLocker.Lock()
	defer c.rwLocker.Unlock()
	c.config = config
}

func NewConfiguration(config *parser.Config, preventer loop_preventer.LoopPreventer) Configuration {
	return &configuration{
		rwLocker:      new(sync.RWMutex),
		loopPreventer: preventer,
		config:        config,
	}
}
