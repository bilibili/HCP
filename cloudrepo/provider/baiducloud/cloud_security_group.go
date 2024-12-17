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

	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/bilibili/HCP/app/interface/v1/configs"
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
	cli, err := c.clientBcc()
	if err != nil {
		return nil, err
	}
	maxItems := 1000
	var marker = new(string)
	*marker = ""
	request := &api.ListSecurityGroupArgs{
		// 指定每页包含的最大数量(主实例)，最大数量不超过1000，缺省值为1000，可选
		MaxKeys: maxItems,
	}
	var data []api.SecurityGroupModel
	for {
		request.Marker = *marker
		resp, err := cli.ListSecurityGroup(request)
		if err != nil {
			return nil, err
		}
		data = append(data, resp.SecurityGroups...)
		marker = &resp.NextMarker
		if !resp.IsTruncated {
			break
		}
	}
	return list2DoSecurityGroup(&api.ListSecurityGroupResult{
		SecurityGroups: data,
	})
}

func list2DoSecurityGroup(resp *api.ListSecurityGroupResult) (list []cloudrepo.CloudSecurityGroup, err error) {
	for _, v := range resp.SecurityGroups {
		_ = v
		var projectCid string
		for _, tag := range v.Tags {
			if tag.TagKey == configs.Conf.CloudConf.TagProjectKey {
				projectCid = fmt.Sprintf("%s-%s", tag.TagKey, tag.TagValue)
			}
		}
		list = append(list, &CloudSecurityGroup{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.Id,
				Name:   v.Name,
				Status: "available",
			},
			Project: projectCid,
			Vpc:     v.VpcId,
		})
	}
	return
}
