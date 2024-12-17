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
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/eps/v1/model"
	iamModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	"github.com/pkg/errors"
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

// ListProject 获取云项目列表
func (c *Client) ListProject(_ *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudProject, error) {
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
	var data []model.EpDetail
	pageNum, pageSize := 1, 100
	cli := c.clientProject()
	request := model.ListEnterpriseProjectRequest{}
	for {
		offset := int32(pageNum-1) * int32(pageSize)
		limit := int32(pageSize)
		request.Limit = &limit
		request.Offset = &offset
		resp, err := cli.ListEnterpriseProject(&request)
		if err != nil {
			return nil, err
		}
		if resp == nil || len(*resp.EnterpriseProjects) == 0 {
			break
		}
		data = append(data, *resp.EnterpriseProjects...)
		if len(data) >= int(*resp.TotalCount) {
			break
		}
		pageNum += 1
	}
	return list2Do(data)
}

func list2Do(resp []model.EpDetail) (list []cloudrepo.CloudProject, err error) {
	for _, v := range resp {
		_ = v
		list = append(list, &CloudProject{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:  v.Id,
				Name: v.Name,
			},
		})
	}
	return
}

// ListProjectIam 获取云项目Iam列表
func (c *Client) ListProjectIam(req *cloudrepo.GetCloudProductReq) (data []iamModel.AuthProjectResult, err error) {
	req.DisablePage = true
	request := iamModel.KeystoneListAuthProjectsRequest{}
	resp, err := c.clientIam().KeystoneListAuthProjects(&request)
	if err != nil {
		return nil, err
	}
	if resp != nil && resp.Projects != nil {
		return *resp.Projects, nil
	}
	return nil, nil
}
