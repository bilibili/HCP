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

import "fmt"

// CloudVpcRepo 云专有网路接口
type CloudVpcRepo interface {
	// QueryCloudVpc 查询云专有网路
	QueryCloudVpc(where *CloudVpcWhere, output *CloudVpcOutput) ([]*CloudVpc, error)
	// CountCloudVpc 查询云专有网路数量
	CountCloudVpc(where *CloudVpcWhere) (int64, error)
	// CreateCloudVpc 创建云专有网路
	CreateCloudVpc(create []*CloudVpc) ([]*CloudVpc, error)
	// UpdateCloudVpc 更新云专有网路
	UpdateCloudVpc(where *CloudVpcWhere, update *CloudVpc) error
	// DeleteCloudVpc 删除云专有网路
	DeleteCloudVpc(deleteID []int) error
	// UpsertCloudVpc 更新或插入云专有网路
	UpsertCloudVpc(upsert []*CloudVpcUpsert) error
}

// CloudVpcWhere 云专有网路查询条件
type CloudVpcWhere struct {
	// 查询语句
	Query string
	// 查询参数
	Arg interface{}
	// 查询条件
	Conditions map[string]interface{}
}

// CloudVpcOutput 云专有网路查询输出条件参数
type CloudVpcOutput struct {
	OutPutCommon
}

// CloudVpc 云专有网路参数结构体
type CloudVpc struct {
	CloudProductCommon
	// 云项目ID
	ProjectID int `gorm:"column:project_id;" json:"project_id"`
	// 云项目实体
	Project *CloudProject `json:"project"`
	// 云专有网路Cidr
	Cidr string `gorm:"column:cidr" json:"cidr"`
	// 状态
	Status string `gorm:"column:status" json:"status"`
	// 云地域ID
	RegionID int `gorm:"column:region_id" json:"region_id"`
	// 云地域实体
	Region *CloudRegion `json:"region"`
	// 云账号ID
	AccountID int `gorm:"column:account_id" json:"account_id"`
	// 云账号实体
	Account *Account `json:"account"`
	// 云账号和项目关联关系实体
	ProjectAccountConfig *ProjectAccountConfig `gorm:"foreignKey:ProjectID;references:ProjectID" gorm:"foreignKey:AccountID;references:AccountID" json:"project_account_config"`
	// 本地项目实体
	ProjectConfig *ProjectConfig `gorm:"-" json:"project_config"`
}

// CloudVpcUpsert 云专有网路更新或插入结构体
type CloudVpcUpsert struct {
	CloudVpc
}

// CloudVpcUseCase 云专有网路业务逻辑
type CloudVpcUseCase struct {
	repo CloudVpcRepo
}

// NewCloudVpcUseCase 创建云专有网路业务逻辑
func NewCloudVpcUseCase(repo CloudVpcRepo) *CloudVpcUseCase {
	return &CloudVpcUseCase{repo: repo}
}

// QueryCloudVpc 查询云专有网路
func (c *CloudVpcUseCase) QueryCloudVpc(where *CloudVpcWhere, output *CloudVpcOutput) ([]*CloudVpc, error) {
	return c.repo.QueryCloudVpc(where, output)
}

// CountCloudVpc 查询云专有网路数量
func (c *CloudVpcUseCase) CountCloudVpc(where *CloudVpcWhere) (int64, error) {
	return c.repo.CountCloudVpc(where)
}

// CreateCloudVpc 创建云专有网路
func (c *CloudVpcUseCase) CreateCloudVpc(create []*CloudVpc) ([]*CloudVpc, error) {
	return c.repo.CreateCloudVpc(create)
}

// UpdateCloudVpc 更新云专有网路
func (c *CloudVpcUseCase) UpdateCloudVpc(where *CloudVpcWhere, update *CloudVpc) error {
	return c.repo.UpdateCloudVpc(where, update)
}

// DeleteCloudVpc 删除云专有网路
func (c *CloudVpcUseCase) DeleteCloudVpc(deleteID []int) error {
	return c.repo.DeleteCloudVpc(deleteID)
}

// UpsertCloudVpc 更新或插入云专有网路
func (c *CloudVpcUseCase) UpsertCloudVpc(upsert []*CloudVpcUpsert) error {
	return c.repo.UpsertCloudVpc(upsert)
}

// DiffCloudVpc 云专有网路输入数据与存量数据差异
func (c *CloudVpcUseCase) DiffCloudVpc(inputs []*CloudVpc, where *CloudVpcWhere) (creates []*CloudVpc, updates []*CloudVpc, deletes []*CloudVpc, err error) {
	lists, err := c.QueryCloudVpc(where, &CloudVpcOutput{})
	if err != nil {
		return nil, nil, nil, err
	}
	existMap := map[string]CloudVpc{}
	for _, find := range lists {
		existMap[fmt.Sprintf("%v-%v", find.CID, find.AccountID)] = *find
	}
	currentMap := map[string]CloudVpc{}
	for _, v := range inputs {
		currentMap[fmt.Sprintf("%v-%v", v.CID, v.AccountID)] = *v

		found, exist := existMap[fmt.Sprintf("%v-%v", v.CID, v.AccountID)]
		if exist {
			tempID := found.CloudProductCommon.ID
			found.CloudProductCommon.ID = 0
			found.Project = nil
			found.Region = nil
			found.Account = nil
			diff := false
			diff, err = DiffNew(found, *v)
			if err != nil {
				continue
			}
			if !diff {
				continue
			}
			updates = append(updates, &CloudVpc{
				CloudProductCommon: CloudProductCommon{
					ID:   tempID,
					Name: v.Name,
					CID:  v.CID,
				},
				ProjectID: v.ProjectID,
				Cidr:      v.Cidr,
				Status:    v.Status,
				RegionID:  v.RegionID,
				AccountID: v.AccountID,
			})
		} else {
			creates = append(creates, &CloudVpc{
				CloudProductCommon: CloudProductCommon{
					Name: v.Name,
					CID:  v.CID,
				},
				ProjectID: v.ProjectID,
				Cidr:      v.Cidr,
				Status:    v.Status,
				RegionID:  v.RegionID,
				AccountID: v.AccountID,
			})
		}
	}
	for _, v := range lists {
		if _, exist := currentMap[fmt.Sprintf("%v-%v", v.CID, v.AccountID)]; !exist {
			deletes = append(deletes, &CloudVpc{
				CloudProductCommon: CloudProductCommon{ID: v.ID},
			})
		}
	}
	return
}
