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

// SyncCloudProject 同步云项目
func (c *Client) SyncCloudProject(ctx context.Context, req *common.QueryCloudResourceRequest) error {
	commonRequest := c.createSyncCloudRequest(req)
	if commonRequest == nil {
		log.Error("SyncCloudProject CreateSyncCloudRequest err")
		return common.ReturnErr(common.ErrDatabaseQueryCode, "SyncCloudProject CreateSyncCloudRequest err", common.ErrDatabaseQueryMessage)
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
	projects, err := provider.ListProject(&cloudrepo.GetCloudProductReq{
		Region:       req.Region,
		ResourceID:   req.ResourceID,
		ResourceName: req.ResourceName,
	})
	if err != nil {
		return err
	}
	if len(projects) < 1 {
		return err
	}
	// 通过标签拼接云项目(百度云部分产品接口未返回项目参数)
	tags, err := provider.ListTag(&cloudrepo.GetCloudProductReq{})
	projectTags := c.CloudProject.ConvertTagToProject(tags, commonRequest.ProviderID, commonRequest.AccountID)

	data := make([]*biz.CloudProject, len(projects))
	for i := range projects {
		data[i] = &biz.CloudProject{
			CloudProductCommon: biz.CloudProductCommon{
				Name: projects[i].GetName(),
				CID:  projects[i].GetCID(),
			},
		}
		data[i].ProviderID = commonRequest.ProviderID
		data[i].AccountID = commonRequest.AccountID
	}
	for _, v := range projectTags {
		data = append(data, v)
	}
	conditions := map[string]interface{}{
		"provider_id": commonRequest.ProviderID,
		"account_id":  commonRequest.AccountID,
	}
	creates, updates, deletes, err := c.CloudProject.DiffCloudProject(commonRequest.ProviderID, data, &biz.CloudProjectWhere{Conditions: conditions})
	if len(creates) > 0 {
		_, err = c.CloudProject.CreateCloudProject(creates)
		if err != nil {
			return err
		}
	}

	if len(updates) > 0 {
		for _, v := range updates {
			err = c.CloudProject.UpdateCloudProject(&biz.CloudProjectWhere{
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
		err = c.CloudProject.DeleteCloudProject(deleteIds)
		if err != nil {
			return err
		}
	}
	return nil
}
