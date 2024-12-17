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
	"time"
)

// CloudSecurityGroupRepo 云安全组接口
type CloudSecurityGroupRepo interface {
	// QueryCloudSecurityGroup 查询云安全组
	QueryCloudSecurityGroup(where *CloudSecurityGroupWhere, output *CloudSecurityGroupOutput) ([]*CloudSecurityGroup, error)
	// CountCloudSecurityGroup 查询云安全组数量
	CountCloudSecurityGroup(where *CloudSecurityGroupWhere) (int64, error)
	// CreateCloudSecurityGroup 创建云安全组
	CreateCloudSecurityGroup(create []*CloudSecurityGroup) ([]*CloudSecurityGroup, error)
	// UpdateCloudSecurityGroup 更新云安全组
	UpdateCloudSecurityGroup(where *CloudSecurityGroupWhere, update *CloudSecurityGroup) error
	// DeleteCloudSecurityGroup 删除云安全组
	DeleteCloudSecurityGroup(deleteID []int) error
	// UpsertCloudSecurityGroup 更新云安全组
	UpsertCloudSecurityGroup(upsert []*CloudSecurityGroupUpsert) error
}

// CloudSecurityGroupWhere 云安全组查询条件
type CloudSecurityGroupWhere struct {
	// 查询语句
	Query string
	// 查询参数
	Arg interface{}
	// 查询条件
	Conditions map[string]interface{}
}

// CloudSecurityGroupOutput 云安全组查询输出条件参数
type CloudSecurityGroupOutput struct {
	OutPutCommon
}

// CloudSecurityGroup 云安全组参数结构体
type CloudSecurityGroup struct {
	CloudProductCommon
	// 云项目ID
	ProjectID int `gorm:"column:project_id;" json:"project_id"`
	// 云项目实体
	Project *CloudProject `json:"project"`
	// 云地域ID
	RegionID int `gorm:"column:region_id" json:"region_id"`
	// 云地域实体
	Region *CloudRegion `json:"region"`
	// 云账号ID
	AccountID int `gorm:"column:account_id" json:"account_id"`
	// 云账号实体
	Account *Account `json:"account"`
	// 云专有网络ID
	VpcID int `gorm:"column:vpc_id" json:"vpc_id"`
	// 云专有网络实体
	Vpc *CloudVpc `json:"vpc"`
	// 本地项目云云账号关联关系实体
	ProjectAccountConfig *ProjectAccountConfig `gorm:"foreignKey:ProjectID;references:ProjectID" gorm:"foreignKey:AccountID;references:AccountID" json:"project_account_config"`
	// 本地项目实体
	ProjectConfig *ProjectConfig `gorm:"-" json:"project_config"`
	// 状态
	Status string `gorm:"column:status" json:"status"`
	// 创建时间
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	// 更新时间
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

// CloudSecurityGroupUpsert 云安全组更新或插入结构体
type CloudSecurityGroupUpsert struct {
	CloudSecurityGroup
}

// CloudSecurityGroupUseCase 云安全组业务逻辑
type CloudSecurityGroupUseCase struct {
	repo CloudSecurityGroupRepo
}

// NewCloudSecurityGroupUseCase 创建云安全组业务逻辑
func NewCloudSecurityGroupUseCase(repo CloudSecurityGroupRepo) *CloudSecurityGroupUseCase {
	return &CloudSecurityGroupUseCase{repo: repo}
}

// QueryCloudSecurityGroup 查询云安全组
func (c *CloudSecurityGroupUseCase) QueryCloudSecurityGroup(where *CloudSecurityGroupWhere, output *CloudSecurityGroupOutput) ([]*CloudSecurityGroup, error) {
	return c.repo.QueryCloudSecurityGroup(where, output)
}

// CountCloudSecurityGroup 查询云安全组数量
func (c *CloudSecurityGroupUseCase) CountCloudSecurityGroup(where *CloudSecurityGroupWhere) (int64, error) {
	return c.repo.CountCloudSecurityGroup(where)
}

// CreateCloudSecurityGroup 创建云安全组
func (c *CloudSecurityGroupUseCase) CreateCloudSecurityGroup(create []*CloudSecurityGroup) ([]*CloudSecurityGroup, error) {
	return c.repo.CreateCloudSecurityGroup(create)
}

// UpdateCloudSecurityGroup 更新云安全组
func (c *CloudSecurityGroupUseCase) UpdateCloudSecurityGroup(where *CloudSecurityGroupWhere, update *CloudSecurityGroup) error {
	return c.repo.UpdateCloudSecurityGroup(where, update)
}

// DeleteCloudSecurityGroup 删除云安全组
func (c *CloudSecurityGroupUseCase) DeleteCloudSecurityGroup(deleteID []int) error {
	return c.repo.DeleteCloudSecurityGroup(deleteID)
}

// UpsertCloudSecurityGroup 更新或插入云安全组
func (c *CloudSecurityGroupUseCase) UpsertCloudSecurityGroup(upsert []*CloudSecurityGroupUpsert) error {
	return c.repo.UpsertCloudSecurityGroup(upsert)
}

// DiffCloudSecurityGroup 云安全组输入数据与存量数据差异
func (c *CloudSecurityGroupUseCase) DiffCloudSecurityGroup(inputs []*CloudSecurityGroup, where *CloudSecurityGroupWhere) (creates []*CloudSecurityGroup, updates []*CloudSecurityGroup, deletes []*CloudSecurityGroup, err error) {
	lists, err := c.QueryCloudSecurityGroup(where, &CloudSecurityGroupOutput{})
	if err != nil {
		return nil, nil, nil, err
	}
	existsMap := map[string]CloudSecurityGroup{}
	for _, find := range lists {
		existsMap[fmt.Sprintf("%v-%v", find.CID, find.AccountID)] = *find
	}
	currentMap := map[string]CloudSecurityGroup{}
	for _, v := range inputs {
		currentMap[fmt.Sprintf("%v-%v", v.CID, v.AccountID)] = *v

		found, exist := existsMap[fmt.Sprintf("%v-%v", v.CID, v.AccountID)]
		if exist {
			tempID := found.CloudProductCommon.ID
			found.CloudProductCommon.ID = 0
			found.Region = nil
			found.Account = nil
			found.Project = nil
			found.Vpc = nil
			diff := false
			diff, err = DiffNew(found, *v)
			if err != nil {
				continue
			}
			if !diff {
				continue
			}
			updates = append(updates, &CloudSecurityGroup{
				CloudProductCommon: CloudProductCommon{
					ID:   tempID,
					Name: v.Name,
					CID:  v.CID,
				},
				ProjectID: v.ProjectID,
				RegionID:  v.RegionID,
				AccountID: v.AccountID,
				Status:    v.Status,
				VpcID:     v.VpcID,
			})
		} else {
			creates = append(creates, &CloudSecurityGroup{
				CloudProductCommon: CloudProductCommon{
					Name: v.Name,
					CID:  v.CID,
				},
				ProjectID: v.ProjectID,
				RegionID:  v.RegionID,
				AccountID: v.AccountID,
				Status:    v.Status,
				VpcID:     v.VpcID,
			})
		}
	}
	for _, v := range lists {
		if _, exist := currentMap[fmt.Sprintf("%v-%v", v.CID, v.AccountID)]; !exist {
			deletes = append(deletes, &CloudSecurityGroup{
				CloudProductCommon: CloudProductCommon{ID: v.ID},
			})
		}
	}
	return
}
