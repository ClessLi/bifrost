package nginx

import (
	"fmt"
	"path/filepath"
)

type Include struct {
	BasicContext `json:"include"`
	Key          *Key     `json:"tags"`
	Comment      *Comment `json:"comments"`
	ConfPWD      string   `json:"conf_pwd"`
}

func (i *Include) QueryAll(pType parserType, isRec bool, values ...string) []Parser {
	kw, err := newKW(pType, values...)
	if err != nil {
		return nil
	}
	kw.IsRec = isRec
	if isRec {
		return i.subQueryAll([]Parser{}, *kw)
	} else {
		parsers := make([]Parser, 0)
		for _, child := range i.Children {
			parsers = append(parsers, child.QueryAllByKeywords(*kw)...)
		}
		return parsers
	}
}

func (i *Include) QueryAllByKeywords(kw Keywords) (parsers []Parser) {
	//for _, conf := range i.Children {
	//	for _, child := range conf.(*Config).Children {
	//		if tmpParsers := child.QueryAllByKeywords(kw); tmpParsers != nil {
	//			parsers = append(parsers, tmpParsers...)
	//		}
	//	}
	//}
	//return
	if i.filter(kw) {
		parsers = append(parsers, i)
	}
	return i.subQueryAll(parsers, kw)
	//if kw.IsRec {
	//	return i.subQueryAll(parsers, kw)
	//} else {
	//	for _, child := range i.Children {
	//		parsers = append(parsers, child.QueryAllByKeywords(kw)...)
	//	}
	//	return
	//}
}

func (i *Include) Query(pType parserType, isRec bool, values ...string) Parser {
	kw, err := newKW(pType, values...)
	if err != nil {
		return nil
	}
	kw.IsRec = isRec
	return i.subQuery(*kw)
	//if isRec {
	//	return i.subQuery(*kw)
	//} else {
	//	for _, child := range i.Children {
	//		if parser := child.QueryByKeywords(*kw); parser != nil {
	//			return parser
	//		}
	//	}
	//	return nil
	//}
}

func (i *Include) QueryByKeywords(kw Keywords) Parser {
	if i.filter(kw) {
		return i
	}
	return i.subQuery(kw)
}

func (i *Include) Insert(indexParser Parser, pType parserType, values ...string) error {
	if values != nil {
		p, err := newParser(pType, values...)
		if err != nil {
			return err
		}

		return i.InsertByParser(indexParser, p)
	}
	return ParserControlNoParamError
}

func (i *Include) InsertByParser(indexParser Parser, contents ...Parser) error {
	for _, child := range i.Children {
		if err := child.(*Config).InsertByParser(indexParser, contents...); err == ParserControlIndexNotFoundError {
			continue
		} else {
			return err
		}
	}
	return ParserControlIndexNotFoundError
}

func (i *Include) Add(pType parserType, values ...string) error {
	if values != nil {
		parser, err := newParser(pType, values...)
		if err != nil {
			return err
		}
		i.AddByParser(parser)
		return nil

	}
	return ParserControlNoParamError
}

func (i *Include) AddByParser(contents ...Parser) {
	i.Children[len(i.Children)-1].(*Config).AddByParser(contents...)
}

func (i *Include) Remove(pType parserType, values ...string) error {
	i.RemoveByParser(i.QueryAll(pType, false, values...)...)
	return nil
}

func (i *Include) RemoveByParser(contents ...Parser) {
	for _, child := range i.Children {
		child.(*Config).RemoveByParser(contents...)
	}
}

func (i *Include) Modify(indexParser Parser, pType parserType, values ...string) error {
	if values != nil {

		ctx, err := newParser(pType, values...)
		if err != nil {
			return err
		}

		return i.ModifyByParser(indexParser, ctx)
	}
	return ParserControlNoParamError
}

func (i *Include) ModifyByParser(indexParser Parser, content Parser) error {
	for _, child := range i.Children {
		if err := child.(*Config).ModifyByParser(indexParser, content); err == ParserControlIndexNotFoundError {
			continue
		} else {
			return err
		}
	}
	return ParserControlIndexNotFoundError
}

func (i *Include) Params() (parsers []Parser) {
	for _, child := range i.Children {
		switch child.(type) {
		case *Config:
			parsers = append(parsers, child.(*Config).Params()...)
		}
	}
	return
}

func (i *Include) BitSize(_ Order, _ int) byte {
	return 0
}

func (i *Include) BitLen(_ Order) int {
	return 0
}

func (i *Include) Size(_ Order) int {
	return 0
}

func (i *Include) String() []string {
	var strs []string
	// 暂取消include对象自身信息
	//strs = append(strs, i.Comment.String()[0])
	for _, child := range i.Children {
		strs = append(strs, child.String()...)
	}
	// 暂取消include对象自身信息
	//strs = append(strs, "#End"+i.Comment.String()[0])

	return strs
}

func (i *Include) load(allConfigs *map[string]*Config, configCaches *[]string) error {
	//var newCaches []string
	//copy(newCaches, *configCaches)
	paths, err := filepath.Glob(filepath.Join(i.ConfPWD, i.Value))
	if err != nil {
		return err
	}

	for _, path := range paths {

		//conf, lerr := load(path, newCaches)
		relPath, relErr := filepath.Rel(i.ConfPWD, path)
		if relErr != nil {
			return relErr
		}

		conf, lerr := load(i.ConfPWD, relPath, allConfigs, configCaches)
		if lerr != nil {
			return lerr
		}

		//if len(paths) == 1 && paths[0] == path {
		//	i.AddByParser(conf.Children...)
		//} else {
		//	sub := NewInclude("", path)
		//	sub.AddByParser(conf.Children...)
		//	i.AddByParser(sub)
		//}

		i.AddByParser(conf)

	}

	return nil
}

func (i *Include) dump() ([]string, error) {
	for _, child := range i.Children {
		err := Save(child.(*Config))
		if err != nil {
			return nil, err
		}
	}

	return i.Key.String(), nil
}

func NewInclude(dir, paths string, allConfigs *map[string]*Config, configCaches *[]string) (*Include, error) {
	include := &Include{
		BasicContext: BasicContext{
			Name:     TypeInclude,
			Value:    paths,
			Children: nil,
		},
		Key:     NewKey(fmt.Sprintf("%s", TypeInclude), paths),
		Comment: NewComment(fmt.Sprintf("%s %s", TypeInclude, paths), false),
		ConfPWD: dir,
	}

	err := include.load(allConfigs, configCaches)
	if err != nil {
		return nil, err
	}

	return include, nil
}
