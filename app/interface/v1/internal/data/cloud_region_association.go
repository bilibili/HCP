//Package data
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

// CloudRegionAssociationRepo 云地域关联关系数据DB
type CloudRegionAssociationRepo struct {
	db *biz.DB
}

// NewCloudRegionAssociationRepo 初始化云地域关联关系数据DB
func NewCloudRegionAssociationRepo(db *biz.DB) biz.CloudRegionAssociationRepo {
	return &CloudRegionAssociationRepo{
		db: db,
	}
}

// ParseCloudRegionAssociation 云地域关联关系查询条件处理
func (c *CloudRegionAssociationRepo) ParseCloudRegionAssociation(where *biz.CloudRegionAssociationWhere) (conditions map[string]interface{}, err error) {
	conditionMaps, conditionSources := make(map[string]interface{}), make(map[string]interface{})
	if account, ok := where.Conditions["account"]; ok && account != nil {
		jsonBytes, err := json.Marshal(account)
		if err != nil {
			return nil, err
		}
		accountConditions, err := biz.ParseCloudData(jsonBytes)
		if err != nil {
			log.Error("cloud_region_association account table parse data error(%v)", err)
			return nil, err
		}
		for key, value := range accountConditions {
			conditionMaps[fmt.Sprintf("Account.%v", key)] = value
		}
		conditionSources["account"] = true
	}
	if cloudRegion, ok := where.Conditions["region"]; ok && cloudRegion != nil {
		jsonBytes, err := json.Marshal(cloudRegion)
		if err != nil {
			return nil, err
		}
		cloudRegionConditions, err := biz.ParseCloudData(jsonBytes)
		if err != nil {
			log.Error("cloud_region_association cloud_region table parse data error(%v)", err)
		}
		for key, value := range cloudRegionConditions {
			conditionMaps[fmt.Sprintf("Region.%v", key)] = value
		}
		conditionSources["region"] = true
	}
	if provider, ok := where.Conditions["provider"]; ok && provider != nil {
		jsonBytes, err := json.Marshal(provider)
		if err != nil {
			return nil, err
		}
		providerConditions, err := biz.ParseCloudData(jsonBytes)
		if err != nil {
			log.Error("cloud_region_association provider table parse data error(%v)", err)
		}
		for key, value := range providerConditions {
			conditionMaps[fmt.Sprintf("Provider.%v", key)] = value
		}
		conditionSources["provider"] = true
	}
	conditionDiff := biz.DifferenceMap(where.Conditions, conditionSources)
	if conditionDiff != nil {
		jsonBytes, err := json.Marshal(conditionDiff)
		if err != nil {
			return nil, err
		}
		conditionDos, err := biz.ParseCloudData(jsonBytes)
		if err != nil {
			log.Error("cloud_region_association table parse data error(%v)", err)
			return nil, err
		}
		conditions = biz.GetHandleConditions(conditionDos, "cloud_region_association")
		for k, v := range conditionMaps {
			conditions[k] = v
		}
	}
	return conditions, nil
}

// QueryCloudRegionAssociation 云地域关联关系查询
func (c *CloudRegionAssociationRepo) QueryCloudRegionAssociation(where *biz.CloudRegionAssociationWhere) (list []*biz.CloudRegionAssociation, err error) {
	conditions, err := c.ParseCloudRegionAssociation(where) //条件转换
	if err != nil {
		return nil, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("cloud_region_association").Model(&biz.CloudRegionAssociation{})
	if where != nil {
		for key, value := range conditions {
			session.Where(key, value)
		}
	}
	err = session.Joins("Account", c.db.GormDB.Select("id", "name", "cid", "provider_id")).
		Joins("Region", c.db.GormDB.Select("id", "name")).
		Find(&list).Error
	if err != nil {
		log.Error("cloud_region_association table query data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// CreateCloudRegionAssociation 云地域关联关系创建
func (c *CloudRegionAssociationRepo) CreateCloudRegionAssociation(create []*biz.CloudRegionAssociation) (list []*biz.CloudRegionAssociation, err error) {
	for _, v := range create {
		err = c.db.GormDB.Create(v).Error
		if err != nil {
			log.Error("cloud_region_association table create data error(%v)", err)
			return nil, err
		}
		list = append(list, &biz.CloudRegionAssociation{
			ID: v.ID,
		})
	}
	return list, nil
}

// UpdateCloudRegionAssociation 云地域关联关系更新
func (c *CloudRegionAssociationRepo) UpdateCloudRegionAssociation(where *biz.CloudRegionAssociationWhere, update *biz.CloudRegionAssociation) error {
	var updateMap map[string]interface{}
	if update != nil {
		updateMap["account_id"] = update.AccountID
		updateMap["region_id"] = update.RegionID
	}
	if where != nil && len(updateMap) > 0 {
		err := c.db.GormDB.Where(where.Query, where.Arg).Updates(updateMap).Error
		if err != nil {
			log.Error("cloud_region_association table update data error(%v)", err)
			return err
		}
	}
	return nil
}

// DeleteCloudRegionAssociation 云地域关联关系删除
func (c *CloudRegionAssociationRepo) DeleteCloudRegionAssociation(deleteID []int) error {
	err := c.db.GormDB.Table("cloud_region_association").Where("id IN ?", deleteID).Delete(&biz.CloudProject{}).Error
	if err != nil {
		log.Error("cloud_region_association table update data error(%v)", err)
		return err
	}
	return nil
}

// UpsertCloudRegionAssociation 云地域关联关系更新或创建
func (c *CloudRegionAssociationRepo) UpsertCloudRegionAssociation(upsert []*biz.CloudRegionAssociationUpsert) error {
	return nil
}
