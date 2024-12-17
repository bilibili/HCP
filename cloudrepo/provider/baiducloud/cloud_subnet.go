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
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/bilibili/HCP/app/interface/v1/configs"
	"github.com/bilibili/HCP/cloudrepo"
)

// CloudSubnet 云子网
type CloudSubnet struct {
	cloudrepo.CloudProductCommon

	Cidr    string
	Project string
	Vpc     string
	Zone    string
}

// GetCID 获取云子网CID
func (c *CloudSubnet) GetCID() string {
	return c.CID
}

// GetName 获取云子网名称
func (c *CloudSubnet) GetName() string {
	return c.Name
}

// GetStatus 获取云子网状态
func (c *CloudSubnet) GetStatus() string {
	return c.Status
}

// GetCidr 获取云子网CIDR
func (c *CloudSubnet) GetCidr() string {
	return c.Cidr
}

// GetProject 获取云子网项目
func (c *CloudSubnet) GetProject() string {
	return c.Project
}

// GetVpc 获取云子网专有网络
func (c *CloudSubnet) GetVpc() string {
	return c.Vpc
}

// GetZone 获取云子网可用区
func (c *CloudSubnet) GetZone() string {
	return c.Zone
}

// ListSubnet 获取子网列表
func (c *Client) ListSubnet(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudSubnet, error) {
	cli, err := c.clientVpc()
	if err != nil {
		return nil, err
	}
	maxItems := 1000
	var marker = new(string)
	*marker = ""
	request := &vpc.ListSubnetArgs{
		MaxKeys: maxItems,
	}
	var data []vpc.Subnet
	for {
		request.Marker = *marker
		resp, err := cli.ListSubnets(request)
		if err != nil {
			return nil, err
		}
		data = append(data, resp.Subnets...)
		marker = &resp.NextMarker
		if !resp.IsTruncated {
			break
		}
	}
	return list2DoSubnet(data)
}

func list2DoSubnet(resp []vpc.Subnet) (list []cloudrepo.CloudSubnet, err error) {
	for _, v := range resp {
		_ = v
		var projectCid string
		for _, tag := range v.Tags {
			if tag.TagKey == configs.Conf.CloudConf.TagProjectKey {
				projectCid = fmt.Sprintf("%s-%s", tag.TagKey, tag.TagValue)
			}
		}
		list = append(list, &CloudSubnet{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.SubnetId,
				Name:   v.Name,
				Status: "available",
			},
			Cidr:    v.Cidr,
			Project: projectCid,
			Vpc:     v.VPCId,
			Zone:    v.ZoneName,
		})
	}
	return
}
