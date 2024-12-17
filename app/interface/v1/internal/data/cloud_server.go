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
	"strings"
	"time"

	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/go-kratos/kratos/pkg/log"
	"gorm.io/gorm"
)

// CloudServerRepo 云服务器DB
type CloudServerRepo struct {
	db *biz.DB
}

// NewCloudServerRepo 初始化云服务器DB
func NewCloudServerRepo(db *biz.DB) biz.CloudServerRepo {
	return &CloudServerRepo{
		db: db,
	}
}

// ParseCloudServer 云服务器查询条件处理
func (c *CloudServerRepo) ParseCloudServer(where *biz.CloudServerWhere) (conditions map[string]interface{}, err error) {
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
	if serverSpec, ok := where.Conditions["server_spec"]; ok && serverSpec != nil {
		if serverSpecMap, okMap := serverSpec.(map[string]interface{}); okMap {
			for key, value := range serverSpecMap {
				if mapValues, ok := value.(map[string]interface{}); !ok {
					conditionMaps[fmt.Sprintf("ServerSpec.%v", key)] = value
				} else {
					for subKey, subValue := range mapValues {
						conditionMaps[fmt.Sprintf("ServerSpec__%v.%v", biz.FirstStrUpper(key), subKey)] = subValue
					}
				}
			}
		}
		conditionSources["server_spec"] = true
	}
	if serverImage, ok := where.Conditions["server_image"]; ok && serverImage != nil {
		if serverImageMap, okMap := serverImage.(map[string]interface{}); okMap {
			for key, value := range serverImageMap {
				if mapValues, ok := value.(map[string]interface{}); !ok {
					conditionMaps[fmt.Sprintf("ServerImage.%v", key)] = value
				} else {
					for subKey, subValue := range mapValues {
						conditionMaps[fmt.Sprintf("ServerImage__%v.%v", biz.FirstStrUpper(key), subKey)] = subValue
					}
				}
			}
		}
		conditionSources["server_image"] = true
	}
	if projectConfig, ok := where.Conditions["project_config"]; ok && projectConfig != nil {
		if projectConfigMap, okMap := projectConfig.(map[string]interface{}); okMap {
			for key, value := range projectConfigMap {
				if mapValues, ok := value.(map[string]interface{}); !ok {
					conditionMaps[fmt.Sprintf("ProjectConfig.%v", key)] = value
				} else {
					for subKey, subValue := range mapValues {
						conditionMaps[fmt.Sprintf("ProjectConfig__%v.%v", biz.FirstStrUpper(key), subKey)] = subValue
					}
				}
			}
		}
		conditionSources["project_config"] = true
	}
	if zone, ok := where.Conditions["zone"]; ok && zone != nil {
		if zoneMap, okMap := zone.(map[string]interface{}); okMap {
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
	if subnet, ok := where.Conditions["subnet"]; ok && subnet != nil {
		if subnetMap, okMap := subnet.(map[string]interface{}); okMap {
			for key, value := range subnetMap {
				if mapValues, ok := value.(map[string]interface{}); !ok {
					conditionMaps[fmt.Sprintf("Subnet.%v", key)] = value
				} else {
					for subKey, subValue := range mapValues {
						conditionMaps[fmt.Sprintf("Subnet__%v.%v", biz.FirstStrUpper(key), subKey)] = subValue
					}
				}
			}
		}
		conditionSources["account"] = true
	}
	jsonBytes, err := json.Marshal(conditionMaps)
	if err != nil {
		return nil, err
	}
	conditionMaps, err = biz.ParseCloudData(jsonBytes)
	if err != nil {
		log.Error("cloud_server table parse data error(%v)", err)
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
			log.Error("cloud_server table parse data error(%v)", err)
			return nil, err
		}
		conditions = biz.GetHandleConditions(conditionDos, "cloud_server")
		for k, v := range conditionMaps {
			conditions[k] = v
		}
	}
	return conditions, nil
}

// CountCloudServer 查询云服务器数量
func (c *CloudServerRepo) CountCloudServer(where *biz.CloudServerWhere) (total int64, err error) {
	conditions, err := c.ParseCloudServer(where) //条件转换
	if err != nil {
		return 0, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("cloud_server").Model(&biz.CloudServer{})
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
		Joins("Project").
		Joins("ProjectConfig").
		Joins("Region").
		Joins("Zone").
		Joins("Vpc").
		Joins("Subnet").
		Joins("ServerImage").
		Joins("ServerSpec").
		Joins("LEFT JOIN provider Account_Provider ON Account.provider_id = Account_Provider.id").
		Count(&total).Error
	if err != nil {
		log.Error("cloud_server table query data count error (%v)", err)
		return 0, err
	}
	return total, nil
}

// QueryCloudServer 查询云服务器
func (c *CloudServerRepo) QueryCloudServer(where *biz.CloudServerWhere, output *biz.CloudServerOutput) (list []*biz.CloudServer, err error) {
	conditions, err := c.ParseCloudServer(where) //条件转换
	if err != nil {
		return nil, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("cloud_server").Model(&biz.CloudServer{})
	if where != nil {
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

		if or, ok := conditions["or"]; ok && or != nil { // or条件
			if mapValues, isMap := or.(map[string]interface{}); isMap {
				for k, v := range mapValues {
					session = session.Or(k, v)
				}
			}
			delete(conditions, "or")
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
		Joins("Project", c.db.GormDB.Select("id", "name", "cid")).
		Joins("Region", c.db.GormDB.Select("id", "name", "cid", "provider_id")).
		Joins("Zone", c.db.GormDB.Select("id", "name", "cid")).
		Joins("Vpc", c.db.GormDB.Select("id", "name", "cid")).
		Joins("Subnet", c.db.GormDB.Select("id", "name", "cid")).
		Joins("ServerImage", c.db.GormDB.Select("id", "name", "cid", "os_name", "type")).
		Joins("ServerSpec", c.db.GormDB.Select("id", "name", "cid", "cpu", "gpu", "gpu_model", "memory")).
		Joins("ProjectConfig", c.db.GormDB.Select("id", "name")).
		Joins("LEFT JOIN provider Account_Provider ON Account.provider_id = Account_Provider.id").
		Preload("Account.Provider").
		Find(&list).Error
	if err != nil {
		log.Error("cloud_server table query data error (%v)", err)
		return nil, err
	}
	for _, v := range list {
		var (
			securityGroups []*biz.CloudSecurityGroup
		)
		if v.SecurityGroupCID != "" {
			err = c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("cloud_security_group").Model(&biz.CloudSecurityGroup{}).
				Where("cid IN ?", strings.Split(v.SecurityGroupCID, ",")).
				Find(&securityGroups).Error
			if err != nil {
				log.Error("cloud_server cloud_security_group table query data error (%v)", err)
				return nil, err
			}
			v.SecurityGroup = securityGroups
		}
		if v.ServerSpec != nil && v.ServerSpec.Memory > 0 {
			v.ServerSpec.Memory = v.ServerSpec.Memory / 1024
		}
	}
	return list, nil
}

// CreateCloudServer 创建云服务器
func (c *CloudServerRepo) CreateCloudServer(create []*biz.CloudServer) (list []*biz.CloudServer, err error) {
	for _, v := range create {
		v.CreateTime = time.Now()
		v.UpdateTime = time.Now()
		err = c.db.GormDB.Create(v).Error
		if err != nil {
			log.Error("cloud_Server table create data error(%v)", err)
			return nil, err
		}
		list = append(list, &biz.CloudServer{
			CloudProductCommon: biz.CloudProductCommon{
				ID: v.ID,
			},
		})
	}
	return list, nil
}

// UpdateCloudServer 更新云服务器
func (c *CloudServerRepo) UpdateCloudServer(where *biz.CloudServerWhere, update *biz.CloudServer) error {
	updateMap := make(map[string]interface{})
	if update != nil {
		updateMap["name"] = update.Name
		updateMap["cid"] = update.CID
		updateMap["account_id"] = update.AccountID
		updateMap["project_config_id"] = update.ProjectConfigID
		updateMap["project_id"] = update.ProjectID
		updateMap["region_id"] = update.RegionID
		updateMap["zone_id"] = update.ZoneID
		updateMap["vpc_id"] = update.VpcID
		updateMap["subnet_id"] = update.SubnetID
		updateMap["subnet_cid"] = update.SubnetCID
		updateMap["security_group_cid"] = update.SecurityGroupCID
		updateMap["server_image_id"] = update.ServerImageID
		updateMap["charge_type"] = update.ChargeType
		updateMap["renew_status"] = update.RenewStatus
		updateMap["private_ip"] = update.PrivateIp
		updateMap["public_ip"] = update.PublicIp
		updateMap["status"] = update.Status
		updateMap["expire_time"] = update.ExpireTime
		updateMap["server_spec_id"] = update.ServerSpecID
		updateMap["update_time"] = time.Now()
	}
	if where != nil && len(updateMap) > 0 {
		err := c.db.GormDB.Table("cloud_server").Where(where.Query, where.Arg).Updates(updateMap).Error
		if err != nil {
			log.Error("cloud_server table update data error(%v)", err)
			return err
		}
	}
	return nil
}

// DeleteCloudServer 删除云服务器
func (c *CloudServerRepo) DeleteCloudServer(deleteID []int) error {
	err := c.db.GormDB.Table("cloud_server").Where("id IN ?", deleteID).Delete(&biz.CloudProject{}).Error
	if err != nil {
		log.Error("cloud_server table update data error(%v)", err)
		return err
	}
	return nil
}

// UpsertCloudServer 云服务器更新或创建
func (c *CloudServerRepo) UpsertCloudServer(upsert []*biz.CloudServerUpsert) error {
	return nil
}
