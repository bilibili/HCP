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

// SyncCloudServerSpec 同步云服务器规格
func (c *Client) SyncCloudServerSpec(ctx context.Context, req *common.QueryCloudResourceRequest) error {
	commonRequest := c.createSyncCloudRequest(req)
	if commonRequest == nil {
		log.Error("SyncCloudServerSpec CreateSyncCloudRequest err")
		return common.ReturnErr(common.ErrDatabaseQueryCode, "SyncCloudServerSpec CreateSyncCloudRequest err", common.ErrDatabaseQueryMessage)
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
	serverSpecs, err := provider.ListServerSpec(&cloudrepo.GetCloudProductReq{
		ResourceID:   req.ResourceID,
		ResourceName: req.ResourceName,
	})
	if err != nil {
		return err
	}
	if len(serverSpecs) < 1 {
		return err
	}
	existMap := map[string]cloudrepo.CloudServerSpec{}
	for _, v := range serverSpecs {
		existMap[fmt.Sprintf("%v-%v-%v", v.GetCID(), commonRequest.AccountID, v.GetName())] = v
	}

	var data []*biz.CloudServerSpec
	for _, v := range existMap {
		data = append(data, &biz.CloudServerSpec{
			CloudProductCommon: biz.CloudProductCommon{
				Name: v.GetName(),
				CID:  v.GetCID(),
			},
			AccountID:  commonRequest.AccountID,
			ProviderID: commonRequest.ProviderID,
			Bandwidth:  v.GetBandWidth(),
			Category:   v.GetCategory(),
			CPU:        v.GetCPU(),
			Family:     v.GetFamily(),
			GPU:        v.GetGPU(),
			GPUModel:   v.GetGPUModel(),
			Memory:     v.GetMemory(),
			PPS:        v.GetPPS(),
			Status:     v.GetStatus(),
		})
	}
	conditions := map[string]interface{}{
		"account_id": commonRequest.AccountID,
	}
	creates, updates, _, err := c.CloudServerSpec.DiffCloudServerSpec(data, &biz.CloudServerSpecWhere{Conditions: conditions})
	if len(creates) > 0 {
		_, err = c.CloudServerSpec.CreateCloudServerSpec(creates)
		if err != nil {
			return err
		}
	}

	if len(updates) > 0 {
		for _, v := range updates {
			err = c.CloudServerSpec.UpdateCloudServerSpec(&biz.CloudServerSpecWhere{
				Query: "id = ?",
				Arg:   v.ID,
			}, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
