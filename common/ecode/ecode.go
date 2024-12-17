package ecode

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/pkg/ecode"
)

var errcode = map[int]string{
	0:    "success",
	-401: "用户未登录",
	-403: "权限不足",
	-404: "请求参数错误",
	-405: "校验未通过",
	-500: "服务器内部错误",
}

func init() {
	ecode.Register(errcode)
}

// Code error struct
type Code struct {
	// 信息
	msg string
	// 编码
	code int
	// 详情
	details []interface{}
}

// Error return error message
func (e Code) Error() string {
	return fmt.Sprintf("%v: %v", e.code, e.msg)
}

// Code return error code
func (e Code) Code() int { return int(e.code) }

// Message return error message
func (e Code) Message() string {
	return fmt.Sprintf("%v", e.msg)
}

// Details return details.
func (e Code) Details() []interface{} {
	return e.details
}

// Equal for compatible.
// Deprecated: please use ecode.EqualError.
func (e Code) Equal(err error) bool { return EqualError(e, err) }

// EqualError equal error
func EqualError(code ecode.Codes, err error) bool {
	return ecode.Cause(err).Code() == code.Code()
}

// Int parse code int to error.
func NewECode(msgs ...string) *Code {
	code := -500
	size := len(msgs)
	switch size {
	case 1:
		msg := msgs[0]
		if cd, err := strconv.Atoi(msg); err == nil {
			code = cd
			if msgTmp, ok := errcode[code]; ok {
				msg = msgTmp
			}
		}
		//解析"-60001: NotFoundError"这种错误格式
		if smsg := strings.Split(msg, ": "); len(smsg) > 1 {
			if cd, err := strconv.Atoi(smsg[0]); err == nil {
				code = cd
				msg = smsg[1]
			}
		}
		return &Code{msg: msg, code: code}
	case 2:
		if cd, err := strconv.Atoi(msgs[0]); err == nil {
			code = cd
		}
		return &Code{msg: msgs[1], code: code}
	default:
		return &Code{msg: "服务器内部错误", code: code}
	}
}
