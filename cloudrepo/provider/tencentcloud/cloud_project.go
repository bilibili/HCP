// Package tencentcloud
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
package tencentcloud

import (
	"encoding/json"
	"fmt"

	"github.com/bilibili/HCP/cloudrepo"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
)

// CloudProject 云项目
type CloudProject struct {
	cloudrepo.CloudProductCommon
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
func (c *Client) ListProject(_ *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudProject, error) {
	cli, err := c.clientProject()
	if err != nil {
		return nil, err
	}
	request := dcdb.NewDescribeProjectsRequest()
	resp, err := cli.DescribeProjects(request)
	if err != nil {
		return nil, err
	}
	response := &struct {
		Response struct {
			Projects  []Project `json:"Projects,omitempty" name:"Projects"`
			RequestId string    `json:"RequestId,omitempty" name:"RequestId"`
		} `json:"Response"`
	}{}
	err = json.Unmarshal([]byte(resp.ToJsonString()), response)
	if err != nil {
		return nil, err
	}
	return list2Do(response.Response.Projects)
}

func list2Do(resp []Project) (list []cloudrepo.CloudProject, err error) {
	for _, v := range resp {
		_ = v
		list = append(list, &CloudProject{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:  fmt.Sprintf("%v", v.ProjectId),
				Name: v.Name,
			},
		})
	}
	return
}
