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

	common "github.com/bilibili/HCP/common/models"
)

// CloudServerRepo 云服务器接口
type CloudServerRepo interface {
	// CountCloudServer 查询云服务器数量
	CountCloudServer(where *CloudServerWhere) (total int64, err error)
	// QueryCloudServer 查询云服务器
	QueryCloudServer(where *CloudServerWhere, output *CloudServerOutput) ([]*CloudServer, error)
	// CreateCloudServer 创建云服务器
	CreateCloudServer(create []*CloudServer) ([]*CloudServer, error)
	// UpdateCloudServer 更新云服务器
	UpdateCloudServer(where *CloudServerWhere, update *CloudServer) error
	// DeleteCloudServer 删除云服务器
	DeleteCloudServer(deleteID []int) error
	// UpsertCloudServer 更新或插入云服务器
	UpsertCloudServer(upsert []*CloudServerUpsert) error
}

// CloudServerWhere 云服务器查询条件
type CloudServerWhere struct {
	// 查询语句
	Query string
	// 查询参数
	Arg interface{}
	// 查询条件
	Conditions map[string]interface{}
}

// CloudServerOutput 云服务器查询输出条件参数
type CloudServerOutput struct {
	OutPutCommon
}

// CreateCloudServer 创建云服务器参数结构体
type CreateCloudServer struct {
	// 是否绑定公网 IP
	BindPublicIP string `gorm:"-" json:"bind_public_ip"`
	// 名称前缀
	NamePrefix string `gorm:"-" json:"name_prefix"`
	// 名称后缀
	NameSuffix string `gorm:"-" json:"name_suffix"`
	// RDS类型
	RdsArch string `gorm:"-" json:"rds_arch"`
	// Redis类型
	RedisArch string `gorm:"-" json:"redis_suffix"`
	// 实例数量
	Count int `gorm:"-" json:"count"`
	// 磁盘大小
	DiskSize int `gorm:"-" json:"disk_size"`
	// 是否开启 DNS
	EnableDNS string `gorm:"-" json:"enable_dns"`
	// 带宽
	Bandwidth int `gorm:"-" json:"bandwidth"`
	// 计费类型
	ChargeType *ChargeType `gorm:"-" json:"charge_type"`
	// 续费类型
	RenewalType *ChargeType `gorm:"-" json:"renewal_type"`
	// 磁盘类型
	DiskType *DiskType `gorm:"-" json:"disk_type"`
	// 登录密码
	LoginPassword string `gorm:"-" json:"login_password"`
	// 性能参数（仅阿里云 ECS 需要配置）
	PerformanceLevel string `gorm:"-" json:"performance_level"`
	CloudServer
}

// CloudServer 云服务器参数结构体
type CloudServer struct {
	CloudProductCommon
	// 需要操作的实例id数组
	FormObjects []common.FormObject `gorm:"-" json:"formObjects"`
	// 需要操作的实例id
	FormObject common.FormObject `gorm:"-" json:"formObject"`
	// 要更新的实例名称
	NewName string `gorm:"-" json:"newName"`
	//关机保留时间
	ReservationPeriod string `gorm:"-" json:"reservation_period"`
	// 云账号ID
	AccountID int `gorm:"column:account_id" json:"account_id"`
	// 云账号实体
	Account *Account `json:"account"`
	// 云项目ID
	ProjectID int `gorm:"column:project_id" json:"project_id"`
	// 云项目实体
	Project *CloudProject `json:"project"`
	// 本地项目ID
	ProjectConfigID int `gorm:"column:project_config_id" json:"project_config_id"`
	// 本地项目实体
	ProjectConfig *ProjectConfig `json:"project_config"`
	// 云地域ID
	RegionID int `gorm:"column:region_id" json:"region_id"`
	// 云地域实体
	Region *CloudRegion `json:"region"`
	// 云可用区ID
	ZoneID int `gorm:"column:zone_id" json:"zone_id"`
	// 云可用区实体
	Zone *CloudZone `json:"zone"`
	// 云专有网路ID
	VpcID int `gorm:"column:vpc_id" json:"vpc_id"`
	// 云专有网路实体
	Vpc *CloudVpc `json:"vpc"`
	// 云子网ID
	SubnetID int `gorm:"column:subnet_id" json:"subnet_id"`
	// 云子网实体
	Subnet *CloudSubnet `json:"subnet"`
	// 云子网CID
	SubnetCID string `gorm:"column:subnet_cid" json:"subnet_cid"`
	// 云安全组CID
	SecurityGroupCID string `gorm:"column:security_group_cid" json:"security_group_cid"`
	// 云镜像ID
	ServerImageID int `gorm:"column:server_image_id" json:"server_image_id"`
	// 云镜像实体
	ServerImage *CloudServerImage `json:"server_image"`
	// 云实例规格ID
	ServerSpecID int `gorm:"column:server_spec_id" json:"server_spec_id"`
	// 云实例规格实体
	ServerSpec *CloudServerSpec `json:"server_spec"`
	// 计费方式
	ChargeType string `gorm:"column:charge_type" json:"charge_type"`
	// 续费方式
	RenewStatus string `gorm:"column:renew_status" json:"renew_status"`
	// 内网IP
	PrivateIp string `gorm:"column:private_ip" json:"private_ip"`
	// 外网IP
	PublicIp string `gorm:"column:public_ip" json:"public_ip"`
	// 状态
	Status string `gorm:"column:status" json:"status"`
	// 过期时间
	ExpireTime string `gorm:"column:expire_time" json:"expire_time"`
	// 创建时间
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	// 更新时间
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
	// 云安全组
	SecurityGroup interface{} `gorm:"-" json:"security_group"`
}

