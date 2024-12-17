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
	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/bilibili/HCP/cloudrepo"
)

// CloudProject 云项目
type CloudProject struct {
	cloudrepo.CloudProductCommon
	client *Client

	ProjectID string
}

// GetCID 获取云项目CID
func (c *CloudProject) GetCID() string {
	return c.CID
}

// GetName 获取云项目名称
func (c *CloudProject) GetName() string {
	return c.Name
}

// GetStatus 获取云项目状态
func (c *CloudProject) GetStatus() string {
	return c.Status
}

// ListProject 获取项目列表
func (c *Client) ListProject(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudProject, error) {
	cli, err := c.clientProject()
	if err != nil {
		return nil, err
	}
	var data []resourcemanager.ResourceGroup
	pageNum, pageSize := 1, 50
	request := resourcemanager.CreateListResourceGroupsRequest()
	request.SetScheme("https")
	request.RegionId = req.Region
	request.DisplayName = req.ResourceName
	for {
		request.PageNumber = requests.Integer(fmt.Sprintf("%d", pageNum))
		request.PageSize = requests.Integer(fmt.Sprintf("%d", pageSize))
		resp, err := cli.ListResourceGroups(request)
		if err != nil {
			return nil, err
		}
		data = append(data, resp.ResourceGroups.ResourceGroup...)
		if len(data) >= resp.TotalCount {
			break
		}
		pageNum += 1
	}

	return list2Do(data), nil
}

func list2Do(resp []resourcemanager.ResourceGroup) (list []cloudrepo.CloudProject) {
	for _, v := range resp {
		_ = v
		list = append(list, &CloudProject{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:  v.Id,
				Name: v.DisplayName,
			},
			//ProjectID: v.Id,
		})
	}
	return
}
