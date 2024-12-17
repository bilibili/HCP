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

package common

const (
	TypeImage                   = "type_image"
	PayType                     = "pay_type"
	CloudPayResourceRenewStatus = "cloud_renew_status"

	CloudServerSpecGeneral = "通用型"
	CloudServerSpecCompute = "计算型"
	CloudServerSpecMemory  = "内存型"
	CloudDBStatusPending   = "pending"
	CloudDBStatusStopped   = "stopped"
	CloudDBStatusStopping  = "stopping"
	CloudDBStatusRunning   = "running"
	CloudPayTypeByFlow     = "postpaid"
	CloudPayTypeByPeriod   = "prepaid"

	CloudRegion         = "cloud_region"
	CloudRegionLocation = "cloud_region_location"
	CloudZone           = "cloud_zone"
	CloudProject        = "cloud_project"
	CloudSecurityGroup  = "cloud_security_group"
	CloudVpc            = "cloud_vpc"
	CloudSubnet         = "cloud_subnet"
	CloudServerImage    = "cloud_server_image"
	CloudServer         = "cloud_server"
	CloudServerSpec     = "cloud_server_spec"
)

var (
	SyncCloudProducts = []string{
		CloudRegion,
		CloudServerImage,
		CloudServer,
		CloudZone,
	}

	SyncCloudProductsQuick = []string{
		CloudProject,
		CloudVpc,
		CloudSubnet,
		CloudSecurityGroup,
	}

	SyncCloudProductsSlow = []string{
		CloudServerSpec,
	}
)

// NewCloudClientRequest 云客户端请求
type NewCloudClientRequest struct {
	// 云地域
	Region string `json:"region"`
	// 密钥ID
	SecretId string `json:"secret_id"`
	// 密钥
	SecretKey string `json:"secret_key"`
	// AK
	SecretAk string `json:"secret_ak"`
	// SK
	SecretSk string `json:"secret_sk"`
	// 公钥
	PublicKey string `json:"public_key"`
	// 私钥
	PrivateKey string `json:"private_key"`
}

// QueryCloudResourceRequest 查询云资源请求
type QueryCloudResourceRequest struct {
	NewClientReq NewCloudClientRequest `json:"new_client_req"`
	// 云ID
	CloudID string `json:"cloud_id"`
	// 资源类型
	ResourceType string `json:"resource_type"`
	// 云厂商ID
	ProviderID int `json:"provider_id"`
	// 云账号ID
	AccountID int `json:"account_id"`
	// 云账号CID
	AccountCID string `json:"account_cid"`
	// 云账号名称
	AccountName string `json:"account_name"`
	// 云账号别名
	AccountAlias string `json:"account_alias"`
	// 云账号唯一ID
	AccountUnionID string `json:"account_union_id"`
	// 云地域
	Region string `json:"region"`
	// 云资源ID
	ResourceID string `json:"id"`
	// 云资源名称
	ResourceName string `json:"name"`
	// 计费方式
	ChargeType string `json:"charge_type"`
	// 云项目ID
	ProjectID int `json:"project_id"`
	// 云地域ID
	RegionID int `json:"region_id"`
}

// SyncCloudResourceRequest 同步云资源请求
type SyncCloudResourceRequest struct {
	// 云地域
	Region string `json:"region"`
	// 云ID
	CloudID string `json:"cloud_id"`
	// 云账号ID
	AccountID int64 `json:"account_id"`
	// 云账号CID
	AccountCID string `json:"account_cid"`
	// 云账号名称
	AccountName string `json:"account_name"`
	// 云账号别名
	AccountAlias string `json:"account_alias"`
	// 云资源类型
	ResourceType string `json:"resource_type"`
	// 云资源名称
	ResourceName string `json:"resource_name"`
	// 云资源ID
	ResourceID string `json:"resource_id"`
}

// FormObject 对象
type FormObject struct {
	// ID
	ID int `json:"id"`
}
