// Package tencentcloud
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
package tencentcloud

// Zone 可用区
type Zone struct {
	Zone     string `json:"Zone,omitempty" name:"Zone"`
	ZoneName string `json:"ZoneName,omitempty" name:"ZoneName"`
	//ZoneId    string `json:"ZoneId,omitempty" name:"ZoneId"`
	ZoneState string `json:"ZoneState,omitempty" name:"ZoneState"`
}

// Vpc 专有网络
type Vpc struct {
	// `VPC`名称。
	VpcName string `json:"VpcName,omitempty" name:"VpcName"`
	// `VPC`实例`ID`，例如：vpc-azd4dt1c。
	VpcId string `json:"VpcId,omitempty" name:"VpcId"`
	// `VPC`的`IPv4` `CIDR`。
	CidrBlock string `json:"CidrBlock,omitempty" name:"CidrBlock"`
	// 是否默认`VPC`。
	IsDefault bool `json:"IsDefault,omitempty" name:"IsDefault"`
	// 是否开启组播。
	EnableMulticast bool `json:"EnableMulticast,omitempty" name:"EnableMulticast"`
	// 创建时间。
	CreatedTime string `json:"CreatedTime,omitempty" name:"CreatedTime"`
	// `DNS`列表。
	DnsServerSet []string `json:"DnsServerSet,omitempty" name:"DnsServerSet"`
	// `DHCP`域名选项值。
	DomainName string `json:"DomainName,omitempty" name:"DomainName"`
	// `DHCP`选项集`ID`。
	DhcpOptionsId string `json:"DhcpOptionsId,omitempty" name:"DhcpOptionsId"`
	// 是否开启`DHCP`。
	EnableDhcp bool `json:"EnableDhcp,omitempty" name:"EnableDhcp"`
	// `VPC`的`IPv6` `CIDR`。
	Ipv6CidrBlock string `json:"Ipv6CidrBlock,omitempty" name:"Ipv6CidrBlock"`
	// 标签键值对
	TagSet []Tag `json:"TagSet,omitempty" name:"TagSet"`
	// 辅助CIDR
	// 注意：此字段可能返回 null，表示取不到有效值。
	AssistantCidrSet []AssistantCidr `json:"AssistantCidrSet,omitempty" name:"AssistantCidrSet"`
}

// AssistantCidr 辅助CIDR
type AssistantCidr struct {
	// `VPC`实例`ID`。形如：`vpc-6v2ht8q5`
	VpcId string `json:"VpcId,omitempty" name:"VpcId"`
	// 辅助CIDR。形如：`172.16.0.0/16`
	CidrBlock string `json:"CidrBlock,omitempty" name:"CidrBlock"`
	// 辅助CIDR类型（0：普通辅助CIDR，1：容器辅助CIDR），默认都是0。
	AssistantType int64 `json:"AssistantType,omitempty" name:"AssistantType"`
	// 辅助CIDR拆分的子网。
	// 注意：此字段可能返回 null，表示取不到有效值。
	SubnetSet []Subnet `json:"SubnetSet,omitempty" name:"SubnetSet"`
}

