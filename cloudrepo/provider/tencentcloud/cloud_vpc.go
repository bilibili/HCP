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
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
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
	var data []Vpc
	pageNum, pageSize := 1, 100
	request := vpc.NewDescribeVpcsRequest()
	if req.ResourceID != "" {
		request.VpcIds = []*string{&req.ResourceID}
	}
	for {
		offset := fmt.Sprintf("%v", uint64((pageNum-1)*pageSize))
		limit := fmt.Sprintf("%v", uint64(pageSize))
		request.Offset = &offset
		request.Limit = &limit
		resp, err := cli.DescribeVpcs(request)
		if err != nil {
			return nil, err
		}
		response := &struct {
			Response *struct {
				// 符合条件的对象数。
				TotalCount uint64 `json:"TotalCount,omitempty" name:"TotalCount"`
				// VPC对象。
				VpcSet []Vpc `json:"VpcSet,omitempty" name:"VpcSet"`
				// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
				RequestId string `json:"RequestId,omitempty" name:"RequestId"`
			} `json:"Response"`
		}{}
		err = json.Unmarshal([]byte(resp.ToJsonString()), response)
		if err != nil {
			return nil, err
		}
		data = append(data, response.Response.VpcSet...)
		if len(data) >= int(response.Response.TotalCount) {
			break
		}
	}
	return c.list2DoVpc(data)
}

func (c *Client) list2DoVpc(resp []Vpc) (list []cloudrepo.CloudVpc, err error) {
	for _, v := range resp {
		_ = v
		list = append(list, &CloudVpc{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.VpcId,
				Name:   v.VpcName,
				Status: "available",
			},
			Cidr:    v.CidrBlock,
			Project: c.GetProject(v.TagSet),
		})
	}
	return
}
