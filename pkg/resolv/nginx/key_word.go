package nginx

import (
	"fmt"
	"regexp"
	"strings"
)

type Keyword struct {
	Type  parserType
	Name  string
	Value string
	IsReg bool
}

type Keywords struct {
	Keyword
	ChildKWs []Keywords
	IsRec    bool
}

func NewKeyWords(contextType parserType, name, value string, isReg, isRec bool, subKWs ...interface{}) Keywords {
	switch contextType {
	case TypeKey, TypeComment:
	default:
		name = fmt.Sprintf("%s", contextType)
	}
	childKWs := make([]Keywords, 0)
	if subKWs != nil {
		for _, kw := range subKWs {
			switch kw.(type) {
			case Keywords:
				childKWs = append(childKWs, kw.(Keywords))
			}
		}
	} else {
		childKWs = nil
	}

	return Keywords{
		Keyword: Keyword{
			Type:  contextType,
			Name:  name,
			Value: value,
			IsReg: isReg,
		},
		ChildKWs: childKWs,
		IsRec:    isRec,
	}
}

func newKW(pType parserType, values ...string) (*Keywords, error) {
	var kws Keywords
	if values != nil {
		switch pType {
		case TypeComment:
			if ms := regexp.MustCompile(`^#+[ \r\t\f]*(.*)$`).FindStringSubmatch(values[0]); len(ms) == 2 {
				kws = NewKeyWords(pType, "", ms[1], true, false)
				return &kws, nil
			} else {
				return nil, ParserControlParamsError
			}
		case TypeKey:
			keyValue := ""
			kv := strings.Split(values[0], ":")
			if len(kv) > 1 {
				keyValue = strings.Join(kv[1:], ":")
			}
			keyName := kv[0]
			kws = NewKeyWords(pType, keyName, keyValue, true, false)
			return &kws, nil
		case TypeGeo, TypeIf, TypeLimitExcept, TypeLocation, TypeMap, TypeUpstream:
			kws = NewKeyWords(pType, "", values[0], true, false)
			if len(values) > 1 {
				values = values[1:]
			} else {
				values = nil
			}
		case TypeEvents, TypeHttp, TypeServer, TypeStream, TypeTypes:
			kws = NewKeyWords(pType, "", "", false, false)
		default:
			return nil, fmt.Errorf("unknown nginx context type: %s", pType)
		}
	} else {
		switch pType {
		case TypeEvents, TypeHttp, TypeServer, TypeStream, TypeTypes:
			kws = NewKeyWords(pType, "", "", false, false)
		default:
			return nil, fmt.Errorf("unknown nginx context type: %s", pType)
		}
	}

	if values != nil {
		for _, value := range values {
			if ms := regexp.MustCompile(`#+[ \r\t\f]*(.*?)`).FindStringSubmatch(value); len(ms) == 2 {
				kw := NewKeyWords(TypeComment, "", ms[1], true, false)
				kws.ChildKWs = append(kws.ChildKWs, kw)
			} else {
				kw, err := newKW(TypeKey, value)
				if err != nil {
					return nil, err
				}
				kws.ChildKWs = append(kws.ChildKWs, *kw)
			}
		}
	}
	return &kws, nil
}
