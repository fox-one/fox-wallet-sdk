package sdk

import (
	"encoding/json"
	"log"
)

const (
	// RequestFailed request failed
	RequestFailed = 3000000
)

// Error error
type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Hint string `json:"hint,omitempty"`
}

// Error error message
func (sessionError Error) Error() string {
	str, err := json.Marshal(sessionError)
	if err != nil {
		log.Panicln(err)
	}
	return string(str)
}

func requestError(err error) *Error {
	e := &Error{
		Code: RequestFailed,
		Msg:  "request failed",
	}
	if err != nil {
		e.Hint = err.Error()
	}
	return e
}
