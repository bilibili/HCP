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

// SyncCloudServerImage 同步云镜像
func (c *Client) SyncCloudServerImage(ctx context.Context, req *common.QueryCloudResourceRequest) error {
	commonRequest := c.createSyncCloudRequest(req)
	if commonRequest == nil {
		log.Error("SyncCloudServerImage CreateSyncCloudRequest err")
		return common.ReturnErr(common.ErrDatabaseQueryCode, "SyncCloudServerImage CreateSyncCloudRequest err", common.ErrDatabaseQueryMessage)
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
	serverImages, err := provider.ListServerImage(&cloudrepo.GetCloudProductReq{
		ResourceID:   req.ResourceID,
		ResourceName: req.ResourceName,
		Region:       req.Region,
	})
	if err != nil {
		return err
	}
	if len(serverImages) < 1 {
		return err
	}
	existMap := map[string]cloudrepo.CloudServerImage{}
	for _, v := range serverImages {
		existMap[fmt.Sprintf("%v-%v-%v", v.GetCID(), commonRequest.AccountID, commonRequest.RegionID)] = v
	}
	var data []*biz.CloudServerImage
	for _, v := range existMap {
		data = append(data, &biz.CloudServerImage{
			CloudProductCommon: biz.CloudProductCommon{
				Name: v.GetName(),
				CID:  v.GetCID(),
			},
			AccountID: commonRequest.AccountID,
			RegionID:  commonRequest.RegionID,
			OsName:    v.GetOsName(),
			Type:      v.GetType(),
			Status:    v.GetStatus(),
		})
	}
	conditions := map[string]interface{}{
		"account_id": commonRequest.AccountID,
		"region_id":  commonRequest.RegionID,
	}
	creates, updates, deletes, err := c.CloudServerImage.DiffCloudServerImage(data, &biz.CloudServerImageWhere{Conditions: conditions})
	if err != nil {
		return err
	}
	if len(creates) > 0 {
		_, err = c.CloudServerImage.CreatBatchesCloudServerImage(creates)
		if err != nil {
			return err
		}
	}

	if len(updates) > 0 {
		for _, v := range updates {
			err = c.CloudServerImage.UpdateCloudServerImage(&biz.CloudServerImageWhere{
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
		err = c.CloudServerImage.DeleteCloudServerImage(deleteIds)
		if err != nil {
			return err
		}
	}
	return nil
}
