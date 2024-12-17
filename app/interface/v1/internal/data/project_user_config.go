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
	"time"

	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/go-kratos/kratos/pkg/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProjectUserConfigRepo 本地项目和用户关联关系DB
type ProjectUserConfigRepo struct {
	db *biz.DB
}

// NewProjectUserConfigRepo 实例化项目用户关联关系DB
func NewProjectUserConfigRepo(db *biz.DB) biz.ProjectUserConfigRepo {
	return &ProjectUserConfigRepo{db: db}
}

// ParseProjectUserConfig 项目用户关联关系查询条件处理
func (c *ProjectUserConfigRepo) ParseProjectUserConfig(where *biz.ProjectUserConfigWhere) (conditions map[string]interface{}, err error) {
	if where.Conditions != nil {
		jsonBytes, err := json.Marshal(where.Conditions)
		if err != nil {
			return nil, err
		}
		conditionDos, err := biz.ParseCloudData(jsonBytes)
		if err != nil {
			log.Error("project_user_config table parse data error(%v)", err)
			return nil, err
		}
		conditions = biz.GetHandleConditions(conditionDos, "")
	}
	return conditions, nil
}

// CountProjectUserConfig 项目用户关联关系查询数量
func (c *ProjectUserConfigRepo) CountProjectUserConfig(where *biz.ProjectUserConfigWhere, group string) (total int64, err error) {
	conditions, err := c.ParseProjectUserConfig(where) //条件转换
	if err != nil {
		return 0, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("project_user_config").Model(&biz.ProjectUserConfig{})
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
	if group != "" {
		session.Group(group)
	}
	err = session.Count(&total).Error
	if err != nil {
		log.Error("project_user_config table query data count error (%v)", err)
		return 0, err
	}
	return total, nil
}

// QueryProjectUserConfig 项目用户关联关系查询
func (c *ProjectUserConfigRepo) QueryProjectUserConfig(where *biz.ProjectUserConfigWhere, output *biz.ProjectUserConfigOutput) (list []*biz.ProjectUserConfig, err error) {
	conditions, err := c.ParseProjectUserConfig(where) //条件转换
	if err != nil {
		return nil, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("project_user_config").Model(&biz.ProjectUserConfig{})
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
	err = session.Find(&list).Error
	if err != nil {
		log.Error("project_user_config table query data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// CreateProjectUserConfig 项目用户关联关系创建
func (c *ProjectUserConfigRepo) CreateProjectUserConfig(create []*biz.ProjectUserConfig) (list []*biz.ProjectUserConfig, err error) {
	for _, v := range create {
		v.CreateTime = time.Now()
		v.UpdateTime = time.Now()
		v.DeleteTime = timeParse("0001-01-01T00:00:00+08:00")
		v.User = nil
		v.ProjectConfig = nil
		err = c.db.GormDB.Table("project_user_config").Create(v).Error
		if err != nil {
			log.Error("project_user_config table create data error(%v)", err)
			return nil, err
		}
		list = append(list, &biz.ProjectUserConfig{
			ID: v.ID,
		})
	}
	return list, nil
}

// UpdateProjectUserConfig 项目用户关联关系更新
func (c *ProjectUserConfigRepo) UpdateProjectUserConfig(where *biz.ProjectUserConfigWhere, update *biz.ProjectUserConfig) error {
	updateMap := make(map[string]interface{})
	if update != nil {
		updateMap["user_id"] = update.UserID
		updateMap["project_config_id"] = update.ProjectConfigID
		updateMap["role"] = update.Role
		updateMap["is_delete"] = update.IsDelete
	}
	if update.IsDelete == 1 {
		updateMap["delete_time"] = time.Now()
	}
	if where != nil && len(updateMap) > 0 {
		err := c.db.GormDB.Table("project_user_config").Where(where.Query, where.Arg).Updates(updateMap).Error
		if err != nil {
			log.Error("project_user_config table update data error(%v)", err)
			return err
		}
	}
	return nil
}

// DeleteProjectUserConfig 项目用户关联关系删除
func (c *ProjectUserConfigRepo) DeleteProjectUserConfig(softDelete int, deleteID []int) error {
	if softDelete == 1 { // 软删除
		now := time.Now()
		deleteTime := now.Format("2006-01-02 15:04:05")
		var updateMap map[string]interface{}
		updateMap["is_delete"] = 1
		updateMap["delete_time"] = deleteTime
		err := c.db.GormDB.Table("project_user_config").Where("id IN ?", deleteID).Updates(updateMap).Error
		if err != nil {
			log.Error("project_user_config table delete data error(%v)", err)
			return nil
		}
	} else {
		err := c.db.GormDB.Table("project_user_config").Where("id IN ?", deleteID).Delete(&biz.ProjectUserConfig{}).Error
		if err != nil {
			log.Error("project_user_config table update data error(%v)", err)
			return err
		}
	}
	return nil
}

// DeleteProjectUserConfigByWhere 项目用户关联关系删除
func (c *ProjectUserConfigRepo) DeleteProjectUserConfigByWhere(where *biz.ProjectUserConfigWhere) error {
	if where != nil {
		session := c.db.GormDB.Table("project_user_config")
		for key, value := range where.Conditions {
			session.Where(key, value)
		}
		err := session.Delete(&biz.ProjectUserConfig{}).Error
		if err != nil {
			log.Error("project_user_config table delete data by where error(%v)", err)
			return err
		}
	}
	return nil
}

// UpsertProjectUserConfig 项目用户关联关系更新或创建
func (c *ProjectUserConfigRepo) UpsertProjectUserConfig(upsert *biz.ProjectUserConfigUpsert) error {
	err := c.db.GormDB.Table("project_user_config").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "project_config_id"}, {Name: "role"}},
		DoUpdates: clause.AssignmentColumns([]string{"user_id", "project_config_id", "role"}),
	}).Create(&biz.ProjectUserConfig{
		UserID:          upsert.UserID,
		ProjectConfigID: upsert.ProjectConfigID,
		Role:            upsert.Role,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
		DeleteTime:      timeParse("0001-01-01T00:00:00+08:00"),
	}).Error
	if err != nil {
		log.Error("project_user_config table upsert data error(%v)", err)
		return err
	}
	return nil
}
