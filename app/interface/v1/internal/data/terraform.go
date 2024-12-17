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

	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/go-kratos/kratos/pkg/log"
	"gorm.io/gorm"
)

// TerraformRepo terraform模版数据DB
type TerraformRepo struct {
	db *biz.DB
}

// NewTerraformRepo 初始化terraform模版数据DB
func NewTerraformRepo(db *biz.DB) biz.TerraformRepo {
	return &TerraformRepo{db: db}
}

// ParseTerraform terraform模版查询条件处理
func (c *TerraformRepo) ParseTerraform(where *biz.TerraformWhere) (conditions map[string]interface{}, err error) {
	if where.Conditions != nil {
		jsonBytes, err := json.Marshal(where.Conditions)
		if err != nil {
			return nil, err
		}
		conditionDos, err := biz.ParseCloudData(jsonBytes)
		if err != nil {
			log.Error("terraform_template_data table parse data error(%v)", err)
			return nil, err
		}
		conditions = biz.GetHandleConditions(conditionDos, "")
	}
	return conditions, nil
}

// CountTerraform terraform模版查询数量
func (c *TerraformRepo) CountTerraform(where *biz.TerraformWhere) (total int64, err error) {
	conditions, err := c.ParseTerraform(where) //条件转换
	if err != nil {
		return 0, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("terraform_template_data").Model(&biz.Terraform{})
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
	err = session.Count(&total).Error
	if err != nil {
		log.Error("terraform_template_data table count error (%v)", err)
		return 0, err
	}
	return total, nil
}

// QueryTerraform terraform模版查询
func (c *TerraformRepo) QueryTerraform(where *biz.TerraformWhere, output *biz.TerraformOutput) (list []*biz.Terraform, err error) {
	conditions, err := c.ParseTerraform(where) //条件转换
	if err != nil {
		return nil, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("terraform_template_data").Model(&biz.Terraform{})
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
	err = session.Find(&list).Error
	if err != nil {
		log.Error("terraform_template_data table query data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// CreateTerraform terraform模版创建
func (c *TerraformRepo) CreateTerraform(create []*biz.Terraform) (list []*biz.Terraform, err error) {
	for _, v := range create {
		err = c.db.GormDB.Table("terraform_template_data").Create(v).Error
		if err != nil {
			log.Error("terraform_template_data table create data error(%v)", err)
			return nil, err
		}
		list = append(list, &biz.Terraform{
			ID: v.ID,
		})
	}
	return list, nil
}

// UpdateTerraform terraform模版更新
func (c *TerraformRepo) UpdateTerraform(where *biz.TerraformWhere, update *biz.Terraform) error {
	updateMap := make(map[string]interface{})
	if update != nil {
		updateMap["name"] = update.Name
		updateMap["provider_id"] = update.ProviderID
		updateMap["operate"] = update.Operate
		updateMap["dependent_parameters"] = update.DependentParameters
		updateMap["data"] = update.Data
	}
	if where != nil && len(updateMap) > 0 {
		err := c.db.GormDB.Table("terraform_template_data").Where(where.Query, where.Arg).Updates(updateMap).Error
		if err != nil {
			log.Error("terraform_template_data table update data error(%v)", err)
			return err
		}
	}
	return nil
}

// DeleteTerraform terraform模版删除
func (c *TerraformRepo) DeleteTerraform(deleteID []int) error {
	err := c.db.GormDB.Table("terraform_template_data").Where("id IN ?", deleteID).Delete(&biz.Terraform{}).Error
	if err != nil {
		log.Error("terraform_template_data table update data error(%v)", err)
		return err
	}
	return nil
}

// UpsertTerraform terraform模版更新或创建
func (c *TerraformRepo) UpsertTerraform(upsert []*biz.TerraformUpsert) error {
	return nil
}
