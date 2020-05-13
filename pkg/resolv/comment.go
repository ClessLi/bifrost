package resolv

import (
	"regexp"
)

type Comment struct {
	Comments string `json:"comments"`
	Inline   bool   `json:"inline"`
}

func (cmt *Comment) String() []string {
	return []string{"# " + cmt.Comments + "\n"}
}

func (cmt *Comment) QueryAll(kw KeyWords) (parsers []Parser) {
	if !kw.IsReg {
		if kw.Type == TypeComment && kw.Value == cmt.Comments {
			parsers = append(parsers, cmt)
		} else {
			parsers = nil
		}
	} else {
		if kw.Type == TypeComment && regexp.MustCompile(kw.Value).MatchString(cmt.Comments) {
			parsers = append(parsers, cmt)
		} else {
			parsers = nil
		}
	}
	return
}

func (cmt *Comment) Query(kw KeyWords) (parser Parser) {
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

func (cmt *Comment) BitSize(order Order, bit int) byte {
	return 0
}

func (cmt *Comment) BitLen(order Order) int {
	return 0
}

func (cmt *Comment) Size(order Order) int {
	return 0
}

func NewComment(value string, inline bool) *Comment {
	return &Comment{
		Comments: value,
		Inline:   inline,
	}
}
