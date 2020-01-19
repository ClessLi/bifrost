package resolv

type Comment struct {
	Comments string
	Inline   bool
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
