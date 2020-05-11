package resolv

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

func (i *Include) Filter(kw KeyWords) (parsers []Parser) {
	//for _, conf := range i.Children {
	//	for _, child := range conf.(*Config).Children {
	//		if tmpParsers := child.Filter(kw); tmpParsers != nil {
	//			parsers = append(parsers, tmpParsers...)
	//		}
	//	}
	//}
	//return
	return i.subFilter(parsers, kw)
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

func (i *Include) load() error {

	paths, err := filepath.Glob(filepath.Join(i.ConfPWD, i.Value))
	if err != nil {
		return err
	}

	for _, path := range paths {

		conf, lerr := load(path)
		if lerr != nil {
			return lerr
		}

		//if len(paths) == 1 && paths[0] == path {
		//	i.Add(conf.Children...)
		//} else {
		//	sub := NewInclude("", path)
		//	sub.Add(conf.Children...)
		//	i.Add(sub)
		//}

		i.Add(conf)

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

func NewInclude(dir, paths string) *Include {
	include := &Include{
		BasicContext: BasicContext{
			Name:     TypeInclude,
			Value:    paths,
			Children: nil,
		},
		Key:     NewKey(TypeInclude, paths),
		Comment: NewComment(fmt.Sprintf("%s %s", TypeInclude, paths), false),
		ConfPWD: dir,
	}

	err := include.load()
	if err != nil {
		return nil
	}

	return include
}
