package resolv

import (
	"path/filepath"
)

type Include struct {
	BasicContext
	Key     *Key
	Comment *Comment
	confDir string
}

func (i *Include) String() []string {
	var strs []string
	strs = append(strs, i.Comment.String()[0])
	for _, child := range i.Children {
		strs = append(strs, child.String()...)
	}
	strs = append(strs, "#End"+i.Comment.String()[0])

	return strs
}

func (i *Include) load() error {

	paths, err := filepath.Glob(filepath.Join(i.confDir, i.Value))
	if err != nil {
		return err
	}

	for _, path := range paths {

		conf, lerr := load(path)
		if lerr != nil {
			return lerr
		}

		if len(paths) == 1 && paths[0] == path {
			i.Add(conf.Children...)
		} else {
			sub := NewInclude("", path)
			sub.Add(conf.Children...)
			i.Add(sub)
		}

	}

	return nil
}

func NewInclude(dir, value string) *Include {
	include := &Include{
		BasicContext: BasicContext{
			Name:     "include",
			Value:    value,
			depth:    0,
			Children: nil,
		},
		Key: &Key{
			Name:  "include",
			Value: value,
		},
		Comment: &Comment{
			Comments: "include " + value,
			Inline:   false,
		},
		confDir: dir,
	}

	err := include.load()
	if err != nil {
		return nil
	}

	return include
}
