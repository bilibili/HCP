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
	sg "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

// CloudSecurityGroup 云安全组
type CloudSecurityGroup struct {
	cloudrepo.CloudProductCommon

	Project string
	Vpc     string
}

// Tag 标签
type Tag struct {
	Key   string `json:"Key,omitempty" name:"Key"`
	Value string `json:"Value,omitempty" name:"Value"`
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
	cli, err := c.clientSecurityGroup()
	if err != nil {
		return nil, err
	}
	var data []SecurityGroup
	pageNum, pageSize := 1, 100
	request := sg.NewDescribeSecurityGroupsRequest()
	if req.ResourceID != "" {
		request.SecurityGroupIds = []*string{&req.ResourceID}
	}
	for {
		offset := fmt.Sprintf("%v", uint64((pageNum-1)*pageSize))
		limit := fmt.Sprintf("%v", uint64(pageSize))
		request.Offset = &offset
		request.Limit = &limit
		resp, err := cli.DescribeSecurityGroups(request)
		if err != nil {
			return nil, err
		}
		response := &struct {
			Response struct {
				// 安全组对象。
				// 注意：此字段可能返回 null，表示取不到有效值。
				SecurityGroupSet []SecurityGroup `json:"SecurityGroupSet,omitempty" name:"SecurityGroupSet"`
				// 符合条件的实例数量。
				TotalCount uint64 `json:"TotalCount,omitempty" name:"TotalCount"`
				// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
				RequestId string `json:"RequestId,omitempty" name:"RequestId"`
			} `json:"Response"`
		}{}
		err = json.Unmarshal([]byte(resp.ToJsonString()), response)
		if err != nil {
			return nil, err
		}
		data = append(data, response.Response.SecurityGroupSet...)
		if len(data) >= int(*resp.Response.TotalCount) {
			break
		}
		pageNum += 1
	}
	return list2DoSecurityGroup(data), nil
}

func list2DoSecurityGroup(resp []SecurityGroup) (list []cloudrepo.CloudSecurityGroup) {
	for _, v := range resp {
		_ = v
		list = append(list, &CloudSecurityGroup{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.SecurityGroupId,
				Name:   v.SecurityGroupName,
				Status: "available",
			},
			Project: v.ProjectId,
			Vpc:     "",
		})
	}
	return
}
