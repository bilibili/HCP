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

// CloudRegionRepo 云地域DB
type CloudRegionRepo struct {
	db *biz.DB
}

// NewCloudRegionRepo 初始化云地域DB
func NewCloudRegionRepo(db *biz.DB) biz.CloudRegionRepo {
	return &CloudRegionRepo{
		db: db,
	}
}

// ParseCloudRegion 云地域查询条件处理
func (c *CloudRegionRepo) ParseCloudRegion(where *biz.CloudRegionWhere) (conditions map[string]interface{}, err error) {
	conditionMaps, conditionSources := make(map[string]interface{}), make(map[string]interface{})
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
	if account, ok := where.Conditions["account"]; ok && account != nil {
		jsonBytes, err := json.Marshal(account)
		if err != nil {
			return nil, err
		}
		accountConditions, err := biz.ParseCloudData(jsonBytes)
		if err != nil {
			log.Error("cloud_region account table parse data error(%v)", err)
			return nil, err
		}
		//查询地域ID
		var regionIds []int
		var regionAssociation []*biz.CloudRegionAssociation
		regionSession := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("cloud_region_association").Model(&biz.CloudRegionAssociation{})
		if where != nil {
			for key, value := range accountConditions {
				field := fmt.Sprintf("account_%v", key)
				regionSession.Where(field, value)
			}
		}
		err = regionSession.Find(&regionAssociation).Error
		if err != nil {
			log.Error("cloud_region cloud_region_association table query data error(%v)", err)
			return nil, err
		}
		for _, v := range regionAssociation {
			regionIds = append(regionIds, v.RegionID)
		}
		if len(regionIds) > 0 {
			conditionMaps["cloud_region.id_IN ?"] = regionIds
		}
		conditionSources["account"] = true
	}
	jsonBytes, err := json.Marshal(conditionMaps)
	if err != nil {
		return nil, err
	}
	conditionMaps, err = biz.ParseCloudData(jsonBytes)
	if err != nil {
		log.Error("cloud_region table parse data error(%v)", err)
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
		conditions = biz.GetHandleConditions(conditionDos, "cloud_region")
		for k, v := range conditionMaps {
			conditions[k] = v
		}
	}
	return conditions, nil
}

// CountCloudRegion 云地域查询数量
func (c *CloudRegionRepo) CountCloudRegion(where *biz.CloudRegionWhere) (total int64, err error) {
	conditions, err := c.ParseCloudRegion(where) //条件转换
	if err != nil {
		return 0, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("cloud_region").Model(&biz.CloudRegion{})
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
		log.Error("cloud_region table query data count error (%v)", err)
		return 0, err
	}
	return total, nil
}

// QueryCloudRegion 云地域查询
func (c *CloudRegionRepo) QueryCloudRegion(where *biz.CloudRegionWhere, output *biz.CloudRegionOutput) (list []*biz.CloudRegion, err error) {
	conditions, err := c.ParseCloudRegion(where) //条件转换
	if err != nil {
		return nil, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("cloud_region").Model(&biz.CloudRegion{})
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
	err = session.Joins("Provider", c.db.GormDB.Select("id", "name")).Find(&list).Error
	if err != nil {
		log.Error("cloud_region table query data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// CreateCloudRegion 云地域创建
func (c *CloudRegionRepo) CreateCloudRegion(create []*biz.CloudRegion) (list []*biz.CloudRegion, err error) {
	for _, v := range create {
		err = c.db.GormDB.Table("cloud_region").Create(v).Error
		if err != nil {
			log.Error("cloud_region table create data error(%v)", err)
			return nil, err
		}
		list = append(list, &biz.CloudRegion{
			CloudProductCommon: biz.CloudProductCommon{
				ID: v.ID,
			},
		})
	}
	return list, nil
}

// UpdateCloudRegion 云地域更新
func (c *CloudRegionRepo) UpdateCloudRegion(where *biz.CloudRegionWhere, update *biz.CloudRegion) error {
	updateMap := make(map[string]interface{})
	if update != nil {
		updateMap["name"] = update.CloudProductCommon.Name
		updateMap["cid"] = update.CloudProductCommon.CID
		updateMap["provider_id"] = update.ProviderID
	}
	if where != nil && len(updateMap) > 0 {
		err := c.db.GormDB.Table("cloud_region").Where(where.Query, where.Arg).Updates(updateMap).Error
		if err != nil {
			log.Error("cloud_region table update data error(%v)", err)
			return err
		}
	}
	return nil
}

// DeleteCloudRegion 云地域删除
func (c *CloudRegionRepo) DeleteCloudRegion(deleteID []int) error {
	err := c.db.GormDB.Table("cloud_region").Where("id IN ?", deleteID).Delete(&biz.CloudRegion{}).Error
	if err != nil {
		log.Error("cloud_region table update data error(%v)", err)
		return err
	}
	return nil
}

// UpsertCloudRegion 云地域更新或创建
func (c *CloudRegionRepo) UpsertCloudRegion(upsert []*biz.CloudRegionUpsert) error {
	return nil
}
