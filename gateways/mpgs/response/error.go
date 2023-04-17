package response

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	Err    *ErrorBody `json:"error,omitempty"`
	Result string     `json:"result,omitempty"`
}

// Error implements error
func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return fmt.Sprintf(`%s`, string(b))
}

type ErrorBody struct {
	Cause          string `json:"cause,omitempty"`
	Explanation    string `json:"explanation,omitempty"`
	Field          string `json:"field,omitempty"`
	SupportCode    string `json:"supportCode,omitempty"`
	ValidationType string `json:"validationType,omitempty"`
}

var _ error = (*Error)(nil)
