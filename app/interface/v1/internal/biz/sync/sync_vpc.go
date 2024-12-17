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

// SyncCloudVpc 同步云专有网络
func (c *Client) SyncCloudVpc(ctx context.Context, req *common.QueryCloudResourceRequest) error {
	commonRequest := c.createSyncCloudRequest(req)
	if commonRequest == nil {
		log.Error("SyncCloudVpc CreateSyncCloudRequest err")
		return common.ReturnErr(common.ErrDatabaseQueryCode, "SyncCloudVpc CreateSyncCloudRequest err", common.ErrDatabaseQueryMessage)
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
	vpcs, err := provider.ListVpc(&cloudrepo.GetCloudProductReq{
		ResourceID:   req.ResourceID,
		ResourceName: req.ResourceName,
		Region:       req.Region,
	})
	if err != nil {
		return err
	}
	if len(vpcs) < 1 {
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

	data := make([]*biz.CloudVpc, len(vpcs))
	for i := range vpcs {
		data[i] = &biz.CloudVpc{
			CloudProductCommon: biz.CloudProductCommon{
				Name: vpcs[i].GetName(),
				CID:  vpcs[i].GetCID(),
			},
			AccountID: commonRequest.AccountID,
			RegionID:  commonRequest.RegionID,
			Cidr:      vpcs[i].GetCidr(),
			Status:    vpcs[i].GetStatus(),
		}
		project, exist := existProjectMap[fmt.Sprintf("%v-%v", vpcs[i].GetProject(), commonRequest.ProviderID)]
		if exist {
			data[i].ProjectID = project.ID
		}
	}
	conditions := map[string]interface{}{
		"account_id": commonRequest.AccountID,
		"region_id":  commonRequest.RegionID,
	}
	creates, updates, deletes, err := c.CloudVpc.DiffCloudVpc(data, &biz.CloudVpcWhere{Conditions: conditions})
	if err != nil {
		return err
	}
	if len(creates) > 0 {
		_, err = c.CloudVpc.CreateCloudVpc(creates)
		if err != nil {
			return err
		}
	}

	if len(updates) > 0 {
		for _, v := range updates {
			err = c.CloudVpc.UpdateCloudVpc(&biz.CloudVpcWhere{
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
		err = c.CloudVpc.DeleteCloudVpc(deleteIds)
		if err != nil {
			return err
		}
	}
	return nil
}
