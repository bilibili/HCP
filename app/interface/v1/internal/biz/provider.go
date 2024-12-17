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

import common "github.com/bilibili/HCP/common/models"

// ProviderRepo 厂商接口
type ProviderRepo interface {
	// CountProvider 查询厂商数量
	CountProvider(where *ProviderWhere) (int64, error)
	// QueryProvider 查询厂商
	QueryProvider(where *ProviderWhere, output *ProviderOutput) ([]*Provider, error)
	// CreateProvider 创建厂商
	CreateProvider(create []*Provider) ([]*Provider, error)
	// UpdateProvider 更新厂商
	UpdateProvider(where *ProviderWhere, update *Provider) error
	// DeleteProvider 删除厂商
	DeleteProvider(deleteID []int) error
	// UpsertProvider 更新或插入厂商
	UpsertProvider(upsert []*ProviderUpsert) error
}

// ProviderWhere 厂商查询条件
type ProviderWhere struct {
	// 查询语句
	Query string
	// 查询参数
	Arg interface{}
	// 查询条件
	Conditions map[string]interface{}
}

// ProviderOutput 厂商查询输出条件参数
type ProviderOutput struct {
	OutPutCommon
}

// Provider 厂商参数结构体
type Provider struct {
	// 操作的对象ID
	FormObject common.FormObject `gorm:"-" json:"formObject"`
	// ID
	ID int `gorm:"column:id" json:"id"`
	// 名称
	Name string `gorm:"column:name" json:"name"`
	// 别名
	Alias string `gorm:"column:alias" json:"alias"`
	// 售前联系人
	PreSaleContacts string `gorm:"column:pre_sale_contacts" json:"pre_sale_contacts"`
	// 售后联系人
	AfterSaleContacts string `gorm:"column:after_sale_contacts" json:"after_sale_contacts"`
	// 商务联系人
	BusinessContacts string `gorm:"column:business_contacts" json:"business_contacts"`
	// 关联账号数量
	AccountNumber int `gorm:"-" json:"account_number"`
}

// ProviderUpsert 厂商更新或插入结构体
type ProviderUpsert struct {
	Provider
}

// ProviderUseCase 厂商业务逻辑
type ProviderUseCase struct {
	repo ProviderRepo
}

// NewProviderUseCase 创建厂商业务逻辑
func NewProviderUseCase(repo ProviderRepo) *ProviderUseCase {
	return &ProviderUseCase{repo: repo}
}

// QueryProvider 查询厂商
func (c *ProviderUseCase) QueryProvider(where *ProviderWhere, output *ProviderOutput) ([]*Provider, error) {
	return c.repo.QueryProvider(where, output)
}

// CountProvider 查询厂商数量
func (c *ProviderUseCase) CountProvider(where *ProviderWhere) (int64, error) {
	return c.repo.CountProvider(where)
}

// CreateProvider 创建厂商
func (c *ProviderUseCase) CreateProvider(create []*Provider) ([]*Provider, error) {
	return c.repo.CreateProvider(create)
}

// UpdateProvider 更新厂商
func (c *ProviderUseCase) UpdateProvider(where *ProviderWhere, update *Provider) error {
	return c.repo.UpdateProvider(where, update)
}

// DeleteProvider 删除厂商
func (c *ProviderUseCase) DeleteProvider(deleteID []int) error {
	return c.repo.DeleteProvider(deleteID)
}

// UpsertProvider 更新或插入厂商
func (c *ProviderUseCase) UpsertProvider(upsert []*ProviderUpsert) error {
	return c.repo.UpsertProvider(upsert)
}
