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
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/bilibili/HCP/cloudrepo"
	"github.com/go-kratos/kratos/pkg/log"
)

// CloudServer 云服务器
type CloudServer struct {
	cloudrepo.CloudProductCommon

	SeverSecurityGroupCid []string
	SubnetCid             string
	VpcCid                string
	ImageCid              string
	ZoneCid               string
	ProjectCid            string
	ServerSpec            string
	ChangeType            string
	RenewStatus           string
	PrivateIP             string
	PublicIP              string
	ExpireTime            string
}

// GetCID 获取云服务器CID
func (c *CloudServer) GetCID() string {
	return c.CID
}

// GetName 获取云服务器名称
func (c *CloudServer) GetName() string {
	return c.Name
}

// GetStatus 获取云服务器状态
func (c *CloudServer) GetStatus() string {
	return c.Status
}

// GetServerSpec 获取云服务器规格
func (c *CloudServer) GetServerSpec() string {
	return c.ServerSpec
}

// GetChangeType 获取云服务器计费方式
func (c *CloudServer) GetChangeType() string {
	return c.ChangeType
}

// GetRenewStatus 获取云服务器续费状态
func (c *CloudServer) GetRenewStatus() string {
	return c.RenewStatus
}

// GetPrivateIP 获取云服务器内网ip
func (c *CloudServer) GetPrivateIP() string {
	return c.PrivateIP
}

// GetPublicIP 获取云服务器外网ip
func (c *CloudServer) GetPublicIP() string {
	return c.PublicIP
}

// GetExpireTime 获取云服务器过期时间
func (c *CloudServer) GetExpireTime() string {
	return c.ExpireTime
}

// GetSecurityGroupCid 获取云服务器安全组
func (c *CloudServer) GetSecurityGroupCid() []string {
	return c.SeverSecurityGroupCid
}

// GetSubnetCid 获取云服务器子网
func (c *CloudServer) GetSubnetCid() string {
	return c.SubnetCid
}

// GetVpcCid 获取云服务器专有网络
func (c *CloudServer) GetVpcCid() string {
	return c.VpcCid
}

// GetImageCid 获取云服务器镜像
func (c *CloudServer) GetImageCid() string {
	return c.ImageCid
}

// GetZoneCid 获取云服务器可用区
func (c *CloudServer) GetZoneCid() string {
	return c.ZoneCid
}

// GetProjectCid 获取云服务器项目
func (c *CloudServer) GetProjectCid() string {
	return c.ProjectCid
}

// ListServer 获取服务器列表
func (c *Client) ListServer(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServer, error) {
	cli := c.clientEc2()
	req.DisablePage = true
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFn()
	var data []types.Reservation
	pageSize := 100
	maxResult := int32(pageSize)
	request := &ec2.DescribeInstancesInput{MaxResults: &maxResult}
	for {
		if req.Cursor != "" {
			request.NextToken = &req.Cursor
		}
		resp, err := cli.DescribeInstances(ctx, nil)
		if err != nil {
			if strings.Contains(err.Error(), "InvalidOperation.NotSupportedEndpoint") {
				log.Error("queryServerAws QueryEc2 error (%s)", err.Error())
				return nil, nil
			}
			return nil, err
		}
		if resp.Reservations != nil && len(resp.Reservations) > 0 {
			data = append(data, resp.Reservations...)
		}
		if resp.NextToken != nil {
			req.Cursor = *resp.NextToken
		}
		if resp.NextToken == nil {
			break
		}
	}
	return c.list2DoServer(data)
}

func (c *Client) list2DoServer(resp []types.Reservation) (list []cloudrepo.CloudServer, err error) {
	for _, reservation := range resp {
		for _, v := range reservation.Instances {
			_ = v
			name := c.GetNameTag(v.Tags)
			var (
				intraIps          []string
				extraIps          []string
				securityGroupCids []string
				subnetCid         string
				vpcCid            string
				imageCid          string
				projectCid        string
				zoneCid           string
			)
			if v.PrivateIpAddress != nil {
				intraIps = []string{*v.PrivateIpAddress}
			}
			if v.PublicIpAddress != nil {
				extraIps = []string{*v.PublicIpAddress}
			}
			for _, sg := range v.SecurityGroups {
				securityGroupCids = append(securityGroupCids, *sg.GroupName)
			}
			if v.SubnetId != nil {
				subnetCid = *v.SubnetId
			}
			if v.VpcId != nil {
				vpcCid = *v.VpcId
			}
			if v.ImageId != nil {
				imageCid = *v.ImageId
			}
			if v.Placement != nil && v.Placement.AvailabilityZone != nil {
				zoneCid = *v.Placement.AvailabilityZone
			}
			projectCid = c.GetEc2ProjectTag(v.Tags)
			list = append(list, &CloudServer{
				CloudProductCommon: cloudrepo.CloudProductCommon{
					CID:    *v.InstanceId,
					Name:   name,
					Status: cloudrepo.GetCloudEnum(cloudrepo.CloudAws, "cloud_server", string(v.State.Name)),
				},
				ServerSpec:            string(v.InstanceType),
				ChangeType:            cloudrepo.GetCloudEnum(cloudrepo.CloudAws, "pay_type", ""),
				RenewStatus:           "",
				PrivateIP:             strings.Join(intraIps, ","),
				PublicIP:              strings.Join(extraIps, ","),
				ExpireTime:            "",
				SeverSecurityGroupCid: securityGroupCids,
				SubnetCid:             subnetCid,
				VpcCid:                vpcCid,
				ImageCid:              imageCid,
				ZoneCid:               zoneCid,
				ProjectCid:            projectCid,
			})
		}
	}
	return
}

