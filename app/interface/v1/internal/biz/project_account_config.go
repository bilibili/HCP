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
	"time"

	common "github.com/bilibili/HCP/common/models"
)

// ProjectAccountConfigRepo 云账号和本地项目关联关系接口
type ProjectAccountConfigRepo interface {
	// CountProjectAccountConfig 查询云账号和本地项目关联关系数量
	CountProjectAccountConfig(where *ProjectAccountConfigWhere) (int64, error)
	// QueryProjectAccountConfig 查询云账号和本地项目关联关系
	QueryProjectAccountConfig(where *ProjectAccountConfigWhere, output *ProjectAccountConfigOutput) ([]*ProjectAccountConfig, error)
	// CreateProjectAccountConfig 创建云账号和本地项目关联关系
	CreateProjectAccountConfig(create []*ProjectAccountConfig) ([]*ProjectAccountConfig, error)
	// UpdateProjectAccountConfig 更新云账号和本地项目关联关系
	UpdateProjectAccountConfig(where *ProjectAccountConfigWhere, update *ProjectAccountConfig) error
	// DeleteProjectAccountConfig 删除云账号和本地项目关联关系
	DeleteProjectAccountConfig(deleteID []int) error
	// UpsertProjectAccountConfig 更新或插入云账号和本地项目关联关系
	UpsertProjectAccountConfig(upsert []*ProjectAccountConfigUpsert) error
}

// ProjectAccountConfigWhere 云账号和本地项目关联关系查询条件
type ProjectAccountConfigWhere struct {
	// Query 查询语句
	Query string
	// Arg 查询参数
	Arg interface{}
	// Conditions 查询条件
	Conditions map[string]interface{}
}

// ProjectAccountConfigOutput 云账号和本地项目关联关系查询输出条件参数
type ProjectAccountConfigOutput struct {
	OutPutCommon
}

// ProjectAccountConfig 云账号和本地项目关联关系结构体
type ProjectAccountConfig struct {
	// 操作的对象ID
	FormObject common.FormObject `gorm:"-" json:"formObject"`
	// 关联的云项目
	ConnectCloudProject string `gorm:"-" json:"connectCloudProject"`
	// ID
	ID int `gorm:"column:id" json:"id"`
	// 云账号ID
	AccountID int `gorm:"column:account_id" json:"account_id"`
	// 云账号实体
	Account *Account `json:"account"`
	// 云项目ID
	ProjectID int `gorm:"column:project_id" json:"project_id"`
	// 云项目实体
	Project *CloudProject `json:"project"`
	// 本地项目ID
	ProjectConfigID int `gorm:"column:project_config_id" json:"project_config_id"`
	// 本地项目实体
	ProjectConfig *ProjectConfig `json:"project_config"`
	// 是否删除
	IsDelete int `gorm:"column:is_delete" json:"is_delete"`
	// 创建时间
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	// 删除时间
	DeleteTime time.Time `gorm:"column:delete_time" json:"delete_time"`
	// 更新时间
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

// ProjectAccountConfigUpsert 云账号和本地项目关联关系更新或插入结构体
type ProjectAccountConfigUpsert struct {
	ProjectAccountConfig
}

// ProjectAccountConfigUseCase 云账号和本地项目关联关系业务逻辑
type ProjectAccountConfigUseCase struct {
	repo ProjectAccountConfigRepo
}

// NewProjectAccountConfigUseCase 创建云账号和本地项目关联关系业务逻辑
func NewProjectAccountConfigUseCase(repo ProjectAccountConfigRepo) *ProjectAccountConfigUseCase {
	return &ProjectAccountConfigUseCase{repo: repo}
}

// QueryProjectAccountConfig 查询云账号和本地项目关联关系
func (c *ProjectAccountConfigUseCase) QueryProjectAccountConfig(where *ProjectAccountConfigWhere, output *ProjectAccountConfigOutput) ([]*ProjectAccountConfig, error) {
	return c.repo.QueryProjectAccountConfig(where, output)
}

// CountProjectAccountConfig 查询云账号和本地项目关联关系数量
func (c *ProjectAccountConfigUseCase) CountProjectAccountConfig(where *ProjectAccountConfigWhere) (int64, error) {
	return c.repo.CountProjectAccountConfig(where)
}

// CreateProjectAccountConfig 创建云账号和本地项目关联关系
func (c *ProjectAccountConfigUseCase) CreateProjectAccountConfig(create []*ProjectAccountConfig) ([]*ProjectAccountConfig, error) {
	return c.repo.CreateProjectAccountConfig(create)
}

// UpdateProjectAccountConfig 更新云账号和本地项目关联关系
func (c *ProjectAccountConfigUseCase) UpdateProjectAccountConfig(where *ProjectAccountConfigWhere, update *ProjectAccountConfig) error {
	return c.repo.UpdateProjectAccountConfig(where, update)
}

// DeleteProjectAccountConfig 删除云账号和本地项目关联关系
func (c *ProjectAccountConfigUseCase) DeleteProjectAccountConfig(deleteID []int) error {
	return c.repo.DeleteProjectAccountConfig(deleteID)
}

// UpsertProjectAccountConfig 更新或插入云账号和本地项目关联关系
func (c *ProjectAccountConfigUseCase) UpsertProjectAccountConfig(upsert []*ProjectAccountConfigUpsert) error {
	return c.repo.UpsertProjectAccountConfig(upsert)
}
