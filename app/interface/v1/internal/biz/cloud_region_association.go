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

// CloudRegionAssociationRepo 云地域关联关系接口
type CloudRegionAssociationRepo interface {
	// QueryCloudRegionAssociation 查询云地域关联关系
	QueryCloudRegionAssociation(where *CloudRegionAssociationWhere) ([]*CloudRegionAssociation, error)
	// CreateCloudRegionAssociation 创建云地域关联关系
	CreateCloudRegionAssociation(create []*CloudRegionAssociation) ([]*CloudRegionAssociation, error)
	// UpdateCloudRegionAssociation 更新云地域关联关系
	UpdateCloudRegionAssociation(where *CloudRegionAssociationWhere, update *CloudRegionAssociation) error
	// DeleteCloudRegionAssociation 删除云地域关联关系
	DeleteCloudRegionAssociation(deleteID []int) error
	// UpsertCloudRegionAssociation 更新云或插入地域关联关系
	UpsertCloudRegionAssociation(upsert []*CloudRegionAssociationUpsert) error
}

// CloudRegionAssociationWhere 云地域关联关系查询条件
type CloudRegionAssociationWhere struct {
	// 查询语句
	Query string
	// 查询参数
	Arg interface{}
	// 查询条件
	Conditions map[string]interface{}
}

// CloudRegionAssociation 云地域关联关系参数结构体
type CloudRegionAssociation struct {
	// ID
	ID int `gorm:"column:id" json:"id"`
	// 云账号ID
	AccountID int `gorm:"column:account_id" json:"account_id"`
	// 云账号实体
	Account *Account `json:"account"`
	// 云地域ID
	RegionID int `gorm:"column:region_id" json:"region_id"`
	// 云地域实体
	Region *CloudRegion `json:"region"`
}

// CloudRegionAssociationUpsert 云地域关联关系更新或插入结构体
type CloudRegionAssociationUpsert struct {
	CloudRegionAssociation
}

// CloudRegionAssociationUseCase 云地域关联关系业务逻辑接口
type CloudRegionAssociationUseCase struct {
	repo CloudRegionAssociationRepo
}

// NewCloudRegionAssociationUseCase 创建云地域关联关系业务逻辑接口
func NewCloudRegionAssociationUseCase(repo CloudRegionAssociationRepo) *CloudRegionAssociationUseCase {
	return &CloudRegionAssociationUseCase{repo: repo}
}

// QueryCloudRegionAssociation 查询云地域关联关系
func (c *CloudRegionAssociationUseCase) QueryCloudRegionAssociation(where *CloudRegionAssociationWhere) ([]*CloudRegionAssociation, error) {
	return c.repo.QueryCloudRegionAssociation(where)
}

// CreateCloudRegionAssociation 创建云地域关联关系
func (c *CloudRegionAssociationUseCase) CreateCloudRegionAssociation(create []*CloudRegionAssociation) ([]*CloudRegionAssociation, error) {
	return c.repo.CreateCloudRegionAssociation(create)
}

// UpdateCloudRegionAssociation 更新云地域关联关系
func (c *CloudRegionAssociationUseCase) UpdateCloudRegionAssociation(where *CloudRegionAssociationWhere, update *CloudRegionAssociation) error {
	return c.repo.UpdateCloudRegionAssociation(where, update)
}

// DeleteCloudRegionAssociation 删除云地域关联关系
func (c *CloudRegionAssociationUseCase) DeleteCloudRegionAssociation(deleteID []int) error {
	return c.repo.DeleteCloudRegionAssociation(deleteID)
}

// UpsertCloudRegionAssociation 更新云或插入地域关联关系
func (c *CloudRegionAssociationUseCase) UpsertCloudRegionAssociation(upsert []*CloudRegionAssociationUpsert) error {
	return c.repo.UpsertCloudRegionAssociation(upsert)
}
