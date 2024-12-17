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

// FormTemplateRepo 表单模版数据DB
type FormTemplateRepo struct {
	db *biz.DB
}

// NewFormTemplateRepo 初始化表单模版数据DB
func NewFormTemplateRepo(db *biz.DB) biz.FormTemplateRepo {
	return &FormTemplateRepo{
		db: db,
	}
}

// ParseFormTemplate 表单模版数据查询条件处理
func (f *FormTemplateRepo) ParseFormTemplate(where *biz.FormTemplateWhere) (conditions map[string]interface{}, err error) {
	jsonBytes, err := json.Marshal(where.Conditions)
	if err != nil {
		return nil, err
	}
	conditions, err = biz.ParseCloudData(jsonBytes)
	if err != nil {
		log.Error("form_template table parse data error(%v)", err)
		return nil, err
	}
	for k, v := range conditions {
		conditions[k] = v
	}
	return conditions, nil
}

// QueryFormTemplate 表单模版数据查询
func (f *FormTemplateRepo) QueryFormTemplate(where *biz.FormTemplateWhere) (list []*biz.FormTemplateData, err error) {
	conditions, err := f.ParseFormTemplate(where) //条件转换
	if err != nil {
		return nil, err
	}
	session := f.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("form_template_data").Model(&biz.FormTemplate{})
	if where != nil {
		for key, value := range conditions {
			session.Where(key, value)
		}
	}
	err = session.Find(&list).Error
	if err != nil {
		log.Error("form_template_data table query data error (%v)", err)
		return nil, err
	}
	return list, nil
}
