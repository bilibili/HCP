// Package cloudrepo
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
package cloudrepo

import "github.com/bilibili/HCP/common/ecode"

const (
	CloudAli     = "alicloud"     // 阿里云
	CloudTencent = "tencentcloud" // 腾讯云
	CloudHuawei  = "huaweicloud"  // 华为云
	CloudAws     = "awscloud"     // AWS
	CloudBaidu   = "baiducloud"   // 百度云
)

var (
	NotFoundError = ecode.NewECode("-60001")
)

// CloudProductCommon 云产品通用
type CloudProductCommon struct {
	// 云ID
	CID string
	// 云产品名称
	Name string
	// 云产品状态
	Status string
}

// GetCloudProductReq 获取云产品列表请求
type GetCloudProductReq struct {
	GetProviderReq
	// 云产品名称
	ResourceName string
	// 云产品ID
	ResourceID string
	// 云地域
	Region string
	// 标识
	Cursor string
	// 是否分页
	DisablePage bool
	// 条数
	PageSize int
	// 页码
	PageNum int
}

// RebootCloudServerReq 重启云服务器请求参数
type RebootCloudServerReq struct {
	// 云服务器ID
	InstanceID string
	// 是否强制
	IsForce bool
	// 云项目ID
	ProjectID string
}

// RebootCloudServerReply 重启云服务器响应参数
type RebootCloudServerReply struct {
	// 请求ID
	RequestID string
}

// ReinstallCloudServerReq 重装云服务器请求参数
type ReinstallCloudServerReq struct {
	// 云服务器ID
	InstanceID string
	// 镜像ID
	ImageID string
	// 云项目ID
	ProjectID string
	// 云账号ID
	AccountID string
	// 密钥对ID
	KeyPairID string
	// 密钥对名称
	KeyPairName string
}

// ReinstallCloudServerReply 重装云服务器响应参数
type ReinstallCloudServerReply struct {
	// 请求ID
	RequestID string
}

// RenameCloudServerReq 重命名云服务器名称请求参数
type RenameCloudServerReq struct {
	// 云服务器ID
	InstanceID string
	// 新名称
	NewName string
	// 云项目ID
	ProjectID string
}

// RenameCloudServerReply 重命名云服务器名称响应参数
type RenameCloudServerReply struct {
	// 请求ID
	RequestID string
}

// ChangeConfigCloudServerReq 变更云服务器配置请求参数
type ChangeConfigCloudServerReq struct {
	// 云服务器ID
	InstanceID string
	// 云服务器类型
	InstanceType string
	// 云项目ID
	ProjectID string
}

// ChangeConfigCloudServerReply 变更云服务器配置响应参数
type ChangeConfigCloudServerReply struct {
	// 请求ID
	RequestID string
}

// DeleteCloudServerReq 删除云服务器请求参数
type DeleteCloudServerReq struct {
	// 云服务器ID
	InstanceID string
	// 是否强制
	IsForce bool
	// 云项目ID
	ProjectID string
}

// DeleteCloudServerReply 删除云服务器响应参数
type DeleteCloudServerReply struct {
	// 请求ID
	RequestID string
}

// StartCloudServerReq 启动云服务器请求参数
type StartCloudServerReq struct {
	// 云服务器ID
	InstanceID string
	// 云项目ID
	ProjectID string
	// 是否强制
	IsForce bool
}

// StartCloudServerReply 启动云服务器响应参数
type StartCloudServerReply struct {
	// 请求ID
	RequestID string
}

// StopCloudServerReq 停止云服务器请求参数
type StopCloudServerReq struct {
	// 云服务器ID
	InstanceID string
	// 云项目ID
	ProjectID string
	// 是否强制
	IsForce bool
}

// StopCloudServerReply 停止云服务器响应参数
type StopCloudServerReply struct {
	// 请求ID
	RequestID string
}

// DescribeCloudServerReq 获取云服务器详情请求参数
type DescribeCloudServerReq struct {
	// 云服务器ID
	InstanceID string
	// 云项目ID
	ProjectID string
}

// DescribeCloudServerReply 获取云服务器详情响应参数
type DescribeCloudServerReply struct {
	// 云服务器ID
	InstanceID string
	// 云服务器名称
	InstanceName string
	// 镜像ID
	ImageID string
	// 云服务器状态
	InstanceState string
	// 云服务器类型
	InstanceType string
	// 计费方式
	ChargeType string
	// 实例的最新操作状态
	LatestOperationState string
	// 实例是否操作完成
	IsOperation bool
}

// ChangeServerChargeTypeReq 变更云服务器计费方式请求参数
type ChangeServerChargeTypeReq struct {
	// 云服务器ID
	InstanceID string
	// 云项目ID
	ProjectID string
	// 计费方式
	ChargeType string
}

// ChangeServerChargeTypeReply 变更云服务器计费方式响应参数
type ChangeServerChargeTypeReply struct {
	// 请求ID
	RequestID string
}
