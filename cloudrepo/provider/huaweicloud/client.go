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
	"strings"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	bss "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/bss/v2"
	bssRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/bss/v2/region"
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	ecsRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/region"
	eps "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/eps/v1"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/eps/v1/region"
	iam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	iamRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/region"
	ims "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2"
	imsRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2/region"
	vpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2"
	vpcRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2/region"
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

func (c *Client) clientProject() *eps.EpsClient {
	au := global.NewCredentialsBuilder().WithAk(c.SecretId).WithSk(c.SecretKey).Build()
	cli := eps.NewEpsClient(
		eps.EpsClientBuilder().
			WithRegion(region.ValueOf("cn-north-4")).
			WithCredential(au).
			Build())
	return cli
}

func (c *Client) clientEcs() *ecs.EcsClient {
	au := basic.NewCredentialsBuilder().WithAk(c.SecretId).WithSk(c.SecretKey).Build()
	builder := ecs.EcsClientBuilder().
		//WithRegion(ecsRegion.ValueOf("cn-south-1")).
		WithRegion(ecsRegion.ValueOf(c.Region)).
		WithCredential(au).
		WithHttpConfig(config.DefaultHttpConfig().WithTimeout(time.Second * 16)).
		Build()
	cli := ecs.NewEcsClient(builder)
	return cli
}

func (c *Client) clientIam() *iam.IamClient {
	au := global.NewCredentialsBuilder().
		WithAk(c.SecretId).
		WithSk(c.SecretKey).
		Build()
	cli := iam.NewIamClient(
		iam.IamClientBuilder().
			WithRegion(iamRegion.ValueOf("cn-east-2")).
			WithCredential(au).
			Build())
	return cli
}

func (c *Client) clientVpc() *vpc.VpcClient {
	au := basic.NewCredentialsBuilder().WithAk(c.SecretId).WithSk(c.SecretKey).Build()
	cli := vpc.NewVpcClient(
		vpc.VpcClientBuilder().
			//WithRegion(vpcRegion.ValueOf("cn-north-1")).
			WithRegion(vpcRegion.ValueOf(c.Region)).
			WithCredential(au).
			WithHttpConfig(config.DefaultHttpConfig().WithTimeout(time.Second * 30)).
			Build())
	return cli
}

func (c *Client) clientVpcWithEndpoint(projectId, endpoint string) *vpc.VpcClient {
	au := basic.NewCredentialsBuilder().WithAk(c.SecretId).WithSk(c.SecretKey).WithProjectId(projectId).Build()
	cli := vpc.NewVpcClient(
		vpc.VpcClientBuilder().
			WithEndpoint(endpoint).
			WithCredential(au).
			WithHttpConfig(config.DefaultHttpConfig().WithTimeout(time.Second * 8)).
			Build())
	return cli
}

func (c *Client) clientEcsWithEndPoint(projectId string, endpoint string) *ecs.EcsClient {
	au := basic.NewCredentialsBuilder().WithAk(c.SecretId).WithSk(c.SecretKey).
		WithProjectId(projectId).
		Build()
	builder := ecs.EcsClientBuilder().
		WithEndpoint(endpoint).
		WithCredential(au).
		WithHttpConfig(config.DefaultHttpConfig().WithTimeout(time.Second * 16)).
		Build()
	cli := ecs.NewEcsClient(builder)
	return cli
}

func (c *Client) clientImage() *ims.ImsClient {
	au := basic.NewCredentialsBuilder().WithAk(c.SecretId).WithSk(c.SecretKey).Build()
	cli := ims.NewImsClient(
		ims.ImsClientBuilder().
			//WithRegion(imsRegion.ValueOf("cn-south-1")).
			WithRegion(imsRegion.ValueOf(c.Region)).
			WithCredential(au).
			WithHttpConfig(config.DefaultHttpConfig().WithTimeout(time.Second * 9)).
			Build())
	return cli
}

func (c *Client) clientBssWithProject() *bss.BssClient {
	au := global.NewCredentialsBuilder().WithAk(c.SecretId).WithSk(c.SecretKey).Build()
	cli := bss.NewBssClient(
		bss.BssClientBuilder().
			WithRegion(bssRegion.ValueOf("cn-north-1")).
			WithCredential(au).
			WithHttpConfig(config.DefaultHttpConfig().WithTimeout(time.Second * 30)).
			Build())
	return cli
}

func (c *Client) clientBss(regionVal string) *bss.BssClient {
	regionCheck := func(region string) bool {
		if strings.HasPrefix(region, "cn") {
			return false
		} else {
			return true
		}
	}
	auth := global.NewCredentialsBuilder().
		WithAk(c.SecretId).
		WithSk(c.SecretKey).
		Build()
	builder := bss.BssClientBuilder().
		WithRegion(bssRegion.ValueOf("cn-north-1")).
		WithCredential(auth).
		Build()
	if regionCheck(regionVal) {
		builder.WithEndpoints([]string{"https://bss-intl.myhuaweicloud.com"})
	}
	return bss.NewBssClient(builder)
}

func (c *Client) clientNewEcs(projectId string) *ecs.EcsClient {
	var (
		auth    *basic.Credentials
		builder *core.HcHttpClient
	)

	if ecsEndPointMap[c.Region] != "" && projectId != "" {
		auth = basic.NewCredentialsBuilder().WithAk(c.SecretId).WithSk(c.SecretKey).WithProjectId(projectId).Build()
		builder = ecs.EcsClientBuilder().
			WithEndpoint(ecsEndPointMap[c.Region]).
			WithCredential(auth).
			WithHttpConfig(config.DefaultHttpConfig().WithTimeout(time.Second * 16)).
			Build()

	} else {
		auth := basic.NewCredentialsBuilder().
			WithAk(c.SecretId).
			WithSk(c.SecretKey).
			//WithProjectId(projectId).
			Build()
		builder = ecs.EcsClientBuilder().
			WithRegion(ecsRegion.ValueOf(c.Region)).
			WithCredential(auth).
			WithHttpConfig(config.DefaultHttpConfig().WithTimeout(time.Second * 16)).
			Build()
	}
	return ecs.NewEcsClient(builder)
}
