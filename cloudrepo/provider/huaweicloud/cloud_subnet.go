// Package huaweicloud
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
package huaweicloud

import (
	"fmt"
	"strings"

	"github.com/bilibili/HCP/cloudrepo"
	vpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2/model"
	"github.com/pkg/errors"
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
	// 部分Region实例化client会报错，这里进行异常捕获
	defer func() {
		e := recover()
		if e != nil {
			err := errors.New(fmt.Sprintf("%v", e))
			if err != nil {
				return
			}
		}
	}()
	var cli *vpc.VpcClient
	var data []model.Subnet
	ep, ok := vpcEndpointMap[req.Region]
	if ok {
		projects, err := c.ListProjectIam(req)
		if err != nil {
			return nil, err
		}
		for _, p := range projects {
			if p.Name != req.Region {
				continue
			}
			cli = c.clientVpcWithEndpoint(p.Id, ep)
			servers, err := c.querySubnet(req, cli)
			if err != nil {
				return nil, err
			}
			for idx, _ := range servers {
				data = append(data, servers[idx])
			}
		}
		return list2DoSubnet(data)
	} else {
		cli = c.clientVpc()
		resp, err := c.querySubnet(req, cli)
		if err != nil {
			return nil, err
		}
		return list2DoSubnet(resp)
	}
}

func list2DoSubnet(resp []model.Subnet) (list []cloudrepo.CloudSubnet, err error) {
	for _, v := range resp {
		_ = v
		statusBts, _ := v.Status.MarshalJSON()
		status := string(statusBts)
		status = strings.Replace(status, `"`, "", -1)
		status = strings.Replace(status, "\n", "", -1)
		list = append(list, &CloudSubnet{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.Id,
				Name:   v.Name,
				Status: cloudrepo.GetCloudEnum(cloudrepo.CloudHuawei, "cloud_subnet", status),
			},
			Cidr:    v.Cidr,
			Project: "",
			Vpc:     v.VpcId,
			Zone:    v.AvailabilityZone,
		})
	}
	return
}

func (c *Client) querySubnet(req *cloudrepo.GetCloudProductReq, cli *vpc.VpcClient) (data []model.Subnet, err error) {
	var marker = new(string)
	*marker = ""
	pageSize := 100
	request := model.ListSubnetsRequest{}
	if req.ResourceID != "" {
		request.VpcId = &req.ResourceID
	}
	for {
		limit := int32(pageSize)
		request.Limit = &limit
		request.Marker = marker
		resp, err := cli.ListSubnets(&request)
		if err != nil {
			return nil, err
		}
		data = append(data, *resp.Subnets...)
		if resp != nil && resp.Subnets != nil && len(*resp.Subnets) > 0 {
			last := (*resp.Subnets)[len(*resp.Subnets)-1]
			marker = &last.Id
			if &last.Id == nil || &last.Id == marker {
				break
			}
		} else {
			break
		}
	}

	return data, nil
}
