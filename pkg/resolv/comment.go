package resolv

import "regexp"

type Comment struct {
	Comments string `json:"comments"`
	Inline   bool   `json:"inline"`
}

func (cmt *Comment) String() []string {
	return []string{"# " + cmt.Comments + "\n"}
}

func (cmt *Comment) Filter(kw KeyWords) (parsers []Parser) {
	if !kw.IsReg {
		if kw.Type == "comments" && kw.Value == cmt.Comments {
			parsers = append(parsers, cmt)
		} else {
			parsers = nil
		}
	} else {
		if kw.Type == "comments" && regexp.MustCompile(kw.Value).MatchString(cmt.Comments) {
			parsers = append(parsers, cmt)
		} else {
			parsers = nil
		}
	}
	return
}

func NewComment(value string, inline bool) *Comment {
	return &Comment{
		Comments: value,
		Inline:   inline,
	}
}
