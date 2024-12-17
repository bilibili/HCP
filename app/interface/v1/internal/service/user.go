// Package service
/*
 * Copyright 2024-2025 Bilibili Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package service

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/bilibili/HCP/app/interface/v1/configs"
	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/bilibili/HCP/common/ecode"
	common "github.com/bilibili/HCP/common/models"
	"github.com/bilibili/HCP/utils"
	"github.com/go-kratos/kratos/pkg/log"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster/binding"
)

// UserLoginReq 用户登录请求参数
type UserLoginReq struct {
	// 用户名
	UserName string `json:"username" form:"username" validate:"required" `
	// 密码
	Password string `json:"password" form:"password" validate:"required"`
	// 签名
	Sign string `json:"sign" form:"timestamp" validate:"required"`
	// 标记
	Marks string `json:"marks" form:"marks" validate:"required"`
}

// UserOverview 用户概览
type UserOverview struct {
	// 关联项目数
	ProjectNum int `json:"project_num"`
	// 待办事项数
	TodoNum int `json:"todo_num"`
	// 申请数
	ApplicationNum int `json:"application_num"`
}

// Login 用户登录
func Login(c *bm.Context) {
	req := &UserLoginReq{}
	if err := c.BindWith(&req, binding.JSON); err != nil {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, common.ErrLoginAccountPwdEmptyMessage))
		return
	}
	// 校验marks sha256
	if flag := utils.SHA256(fmt.Sprintf("%v%v%v", req.UserName, req.Password, req.Sign)); flag != req.Marks {
		err := errors.New(fmt.Sprintf("sha256 not same,%v,%v,%v,%v", req.UserName, req.Password, req.Sign, req.Marks))
		log.Error(err.Error())
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, common.ErrLoginInformationIncorrectMessage))
		return
	}
	//查询数据库验证账户密码
	conditions := map[string]interface{}{
		"name":      req.UserName,
		"is_delete": 0,
	}
	field := "id, name, role, password, salt, create_time, update_time"
	user, err := Svc.User.QueryUser(&biz.UserWhere{Conditions: conditions}, &biz.UserOutput{}, field)
	if err != nil {
		err = errors.New(fmt.Sprintf("query user error or user is empty. userName:%v", req.UserName))
		log.Error(err.Error())
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, common.ErrLoginUserEmptyMessage))
		return
	}
	if len(user) == 0 {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, common.ErrLoginUserEmptyMessage))
		return
	}
	if len(user) > 1 {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, common.ErrLoginUserNotUniqueMessage))
		return
	}
	loginPassword := fmt.Sprintf("%v%v", req.Password, user[0].Salt)
	if user[0].Password != utils.SHA256(loginPassword) {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, common.ErrLoginUserNamePwdMessage))
		return
	}
	//生成JWT
	jwt, err := utils.GenerateJWT(user[0].ID, user[0].Name)
	//生成Cookie
	cookie := http.Cookie{
		Name:     "token",
		Value:    jwt,
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &cookie)
	log.Info("user:%v login success", req.UserName)
	// 获取设置的cookie
	c.JSON(nil, ecode.NewECode(common.SuccessCode, common.SuccessLoginSuccessMessage))
	return
}

// Logout 用户退出登录
func Logout(c *bm.Context) {
	cookie, err := c.Request.Cookie("token")
	if err != nil || cookie == nil || cookie.Value == "" {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, common.ErrDataEmptyMessage))
		c.Abort()
		return
	}
	//解析token
	claims, err := utils.ParseToken(cookie.Value)
	userName := claims.Issuer
	if err != nil {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, common.ErrDataEmptyMessage))
		c.Abort()
		return
	}
	outLoginCookie := http.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &outLoginCookie)
	log.Info("user:%v login successfully logged out", userName)
	// 获取设置的cookie
	c.JSON(nil, ecode.NewECode(common.SuccessCode, common.SuccessOutLoginSuccessMessage))
	return
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *bm.Context) {
	cookie, err := c.Request.Cookie("token")
	if err != nil || cookie == nil || cookie.Value == "" {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, common.ErrDataEmptyMessage))
		c.Abort()
		return
	}
	//解析token
	claims, err := utils.ParseToken(cookie.Value)
	if err != nil {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, common.ErrDataEmptyMessage))
		c.Abort()
		return
	}
	// 获取用户信息
	userConditions := map[string]interface{}{
		"id": claims.UserId,
	}
	users, err := Svc.User.QueryUser(&biz.UserWhere{Conditions: userConditions}, &biz.UserOutput{}, "id, name, role")
	if err != nil || len(users) == 0 {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, common.ErrDataEmptyMessage))
		c.Abort()
		return
	}
	roleMap := map[int]string{
		1: "超级管理员",
		2: "管理员",
		3: "用户",
	}
	role := "用户"
	if roleVal, ok := roleMap[users[0].Role]; ok {
		role = roleVal
	}
	resp := &biz.UserValues{
		ID:       users[0].ID,
		Name:     users[0].Name,
		Role:     users[0].Role,
		RoleName: role,
	}
	c.JSON(resp, err)
	return
}

// GetLoginSign 获取登录签名
func GetLoginSign(c *bm.Context) {
	randStr := randStringBytes(128)
	c.JSON(randStr, nil)
	return
}

// randStringBytes 生成随机字符串
func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = configs.Conf.LoginConf.LetterBytes[rand.Intn(len(configs.Conf.LoginConf.LetterBytes))]
	}
	return string(b)
}

// QueryUserOverview 查询用户总览
func QueryUserOverview(c *bm.Context, w http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value("userInfo").(biz.UserValues)
	projectUserConditions := map[string]interface{}{
		"user_id": userInfo.ID,
	}
	projectUserConfigNum, err := Svc.ProjectUserConfig.CountProjectUserConfig(&biz.ProjectUserConfigWhere{Conditions: projectUserConditions}, "project_config_id")
	if err != nil {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, common.ErrLoginInformationIncorrectMessage))
		c.Abort()
		return
	}
	jobConditions := map[string]interface{}{
		"user_id": userInfo.ID,
	}
	jobNum, err := Svc.Job.CountJob(&biz.JobWhere{Conditions: jobConditions})
	if err != nil {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, common.ErrLoginInformationIncorrectMessage))
		c.Abort()
		return
	}
	resp := &UserOverview{
		ProjectNum:     int(projectUserConfigNum),
		TodoNum:        0,
		ApplicationNum: int(jobNum),
	}
	c.JSON(resp, err)
	return
}

// QueryUserPermission 根据用户权限获取关联的项目
func QueryUserPermission(userInfo biz.UserValues, schema string) (map[string]interface{}, error) {
	conditions := map[string]interface{}{}
	if userInfo.Role == 3 { //普通用户增加项目条件筛选
		if len(userInfo.ProjectConfigID) > 0 {
			switch schema {
			case "cloud_server": //云服务器可直接搜索项目字段匹配
				conditions["project_config_id_IN"] = userInfo.ProjectConfigID
				break
			case "cloud_vpc", "cloud_subnet", "cloud_security_group":
				conditions["project_config"] = map[string]interface{}{
					"id_IN": userInfo.ProjectConfigID,
				}
				break
			case "project_config":
				conditions["id_IN"] = userInfo.ProjectConfigID
				break
			}
		} else { //无项目权限
			return conditions, ecode.NewECode(common.ErrNotJoinProjectMessage)
		}
	}
	return conditions, nil
}
