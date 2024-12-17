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
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/tag"
	"github.com/baidubce/bce-sdk-go/services/vpc"
)

// Client 百度云客户端
type Client struct {
	// 云地域
	Region string
	// 密钥ID
	SecretId string
	// 密钥
	SecretKey string
}

// NewClient 创建百度云客户端
func NewClient(region, secretId, secretKey string) *Client {
	return &Client{
		Region:    region,
		SecretId:  secretId,
		SecretKey: secretKey,
	}
}

func getBaiduRegionEndpointName(region string) string {
	baiduRegionEndpointName := map[string]string{
		"cn-bj":  "bj",
		"cn-gz":  "gz",
		"cn-su":  "su",
		"cn-hkg": "hkg",
		"cn-fwh": "fwh",
		"cn-bd":  "bd",
		"cn-sin": "sin",
		"cn-fsh": "fsh",
	}
	return baiduRegionEndpointName[region]
}

func (c *Client) clientBcc() (*bcc.Client, error) {
	return bcc.NewClient(c.SecretId, c.SecretKey, fmt.Sprintf("bcc.%v.baidubce.com", getBaiduRegionEndpointName(c.Region)))
}

func (c *Client) clientVpc() (*vpc.Client, error) {
	return vpc.NewClient(c.SecretId, c.SecretKey, fmt.Sprintf("bcc.%v.baidubce.com", getBaiduRegionEndpointName(c.Region)))
}

func (c *Client) clientTag() (*tag.Client, error) {
	return tag.NewClient(c.SecretId, c.SecretKey, "tag.baidubce.com")
}
