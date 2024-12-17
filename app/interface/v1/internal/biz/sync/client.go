// Package sync
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
package sync

import (
	"context"
	"fmt"

	"github.com/bilibili/HCP/app/interface/v1/configs"
	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	common "github.com/bilibili/HCP/common/models"

	_ "github.com/bilibili/HCP/cloudrepo/provider/alicloud"
	_ "github.com/bilibili/HCP/cloudrepo/provider/awscloud"
	_ "github.com/bilibili/HCP/cloudrepo/provider/baiducloud"
	_ "github.com/bilibili/HCP/cloudrepo/provider/huaweicloud"
	_ "github.com/bilibili/HCP/cloudrepo/provider/tencentcloud"
)

// syncJob 资源类型结构体
type syncJob struct {
	ResourceType string
}

// Client 资源同步客户端结构体
type Client struct {
	Account                *biz.AccountUseCase
	Provider               *biz.ProviderUseCase
	CloudProject           *biz.CloudProjectUseCase
	CloudRegion            *biz.CloudRegionUseCase
	CloudRegionAssociation *biz.CloudRegionAssociationUseCase
	CloudSecurityGroup     *biz.CloudSecurityGroupUseCase
	CloudZone              *biz.CloudZoneUseCase
	CloudVpc               *biz.CloudVpcUseCase
	CloudSubnet            *biz.CloudSubnetUseCase
	CloudServerImage       *biz.CloudServerImageUseCase
	CloudServerSpec        *biz.CloudServerSpecUseCase
	CloudServer            *biz.CloudServerUseCase
	ProjectAccountConfig   *biz.ProjectAccountConfigUseCase
	syncCloudWorkCh        chan syncJob
	syncCloudWorkQuickCh   chan syncJob
	syncCloudWorkSlowCh    chan syncJob
	ResourceSyncFunc       map[string]func(ctx context.Context, req *common.QueryCloudResourceRequest) error
}

// NewClient 创建资源同步客户端
func NewClient(account *biz.AccountUseCase, provider *biz.ProviderUseCase, cloudProject *biz.CloudProjectUseCase, cloudRegion *biz.CloudRegionUseCase, cloudRegionAssociation *biz.CloudRegionAssociationUseCase,
	cloudSecurityGroup *biz.CloudSecurityGroupUseCase, cloudZone *biz.CloudZoneUseCase, cloudVpc *biz.CloudVpcUseCase, cloudSubnet *biz.CloudSubnetUseCase, cloudServerImage *biz.CloudServerImageUseCase,
	cloudServerSpec *biz.CloudServerSpecUseCase, cloudServer *biz.CloudServerUseCase, projectAccountConfig *biz.ProjectAccountConfigUseCase) *Client {
	c := &Client{
		Account:                account,
		Provider:               provider,
		CloudProject:           cloudProject,
		CloudRegion:            cloudRegion,
		CloudRegionAssociation: cloudRegionAssociation,
		CloudSecurityGroup:     cloudSecurityGroup,
		CloudZone:              cloudZone,
		CloudVpc:               cloudVpc,
		CloudSubnet:            cloudSubnet,
		CloudServerImage:       cloudServerImage,
		CloudServerSpec:        cloudServerSpec,
		CloudServer:            cloudServer,
		ProjectAccountConfig:   projectAccountConfig,
		syncCloudWorkCh:        make(chan syncJob, configs.Conf.CloudSync.ConcurrencyAccount),
		syncCloudWorkQuickCh:   make(chan syncJob, configs.Conf.CloudSync.ConcurrencyAccount),
		syncCloudWorkSlowCh:    make(chan syncJob),
	}
	c.ResourceSyncFunc = map[string]func(ctx context.Context, req *common.QueryCloudResourceRequest) error{
		common.CloudServer:        c.SyncCloudServer,        //云服务器
		common.CloudProject:       c.SyncCloudProject,       //云项目
		common.CloudRegion:        c.SyncCloudRegion,        //云地域
		common.CloudSecurityGroup: c.SyncCloudSecurityGroup, //云安全组
		common.CloudSubnet:        c.SyncCloudSubnet,        //云子网
		common.CloudVpc:           c.SyncCloudVpc,           //云专有网路
		common.CloudZone:          c.SyncCloudZone,          //云可用区
		common.CloudServerImage:   c.SyncCloudServerImage,   //云镜像
		common.CloudServerSpec:    c.SyncCloudServerSpec,    //云服务器规格
	}
	//启动资源同步任务
	c.StartSync()
	c.StartSlowWorker()
	c.StartQuickWorker()
	c.StartWorker()
	return c
}

// createSyncCloudRequest 获取同步云资源请求参数
func (c *Client) createSyncCloudRequest(req *common.QueryCloudResourceRequest) *common.QueryCloudResourceRequest {
	var (
		accountID     int
		accountCID    string
		providerID    int
		providerAlias string
		regionID      int
	)
	//获取云账号数据
	if req.AccountCID != "" && req.AccountName != "" {
		conditions := map[string]interface{}{
			"cid":  req.AccountCID,
			"name": req.AccountName,
		}
		accountData, err := c.Account.QueryAccount(&biz.AccountWhere{Conditions: conditions}, &biz.AccountOutput{})
		if err != nil {
			return &common.QueryCloudResourceRequest{}
		}
		if accountData != nil {
			accountID = accountData[0].ID
			accountCID = accountData[0].CID
		}
	}
	//获取云厂商数据
	if req.CloudID != "" {
		conditions := map[string]interface{}{
			"alias": req.CloudID,
		}
		providerData, err := c.Provider.QueryProvider(&biz.ProviderWhere{Conditions: conditions}, &biz.ProviderOutput{})
		if err != nil {
			return &common.QueryCloudResourceRequest{}
		}
		if providerData != nil {
			providerID = providerData[0].ID
			providerAlias = providerData[0].Alias
		}
	}
	//获取地域数据
	if req.ResourceType != common.CloudRegion {
		conditions := map[string]interface{}{
			"cid":         req.Region,
			"provider_id": providerID,
		}
		regionData, err := c.CloudRegion.QueryCloudRegion(&biz.CloudRegionWhere{Conditions: conditions}, &biz.CloudRegionOutput{})
		if err != nil {
			return &common.QueryCloudResourceRequest{}
		}
		if regionData != nil && len(regionData) > 0 {
			regionID = regionData[0].ID
		}
	}

	return &common.QueryCloudResourceRequest{
		NewClientReq: common.NewCloudClientRequest{
			Region:    req.Region,
			SecretId:  req.NewClientReq.SecretId,
			SecretKey: req.NewClientReq.SecretKey,
		},
		CloudID:        req.CloudID,
		AccountID:      accountID,
		AccountUnionID: fmt.Sprintf("%s-%s", providerAlias, accountCID),
		AccountName:    req.AccountName,
		AccountAlias:   req.AccountAlias,
		Region:         req.Region,
		RegionID:       regionID,
		ResourceID:     req.ResourceID,
		ResourceName:   req.ResourceName,
		ResourceType:   req.ResourceType,
		ProviderID:     providerID,
	}
}
