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

// UserRepo 用户接口
type UserRepo interface {
	// CountUser 查询用户数量
	CountUser(where *UserWhere) (int64, error)
	// QueryUser 查询用户
	QueryUser(where *UserWhere, output *UserOutput, field string) ([]*User, error)
	// CreateUser 创建用户
	CreateUser(create []*User) ([]*User, error)
	// UpdateUser 更新用户
	UpdateUser(where *UserWhere, update *User) error
	// DeleteUser 删除用户
	DeleteUser(softDelete int, deleteID []int) error
	// UpsertUser 更新或插入用户
	UpsertUser(upsert []*UserUpsert) error
}

// UserWhere 用户查询条件
type UserWhere struct {
	// 查询语句
	Query string
	// 查询参数
	Arg interface{}
	// 查询条件
	Conditions map[string]interface{}
}

// UserOutput 用户查询输出条件参数
type UserOutput struct {
	OutPutCommon
}

// UserDelete 用户删除参数ID
type UserDelete struct {
	FormObject []common.FormObject `gorm:"-" json:"formObject"`
}

// User 用户参数结构体
type User struct {
	// 要操作的对象ID
	FormObject common.FormObject `gorm:"-" json:"formObject"`
	// ID
	ID int `gorm:"primary_key;column:id" json:"id"`
	// 名称
	Name string `gorm:"column:name" json:"name"`
	// 密码
	Password string `gorm:"column:password" json:"password"`
	// 角色
	Role int `gorm:"column:role" json:"role"`
	// 盐
	Salt string `gorm:"column:salt" json:"salt"`
	// 是否删除
	IsDelete int `gorm:"column:is_delete" json:"is_delete"`
	// 创建时间
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	// 删除时间
	DeleteTime time.Time `gorm:"column:delete_time" json:"delete_time"`
	// 更新时间
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

// UserUpsert 用户更新或插入结构体
type UserUpsert struct {
	User
}

// UserUseCase 用户业务逻辑
type UserUseCase struct {
	repo UserRepo
}

// NewUserUseCase 创建用户业务逻辑
func NewUserUseCase(repo UserRepo) *UserUseCase {
	return &UserUseCase{repo: repo}
}

// QueryUser 查询用户
func (c *UserUseCase) QueryUser(where *UserWhere, output *UserOutput, field string) ([]*User, error) {
	return c.repo.QueryUser(where, output, field)
}

// CountUser 查询用户数量
func (c *UserUseCase) CountUser(where *UserWhere) (int64, error) {
	return c.repo.CountUser(where)
}

// CreateUser 创建用户
func (c *UserUseCase) CreateUser(create []*User) ([]*User, error) {
	return c.repo.CreateUser(create)
}

// UpdateUser 更新用户
func (c *UserUseCase) UpdateUser(where *UserWhere, update *User) error {
	return c.repo.UpdateUser(where, update)
}

// DeleteUser 删除用户
func (c *UserUseCase) DeleteUser(softDelete int, deleteID []int) error {
	return c.repo.DeleteUser(softDelete, deleteID)
}

// UpsertUser 更新或插入用户
func (c *UserUseCase) UpsertUser(upsert []*UserUpsert) error {
	return c.repo.UpsertUser(upsert)
}
