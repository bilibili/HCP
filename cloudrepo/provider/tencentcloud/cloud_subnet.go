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

// CloudSubnet 云子网
type CloudSubnet struct {
	cloudrepo.CloudProductCommon

	Cidr    string `json:"cidr,omitempty" name:"Cidr"`
	Project string `json:"project,omitempty" name:"Project"`
	Vpc     string `json:"vpc,omitempty" name:"Vpc"`
	Zone    string `json:"zone,omitempty" name:"Zone"`
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

// GetCidr 获取云子网CIDR
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
	cli, err := c.clientVpc()
	if err != nil {
		return nil, err
	}
	var data []Subnet
	pageNum, pageSize := 1, 100
	request := vpc.NewDescribeSubnetsRequest()
	if req.ResourceID != "" {
		request.SubnetIds = []*string{&req.ResourceID}
	}
	for {
		offset := fmt.Sprintf("%v", uint64((pageNum-1)*pageSize))
		limit := fmt.Sprintf("%v", uint64(pageSize))
		request.Offset = &offset
		request.Limit = &limit
		resp, err := cli.DescribeSubnets(request)
		if err != nil {
			return nil, err
		}
		response := &struct {
			Response struct {
				// 子网列表信息
				SubnetSet []Subnet `json:"SubnetSet,omitempty" name:"SubnetSet"`
				// 返回的子网总数
				TotalCount uint64 `json:"TotalCount,omitempty" name:"TotalCount"`
				// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
				RequestId string `json:"RequestId,omitempty" name:"RequestId"`
			} `json:"Response"`
		}{}
		err = json.Unmarshal([]byte(resp.ToJsonString()), response)
		if err != nil {
			return nil, err
		}
		data = append(data, response.Response.SubnetSet...)
		if len(data) >= int(response.Response.TotalCount) {
			break
		}
		pageNum += 1
	}

	return c.list2DoSubnet(data)
}

func (c *Client) list2DoSubnet(resp []Subnet) (list []cloudrepo.CloudSubnet, err error) {
	for _, v := range resp {
		_ = v
		list = append(list, &CloudSubnet{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.SubnetId,
				Name:   v.SubnetName,
				Status: "available",
			},
			Cidr:    v.CidrBlock,
			Project: c.GetProject(v.TagSet),
			Vpc:     v.VpcId,
			Zone:    v.Zone,
		})
	}
	return
}
