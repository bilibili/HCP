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

// ProjectConfigRepo 本地项目接口
type ProjectConfigRepo interface {
	// CountProjectConfig 查询本地项目数量
	CountProjectConfig(where *ProjectConfigWhere) (int64, error)
	// QueryProjectConfig 查询本地项目
	QueryProjectConfig(where *ProjectConfigWhere, output *ProjectConfigOutput) ([]*ProjectConfig, error)
	// CreateProjectConfig 创建本地项目
	CreateProjectConfig(create []*ProjectConfig) ([]*ProjectConfig, error)
	// UpdateProjectConfig 更新本地项目
	UpdateProjectConfig(where *ProjectConfigWhere, update *ProjectConfig) error
	// DeleteProjectConfig 删除本地项目
	DeleteProjectConfig(softDelete int, deleteID []int) error
	// UpsertProjectConfig 更新或插入本地项目
	UpsertProjectConfig(upsert []*ProjectConfigUpsert) error
}

// ProjectConfigWhere 查询条件
type ProjectConfigWhere struct {
	// 查询语句
	Query string
	// 查询参数
	Arg interface{}
	// 查询条件
	Conditions map[string]interface{}
}

// ProjectConfigOutput 查询输出条件参数
type ProjectConfigOutput struct {
	OutPutCommon
}

// ProjectConfigDelete 删除ID列表
type ProjectConfigDelete struct {
	FormObject []common.FormObject `gorm:"-" json:"formObject"`
}

// ProjectConfig 本地项目参数结构体
type ProjectConfig struct {
	// 操作的对象ID
	FormObject common.FormObject `gorm:"-" json:"formObject"`
	// ID
	ID int `gorm:"column:id" json:"id"`
	// 名称
	Name string `gorm:"column:name" json:"name"`
	// 别名
	Alias string `gorm:"column:alias" json:"alias"`
	// 描述
	Description string `gorm:"column:description" json:"description"`
	// 状态
	Status string `gorm:"column:status" json:"status"`
	// 是否删除
	IsDelete int `gorm:"column:is_delete" json:"is_delete"`
	// 创建时间
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	// 删除时间
	DeleteTime time.Time `gorm:"column:delete_time" json:"delete_time"`
	// 更新时间
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
	// 研发负责人
	RdLeader interface{} `gorm:"-" json:"rd_leader"`
	// 研发成员
	RdMember interface{} `gorm:"-" json:"rd_member"`
	// 运维负责人
	OpLeader interface{} `gorm:"-" json:"op_leader"`
	// 运维成员
	OpMember interface{} `gorm:"-" json:"op_member"`
}

// ProjectConfigUpsert 更新或插入本地项目结构体
type ProjectConfigUpsert struct {
	ProjectConfig
}

// ProjectConfigUseCase 本地项目业务逻辑
type ProjectConfigUseCase struct {
	repo ProjectConfigRepo
}

// NewProjectConfigUseCase 创建本地项目业务逻辑
func NewProjectConfigUseCase(repo ProjectConfigRepo) *ProjectConfigUseCase {
	return &ProjectConfigUseCase{repo: repo}
}

// QueryProjectConfig 查询本地项目
func (c *ProjectConfigUseCase) QueryProjectConfig(where *ProjectConfigWhere, output *ProjectConfigOutput) ([]*ProjectConfig, error) {
	return c.repo.QueryProjectConfig(where, output)
}

// CountProjectConfig 查询本地项目数量
func (c *ProjectConfigUseCase) CountProjectConfig(where *ProjectConfigWhere) (int64, error) {
	return c.repo.CountProjectConfig(where)
}

// CreateProjectConfig 创建本地项目
func (c *ProjectConfigUseCase) CreateProjectConfig(create []*ProjectConfig) ([]*ProjectConfig, error) {
	return c.repo.CreateProjectConfig(create)
}

// UpdateProjectConfig 更新本地项目
func (c *ProjectConfigUseCase) UpdateProjectConfig(where *ProjectConfigWhere, update *ProjectConfig) error {
	return c.repo.UpdateProjectConfig(where, update)
}

// DeleteProjectConfig 删除本地项目
func (c *ProjectConfigUseCase) DeleteProjectConfig(softDelete int, deleteID []int) error {
	return c.repo.DeleteProjectConfig(softDelete, deleteID)
}

// UpsertProjectConfig 更新或插入本地项目
func (c *ProjectConfigUseCase) UpsertProjectConfig(upsert []*ProjectConfigUpsert) error {
	return c.repo.UpsertProjectConfig(upsert)
}
