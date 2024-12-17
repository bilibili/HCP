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

	"github.com/bilibili/HCP/cloudrepo"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
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

// ListTag 获取标签列表
func (c *Client) ListTag(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudTag, error) {
	cli, err := c.clientTag()
	if err != nil {
		return nil, err
	}
	request := tag.NewDescribeTagsRequest()
	resp, err := cli.DescribeTags(request)
	if err != nil {
		return nil, err
	}
	response := &struct {
		Response *struct {
			// 结果总数
			TotalCount uint64 `json:"TotalCount,omitempty" name:"TotalCount"`
			// 数据位移偏量
			Offset uint64 `json:"Offset,omitempty" name:"Offset"`
			// 每页大小
			Limit uint64 `json:"Limit,omitempty" name:"Limit"`
			// 标签列表
			Tags []TagWithDelete `json:"Tags,omitempty" name:"Tags"`
			// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
			RequestId string `json:"RequestId,omitempty" name:"RequestId"`
		} `json:"Response"`
	}{}
	err = json.Unmarshal([]byte(resp.ToJsonString()), response)
	if err != nil {
		return nil, err
	}
	return list2DoTag(response.Response.Tags)
}

func list2DoTag(resp []TagWithDelete) (list []cloudrepo.CloudTag, err error) {
	for _, v := range resp {
		_ = v
		list = append(list, &CloudTag{
			tagKey:   v.TagKey,
			tagValue: v.TagValue,
		})
	}
	return
}
