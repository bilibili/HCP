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

// DiskTypeRepo 数据盘类型DB
type DiskTypeRepo struct {
	db *biz.DB
}

// NewDiskTypeRepo 初始化数据盘类型DB
func NewDiskTypeRepo(db *biz.DB) biz.DiskTypeRepo {
	return &DiskTypeRepo{db: db}
}

// ParseDiskType 数据盘类型查询条件处理
func (c *DiskTypeRepo) ParseDiskType(where *biz.DiskTypeWhere) (conditions map[string]interface{}, err error) {
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
		log.Error("disk_type table parse data error(%v)", err)
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
			log.Error("disk_type table parse data error(%v)", err)
			return nil, err
		}
		conditions = biz.GetHandleConditions(conditionDos, "disk_type")
		for k, v := range conditionMaps {
			conditions[k] = v
		}
	}
	return conditions, nil
}

// CountDiskType 数据盘类型查询数量
func (c *DiskTypeRepo) CountDiskType(where *biz.DiskTypeWhere) (total int64, err error) {
	conditions, err := c.ParseDiskType(where) //条件转换
	if err != nil {
		return 0, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("disk_type").Model(&biz.DiskType{})
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
		log.Error("disk_type table query data count error (%v)", err)
		return 0, err
	}
	return total, nil
}

// QueryDiskType 数据盘类型查询
func (c *DiskTypeRepo) QueryDiskType(where *biz.DiskTypeWhere, output *biz.DiskTypeOutput) (list []*biz.DiskType, err error) {
	conditions, err := c.ParseDiskType(where) //条件转换
	if err != nil {
		return nil, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("disk_type").Model(&biz.DiskType{})
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
		log.Error("disk_type table query data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// CreateDiskType 数据盘类型创建
func (c *DiskTypeRepo) CreateDiskType(create []*biz.DiskType) (list []*biz.DiskType, err error) {
	for _, v := range create {
		err = c.db.GormDB.Table("disk_type").Create(v).Error
		if err != nil {
			log.Error("disk_type table create data error(%v)", err)
			return nil, err
		}
		list = append(list, &biz.DiskType{
			ID: v.ID,
		})
	}
	return list, nil
}

// UpdateDiskType 数据盘类型更新
func (c *DiskTypeRepo) UpdateDiskType(where *biz.DiskTypeWhere, update *biz.DiskType) error {
	updateMap := make(map[string]interface{})
	if update != nil {
		updateMap["name"] = update.Name
		updateMap["cid"] = update.CID
		updateMap["provider_id"] = update.ProviderID
		updateMap["product_type"] = update.ProductType
	}
	if where != nil && len(updateMap) > 0 {
		err := c.db.GormDB.Table("disk_type").Where(where.Query, where.Arg).Updates(updateMap).Error
		if err != nil {
			log.Error("disk_type table update data error(%v)", err)
			return err
		}
	}
	return nil
}

// DeleteDiskType 数据盘类型删除
func (c *DiskTypeRepo) DeleteDiskType(deleteID []int) error {
	err := c.db.GormDB.Table("disk_type").Where("id IN ?", deleteID).Delete(&biz.DiskType{}).Error
	if err != nil {
		log.Error("disk_type table update data error(%v)", err)
		return err
	}
	return nil
}

// UpsertDiskType 数据盘类型更新或创建
func (c *DiskTypeRepo) UpsertDiskType(upsert []*biz.DiskTypeUpsert) error {
	return nil
}
