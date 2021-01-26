package configuration

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration/parser"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_type"
	"strings"
)

// keyword string: <parser type>[':sep: <value string>', ':sep: :reg: <value regexp>']
//
// e.g. for Nginx Config keyword string:
//     1) server
//     2) location:sep: :reg: \^\~\s+\/
//     3) key:sep: server_name test1\.com
//     4) comment:sep: :reg: .*
type Queryer interface {
	// keyword string: <parser type>[':sep: <value string>', ':sep: :reg: <value regexp>']
	//
	// e.g. for Nginx Config keyword string:
	//     1) server
	//     2) location:sep: :reg: \^\~\s+\/
	//     3) key:sep: server_name test1\.com
	//     4) comment:sep: :reg: .*
	Query(keyword string) (Queryer, error)
	// keyword string: <parser type>[':sep: <value string>', ':sep: :reg: <value regexp>']
	//
	// e.g. for Nginx Config keyword string:
	//     1) server
	//     2) location:sep: :reg: \^\~\s+\/
	//     3) key:sep: server_name test1\.com
	//     4) comment:sep: :reg: .*
	QueryAll(keyword string) ([]Queryer, error)
	Self() parser.Parser
	fatherContext() parser.Context
	index() int
}

type queryer struct {
	parser.Parser
	fatherCtx parser.Context
	selfIndex int
}

func (q queryer) Query(keyword string) (Queryer, error) {
	ctx, ok := q.Parser.(parser.Context)
	if !ok {
		return nil, nil
	}
	parserKeyword, err := parseKeyword(keyword)
	if err != nil {
		return nil, err
	}
	foundCtx, index := ctx.Query(parserKeyword)
	if foundCtx == nil {
		return nil, nginx.ErrNotFound
	}
	foundParser, err := foundCtx.GetChild(index)
	if err != nil {
		return nil, err
	}

	return &queryer{
		Parser:    foundParser,
		fatherCtx: foundCtx,
		selfIndex: index,
	}, nil
}

func (q queryer) QueryAll(keyword string) ([]Queryer, error) {
	ctx, ok := q.Parser.(parser.Context)
	if !ok {
		return nil, nil
	}
	parserKeyword, err := parseKeyword(keyword)
	if err != nil {
		return nil, err
	}
	queryers := make([]Queryer, 0)
	for foundCtx, indexes := range ctx.QueryAll(parserKeyword) {
		if indexes == nil {
			continue
		}
		for _, index := range indexes {
			foundParser, err := foundCtx.GetChild(index)
			if err != nil {
				return nil, err
			}
			queryers = append(queryers, &queryer{
				Parser:    foundParser,
				fatherCtx: foundCtx,
				selfIndex: index,
			})
		}
	}
	return queryers, nil
}

func (q queryer) Self() parser.Parser {
	return q.Parser
}

func (q queryer) fatherContext() parser.Context {
	return q.fatherCtx
}

func (q queryer) index() int {
	return q.selfIndex
}

func NewQueryer(fatherContext parser.Context, selfIndex int) (Queryer, error) {
	p, err := fatherContext.GetChild(selfIndex)
	if err != nil {
		return nil, err
	}
	return &queryer{
		Parser:    p,
		fatherCtx: fatherContext,
		selfIndex: selfIndex,
	}, nil
}

func parseKeyword(keyword string) (parser.KeyWords, error) {
	var (
		parserType parser_type.ParserType
		keyValue   string
		isReg      bool
	)

	// keyword
	//     "<ParserType>:sep: <key and value string>"
	//     "<ParserType>:sep: :reg: <key and value regexp>"
	kw := strings.Split(keyword, ":sep:")
	if len(kw) == 2 {
		parserType = parser_type.ParserType(kw[0])

		kv := strings.TrimSpace(kw[1])
		if kv[:5] == ":reg:" {
			isReg = true
			keyValue = strings.TrimSpace(kv[5:])
		} else {
			keyValue = strings.TrimSpace(kv)
		}
	} else if len(kw) == 1 {
		parserType = parser_type.ParserType(kw[0])
	} else {
		return nil, nginx.KeywordStringError
	}

	return parser.NewKeyWords(parserType, isReg, keyValue)
}
