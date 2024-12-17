// Package biz
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
package biz

import "fmt"

// CloudRegionRepo 云地域接口
type CloudRegionRepo interface {
	// QueryCloudRegion 查询云地域
	QueryCloudRegion(where *CloudRegionWhere, output *CloudRegionOutput) ([]*CloudRegion, error)
	// CountCloudRegion 查询云地域数量
	CountCloudRegion(where *CloudRegionWhere) (int64, error)
	// CreateCloudRegion 创建云地域
	CreateCloudRegion(create []*CloudRegion) ([]*CloudRegion, error)
	// UpdateCloudRegion 更新云地域
	UpdateCloudRegion(where *CloudRegionWhere, update *CloudRegion) error
	// DeleteCloudRegion 删除云地域
	DeleteCloudRegion(deleteID []int) error
	// UpsertCloudRegion 更新或插入云地域
	UpsertCloudRegion(upsert []*CloudRegionUpsert) error
}

// CloudRegionWhere 云地域查询条件
type CloudRegionWhere struct {
	// 查询语句
	Query string
	// 查询参数
	Arg interface{}
	// 查询条件
	Conditions map[string]interface{}
}

// CloudRegionOutput 云地域查询输出条件参数
type CloudRegionOutput struct {
	OutPutCommon
}

// CloudRegion 云地域参数结果体
type CloudRegion struct {
	CloudProductCommon
	// 地域ID
	ProviderID int `gorm:"column:provider_id" json:"provider_id"`
	// 地域实体
	Provider *Provider `json:"provider"`
	// 状态
	Status string `gorm:"column:status" json:"status"`
}

// CloudRegionUpsert 云地域更新或插入数据
type CloudRegionUpsert struct {
	CloudRegion
}

// CloudRegionUseCase 云地域业务逻辑
type CloudRegionUseCase struct {
	repo CloudRegionRepo
}

// NewCloudRegionUseCase 创建云地域业务逻辑
func NewCloudRegionUseCase(repo CloudRegionRepo) *CloudRegionUseCase {
	return &CloudRegionUseCase{repo: repo}
}

// QueryCloudRegion 查询云地域
func (c *CloudRegionUseCase) QueryCloudRegion(where *CloudRegionWhere, output *CloudRegionOutput) ([]*CloudRegion, error) {
	return c.repo.QueryCloudRegion(where, output)
}

// CountCloudRegion 查询云地域数量
func (c *CloudRegionUseCase) CountCloudRegion(where *CloudRegionWhere) (int64, error) {
	return c.repo.CountCloudRegion(where)
}

// CreateCloudRegion 创建云地域
func (c *CloudRegionUseCase) CreateCloudRegion(create []*CloudRegion) ([]*CloudRegion, error) {
	return c.repo.CreateCloudRegion(create)
}

// UpdateCloudRegion 更新云地域
func (c *CloudRegionUseCase) UpdateCloudRegion(where *CloudRegionWhere, update *CloudRegion) error {
	return c.repo.UpdateCloudRegion(where, update)
}

// DeleteCloudRegion 删除云地域
func (c *CloudRegionUseCase) DeleteCloudRegion(deleteID []int) error {
	return c.repo.DeleteCloudRegion(deleteID)
}

// UpsertCloudRegion 更新或插入云地域
func (c *CloudRegionUseCase) UpsertCloudRegion(upsert []*CloudRegionUpsert) error {
	return c.repo.UpsertCloudRegion(upsert)
}

// DiffCloudRegion 云地域输入数据与存量数据差异
func (c *CloudRegionUseCase) DiffCloudRegion(inputs []*CloudRegion, conditions map[string]interface{}) (creates []*CloudRegion, updates []*CloudRegion, deletes []*CloudRegion, err error) {
	lists, err := c.QueryCloudRegion(&CloudRegionWhere{Conditions: conditions}, &CloudRegionOutput{})
	if err != nil {
		return nil, nil, nil, err
	}
	existMap := map[string]CloudRegion{}
	for _, find := range lists {
		existMap[fmt.Sprintf("%v-%v", find.CID, find.ProviderID)] = *find
	}
	currentMap := map[string]CloudRegion{}
	for _, v := range inputs {
		currentMap[fmt.Sprintf("%v-%v", v.CID, v.ProviderID)] = *v

		found, exist := existMap[fmt.Sprintf("%v-%v", v.CID, v.ProviderID)]
		if exist {
			tempID := found.CloudProductCommon.ID
			found.CloudProductCommon.ID = 0
			found.Provider = nil
			diff := false
			diff, err = DiffNew(found, *v)
			if err != nil {
				continue
			}
			if !diff {
				continue
			}
			updates = append(updates, &CloudRegion{
				CloudProductCommon: CloudProductCommon{
					ID:   tempID,
					Name: v.Name,
					CID:  v.CID,
				},
				Status:     v.Status,
				ProviderID: v.ProviderID,
			})
		} else {
			creates = append(creates, &CloudRegion{
				CloudProductCommon: CloudProductCommon{
					Name: v.Name,
					CID:  v.CID,
				},
				Status:     v.Status,
				ProviderID: v.ProviderID,
			})
		}
	}
	for _, v := range lists {
		if _, exist := currentMap[fmt.Sprintf("%v-%v", v.CID, v.ProviderID)]; !exist {
			deletes = append(deletes, &CloudRegion{
				CloudProductCommon: CloudProductCommon{
					ID: v.ID,
				},
			})
		}
	}
	return
}
