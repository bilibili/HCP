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
	"fmt"

	sdkAws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroups"
	"github.com/bilibili/HCP/app/interface/v1/configs"
)

// Client aws客户端
type Client struct {
	// 云地域
	Region string
	// 密钥ID
	SecretId string
	// 密钥
	SecretKey string
}

// NewClient 创建aws客户端
func NewClient(region, secretId, secretKey string) *Client {
	return &Client{
		Region:    region,
		SecretId:  secretId,
		SecretKey: secretKey,
	}
}

func (c *Client) clientResourceGroups() *resourcegroups.Client {
	customCredentials := credentials.NewStaticCredentialsProvider(c.SecretId, c.SecretKey, "")
	client := resourcegroups.NewFromConfig(sdkAws.Config{
		Region:      c.Region,
		Credentials: customCredentials,
	})
	return client
}

func (c *Client) clientEc2() *ec2.Client {
	customCredentials := credentials.NewStaticCredentialsProvider(c.SecretId, c.SecretKey, "")
	client := ec2.NewFromConfig(sdkAws.Config{
		Region:      c.Region,
		Credentials: customCredentials,
	})
	return client
}

func (c *Client) GetNameTag(Tags []ec2Types.Tag) string {
	if len(Tags) > 0 {
		for _, t := range Tags {
			if *t.Key == "Name" {
				return *t.Value
			}
		}
	}
	return ""
}

func (c *Client) GetEc2ProjectTag(Tags []ec2Types.Tag) string {
	if len(Tags) > 0 {
		for _, t := range Tags {
			if *t.Key == "ProjectName" { //兼容AWS自动标签key为ProjectName
				return fmt.Sprintf("%s-%s", *t.Key, *t.Value)
			}
			if *t.Key == configs.Conf.CloudConf.TagProjectKey {
				return fmt.Sprintf("%s-%s", configs.Conf.CloudConf.TagProjectKey, *t.Value)
			}
		}
	}
	return ""
}
