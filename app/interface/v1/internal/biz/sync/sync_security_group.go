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

	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/bilibili/HCP/cloudrepo"
	common "github.com/bilibili/HCP/common/models"
	"github.com/go-kratos/kratos/pkg/log"
)

// SyncCloudSecurityGroup 同步云安全组
func (c *Client) SyncCloudSecurityGroup(ctx context.Context, req *common.QueryCloudResourceRequest) error {
	commonRequest := c.createSyncCloudRequest(req)
	if commonRequest == nil {
		log.Error("SyncCloudSecurityGroup CreateSyncCloudRequest err")
		return common.ReturnErr(common.ErrDatabaseQueryCode, "SyncCloudSecurityGroup CreateSyncCloudRequest err", common.ErrDatabaseQueryMessage)
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
	securityGroups, err := provider.ListSecurityGroup(&cloudrepo.GetCloudProductReq{
		ResourceID:   req.ResourceID,
		ResourceName: req.ResourceName,
		Region:       req.Region,
	})
	if err != nil {
		return err
	}
	if len(securityGroups) < 1 {
		return err
	}
	projectConditions := map[string]interface{}{
		"provider_id": commonRequest.ProviderID,
	}
	projects, err := c.CloudProject.QueryCloudProject(&biz.CloudProjectWhere{Conditions: projectConditions}, &biz.CloudProjectOutput{})
	if err != nil {
		return err
	}
	existProjectMap := map[string]biz.CloudProject{}
	for _, v := range projects {
		existProjectMap[fmt.Sprintf("%v-%v", v.CID, v.ProviderID)] = *v
	}

	vpcConditions := map[string]interface{}{
		"account_id": commonRequest.AccountID,
	}
	vpcs, err := c.CloudVpc.QueryCloudVpc(&biz.CloudVpcWhere{Conditions: vpcConditions}, &biz.CloudVpcOutput{})
	if err != nil {
		return err
	}
	existVpcMap := map[string]biz.CloudVpc{}
	for _, v := range vpcs {
		existVpcMap[fmt.Sprintf("%v-%v", v.CID, v.AccountID)] = *v
	}

	data := make([]*biz.CloudSecurityGroup, len(securityGroups))
	for i := range securityGroups {
		data[i] = &biz.CloudSecurityGroup{
			CloudProductCommon: biz.CloudProductCommon{
				Name: securityGroups[i].GetName(),
				CID:  securityGroups[i].GetCID(),
			},
		}
		data[i].Status = securityGroups[i].GetStatus()
		project, exist := existProjectMap[fmt.Sprintf("%v-%v", securityGroups[i].GetProject(), commonRequest.ProviderID)]
		if exist {
			data[i].ProjectID = project.ID
		}
		vpc, existVpc := existVpcMap[fmt.Sprintf("%v-%v", securityGroups[i].GetVpc(), commonRequest.AccountID)]
		if existVpc {
			data[i].VpcID = vpc.ID
		}
		data[i].RegionID = commonRequest.RegionID
		data[i].AccountID = commonRequest.AccountID
	}
	conditions := map[string]interface{}{
		"account_id": commonRequest.AccountID,
		"region_id":  commonRequest.RegionID,
	}
	creates, updates, deletes, err := c.CloudSecurityGroup.DiffCloudSecurityGroup(data, &biz.CloudSecurityGroupWhere{Conditions: conditions})
	if err != nil {
		return err
	}
	if len(creates) > 0 {
		_, err = c.CloudSecurityGroup.CreateCloudSecurityGroup(creates)
		if err != nil {
			return err
		}
	}

	if len(updates) > 0 {
		for _, v := range updates {
			err = c.CloudSecurityGroup.UpdateCloudSecurityGroup(&biz.CloudSecurityGroupWhere{
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
		err = c.CloudSecurityGroup.DeleteCloudSecurityGroup(deleteIds)
		if err != nil {
			return err
		}
	}
	return nil
}
