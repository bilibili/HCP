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

import (
	"github.com/bilibili/HCP/cloudrepo"
	"github.com/pkg/errors"
)

func init() {
	cloudrepo.Register(&tencentcloudRepo{})
}

type tencentcloudRepo struct {
}

// GetName 获取云厂商名称
func (t *tencentcloudRepo) GetName() string {
	return cloudrepo.CloudTencent
}

// GetProvider 获取云厂商的client
func (t *tencentcloudRepo) GetProvider(req *cloudrepo.GetProviderReq) (cloudrepo.CloudProvider, error) {
	client := NewClient(req.Region, req.SecretId, req.SecretKey)
	return &tencentcloudProvider{client: client}, nil
}

type tencentcloudProvider struct {
	client *Client
}

// ListRegion 获取云地域
func (t *tencentcloudProvider) ListRegion(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudRegion, error) {
	return t.client.ListRegion(req)
}

// ListProject 获取云项目
func (t *tencentcloudProvider) ListProject(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudProject, error) {
	return t.client.ListProject(req)
}

// ListSecurityGroup 获取云安全组
func (t *tencentcloudProvider) ListSecurityGroup(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudSecurityGroup, error) {
	return t.client.ListCloudSecurityGroup(req)
}

// ListZone 获取云可用区
func (t *tencentcloudProvider) ListZone(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudZone, error) {
	return t.client.ListZone(req)
}

// ListVpc 获取云专有网络
func (t *tencentcloudProvider) ListVpc(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudVpc, error) {
	return t.client.ListVpc(req)
}

// ListSubnet 获取云子网
func (t *tencentcloudProvider) ListSubnet(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudSubnet, error) {
	return t.client.ListSubnet(req)
}

// ListServerSpec 获取云服务器规格
func (t *tencentcloudProvider) ListServerSpec(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServerSpec, error) {
	return t.client.ListServerSpec(req)
}

// ListServerImage 获取云服务器镜像
func (t *tencentcloudProvider) ListServerImage(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServerImage, error) {
	return t.client.ListServerImage(req)
}

// ListServer 获取云服务器
func (t *tencentcloudProvider) ListServer(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServer, error) {
	return t.client.ListServer(req)
}

// ListTag 获取云标签
func (t *tencentcloudProvider) ListTag(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudTag, error) {
	return t.client.ListTag(req)
}

// RebootServer 重启云服务器
func (t *tencentcloudProvider) RebootServer(req *cloudrepo.RebootCloudServerReq) (*cloudrepo.RebootCloudServerReply, error) {
	return t.client.RebootServer(req)
}

// ReinstallServer 重装系统
func (t *tencentcloudProvider) ReinstallServer(req *cloudrepo.ReinstallCloudServerReq) (*cloudrepo.ReinstallCloudServerReply, error) {
	return t.client.ReinstallServer(req)
}

// RenameServer 重命名云服务器
func (t *tencentcloudProvider) RenameServer(req *cloudrepo.RenameCloudServerReq) (*cloudrepo.RenameCloudServerReply, error) {
	return t.client.RenameServer(req)
}

// ChangeServerConfig 变更云服务器配置
func (t *tencentcloudProvider) ChangeServerConfig(req *cloudrepo.ChangeConfigCloudServerReq) (*cloudrepo.ChangeConfigCloudServerReply, error) {
	return t.client.ChangeServerConfig(req)
}

// DeleteServer 删除云服务器
func (t *tencentcloudProvider) DeleteServer(req *cloudrepo.DeleteCloudServerReq) (*cloudrepo.DeleteCloudServerReply, error) {
	return t.client.DeleteServer(req)
}

// StartServer 启动云服务器
func (t *tencentcloudProvider) StartServer(req *cloudrepo.StartCloudServerReq) (*cloudrepo.StartCloudServerReply, error) {
	return t.client.StartServer(req)
}

// StopServer 停止云服务器
func (t *tencentcloudProvider) StopServer(req *cloudrepo.StopCloudServerReq) (*cloudrepo.StopCloudServerReply, error) {
	return t.client.StopServer(req)
}

// DescribeServer 获取云服务器详情
func (t *tencentcloudProvider) DescribeServer(req *cloudrepo.DescribeCloudServerReq) (*cloudrepo.DescribeCloudServerReply, error) {
	return t.client.DescribeServer(req)
}

// ChangeServerChargeType 变更云服务器计费类型
func (t *tencentcloudProvider) ChangeServerChargeType(req *cloudrepo.ChangeServerChargeTypeReq) (*cloudrepo.ChangeServerChargeTypeReply, error) {
	return nil, errors.Errorf("not support charge type server in tencentcloud")
}