// ChangeServerConfig 变更服务器配置
func (c *Client) ChangeServerConfig(req *cloudrepo.ChangeConfigCloudServerReq) (*cloudrepo.ChangeConfigCloudServerReply, error) {
	cli := c.clientEc2()
	request := &ec2.ModifyInstanceAttributeInput{}
	request.InstanceId = &req.InstanceID
	request.InstanceType = &types.AttributeValue{Value: &req.InstanceType}
	_, err := cli.ModifyInstanceAttribute(context.TODO(), request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.ChangeConfigCloudServerReply{
		RequestID: req.InstanceID,
	}, err
}

// DeleteServer 删除服务器
func (c *Client) DeleteServer(req *cloudrepo.DeleteCloudServerReq) (*cloudrepo.DeleteCloudServerReply, error) {
	cli := c.clientEc2()
	// 检查删除保护状态.如果已开启则先关闭删除保护再进行删除操作
	protect, err := c.getProtectStatusServer(req.InstanceID)
	if err != nil {
		return nil, err
	}
	if protect {
		log.Warn("DeleteVM instance %s which termination protect is in open status", req.InstanceID)
		//关闭删除保护
		err = c.closeProtectServer(req.InstanceID, false)
		if err != nil {
			return nil, err
		}
	}
	request := &ec2.TerminateInstancesInput{}
	request.InstanceIds = []string{req.InstanceID}
	_, err = cli.TerminateInstances(context.TODO(), request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.DeleteCloudServerReply{
		RequestID: req.InstanceID,
	}, err
}

// StartServer 启动服务器
func (c *Client) StartServer(req *cloudrepo.StartCloudServerReq) (*cloudrepo.StartCloudServerReply, error) {
	cli := c.clientEc2()
	request := &ec2.StartInstancesInput{}
	request.InstanceIds = []string{req.InstanceID}
	_, err := cli.StartInstances(context.TODO(), request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.StartCloudServerReply{
		RequestID: req.InstanceID,
	}, err
}

// StopServer 停止服务器
func (c *Client) StopServer(req *cloudrepo.StopCloudServerReq) (*cloudrepo.StopCloudServerReply, error) {
	cli := c.clientEc2()
	request := &ec2.StopInstancesInput{}
	request.InstanceIds = []string{req.InstanceID}
	_, err := cli.StopInstances(context.TODO(), request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.StopCloudServerReply{
		RequestID: req.InstanceID,
	}, err
}

// getProtectStatusServer 获取服务器删除保护状态
func (c *Client) getProtectStatusServer(instanceID string) (bool, error) {
	cli := c.clientEc2()
	request := &ec2.DescribeInstanceAttributeInput{}
	request.InstanceId = &instanceID
	request.Attribute = types.InstanceAttributeNameDisableApiTermination
	ret, err := cli.DescribeInstanceAttribute(context.TODO(), request)
	if err != nil {
		return false, err
	}
	return *ret.DisableApiTermination.Value, nil
}

// closeProtectServer 关闭服务器删除保护
func (c *Client) closeProtectServer(instanceID string, disableDelete bool) error {
	cli := c.clientEc2()
	request := &ec2.ModifyInstanceAttributeInput{
		DisableApiTermination: &types.AttributeBooleanValue{Value: &disableDelete},
		InstanceId:            &instanceID,
	}
	_, err := cli.ModifyInstanceAttribute(context.TODO(), request)
	if err != nil {
		return err
	}
	return nil
}

// DescribeServer 获取服务器信息
func (c *Client) DescribeServer(req *cloudrepo.DescribeCloudServerReq) (*cloudrepo.DescribeCloudServerReply, error) {
	cli := c.clientEc2()
	request := &ec2.DescribeInstancesInput{}
	request.InstanceIds = []string{req.InstanceID}
	resp, err := cli.DescribeInstances(context.TODO(), request)
	if err != nil {
		if strings.Contains(err.Error(), "InvalidInstanceID.NotFound") {
			return nil, cloudrepo.NotFoundError
		} else {
			return nil, err
		}
	}
	if resp.Reservations == nil || len(resp.Reservations) == 0 || len(resp.Reservations[0].Instances) == 0 {
		return nil, cloudrepo.NotFoundError
	}
	var state string
	stateMp := map[string]string{
		"pending":       cloudrepo.CloudVMStatusPending,
		"running":       cloudrepo.CloudVMStatusRunning,
		"stopping":      cloudrepo.CloudVMStatusStopping,
		"shutting-down": cloudrepo.CloudVMStatusShuttingDown,
		"terminated":    cloudrepo.CloudVMStatusTerminated,
		"stopped":       cloudrepo.CloudVMStatusStopped,
	}
	instance := resp.Reservations[0].Instances[0]
	state = stateMp[string(instance.State.Name)]
	return &cloudrepo.DescribeCloudServerReply{
		InstanceID:    *instance.InstanceId,
		InstanceState: state,
		ImageID:       *instance.ImageId,
		InstanceType:  string(instance.InstanceType),
	}, nil
}
