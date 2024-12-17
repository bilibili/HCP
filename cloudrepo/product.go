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

// CloudProduct 云产品通用接口
type CloudProduct interface {
	GetCID() string    // 获取产品CID
	GetName() string   // 获取产品名称
	GetStatus() string // 获取产品状态
}

// CloudProject 云项目通用接口
type CloudProject interface {
	CloudProduct
}

// CloudRegion 云地域通用接口
type CloudRegion interface {
	CloudProduct
}

// CloudZone 云可用区通用接口
type CloudZone interface {
	CloudProduct
}

// CloudSecurityGroup 云安全组通用接口
type CloudSecurityGroup interface {
	CloudProduct

	GetProject() string // 获取云项目
	GetVpc() string     // 获取云专有网路
}

// CloudVpc 云专有网络通用接口
type CloudVpc interface {
	CloudProduct

	GetCidr() string    // 获取cidr
	GetProject() string // 获取云项目
}

// CloudSubnet 云子网通用接口
type CloudSubnet interface {
	CloudProduct

	GetCidr() string    // 获取cidr
	GetProject() string // 获取云项目
	GetVpc() string     // 获取云专有网络
	GetZone() string    // 获取云可用区
}

// CloudTag 云标签通用接口
type CloudTag interface {
	GetTagKey() string   // 获取标签key
	GetTagValue() string // 获取标签value
}

// CloudServerSpec 云服务器规格通用接口
type CloudServerSpec interface {
	CloudProduct

	GetBandWidth() float64 // 获取带宽
	GetCategory() string   // 获取规格类型
	GetCPU() int           // 获取CPU核数
	GetGPU() int           // 获取GPU数量
	GetGPUModel() string   // 获取GPU型号
	GetFamily() string     // 获取规格系列
	GetMemory() int        // 获取内存大小
	GetPPS() float64       // 获取每秒处理量
}

// CloudServerImage 云服务器镜像通用接口
type CloudServerImage interface {
	CloudProduct

	GetOsName() string // 获取操作系统名称
	GetType() string   // 获取镜像类型
}

// CloudServer 云服务器通用接口
type CloudServer interface {
	CloudProduct

	GetSecurityGroupCid() []string // 获取云安全组cid
	GetSubnetCid() string          // 获取云子网cid
	GetVpcCid() string             // 获取云专有网络cid
	GetImageCid() string           // 获取云服务器镜像cid
	GetZoneCid() string            // 获取云可用区cid
	GetProjectCid() string         // 获取云项目cid
	GetServerSpec() string         // 获取云服务器规格cid
	GetChangeType() string         // 获取计费方式
	GetRenewStatus() string        // 获取续费状态
	GetPrivateIP() string          // 获取内网IP
	GetPublicIP() string           // 获取外网IP
	GetExpireTime() string         // 获取过期时间
}
