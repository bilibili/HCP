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
	"github.com/go-kratos/kratos/pkg/log"
	"gorm.io/gorm"
)

// ProjectAccountConfigRepo 云账号和本地项目关联关系DB
type ProjectAccountConfigRepo struct {
	db *biz.DB
}

// NewProjectAccountConfigRepo 初始化云账号和本地项目关联关系DB
func NewProjectAccountConfigRepo(db *biz.DB) biz.ProjectAccountConfigRepo {
	return &ProjectAccountConfigRepo{db: db}
}

// ParseProjectAccountConfig 本地项目账号关联关系查询条件处理
func (c *ProjectAccountConfigRepo) ParseProjectAccountConfig(where *biz.ProjectAccountConfigWhere) (conditions map[string]interface{}, err error) {
	conditionMaps, conditionSources := make(map[string]interface{}), make(map[string]interface{})
	if account, ok := where.Conditions["account"]; ok && account != nil {
		if accountMap, okMap := account.(map[string]interface{}); okMap {
			for key, value := range accountMap {
				if mapValues, ok := value.(map[string]interface{}); !ok {
					conditionMaps[fmt.Sprintf("Account.%v", key)] = value
				} else {
					for subKey, subValue := range mapValues {
						conditionMaps[fmt.Sprintf("Account__%v.%v", biz.FirstStrUpper(key), subKey)] = subValue
					}
				}
			}
		}
		conditionSources["account"] = true
	}
	if cloudProject, ok := where.Conditions["project"]; ok && cloudProject != nil {
		if projectMap, okMap := cloudProject.(map[string]interface{}); okMap {
			for key, value := range projectMap {
				if mapValues, ok := value.(map[string]interface{}); !ok {
					conditionMaps[fmt.Sprintf("Project.%v", key)] = value
				} else {
					for subKey, subValue := range mapValues {
						conditionMaps[fmt.Sprintf("Project__%v.%v", biz.FirstStrUpper(key), subKey)] = subValue
					}
				}
			}
		}
		conditionSources["project"] = true
	}
	if projectConfig, ok := where.Conditions["project_config"]; ok && projectConfig != nil {
		if projectConfigMap, okMap := projectConfig.(map[string]interface{}); okMap {
			for key, value := range projectConfigMap {
				if mapValues, ok := value.(map[string]interface{}); !ok {
					conditionMaps[fmt.Sprintf("ProjectConfig.%v", key)] = value
				} else {
					for subKey, subValue := range mapValues {
						conditionMaps[fmt.Sprintf("ProjectConfig__%v.%v", biz.FirstStrUpper(key), subKey)] = subValue
					}
				}
			}
		}
		conditionSources["project_config"] = true
	}
	jsonBytes, err := json.Marshal(conditionMaps)
	if err != nil {
		return nil, err
	}
	conditionMaps, err = biz.ParseCloudData(jsonBytes)
	if err != nil {
		log.Error("project_account_config table parse data error(%v)", err)
		return nil, err
	}
	conditionDiff := biz.DifferenceMap(where.Conditions, conditionSources)
	if conditionDiff != nil {
		jsonByteDiffs, err := json.Marshal(conditionDiff)
		if err != nil {
			return nil, err
		}
		conditionDos, err := biz.ParseCloudData(jsonByteDiffs)
		if err != nil {
			log.Error("project_account_config table parse data error(%v)", err)
			return nil, err
		}
		conditions = biz.GetHandleConditions(conditionDos, "project_account_config")
		for k, v := range conditionMaps {
			conditions[k] = v
		}
	}
	return conditions, nil
}

// CountProjectAccountConfig 本地项目账号关联关系查询数量
func (c *ProjectAccountConfigRepo) CountProjectAccountConfig(where *biz.ProjectAccountConfigWhere) (total int64, err error) {
	conditions, err := c.ParseProjectAccountConfig(where) //条件转换
	if err != nil {
		return 0, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("project_account_config").Model(&biz.ProjectAccountConfig{})
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
	err = session.Joins("Account").
		Joins("Project").
		Joins("ProjectConfig").
		Joins("LEFT JOIN provider Account_Provider ON Account.provider_id = Account_Provider.id").
		Count(&total).Error
	if err != nil {
		log.Error("project_account_config table query data count error (%v)", err)
		return 0, err
	}
	return total, nil
}

// QueryProjectAccountConfig 本地项目账号关联关系查询
func (c *ProjectAccountConfigRepo) QueryProjectAccountConfig(where *biz.ProjectAccountConfigWhere, output *biz.ProjectAccountConfigOutput) (list []*biz.ProjectAccountConfig, err error) {
	conditions, err := c.ParseProjectAccountConfig(where) //条件转换
	if err != nil {
		return nil, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("project_account_config").Model(&biz.ProjectAccountConfig{})
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
	err = session.Joins("Account", c.db.GormDB.Select("id", "name", "alias", "cid", "provider_id")).
		Joins("Project", c.db.GormDB.Select("id", "name", "cid")).
		Joins("ProjectConfig", c.db.GormDB.Select("id", "name", "alias")).
		Joins("LEFT JOIN provider Account_Provider ON Account.provider_id = Account_Provider.id").
		Preload("Account.Provider").
		Find(&list).Error
	if err != nil {
		log.Error("project_account_config table query data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// CreateProjectAccountConfig 本地项目账号关联关系创建
func (c *ProjectAccountConfigRepo) CreateProjectAccountConfig(create []*biz.ProjectAccountConfig) (list []*biz.ProjectAccountConfig, err error) {
	for _, v := range create {
		v.Account = nil
		v.Project = nil
		v.ProjectConfig = nil
		v.CreateTime = time.Now()
		v.UpdateTime = time.Now()
		v.DeleteTime = timeParse("0001-01-01T00:00:00+08:00")
		err = c.db.GormDB.Table("project_account_config").Create(v).Error
		if err != nil {
			log.Error("project_account_config table create data error(%v)", err)
			return nil, err
		}
		list = append(list, &biz.ProjectAccountConfig{
			ID: v.ID,
		})
	}
	return list, nil
}

// UpdateProjectAccountConfig 本地项目账号关联关系更新
func (c *ProjectAccountConfigRepo) UpdateProjectAccountConfig(where *biz.ProjectAccountConfigWhere, update *biz.ProjectAccountConfig) error {
	updateMap := make(map[string]interface{})
	if update != nil {
		updateMap["project_id"] = update.ProjectID
		updateMap["is_delete"] = update.IsDelete
	}
	if update.IsDelete == 1 {
		updateMap["delete_time"] = time.Now()
	}
	if where != nil && len(updateMap) > 0 {
		err := c.db.GormDB.Table("project_account_config").Where(where.Query, where.Arg).Updates(updateMap).Error
		if err != nil {
			log.Error("project_account_config table update data error(%v)", err)
			return err
		}
	}
	return nil
}

// DeleteProjectAccountConfig 本地项目账号关联关系删除
func (c *ProjectAccountConfigRepo) DeleteProjectAccountConfig(deleteID []int) error {
	err := c.db.GormDB.Table("project_account_config").Where("id IN ?", deleteID).Delete(&biz.ProjectAccountConfig{}).Error
	if err != nil {
		log.Error("project_account_config table update data error(%v)", err)
		return err
	}
	return nil
}

// UpsertProjectAccountConfig 本地项目账号关联关系更新或创建
func (c *ProjectAccountConfigRepo) UpsertProjectAccountConfig(upsert []*biz.ProjectAccountConfigUpsert) error {
	return nil
}
