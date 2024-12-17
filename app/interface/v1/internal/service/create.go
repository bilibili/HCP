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

// CreateReq 创建数据请求
type CreateReq struct {
	// 表名
	Schema string `json:"schema"`
	// 输入数据
	Inputs string `json:"inputs"`
}

// CreateResp 创建数据响应
type CreateResp struct {
	// 响应数据
	Data string `json:"data"`
}

// Create 创建数据
func Create(c *bm.Context) {
	var s *Service
	req := &CreateReq{}
	if err := c.BindWith(req, binding.JSON); err != nil {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
		return
	}
	switch req.Schema {
	case "cloud_region":
		var cloudRegion []*biz.CloudRegion
		// 解析JSON数据到结构体切片
		err := json.Unmarshal([]byte(req.Inputs), &cloudRegion)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		resp, err := s.CloudRegion.CreateCloudRegion(cloudRegion)
		// 将数据转换为JSON
		jsonData, err := json.Marshal(resp)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		data := &CreateResp{
			Data: string(jsonData),
		}
		c.JSON(data, err)
		return
	case "provider":
		var provider []*biz.Provider
		// 解析JSON数据到结构体切片
		err := json.Unmarshal([]byte(req.Inputs), &provider)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		resp, err := s.Provider.CreateProvider(provider)
		// 将数据转换为JSON
		jsonData, err := json.Marshal(resp)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		data := &CreateResp{
			Data: string(jsonData),
		}
		c.JSON(data, err)
	default:

	}
	c.JSON(nil, ecode.NewECode(common.ErrServiceCode, fmt.Sprintf("create operation unknown schema %s", req.Schema)))
	return
}
