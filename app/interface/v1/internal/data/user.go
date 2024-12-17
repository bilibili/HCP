// Package data
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
package data

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/bilibili/HCP/utils"
	"github.com/go-kratos/kratos/pkg/log"
	"gorm.io/gorm"
)

// UserRepo 用户DB
type UserRepo struct {
	db *biz.DB
}

// NewUserRepo 初始化用户DB
func NewUserRepo(db *biz.DB) biz.UserRepo {
	return &UserRepo{db: db}
}

// ParseUser 用户查询条件处理
func (c *UserRepo) ParseUser(where *biz.UserWhere) (conditions map[string]interface{}, err error) {
	if where.Conditions != nil {
		jsonBytes, err := json.Marshal(where.Conditions)
		if err != nil {
			return nil, err
		}
		conditionDos, err := biz.ParseCloudData(jsonBytes)
		if err != nil {
			log.Error("user table parse data error(%v)", err)
			return nil, err
		}
		conditions = biz.GetHandleConditions(conditionDos, "")
	}
	return conditions, nil
}

// CountUser 用户查询数量
func (c *UserRepo) CountUser(where *biz.UserWhere) (total int64, err error) {
	conditions, err := c.ParseUser(where) //条件转换
	if err != nil {
		return 0, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("user").Model(&biz.User{})
	if where != nil {
		if or, ok := conditions["or"]; ok && or != nil { // or条件
			if mapValues, isMap := or.(map[string]interface{}); isMap {
				for k, v := range mapValues {
					session = session.Or(k, v)
				}
			}
			delete(conditions, "or")
		}
		if and, ok := conditions["and"]; ok && and != nil { // and条件
			if mapValues, isMap := and.([]map[string]interface{}); isMap {
				for _, v := range mapValues {
					for k1, v1 := range v {
						session = session.Where(k1, v1)
					}
				}
			}
			delete(conditions, "and")
		}
		for key, value := range conditions {
			session.Where(key, value)
		}
	}
	err = session.Count(&total).Error
	if err != nil {
		log.Error("user table query data count error (%v)", err)
		return 0, err
	}
	return total, nil
}

// QueryUser 用户查询
func (c *UserRepo) QueryUser(where *biz.UserWhere, output *biz.UserOutput, field string) (list []*biz.User, err error) {
	conditions, err := c.ParseUser(where) //条件转换
	if err != nil {
		return nil, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("user").Model(&biz.User{})
	if where != nil {
		if or, ok := conditions["or"]; ok && or != nil { // or条件
			if mapValues, isMap := or.(map[string]interface{}); isMap {
				for k, v := range mapValues {
					session = session.Or(k, v)
				}
			}
			delete(conditions, "or")
		}
		if and, ok := conditions["and"]; ok && and != nil { // and条件
			if mapValues, isMap := and.([]map[string]interface{}); isMap {
				for _, v := range mapValues {
					for k1, v1 := range v {
						session = session.Where(k1, v1)
					}
				}
			}
			delete(conditions, "and")
		}
		for key, value := range conditions {
			session.Where(key, value)
		}
	}
	if output.PageSize != 0 && output.PageNum != 0 {
		session.Limit(output.PageSize).Offset((output.PageNum - 1) * output.PageSize)
	}
	if output.Order != "" {
		session.Order(output.Order)
	}
	err = session.Select(field).Find(&list).Error
	if err != nil {
		log.Error("user table query data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// CreateUser 用户创建
func (c *UserRepo) CreateUser(create []*biz.User) (list []*biz.User, err error) {
	for _, v := range create {
		v.CreateTime = time.Now()
		v.UpdateTime = time.Now()
		v.DeleteTime = timeParse("0001-01-01T00:00:00+08:00")
		// 生成盐
		salt, err := utils.GenerateSalt(16)
		if err != nil {
			log.Error("user table create salt error(%v)", err)
			return nil, err
		}
		v.Salt = salt
		v.Password = utils.SHA256(fmt.Sprintf("%v%v", utils.SHA256(v.Password), salt))
		err = c.db.GormDB.Table("user").Create(v).Error
		if err != nil {
			log.Error("user table create data error(%v)", err)
			return nil, err
		}
		list = append(list, &biz.User{
			ID: v.ID,
		})
	}
	return list, nil
}

// UpdateUser 用户更新
func (c *UserRepo) UpdateUser(where *biz.UserWhere, update *biz.User) error {
	updateMap := make(map[string]interface{})
	if update != nil {
		updateMap["name"] = update.Name
		updateMap["role"] = update.Role
		updateMap["is_delete"] = update.IsDelete
	}
	if update.Password != "" {
		// 生成盐
		salt, err := utils.GenerateSalt(16)
		if err != nil {
			log.Error("user table create salt error(%v)", err)
			return err
		}
		updateMap["salt"] = salt
		updateMap["password"] = utils.SHA256(fmt.Sprintf("%v%v", utils.SHA256(update.Password), salt))
	}
	if update.IsDelete == 1 {
		updateMap["delete_time"] = time.Now()
	}
	if where != nil && len(updateMap) > 0 {
		err := c.db.GormDB.Table("user").Where(where.Query, where.Arg).Updates(updateMap).Error
		if err != nil {
			log.Error("user table update data error(%v)", err)
			return err
		}
	}
	return nil
}

// DeleteUser 用户删除
func (c *UserRepo) DeleteUser(softDelete int, deleteID []int) error {
	if softDelete == 1 { // 软删除
		now := time.Now()
		deleteTime := now.Format("2006-01-02 15:04:05")
		var updateMap map[string]interface{}
		updateMap["is_delete"] = 1
		updateMap["delete_time"] = deleteTime
		err := c.db.GormDB.Table("user").Where("id IN ?", deleteID).Updates(updateMap).Error
		if err != nil {
			log.Error("user table delete data error(%v)", err)
			return nil
		}
	} else {
		err := c.db.GormDB.Table("user").Where("id IN ?", deleteID).Delete(&biz.User{}).Error
		if err != nil {
			log.Error("user table update data error(%v)", err)
			return err
		}
	}
	return nil
}

// UpsertUser 用户更新或创建
func (c *UserRepo) UpsertUser(upsert []*biz.UserUpsert) error {
	return nil
}
