package resolv

type Comment struct {
	Comments string `json:"comments"`
	Inline   bool   `json:"inline"`
}

func (cmt *Comment) String() []string {
	return []string{"# " + cmt.Comments + "\n"}
}

func NewComment(value string, inline bool) *Comment {
	return &Comment{
		Comments: value,
		Inline:   inline,
	}
}
