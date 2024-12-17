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
	cli, err := c.clientVpc()
	if err != nil {
		return nil, err
	}
	maxItems := 1000
	var marker = new(string)
	*marker = ""
	request := &vpc.ListVPCArgs{
		MaxKeys: maxItems,
	}
	var data []vpc.VPC
	for {
		request.Marker = *marker
		resp, err := cli.ListVPC(request)
		if err != nil {
			return nil, err
		}
		data = append(data, resp.VPCs...)
		marker = &resp.NextMarker
		if !resp.IsTruncated {
			break
		}
	}
	return list2DoVpc(data)
}

func list2DoVpc(resp []vpc.VPC) (list []cloudrepo.CloudVpc, err error) {
	for _, v := range resp {
		_ = v
		var projectCid string
		for _, tag := range v.Tags {
			if tag.TagKey == configs.Conf.CloudConf.TagProjectKey {
				projectCid = fmt.Sprintf("%s-%s", tag.TagKey, tag.TagValue)
			}
		}
		list = append(list, &CloudVpc{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.VPCID,
				Name:   v.Name,
				Status: "available",
			},
			Cidr:    v.Cidr,
			Project: projectCid,
		})
	}
	return
}
