package local

import (
	"encoding/json"
	"fmt"

	"github.com/ClessLi/bifrost/internal/pkg/code"
)

type JSONError struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Code    int    `json:"code,omitempty"`
	// Caller  string `json:"caller,omitempty"`
}

func ToJSONError(err error) *JSONError {
	if err == nil {
		return nil
	}
	j := &JSONError{}
	e := json.Unmarshal([]byte(fmt.Sprintf("%#-v", err)), &[]*JSONError{j})
	if e != nil {
		j.Message = ""
		j.Error = err.Error()
		j.Code = 0
	}

	return j
}

func (e *JSONError) ToError() error {
	if e.Code == 0 {
		return code.New(e.Message)
	}

	return code.WithCode(e.Code, "%s", e.Message)
}
