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

// GetCidr 获取云子网CID
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
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFn()
	cli := c.clientEc2()
	var data []types.Subnet
	pageSize := 100
	maxResult := int32(pageSize)
	request := &ec2.DescribeSubnetsInput{MaxResults: &maxResult}
	for {
		if req.Cursor != "" {
			request.NextToken = &req.Cursor
		}
		resp, err := cli.DescribeSubnets(ctx, request)
		if err != nil {
			return nil, err
		}
		if resp.Subnets != nil && len(resp.Subnets) > 0 {
			data = append(data, resp.Subnets...)
		}
		if resp.NextToken != nil {
			req.Cursor = *resp.NextToken
		}
		if resp.NextToken == nil {
			break
		}
	}
	return c.list2DoSubnet(data)
}

func (c *Client) list2DoSubnet(resp []types.Subnet) (list []cloudrepo.CloudSubnet, err error) {
	for _, v := range resp {
		_ = v
		list = append(list, &CloudSubnet{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    *v.SubnetId,
				Name:   c.GetNameTag(v.Tags),
				Status: cloudrepo.GetCloudEnum(cloudrepo.CloudAws, "cloud_subnet", string(v.State)),
			},
			Cidr:    *v.CidrBlock,
			Project: c.GetEc2ProjectTag(v.Tags),
			Vpc:     *v.VpcId,
			Zone:    *v.AvailabilityZoneId,
		})
	}
	return
}
