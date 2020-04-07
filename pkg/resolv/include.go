package resolv

import (
	"encoding/json"
	"path/filepath"
)

type Include struct {
	BasicContext
	Key     *Key     `json:"-"`
	Comment *Comment `json:"-"`
	confPWD string   `json:"conf_pwd"`
}

func (i *Include) MarshalJSON() ([]byte, error) {
	includes := struct {
		ConfPWD  string `json:"conf_pwd"`
		Paths    string `json:"paths"`
		Includes []struct {
			Path    string   `json:"path"`
			Include []Parser `json:"include"`
		} `json:"includes"`
	}{ConfPWD: i.confPWD, Paths: i.Key.Value}

	for _, child := range i.Children {
		includes.Includes = append(includes.Includes, struct {
			Path    string   `json:"path"`
			Include []Parser `json:"include"`
		}{Path: child.(*Config).Value, Include: child.(*Config).Children})
	}

	return json.Marshal(includes)
}

func (i *Include) UnmarshalJSON(b []byte) error {
	includes := &struct {
		ConfPWD  string `json:"conf_pwd"`
		Paths    string `json:"paths"`
		Includes []struct {
			Path    string   `json:"path"`
			Include []Parser `json:"include"`
		} `json:"includes"`
	}{}

	err := json.Unmarshal(b, &includes)
	if err != nil {
		return err
	}

	for _, include := range includes.Includes {
		i.Add(NewConf(include.Include, include.Path))
	}

	i.Name = "include"
	i.Value = includes.Paths
	i.confPWD = includes.ConfPWD
	i.Key = NewKey("include", i.Value)
	i.Comment = NewComment("include "+i.Value, false)
	return nil
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

	paths, err := filepath.Glob(filepath.Join(i.confPWD, i.Value))
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
		Key:     NewKey("include", path),
		Comment: NewComment("include "+path, false),
		confPWD: dir,
	}

	err := include.load()
	if err != nil {
		return nil
	}

	return include
}
