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

// CloudZoneRepo 云可用区接口
type CloudZoneRepo interface {
	// QueryCloudZone 查询云可用区
	QueryCloudZone(where *CloudZoneWhere, output *CloudZoneOutput) ([]*CloudZone, error)
	// CountCloudZone 查询云可用区数量
	CountCloudZone(where *CloudZoneWhere) (int64, error)
	// CreateCloudZone 创建云可用区
	CreateCloudZone(create []*CloudZone) ([]*CloudZone, error)
	// UpdateCloudZone 更新云可用区
	UpdateCloudZone(where *CloudZoneWhere, update *CloudZone) error
	// DeleteCloudZone 删除云可用区
	DeleteCloudZone(deleteID []int) error
	// UpsertCloudZone 更新或插入云可用区
	UpsertCloudZone(upsert []*CloudZoneUpsert) error
}

// CloudZoneWhere 云可用区查询条件
type CloudZoneWhere struct {
	// 查询语句
	Query string
	// 查询参数
	Arg interface{}
	// 查询条件
	Conditions map[string]interface{}
}

// CloudZoneOutput 云可用区查询输出条件参数
type CloudZoneOutput struct {
	OutPutCommon
}

// CloudZone 云可用区参数结构体
type CloudZone struct {
	CloudProductCommon
	// 云厂商ID
	ProviderID int `gorm:"column:provider_id" json:"provider_id"`
	// 云厂商实体
	Provider *Provider `json:"provider"`
	// 云地域ID
	RegionID int `gorm:"column:region_id" json:"region_id"`
	// 云地域实体
	Region *CloudRegion `json:"region"`
	// 状态
	Status string `gorm:"column:status" json:"status"`
}

// CloudZoneUpsert 云可用区更新或插入结构体
type CloudZoneUpsert struct {
	CloudZone
}

// CloudZoneUseCase 云可用区业务逻辑
type CloudZoneUseCase struct {
	repo CloudZoneRepo
}

// NewCloudZoneUseCase 创建云可用区业务逻辑
func NewCloudZoneUseCase(repo CloudZoneRepo) *CloudZoneUseCase {
	return &CloudZoneUseCase{repo: repo}
}

// QueryCloudZone 查询云可用区
func (c *CloudZoneUseCase) QueryCloudZone(where *CloudZoneWhere, output *CloudZoneOutput) ([]*CloudZone, error) {
	return c.repo.QueryCloudZone(where, output)
}

// CountCloudZone 查询云可用区数量
func (c *CloudZoneUseCase) CountCloudZone(where *CloudZoneWhere) (int64, error) {
	return c.repo.CountCloudZone(where)
}

// CreateCloudZone 创建云可用区
func (c *CloudZoneUseCase) CreateCloudZone(create []*CloudZone) ([]*CloudZone, error) {
	return c.repo.CreateCloudZone(create)
}

// UpdateCloudZone 更新云可用区
func (c *CloudZoneUseCase) UpdateCloudZone(where *CloudZoneWhere, update *CloudZone) error {
	return c.repo.UpdateCloudZone(where, update)
}

// DeleteCloudZone 删除云可用区
func (c *CloudZoneUseCase) DeleteCloudZone(deleteID []int) error {
	return c.repo.DeleteCloudZone(deleteID)
}

// UpsertCloudZone 更新或插入云可用区
func (c *CloudZoneUseCase) UpsertCloudZone(upsert []*CloudZoneUpsert) error {
	return c.repo.UpsertCloudZone(upsert)
}

// DiffCloudZone 云可用区输入数据与存量数据差异
func (c *CloudZoneUseCase) DiffCloudZone(inputs []*CloudZone, where *CloudZoneWhere) (creates []*CloudZone, updates []*CloudZone, deletes []*CloudZone, err error) {
	lists, err := c.QueryCloudZone(where, &CloudZoneOutput{})
	if err != nil {
		return nil, nil, nil, err
	}
	existMap := map[string]CloudZone{}
	for _, find := range lists {
		existMap[fmt.Sprintf("%v-%v-%v", find.CID, find.ProviderID, find.RegionID)] = *find
	}
	currentMap := map[string]CloudZone{}
	for _, v := range inputs {
		currentMap[fmt.Sprintf("%v-%v-%v", v.CID, v.ProviderID, v.RegionID)] = *v

		found, exist := existMap[fmt.Sprintf("%v-%v-%v", v.CID, v.ProviderID, v.RegionID)]
		if exist {
			tempID := found.CloudProductCommon.ID
			found.CloudProductCommon.ID = 0
			found.Provider = nil
			found.Region = nil
			diff := false
			diff, err = DiffNew(found, *v)
			if err != nil {
				continue
			}
			if !diff {
				continue
			}
			updates = append(updates, &CloudZone{
				CloudProductCommon: CloudProductCommon{
					ID:   tempID,
					Name: v.Name,
					CID:  v.CID,
				},
				ProviderID: v.ProviderID,
				RegionID:   v.RegionID,
				Status:     v.Status,
			})
		} else {
			creates = append(creates, &CloudZone{
				CloudProductCommon: CloudProductCommon{
					Name: v.Name,
					CID:  v.CID,
				},
				ProviderID: v.ProviderID,
				RegionID:   v.RegionID,
				Status:     v.Status,
			})
		}
	}
	for _, v := range lists {
		if _, exist := currentMap[fmt.Sprintf("%v-%v-%v", v.CID, v.ProviderID, v.RegionID)]; !exist {
			deletes = append(deletes, &CloudZone{
				CloudProductCommon: CloudProductCommon{ID: v.ID},
			})
		}
	}
	return
}
