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
	"fmt"

	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/bilibili/HCP/common/ecode"
	common "github.com/bilibili/HCP/common/models"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster/binding"
)

// DeleteReq 删除数据请求
type DeleteReq struct {
	// 表名
	Schema string `json:"schema"`
	// 删除的id
	Ids string `json:"ids"`
}

// DeleteResp 删除数据返回
type DeleteResp struct {
	// 响应的数据
	Data string `json:"data"`
}

// Delete 删除数据
func Delete(c *bm.Context) {
	req := &DeleteReq{}
	if err := c.BindWith(req, binding.JSON); err != nil {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
		return
	}
	ids := biz.GetStringToIntSlice(req.Ids)
	switch req.Schema {
	case "cloud_region":
		err := Svc.CloudRegion.DeleteCloudRegion(ids)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
	case "account":
		err := Svc.Account.DeleteAccount(0, ids)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
	case "user":
		err := Svc.User.DeleteUser(0, ids)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
	case "project_config":
		err := Svc.ProjectConfig.DeleteProjectConfig(0, ids)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
	default:
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, fmt.Sprintf("delete operation unknown schema %s", req.Schema)))
		return
	}
	c.JSON(nil, ecode.NewECode(common.SuccessCode, common.SuccessMessage))
	return
}
