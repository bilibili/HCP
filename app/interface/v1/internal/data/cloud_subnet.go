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

// CloudSubnetRepo 云子网DB
type CloudSubnetRepo struct {
	db *biz.DB
}

// NewCloudSubnetRepo 初始化云子网DB
func NewCloudSubnetRepo(db *biz.DB) biz.CloudSubnetRepo {
	return &CloudSubnetRepo{db: db}
}

// ParseCloudSubnet 云子网查询条件处理
func (c *CloudSubnetRepo) ParseCloudSubnet(where *biz.CloudSubnetWhere) (conditions map[string]interface{}, err error) {
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
	if cloudZone, ok := where.Conditions["zone"]; ok && cloudZone != nil {
		if zoneMap, okMap := cloudZone.(map[string]interface{}); okMap {
			for key, value := range zoneMap {
				if mapValues, ok := value.(map[string]interface{}); !ok {
					conditionMaps[fmt.Sprintf("Zone.%v", key)] = value
				} else {
					for subKey, subValue := range mapValues {
						conditionMaps[fmt.Sprintf("Zone__%v.%v", biz.FirstStrUpper(key), subKey)] = subValue
					}
				}
			}
		}
		conditionSources["zone"] = true
	}
	if cloudVpc, ok := where.Conditions["vpc"]; ok && cloudVpc != nil {
		if vpcMap, okMap := cloudVpc.(map[string]interface{}); okMap {
			for key, value := range vpcMap {
				if mapValues, ok := value.(map[string]interface{}); !ok {
					conditionMaps[fmt.Sprintf("Vpc.%v", key)] = value
				} else {
					for subKey, subValue := range mapValues {
						conditionMaps[fmt.Sprintf("Vpc__%v.%v", biz.FirstStrUpper(key), subKey)] = subValue
					}
				}
			}
		}
		conditionSources["vpc"] = true
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
	jsonBytes, err := json.Marshal(conditionMaps)
	if err != nil {
		return nil, err
	}
	conditionMaps, err = biz.ParseCloudData(jsonBytes)
	if err != nil {
		log.Error("cloud_subnet table parse data error(%v)", err)
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
			log.Error("cloud_subnet table parse data error(%v)", err)
			return nil, err
		}
		conditions = biz.GetHandleConditions(conditionDos, "cloud_subnet")
		for k, v := range conditionMaps {
			conditions[k] = v
		}
	}
	return conditions, nil
}

// CountCloudSubnet 云子网查询数量
func (c *CloudSubnetRepo) CountCloudSubnet(where *biz.CloudSubnetWhere) (total int64, err error) {
	conditions, err := c.ParseCloudSubnet(where) //条件转换
	if err != nil {
		return 0, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("cloud_subnet").Model(&biz.CloudSubnet{})
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
	err = session.Joins("Account").Joins("Zone").Joins("Vpc").Joins("Project").
		Joins("LEFT JOIN project_account_config ProjectAccountConfig ON ProjectAccountConfig.account_id = cloud_subnet.account_id AND ProjectAccountConfig.project_id = cloud_subnet.project_id").
		Joins("LEFT JOIN project_config ProjectConfig ON ProjectAccountConfig.project_config_id = ProjectConfig.id").
		Joins("LEFT JOIN provider Account_Provider ON Account.provider_id = Account_Provider.id").
		Joins("LEFT JOIN cloud_region Vpc_Region ON Vpc.region_id = Vpc_Region.id").
		Count(&total).Error
	if err != nil {
		log.Error("cloud_subnet table query data count error (%v)", err)
		return 0, err
	}
	return total, nil
}

// QueryCloudSubnet 云子网查询
func (c *CloudSubnetRepo) QueryCloudSubnet(where *biz.CloudSubnetWhere, output *biz.CloudSubnetOutput) (list []*biz.CloudSubnet, err error) {
	conditions, err := c.ParseCloudSubnet(where) //条件转换
	if err != nil {
		return nil, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("cloud_subnet").Model(&biz.CloudSubnet{})
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
		Joins("Zone", c.db.GormDB.Select("id", "name", "cid")).
		Joins("Vpc", c.db.GormDB.Select("id", "name", "cid", "region_id")).
		Joins("Project", c.db.GormDB.Select("id", "name", "cid")).
		Joins("LEFT JOIN project_account_config ProjectAccountConfig ON ProjectAccountConfig.account_id = cloud_subnet.account_id AND ProjectAccountConfig.project_id = cloud_subnet.project_id").
		Joins("LEFT JOIN project_config ProjectConfig ON ProjectAccountConfig.project_config_id = ProjectConfig.id").
		Joins("LEFT JOIN provider Account_Provider ON Account.provider_id = Account_Provider.id").
		Joins("LEFT JOIN cloud_region Vpc_Region ON Vpc.region_id = Vpc_Region.id").
		Preload("ProjectAccountConfig.ProjectConfig").
		Preload("Vpc.Region").
		Preload("Account.Provider").
		Find(&list).Error
	for k, v := range list {
		if v.ProjectAccountConfig != nil && v.ProjectAccountConfig.ProjectConfig != nil {
			list[k].ProjectConfig = v.ProjectAccountConfig.ProjectConfig
		}
		list[k].ProjectAccountConfig = nil
	}
	if err != nil {
		log.Error("cloud_subnet table query data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// CreateCloudSubnet 云子网创建
func (c *CloudSubnetRepo) CreateCloudSubnet(create []*biz.CloudSubnet) (list []*biz.CloudSubnet, err error) {
	for _, v := range create {
		err = c.db.GormDB.Table("cloud_subnet").Create(v).Error
		if err != nil {
			log.Error("cloud_subnet table create data error(%v)", err)
			return nil, err
		}
		list = append(list, &biz.CloudSubnet{
			CloudProductCommon: biz.CloudProductCommon{
				ID: v.ID,
			},
		})
	}
	return list, nil
}

// UpdateCloudSubnet 云子网更新
func (c *CloudSubnetRepo) UpdateCloudSubnet(where *biz.CloudSubnetWhere, update *biz.CloudSubnet) error {
	updateMap := make(map[string]interface{})
	if update != nil {
		updateMap["name"] = update.Name
		updateMap["cid"] = update.CID
		updateMap["cidr"] = update.Cidr
		updateMap["zone_id"] = update.ZoneID
		updateMap["vpc_id"] = update.VpcID
		updateMap["account_id"] = update.AccountID
		updateMap["region_id"] = update.RegionID
		updateMap["project_id"] = update.ProjectID
		updateMap["status"] = update.Status
	}
	if where != nil && len(updateMap) > 0 {
		err := c.db.GormDB.Table("cloud_subnet").Where(where.Query, where.Arg).Updates(updateMap).Error
		if err != nil {
			log.Error("cloud_subnet table update data error(%v)", err)
			return err
		}
	}
	return nil
}

// DeleteCloudSubnet 云子网删除
func (c *CloudSubnetRepo) DeleteCloudSubnet(deleteID []int) error {
	err := c.db.GormDB.Table("cloud_subnet").Where("id IN ?", deleteID).Delete(&biz.CloudSubnet{}).Error
	if err != nil {
		log.Error("cloud_subnet table delete data error(%v)", err)
		return err
	}
	return nil
}

// UpsertCloudSubnet 云子网更新或创建
func (c *CloudSubnetRepo) UpsertCloudSubnet(upsert []*biz.CloudSubnetUpsert) error {
	return nil
}
