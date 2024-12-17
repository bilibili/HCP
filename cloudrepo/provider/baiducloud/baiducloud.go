// Package baiducloud
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
package baiducloud

import (
	"github.com/bilibili/HCP/cloudrepo"
	"github.com/pkg/errors"
)

func init() {
	cloudrepo.Register(&baiducloudRepo{})
}

type baiducloudRepo struct {
}

// GetName 获取云厂商名称
func (b *baiducloudRepo) GetName() string {
	return cloudrepo.CloudBaidu
}

// GetProvider 获取云厂商client
func (b *baiducloudRepo) GetProvider(req *cloudrepo.GetProviderReq) (cloudrepo.CloudProvider, error) {
	client := NewClient(req.Region, req.SecretId, req.SecretKey)
	return &baiducloudProvider{client: client}, nil
}

type baiducloudProvider struct {
	client *Client
}

// ListRegion 获取云地域
func (b *baiducloudProvider) ListRegion(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudRegion, error) {
	return b.client.ListRegion(req)
}

// ListProject 获取云项目
func (b *baiducloudProvider) ListProject(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudProject, error) {
	return b.client.ListProject(req)
}

// ListSecurityGroup 获取云安全组
func (b *baiducloudProvider) ListSecurityGroup(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudSecurityGroup, error) {
	return b.client.ListCloudSecurityGroup(req)
}

// ListZone 获取云可用区
func (b *baiducloudProvider) ListZone(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudZone, error) {
	return b.client.ListZone(req)
}

// ListVpc 获取云专有网络
func (b *baiducloudProvider) ListVpc(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudVpc, error) {
	return b.client.ListVpc(req)
}

// ListSubnet 获取云子网
func (b *baiducloudProvider) ListSubnet(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudSubnet, error) {
	return b.client.ListSubnet(req)
}

// ListServerSpec 获取云服务器规格
func (b *baiducloudProvider) ListServerSpec(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServerSpec, error) {
	return b.client.ListServerSpec(req)
}

// ListServerImage 获取云服务器镜像
func (b *baiducloudProvider) ListServerImage(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServerImage, error) {
	return b.client.ListServerImage(req)
}

// ListServer 获取云服务器
func (b *baiducloudProvider) ListServer(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServer, error) {
	return b.client.ListServer(req)
}

// ListTag 获取云标签
func (b *baiducloudProvider) ListTag(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudTag, error) {
	return b.client.ListTag(req)
}

// RebootServer 重启云服务器
func (b *baiducloudProvider) RebootServer(req *cloudrepo.RebootCloudServerReq) (*cloudrepo.RebootCloudServerReply, error) {
	return b.client.RebootServer(req)
}

// ReinstallServer 重装系统
func (b *baiducloudProvider) ReinstallServer(req *cloudrepo.ReinstallCloudServerReq) (*cloudrepo.ReinstallCloudServerReply, error) {
	return b.client.ReinstallServer(req)
}

// RenameServer 重命名云服务器
func (b *baiducloudProvider) RenameServer(req *cloudrepo.RenameCloudServerReq) (*cloudrepo.RenameCloudServerReply, error) {
	return b.client.RenameServer(req)
}

// ChangeServerConfig 变更云服务器配置
func (b *baiducloudProvider) ChangeServerConfig(req *cloudrepo.ChangeConfigCloudServerReq) (*cloudrepo.ChangeConfigCloudServerReply, error) {
	return b.client.ChangeServerConfig(req)
}

// DeleteServer 删除云服务器
func (b *baiducloudProvider) DeleteServer(req *cloudrepo.DeleteCloudServerReq) (*cloudrepo.DeleteCloudServerReply, error) {
	return b.client.DeleteServer(req)
}

// StartServer 启动云服务器
func (b *baiducloudProvider) StartServer(req *cloudrepo.StartCloudServerReq) (*cloudrepo.StartCloudServerReply, error) {
	return b.client.StartServer(req)
}

// StopServer 停止云服务器
func (b *baiducloudProvider) StopServer(req *cloudrepo.StopCloudServerReq) (*cloudrepo.StopCloudServerReply, error) {
	return b.client.StopServer(req)
}

// DescribeServer 获取云服务器详情
func (b *baiducloudProvider) DescribeServer(req *cloudrepo.DescribeCloudServerReq) (*cloudrepo.DescribeCloudServerReply, error) {
	return b.client.DescribeServer(req)
}

// ChangeServerChargeType 变更云服务器计费方式
func (b *baiducloudProvider) ChangeServerChargeType(req *cloudrepo.ChangeServerChargeTypeReq) (*cloudrepo.ChangeServerChargeTypeReply, error) {
	return nil, errors.Errorf("not support charge type server in baiducloud")
}
