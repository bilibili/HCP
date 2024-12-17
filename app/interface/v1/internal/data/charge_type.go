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

// ChargeTypeRepo 计费方式DB
type ChargeTypeRepo struct {
	db *biz.DB
}

// NewChargeTypeRepo 初始化计费方式DB
func NewChargeTypeRepo(db *biz.DB) biz.ChargeTypeRepo {
	return &ChargeTypeRepo{db: db}
}

// ParseChargeType 计费方式查询条件处理
func (c *ChargeTypeRepo) ParseChargeType(where *biz.ChargeTypeWhere) (conditions map[string]interface{}, err error) {
	conditionMaps, conditionSources := make(map[string]interface{}), make(map[string]interface{})
	if account, ok := where.Conditions["provider"]; ok && account != nil {
		if accountMap, okMap := account.(map[string]interface{}); okMap {
			for key, value := range accountMap {
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
		log.Error("charge_type table parse data error(%v)", err)
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
			log.Error("charge_type table parse data error(%v)", err)
			return nil, err
		}
		conditions = biz.GetHandleConditions(conditionDos, "charge_type")
		for k, v := range conditionMaps {
			conditions[k] = v
		}
	}
	return conditions, nil
}

// CountChargeType 查询计费方式数量
func (c *ChargeTypeRepo) CountChargeType(where *biz.ChargeTypeWhere) (total int64, err error) {
	conditions, err := c.ParseChargeType(where) //条件转换
	if err != nil {
		return 0, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("charge_type").Model(&biz.ChargeType{})
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
	err = session.Joins("Provider").Count(&total).Error
	if err != nil {
		log.Error("charge_type table query data count error (%v)", err)
		return 0, err
	}
	return total, nil
}

// QueryChargeType 查询计费方式
func (c *ChargeTypeRepo) QueryChargeType(where *biz.ChargeTypeWhere, output *biz.ChargeTypeOutput) (list []*biz.ChargeType, err error) {
	conditions, err := c.ParseChargeType(where) //条件转换
	if err != nil {
		return nil, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("charge_type").Model(&biz.ChargeType{})
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
	err = session.Joins("Provider").Find(&list).Error
	if err != nil {
		log.Error("charge_type table query data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// CreateChargeType 创建计费方式
func (c *ChargeTypeRepo) CreateChargeType(create []*biz.ChargeType) (list []*biz.ChargeType, err error) {
	for _, v := range create {
		err = c.db.GormDB.Table("charge_type").Create(v).Error
		if err != nil {
			log.Error("charge_type table create data error(%v)", err)
			return nil, err
		}
		list = append(list, &biz.ChargeType{
			ID: v.ID,
		})
	}
	return list, nil
}

// UpdateChargeType 更新计费方式
func (c *ChargeTypeRepo) UpdateChargeType(where *biz.ChargeTypeWhere, update *biz.ChargeType) error {
	updateMap := make(map[string]interface{})
	if update != nil {
		updateMap["name"] = update.Name
		updateMap["cid"] = update.CID
		updateMap["type"] = update.Type
		updateMap["provider_id"] = update.ProviderID
		updateMap["product_type"] = update.ProductType
	}
	if where != nil && len(updateMap) > 0 {
		err := c.db.GormDB.Table("charge_type").Where(where.Query, where.Arg).Updates(updateMap).Error
		if err != nil {
			log.Error("charge_type table update data error(%v)", err)
			return err
		}
	}
	return nil
}

// DeleteChargeType 删除计费方式
func (c *ChargeTypeRepo) DeleteChargeType(deleteID []int) error {
	err := c.db.GormDB.Table("charge_type").Where("id IN ?", deleteID).Delete(&biz.ChargeType{}).Error
	if err != nil {
		log.Error("charge_type table update data error(%v)", err)
		return err
	}
	return nil
}

// UpsertChargeType 计费方式更新或创建
func (c *ChargeTypeRepo) UpsertChargeType(upsert []*biz.ChargeTypeUpsert) error {
	return nil
}
