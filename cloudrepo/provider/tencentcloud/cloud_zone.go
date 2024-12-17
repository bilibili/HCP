// Package tencentcloud
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
package tencentcloud

import (
	"encoding/json"

	"github.com/bilibili/HCP/cloudrepo"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
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
	request := cvm.NewDescribeZonesRequest()
	resp, err := cli.DescribeZones(request)
	if err != nil {
		return nil, err
	}
	response := &struct {
		Response *struct {
			TotalCount uint64 `json:"TotalCount,omitempty" name:"TotalCount"`
			ZoneSet    []Zone `json:"ZoneSet,omitempty" name:"ZoneSet"`
			RequestId  string `json:"RequestId,omitempty" name:"RequestId"`
		} `json:"Response"`
	}{}
	err = json.Unmarshal([]byte(resp.ToJsonString()), response)
	if err != nil {
		return nil, err
	}
	return list2DoZone(response.Response.ZoneSet)
}

func list2DoZone(resp []Zone) (list []cloudrepo.CloudZone, err error) {
	for _, v := range resp {
		_ = v
		list = append(list, &CloudZone{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.Zone,
				Name:   v.ZoneName,
				Status: cloudrepo.GetCloudEnum(cloudrepo.CloudTencent, "cloud_zone", v.ZoneState),
			},
		})
	}
	return
}
