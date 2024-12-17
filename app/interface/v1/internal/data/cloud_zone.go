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

// CloudZoneRepo 云可用区DB
type CloudZoneRepo struct {
	db *biz.DB
}

// NewCloudZoneRepo 初始化云可用区DB
func NewCloudZoneRepo(db *biz.DB) biz.CloudZoneRepo {
	return &CloudZoneRepo{db: db}
}

// ParseCloudZone 云可用区查询条件处理
func (c *CloudZoneRepo) ParseCloudZone(where *biz.CloudZoneWhere) (conditions map[string]interface{}, err error) {
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
	if cloudRegion, ok := where.Conditions["region"]; ok && cloudRegion != nil {
		if regionMap, okMap := cloudRegion.(map[string]interface{}); okMap {
			for key, value := range regionMap {
				if mapValues, ok := value.(map[string]interface{}); !ok {
					conditionMaps[fmt.Sprintf("Region.%v", key)] = value
				} else {
					for subKey, subValue := range mapValues {
						conditionMaps[fmt.Sprintf("Region__%v.%v", biz.FirstStrUpper(key), subKey)] = subValue
					}
				}
			}
		}
		conditionSources["region"] = true
	}
	jsonBytes, err := json.Marshal(conditionMaps)
	if err != nil {
		return nil, err
	}
	conditionMaps, err = biz.ParseCloudData(jsonBytes)
	if err != nil {
		log.Error("cloud_zone table parse data error(%v)", err)
		return nil, err
	}
	conditionDiff := biz.DifferenceMap(where.Conditions, conditionSources)
	if conditionDiff != nil {
		jsonBytesDiff, err := json.Marshal(conditionDiff)
		if err != nil {
			return nil, err
		}
		conditionDos, err := biz.ParseCloudData(jsonBytesDiff)
		if err != nil {
			log.Error("cloud_zone table parse data error(%v)", err)
			return nil, err
		}
		conditions = biz.GetHandleConditions(conditionDos, "cloud_zone")
		for k, v := range conditionMaps {
			conditions[k] = v
		}
	}
	return conditions, nil
}

// CountCloudZone 云可用区查询数量
func (c *CloudZoneRepo) CountCloudZone(where *biz.CloudZoneWhere) (total int64, err error) {
	conditions, err := c.ParseCloudZone(where) //条件转换
	if err != nil {
		return 0, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("cloud_zone").Model(&biz.CloudZone{})
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
	err = session.Joins("Provider").Joins("Region").Count(&total).Error
	if err != nil {
		log.Error("cloud_zone table query data count error (%v)", err)
		return 0, err
	}
	return total, nil
}

// QueryCloudZone 云可用区查询
func (c *CloudZoneRepo) QueryCloudZone(where *biz.CloudZoneWhere, output *biz.CloudZoneOutput) (list []*biz.CloudZone, err error) {
	conditions, err := c.ParseCloudZone(where) //条件转换
	if err != nil {
		return nil, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("cloud_zone").Model(&biz.CloudZone{})
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
		Joins("Region", c.db.GormDB.Select("id", "name", "cid", "provider_id")).Find(&list).Error
	if err != nil {
		log.Error("cloud_zone table query data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// CreateCloudZone 云可用区创建
func (c *CloudZoneRepo) CreateCloudZone(create []*biz.CloudZone) (list []*biz.CloudZone, err error) {
	for _, v := range create {
		err = c.db.GormDB.Table("cloud_zone").Create(v).Error
		if err != nil {
			log.Error("cloud_zone table create data error(%v)", err)
			return nil, err
		}
		list = append(list, &biz.CloudZone{
			CloudProductCommon: biz.CloudProductCommon{
				ID: v.ID,
			},
		})
	}
	return list, nil
}

// UpdateCloudZone 云可用区更新
func (c *CloudZoneRepo) UpdateCloudZone(where *biz.CloudZoneWhere, update *biz.CloudZone) error {
	updateMap := make(map[string]interface{})
	if update != nil {
		updateMap["name"] = update.Name
		updateMap["cid"] = update.CID
		updateMap["region"] = update.RegionID
		updateMap["provider_id"] = update.ProviderID
		updateMap["status"] = update.Status
	}
	if where != nil && len(updateMap) > 0 {
		err := c.db.GormDB.Table("cloud_zone").Where(where.Query, where.Arg).Updates(updateMap).Error
		if err != nil {
			log.Error("cloud_zone table update data error(%v)", err)
			return err
		}
	}
	return nil
}

// DeleteCloudZone 云可用区删除
func (c *CloudZoneRepo) DeleteCloudZone(deleteID []int) error {
	err := c.db.GormDB.Table("cloud_zone").Where("id IN ?", deleteID).Delete(&biz.CloudZone{}).Error
	if err != nil {
		log.Error("cloud_zone table update data error(%v)", err)
		return err
	}
	return nil
}

// UpsertCloudZone 云可用区更新或创建
func (c *CloudZoneRepo) UpsertCloudZone(upsert []*biz.CloudZoneUpsert) error {
	return nil
}