// Subnet 子网
type Subnet struct {
	// `VPC`实例`ID`。
	VpcId string `json:"VpcId,omitempty" name:"VpcId"`
	// 子网实例`ID`，例如：subnet-bthucmmy。
	SubnetId string `json:"SubnetId,omitempty" name:"SubnetId"`
	// 子网名称。
	SubnetName string `json:"SubnetName,omitempty" name:"SubnetName"`
	// 子网的 `IPv4` `CIDR`。
	CidrBlock string `json:"CidrBlock,omitempty" name:"CidrBlock"`
	// 是否默认子网。
	IsDefault bool `json:"IsDefault,omitempty" name:"IsDefault"`
	// 是否开启广播。
	EnableBroadcast bool `json:"EnableBroadcast,omitempty" name:"EnableBroadcast"`
	// 可用区。
	Zone string `json:"Zone,omitempty" name:"Zone"`
	// 路由表实例ID，例如：rtb-l2h8d7c2。
	RouteTableId string `json:"RouteTableId,omitempty" name:"RouteTableId"`
	// 创建时间。
	CreatedTime string `json:"CreatedTime,omitempty" name:"CreatedTime"`
	// 可用`IPv4`数。
	AvailableIpAddressCount uint64 `json:"AvailableIpAddressCount,omitempty" name:"AvailableIpAddressCount"`
	// 子网的 `IPv6` `CIDR`。
	Ipv6CidrBlock string `json:"Ipv6CidrBlock,omitempty" name:"Ipv6CidrBlock"`
	// 关联`ACL`ID
	NetworkAclId string `json:"NetworkAclId,omitempty" name:"NetworkAclId"`
	// 是否为 `SNAT` 地址池子网。
	IsRemoteVpcSnat bool `json:"IsRemoteVpcSnat,omitempty" name:"IsRemoteVpcSnat"`
	// 子网`IPv4`总数。
	TotalIpAddressCount uint64 `json:"TotalIpAddressCount,omitempty" name:"TotalIpAddressCount"`
	// 标签键值对。
	TagSet []Tag `json:"TagSet,omitempty" name:"TagSet"`
	// CDC实例ID。
	// 注意：此字段可能返回 null，表示取不到有效值。
	CdcId string `json:"CdcId,omitempty" name:"CdcId"`
	// 是否是CDC所属子网。0:否 1:是
	// 注意：此字段可能返回 null，表示取不到有效值。
	IsCdcSubnet int64 `json:"IsCdcSubnet,omitempty" name:"IsCdcSubnet"`
}

// SecurityGroup 安全组
type SecurityGroup struct {
	// 安全组实例ID，例如：sg-ohuuioma。
	SecurityGroupId string `json:"SecurityGroupId,omitempty" name:"SecurityGroupId"`
	// 安全组名称，可任意命名，但不得超过60个字符。
	SecurityGroupName string `json:"SecurityGroupName,omitempty" name:"SecurityGroupName"`
	// 安全组备注，最多100个字符。
	SecurityGroupDesc string `json:"SecurityGroupDesc,omitempty" name:"SecurityGroupDesc"`
	// 项目id，默认0。可在qcloud控制台项目管理页面查询到。
	ProjectId string `json:"ProjectId,omitempty" name:"ProjectId"`
	// 是否是默认安全组，默认安全组不支持删除。
	IsDefault bool `json:"IsDefault,omitempty" name:"IsDefault"`
	// 安全组创建时间。
	CreatedTime string `json:"CreatedTime,omitempty" name:"CreatedTime"`
	// 标签键值对。
	TagSet []Tag `json:"TagSet,omitempty" name:"TagSet"`
	// 安全组更新时间。
	// 注意：此字段可能返回 null，表示取不到有效值。
	UpdateTime string `json:"UpdateTime,omitempty" name:"UpdateTime"`
}

// Project	项目
type Project struct {
	// 项目ID
	ProjectId int64 `json:"ProjectId,omitempty" name:"ProjectId"`
	// 资源拥有者（主账号）uin
	OwnerUin int64 `json:"OwnerUin,omitempty" name:"OwnerUin"`
	// 应用Id
	AppId int64 `json:"AppId,omitempty" name:"AppId"`
	// 项目名称
	Name string `json:"Name,omitempty" name:"Name"`
	// 创建者uin
	CreatorUin int64 `json:"CreatorUin,omitempty" name:"CreatorUin"`
	// 来源平台
	SrcPlat string `json:"SrcPlat,omitempty" name:"SrcPlat"`
	// 来源AppId
	SrcAppId int64 `json:"SrcAppId,omitempty" name:"SrcAppId"`
	// 项目状态,0正常，-1关闭。默认项目为3
	Status     int64  `json:"Status,omitempty" name:"Status"`
	CreateTime string `json:"CreateTime,omitempty" name:"CreateTime"`
	IsDefault  int64  `json:"IsDefault,omitempty" name:"IsDefault"`
	Info       string `json:"Info,omitempty" name:"Info"`
}

