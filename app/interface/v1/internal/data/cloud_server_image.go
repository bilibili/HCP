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

// CloudServerImageRepo 云镜像DB
type CloudServerImageRepo struct {
	db *biz.DB
}

// NewCloudServerImageRepo 初始化云镜像DB
func NewCloudServerImageRepo(db *biz.DB) biz.CloudServerImageRepo {
	return &CloudServerImageRepo{db: db}
}

// ParseCloudServerImage 云镜像查询条件处理
func (c *CloudServerImageRepo) ParseCloudServerImage(where *biz.CloudServerImageWhere) (conditions map[string]interface{}, err error) {
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
		log.Error("cloud_server_image table parse data error(%v)", err)
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
			log.Error("cloud_server_image table parse data error(%v)", err)
			return nil, err
		}
		conditions = biz.GetHandleConditions(conditionDos, "cloud_server_image")
		for k, v := range conditionMaps {
			conditions[k] = v
		}
	}
	return conditions, nil
}

// CountCloudServerImage 云镜像查询数量
func (c *CloudServerImageRepo) CountCloudServerImage(where *biz.CloudServerImageWhere) (total int64, err error) {
	conditions, err := c.ParseCloudServerImage(where) //条件转换
	if err != nil {
		return 0, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("cloud_server_image").Model(&biz.CloudServerImage{})
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
		Joins("Region").
		Joins("LEFT JOIN provider Account_Provider ON Account.provider_id = Account_Provider.id").
		Count(&total).Error
	if err != nil {
		log.Error("account table query data count error (%v)", err)
		return 0, err
	}
	return total, nil
}

// QueryCloudServerImage 云镜像查询
func (c *CloudServerImageRepo) QueryCloudServerImage(where *biz.CloudServerImageWhere, output *biz.CloudServerImageOutput) (list []*biz.CloudServerImage, err error) {
	conditions, err := c.ParseCloudServerImage(where) //条件转换
	if err != nil {
		return nil, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("cloud_server_image").Model(&biz.CloudServerImage{})
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
		Joins("Region", c.db.GormDB.Select("id", "name", "cid", "provider_id")).
		Joins("LEFT JOIN provider Account_Provider ON Account.provider_id = Account_Provider.id").
		Preload("Account.Provider").
		Find(&list).Error
	if err != nil {
		log.Error("cloud_server_image table query data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// CreateCloudServerImage 云镜像创建
func (c *CloudServerImageRepo) CreateCloudServerImage(create []*biz.CloudServerImage) (list []*biz.CloudServerImage, err error) {
	for _, v := range create {
		err = c.db.GormDB.Table("cloud_server_image").Create(v).Error
		if err != nil {
			log.Error("cloud_server_image table create data error(%v)", err)
			return nil, err
		}
		list = append(list, &biz.CloudServerImage{
			CloudProductCommon: biz.CloudProductCommon{
				ID: v.ID,
			},
		})
	}
	return list, nil
}

// CreatBatchesCloudServerImage 云镜像批量创建
func (c *CloudServerImageRepo) CreatBatchesCloudServerImage(create []*biz.CloudServerImage) (list []*biz.CloudServerImage, err error) {
	err = c.db.GormDB.Table("cloud_server_image").CreateInBatches(create, 100).Error
	if err != nil {
		log.Error("cloud_server_image table create batches data error(%v)", err)
		return nil, err
	}
	return list, nil
}

// UpdateCloudServerImage 云镜像更新
func (c *CloudServerImageRepo) UpdateCloudServerImage(where *biz.CloudServerImageWhere, update *biz.CloudServerImage) error {
	updateMap := make(map[string]interface{})
	if update != nil {
		updateMap["name"] = update.Name
		updateMap["cid"] = update.CID
		updateMap["account_id"] = update.AccountID
		updateMap["region_id"] = update.RegionID
		updateMap["os_name"] = update.OsName
		updateMap["type"] = update.Type
		updateMap["status"] = update.Status
	}
	if where != nil && len(updateMap) > 0 {
		err := c.db.GormDB.Table("cloud_server_image").Where(where.Query, where.Arg).Updates(updateMap).Error
		if err != nil {
			log.Error("cloud_project table update data error(%v)", err)
			return err
		}
	}
	return nil
}

// DeleteCloudServerImage 云镜像删除
func (c *CloudServerImageRepo) DeleteCloudServerImage(deleteID []int) error {
	err := c.db.GormDB.Table("cloud_server_image").Where("id IN ?", deleteID).Delete(&biz.CloudProject{}).Error
	if err != nil {
		log.Error("cloud_server_image table delete data error(%v)", err)
		return err
	}
	return nil
}

// UpsertCloudServerImage 云镜像更新或创建
func (c *CloudServerImageRepo) UpsertCloudServerImage(upsert []*biz.CloudServerImageUpsert) error {
	return nil
}
