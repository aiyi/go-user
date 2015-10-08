package errors

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("item not found")

const (
	ErrCodeOK                  = 0
	ErrCodeInternalServerError = 10000 // 内部服务器出错
	ErrCodeBadRequest          = 10001 // 输入参数不正确
	ErrCodeAuthTypeMissing     = 20000 // auth_type 缺失
	ErrCodeAuthTypeUnknown     = 20001 // auth_type 无法识别
	ErrCodeSessionTokenEncode  = 20002 // sessiontoken 编码出错
	ErrCodeAuthFailed          = 20003 // 认证失败
)

var ErrOK = &Error{
	ErrCode: ErrCodeOK,
	ErrMsg:  "success",
}
var ErrInternalServerError = &Error{
	ErrCode: ErrCodeInternalServerError,
	ErrMsg:  "internal server error",
}
var ErrBadRequest = &Error{
	ErrCode: ErrCodeBadRequest,
	ErrMsg:  "bad request",
}
var ErrAuthTypeMissing = &Error{
	ErrCode: ErrCodeAuthTypeMissing,
	ErrMsg:  "x-auth-type missing",
}
var ErrAuthTypeUnknown = &Error{
	ErrCode: ErrCodeAuthTypeUnknown,
	ErrMsg:  "x-auth-type unknown",
}
var ErrSessionTokenEncode = &Error{
	ErrCode: ErrCodeSessionTokenEncode,
	ErrMsg:  "token encoding failure",
}
var ErrAuthFailed = &Error{
	ErrCode: ErrCodeAuthFailed,
	ErrMsg:  "authentication failed",
}

type Error struct {
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

func NewError(ErrCode int, ErrMsg string) *Error {
	return &Error{
		ErrCode: ErrCode,
		ErrMsg:  ErrMsg,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("err_code: %d, err_msg: %s", e.ErrCode, e.ErrMsg)
}
