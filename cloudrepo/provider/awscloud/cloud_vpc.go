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

// CloudVpc 云专有网络
type CloudVpc struct {
	cloudrepo.CloudProductCommon
	Cidr    string
	Project string
}

// GetCID 获取云专有网络CID
func (c *CloudVpc) GetCID() string {
	return c.CID
}

// GetName 获取云专有网络名称
func (c *CloudVpc) GetName() string {
	return c.Name
}

// GetStatus 获取云专有网络状态
func (c *CloudVpc) GetStatus() string {
	return c.Status
}

// GetCidr 获取云专有网络CIDR
func (c *CloudVpc) GetCidr() string {
	return c.Cidr
}

// GetProject 获取云专有网络项目
func (c *CloudVpc) GetProject() string {
	return c.Project
}

// ListVpc 获取专有网路列表
func (c *Client) ListVpc(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudVpc, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFn()
	cli := c.clientEc2()
	var data []types.Vpc
	pageSize := 100
	maxResult := int32(pageSize)
	request := &ec2.DescribeVpcsInput{MaxResults: &maxResult}
	for {
		if req.Cursor != "" {
			request.NextToken = &req.Cursor
		}
		resp, err := cli.DescribeVpcs(ctx, request)
		if err != nil {
			return nil, err
		}
		if resp.Vpcs != nil && len(resp.Vpcs) > 0 {
			data = append(data, resp.Vpcs...)
		}
		if resp.NextToken != nil {
			req.Cursor = *resp.NextToken
		}
		if resp.NextToken == nil {
			break
		}
	}
	return c.list2DoVpc(data)
}

func (c *Client) list2DoVpc(resp []types.Vpc) (list []cloudrepo.CloudVpc, err error) {
	for _, v := range resp {
		_ = v
		list = append(list, &CloudVpc{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    *v.VpcId,
				Name:   c.GetNameTag(v.Tags),
				Status: cloudrepo.GetCloudEnum(cloudrepo.CloudAws, "cloud_vpc", string(v.State)),
			},
			Cidr:    *v.CidrBlock,
			Project: c.GetEc2ProjectTag(v.Tags),
		})
	}
	return
}
