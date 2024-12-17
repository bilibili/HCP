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

// CloudSecurityGroupRepo 云安全组数据DB
type CloudSecurityGroupRepo struct {
	db *biz.DB
}

// NewCloudSecurityGroupRepo 初始化云安全组数据DB
func NewCloudSecurityGroupRepo(db *biz.DB) biz.CloudSecurityGroupRepo {
	return &CloudSecurityGroupRepo{db: db}
}

// ParseCloudSecurityGroup 云安全组查询条件处理
func (c *CloudSecurityGroupRepo) ParseCloudSecurityGroup(where *biz.CloudSecurityGroupWhere) (conditions map[string]interface{}, err error) {
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
					mode := biz.GetFuzzyOrPrecise(key) //获取精确查询、IN查询（兼容关联项目查询）
					if mode == 1 {
						conditionMaps[fmt.Sprintf("ProjectConfig.%v", key)] = value
					} else if mode == 2 {
						conditionMaps[fmt.Sprintf("ProjectConfig.%s_IN ?", biz.GetLastUnderscoreString(key))] = value
					}
				} else {
					for subKey, subValue := range mapValues {
						conditionMaps[fmt.Sprintf("ProjectConfig__%v.%v", biz.FirstStrUpper(key), subKey)] = subValue
					}
				}
			}
		}
		conditionSources["project_config"] = true
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
		log.Error("cloud_security_group table parse data error(%v)", err)
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
			log.Error("cloud_security_group table parse data error(%v)", err)
			return nil, err
		}
		conditions = biz.GetHandleConditions(conditionDos, "cloud_security_group")
		for k, v := range conditionMaps {
			conditions[k] = v
		}
	}
	return conditions, nil
}

// CountCloudSecurityGroup 云安全组查询数量
func (c *CloudSecurityGroupRepo) CountCloudSecurityGroup(where *biz.CloudSecurityGroupWhere) (total int64, err error) {
	conditions, err := c.ParseCloudSecurityGroup(where) //条件转换
	if err != nil {
		return 0, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("cloud_security_group").Model(&biz.CloudSecurityGroup{})
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
	err = session.Joins("Account").Joins("Project").Joins("Region").Joins("Vpc").
		Joins("LEFT JOIN project_account_config ProjectAccountConfig ON ProjectAccountConfig.account_id = cloud_security_group.account_id AND ProjectAccountConfig.project_id = cloud_security_group.project_id").
		Joins("LEFT JOIN project_config ProjectConfig ON ProjectAccountConfig.project_config_id = ProjectConfig.id").
		Joins("LEFT JOIN provider Account_Provider ON Account.provider_id = Account_Provider.id").
		Count(&total).Error
	if err != nil {
		log.Error("cloud_security_group table query data count error (%v)", err)
		return 0, err
	}
	return total, nil
}

// QueryCloudSecurityGroup 云安全组查询
func (c *CloudSecurityGroupRepo) QueryCloudSecurityGroup(where *biz.CloudSecurityGroupWhere, output *biz.CloudSecurityGroupOutput) (list []*biz.CloudSecurityGroup, err error) {
	conditions, err := c.ParseCloudSecurityGroup(where) //条件转换
	if err != nil {
		return nil, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("cloud_security_group").Model(&biz.CloudSecurityGroup{})
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
	err = session.
		Joins("Account", c.db.GormDB.Select("id", "name", "cid", "alias", "provider_id")).
		Joins("Project", c.db.GormDB.Select("id", "name")).
		Joins("Region", c.db.GormDB.Select("id", "name")).
		Joins("Vpc", c.db.GormDB.Select("id", "name", "cid")).
		Joins("LEFT JOIN project_account_config ProjectAccountConfig ON ProjectAccountConfig.account_id = cloud_security_group.account_id AND ProjectAccountConfig.project_id = cloud_security_group.project_id").
		Joins("LEFT JOIN project_config ProjectConfig ON ProjectAccountConfig.project_config_id = ProjectConfig.id").
		Joins("LEFT JOIN provider Account_Provider ON Account.provider_id = Account_Provider.id").
		Preload("Account.Provider").
		Preload("ProjectAccountConfig.ProjectConfig").
		Find(&list).Error
	for k, v := range list {
		if v.ProjectAccountConfig != nil && v.ProjectAccountConfig.ProjectConfig != nil {
			list[k].ProjectConfig = v.ProjectAccountConfig.ProjectConfig
		}
		list[k].ProjectAccountConfig = nil
	}
	if err != nil {
		log.Error("cloud_security_group table query data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// CreateCloudSecurityGroup 云安全组创建
func (c *CloudSecurityGroupRepo) CreateCloudSecurityGroup(create []*biz.CloudSecurityGroup) (list []*biz.CloudSecurityGroup, err error) {
	for _, v := range create {
		v.CreateTime = time.Now()
		v.UpdateTime = time.Now()
		err = c.db.GormDB.Table("cloud_security_group").Create(v).Error
		if err != nil {
			log.Error("cloud_security_group table create data error(%v)", err)
			return nil, err
		}
		list = append(list, &biz.CloudSecurityGroup{
			CloudProductCommon: biz.CloudProductCommon{
				ID: v.ID,
			},
		})
	}
	return list, nil
}

// UpdateCloudSecurityGroup 云安全组更新
func (c *CloudSecurityGroupRepo) UpdateCloudSecurityGroup(where *biz.CloudSecurityGroupWhere, update *biz.CloudSecurityGroup) error {
	updateMap := make(map[string]interface{})
	if update != nil {
		updateMap["name"] = update.Name
		updateMap["cid"] = update.CID
		updateMap["project_id"] = update.ProjectID
		updateMap["vpc_id"] = update.VpcID
		updateMap["region_id"] = update.RegionID
		updateMap["account_id"] = update.AccountID
		updateMap["status"] = update.Status
		updateMap["update_time"] = time.Now()
	}
	if where != nil && len(updateMap) > 0 {
		err := c.db.GormDB.Table("cloud_security_group").Where(where.Query, where.Arg).Updates(updateMap).Error
		if err != nil {
			log.Error("cloud_security_group table update data error(%v)", err)
			return err
		}
	}
	return nil
}

// DeleteCloudSecurityGroup 云安全组删除
func (c *CloudSecurityGroupRepo) DeleteCloudSecurityGroup(deleteID []int) error {
	err := c.db.GormDB.Table("cloud_security_group").Where("id IN ?", deleteID).Delete(&biz.CloudSecurityGroup{}).Error
	if err != nil {
		log.Error("cloud_security_group table update data error(%v)", err)
		return err
	}
	return nil
}

// UpsertCloudSecurityGroup 云安全组更新或创建
func (c *CloudSecurityGroupRepo) UpsertCloudSecurityGroup(upsert []*biz.CloudSecurityGroupUpsert) error {
	return nil
}
