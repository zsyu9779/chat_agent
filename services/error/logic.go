package error

import (
	"encoding/json"
	"fmt"
)

type LogicError struct {
	code string
	msg  string
}

func (e LogicError) Error() string {
	b, _ := json.Marshal(map[string]interface{}{
		"code": e.code,
		"msg":  e.msg,
	})
	return string(b)
}

func (e LogicError) Code() string {
	return e.code
}

func (e LogicError) Msg() string {
	return e.msg
}

func NewLogicErr(code string, msg ...interface{}) error {
	return LogicError{
		code: code,
		msg:  fmt.Sprint(msg...),
	}
}
