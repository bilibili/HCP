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

var vpcEndpointMap = map[string]string{
	"cn-south-2": "https://vpc.cn-south-2.myhuaweicloud.com",
}

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
	var data []model.Vpc
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
			servers, err := c.queryVpc(req, cli)
			if err != nil {
				return nil, err
			}
			for idx, _ := range servers {
				data = append(data, servers[idx])
			}
		}
		return list2DoVpc(data)
	} else {
		cli = c.clientVpc()
		reps, err := c.queryVpc(req, cli)
		if err != nil {
			return nil, err
		}
		return list2DoVpc(reps)
	}
}

func list2DoVpc(resp []model.Vpc) (list []cloudrepo.CloudVpc, err error) {
	statusMap := map[string]string{
		"OK": "available",
	}
	for _, v := range resp {
		_ = v
		marshalJSON, _ := v.Status.MarshalJSON()
		status := strings.Trim(string(marshalJSON), "\n\"")
		list = append(list, &CloudVpc{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.Id,
				Name:   v.Name,
				Status: statusMap[status],
			},
			Cidr:    v.Cidr,
			Project: v.EnterpriseProjectId,
		})
	}
	return
}

func (c *Client) queryVpc(req *cloudrepo.GetCloudProductReq, cli *vpc.VpcClient) (data []model.Vpc, err error) {
	var marker = new(string)
	*marker = ""
	pageSize := 100
	request := model.ListVpcsRequest{}
	if req.ResourceID != "" {
		request.Id = &req.ResourceID
	}
	for {
		limit := int32(pageSize)
		request.Limit = &limit
		request.Marker = marker
		resp, err := cli.ListVpcs(&request)
		if err != nil {
			return nil, err
		}
		if resp != nil && resp.Vpcs != nil && len(*resp.Vpcs) != 0 {
			data = append(data, *resp.Vpcs...)
			last := (*resp.Vpcs)[len(*resp.Vpcs)-1]
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
