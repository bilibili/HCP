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

import (
	common "github.com/bilibili/HCP/common/models"
)

// TerraformRepo terraform 模版接口
type TerraformRepo interface {
	// CountTerraform 查询terraform模版数量
	CountTerraform(where *TerraformWhere) (int64, error)
	// QueryTerraform 查询terraform模版
	QueryTerraform(where *TerraformWhere, output *TerraformOutput) ([]*Terraform, error)
	// CreateTerraform 创建terraform模版
	CreateTerraform(create []*Terraform) ([]*Terraform, error)
	// UpdateTerraform 更新terraform模版
	UpdateTerraform(where *TerraformWhere, update *Terraform) error
	// DeleteTerraform 删除terraform模版
	DeleteTerraform(deleteID []int) error
	// UpsertTerraform 更新或插入terraform模版
	UpsertTerraform(upsert []*TerraformUpsert) error
}

// TerraformWhere terraform 查询条件
type TerraformWhere struct {
	// 查询条件
	Query string
	// 查询参数
	Arg interface{}
	// 查询条件
	Conditions map[string]interface{}
}

// TerraformOutput terraform 查询输出条件参数
type TerraformOutput struct {
	OutPutCommon
}

// TerraformDelete terraform 删除数据ID
type TerraformDelete struct {
	FormObject []common.FormObject `gorm:"-" json:"formObject"`
}

// Terraform terraform 模版参数结构体
type Terraform struct {
	// ID
	ID int `gorm:"primary_key;column:id" json:"id"`
	// 名称
	Name string `gorm:"column:name" json:"name"`
	// 厂商ID
	ProviderID int `gorm:"column:provider_id" json:"provider_id"`
	// 操作
	Operate string `gorm:"column:operate" json:"operate"`
	// 依赖参数
	DependentParameters string `gorm:"column:dependent_parameters" json:"dependent_parameters"`
	// 模版数据
	Data string `gorm:"column:data" json:"data"`
}

// TerraformUpsert terraform 更新或插入结构体
type TerraformUpsert struct {
	Terraform
}

// TerraformUseCase terraform 业务逻辑接口
type TerraformUseCase struct {
	repo TerraformRepo
}

// NewTerraformUseCase 创建terraform业务逻辑
func NewTerraformUseCase(repo TerraformRepo) *TerraformUseCase {
	return &TerraformUseCase{repo: repo}
}

// QueryTerraform 查询terraform
func (c *TerraformUseCase) QueryTerraform(where *TerraformWhere, output *TerraformOutput) ([]*Terraform, error) {
	return c.repo.QueryTerraform(where, output)
}

// CountTerraform 查询terraform数量
func (c *TerraformUseCase) CountTerraform(where *TerraformWhere) (int64, error) {
	return c.repo.CountTerraform(where)
}

// CreateTerraform 创建terraform
func (c *TerraformUseCase) CreateTerraform(create []*Terraform) ([]*Terraform, error) {
	return c.repo.CreateTerraform(create)
}

// UpdateTerraform 更新terraform
func (c *TerraformUseCase) UpdateTerraform(where *TerraformWhere, update *Terraform) error {
	return c.repo.UpdateTerraform(where, update)
}

// DeleteTerraform 删除terraform
func (c *TerraformUseCase) DeleteTerraform(deleteID []int) error {
	return c.repo.DeleteTerraform(deleteID)
}

// UpsertTerraform 更新或插入terraform
func (c *TerraformUseCase) UpsertTerraform(upsert []*TerraformUpsert) error {
	return c.repo.UpsertTerraform(upsert)
}
