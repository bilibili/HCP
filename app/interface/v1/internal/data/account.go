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

	"github.com/bilibili/HCP/app/interface/v1/configs"
	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/bilibili/HCP/utils/aes"
	"github.com/go-kratos/kratos/pkg/log"
	"gorm.io/gorm"
)

// AccountRepo 云账号DB
type AccountRepo struct {
	db *biz.DB
}

// NewAccountRepo 初始化云账号DB
func NewAccountRepo(db *biz.DB) biz.AccountRepo {
	return &AccountRepo{
		db: db,
	}
}

// ParseAccount 账号查询条件处理
func (c *AccountRepo) ParseAccount(where *biz.AccountWhere) (conditions map[string]interface{}, err error) {
	conditionMaps, conditionSources := make(map[string]interface{}), make(map[string]interface{})
	if provider, ok := where.Conditions["provider"]; ok && provider != nil {
		jsonBytes, err := json.Marshal(provider)
		if err != nil {
			return nil, err
		}
		providerConditions, err := biz.ParseCloudData(jsonBytes)
		if err != nil {
			log.Error("account provider table parse data error(%v)", err)
			return nil, err
		}
		for key, value := range providerConditions {
			conditionMaps[fmt.Sprintf("Provider.%v", key)] = value
		}
		conditionSources["provider"] = true
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
	conditionDiff := biz.DifferenceMap(where.Conditions, conditionSources)
	if conditionDiff != nil {
		jsonBytes, err := json.Marshal(conditionDiff)
		if err != nil {
			return nil, err
		}
		conditionDos, err := biz.ParseCloudData(jsonBytes)
		if err != nil {
			log.Error("account table parse data error(%v)", err)
			return nil, err
		}
		conditions = biz.GetHandleConditions(conditionDos, "account")
		for k, v := range conditionMaps {
			conditions[k] = v
		}
	}
	return conditions, nil
}

// CountAccount 账号查询数量
func (c *AccountRepo) CountAccount(where *biz.AccountWhere) (total int64, err error) {
	conditions, err := c.ParseAccount(where) //条件转换
	if err != nil {
		return 0, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("account").Model(&biz.Account{})
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
	err = session.Joins("Provider").
		Joins("LEFT JOIN project_account_config ProjectAccountConfig ON ProjectAccountConfig.account_id = account.id").
		Joins("LEFT JOIN project_config ProjectConfig ON ProjectAccountConfig.project_config_id = ProjectConfig.id").
		Group("account.id").
		Count(&total).Error
	if err != nil {
		log.Error("account table query data count error (%v)", err)
		return 0, err
	}
	return total, nil
}

// QueryAccount 账号查询
func (c *AccountRepo) QueryAccount(where *biz.AccountWhere, output *biz.AccountOutput) (list []*biz.Account, err error) {
	conditions, err := c.ParseAccount(where) //条件转换
	if err != nil {
		return nil, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("account").Model(&biz.Account{})
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
	err = session.Joins("Provider", c.db.GormDB.Select("id", "name", "alias")).
		Joins("LEFT JOIN project_account_config ProjectAccountConfig ON ProjectAccountConfig.account_id = account.id").
		Joins("LEFT JOIN project_config ProjectConfig ON ProjectAccountConfig.project_config_id = ProjectConfig.id").
		Group("account.id").
		Find(&list).Error
	if err != nil {
		log.Error("account table query data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// CreateAccount 账号创建
func (c *AccountRepo) CreateAccount(create []*biz.Account) (list []*biz.Account, err error) {
	for _, v := range create {
		v.CreateTime = time.Now()
		v.UpdateTime = time.Now()
		v.DeleteTime = timeParse("0001-01-01T00:00:00+08:00")
		if v.SyncSecretKey != "" {
			enSecretKey, err := aes.Encrypt(configs.Conf.CloudSecret.SecretAesKey, v.SyncSecretKey)
			if err != nil {
				log.Error("aes secretKey encrypt error(%v)", err)
				return nil, err
			}
			v.SyncSecretKey = enSecretKey
		}
		if v.OperateSecretKey != "" {
			enOperateSecretKey, err := aes.Encrypt(configs.Conf.CloudSecret.SecretAesKey, v.OperateSecretKey)
			if err != nil {
				log.Error("aes operateSecretKey encrypt error(%v)", err)
				return nil, err
			}
			v.OperateSecretKey = enOperateSecretKey
		}
		err = c.db.GormDB.Table("account").Create(v).Error
		if err != nil {
			log.Error("account table create data error(%v)", err)
			return nil, err
		}
		list = append(list, &biz.Account{
			CloudProductCommon: biz.CloudProductCommon{
				ID: v.ID,
			},
		})
	}
	return list, nil
}

// UpdateAccount 账号更新
func (c *AccountRepo) UpdateAccount(where *biz.AccountWhere, update *biz.Account) error {
	updateMap := make(map[string]interface{})
	if update != nil {
		updateMap["name"] = update.Name
		updateMap["cid"] = update.CID
		updateMap["alias"] = update.Alias
		updateMap["provider_id"] = update.ProviderID
		updateMap["sync_secret_id"] = update.SyncSecretID
		updateMap["operate_secret_id"] = update.OperateSecretID
		updateMap["email"] = update.Email
		updateMap["phone"] = update.Phone
		updateMap["description"] = update.Description
		updateMap["subject"] = update.Subject
		updateMap["key_pair_id"] = update.KeyPairID
		updateMap["key_pair_name"] = update.KeyPairName
	}
	if update.OperateSecretKey != "" {
		enOperateSecretKey, err := aes.Encrypt(configs.Conf.CloudSecret.SecretAesKey, update.OperateSecretKey)
		if err != nil {
			log.Error("aes operateSecretKey encrypt error(%v)", err)
			return err
		}
		updateMap["operate_secret_key"] = enOperateSecretKey
	}
	if update.SyncSecretKey != "" {
		enSecretKey, err := aes.Encrypt(configs.Conf.CloudSecret.SecretAesKey, update.SyncSecretKey)
		if err != nil {
			log.Error("aes secretKey encrypt error(%v)", err)
		}
		updateMap["sync_secret_key"] = enSecretKey
	}
	if update.IsDelete == 1 {
		updateMap["delete_time"] = time.Now()
	}
	if where != nil && len(updateMap) > 0 {
		err := c.db.GormDB.Table("account").Where(where.Query, where.Arg).Updates(updateMap).Error
		if err != nil {
			log.Error("account table update data error(%v)", err)
			return nil
		}
	}
	return nil
}

// DeleteAccount 账号删除
func (c *AccountRepo) DeleteAccount(softDelete int, deleteID []int) error {
	if softDelete == 1 { // 软删除
		now := time.Now()
		deleteTime := now.Format("2006-01-02 15:04:05")
		var updateMap map[string]interface{}
		updateMap["is_delete"] = 1
		updateMap["delete_time"] = deleteTime
		err := c.db.GormDB.Table("account").Where("id IN ?", deleteID).Updates(updateMap).Error
		if err != nil {
			log.Error("account table delete data error(%v)", err)
			return nil
		}
	}
	if softDelete == 0 { // 硬删除
		err := c.db.GormDB.Table("account").Where("id IN ?", deleteID).Delete(&biz.Account{}).Error
		if err != nil {
			log.Error("account table update data error(%v)", err)
			return nil
		}
	}
	return nil
}

// UpsertAccount 账号更新或创建
func (c *AccountRepo) UpsertAccount(upsert []*biz.AccountUpsert) error {
	return nil
}