// InstanceTypeQuotaItem 可用区机型配置列表
type InstanceTypeQuotaItem struct {
	// 可用区。
	Zone string `json:"Zone,omitempty" name:"Zone"`
	// 实例机型。
	InstanceType string `json:"InstanceType,omitempty" name:"InstanceType"`
	// 实例计费模式。取值范围： <br><li>PREPAID：表示预付费，即包年包月<br><li>POSTPAID_BY_HOUR：表示后付费，即按量计费<br><li>CDHPAID：表示[CDH](https://cloud.tencent.com/document/product/416)付费，即只对CDH计费，不对CDH上的实例计费。<br><li>`SPOTPAID`：表示竞价实例付费。
	InstanceChargeType string `json:"InstanceChargeType,omitempty" name:"InstanceChargeType"`
	// 网卡类型，例如：25代表25G网卡
	NetworkCard int64 `json:"NetworkCard,omitempty" name:"NetworkCard"`
	// 扩展属性。
	// 注意：此字段可能返回 null，表示取不到有效值。
	//Externals Externals `json:"Externals,omitempty" name:"Externals"`
	// 实例的CPU核数，单位：核。
	Cpu int64 `json:"Cpu,omitempty" name:"Cpu"`
	// 实例内存容量，单位：`GB`。
	Memory int64 `json:"Memory,omitempty" name:"Memory"`
	// 实例机型系列。
	InstanceFamily string `json:"InstanceFamily,omitempty" name:"InstanceFamily"`
	// 机型名称。
	TypeName string `json:"TypeName,omitempty" name:"TypeName"`
	// 本地磁盘规格列表。当该参数返回为空值时，表示当前情况下无法创建本地盘。
	LocalDiskTypeList []LocalDiskType `json:"LocalDiskTypeList,omitempty" name:"LocalDiskTypeList"`
	// 实例是否售卖。取值范围： <br><li>SELL：表示实例可购买<br><li>SOLD_OUT：表示实例已售罄。
	Status string `json:"Status,omitempty" name:"Status"`
	// 实例的售卖价格。
	Price ItemPrice `json:"Price,omitempty" name:"Price"`
	// 售罄原因。
	// 注意：此字段可能返回 null，表示取不到有效值。
	SoldOutReason string `json:"SoldOutReason,omitempty" name:"SoldOutReason"`
	// 内网带宽，单位Gbps。
	InstanceBandwidth float64 `json:"InstanceBandwidth,omitempty" name:"InstanceBandwidth"`
	// 网络收发包能力，单位万PPS。
	InstancePps int64 `json:"InstancePps,omitempty" name:"InstancePps"`
	// 本地存储块数量。
	StorageBlockAmount int64 `json:"StorageBlockAmount,omitempty" name:"StorageBlockAmount"`
	// 处理器型号。
	CpuType string `json:"CpuType,omitempty" name:"CpuType"`
	// 实例的GPU数量。
	Gpu int64 `json:"Gpu,omitempty" name:"Gpu"`
	// 实例的FPGA数量。
	Fpga int64 `json:"Fpga,omitempty" name:"Fpga"`
	// 实例备注信息。
	Remark string `json:"Remark,omitempty" name:"Remark"`
}

