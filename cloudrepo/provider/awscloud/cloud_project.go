// Package awscloud
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
package awscloud

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/resourcegroups"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroups/types"
	"github.com/bilibili/HCP/cloudrepo"
)

const AwsDefaultProject = "默认项目"

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
func (c *Client) ListProject(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudProject, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFn()
	cli := c.clientResourceGroups()
	var data []types.GroupIdentifier
	pageSize := 50
	maxResult := int32(pageSize)
	request := &resourcegroups.ListGroupsInput{MaxResults: &maxResult}
	for {
		if req.Cursor != "" {
			request.NextToken = &req.Cursor
		}
		resp, err := cli.ListGroups(ctx, request)
		if err != nil {
			return nil, err
		}
		if resp.GroupIdentifiers != nil && len(resp.GroupIdentifiers) > 0 {
			data = append(data, resp.GroupIdentifiers...)
		}
		if resp.NextToken != nil {
			req.Cursor = *resp.NextToken
		}
		if resp.NextToken == nil {
			break
		}
	}
	return list2Do(data)
}

func list2Do(resp []types.GroupIdentifier) (list []cloudrepo.CloudProject, err error) {
	for _, v := range resp {
		_ = v
		list = append(list, &CloudProject{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				Name: *v.GroupName,
				CID:  *v.GroupName,
			},
		})
	}
	list = append(list, &CloudProject{
		CloudProductCommon: cloudrepo.CloudProductCommon{
			Name: AwsDefaultProject,
			CID:  AwsDefaultProject,
		},
	})
	return
}
