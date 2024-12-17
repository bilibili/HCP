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
	"strings"
	"time"

	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/bilibili/HCP/cloudrepo"
	common "github.com/bilibili/HCP/common/models"
	"github.com/go-kratos/kratos/pkg/log"
)

// SyncCloudServer 同步云服务器
func (c *Client) SyncCloudServer(ctx context.Context, req *common.QueryCloudResourceRequest) error {
	commonRequest := c.createSyncCloudRequest(req)
	if commonRequest == nil {
		log.Error("SyncCloudServer CreateSyncCloudRequest err")
		return common.ReturnErr(common.ErrDatabaseQueryCode, "SyncCloudServer CreateSyncCloudRequest err", common.ErrDatabaseQueryMessage)
	}
	repo, err := cloudrepo.GetRepo(req.CloudID)
	if err != nil {
		return err
	}
	provider, err := repo.GetProvider(&cloudrepo.GetProviderReq{
		Region:    req.Region,
		SecretId:  req.NewClientReq.SecretId,
		SecretKey: req.NewClientReq.SecretKey,
	})
	servers, err := provider.ListServer(&cloudrepo.GetCloudProductReq{
		ResourceID:   req.ResourceID,
		ResourceName: req.ResourceName,
	})
	if err != nil {
		return err
	}
	if len(servers) < 1 {
		return err
	}
	projectZoneSpecConditions := map[string]interface{}{
		"provider_id": commonRequest.ProviderID,
	}
	projects, err := c.CloudProject.QueryCloudProject(&biz.CloudProjectWhere{Conditions: projectZoneSpecConditions}, &biz.CloudProjectOutput{})
	if err != nil {
		return err
	}
	existProjectMap := map[string]biz.CloudProject{}
	for _, v := range projects {
		existProjectMap[fmt.Sprintf("%v-%v", v.CID, v.ProviderID)] = *v
	}

	zones, err := c.CloudZone.QueryCloudZone(&biz.CloudZoneWhere{Conditions: projectZoneSpecConditions}, &biz.CloudZoneOutput{})
	if err != nil {
		return err
	}
	existZoneMap := map[string]biz.CloudZone{}
	for _, v := range zones {
		existZoneMap[fmt.Sprintf("%v-%v", v.CID, v.ProviderID)] = *v
	}

	vpcSubnetImageConditions := map[string]interface{}{
		"account_id": commonRequest.AccountID,
	}
	vpcs, err := c.CloudVpc.QueryCloudVpc(&biz.CloudVpcWhere{Conditions: vpcSubnetImageConditions}, &biz.CloudVpcOutput{})
	if err != nil {
		return err
	}
	existVpcMap := map[string]biz.CloudVpc{}
	for _, v := range vpcs {
		existVpcMap[fmt.Sprintf("%v-%v", v.CID, v.AccountID)] = *v
	}

	subnets, err := c.CloudSubnet.QueryCloudSubnet(&biz.CloudSubnetWhere{Conditions: vpcSubnetImageConditions}, &biz.CloudSubnetOutput{})
	if err != nil {
		return err
	}
	existSubnetMap := map[string]biz.CloudSubnet{}
	for _, v := range subnets {
		existSubnetMap[fmt.Sprintf("%v-%v", v.CID, v.AccountID)] = *v
	}

	serverImages, err := c.CloudServerImage.QueryCloudServerImage(&biz.CloudServerImageWhere{Conditions: vpcSubnetImageConditions}, &biz.CloudServerImageOutput{})
	if err != nil {
		return err
	}
	existImageMap := map[string]biz.CloudServerImage{}
	for _, v := range serverImages {
		existImageMap[fmt.Sprintf("%v-%v", v.CID, v.AccountID)] = *v
	}

	projectAccountConfigs, err := c.ProjectAccountConfig.QueryProjectAccountConfig(&biz.ProjectAccountConfigWhere{Conditions: vpcSubnetImageConditions}, &biz.ProjectAccountConfigOutput{})
	if err != nil {
		return err
	}
	existProjectAccountConfigMap := map[string]biz.ProjectAccountConfig{}
	for _, v := range projectAccountConfigs {
		existProjectAccountConfigMap[fmt.Sprintf("%v-%v", v.ProjectID, v.AccountID)] = *v
	}

	serverSpecs, err := c.CloudServerSpec.QueryCloudServerSpec(&biz.CloudServerSpecWhere{Conditions: projectZoneSpecConditions}, &biz.CloudServerSpecOutput{})
	if err != nil {
		return err
	}
	existSpecMap := map[string]biz.CloudServerSpec{}
	for _, v := range serverSpecs {
		existSpecMap[fmt.Sprintf("%v-%v", v.CID, v.AccountID)] = *v
	}

	data := make([]*biz.CloudServer, len(servers))
	for i := range servers {
		data[i] = &biz.CloudServer{
			CloudProductCommon: biz.CloudProductCommon{
				Name: servers[i].GetName(),
				CID:  servers[i].GetCID(),
			},
			AccountID: commonRequest.AccountID,
			RegionID:  commonRequest.RegionID,
		}
		project, existProject := existProjectMap[fmt.Sprintf("%v-%v", servers[i].GetProjectCid(), commonRequest.ProviderID)]
		if existProject {
			data[i].ProjectID = project.ID

			//获取本地项目
			projectConfig, existProjectConfig := existProjectAccountConfigMap[fmt.Sprintf("%v-%v", project.ID, commonRequest.AccountID)]
			if existProjectConfig {
				data[i].ProjectConfigID = projectConfig.ProjectConfigID
			}
		}
		zone, existZone := existZoneMap[fmt.Sprintf("%v-%v", servers[i].GetZoneCid(), commonRequest.ProviderID)]
		if existZone {
			data[i].ZoneID = zone.ID
		}
		vpc, existVpc := existVpcMap[fmt.Sprintf("%v-%v", servers[i].GetVpcCid(), commonRequest.AccountID)]
		if existVpc {
			data[i].VpcID = vpc.ID
		}
		subnet, existSubnet := existSubnetMap[fmt.Sprintf("%v-%v", servers[i].GetSubnetCid(), commonRequest.AccountID)]
		if existSubnet {
			data[i].SubnetID = subnet.ID
		}
		serverImage, existServerImage := existImageMap[fmt.Sprintf("%v-%v", servers[i].GetImageCid(), commonRequest.AccountID)]
		if existServerImage {
			data[i].ServerImageID = serverImage.ID
		}
		serverSpec, existServerSpec := existSpecMap[fmt.Sprintf("%v-%v", servers[i].GetServerSpec(), commonRequest.AccountID)]
		if existServerSpec {
			data[i].ServerSpecID = serverSpec.ID
		}
		data[i].AccountID = commonRequest.AccountID
		data[i].RegionID = commonRequest.RegionID
		data[i].ChargeType = servers[i].GetChangeType()
		data[i].RenewStatus = servers[i].GetRenewStatus()
		data[i].PrivateIp = servers[i].GetPrivateIP()
		data[i].PublicIp = servers[i].GetPublicIP()
		data[i].SecurityGroupCID = strings.Join(servers[i].GetSecurityGroupCid(), ",")
		data[i].SubnetCID = servers[i].GetSubnetCid()
		data[i].Status = servers[i].GetStatus()
		expireTime := servers[i].GetExpireTime()
		if expireTime != "" {
			t, err := time.Parse(time.RFC3339, expireTime)
			if err != nil {
				return err
			}
			expireTime = t.Format("2006-01-02 15:04:05")
		}
		data[i].ExpireTime = expireTime
	}
	conditions := map[string]interface{}{
		"account_id": commonRequest.AccountID,
		"region_id":  commonRequest.RegionID,
	}
	if req.ResourceID != "" {
		conditions["cid"] = req.ResourceID
	}
	creates, updates, deletes, err := c.CloudServer.DiffCloudServer(data, &biz.CloudServerWhere{Conditions: conditions})
	if len(creates) > 0 {
		_, err = c.CloudServer.CreateCloudServer(creates)
		if err != nil {
			return err
		}
	}

	if len(updates) > 0 {
		for _, v := range updates {
			err = c.CloudServer.UpdateCloudServer(&biz.CloudServerWhere{
				Query: "id = ?",
				Arg:   v.ID,
			}, v)
			if err != nil {
				return err
			}
		}
	}
	if len(deletes) > 0 {
		var deleteIds []int
		for _, v := range deletes {
			deleteIds = append(deleteIds, v.ID)
		}
		err = c.CloudServer.DeleteCloudServer(deleteIds)
		if err != nil {
			return err
		}
	}
	return nil
}