// ItemPrice 实例的售卖价格
type ItemPrice struct {

	// 后续合计费用的原价，后付费模式使用，单位：元。<br><li>如返回了其他时间区间项，如UnitPriceSecondStep，则本项代表时间区间在(0, 96)小时；若未返回其他时间区间项，则本项代表全时段，即(0, ∞)小时
	// 注意：此字段可能返回 null，表示取不到有效值。
	UnitPrice *float64 `json:"UnitPrice,omitempty" name:"UnitPrice"`

	// 后续计价单元，后付费模式使用，可取值范围： <br><li>HOUR：表示计价单元是按每小时来计算。当前涉及该计价单元的场景有：实例按小时后付费（POSTPAID_BY_HOUR）、带宽按小时后付费（BANDWIDTH_POSTPAID_BY_HOUR）：<br><li>GB：表示计价单元是按每GB来计算。当前涉及该计价单元的场景有：流量按小时后付费（TRAFFIC_POSTPAID_BY_HOUR）。
	// 注意：此字段可能返回 null，表示取不到有效值。
	ChargeUnit *string `json:"ChargeUnit,omitempty" name:"ChargeUnit"`

	// 预支合计费用的原价，预付费模式使用，单位：元。
	// 注意：此字段可能返回 null，表示取不到有效值。
	OriginalPrice *float64 `json:"OriginalPrice,omitempty" name:"OriginalPrice"`

	// 预支合计费用的折扣价，预付费模式使用，单位：元。
	// 注意：此字段可能返回 null，表示取不到有效值。
	DiscountPrice *float64 `json:"DiscountPrice,omitempty" name:"DiscountPrice"`

	// 折扣，如20.0代表2折。
	// 注意：此字段可能返回 null，表示取不到有效值。
	Discount *float64 `json:"Discount,omitempty" name:"Discount"`

	// 后续合计费用的折扣价，后付费模式使用，单位：元<br><li>如返回了其他时间区间项，如UnitPriceDiscountSecondStep，则本项代表时间区间在(0, 96)小时；若未返回其他时间区间项，则本项代表全时段，即(0, ∞)小时
	// 注意：此字段可能返回 null，表示取不到有效值。
	UnitPriceDiscount *float64 `json:"UnitPriceDiscount,omitempty" name:"UnitPriceDiscount"`

	// 使用时间区间在(96, 360)小时的后续合计费用的原价，后付费模式使用，单位：元。
	// 注意：此字段可能返回 null，表示取不到有效值。
	UnitPriceSecondStep *float64 `json:"UnitPriceSecondStep,omitempty" name:"UnitPriceSecondStep"`

	// 使用时间区间在(96, 360)小时的后续合计费用的折扣价，后付费模式使用，单位：元
	// 注意：此字段可能返回 null，表示取不到有效值。
	UnitPriceDiscountSecondStep *float64 `json:"UnitPriceDiscountSecondStep,omitempty" name:"UnitPriceDiscountSecondStep"`

	// 使用时间区间在(360, ∞)小时的后续合计费用的原价，后付费模式使用，单位：元。
	// 注意：此字段可能返回 null，表示取不到有效值。
	UnitPriceThirdStep *float64 `json:"UnitPriceThirdStep,omitempty" name:"UnitPriceThirdStep"`

	// 使用时间区间在(360, ∞)小时的后续合计费用的折扣价，后付费模式使用，单位：元
	// 注意：此字段可能返回 null，表示取不到有效值。
	UnitPriceDiscountThirdStep *float64 `json:"UnitPriceDiscountThirdStep,omitempty" name:"UnitPriceDiscountThirdStep"`

	// 预支三年合计费用的原价，预付费模式使用，单位：元。
	// 注意：此字段可能返回 null，表示取不到有效值。
	// 注意：此字段可能返回 null，表示取不到有效值。
	OriginalPriceThreeYear *float64 `json:"OriginalPriceThreeYear,omitempty" name:"OriginalPriceThreeYear"`

	// 预支三年合计费用的折扣价，预付费模式使用，单位：元。
	// 注意：此字段可能返回 null，表示取不到有效值。
	// 注意：此字段可能返回 null，表示取不到有效值。
	DiscountPriceThreeYear *float64 `json:"DiscountPriceThreeYear,omitempty" name:"DiscountPriceThreeYear"`

	// 预支三年应用的折扣，如20.0代表2折。
	// 注意：此字段可能返回 null，表示取不到有效值。
	// 注意：此字段可能返回 null，表示取不到有效值。
	DiscountThreeYear *float64 `json:"DiscountThreeYear,omitempty" name:"DiscountThreeYear"`

	// 预支五年合计费用的原价，预付费模式使用，单位：元。
	// 注意：此字段可能返回 null，表示取不到有效值。
	// 注意：此字段可能返回 null，表示取不到有效值。
	OriginalPriceFiveYear *float64 `json:"OriginalPriceFiveYear,omitempty" name:"OriginalPriceFiveYear"`

	// 预支五年合计费用的折扣价，预付费模式使用，单位：元。
	// 注意：此字段可能返回 null，表示取不到有效值。
	// 注意：此字段可能返回 null，表示取不到有效值。
	DiscountPriceFiveYear *float64 `json:"DiscountPriceFiveYear,omitempty" name:"DiscountPriceFiveYear"`

	// 预支五年应用的折扣，如20.0代表2折。
	// 注意：此字段可能返回 null，表示取不到有效值。
	// 注意：此字段可能返回 null，表示取不到有效值。
	DiscountFiveYear *float64 `json:"DiscountFiveYear,omitempty" name:"DiscountFiveYear"`

	// 预支一年合计费用的原价，预付费模式使用，单位：元。
	// 注意：此字段可能返回 null，表示取不到有效值。
	// 注意：此字段可能返回 null，表示取不到有效值。
	OriginalPriceOneYear *float64 `json:"OriginalPriceOneYear,omitempty" name:"OriginalPriceOneYear"`

	// 预支一年合计费用的折扣价，预付费模式使用，单位：元。
	// 注意：此字段可能返回 null，表示取不到有效值。
	// 注意：此字段可能返回 null，表示取不到有效值。
	DiscountPriceOneYear *float64 `json:"DiscountPriceOneYear,omitempty" name:"DiscountPriceOneYear"`

	// 预支一年应用的折扣，如20.0代表2折。
	// 注意：此字段可能返回 null，表示取不到有效值。
	// 注意：此字段可能返回 null，表示取不到有效值。
	DiscountOneYear *float64 `json:"DiscountOneYear,omitempty" name:"DiscountOneYear"`
}

