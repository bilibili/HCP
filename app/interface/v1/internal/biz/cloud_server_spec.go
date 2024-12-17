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

// CloudServerSpecRepo 云服务器规格接口
type CloudServerSpecRepo interface {
	QueryCloudServerSpec(where *CloudServerSpecWhere, output *CloudServerSpecOutput) ([]*CloudServerSpec, error)
	CountCloudServerSpec(where *CloudServerSpecWhere) (int64, error)
	CreateCloudServerSpec(create []*CloudServerSpec) ([]*CloudServerSpec, error)
	UpdateCloudServerSpec(where *CloudServerSpecWhere, update *CloudServerSpec) error
	DeleteCloudServerSpec(deleteID []int) error
	UpsertCloudServerSpec(upsert []*CloudServerSpecUpsert) error
}

// CloudServerSpecWhere 云服务器规格查询条件
type CloudServerSpecWhere struct {
	// 查询语句
	Query string
	// 查询参数
	Arg interface{}
	// 查询条件
	Conditions map[string]interface{}
}

// CloudServerSpecOutput 云服务器规格查询输出条件参数
type CloudServerSpecOutput struct {
	OutPutCommon
}

// CloudServerSpec 云服务器规格参数结构体
type CloudServerSpec struct {
	CloudProductCommon
	// 云厂商ID
	ProviderID int `gorm:"column:provider_id" json:"provider_id"`
	// 云厂商实体
	Provider *Provider `json:"provider"`
	// 云账号ID
	AccountID int `gorm:"column:account_id" json:"account_id"`
	// 云账号实体
	Account *Account `json:"account"`
	// 内网宽带
	Bandwidth float64 `gorm:"column:bandwidth" json:"bandwidth"`
	// 规格类型
	Category string `gorm:"column:category" json:"category"`
	// CPU数量
	CPU int `gorm:"column:cpu" json:"cpu"`
	// 规格族
	Family string `gorm:"column:family" json:"family"`
	// GPU数量
	GPU int `gorm:"column:gpu" json:"gpu"`
	// GPU型号
	GPUModel string `gorm:"column:gpu_model" json:"gpu_model"`
	// 内存大小
	Memory int `gorm:"column:memory" json:"memory"`
	// 内网收发包PPS
	PPS float64 `gorm:"column:pps" json:"pps"`
	// 状态
	Status string `gorm:"column:status" json:"status"`
}

// CloudServerSpecUpsert 云服务器规格更新或插入结构体
type CloudServerSpecUpsert struct {
	CloudServerSpec
}

// CloudServerSpecUseCase 云服务器规格业务逻辑接口
type CloudServerSpecUseCase struct {
	repo CloudServerSpecRepo
}

// NewCloudServerSpecUseCase 云服务器规格业务逻辑接口实现
func NewCloudServerSpecUseCase(repo CloudServerSpecRepo) *CloudServerSpecUseCase {
	return &CloudServerSpecUseCase{repo: repo}
}

// QueryCloudServerSpec 查询云服务器规格
func (c *CloudServerSpecUseCase) QueryCloudServerSpec(where *CloudServerSpecWhere, output *CloudServerSpecOutput) ([]*CloudServerSpec, error) {
	return c.repo.QueryCloudServerSpec(where, output)
}

// CountCloudServerSpec 查询云服务器规格数量
func (c *CloudServerSpecUseCase) CountCloudServerSpec(where *CloudServerSpecWhere) (int64, error) {
	return c.repo.CountCloudServerSpec(where)
}

// CreateCloudServerSpec 创建云服务器规格
func (c *CloudServerSpecUseCase) CreateCloudServerSpec(create []*CloudServerSpec) ([]*CloudServerSpec, error) {
	return c.repo.CreateCloudServerSpec(create)
}

// UpdateCloudServerSpec 更新云服务器规格
func (c *CloudServerSpecUseCase) UpdateCloudServerSpec(where *CloudServerSpecWhere, update *CloudServerSpec) error {
	return c.repo.UpdateCloudServerSpec(where, update)
}

// DeleteCloudServerSpec 删除云服务器规格
func (c *CloudServerSpecUseCase) DeleteCloudServerSpec(deleteID []int) error {
	return c.repo.DeleteCloudServerSpec(deleteID)
}

// UpsertCloudServerSpec 更新或插入云服务器规格
func (c *CloudServerSpecUseCase) UpsertCloudServerSpec(upsert []*CloudServerSpecUpsert) error {
	return c.repo.UpsertCloudServerSpec(upsert)
}

// DiffCloudServerSpec 云服务器规格输入数据与存量数据差异
func (c *CloudServerSpecUseCase) DiffCloudServerSpec(inputs []*CloudServerSpec, where *CloudServerSpecWhere) (creates []*CloudServerSpec, updates []*CloudServerSpec, deletes []*CloudServerSpec, err error) {
	lists, err := c.QueryCloudServerSpec(where, &CloudServerSpecOutput{})
	if err != nil {
		return nil, nil, nil, err
	}
	existMap := map[string]CloudServerSpec{}
	for _, find := range lists {
		existMap[fmt.Sprintf("%v-%v-%v", find.CID, find.AccountID, find.Name)] = *find
	}
	currentMap := map[string]CloudServerSpec{}
	for _, v := range inputs {
		currentMap[fmt.Sprintf("%v-%v-%v", v.CID, v.AccountID, v.Name)] = *v

		found, exist := existMap[fmt.Sprintf("%v-%v-%v", v.CID, v.AccountID, v.Name)]
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
			updates = append(updates, &CloudServerSpec{
				CloudProductCommon: CloudProductCommon{
					ID:   tempID,
					Name: v.Name,
					CID:  v.CID,
				},
				ProviderID: v.ProviderID,
				AccountID:  v.AccountID,
				Bandwidth:  v.Bandwidth,
				Category:   v.Category,
				CPU:        v.CPU,
				Family:     v.Family,
				GPU:        v.GPU,
				GPUModel:   v.GPUModel,
				Memory:     v.Memory,
				PPS:        v.PPS,
				Status:     v.Status,
			})
		} else {
			creates = append(creates, &CloudServerSpec{
				CloudProductCommon: CloudProductCommon{
					Name: v.Name,
					CID:  v.CID,
				},
				ProviderID: v.ProviderID,
				AccountID:  v.AccountID,
				Bandwidth:  v.Bandwidth,
				Category:   v.Category,
				CPU:        v.CPU,
				Family:     v.Family,
				GPU:        v.GPU,
				GPUModel:   v.GPUModel,
				Memory:     v.Memory,
				PPS:        v.PPS,
				Status:     v.Status,
			})
		}
	}
	return
}
