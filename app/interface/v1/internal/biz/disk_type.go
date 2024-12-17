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

// DiskTypeRepo 数据盘类型接口
type DiskTypeRepo interface {
	// QueryDiskType 查询数据盘类型
	QueryDiskType(where *DiskTypeWhere, output *DiskTypeOutput) ([]*DiskType, error)
	// CountDiskType 查询数据盘类型数量
	CountDiskType(where *DiskTypeWhere) (int64, error)
	// CreateDiskType 创建数据盘类型
	CreateDiskType(create []*DiskType) ([]*DiskType, error)
	// UpdateDiskType 更新数据盘类型
	UpdateDiskType(where *DiskTypeWhere, update *DiskType) error
	// DeleteDiskType 删除数据盘类型
	DeleteDiskType(deleteID []int) error
	// UpsertDiskType 更新或插入数据盘类型
	UpsertDiskType(upsert []*DiskTypeUpsert) error
}

// DiskTypeWhere 查询条件
type DiskTypeWhere struct {
	// 查询语句
	Query string
	// 查询参数
	Arg interface{}
	// 	查询条件
	Conditions map[string]interface{}
}

// DiskTypeOutput 数据盘类型查询输出条件参数
type DiskTypeOutput struct {
	OutPutCommon
}

// DiskType 数据盘类型参数结构体
type DiskType struct {
	// ID
	ID int `gorm:"column:id" json:"id"`
	// 数据盘CID
	CID string `gorm:"column:cid" json:"cid"`
	// 数据盘名称
	Name string `gorm:"column:name" json:"name"`
	// 云厂商ID
	ProviderID int `gorm:"column:provider_id" json:"provider_id"`
	// 云厂商实体
	Provider *Provider `json:"provider"`
	// 数据盘支持产品类型
	ProductType string `gorm:"column:product_type" json:"product_type"`
}

// DiskTypeUpsert 数据盘类型更新或插入结构体
type DiskTypeUpsert struct {
	DiskType
}

// DiskTypeUseCase 数据盘类型业务逻辑
type DiskTypeUseCase struct {
	repo DiskTypeRepo
}

// NewDiskTypeUseCase 创建数据盘类型业务逻辑
func NewDiskTypeUseCase(repo DiskTypeRepo) *DiskTypeUseCase {
	return &DiskTypeUseCase{repo: repo}
}

// QueryDiskType 查询数据盘类型
func (c *DiskTypeUseCase) QueryDiskType(where *DiskTypeWhere, output *DiskTypeOutput) ([]*DiskType, error) {
	return c.repo.QueryDiskType(where, output)
}

// CountDiskType 查询数据盘类型数量
func (c *DiskTypeUseCase) CountDiskType(where *DiskTypeWhere) (int64, error) {
	return c.repo.CountDiskType(where)
}

// CreateDiskType 创建数据盘类型
func (c *DiskTypeUseCase) CreateDiskType(create []*DiskType) ([]*DiskType, error) {
	return c.repo.CreateDiskType(create)
}

// UpdateDiskType 更新数据盘类型
func (c *DiskTypeUseCase) UpdateDiskType(where *DiskTypeWhere, update *DiskType) error {
	return c.repo.UpdateDiskType(where, update)
}

// DeleteDiskType 删除数据盘类型
func (c *DiskTypeUseCase) DeleteDiskType(deleteID []int) error {
	return c.repo.DeleteDiskType(deleteID)
}

// UpsertDiskType 更新或插入数据盘类型
func (c *DiskTypeUseCase) UpsertDiskType(upsert []*DiskTypeUpsert) error {
	return c.repo.UpsertDiskType(upsert)
}
