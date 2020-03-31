package resolv

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type Key struct {
	Name  string
	Value string
}

func (k *Key) String() []string {
	if k.Value == "" {
		return []string{k.Name + ";\n"}
	} else if !inString(k.Value, "\"") && (inString(k.Value, ";") || inString(k.Value, "#")) {
		return []string{k.Name + " \"" + k.Value + "\";\n"}
	}
	return []string{k.Name + " " + k.Value + ";\n"}
}

func (k *Key) UnmarshalJSON(b []byte) error {
	s := string(bytes.Trim(b, "\""))
	reg := regexp.MustCompile(`\s*{\s*\"(\S+)\"\s*:\s*\"(\S*)\"\s*}\s*`)
	m := reg.FindStringIndex(s)
	if m != nil {
		ms := reg.FindStringSubmatch(s)
		k.Name = ms[1]
		if len(ms) > 2 {
			k.Value = ms[2]
		} else {
			k.Value = ""
		}
		return nil
	} else {
		return fmt.Errorf("%s is not a key", s)
	}
}

func (k *Key) MarshalJSON() ([]byte, error) {
	key := jsonFormat(k.Name)
	value := jsonFormat(k.Value)
	//var stamp = fmt.Sprintf("{\"%s\": \"%s\"}", k.Name, k.Value)
	var stamp = fmt.Sprintf("{\"%s\": \"%s\"}", key, value)
	//fmt.Print(stamp)
	return []byte(stamp), nil
}

func jsonFormat(s string) string {
	var i int
	for i = 0; i < len(s); i++ {
		if isSpecial(s[i]) || isNewLine(s, i) {
			break
		}
	}

	if i >= len(s) {
		return s
	}

	b := make([]byte, 2*len(s)-i)
	copy(b, s[:i])
	j := i
	for ; i < len(s); i++ {
		if isSpecial(s[i]) {
			b[j] = '\\'
			j++
		} else if isNewLine(s, i) {
			b[j] = '\\'
			j++
			b[j] = 'r'
			j++
			i++
		}
		b[j] = s[i]
		j++
	}
	return string(b[:j])
}

func isNewLine(s string, i int) bool {
	if i+1 >= len(s) {
		return false
	}

	return s[i:i+1] == "\n"
}

func isSpecial(b byte) bool {
	//return b == []byte("\"")[0] || b == []byte("'")[0] || b == []byte("`")[0]
	//return b == []byte("\"")[0] || b == []byte(":")[0]
	return b == []byte("\"")[0]
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
