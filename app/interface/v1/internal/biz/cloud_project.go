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

import (
	"fmt"
	"strings"

	"github.com/bilibili/HCP/app/interface/v1/configs"
	"github.com/bilibili/HCP/cloudrepo"
)

// CloudProjectRepo 云项目接口
type CloudProjectRepo interface {
	// CountCloudProject 查询云项目数量
	CountCloudProject(where *CloudProjectWhere) (int64, error)
	// QueryCloudProject 查询云项目
	QueryCloudProject(where *CloudProjectWhere, output *CloudProjectOutput) ([]*CloudProject, error)
	// CreateCloudProject 创建云项目
	CreateCloudProject(create []*CloudProject) ([]*CloudProject, error)
	// UpdateCloudProject 更新云项目
	UpdateCloudProject(where *CloudProjectWhere, update *CloudProject) error
	// DeleteCloudProject 删除云项目
	DeleteCloudProject(deleteID []int) error
	// UpsertCloudProject 更新或插入云项目
	UpsertCloudProject(upsert []*CloudProjectUpsert) error
}

// CloudProjectWhere 云项目查询条件
type CloudProjectWhere struct {
	// 查询语句
	Query string
	// 查询参数
	Arg interface{}
	// 查询条件
	Conditions map[string]interface{}
}

// CloudProjectOutput 云项目查询输出条件参数
type CloudProjectOutput struct {
	OutPutCommon
}

// CloudProject 云项目参数结构体
type CloudProject struct {
	CloudProductCommon
	// 云厂商ID
	ProviderID int `gorm:"column:provider_id" json:"provider_id"`
	// 云厂商实体
	Provider *Provider `json:"provider"`
	// 云账号ID
	AccountID int `gorm:"column:account_id" json:"account_id"`
	// 云账号实体
	Account *Account `json:"account"`
}

// CloudProjectUpsert 云项目更新或插入结构体
type CloudProjectUpsert struct {
	CloudProject
}

// CloudProjectUseCase 云项目业务逻辑
type CloudProjectUseCase struct {
	repo CloudProjectRepo
}

// NewCloudProjectUseCase 创建云项目业务逻辑
func NewCloudProjectUseCase(repo CloudProjectRepo) *CloudProjectUseCase {
	return &CloudProjectUseCase{repo: repo}
}

// QueryCloudProject 查询云项目
func (c *CloudProjectUseCase) QueryCloudProject(where *CloudProjectWhere, output *CloudProjectOutput) ([]*CloudProject, error) {
	return c.repo.QueryCloudProject(where, output)
}

// CountCloudProject 查询云项目数量
func (c *CloudProjectUseCase) CountCloudProject(where *CloudProjectWhere) (int64, error) {
	return c.repo.CountCloudProject(where)
}

// CreateCloudProject 创建云项目
func (c *CloudProjectUseCase) CreateCloudProject(create []*CloudProject) ([]*CloudProject, error) {
	return c.repo.CreateCloudProject(create)
}

// UpdateCloudProject 更新云项目
func (c *CloudProjectUseCase) UpdateCloudProject(where *CloudProjectWhere, update *CloudProject) error {
	return c.repo.UpdateCloudProject(where, update)
}

// DeleteCloudProject 删除云项目
func (c *CloudProjectUseCase) DeleteCloudProject(deleteID []int) error {
	return c.repo.DeleteCloudProject(deleteID)
}

// UpsertCloudProject 更新或插入云项目
func (c *CloudProjectUseCase) UpsertCloudProject(upsert []*CloudProjectUpsert) error {
	return c.repo.UpsertCloudProject(upsert)
}

// ConvertTagToProject 标签转换为云项目
func (c *CloudProjectUseCase) ConvertTagToProject(inputs []cloudrepo.CloudTag, ProviderID, AccountID int) []*CloudProject {
	var projects []*CloudProject
	if ProviderID == 1 || ProviderID == 3 || ProviderID == 5 { //只处理百度云、腾讯云、AWS(项目标签Key为ProjectName)
		for _, v := range inputs {
			if v.GetTagKey() == configs.Conf.CloudConf.TagProjectKey || v.GetTagKey() == "ProjectName" {
				projects = append(projects, &CloudProject{
					CloudProductCommon: CloudProductCommon{
						Name: v.GetTagValue(),
						CID:  fmt.Sprintf("%v-%v", v.GetTagKey(), v.GetTagValue()),
					},
					ProviderID: ProviderID,
					AccountID:  AccountID,
				})
			}
		}
	}
	return projects
}

// DiffCloudProject 云项目输入数据与存量数据差异
func (c *CloudProjectUseCase) DiffCloudProject(ProviderID int, inputs []*CloudProject, where *CloudProjectWhere) (creates []*CloudProject, updates []*CloudProject, deletes []*CloudProject, err error) {
	lists, err := c.QueryCloudProject(where, &CloudProjectOutput{})
	if err != nil {
		return nil, nil, nil, err
	}
	existMap := map[string]CloudProject{}
	for _, find := range lists {
		existMap[fmt.Sprintf("%v-%v", find.CID, find.ProviderID)] = *find
	}
	currentMap := map[string]CloudProject{}
	for _, v := range inputs {
		currentMap[fmt.Sprintf("%v-%v", v.CID, v.ProviderID)] = *v

		found, exist := existMap[fmt.Sprintf("%v-%v", v.CID, v.ProviderID)]
		if exist {
			tempID := found.CloudProductCommon.ID
			found.CloudProductCommon.ID = 0
			found.Provider = nil
			found.Account = nil
			diff := false
			diff, err = DiffNew(found, *v)
			if err != nil {
				continue
			}
			if !diff {
				continue
			}
			updates = append(updates, &CloudProject{
				CloudProductCommon: CloudProductCommon{
					ID:   tempID,
					Name: v.Name,
					CID:  v.CID,
				},
				ProviderID: v.ProviderID,
				AccountID:  v.AccountID,
			})
		} else {
			creates = append(creates, &CloudProject{
				CloudProductCommon: CloudProductCommon{
					Name: v.Name,
					CID:  v.CID,
				},
				ProviderID: v.ProviderID,
				AccountID:  v.AccountID,
			})
		}
	}
	for _, v := range lists {
		if _, exist := currentMap[fmt.Sprintf("%v-%v", v.CID, v.ProviderID)]; !exist {
			//AWS 不删除标签拼接项目（AWS不同区域标签返回不同，导致会异常删除）
			if ProviderID == 5 && strings.Contains(v.CID, configs.Conf.CloudConf.TagProjectKey) {
				continue
			}
			deletes = append(deletes, &CloudProject{
				CloudProductCommon: CloudProductCommon{
					ID: v.ID,
				},
			})
		}
	}
	return
}
