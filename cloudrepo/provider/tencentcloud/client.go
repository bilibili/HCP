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
	"fmt"

	"github.com/bilibili/HCP/app/interface/v1/configs"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
	sg "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

type Client struct {
	Region    string
	SecretId  string
	SecretKey string
}

func NewClient(region, secretId, secretKey string) *Client {
	return &Client{
		Region:    region,
		SecretId:  secretId,
		SecretKey: secretKey,
	}
}

func (c *Client) clientProject() (*dcdb.Client, error) {
	cred := common.NewCredential(c.SecretId, c.SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dcdb.tencentcloudapi.com"
	client, err := dcdb.NewClient(cred, c.Region, cpf)
	return client, err
}

func (c *Client) clientTag() (*tag.Client, error) {
	cred := common.NewCredential(c.SecretId, c.SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tag.tencentcloudapi.com"
	client, err := tag.NewClient(cred, c.Region, cpf)
	return client, err
}

func (c *Client) clientEcs() (*cvm.Client, error) {
	cred := common.NewCredential(c.SecretId, c.SecretKey)
	client, err := cvm.NewClient(cred, c.Region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) clientCvm() (*cvm.Client, error) {
	credential := common.NewCredential(c.SecretId, c.SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	client, err := cvm.NewClient(credential, c.Region, cpf)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) clientSecurityGroup() (*sg.Client, error) {
	cred := common.NewCredential(c.SecretId, c.SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"
	client, err := sg.NewClient(cred, c.Region, cpf)
	return client, err
}

func (c *Client) clientVpc() (*vpc.Client, error) {
	cred := common.NewCredential(c.SecretId, c.SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"
	client, err := vpc.NewClient(cred, c.Region, cpf)
	return client, err
}

func (c *Client) GetProject(tags []Tag) string {
	for _, v := range tags {
		if v.Key == configs.Conf.CloudConf.TagProjectKey {
			return fmt.Sprintf("%s-%s", v.Key, v.Value)
		}
	}
	return ""
}
