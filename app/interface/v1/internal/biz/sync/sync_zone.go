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

	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/bilibili/HCP/cloudrepo"
	common "github.com/bilibili/HCP/common/models"
	"github.com/go-kratos/kratos/pkg/log"
)

// SyncCloudZone 同步云可用区
func (c *Client) SyncCloudZone(ctx context.Context, req *common.QueryCloudResourceRequest) error {
	commonRequest := c.createSyncCloudRequest(req)
	if commonRequest == nil {
		log.Error("SyncCloudZone CreateSyncCloudRequest err")
		return common.ReturnErr(common.ErrDatabaseQueryCode, "SyncCloudZone CreateSyncCloudRequest err", common.ErrDatabaseQueryMessage)
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
	zones, err := provider.ListZone(&cloudrepo.GetCloudProductReq{
		ResourceID:   req.ResourceID,
		ResourceName: req.ResourceName,
		Region:       req.Region,
	})
	if err != nil {
		return err
	}
	if len(zones) < 1 {
		return err
	}
	data := make([]*biz.CloudZone, len(zones))
	for i := range zones {
		data[i] = &biz.CloudZone{
			CloudProductCommon: biz.CloudProductCommon{
				Name: zones[i].GetName(),
				CID:  zones[i].GetCID(),
			},
			ProviderID: commonRequest.ProviderID,
			RegionID:   commonRequest.RegionID,
			Status:     zones[i].GetStatus(),
		}
	}
	conditions := map[string]interface{}{
		"provider_id": commonRequest.ProviderID,
		"region_id":   commonRequest.RegionID,
	}
	creates, updates, deletes, err := c.CloudZone.DiffCloudZone(data, &biz.CloudZoneWhere{Conditions: conditions})
	if err != nil {
		return err
	}
	if len(creates) > 0 {
		_, err = c.CloudZone.CreateCloudZone(creates)
		if err != nil {
			return err
		}
	}

	if len(updates) > 0 {
		for _, v := range updates {
			err = c.CloudZone.UpdateCloudZone(&biz.CloudZoneWhere{
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
		err = c.CloudZone.DeleteCloudZone(deleteIds)
		if err != nil {
			return err
		}
	}
	return nil
}
