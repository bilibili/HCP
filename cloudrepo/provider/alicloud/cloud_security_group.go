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
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/bilibili/HCP/cloudrepo"
)

// CloudSecurityGroup 云安全组
type CloudSecurityGroup struct {
	cloudrepo.CloudProductCommon

	Project string
	Vpc     string
}

// GetCID 云安全组CID
func (c *CloudSecurityGroup) GetCID() string {
	return c.CID
}

// GetName 云安全组名称
func (c *CloudSecurityGroup) GetName() string {
	return c.Name
}

// GetStatus 云安全组状态
func (c *CloudSecurityGroup) GetStatus() string {
	return c.Status
}

// GetProject 云安全组项目
func (c *CloudSecurityGroup) GetProject() string {
	return c.Project
}

// GetVpc 云安全组专有网路
func (c *CloudSecurityGroup) GetVpc() string {
	return c.Vpc
}

// ListCloudSecurityGroup 获取安全组列表
func (c *Client) ListCloudSecurityGroup(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudSecurityGroup, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	var data []ecs.SecurityGroup
	pageNum, pageSize := 1, 50
	request := ecs.CreateDescribeSecurityGroupsRequest()
	request.RegionId = req.Region
	request.SecurityGroupId = req.ResourceID
	request.SecurityGroupName = req.ResourceName
	for {
		request.PageNumber = requests.Integer(fmt.Sprintf("%d", pageNum))
		request.PageSize = requests.Integer(fmt.Sprintf("%d", pageSize))
		resp, err := cli.DescribeSecurityGroups(request)
		if err != nil {
			return nil, err
		}
		data = append(data, resp.SecurityGroups.SecurityGroup...)
		if len(data) >= resp.TotalCount {
			break
		}
		pageNum += 1
	}
	return list2DoSecurityGroup(data)
}

func list2DoSecurityGroup(resp []ecs.SecurityGroup) (list []cloudrepo.CloudSecurityGroup, err error) {
	for _, v := range resp {
		_ = v
		list = append(list, &CloudSecurityGroup{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.SecurityGroupId,
				Name:   v.SecurityGroupName,
				Status: "available",
			},
			Project: v.ResourceGroupId,
			Vpc:     v.VpcId,
		})
	}
	return
}
