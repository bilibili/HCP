// Package huaweicloud
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
package huaweicloud

import (
	"fmt"

	"github.com/bilibili/HCP/cloudrepo"
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	modelV2 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	"github.com/pkg/errors"
)

var ecsEndPointMap = map[string]string{
	"cn-south-2": "https://ecs.cn-south-2.myhuaweicloud.com",
}

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
	// 部分Region实例化client会报错，这里进行异常捕获
	defer func() {
		e := recover()
		if e != nil {
			err := errors.New(fmt.Sprintf("%v", e))
			if err != nil {
				return
			}
		}
	}()
	var cli *ecs.EcsClient
	var data []modelV2.NovaAvailabilityZone
	ep, ok := ecsEndPointMap[req.Region]
	if ok {
		projects, err := c.ListProjectIam(req)
		if err != nil {
			return nil, err
		}
		for _, p := range projects {
			if p.Name != req.Region {
				continue
			}
			cli = c.clientEcsWithEndPoint(p.Id, ep)
			servers, err := c.queryZone(req, cli)
			if err != nil {
				return nil, err
			}
			for idx, _ := range servers {
				data = append(data, servers[idx])
			}
		}
		return list2DoZone(data)
	} else {
		cli = c.clientEcs()
		resp, err := c.queryZone(req, cli)
		if err != nil {
			return nil, err
		}
		return list2DoZone(resp)
	}
}

func list2DoZone(resp []modelV2.NovaAvailabilityZone) (list []cloudrepo.CloudZone, err error) {
	for _, v := range resp {
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

func (c *Client) queryZone(req *cloudrepo.GetCloudProductReq, cli *ecs.EcsClient) (data []modelV2.NovaAvailabilityZone, err error) {
	req.DisablePage = true
	request := modelV2.NovaListAvailabilityZonesRequest{}
	resp, err := cli.NovaListAvailabilityZones(&request)
	if err != nil {
		return nil, err
	}
	return *resp.AvailabilityZoneInfo, nil
}
