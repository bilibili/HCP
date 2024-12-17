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

// ChargeTypeRepo 计费方式接口
type ChargeTypeRepo interface {
	// QueryChargeType 查询计费方式
	QueryChargeType(where *ChargeTypeWhere, output *ChargeTypeOutput) ([]*ChargeType, error)
	// CountChargeType 查询计费方式数量
	CountChargeType(where *ChargeTypeWhere) (int64, error)
	// CreateChargeType 创建计费方式
	CreateChargeType(create []*ChargeType) ([]*ChargeType, error)
	// UpdateChargeType 更新计费方式
	UpdateChargeType(where *ChargeTypeWhere, update *ChargeType) error
	// DeleteChargeType 删除计费方式
	DeleteChargeType(deleteID []int) error
	// UpsertChargeType 更新或插入计费方式
	UpsertChargeType(upsert []*ChargeTypeUpsert) error
}

// ChargeTypeWhere 计费方式查询条件
type ChargeTypeWhere struct {
	// 查询语句
	Query string
	// 查询参数
	Arg interface{}
	// 查询条件
	Conditions map[string]interface{}
}

// ChargeTypeOutput 计费方式输出条件参数
type ChargeTypeOutput struct {
	OutPutCommon
}

// ChargeType 计费方式参数结构体
type ChargeType struct {
	// ID
	ID int `gorm:"column:id" json:"id"`
	// 云ID
	CID string `gorm:"column:cid" json:"cid"`
	// 名称
	Name string `gorm:"column:name" json:"name"`
	// 类型
	Type int `gorm:"column:type" json:"type"`
	// 云厂商ID
	ProviderID int `gorm:"column:provider_id" json:"provider_id"`
	// 云厂商实体
	Provider *Provider `json:"provider"`
	// 产品类型
	ProductType string `gorm:"column:product_type" json:"product_type"`
}

// ChargeTypeUpsert 计费方式更新或插入结构体
type ChargeTypeUpsert struct {
	ChargeType
}

// ChargeTypeUseCase 计费方式业务逻辑
type ChargeTypeUseCase struct {
	repo ChargeTypeRepo
}

// NewChargeTypeUseCase 创建计费方式业务逻辑
func NewChargeTypeUseCase(repo ChargeTypeRepo) *ChargeTypeUseCase {
	return &ChargeTypeUseCase{repo: repo}
}

// QueryChargeType 查询计费方式
func (c *ChargeTypeUseCase) QueryChargeType(where *ChargeTypeWhere, output *ChargeTypeOutput) ([]*ChargeType, error) {
	return c.repo.QueryChargeType(where, output)
}

// CountChargeType 查询计费方式数量
func (c *ChargeTypeUseCase) CountChargeType(where *ChargeTypeWhere) (int64, error) {
	return c.repo.CountChargeType(where)
}

// CreateChargeType 创建计费方式
func (c *ChargeTypeUseCase) CreateChargeType(create []*ChargeType) ([]*ChargeType, error) {
	return c.repo.CreateChargeType(create)
}

// UpdateChargeType 更新计费方式
func (c *ChargeTypeUseCase) UpdateChargeType(where *ChargeTypeWhere, update *ChargeType) error {
	return c.repo.UpdateChargeType(where, update)
}

// DeleteChargeType 删除计费方式
func (c *ChargeTypeUseCase) DeleteChargeType(deleteID []int) error {
	return c.repo.DeleteChargeType(deleteID)
}

// UpsertChargeType 更新或插入计费方式
func (c *ChargeTypeUseCase) UpsertChargeType(upsert []*ChargeTypeUpsert) error {
	return c.repo.UpsertChargeType(upsert)
}
