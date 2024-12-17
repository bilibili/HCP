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

// CloudVpc 云专有网络
type CloudVpc struct {
	cloudrepo.CloudProductCommon
	Cidr    string
	Project string
}

// GetCID 获取专有网络CID
func (c *CloudVpc) GetCID() string {
	return c.CID
}

// GetName 获取专有网络名称
func (c *CloudVpc) GetName() string {
	return c.Name
}

// GetStatus 获取专有网络状态
func (c *CloudVpc) GetStatus() string {
	return c.Status
}

// GetCidr 获取专有网络CIDR
func (c *CloudVpc) GetCidr() string {
	return c.Cidr
}

// GetProject 获取专有网络所属项目
func (c *CloudVpc) GetProject() string {
	return c.Project
}

// ListVpc 获取专有网络列表
func (c *Client) ListVpc(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudVpc, error) {
	cli, err := c.clientVpc()
	if err != nil {
		return nil, err
	}
	var data []vpc.Vpc
	pageNum, pageSize := 1, 50
	request := vpc.CreateDescribeVpcsRequest()
	request.RegionId = req.Region
	request.VpcId = req.ResourceID
	request.VpcName = req.ResourceName
	for {
		request.PageNumber = requests.Integer(fmt.Sprintf("%d", pageNum))
		request.PageSize = requests.Integer(fmt.Sprintf("%d", pageSize))
		resp, err := cli.DescribeVpcs(request)
		if err != nil {
			return nil, err
		}
		data = append(data, resp.Vpcs.Vpc...)
		if len(data) >= resp.TotalCount {
			break
		}
		pageNum += 1
	}
	return list2DoVpc(data)
}

func list2DoVpc(resp []vpc.Vpc) (list []cloudrepo.CloudVpc, err error) {
	statusMap := map[string]string{
		"Available": "available",
		"Pending":   "pending",
	}
	for _, v := range resp {
		_ = v
		list = append(list, &CloudVpc{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.VpcId,
				Name:   v.VpcName,
				Status: statusMap[v.Status],
			},
			Cidr:    v.CidrBlock,
			Project: v.ResourceGroupId,
		})
	}
	return
}
