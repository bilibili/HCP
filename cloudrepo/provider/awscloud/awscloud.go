// Package awscloud
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
package awscloud

import (
	"github.com/bilibili/HCP/cloudrepo"
	"github.com/pkg/errors"
)

func init() {
	cloudrepo.Register(&awscloudRepo{})
}

type awscloudRepo struct {
}

// GetName 获取云厂商名称
func (a *awscloudRepo) GetName() string {
	return cloudrepo.CloudAws
}

// GetProvider 获取云厂商client
func (a *awscloudRepo) GetProvider(req *cloudrepo.GetProviderReq) (cloudrepo.CloudProvider, error) {
	client := NewClient(req.Region, req.SecretId, req.SecretKey)
	return &awscloudProvider{client: client}, nil
}

type awscloudProvider struct {
	client *Client
}

// ListRegion 获取云地域
func (a *awscloudProvider) ListRegion(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudRegion, error) {
	return a.client.ListRegion(req)
}

// ListProject 获取云项目
func (a *awscloudProvider) ListProject(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudProject, error) {
	return a.client.ListProject(req)
}

// ListSecurityGroup 获取云安全组
func (a *awscloudProvider) ListSecurityGroup(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudSecurityGroup, error) {
	return a.client.ListCloudSecurityGroup(req)
}

// ListZone 获取云可用区
func (a *awscloudProvider) ListZone(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudZone, error) {
	return a.client.ListZone(req)
}

// ListVpc 获取云专有网路
func (a *awscloudProvider) ListVpc(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudVpc, error) {
	return a.client.ListVpc(req)
}

// ListSubnet 获取云子网
func (a *awscloudProvider) ListSubnet(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudSubnet, error) {
	return a.client.ListSubnet(req)
}

// ListServerSpec 获取云服务器规格
func (a *awscloudProvider) ListServerSpec(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServerSpec, error) {
	return a.client.ListServerSpec(req)
}

// ListServerImage 获取云服务器镜像
func (a *awscloudProvider) ListServerImage(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServerImage, error) {
	return a.client.ListServerImage(req)
}

// ListServer 获取云服务器
func (a *awscloudProvider) ListServer(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServer, error) {
	return a.client.ListServer(req)
}

// ListTag 获取云标签
func (a *awscloudProvider) ListTag(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudTag, error) {
	return a.client.ListTag(req)
}

// RebootServer 重启云服务器
func (a *awscloudProvider) RebootServer(req *cloudrepo.RebootCloudServerReq) (*cloudrepo.RebootCloudServerReply, error) {
	return nil, errors.Errorf("not support reboot server in awscloud")
}

// ReinstallServer 重装系统
func (a *awscloudProvider) ReinstallServer(req *cloudrepo.ReinstallCloudServerReq) (*cloudrepo.ReinstallCloudServerReply, error) {
	return nil, errors.Errorf("not support reintall server in awscloud")
}

// RenameServer 重命名云服务器
func (a *awscloudProvider) RenameServer(req *cloudrepo.RenameCloudServerReq) (*cloudrepo.RenameCloudServerReply, error) {
	return nil, errors.Errorf("not support rename server in awscloud")
}

// ChangeServerConfig 变更云服务器配置
func (a *awscloudProvider) ChangeServerConfig(req *cloudrepo.ChangeConfigCloudServerReq) (*cloudrepo.ChangeConfigCloudServerReply, error) {
	return a.client.ChangeServerConfig(req)
}

// DeleteServer 删除云服务器
func (a *awscloudProvider) DeleteServer(req *cloudrepo.DeleteCloudServerReq) (*cloudrepo.DeleteCloudServerReply, error) {
	return a.client.DeleteServer(req)
}

// StartServer 启动云服务器
func (a *awscloudProvider) StartServer(req *cloudrepo.StartCloudServerReq) (*cloudrepo.StartCloudServerReply, error) {
	return a.client.StartServer(req)
}

// StopServer 停止云服务器
func (a *awscloudProvider) StopServer(req *cloudrepo.StopCloudServerReq) (*cloudrepo.StopCloudServerReply, error) {
	return a.client.StopServer(req)
}

// DescribeServer 获取云服务器详情
func (a *awscloudProvider) DescribeServer(req *cloudrepo.DescribeCloudServerReq) (*cloudrepo.DescribeCloudServerReply, error) {
	return a.client.DescribeServer(req)
}

// ChangeServerChargeType 变更云服务器计费方式
func (a *awscloudProvider) ChangeServerChargeType(req *cloudrepo.ChangeServerChargeTypeReq) (*cloudrepo.ChangeServerChargeTypeReply, error) {
	return nil, errors.Errorf("not support charge type server in awscloud")
}
