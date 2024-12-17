// Package alicloud
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
package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
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
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	request := ecs.CreateDescribeZonesRequest()
	request.SetConnectTimeout(time.Second * 2)
	request.SetReadTimeout(time.Second * 3)
	request.RegionId = req.Region
	resp, err := cli.DescribeZones(request)
	if err != nil {
		return nil, err
	}
	return List2DoZone(resp)
}

func List2DoZone(resp *ecs.DescribeZonesResponse) (list []cloudrepo.CloudZone, err error) {
	for _, v := range resp.Zones.Zone {
		_ = v
		list = append(list, &CloudZone{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.ZoneId,
				Name:   v.LocalName,
				Status: "available",
			},
		})
	}
	return
}