// CloudServerUpsert 云服务器更新或插入结构体
type CloudServerUpsert struct {
	CloudServer
}

// CloudServerUseCase 云服务器业务逻辑
type CloudServerUseCase struct {
	repo CloudServerRepo
}

// NewCloudServerUseCase 创建云服务器业务逻辑
func NewCloudServerUseCase(repo CloudServerRepo) *CloudServerUseCase {
	return &CloudServerUseCase{repo: repo}
}

// CountCloudServer 查询云服务器数量
func (c *CloudServerUseCase) CountCloudServer(where *CloudServerWhere) (total int64, err error) {
	return c.repo.CountCloudServer(where)
}

// QueryCloudServer 查询云服务器
func (c *CloudServerUseCase) QueryCloudServer(where *CloudServerWhere, output *CloudServerOutput) ([]*CloudServer, error) {
	return c.repo.QueryCloudServer(where, output)
}

// CreateCloudServer 创建云服务器
func (c *CloudServerUseCase) CreateCloudServer(create []*CloudServer) ([]*CloudServer, error) {
	return c.repo.CreateCloudServer(create)
}

// UpdateCloudServer 更新云服务器
func (c *CloudServerUseCase) UpdateCloudServer(where *CloudServerWhere, update *CloudServer) error {
	return c.repo.UpdateCloudServer(where, update)
}

// DeleteCloudServer 删除云服务器
func (c *CloudServerUseCase) DeleteCloudServer(deleteID []int) error {
	return c.repo.DeleteCloudServer(deleteID)
}

// UpsertCloudServer 更新或插入云服务器
func (c *CloudServerUseCase) UpsertCloudServer(upsert []*CloudServerUpsert) error {
	return c.repo.UpsertCloudServer(upsert)
}

// DiffCloudServer 云服务器输入数据与存量数据差异
func (c *CloudServerUseCase) DiffCloudServer(inputs []*CloudServer, where *CloudServerWhere) (creates []*CloudServer, updates []*CloudServer, deletes []*CloudServer, err error) {
	lists, err := c.QueryCloudServer(where, &CloudServerOutput{})
	if err != nil {
		return nil, nil, nil, err
	}
	existMap := map[string]CloudServer{}
	for _, find := range lists {
		existMap[fmt.Sprintf("%v-%v", find.CID, find.AccountID)] = *find
	}
	currentMap := map[string]CloudServer{}
	for _, v := range inputs {
		currentMap[fmt.Sprintf("%v-%v", v.CID, v.AccountID)] = *v

		found, exist := existMap[fmt.Sprintf("%v-%v", v.CID, v.AccountID)]
		if exist {
			tempID := found.CloudProductCommon.ID
			found.CloudProductCommon.ID = 0
			found.Account = nil
			found.Project = nil
			found.ProjectConfig = nil
			found.Region = nil
			found.Zone = nil
			found.Vpc = nil
			found.Subnet = nil
			found.SecurityGroup = nil
			found.ServerImage = nil
			found.ServerSpec = nil
			found.CreateTime = timeParse("0001-01-01 00:00:00")
			found.UpdateTime = timeParse("0001-01-01 00:00:00")
			diff := false
			diff, err = DiffNew(found, *v)
			if err != nil {
				continue
			}
			if !diff {
				continue
			}
			updates = append(updates, &CloudServer{
				CloudProductCommon: CloudProductCommon{
					ID:   tempID,
					Name: v.Name,
					CID:  v.CID,
				},
				AccountID:        v.AccountID,
				ProjectConfigID:  v.ProjectConfigID,
				ProjectID:        v.ProjectID,
				RegionID:         v.RegionID,
				ZoneID:           v.ZoneID,
				VpcID:            v.VpcID,
				SubnetID:         v.SubnetID,
				SubnetCID:        v.SubnetCID,
				SecurityGroupCID: v.SecurityGroupCID,
				ServerImageID:    v.ServerImageID,
				ServerSpecID:     v.ServerSpecID,
				ChargeType:       v.ChargeType,
				RenewStatus:      v.RenewStatus,
				PrivateIp:        v.PrivateIp,
				PublicIp:         v.PublicIp,
				Status:           v.Status,
				ExpireTime:       v.ExpireTime,
			})
		} else {
			creates = append(creates, &CloudServer{
				CloudProductCommon: CloudProductCommon{
					Name: v.Name,
					CID:  v.CID,
				},
				AccountID:        v.AccountID,
				ProjectID:        v.ProjectID,
				ProjectConfigID:  v.ProjectConfigID,
				RegionID:         v.RegionID,
				ZoneID:           v.ZoneID,
				VpcID:            v.VpcID,
				SubnetID:         v.SubnetID,
				SubnetCID:        v.SubnetCID,
				SecurityGroupCID: v.SecurityGroupCID,
				ServerImageID:    v.ServerImageID,
				ServerSpecID:     v.ServerSpecID,
				ChargeType:       v.ChargeType,
				RenewStatus:      v.RenewStatus,
				PrivateIp:        v.PrivateIp,
				PublicIp:         v.PublicIp,
				Status:           v.Status,
				ExpireTime:       v.ExpireTime,
			})
		}
	}
	for _, v := range lists {
		if _, exist := currentMap[fmt.Sprintf("%v-%v", v.CID, v.AccountID)]; !exist {
			deletes = append(deletes, &CloudServer{
				CloudProductCommon: CloudProductCommon{
					ID: v.ID,
				},
			})
		}
	}
	return
}
