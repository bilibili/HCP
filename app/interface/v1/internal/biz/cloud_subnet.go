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

// CloudSubnetRepo 云子网接口
type CloudSubnetRepo interface {
	// QueryCloudSubnet 查询云子网
	QueryCloudSubnet(where *CloudSubnetWhere, output *CloudSubnetOutput) ([]*CloudSubnet, error)
	// CountCloudSubnet 查询云子网数量
	CountCloudSubnet(where *CloudSubnetWhere) (int64, error)
	// CreateCloudSubnet 创建云子网
	CreateCloudSubnet(create []*CloudSubnet) ([]*CloudSubnet, error)
	// UpdateCloudSubnet 更新云子网
	UpdateCloudSubnet(where *CloudSubnetWhere, update *CloudSubnet) error
	// DeleteCloudSubnet 删除云子网
	DeleteCloudSubnet(deleteID []int) error
	// UpsertCloudSubnet 更新或插入云子网
	UpsertCloudSubnet(upsert []*CloudSubnetUpsert) error
}

// CloudSubnetWhere 云子网查询条件
type CloudSubnetWhere struct {
	// 查询语句
	Query string
	// 查询参数
	Arg interface{}
	// 查询条件
	Conditions map[string]interface{}
}

// CloudSubnetOutput 云子网查询输出条件参数
type CloudSubnetOutput struct {
	OutPutCommon
}

// CloudSubnet 云子网参数结构体
type CloudSubnet struct {
	CloudProductCommon
	// 云项目ID
	ProjectID int `gorm:"column:project_id;" json:"project_id"`
	// 云项目实体
	Project *CloudProject `json:"project"`
	// 云子网CIDR
	Cidr string `gorm:"column:cidr" json:"cidr"`
	// 云可用区ID
	ZoneID int `gorm:"column:zone_id" json:"zone_id"`
	// 云可用区实体
	Zone *CloudZone `json:"zone"`
	// 云专有网路ID
	VpcID int `gorm:"column:vpc_id" json:"vpc_id"`
	// 云专有网路实体
	Vpc *CloudVpc `json:"vpc"`
	// 云账号ID
	AccountID int `gorm:"column:account_id" json:"account_id"`
	// 云账号实体
	Account *Account `json:"account"`
	// 云地域ID
	RegionID int `gorm:"column:region_id" json:"region_id"`
	// 云地域实体
	Region *CloudRegion `json:"region"`
	// 状态
	Status string `gorm:"column:status" json:"status"`
	//云账号和本地账号关联关系实体
	ProjectAccountConfig *ProjectAccountConfig `gorm:"foreignKey:ProjectID;references:ProjectID" gorm:"foreignKey:AccountID;references:AccountID" json:"project_account_config"`
	// 本地项目实体
	ProjectConfig *ProjectConfig `gorm:"-" json:"project_config"`
}

// CloudSubnetUpsert 云子网更新或插入结构体
type CloudSubnetUpsert struct {
	CloudSubnet
}

// CloudSubnetUseCase 云子网业务逻辑
type CloudSubnetUseCase struct {
	repo CloudSubnetRepo
}

// NewCloudSubnetUseCase 创建云子网业务逻辑
func NewCloudSubnetUseCase(repo CloudSubnetRepo) *CloudSubnetUseCase {
	return &CloudSubnetUseCase{repo: repo}
}

// QueryCloudSubnet 查询云子网
func (c *CloudSubnetUseCase) QueryCloudSubnet(where *CloudSubnetWhere, output *CloudSubnetOutput) ([]*CloudSubnet, error) {
	return c.repo.QueryCloudSubnet(where, output)
}

// CountCloudSubnet 查询云子网数量
func (c *CloudSubnetUseCase) CountCloudSubnet(where *CloudSubnetWhere) (int64, error) {
	return c.repo.CountCloudSubnet(where)
}

// CreateCloudSubnet 创建云子网
func (c *CloudSubnetUseCase) CreateCloudSubnet(create []*CloudSubnet) ([]*CloudSubnet, error) {
	return c.repo.CreateCloudSubnet(create)
}

// UpdateCloudSubnet 更新云子网
func (c *CloudSubnetUseCase) UpdateCloudSubnet(where *CloudSubnetWhere, update *CloudSubnet) error {
	return c.repo.UpdateCloudSubnet(where, update)
}

// DeleteCloudSubnet 删除云子网
func (c *CloudSubnetUseCase) DeleteCloudSubnet(deleteID []int) error {
	return c.repo.DeleteCloudSubnet(deleteID)
}

// UpsertCloudSubnet 更新或插入云子网
func (c *CloudSubnetUseCase) UpsertCloudSubnet(upsert []*CloudSubnetUpsert) error {
	return c.repo.UpsertCloudSubnet(upsert)
}

// DiffCloudSubnet 云子网输入数据与存量数据差异
func (c *CloudSubnetUseCase) DiffCloudSubnet(inputs []*CloudSubnet, where *CloudSubnetWhere) (creates []*CloudSubnet, updates []*CloudSubnet, deletes []*CloudSubnet, err error) {
	lists, err := c.QueryCloudSubnet(where, &CloudSubnetOutput{})
	if err != nil {
		return nil, nil, nil, err
	}
	existMap := map[string]CloudSubnet{}
	for _, find := range lists {
		existMap[fmt.Sprintf("%v-%v-%v", find.CID, find.AccountID, find.RegionID)] = *find
	}
	currentMap := map[string]CloudSubnet{}
	for _, v := range inputs {
		currentMap[fmt.Sprintf("%v-%v-%v", v.CID, v.AccountID, v.RegionID)] = *v

		found, exist := existMap[fmt.Sprintf("%v-%v-%v", v.CID, v.AccountID, v.RegionID)]
		if exist {
			tempID := found.CloudProductCommon.ID
			found.CloudProductCommon.ID = 0
			found.Project = nil
			found.Zone = nil
			found.Vpc = nil
			found.Account = nil
			found.ProjectAccountConfig = nil
			found.ProjectConfig = nil
			found.Region = nil
			diff := false
			diff, err = DiffNew(found, *v)
			if err != nil {
				continue
			}
			if !diff {
				continue
			}
			updates = append(updates, &CloudSubnet{
				CloudProductCommon: CloudProductCommon{
					ID:   tempID,
					Name: v.Name,
					CID:  v.CID,
				},
				ProjectID: v.ProjectID,
				Cidr:      v.Cidr,
				ZoneID:    v.ZoneID,
				VpcID:     v.VpcID,
				AccountID: v.AccountID,
				RegionID:  v.RegionID,
				Status:    v.Status,
			})
		} else {
			creates = append(creates, &CloudSubnet{
				CloudProductCommon: CloudProductCommon{
					Name: v.Name,
					CID:  v.CID,
				},
				ProjectID: v.ProjectID,
				Cidr:      v.Cidr,
				ZoneID:    v.ZoneID,
				VpcID:     v.VpcID,
				AccountID: v.AccountID,
				RegionID:  v.RegionID,
				Status:    v.Status,
			})
		}
	}
	for _, v := range lists {
		if _, exist := currentMap[fmt.Sprintf("%v-%v-%v", v.CID, v.AccountID, v.RegionID)]; !exist {
			deletes = append(deletes, &CloudSubnet{
				CloudProductCommon: CloudProductCommon{ID: v.ID},
			})
		}
	}
	return
}
