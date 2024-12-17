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

// CloudServerImageRepo 云服务器镜像接口
type CloudServerImageRepo interface {
	// QueryCloudServerImage 查询云服务器镜像
	QueryCloudServerImage(where *CloudServerImageWhere, output *CloudServerImageOutput) ([]*CloudServerImage, error)
	// CountCloudServerImage 查询云服务器镜像数量
	CountCloudServerImage(where *CloudServerImageWhere) (int64, error)
	// CreateCloudServerImage 创建云服务器镜像
	CreateCloudServerImage(create []*CloudServerImage) ([]*CloudServerImage, error)
	// CreatBatchesCloudServerImage 批量创建云服务器镜像
	CreatBatchesCloudServerImage(create []*CloudServerImage) ([]*CloudServerImage, error)
	// UpdateCloudServerImage 更新云服务器镜像
	UpdateCloudServerImage(where *CloudServerImageWhere, update *CloudServerImage) error
	// DeleteCloudServerImage 删除云服务器镜像
	DeleteCloudServerImage(deleteID []int) error
	// UpsertCloudServerImage 更新或插入云服务器镜像
	UpsertCloudServerImage(upsert []*CloudServerImageUpsert) error
}

// CloudServerImageWhere 云服务器镜像查询条件
type CloudServerImageWhere struct {
	// 查询语句
	Query string
	// 查询参数
	Arg interface{}
	// 查询条件
	Conditions map[string]interface{}
}

// CloudServerImageOutput 云服务器镜像查询输出条件参数
type CloudServerImageOutput struct {
	OutPutCommon
}

// CloudServerImage 云服务器镜像参数结构体
type CloudServerImage struct {
	CloudProductCommon
	// 云地域ID
	RegionID int `gorm:"region_id" json:"region_id"`
	// 云地域实体
	Region *CloudRegion `json:"region"`
	// 云账号ID
	AccountID int `gorm:"account_id" json:"account_id"`
	// 云账号实体
	Account *Account `json:"account"`
	// 系统名称
	OsName string `gorm:"os_name" json:"os_name"`
	// 系统类型
	Type string `gorm:"type" json:"type"`
	// 状态
	Status string `gorm:"status" json:"status"`
}

// CloudServerImageUpsert 云服务器镜像更新或插入结构体
type CloudServerImageUpsert struct {
	CloudServerImage
}

// CloudServerImageUseCase 云服务器镜像业务逻辑接口
type CloudServerImageUseCase struct {
	repo CloudServerImageRepo
}

// NewCloudServerImageUseCase 创建云服务器镜像业务逻辑接口
func NewCloudServerImageUseCase(repo CloudServerImageRepo) *CloudServerImageUseCase {
	return &CloudServerImageUseCase{repo: repo}
}

// QueryCloudServerImage 查询云服务器镜像
func (c *CloudServerImageUseCase) QueryCloudServerImage(where *CloudServerImageWhere, output *CloudServerImageOutput) ([]*CloudServerImage, error) {
	return c.repo.QueryCloudServerImage(where, output)
}

// CountCloudServerImage 查询云服务器镜像数量
func (c *CloudServerImageUseCase) CountCloudServerImage(where *CloudServerImageWhere) (int64, error) {
	return c.repo.CountCloudServerImage(where)
}

// CreateCloudServerImage 创建云服务器镜像
func (c *CloudServerImageUseCase) CreateCloudServerImage(create []*CloudServerImage) ([]*CloudServerImage, error) {
	return c.repo.CreateCloudServerImage(create)
}

// CreatBatchesCloudServerImage 批量创建云服务器镜像
func (c *CloudServerImageUseCase) CreatBatchesCloudServerImage(create []*CloudServerImage) ([]*CloudServerImage, error) {
	return c.repo.CreatBatchesCloudServerImage(create)
}

// UpdateCloudServerImage 更新云服务器镜像
func (c *CloudServerImageUseCase) UpdateCloudServerImage(where *CloudServerImageWhere, update *CloudServerImage) error {
	return c.repo.UpdateCloudServerImage(where, update)
}

// DeleteCloudServerImage 删除云服务器镜像
func (c *CloudServerImageUseCase) DeleteCloudServerImage(deleteID []int) error {
	return c.repo.DeleteCloudServerImage(deleteID)
}

// UpsertCloudServerImage 更新或插入云服务器镜像
func (c *CloudServerImageUseCase) UpsertCloudServerImage(upsert []*CloudServerImageUpsert) error {
	return c.repo.UpsertCloudServerImage(upsert)
}

// DiffCloudServerImage 云服务器镜像输入数据与存量数据差异
func (c *CloudServerImageUseCase) DiffCloudServerImage(inputs []*CloudServerImage, where *CloudServerImageWhere) (creates []*CloudServerImage, updates []*CloudServerImage, deletes []*CloudServerImage, err error) {
	lists, err := c.QueryCloudServerImage(where, &CloudServerImageOutput{})
	if err != nil {
		return nil, nil, nil, err
	}
	existMap := map[string]CloudServerImage{}
	for _, find := range lists {
		existMap[fmt.Sprintf("%v-%v-%v", find.CID, find.AccountID, find.RegionID)] = *find
	}
	currentMap := map[string]CloudServerImage{}
	for _, v := range inputs {
		currentMap[fmt.Sprintf("%v-%v-%v", v.CID, v.AccountID, v.RegionID)] = *v

		found, exist := existMap[fmt.Sprintf("%v-%v-%v", v.CID, v.AccountID, v.RegionID)]
		if exist {
			tempID := found.CloudProductCommon.ID
			found.CloudProductCommon.ID = 0
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
			updates = append(updates, &CloudServerImage{
				CloudProductCommon: CloudProductCommon{
					ID:   tempID,
					Name: v.Name,
					CID:  v.CID,
				},
				RegionID:  v.RegionID,
				AccountID: v.AccountID,
				OsName:    v.OsName,
				Type:      v.Type,
				Status:    v.Status,
			})
		} else {
			creates = append(creates, &CloudServerImage{
				CloudProductCommon: CloudProductCommon{
					Name: v.Name,
					CID:  v.CID,
				},
				RegionID:  v.RegionID,
				AccountID: v.AccountID,
				OsName:    v.OsName,
				Type:      v.Type,
				Status:    v.Status,
			})
		}
	}
	for _, v := range lists {
		if _, exist := currentMap[fmt.Sprintf("%v-%v-%v", v.CID, v.AccountID, v.RegionID)]; !exist {
			deletes = append(deletes, &CloudServerImage{
				CloudProductCommon: CloudProductCommon{ID: v.ID},
			})
		}
	}
	return
}
