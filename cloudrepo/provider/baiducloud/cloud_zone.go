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
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/bilibili/HCP/cloudrepo"
)

// CloudZone 云可用区
type CloudZone struct {
	cloudrepo.CloudProductCommon
}

// GetCID 获取云可用区CID
func (c *CloudZone) GetCID() string {
	return c.CID
}

// GetName 获取云可用区名称
func (c *CloudZone) GetName() string {
	return c.Name
}

// GetStatus 获取云可用区状态
func (c *CloudZone) GetStatus() string {
	return c.Status
}

// ListZone 获取可用区列表
func (c *Client) ListZone(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudZone, error) {
	cli, err := c.clientBcc()
	if err != nil {
		return nil, err
	}
	listZoneResult, err := cli.ListZone()
	if err != nil {
		return nil, err
	}
	return c.List2DoZone(listZoneResult)
}

func (c *Client) List2DoZone(resp *api.ListZoneResult) (list []cloudrepo.CloudZone, err error) {
	for _, v := range resp.Zones {
		_ = v
		list = append(list, &CloudZone{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.ZoneName,
				Name:   v.ZoneName,
				Status: "available",
			},
		})
	}
	return
}
