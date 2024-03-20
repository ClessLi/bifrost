package local

const (
	DoubleQuotes = `\s*"[^"]*"`
	SingleQuotes = `\s*'[^']*'`
	Normal       = `\s*[^;\s#]+`
	Abnormal     = `^[\n\t\f ]*;.*`
	LineBreak    = `\n`
	S1           = DoubleQuotes + `|` + SingleQuotes + `|` + Normal
	S            = `^\s*(` + S1 + `)\s+((?:` + S1 + `)+)\s*;`
)
