// Package alicloud
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
package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
)

// Client 阿里云客户端
type Client struct {
	// 云地域
	Region string
	// 密钥ID
	SecretId string
	// 密钥
	SecretKey string
}

// NewClient 创建阿里云客户端
func NewClient(region, secretId, secretKey string) *Client {
	return &Client{
		Region:    region,
		SecretId:  secretId,
		SecretKey: secretKey,
	}
}

func (c *Client) clientProject() (*resourcemanager.Client, error) {
	cli, err := resourcemanager.NewClientWithAccessKey(c.Region, c.SecretId, c.SecretKey)
	return cli, err
}

func (c *Client) clientPvtz() (*pvtz.Client, error) {
	cli, err := pvtz.NewClientWithAccessKey(c.Region, c.SecretId, c.SecretKey)
	return cli, err
}

func (c *Client) clientEcs() (*ecs.Client, error) {
	cli, err := ecs.NewClientWithAccessKey(c.Region, c.SecretId, c.SecretKey)
	return cli, err
}

func (c *Client) clientVpc() (*vpc.Client, error) {
	cli, err := vpc.NewClientWithAccessKey(c.Region, c.SecretId, c.SecretKey)
	return cli, err
}

func (c *Client) clientBss(region string) (*bssopenapi.Client, error) {
	if region == "" {
		region = c.Region
	}
	cli, err := bssopenapi.NewClientWithAccessKey(region, c.SecretId, c.SecretKey)
	return cli, err
}
