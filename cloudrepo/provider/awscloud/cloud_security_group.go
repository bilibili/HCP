// Package awscloud
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
package awscloud

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/bilibili/HCP/cloudrepo"
)

// CloudSecurityGroup 云安全组
type CloudSecurityGroup struct {
	cloudrepo.CloudProductCommon

	Project string
	Vpc     string
}

// GetCID 获取云安全组CID
func (c *CloudSecurityGroup) GetCID() string {
	return c.CID
}

// GetName 获取云安全组名称
func (c *CloudSecurityGroup) GetName() string {
	return c.Name
}

// GetStatus 获取云安全组状态
func (c *CloudSecurityGroup) GetStatus() string {
	return c.Status
}

// GetProject 获取云安全组项目
func (c *CloudSecurityGroup) GetProject() string {
	return c.Project
}

// GetVpc 获取云安全组专有网络
func (c *CloudSecurityGroup) GetVpc() string {
	return c.Vpc
}

// ListCloudSecurityGroup 获取安全组列表
func (c *Client) ListCloudSecurityGroup(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudSecurityGroup, error) {
	req.DisablePage = true
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFn()
	var data []types.SecurityGroup
	pageSize := 1000
	maxResult := int32(pageSize)
	request := &ec2.DescribeSecurityGroupsInput{MaxResults: &maxResult}
	for {
		if req.Cursor != "" {
			request.NextToken = &req.Cursor
		}
		cli := c.clientEc2()
		resp, err := cli.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{})
		if err != nil {
			return nil, err
		}
		if resp.SecurityGroups != nil && len(resp.SecurityGroups) > 0 {
			data = append(data, resp.SecurityGroups...)
		}
		if resp.NextToken != nil {
			req.Cursor = *resp.NextToken
		}
		if resp.NextToken == nil {
			break
		}
	}
	return c.list2DoSecurityGroup(data)
}

func (c *Client) list2DoSecurityGroup(resp []types.SecurityGroup) (list []cloudrepo.CloudSecurityGroup, err error) {
	for _, v := range resp {
		_ = v
		list = append(list, &CloudSecurityGroup{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    *v.GroupId,
				Name:   *v.GroupName,
				Status: "available",
			},
			Project: c.GetEc2ProjectTag(v.Tags),
			Vpc:     *v.VpcId,
		})
	}
	return
}
