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

	"github.com/bilibili/HCP/cloudrepo"
	vpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2/model"
	"github.com/pkg/errors"
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
	var data []model.SecurityGroup
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
			dts, err := c.querySecurityGroup(req, cli)
			if err != nil {
				return nil, err
			}
			for idx, _ := range dts {
				data = append(data, dts[idx])
			}
		}
		return list2DoSecurityGroup(data)
	} else {
		cli = c.clientVpc()
		resp, err := c.querySecurityGroup(req, cli)
		if err != nil {
			return nil, err
		}
		return list2DoSecurityGroup(resp)
	}
}

func list2DoSecurityGroup(resp []model.SecurityGroup) (list []cloudrepo.CloudSecurityGroup, err error) {
	for _, v := range resp {
		_ = v
		vpcCid := ""
		if v.VpcId != nil {
			vpcCid = *v.VpcId
		}
		list = append(list, &CloudSecurityGroup{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.Id,
				Name:   v.Name,
				Status: "available",
			},
			Project: *v.EnterpriseProjectId,
			Vpc:     vpcCid,
		})
	}
	return
}

func (c *Client) querySecurityGroup(req *cloudrepo.GetCloudProductReq, cli *vpc.VpcClient) (data []model.SecurityGroup, err error) {
	var marker = new(string)
	*marker = ""
	pageSize := 100
	request := model.ListSecurityGroupsRequest{}
	for {
		limit := int32(pageSize)
		request.Limit = &limit
		request.Marker = marker
		resp, err := cli.ListSecurityGroups(&request)
		if err != nil {
			return nil, err
		}
		data = append(data, *resp.SecurityGroups...)
		if resp != nil && resp.SecurityGroups != nil && len(*resp.SecurityGroups) > 0 {
			last := (*resp.SecurityGroups)[len(*resp.SecurityGroups)-1]
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
