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
	"github.com/bilibili/HCP/common/ecode"
	common "github.com/bilibili/HCP/common/models"
	"github.com/go-kratos/kratos/pkg/log"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster/binding"
)

// SyncReq 云资源同步请求
type SyncReq struct {
	// 云ID（例如：alicloud）
	CloudID string `json:"cloud_id" validate:"required"`
	// 云地域
	Region string `json:"region" validate:"required"`
	// 云账号CID
	AccountCID string `json:"account_cid" validate:"required"`
	// 云账号名称
	AccountName string `json:"account_name" validate:"required"`
	// 云账号别名
	AccountAlias string `json:"account_alias" validate:"required"`
	// 云资源类型
	ResourceType string `json:"resource_type" validate:"required"`
	// 云资源名称
	ResourceName string `json:"resource_name"`
	// 云资源ID
	ResourceID string `json:"resource_id"`
}

// SyncResp 云资源同步响应
type SyncResp struct {
	// 响应成功
	Success string `json:"success"`
}

// SyncCloudResource 云资源同步
func SyncCloudResource(c *bm.Context) {
	ctx := c.Request.Context()
	req := &SyncReq{}
	if err := c.BindWith(req, binding.JSON); err != nil {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
		return
	}
	request := &common.SyncCloudResourceRequest{
		Region:       req.Region,
		CloudID:      req.CloudID,
		AccountCID:   req.AccountCID,
		AccountName:  req.AccountName,
		AccountAlias: req.AccountAlias,
		ResourceType: req.ResourceType,
		ResourceName: req.ResourceName,
		ResourceID:   req.ResourceID,
	}
	err := Svc.sync.SyncCloudResource(ctx, request)
	if err != nil {
		log.Error("sync cloud resource(%v): %v", req.ResourceType, err)
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
		return
	}
	data := &SyncResp{
		Success: common.SuccessMessage,
	}
	c.JSON(data, err)
	return
}
