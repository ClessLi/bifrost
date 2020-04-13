package resolv

import (
	"strings"
)

type Key struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (k *Key) String() []string {
	if k.Value == "" {
		return []string{k.Name + ";\n"}
	} else if !inString(k.Value, "\"") && (inString(k.Value, ";") || inString(k.Value, "#")) {
		return []string{k.Name + " \"" + k.Value + "\";\n"}
	}
	return []string{k.Name + " " + k.Value + ";\n"}
}

func inString(str string, s string) bool {
	if strings.Index(str, s) >= 0 {
		return true
	}
	return false
}

func NewKey(name, value string) *Key {
	return &Key{
		Name:  name,
		Value: value,
	}
}
