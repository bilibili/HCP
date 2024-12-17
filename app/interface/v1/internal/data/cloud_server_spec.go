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

// CloudServerSpecRepo 云服务器规格DB
type CloudServerSpecRepo struct {
	db *biz.DB
}

// NewCloudServerSpecRepo 初始化云服务器规格DB
func NewCloudServerSpecRepo(db *biz.DB) biz.CloudServerSpecRepo {
	return &CloudServerSpecRepo{db: db}
}

// ParseCloudServerSpec 云服务器规格查询条件处理
func (c *CloudServerSpecRepo) ParseCloudServerSpec(where *biz.CloudServerSpecWhere) (conditions map[string]interface{}, err error) {
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
		log.Error("cloud_server_spec table parse data error(%v)", err)
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
			log.Error("cloud_server_spec table parse data error(%v)", err)
			return nil, err
		}
		conditions = biz.GetHandleConditions(conditionDos, "cloud_server_spec")
		for k, v := range conditionMaps {
			conditions[k] = v
		}
	}
	return conditions, nil
}

// CountCloudServerSpec 云服务器规格查询数量
func (c *CloudServerSpecRepo) CountCloudServerSpec(where *biz.CloudServerSpecWhere) (total int64, err error) {
	conditions, err := c.ParseCloudServerSpec(where) //条件转换
	if err != nil {
		return 0, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("cloud_server_spec").Model(&biz.CloudServerSpec{})
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
		log.Error("cloud_server_spec table query data count error (%v)", err)
		return 0, err
	}
	return total, nil
}

// QueryCloudServerSpec 云服务器规格查询
func (c *CloudServerSpecRepo) QueryCloudServerSpec(where *biz.CloudServerSpecWhere, output *biz.CloudServerSpecOutput) (list []*biz.CloudServerSpec, err error) {
	conditions, err := c.ParseCloudServerSpec(where) //条件转换
	if err != nil {
		return nil, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("cloud_server_spec").Model(&biz.CloudServerSpec{})
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
	err = session.Joins("Account", c.db.GormDB.Select("id", "name", "cid", "provider_id")).
		Joins("Provider", c.db.GormDB.Select("id", "name", "alias")).
		Find(&list).Error
	if err != nil {
		log.Error("cloud_server_spec table query data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// CreateCloudServerSpec 云服务器规格创建
func (c *CloudServerSpecRepo) CreateCloudServerSpec(create []*biz.CloudServerSpec) (list []*biz.CloudServerSpec, err error) {
	for _, v := range create {
		err = c.db.GormDB.Table("cloud_server_spec").Create(v).Error
		if err != nil {
			log.Error("cloud_server_spec table create data error(%v)", err)
			return nil, err
		}
		list = append(list, &biz.CloudServerSpec{
			CloudProductCommon: biz.CloudProductCommon{
				ID: v.ID,
			},
		})
	}
	return list, nil
}

// UpdateCloudServerSpec 云服务器规格更新
func (c *CloudServerSpecRepo) UpdateCloudServerSpec(where *biz.CloudServerSpecWhere, update *biz.CloudServerSpec) error {
	updateMap := make(map[string]interface{})
	if update != nil {
		updateMap["name"] = update.Name
		updateMap["cid"] = update.CID
		updateMap["account_id"] = update.AccountID
		updateMap["provider_id"] = update.ProviderID
		updateMap["bandwidth"] = update.Bandwidth
		updateMap["category"] = update.Category
		updateMap["cpu"] = update.CPU
		updateMap["family"] = update.Family
		updateMap["gpu"] = update.GPU
		updateMap["gpu_model"] = update.GPUModel
		updateMap["memory"] = update.Memory
		updateMap["pps"] = update.PPS
		updateMap["status"] = update.Status
	}
	if where != nil && len(updateMap) > 0 {
		err := c.db.GormDB.Table("cloud_server_spec").Where(where.Query, where.Arg).Updates(updateMap).Error
		if err != nil {
			log.Error("cloud_server_spec table update data error(%v)", err)
			return err
		}
	}
	return nil
}

// DeleteCloudServerSpec 云服务器规格删除
func (c *CloudServerSpecRepo) DeleteCloudServerSpec(deleteID []int) error {
	err := c.db.GormDB.Table("cloud_server_spec").Where("id IN ?", deleteID).Delete(&biz.CloudServerSpec{}).Error
	if err != nil {
		log.Error("cloud_server_spec table delete data error(%v)", err)
		return err
	}
	return nil
}

// UpsertCloudServerSpec 云服务器规格更新或创建
func (c *CloudServerSpecRepo) UpsertCloudServerSpec(upsert []*biz.CloudServerSpecUpsert) error {
	return nil
}
