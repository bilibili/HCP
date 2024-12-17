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
	"encoding/json"
	"fmt"

	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/bilibili/HCP/common/ecode"
	common "github.com/bilibili/HCP/common/models"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster/binding"
)

// UpdateReq 更新数据请求参数
type UpdateReq struct {
	// 表名
	Schema string `json:"schema"`
	// 条件
	Where json.RawMessage `json:"where"`
	// 更新数据
	Inputs string `json:"inputs"`
}

// UpdateResp 更新数据返回参数
type UpdateResp struct {
}

// Update 更新数据
func Update(c *bm.Context) {
	var s *Service
	req := &UpdateReq{}
	if err := c.BindWith(req, binding.JSON); err != nil {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
		return
	}
	switch req.Schema {
	case "cloud_region":
		condition, err := biz.ParseCloudData(req.Where)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		var cloudRegion *biz.CloudRegion
		// 解析JSON数据到结构体切片
		err = json.Unmarshal([]byte(req.Inputs), &cloudRegion)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		err = s.CloudRegion.UpdateCloudRegion(&biz.CloudRegionWhere{Conditions: condition}, cloudRegion)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		data := &CreateResp{}
		c.JSON(data, err)
	default:

	}
	c.JSON(nil, ecode.NewECode(common.ErrServiceCode, fmt.Sprintf("update operation unknown schema %s", req.Schema)))
	return
}
