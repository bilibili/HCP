// Package middleware
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
package middleware

import (
	"context"
	"net/http"

	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/bilibili/HCP/app/interface/v1/internal/service"
	"github.com/bilibili/HCP/common/ecode"
	common "github.com/bilibili/HCP/common/models"
	"github.com/bilibili/HCP/utils"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster"
)

// CheckLoginMiddleware 登录校验中间件
func CheckLoginMiddleware(f func(c *blademaster.Context, w http.ResponseWriter, r *http.Request)) blademaster.HandlerFunc {
	return func(c *blademaster.Context) {
		// 获取请求的 Cookie
		cookie, err := c.Request.Cookie("token")
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, common.ErrLoginInformationIncorrectMessage))
			c.Abort()
			return
		}
		if cookie.Value == "" {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, common.ErrNeedLoginUserMessage))
			c.Abort()
			return
		}
		//解析token
		claims, err := utils.ParseToken(cookie.Value)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, common.ErrNeedLoginUserMessage))
			c.Abort()
			return
		}
		// 获取用户信息
		userConditions := map[string]interface{}{
			"id": claims.UserId,
		}
		users, err := service.Svc.User.QueryUser(&biz.UserWhere{Conditions: userConditions}, &biz.UserOutput{}, "id, name, role")
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, common.ErrLoginInformationIncorrectMessage))
			c.Abort()
			return
		}
		if len(users) == 0 {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, common.ErrLoginUserEmptyMessage))
			c.Abort()
			return
		}
		var values biz.UserValues
		// 获取用户权限 1:超级管理员 2:管理员 0:普通用户
		role := users[0].Role
		values.ID = users[0].ID
		values.Name = users[0].Name
		values.Role = role
		if role == 3 { //获取用户关联的项目
			projectUserConditions := map[string]interface{}{
				"user_id": claims.UserId,
			}
			projectUserConfigs, err := service.Svc.ProjectUserConfig.QueryProjectUserConfig(&biz.ProjectUserConfigWhere{Conditions: projectUserConditions}, &biz.ProjectUserConfigOutput{})
			if err != nil {
				c.JSON(nil, ecode.NewECode(common.ErrServiceCode, common.ErrLoginInformationIncorrectMessage))
				c.Abort()
				return
			}
			var ids []int
			for _, v := range projectUserConfigs {
				ids = append(ids, v.ProjectConfigID)
			}
			values.ProjectConfigID = ids
		}
		// 添加用户信息到 context 中
		ctx := context.WithValue(c.Request.Context(), "userInfo", values)
		c.Request = c.Request.WithContext(ctx)
		f(c, c.Writer, c.Request)
	}
}
