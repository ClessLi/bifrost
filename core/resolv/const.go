package resolv

const (
	DoubleQuotes = `\s*"[^"]*"`
	SingleQuotes = `\s*'[^']*'`
	Normal       = `\s*[^;\s]*`
	S1           = DoubleQuotes + `|` + SingleQuotes + `|` + Normal
	S            = `^\s*(` + S1 + `)\s*((?:` + S1 + `)+);`
)
