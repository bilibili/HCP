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
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
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

// GetProject 获取云子网所属项目
func (c *CloudSubnet) GetProject() string {
	return c.Project
}

// GetVpc 获取云子网所属专有网络
func (c *CloudSubnet) GetVpc() string {
	return c.Vpc
}

// GetZone 获取云子网所属可用区
func (c *CloudSubnet) GetZone() string {
	return c.Zone
}

// ListSubnet 获取子网列表
func (c *Client) ListSubnet(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudSubnet, error) {
	cli, err := c.clientVpc()
	if err != nil {
		return nil, err
	}
	var data []vpc.VSwitch
	pageNum, pageSize := 1, 50
	request := vpc.CreateDescribeVSwitchesRequest()
	request.RegionId = req.Region
	for {
		request.PageNumber = requests.Integer(fmt.Sprintf("%d", pageNum))
		request.PageSize = requests.Integer(fmt.Sprintf("%d", pageSize))
		resp, err := cli.DescribeVSwitches(request)
		if err != nil {
			return nil, err
		}
		data = append(data, resp.VSwitches.VSwitch...)
		if len(data) >= resp.TotalCount {
			break
		}
		pageNum += 1
	}

	return list2DoSubnet(data)
}

func list2DoSubnet(resp []vpc.VSwitch) (list []cloudrepo.CloudSubnet, err error) {
	for _, v := range resp {
		list = append(list, &CloudSubnet{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.VSwitchId,
				Name:   v.VSwitchName,
				Status: cloudrepo.GetCloudEnum(cloudrepo.CloudAli, "cloud_subnet", v.Status),
			},
			Cidr:    v.CidrBlock,
			Project: v.ResourceGroupId,
			Vpc:     v.VpcId,
			Zone:    v.ZoneId,
		})
	}
	return
}
