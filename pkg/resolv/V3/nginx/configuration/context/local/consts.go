package local

const (
	DoubleQuotes = `\s*"[^"]*"`
	SingleQuotes = `\s*'[^']*'`
	Normal       = `\s*[^;\s#]+`
	Abnormal     = `^[\n\r\t\f ]*;.*`
	LineBreak    = `[\n\r]`
	S1           = DoubleQuotes + `|` + SingleQuotes + `|` + Normal
	S            = `^\s*(` + S1 + `)\s+((?:` + S1 + `)+)\s*;`
)
