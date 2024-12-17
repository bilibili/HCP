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

// ProviderRepo 厂商DB
type ProviderRepo struct {
	db *biz.DB
}

// NewProviderRepo 初始化厂商DB
func NewProviderRepo(db *biz.DB) biz.ProviderRepo {
	return &ProviderRepo{
		db: db,
	}
}

// ParseProvider 厂商查询条件处理
func (c *ProviderRepo) ParseProvider(where *biz.ProviderWhere) (conditions map[string]interface{}, err error) {
	conditionMaps, conditionSources := make(map[string]interface{}), make(map[string]interface{})
	if projectConfig, ok := where.Conditions["project_config"]; ok && projectConfig != nil {
		jsonBytes, err := json.Marshal(projectConfig)
		if err != nil {
			return nil, err
		}
		projectConfigConditions, err := biz.ParseCloudData(jsonBytes)
		if err != nil {
			log.Error("provider project_config table parse data error(%v)", err)
			return nil, err
		}
		// 查询云厂商ID
		var providerIds []int
		var projectAccountConfigs []biz.ProjectAccountConfig
		projectAccountConfigSession := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("project_account_config").Model(&biz.ProjectAccountConfig{})
		if where != nil {
			for key, value := range projectConfigConditions {
				field := fmt.Sprintf("project_config_%v", key)
				projectAccountConfigSession.Where(field, value)
			}
		}
		err = projectAccountConfigSession.Select("Account.provider_id").Distinct().Joins("Account", c.db.GormDB.Select("id", "provider_id")).
			Group("account_id").Find(&projectAccountConfigs).Error
		for _, v := range projectAccountConfigs {
			providerIds = append(providerIds, v.Account.ProviderID)
		}
		if len(providerIds) > 0 {
			conditionMaps["id_IN ?"] = providerIds
		}
		conditionSources["project_config"] = true
	}
	jsonBytes, err := json.Marshal(conditionMaps)
	if err != nil {
		return nil, err
	}
	conditionMaps, err = biz.ParseCloudData(jsonBytes)
	if err != nil {
		log.Error("provider table parse data error(%v)", err)
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
			log.Error("provider table parse data error(%v)", err)
			return nil, err
		}
		conditions = biz.GetHandleConditions(conditionDos, "")
		for k, v := range conditionMaps {
			conditions[k] = v
		}
	}
	return conditions, nil
}

// CountProvider 厂商查询数量
func (c *ProviderRepo) CountProvider(where *biz.ProviderWhere) (total int64, err error) {
	conditions, err := c.ParseProvider(where) //条件转换
	if err != nil {
		return 0, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("provider").Model(&biz.Provider{})
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
		log.Error("provider table query data count error (%v)", err)
		return 0, err
	}
	return total, nil
}

// QueryProvider 厂商查询
func (c *ProviderRepo) QueryProvider(where *biz.ProviderWhere, output *biz.ProviderOutput) (list []*biz.Provider, err error) {
	conditions, err := c.ParseProvider(where) //条件转换
	if err != nil {
		return nil, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("provider").Model(&biz.Provider{})
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
	for _, v := range list {
		var total int64
		err = c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("account").Model(&biz.Account{}).
			Where("provider_id = ?", v.ID).
			Count(&total).Error
		if err != nil {
			log.Error("provider account table query data count error (%v)", err)
			return nil, err
		}
		v.AccountNumber = int(total)
	}
	if err != nil {
		log.Error("provider table query data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// CreateProvider 厂商创建
func (c *ProviderRepo) CreateProvider(create []*biz.Provider) (list []*biz.Provider, err error) {
	for _, v := range create {
		err = c.db.GormDB.Table("provider").Create(v).Error
		if err != nil {
			log.Error("provider table create data error(%v)", err)
			return nil, err
		}
		list = append(list, &biz.Provider{
			ID: v.ID,
		})
	}
	return list, nil
}

func (c *ProviderRepo) UpdateProvider(where *biz.ProviderWhere, update *biz.Provider) error {
	updateMap := make(map[string]interface{})
	if update != nil {
		updateMap["name"] = update.Name
		updateMap["alias"] = update.Alias
		updateMap["pre_sale_contacts"] = update.PreSaleContacts
		updateMap["after_sale_contacts"] = update.AfterSaleContacts
		updateMap["business_contacts"] = update.BusinessContacts
	}
	if isDelete, ok := where.Conditions["is_delete"]; ok && isDelete == 1 {
		updateMap["delete_time"] = time.Now()
	}
	if where != nil && len(updateMap) > 0 {
		err := c.db.GormDB.Table("provider").Where(where.Query, where.Arg).Updates(updateMap).Error
		if err != nil {
			log.Error("provider table update data error(%v)", err)
			return err
		}
	}
	return nil
}

// DeleteProvider 厂商删除
func (c *ProviderRepo) DeleteProvider(deleteID []int) error {
	err := c.db.GormDB.Table("provider").Where("id IN ?", deleteID).Delete(&biz.Provider{}).Error
	if err != nil {
		log.Error("provider table update data error(%v)", err)
		return err
	}
	return nil
}

// UpsertProvider 厂商更新或创建
func (c *ProviderRepo) UpsertProvider(upsert []*biz.ProviderUpsert) error {
	return nil
}
