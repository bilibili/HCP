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

// AccountRepo 云账号接口
type AccountRepo interface {
	// QueryAccount 查询云账号
	QueryAccount(where *AccountWhere, output *AccountOutput) ([]*Account, error)
	// CountAccount 查询云账号数量
	CountAccount(where *AccountWhere) (int64, error)
	// CreateAccount 创建云账号
	CreateAccount(create []*Account) ([]*Account, error)
	// UpdateAccount 更新云账号
	UpdateAccount(where *AccountWhere, update *Account) error
	// DeleteAccount 删除云账号
	DeleteAccount(softDelete int, deleteID []int) error
	// UpsertAccount 更新或插入云账号
	UpsertAccount(upsert []*AccountUpsert) error
}

// AccountWhere 云账号查询条件
type AccountWhere struct {
	// 查询语句
	Query string
	// 查询参数
	Arg interface{}
	// 查询条件
	Conditions map[string]interface{}
}

// AccountOutput 云账号查询输出条件参数
type AccountOutput struct {
	OutPutCommon
}

// AccountDelete 云账号删除参数
type AccountDelete struct {
	FormObject []common.FormObject `gorm:"-" json:"formObject"`
}

// Account 云账号参数结构体
type Account struct {
	CloudProductCommon
	// 选中的列表ID
	FormObject common.FormObject `gorm:"-" json:"formObject"`
	// 别名
	Alias string `gorm:"column:alias" json:"alias"`
	// 云厂商ID
	ProviderID int `gorm:"column:provider_id" json:"provider_id"`
	// 云厂商实体
	Provider *Provider `json:"provider"`
	// 云账号操作ID
	OperateSecretID string `gorm:"column:operate_secret_id" json:"operate_secret_id"`
	// 云账号操作密钥
	OperateSecretKey string `gorm:"column:operate_secret_key" json:"operate_secret_key"`
	// 云账号同步ID
	SyncSecretID string `gorm:"column:sync_secret_id" json:"sync_secret_id"`
	// 云账号同步密钥
	SyncSecretKey string `gorm:"column:sync_secret_key" json:"sync_secret_key"`
	// 云账号密钥对ID
	KeyPairID string `gorm:"column:key_pair_id" json:"key_pair_id"`
	// 云账号密钥对名称
	KeyPairName string `gorm:"column:key_pair_name" json:"key_pair_name"`
	// 邮箱
	Email string `gorm:"column:email" json:"email"`
	// 电话
	Phone string `gorm:"column:phone" json:"phone"`
	// 描述
	Description string `gorm:"column:description" json:"description"`
	// 主体
	Subject string `gorm:"column:subject" json:"subject"`
	// 是否删除
	IsDelete int `gorm:"column:is_delete" json:"is_delete"`
	// 创建时间
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	// 删除时间
	DeleteTime time.Time `gorm:"column:delete_time" json:"delete_time"`
	// 更新时间
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

// AccountUpsert 云账号更新或插入结构体
type AccountUpsert struct {
	Account
}

// AccountUseCase 云账号业务逻辑
type AccountUseCase struct {
	repo AccountRepo
}

// NewAccountUseCase 创建云账号业务逻辑
func NewAccountUseCase(repo AccountRepo) *AccountUseCase {
	return &AccountUseCase{repo: repo}
}

// QueryAccount 查询云账号
func (c *AccountUseCase) QueryAccount(where *AccountWhere, output *AccountOutput) ([]*Account, error) {
	return c.repo.QueryAccount(where, output)
}

// CountAccount 查询云账号数量
func (c *AccountUseCase) CountAccount(where *AccountWhere) (int64, error) {
	return c.repo.CountAccount(where)
}

// CreateAccount 创建云账号
func (c *AccountUseCase) CreateAccount(create []*Account) ([]*Account, error) {
	return c.repo.CreateAccount(create)
}

// UpdateAccount 更新云账号
func (c *AccountUseCase) UpdateAccount(where *AccountWhere, update *Account) error {
	return c.repo.UpdateAccount(where, update)
}

// DeleteAccount 删除云账号
func (c *AccountUseCase) DeleteAccount(softDelete int, deleteID []int) error {
	return c.repo.DeleteAccount(softDelete, deleteID)
}

// UpsertAccount 更新或插入云账号
func (c *AccountUseCase) UpsertAccount(upsert []*AccountUpsert) error {
	return c.repo.UpsertAccount(upsert)
}
