package nginx

import (
	"regexp"
	"strings"
)

type Comment struct {
	Comments string `json:"comments"`
	Inline   bool   `json:"inline"`
}

func (cmt Comment) String() []string {
	return cmt.string(nil, 0)
}

func (cmt *Comment) string(_ *Caches, deep int) []string {
	return []string{strings.Repeat(INDENT, deep) + "# " + cmt.Comments + "\n"}
}

func (cmt *Comment) QueryAll(pType parserType, isRec bool, values ...string) []Parser {
	kw, err := newKW(pType, values...)
	if err != nil {
		return nil
	}
	kw.IsRec = isRec
	return cmt.QueryAllByKeywords(*kw)
}

func (cmt *Comment) QueryAllByKeywords(kw Keywords) (parsers []Parser) {
	if parser := cmt.QueryByKeywords(kw); parser != nil {
		parsers = append(parsers, parser)
	}
	return
}

func (cmt *Comment) Query(pType parserType, isRec bool, values ...string) Parser {
	kw, err := newKW(pType, values...)
	if err != nil {
		return nil
	}
	kw.IsRec = isRec
	return cmt.QueryByKeywords(*kw)
}

func (cmt *Comment) QueryByKeywords(kw Keywords) (parser Parser) {
	if !kw.IsReg {
		if kw.Type == TypeComment && kw.Value == cmt.Comments {
			parser = cmt
		} else {
			parser = nil
		}
	} else {
		if kw.Type == TypeComment && regexp.MustCompile(kw.Value).MatchString(cmt.Comments) {
			parser = cmt
		} else {
			parser = nil
		}
	}
	return
}

func (cmt *Comment) BitSize(_ Order, _ int) byte {
	return 0
}

func (cmt *Comment) BitLen(_ Order) int {
	return 0
}

func (cmt *Comment) Size(_ Order) int {
	return 0
}

func NewComment(value string, inline bool) *Comment {
	return &Comment{
		Comments: value,
		Inline:   inline,
	}
}
