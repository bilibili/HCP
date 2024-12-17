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

// SyncCloudRegion 同步云地域
func (c *Client) SyncCloudRegion(ctx context.Context, req *common.QueryCloudResourceRequest) error {
	var (
		cIds      []string
		regionIds []int
	)
	commonRequest := c.createSyncCloudRequest(req)
	if commonRequest == nil {
		log.Error("SyncCloudRegion CreateSyncCloudRequest err")
		return common.ReturnErr(common.ErrDatabaseQueryCode, "SyncCloudRegion CreateSyncCloudRequest err", common.ErrDatabaseQueryMessage)
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
	regions, err := provider.ListRegion(&cloudrepo.GetCloudProductReq{
		ResourceID:   req.ResourceID,
		ResourceName: req.ResourceName,
	})
	if err != nil {
		return err
	}
	if len(regions) < 1 {
		return err
	}
	data := make([]*biz.CloudRegion, len(regions))
	for i := range regions {
		data[i] = &biz.CloudRegion{
			CloudProductCommon: biz.CloudProductCommon{
				Name: regions[i].GetName(),
				CID:  regions[i].GetCID(),
			},
		}
		data[i].Status = regions[i].GetStatus()
		data[i].ProviderID = commonRequest.ProviderID

		cIds = append(cIds, regions[i].GetCID())
	}
	//获取地域关联关系
	existConnMap := map[string]biz.CloudRegionAssociation{}
	if len(cIds) > 0 {
		regionConditions := map[string]interface{}{
			"cid_IN": cIds,
		}
		regionLists, err := c.CloudRegion.QueryCloudRegion(&biz.CloudRegionWhere{Conditions: regionConditions}, &biz.CloudRegionOutput{})
		if err != nil {
			return err
		}
		for _, v := range regionLists {
			regionIds = append(regionIds, v.ID)
		}

		if len(regionIds) > 0 {
			connConditions := map[string]interface{}{
				"region_id_IN": regionIds,
			}
			connLists, err := c.CloudRegionAssociation.QueryCloudRegionAssociation(&biz.CloudRegionAssociationWhere{Conditions: connConditions})
			if err != nil {
				return err
			}
			for _, find := range connLists {
				existConnMap[fmt.Sprintf("%v-%v", find.AccountID, find.RegionID)] = *find
			}
		}
	}
	conditions := map[string]interface{}{
		"provider_id": commonRequest.ProviderID,
	}
	creates, updates, deletes, err := c.CloudRegion.DiffCloudRegion(data, conditions)
	if err != nil {
		return err
	}
	if len(creates) > 0 {
		dts, err := c.CloudRegion.CreateCloudRegion(creates)
		if err != nil {
			return err
		}
		if len(dts) > 0 {
			var connCreates []*biz.CloudRegionAssociation
			for _, v := range dts {
				connCreates = append(connCreates, &biz.CloudRegionAssociation{
					AccountID: commonRequest.AccountID,
					RegionID:  v.ID,
				})
			}
			_, err = c.CloudRegionAssociation.CreateCloudRegionAssociation(connCreates)
			if err != nil {
				return err
			}
		}
	}

	if len(updates) > 0 {
		for k, v := range updates {
			err = c.CloudRegion.UpdateCloudRegion(&biz.CloudRegionWhere{
				Query: "id = ?",
				Arg:   v.ID,
			}, v)
			if err != nil {
				return err
			}
			_, exist := existConnMap[fmt.Sprintf("%v-%v", commonRequest.AccountID, updates[k].ID)]
			if !exist {
				var connCreates []*biz.CloudRegionAssociation
				connCreates = append(connCreates, &biz.CloudRegionAssociation{
					AccountID: commonRequest.AccountID,
					RegionID:  v.ID,
				})
				_, err = c.CloudRegionAssociation.CreateCloudRegionAssociation(connCreates)
				if err != nil {
					return err
				}
			}
		}
	}
	if len(deletes) > 0 {
		var deleteIds []int
		for _, v := range deletes {
			deleteIds = append(deleteIds, v.ID)
		}
		err = c.CloudRegion.DeleteCloudRegion(deleteIds)
		if err != nil {
			return err
		}
	}
	return nil
}
