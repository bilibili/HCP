// Package baiducloud
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
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/tag"
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

// ListTag 获取标签列表
func (c *Client) ListTag(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudTag, error) {
	cli, err := c.clientTag()
	if err != nil {
		return nil, err
	}
	data, err := cli.UserTagList("", "")
	if err != nil {
		return nil, err
	}
	return list2DoTag(data.Tags)
}

func list2DoTag(resp []tag.Tag) (list []cloudrepo.CloudTag, err error) {
	for _, v := range resp {
		_ = v
		list = append(list, &CloudTag{
			tagKey:   v.TagKey,
			tagValue: v.TagValue,
		})
	}
	return
}
