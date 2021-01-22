package parser_position

type ParserPosition interface {
	Id() string
}

type parserPosition struct {
	id string
}

func (p parserPosition) Id() string {
	return p.id
}

func NewPosition(id string) ParserPosition {
	return &parserPosition{
		id: id,
	}
}

//func NewParserPosition(path string, configDeep, globalDeep int) ParserPosition {
//	return &parserPosition{
//		configAbsPath: path,
//		configDeep:    configDeep,
//		globalDeep:    globalDeep,
//		indent:        INDENT,
//		loopPreventer: loop_preventer.NewLoopPreverter(path),
//		locker:        new(sync.RWMutex),
//	}
//}
