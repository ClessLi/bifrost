package resolv

import (
	"path/filepath"
)

type Include struct {
	BasicContext
	Key     *Key     `json:"-"`
	Comment *Comment `json:"-"`
	confDir string   `json:"conf_dir"`
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

func NewInclude(dir, path string) *Include {
	include := &Include{
		BasicContext: BasicContext{
			Name:     "include",
			Value:    path,
			depth:    0,
			Children: nil,
		},
		Key: &Key{
			Name:  "include",
			Value: path,
		},
		Comment: &Comment{
			Comments: "include " + path,
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
