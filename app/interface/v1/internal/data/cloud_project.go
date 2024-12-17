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

	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/go-kratos/kratos/pkg/log"
	"gorm.io/gorm"
)

// CloudProjectRepo 云项目DB
type CloudProjectRepo struct {
	db *biz.DB
}

// NewCloudProjectRepo 初始化云项目DB
func NewCloudProjectRepo(db *biz.DB) biz.CloudProjectRepo {
	return &CloudProjectRepo{db: db}
}

// ParseCloudProject 云项目查询条件处理
func (c *CloudProjectRepo) ParseCloudProject(where *biz.CloudProjectWhere) (conditions map[string]interface{}, err error) {
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
	if provider, ok := where.Conditions["provider"]; ok && provider != nil {
		if providerMap, okMap := provider.(map[string]interface{}); okMap {
			for key, value := range providerMap {
				if mapValues, ok := value.(map[string]interface{}); !ok {
					conditionMaps[fmt.Sprintf("Provider.%v", key)] = value
				} else {
					for subKey, subValue := range mapValues {
						conditionMaps[fmt.Sprintf("Provider__%v.%v", biz.FirstStrUpper(key), subKey)] = subValue
					}
				}
			}
		}
		conditionSources["provider"] = true
	}
	jsonBytes, err := json.Marshal(conditionMaps)
	if err != nil {
		return nil, err
	}
	conditionMaps, err = biz.ParseCloudData(jsonBytes)
	if err != nil {
		log.Error("cloud_project table parse data error(%v)", err)
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
			log.Error("cloud_project table parse data error(%v)", err)
			return nil, err
		}
		conditions = biz.GetHandleConditions(conditionDos, "cloud_project")
		for k, v := range conditionMaps {
			conditions[k] = v
		}
	}
	return conditions, nil
}

// CountCloudProject 云项目查询数量
func (c *CloudProjectRepo) CountCloudProject(where *biz.CloudProjectWhere) (total int64, err error) {
	conditions, err := c.ParseCloudProject(where) //条件转换
	if err != nil {
		return 0, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("cloud_project").Model(&biz.CloudProject{})
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
	err = session.Joins("Account").Joins("Provider").Count(&total).Error
	if err != nil {
		log.Error("cloud_project table query data count error (%v)", err)
		return 0, err
	}
	return total, nil
}

// QueryCloudProject 云项目查询
func (c *CloudProjectRepo) QueryCloudProject(where *biz.CloudProjectWhere, output *biz.CloudProjectOutput) (list []*biz.CloudProject, err error) {
	conditions, err := c.ParseCloudProject(where) //条件转换
	if err != nil {
		return nil, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("cloud_project").Model(&biz.CloudProject{})
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
	err = session.Joins("Provider", c.db.GormDB.Select("id", "name")).
		Joins("Account", c.db.GormDB.Select("id", "name")).
		Find(&list).Error
	if err != nil {
		log.Error("cloud_project table query data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// CreateCloudProject 云项目创建
func (c *CloudProjectRepo) CreateCloudProject(create []*biz.CloudProject) (list []*biz.CloudProject, err error) {
	for _, v := range create {
		err = c.db.GormDB.Table("cloud_project").Create(v).Error
		if err != nil {
			log.Error("cloud_project table create data error(%v)", err)
			//return nil, err //tips：可能会存在重复的标签名拼接的项目会和数据库唯一组合索引冲突，有重复只生成错误日志不返回错误
		}
		list = append(list, &biz.CloudProject{
			CloudProductCommon: biz.CloudProductCommon{
				ID: v.ID,
			},
		})
	}
	return list, nil
}

// UpdateCloudProject 云项目更新
func (c *CloudProjectRepo) UpdateCloudProject(where *biz.CloudProjectWhere, update *biz.CloudProject) error {
	updateMap := make(map[string]interface{})
	if update != nil {
		updateMap["name"] = update.Name
		updateMap["cid"] = update.CID
		updateMap["account_id"] = update.AccountID
		updateMap["provider_id"] = update.ProviderID
	}
	if where != nil && len(updateMap) > 0 {
		err := c.db.GormDB.Table("cloud_project").Where(where.Query, where.Arg).Updates(updateMap).Error
		if err != nil {
			log.Error("cloud_project table update data error(%v)", err)
			return err
		}
	}
	return nil
}

// DeleteCloudProject 云项目删除
func (c *CloudProjectRepo) DeleteCloudProject(deleteID []int) error {
	err := c.db.GormDB.Table("cloud_project").Where("id IN ?", deleteID).Delete(&biz.CloudProject{}).Error
	if err != nil {
		log.Error("cloud_project table update data error(%v)", err)
		return err
	}
	return nil
}

// UpsertCloudProject 云项目更新或创建
func (c *CloudProjectRepo) UpsertCloudProject(upsert []*biz.CloudProjectUpsert) error {
	return nil
}
