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

// ProjectUserConfigRepo 本地项目与用户关联关系接口
type ProjectUserConfigRepo interface {
	// CountProjectUserConfig 查询项目与用户关联关系数量
	CountProjectUserConfig(where *ProjectUserConfigWhere, group string) (int64, error)
	// QueryProjectUserConfig 查询项目与用户关联关系
	QueryProjectUserConfig(where *ProjectUserConfigWhere, output *ProjectUserConfigOutput) ([]*ProjectUserConfig, error)
	// CreateProjectUserConfig 创建项目与用户关联关系
	CreateProjectUserConfig(create []*ProjectUserConfig) ([]*ProjectUserConfig, error)
	// UpdateProjectUserConfig 更新项目与用户关联关系
	UpdateProjectUserConfig(where *ProjectUserConfigWhere, update *ProjectUserConfig) error
	// DeleteProjectUserConfig 删除项目与用户关联关系
	DeleteProjectUserConfig(softDelete int, deleteID []int) error
	// DeleteProjectUserConfigByWhere 删除项目与用户关联关系通过条件
	DeleteProjectUserConfigByWhere(where *ProjectUserConfigWhere) error
	// UpsertProjectUserConfig 更新或创建项目与用户关联关系
	UpsertProjectUserConfig(upsert *ProjectUserConfigUpsert) error
}

// ProjectUserConfigWhere 查询项目与用户关联关系条件
type ProjectUserConfigWhere struct {
	// Query 查询语句
	Query string
	// Arg 查询参数
	Arg interface{}
	// Conditions 查询条件
	Conditions map[string]interface{}
}

// ProjectUserConfigOutput 查询项目与用户关联关系输出条件参数
type ProjectUserConfigOutput struct {
	OutPutCommon
}

// ProjectConfigID 项目配置id
type ProjectConfigID struct {
	ID int `gorm:"column:id" json:"id"`
}

// ProjectUserConfig 项目与用户关联关系参数结构体
type ProjectUserConfig struct {
	// 操作对象ID
	FormObject common.FormObject `gorm:"-" json:"formObject"`
	// 项目角色
	ProjectRole string `gorm:"-" json:"project_role"`
	// 目标项目ID列表
	TargetProject []ProjectConfigID `gorm:"-" json:"target_project"`
	// ID
	ID int `gorm:"column:id" json:"id"`
	// 用户ID
	UserID int `gorm:"column:user_id" json:"user_id"`
	// 用户实体
	User *User `json:"user"`
	// 本地项目ID
	ProjectConfigID int `gorm:"column:project_config_id" json:"project_config_id"`
	// 本地项目实体
	ProjectConfig *ProjectConfig `json:"project_config"`
	// 角色
	Role int `gorm:"column:role" json:"role"`
	// 是否删除
	IsDelete int `gorm:"column:is_delete" json:"is_delete"`
	// 创建时间
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	// 删除时间
	DeleteTime time.Time `gorm:"column:delete_time" json:"delete_time"`
	// 更新时间
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

// ProjectUserConfigUpsert 更新或创建项目与用户关联关系结构体
type ProjectUserConfigUpsert struct {
	ProjectUserConfig
}

// ProjectUserConfigUseCase 项目与用户关联关系业务逻辑
type ProjectUserConfigUseCase struct {
	repo ProjectUserConfigRepo
}

// NewProjectUserConfigUseCase 创建项目与用户关联关系业务逻辑
func NewProjectUserConfigUseCase(repo ProjectUserConfigRepo) *ProjectUserConfigUseCase {
	return &ProjectUserConfigUseCase{repo: repo}
}

// QueryProjectUserConfig 查询项目与用户关联关系
func (c *ProjectUserConfigUseCase) QueryProjectUserConfig(where *ProjectUserConfigWhere, output *ProjectUserConfigOutput) ([]*ProjectUserConfig, error) {
	return c.repo.QueryProjectUserConfig(where, output)
}

// CountProjectUserConfig 查询项目与用户关联关系数量
func (c *ProjectUserConfigUseCase) CountProjectUserConfig(where *ProjectUserConfigWhere, group string) (int64, error) {
	return c.repo.CountProjectUserConfig(where, group)
}

// CreateProjectUserConfig 创建项目与用户关联关系
func (c *ProjectUserConfigUseCase) CreateProjectUserConfig(create []*ProjectUserConfig) ([]*ProjectUserConfig, error) {
	return c.repo.CreateProjectUserConfig(create)
}

// UpdateProjectUserConfig 更新项目与用户关联关系
func (c *ProjectUserConfigUseCase) UpdateProjectUserConfig(where *ProjectUserConfigWhere, update *ProjectUserConfig) error {
	return c.repo.UpdateProjectUserConfig(where, update)
}

// DeleteProjectUserConfig 删除项目与用户关联关系
func (c *ProjectUserConfigUseCase) DeleteProjectUserConfig(softDelete int, deleteID []int) error {
	return c.repo.DeleteProjectUserConfig(softDelete, deleteID)
}

// DeleteProjectUserConfigByWhere 删除项目与用户关联关系通过条件
func (c *ProjectUserConfigUseCase) DeleteProjectUserConfigByWhere(where *ProjectUserConfigWhere) error {
	return c.repo.DeleteProjectUserConfigByWhere(where)
}

// UpsertProjectUserConfig 更新或创建项目与用户关联关系
func (c *ProjectUserConfigUseCase) UpsertProjectUserConfig(upsert *ProjectUserConfigUpsert) error {
	return c.repo.UpsertProjectUserConfig(upsert)
}
