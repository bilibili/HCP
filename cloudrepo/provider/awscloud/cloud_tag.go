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

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/bilibili/HCP/app/interface/v1/configs"
	"github.com/bilibili/HCP/cloudrepo"
)

// CloudTag 云标签
type CloudTag struct {
	ID       int
	tagKey   string
	tagValue string
}

// GetTagKey 获取云标签key
func (c *CloudTag) GetTagKey() string {
	return c.tagKey
}

// GetTagValue 获取云标签value
func (c *CloudTag) GetTagValue() string {
	return c.tagValue
}

// ListTag 获取云标签列表
func (c *Client) ListTag(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudTag, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFn()
	cli := c.clientEc2()
	var data []types.TagDescription
	pageSize := 100
	maxResult := int32(pageSize)
	var filter []types.Filter
	key := "tag-key"
	filter = append(filter, types.Filter{
		Name:   &key,
		Values: []string{configs.Conf.CloudConf.TagProjectKey},
	})
	request := &ec2.DescribeTagsInput{MaxResults: &maxResult, Filters: filter}
	for {
		if req.Cursor != "" {
			request.NextToken = &req.Cursor
		}
		resp, err := cli.DescribeTags(ctx, request)
		if err != nil {
			return nil, err
		}
		if resp.Tags != nil && len(resp.Tags) > 0 {
			data = append(data, resp.Tags...)
		}
		if resp.NextToken == nil {
			break
		}
	}
	return list2DoTag(data)
}

func list2DoTag(resp []types.TagDescription) (list []cloudrepo.CloudTag, err error) {
	for _, v := range resp {
		_ = v
		list = append(list, &CloudTag{
			tagKey:   *v.Key,
			tagValue: *v.Value,
		})
	}
	return
}
