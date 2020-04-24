package resolv

type KeyWord struct {
	Type  string
	Name  string
	Value string
	IsReg bool
}

type KeyWords struct {
	KeyWord
	ChildKWs []KeyWords
}

func NewKeyWords(contextType, name, value string, isReg bool, subKWs ...interface{}) KeyWords {
	switch contextType {
	case "key", "comments":
	default:
		name = contextType
	}
	childKWs := make([]KeyWords, 0)
	if subKWs != nil {
		for _, kw := range subKWs {
			switch kw.(type) {
			case KeyWords:
				childKWs = append(childKWs, kw.(KeyWords))
			}
		}
	} else {
		childKWs = nil
	}

	return KeyWords{
		KeyWord: KeyWord{
			Type:  contextType,
			Name:  name,
			Value: value,
			IsReg: isReg,
		},
		ChildKWs: childKWs,
	}
}
