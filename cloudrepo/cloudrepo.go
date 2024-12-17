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

import (
	"fmt"
	"strings"

	"github.com/araddon/dateparse"
	"github.com/bilibili/HCP/app/interface/v1/configs"
)

// GetProviderReq 获取云厂商请求
type GetProviderReq struct {
	// 云地域
	Region string
	// 密钥ID
	SecretId string
	// 密钥
	SecretKey string
}

// CloudRepo 云厂商数据
type CloudRepo interface {
	GetName() string                                        // 获取云厂商名称
	GetProvider(req *GetProviderReq) (CloudProvider, error) // 获取云厂商client
}

// CloudProvider 获取云厂商产品
type CloudProvider interface {
	ListRegion(req *GetCloudProductReq) ([]CloudRegion, error)                                   //获取地域
	ListProject(req *GetCloudProductReq) ([]CloudProject, error)                                 //获取项目
	ListServer(req *GetCloudProductReq) ([]CloudServer, error)                                   //获取服务器
	ListSecurityGroup(req *GetCloudProductReq) ([]CloudSecurityGroup, error)                     //获取安全组
	ListZone(req *GetCloudProductReq) ([]CloudZone, error)                                       //获取可用区
	ListVpc(req *GetCloudProductReq) ([]CloudVpc, error)                                         //获取专有网路
	ListSubnet(req *GetCloudProductReq) ([]CloudSubnet, error)                                   //获取子网
	ListServerImage(req *GetCloudProductReq) ([]CloudServerImage, error)                         //获取镜像
	ListServerSpec(req *GetCloudProductReq) ([]CloudServerSpec, error)                           //获取服务器规格
	ListTag(req *GetCloudProductReq) ([]CloudTag, error)                                         //获取标签
	RebootServer(req *RebootCloudServerReq) (*RebootCloudServerReply, error)                     //重启
	ReinstallServer(req *ReinstallCloudServerReq) (*ReinstallCloudServerReply, error)            //重装
	RenameServer(req *RenameCloudServerReq) (*RenameCloudServerReply, error)                     //改名
	ChangeServerConfig(req *ChangeConfigCloudServerReq) (*ChangeConfigCloudServerReply, error)   //改配
	DeleteServer(req *DeleteCloudServerReq) (*DeleteCloudServerReply, error)                     //清退
	StartServer(req *StartCloudServerReq) (*StartCloudServerReply, error)                        //开机
	StopServer(req *StopCloudServerReq) (*StopCloudServerReply, error)                           //关机
	DescribeServer(req *DescribeCloudServerReq) (*DescribeCloudServerReply, error)               //查询
	ChangeServerChargeType(req *ChangeServerChargeTypeReq) (*ChangeServerChargeTypeReply, error) //改费
}

var (
	repos        = map[string]CloudRepo{}
	CloudEnumMap = map[string]map[string]map[string]string{
		CloudAli:     statusAli,
		CloudAws:     statusAws,
		CloudTencent: statusTencent,
		CloudBaidu:   statusBaidu,
		CloudHuawei:  statusHuawei,
	}
)

// Register 注册云厂商
func Register(repo CloudRepo) {
	repos[repo.GetName()] = repo
}

// GetRepo 获取云厂商
func GetRepo(name string) (CloudRepo, error) {
	repo, ok := repos[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return repo, nil
}

// GetCloudEnum 获取云厂商
func GetCloudEnum(cloudId string, enumType string, enum string) string {
	cloudMap, ok := CloudEnumMap[cloudId]
	if !ok {
		return "unknown"
	}
	enumMap, ok := cloudMap[enumType]
	if !ok {
		return "unknown"
	}
	e, ok := enumMap[enum]
	if !ok {
		return "unknown"
	}
	return e
}

// GetGpuModel 获取GPU型号
func GetGpuModel(gpuModelRaw string) string {
	if strings.Contains(gpuModelRaw, "T4") {
		return "NVIDIA T4"
	}
	for _, gpu := range configs.Conf.ServerSpec.Gpu.GpuModels {
		if strings.Contains(strings.ToUpper(gpuModelRaw), strings.ToUpper(gpu)) {
			return gpu
		}
	}
	return ""
}

// TimeTrans 时间格式转换
func TimeTrans(tm string) string {
	var tmFormat string
	cTime, err := dateparse.ParseAny(tm)
	if err == nil {
		tmFormat = cTime.Format("2006-01-02T15:04:05Z")
	}
	return tmFormat
}