// LocalDiskType 本地磁盘规格列表
type LocalDiskType struct {
	// 本地磁盘类型。
	Type string `json:"Type,omitempty" name:"Type"`
	// 本地磁盘属性。
	PartitionType string `json:"PartitionType,omitempty" name:"PartitionType"`
	// 本地磁盘最小值。
	MinSize int64 `json:"MinSize,omitempty" name:"MinSize"`
	// 本地磁盘最大值。
	MaxSize int64 `json:"MaxSize,omitempty" name:"MaxSize"`
	// 购买时本地盘是否为必选。取值范围：<br><li>REQUIRED：表示必选<br><li>OPTIONAL：表示可选。
	Required string `json:"Required,omitempty" name:"Required"`
}

// Image 镜像
type Image struct {
	// 镜像ID
	ImageId string `json:"ImageId,omitempty" name:"ImageId"`
	// 镜像操作系统
	OsName string `json:"OsName,omitempty" name:"OsName"`
	// 镜像类型
	ImageType string `json:"ImageType,omitempty" name:"ImageType"`
	// 镜像创建时间
	CreatedTime string `json:"CreatedTime,omitempty" name:"CreatedTime"`
	// 镜像名称
	ImageName string `json:"ImageName,omitempty" name:"ImageName"`
	// 镜像描述
	ImageDescription string `json:"ImageDescription,omitempty" name:"ImageDescription"`
	// 镜像大小
	ImageSize int64 `json:"ImageSize,omitempty" name:"ImageSize"`
	// 镜像架构
	Architecture string `json:"Architecture,omitempty" name:"Architecture"`
	// 镜像状态:
	// CREATING-创建中
	// NORMAL-正常
	// CREATEFAILED-创建失败
	// USING-使用中
	// SYNCING-同步中
	// IMPORTING-导入中
	// IMPORTFAILED-导入失败
	ImageState string `json:"ImageState,omitempty" name:"ImageState"`
	// 镜像来源平台
	Platform string `json:"Platform,omitempty" name:"Platform"`
	// 镜像创建者
	ImageCreator string `json:"ImageCreator,omitempty" name:"ImageCreator"`
	// 镜像来源
	ImageSource string `json:"ImageSource,omitempty" name:"ImageSource"`
	// 同步百分比
	// 注意：此字段可能返回 null，表示取不到有效值。
	SyncPercent int64 `json:"SyncPercent,omitempty" name:"SyncPercent"`
	// 镜像是否支持cloud-init
	// 注意：此字段可能返回 null，表示取不到有效值。
	IsSupportCloudinit bool `json:"IsSupportCloudinit,omitempty" name:"IsSupportCloudinit"`
	// 镜像关联的快照信息
	// 注意：此字段可能返回 null，表示取不到有效值。
	SnapshotSet []Snapshot `json:"SnapshotSet,omitempty" name:"SnapshotSet"`
}

