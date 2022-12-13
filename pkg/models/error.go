package models

import (
	"encoding/json"
	"fmt"
)

type ErrorCode struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

type ErrorMessage struct {
	Msg string `json:"msg"`
}

type ErrorCause struct {
	ErrorType    string `json:"errorType"`
	ErrorMessage string `json:"errorMessage"`
}

func (e ErrorCode) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf(`{"code": %v, "msg": "%s"}`, e.Code, e.Message)
	}
	return string(b)
}

func UnexpectedError(error string) *ErrorCode {
	return &ErrorCode{
		Code:    4000,
		Message: fmt.Sprintf("unexpected error. %v", error),
	}
}
