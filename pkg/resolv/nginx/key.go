package nginx

import (
	"regexp"
	"strings"
)

type Key struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (k Key) String() []string {
	return k.string(nil, 0)
}

func (k *Key) string(_ *Caches, deep int) []string {
	ind := strings.Repeat(INDENT, deep)
	if k.Value == "" {
		return []string{ind + k.Name + ";\n"}
		//} else if !inString(k.Value, "\"") && (inString(k.Value, ";") || inString(k.Value, "#")) {
		//	return []string{k.Name + " \"" + k.Value + "\";\n"}
	}
	return []string{ind + k.Name + " " + k.Value + ";\n"}
}

func (k *Key) QueryAll(pType parserType, isRec bool, values ...string) []Parser {
	kw, err := newKW(pType, values...)
	if err != nil {
		return nil
	}
	kw.IsRec = isRec
	return k.QueryAllByKeywords(*kw)
}

func (k *Key) QueryAllByKeywords(kw Keywords) (parsers []Parser) {
	if parser := k.QueryByKeywords(kw); parser != nil {
		parsers = append(parsers, parser)
	}
	return
}

func (k *Key) Query(pType parserType, isRec bool, values ...string) Parser {
	kw, err := newKW(pType, values...)
	if err != nil {
		return nil
	}
	kw.IsRec = isRec
	return k.QueryByKeywords(*kw)
}

func (k *Key) QueryByKeywords(kw Keywords) (parser Parser) {
	if !kw.IsReg {
		if kw.Type == TypeKey && kw.Name == k.Name && (kw.Value == k.Value || kw.Value == `*`) {
			parser = k
		} else {
			parser = nil
		}
	} else {
		if kw.Type == TypeKey && regexp.MustCompile(kw.Name).MatchString(k.Name) && regexp.MustCompile(kw.Value).MatchString(k.Value) {
			parser = k
		} else {
			parser = nil
		}
	}
	return
}

func (k *Key) BitSize(_ Order, _ int) byte {
	return 0
}

func (k *Key) BitLen(_ Order) int {
	return 0
}

func (k *Key) Size(_ Order) int {
	return 0
}

//func inString(str string, s string) bool {
//	if strings.Index(str, s) >= 0 {
//		return true
//	}
//	return false
//}

func NewKey(name, value string) *Key {
	return &Key{
		Name:  name,
		Value: value,
	}
}