// Snapshot 镜像关联的快照信息
type Snapshot struct {
	// 快照Id。
	SnapshotId string `json:"SnapshotId,omitempty" name:"SnapshotId"`
	// 创建此快照的云硬盘类型。取值范围：
	// SYSTEM_DISK：系统盘
	// DATA_DISK：数据盘。
	DiskUsage string `json:"DiskUsage,omitempty" name:"DiskUsage"`
	// 创建此快照的云硬盘大小，单位GB。
	DiskSize int64 `json:"DiskSize,omitempty" name:"DiskSize"`
}

// Instance 实例
type Instance struct {
	Placement                Placement           `json:"Placement,omitempty" name:"Placement"`
	InstanceId               string              `json:"InstanceId,omitempty" name:"InstanceId"`
	InstanceType             string              `json:"InstanceType,omitempty" name:"InstanceType"`
	CPU                      int64               `json:"CPU,omitempty" name:"CPU"`
	Memory                   int64               `json:"Memory,omitempty" name:"Memory"`
	RestrictState            string              `json:"RestrictState,omitempty" name:"RestrictState"`
	InstanceName             string              `json:"InstanceName,omitempty" name:"InstanceName"`
	InstanceChargeType       string              `json:"InstanceChargeType,omitempty" name:"InstanceChargeType"`
	SystemDisk               SystemDisk          `json:"SystemDisk,omitempty" name:"SystemDisk"`
	DataDisks                []DataDisk          `json:"DataDisks,omitempty" name:"DataDisks"`
	PrivateIpAddresses       []string            `json:"PrivateIpAddresses,omitempty" name:"PrivateIpAddresses"`
	PublicIpAddresses        []string            `json:"PublicIpAddresses,omitempty" name:"PublicIpAddresses"`
	InternetAccessible       InternetAccessible  `json:"InternetAccessible,omitempty" name:"InternetAccessible"`
	VirtualPrivateCloud      VirtualPrivateCloud `json:"VirtualPrivateCloud,omitempty" name:"VirtualPrivateCloud"`
	ImageId                  string              `json:"ImageId,omitempty" name:"ImageId"`
	RenewFlag                string              `json:"RenewFlag,omitempty" name:"RenewFlag"`
	CreatedTime              string              `json:"CreatedTime,omitempty" name:"CreatedTime"`
	ExpiredTime              string              `json:"ExpiredTime,omitempty" name:"ExpiredTime"`
	OsName                   string              `json:"OsName,omitempty" name:"OsName"`
	SecurityGroupIds         []string            `json:"SecurityGroupIds,omitempty" name:"SecurityGroupIds"`
	LoginSettings            LoginSettings       `json:"LoginSettings,omitempty" name:"LoginSettings"`
	InstanceState            string              `json:"InstanceState,omitempty" name:"InstanceState"`
	Tags                     []Tag               `json:"Tags,omitempty" name:"Tags"`
	StopChargingMode         string              `json:"StopChargingMode,omitempty" name:"StopChargingMode"`
	Uuid                     string              `json:"Uuid,omitempty" name:"Uuid"`
	LatestOperation          string              `json:"LatestOperation,omitempty" name:"LatestOperation"`
	LatestOperationState     string              `json:"LatestOperationState,omitempty" name:"LatestOperationState"`
	LatestOperationRequestId string              `json:"LatestOperationRequestId,omitempty" name:"LatestOperationRequestId"`
	DisasterRecoverGroupId   string              `json:"DisasterRecoverGroupId,omitempty" name:"DisasterRecoverGroupId"`
	IPv6Addresses            []string            `json:"IPv6Addresses,omitempty" name:"IPv6Addresses"`
	CamRoleName              string              `json:"CamRoleName,omitempty" name:"CamRoleName"`
	HpcClusterId             string              `json:"HpcClusterId,omitempty" name:"HpcClusterId"`
	RdmaIpAddresses          []string            `json:"RdmaIpAddresses,omitempty" name:"RdmaIpAddresses"`
	IsolatedSource           string              `json:"IsolatedSource,omitempty" name:"IsolatedSource"`
	GPUInfo                  GPUInfo             `json:"GPUInfo,omitempty" name:"GPUInfo"`
}

