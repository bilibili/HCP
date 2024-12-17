// Package huaweicloud
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
package huaweicloud

import (
	"github.com/bilibili/HCP/cloudrepo"
	"github.com/pkg/errors"
)

func init() {
	cloudrepo.Register(&huaweicloudRepo{})
}

type huaweicloudRepo struct {
}

// GetName 获取云厂商名称
func (h *huaweicloudRepo) GetName() string {
	return cloudrepo.CloudHuawei
}

// GetProvider 获取云厂商client
func (h *huaweicloudRepo) GetProvider(req *cloudrepo.GetProviderReq) (cloudrepo.CloudProvider, error) {
	client := NewClient(req.Region, req.SecretId, req.SecretKey)
	return &huaweicloudProvider{client: client}, nil
}

type huaweicloudProvider struct {
	client *Client
}

// ListRegion 获取云地域
func (h *huaweicloudProvider) ListRegion(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudRegion, error) {
	return h.client.ListRegion(req)
}

// ListProject 获取云项目
func (h *huaweicloudProvider) ListProject(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudProject, error) {
	return h.client.ListProject(req)
}

// ListSecurityGroup 获取云安全组
func (h *huaweicloudProvider) ListSecurityGroup(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudSecurityGroup, error) {
	return h.client.ListCloudSecurityGroup(req)
}

// ListZone 获取云可用区
func (h *huaweicloudProvider) ListZone(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudZone, error) {
	return h.client.ListZone(req)
}

// ListVpc 获取云专有网络
func (h *huaweicloudProvider) ListVpc(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudVpc, error) {
	return h.client.ListVpc(req)
}

// ListSubnet 获取云子网
func (h *huaweicloudProvider) ListSubnet(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudSubnet, error) {
	return h.client.ListSubnet(req)
}

// ListServerSpec 获取云服务器规格
func (h *huaweicloudProvider) ListServerSpec(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServerSpec, error) {
	return h.client.ListServerSpec(req)
}

// ListServerImage 获取云服务器镜像
func (h *huaweicloudProvider) ListServerImage(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServerImage, error) {
	return h.client.ListServerImage(req)
}

// ListServer 获取云服务器
func (h *huaweicloudProvider) ListServer(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServer, error) {
	return h.client.ListServer(req)
}

// ListTag 获取云标签
func (h *huaweicloudProvider) ListTag(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudTag, error) {
	return nil, nil
}

// RebootServer 重启云服务器
func (h *huaweicloudProvider) RebootServer(req *cloudrepo.RebootCloudServerReq) (*cloudrepo.RebootCloudServerReply, error) {
	return h.client.RebootServer(req)
}

// ReinstallServer 重装系统
func (h *huaweicloudProvider) ReinstallServer(req *cloudrepo.ReinstallCloudServerReq) (*cloudrepo.ReinstallCloudServerReply, error) {
	return h.client.ReinstallServer(req)
}

// RenameServer 重命名云服务器
func (h *huaweicloudProvider) RenameServer(req *cloudrepo.RenameCloudServerReq) (*cloudrepo.RenameCloudServerReply, error) {
	return h.client.RenameServer(req)
}

// ChangeServerConfig 变更云服务器配置
func (h *huaweicloudProvider) ChangeServerConfig(req *cloudrepo.ChangeConfigCloudServerReq) (*cloudrepo.ChangeConfigCloudServerReply, error) {
	return h.client.ChangeServerConfig(req)
}

// DeleteServer 删除云服务器
func (h *huaweicloudProvider) DeleteServer(req *cloudrepo.DeleteCloudServerReq) (*cloudrepo.DeleteCloudServerReply, error) {
	return h.client.DeleteServer(req)
}

// StartServer 启动云服务器
func (h *huaweicloudProvider) StartServer(req *cloudrepo.StartCloudServerReq) (*cloudrepo.StartCloudServerReply, error) {
	return h.client.StartServer(req)
}

// StopServer 停止云服务器
func (h *huaweicloudProvider) StopServer(req *cloudrepo.StopCloudServerReq) (*cloudrepo.StopCloudServerReply, error) {
	return h.client.StopServer(req)
}

// DescribeServer 获取云服务器详情
func (h *huaweicloudProvider) DescribeServer(req *cloudrepo.DescribeCloudServerReq) (*cloudrepo.DescribeCloudServerReply, error) {
	return h.client.DescribeServer(req)
}

// ChangeServerChargeType 变更云服务器计费类型
func (h *huaweicloudProvider) ChangeServerChargeType(req *cloudrepo.ChangeServerChargeTypeReq) (*cloudrepo.ChangeServerChargeTypeReply, error) {
	return nil, errors.Errorf("not support charge type server in huaweiucloud")
}
