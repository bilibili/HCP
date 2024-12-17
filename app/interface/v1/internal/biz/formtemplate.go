// Package biz
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
package biz

// FormTemplateRepo terraform接口
type FormTemplateRepo interface {
	QueryFormTemplate(where *FormTemplateWhere) ([]*FormTemplateData, error)
}

// FormTemplateWhere 查询条件
type FormTemplateWhere struct {
	// Query 查询语句
	Query string
	// Arg 查询参数
	Arg interface{}
	// Conditions 查询条件
	Conditions map[string]interface{}
}

// FormTemplateData 表单模板参数结构体
type FormTemplateData struct {
	// ID
	ID int `gorm:"column:id" json:"id"`
	// 模版Name
	Name string `gorm:"column:name" json:"name"`
	// 模版Title
	Title string `gorm:"column:title" json:"title"`
	// 操作
	Operate string `gorm:"column:operate" json:"operate"`
	// 模版数据
	Data string `gorm:"column:data" json:"data"`
}

// FormTemplateUseCase terraform模板业务逻辑
type FormTemplateUseCase struct {
	repo FormTemplateRepo
}

// NewFormTemplateUseCase 创建terraform模板业务逻辑
func NewFormTemplateUseCase(repo FormTemplateRepo) *FormTemplateUseCase {
	return &FormTemplateUseCase{repo: repo}
}

// QueryFormTemplate 查询terraform模板
func (c *FormTemplateUseCase) QueryFormTemplate(where *FormTemplateWhere) ([]*FormTemplateData, error) {
	return c.repo.QueryFormTemplate(where)
}
