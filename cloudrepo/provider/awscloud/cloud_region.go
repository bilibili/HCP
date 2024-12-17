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
	"github.com/bilibili/HCP/cloudrepo"
)

// CloudRegion 云地域
type CloudRegion struct {
	cloudrepo.CloudProductCommon
	client *Client
}

// GetCID 获取云地域CID
func (c *CloudRegion) GetCID() string {
	return c.CID
}

// GetName 获取云地域名称
func (c *CloudRegion) GetName() string {
	return c.Name
}

// GetStatus 获取云地域状态
func (c *CloudRegion) GetStatus() string {
	return c.Status
}

// ListRegion 获取地域列表
func (c *Client) ListRegion(_ *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudRegion, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFn()
	c.Region = "us-east-1"
	cli := c.clientEc2()
	request := &ec2.DescribeRegionsInput{}
	reps, err := cli.DescribeRegions(ctx, request)
	if err != nil {
		return nil, err
	}
	return list2DoRegion(reps)
}

func list2DoRegion(resp *ec2.DescribeRegionsOutput) (list []cloudrepo.CloudRegion, err error) {
	for _, v := range resp.Regions {
		_ = v
		list = append(list, &CloudRegion{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    *v.RegionName,
				Name:   cloudrepo.GetCloudEnum(cloudrepo.CloudAws, "cloud_region_location", *v.RegionName),
				Status: "available",
			},
		})
	}
	return
}