// Placement 实例所属地域信息
type Placement struct {
	Zone      string   `json:"Zone,omitempty" name:"Zone"`
	ProjectId int64    `json:"ProjectId,omitempty" name:"ProjectId"`
	HostIds   []string `json:"HostIds,omitempty" name:"HostIds"`
	HostIps   []string `json:"HostIps,omitempty" name:"HostIps"`
	HostId    string   `json:"HostId,omitempty" name:"HostId"`
}

// SystemDisk 系统盘信息
type SystemDisk struct {
	DiskType string `json:"DiskType,omitempty" name:"DiskType"`
	DiskId   string `json:"DiskId,omitempty" name:"DiskId"`
	DiskSize int64  `json:"DiskSize,omitempty" name:"DiskSize"`
	CdcId    string `json:"CdcId,omitempty" name:"CdcId"`
}

// DataDisk 数据盘信息
type DataDisk struct {
	DiskSize              int64  `json:"DiskSize,omitempty" name:"DiskSize"`
	DiskType              string `json:"DiskType,omitempty" name:"DiskType"`
	DiskId                string `json:"DiskId,omitempty" name:"DiskId"`
	DeleteWithInstance    bool   `json:"DeleteWithInstance,omitempty" name:"DeleteWithInstance"`
	SnapshotId            string `json:"SnapshotId,omitempty" name:"SnapshotId"`
	Encrypt               bool   `json:"Encrypt,omitempty" name:"Encrypt"`
	KmsKeyId              string `json:"KmsKeyId,omitempty" name:"KmsKeyId"`
	ThroughputPerformance int64  `json:"ThroughputPerformance,omitempty" name:"ThroughputPerformance"`
	CdcId                 string `json:"CdcId,omitempty" name:"CdcId"`
}

// InternetAccessible 公网带宽信息
type InternetAccessible struct {
	InternetChargeType      string `json:"InternetChargeType,omitempty" name:"InternetChargeType"`
	InternetMaxBandwidthOut int64  `json:"InternetMaxBandwidthOut,omitempty" name:"InternetMaxBandwidthOut"`
	PublicIpAssigned        bool   `json:"PublicIpAssigned,omitempty" name:"PublicIpAssigned"`
	BandwidthPackageId      string `json:"BandwidthPackageId,omitempty" name:"BandwidthPackageId"`
}

// VirtualPrivateCloud 私有网络信息
type VirtualPrivateCloud struct {
	VpcId              string   `json:"VpcId,omitempty" name:"VpcId"`
	SubnetId           string   `json:"SubnetId,omitempty" name:"SubnetId"`
	AsVpcGateway       bool     `json:"AsVpcGateway,omitempty" name:"AsVpcGateway"`
	PrivateIpAddresses []string `json:"PrivateIpAddresses,omitempty" name:"PrivateIpAddresses"`
	Ipv6AddressCount   uint64   `json:"Ipv6AddressCount,omitempty" name:"Ipv6AddressCount"`
}

// LoginSettings 登录设置
type LoginSettings struct {
	Password       string   `json:"Password,omitempty" name:"Password"`
	KeyIds         []string `json:"KeyIds,omitempty" name:"KeyIds"`
	KeepImageLogin string   `json:"KeepImageLogin,omitempty" name:"KeepImageLogin"`
}

// GPUInfo GPU信息
type GPUInfo struct {
	GPUCount float64  `json:"GPUCount,omitempty" name:"GPUCount"`
	GPUId    []string `json:"GPUId,omitempty" name:"GPUId"`
	GPUType  string   `json:"GPUType,omitempty" name:"GPUType"`
}

// TagWithDelete 标签信息
type TagWithDelete struct {
	// 标签键
	TagKey string `json:"TagKey,omitempty" name:"TagKey"`
	// 标签值
	TagValue string `json:"TagValue,omitempty" name:"TagValue"`
	// 是否可以删除
	CanDelete uint64 `json:"CanDelete,omitempty" name:"CanDelete"`
}
