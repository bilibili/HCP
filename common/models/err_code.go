package common

import (
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	SuccessMessage                = "SUCCESS"
	SuccessCode                   = "200"
	SuccessLoginSuccessMessage    = "登录成功"
	SuccessOutLoginSuccessMessage = "退出登录成功"

	ErrDatabaseQueryCode                = 400
	ErrDatabaseQueryMessage             = "数据库查询错误"
	ErrDataEmptyMessage                 = "数据为空"
	ErrServiceCode                      = "-500"
	ErrLoginAccountPwdEmptyMessage      = "账号或密码为空"
	ErrLoginInformationIncorrectMessage = "登录信息有误"
	ErrLoginUserEmptyMessage            = "用户不存在"
	ErrLoginUserNotUniqueMessage        = "用户不唯一"
	ErrLoginUserNamePwdMessage          = "用户名或密码错误"
	ErrNeedLoginUserMessage             = "请先登录"
	ErrNotJoinProjectMessage            = "请先加入项目"
)

// ReturnErr 返回错误
func ReturnErr(code int, reason, message string) error {
	return errors.New(code, reason, message)
}
