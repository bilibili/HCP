// Package baiducloud
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
package baiducloud

import (
	"github.com/bilibili/HCP/cloudrepo"
)

// CloudRegion 云地域
type CloudRegion struct {
	cloudrepo.CloudProductCommon
	client *Client
}

// GetCID 获取云地域CID
func (c *CloudRegion) GetCID() string {
	return c.CID
}

// GetName 获取云地域名称
func (c *CloudRegion) GetName() string {
	return c.Name
}

// GetStatus 获取云地域状态
func (c *CloudRegion) GetStatus() string {
	return c.Status
}

// ListRegion 获取地域列表
func (c *Client) ListRegion(_ *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudRegion, error) {
	var data []cloudrepo.CloudRegion
	baiduRegionList := []string{"cn-bj", "cn-gz", "cn-su", "cn-hkg", "cn-fwh", "cn-bd", "cn-sin", "cn-fsh"}
	baiduRegionName := map[string]string{
		"cn-bj":  "北京",
		"cn-gz":  "广州",
		"cn-su":  "苏州",
		"cn-hkg": "香港",
		"cn-fwh": "武汉",
		"cn-bd":  "保定",
		"cn-sin": "新加坡",
		"cn-fsh": "上海",
	}
	for _, region := range baiduRegionList {
		data = append(data, &CloudRegion{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    region,
				Name:   baiduRegionName[region],
				Status: "available",
			},
		})
	}
	return data, nil
}
