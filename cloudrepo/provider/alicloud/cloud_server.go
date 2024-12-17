// Package alicloud
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
package alicloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/bilibili/HCP/cloudrepo"
	"github.com/go-kratos/kratos/pkg/log"
)

// CloudServer 云服务器
type CloudServer struct {
	cloudrepo.CloudProductCommon

	SubnetCid             string
	SeverSecurityGroupCid []string
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

// GetPrivateIP 获取云服务器内网IP
func (c *CloudServer) GetPrivateIP() string {
	return c.PrivateIP
}

// GetPublicIP 获取云服务器公网IP
func (c *CloudServer) GetPublicIP() string {
	return c.PublicIP
}

// GetExpireTime 获取云服务器过期时间
func (c *CloudServer) GetExpireTime() string {
	return c.ExpireTime
}

// GetSecurityGroupCid 获取云服务器安全组CID
func (c *CloudServer) GetSecurityGroupCid() []string {
	return c.SeverSecurityGroupCid
}

// GetSubnetCid 获取云服务器子网CID
func (c *CloudServer) GetSubnetCid() string {
	return c.SubnetCid
}

// GetVpcCid 获取云服务器专有网络CID
func (c *CloudServer) GetVpcCid() string {
	return c.VpcCid
}

// GetImageCid 获取云服务器镜像CID
func (c *CloudServer) GetImageCid() string {
	return c.ImageCid
}

// GetZoneCid 获取云服务器可用区CID
func (c *CloudServer) GetZoneCid() string {
	return c.ZoneCid
}

// GetProjectCid 获取云服务器项目CID
func (c *CloudServer) GetProjectCid() string {
	return c.ProjectCid
}

// ListServer 获取云服务器列表
func (c *Client) ListServer(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServer, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	var data []ecs.Instance
	pageNum, pageSize := 1, 100
	request := ecs.CreateDescribeInstancesRequest()
	if req.ResourceID != "" {
		ids := []string{req.ResourceID}
		b, _ := json.Marshal(ids)
		request.InstanceIds = string(b)
	}
	request.InstanceName = req.ResourceName
	for {
		request.PageNumber = requests.Integer(fmt.Sprintf("%d", pageNum))
		request.PageSize = requests.Integer(fmt.Sprintf("%d", pageSize))
		resp, err := cli.DescribeInstances(request)
		if err != nil {
			return nil, err
		}
		data = append(data, resp.Instances.Instance...)
		if len(data) >= resp.TotalCount {
			break
		}
		pageNum += 1
	}
	return c.list2DoServer(data, req.Region)
}

func (c *Client) list2DoServer(resp []ecs.Instance, region string) (list []cloudrepo.CloudServer, err error) {
	var ids []string
	for _, dt := range resp {
		ids = append(ids, dt.InstanceId)
	}
	//cli, err := c.clientBss(region)
	instResp, err := c.QueryInstanceInfo(region, ids, "")
	if err != nil {
		log.Error("queryServerAli QueryAvailableInstancesById error(%s)", err.Error())
	}
	var instanceInfoMap = make(map[string]bssopenapi.Instance, 0)
	for _, instance := range instResp {
		instanceInfoMap[instance.InstanceID] = instance
	}
	for _, v := range resp {
		_ = v
		renewStatus := ""
		if i, ok := instanceInfoMap[v.InstanceId]; ok {
			renewStatus = i.RenewStatus
		}
		var (
			publicIp          []string
			privateIp         []string
			securityGroupCids []string
			subnetCid         string
			vpcCid            string
			imageCid          string
			zoneCid           string
			projectCid        string
		)
		publicIpMap := map[string]struct{}{}
		privateIpMap := map[string]struct{}{}
		for _, i := range v.PublicIpAddress.IpAddress {
			if i != "" {
				publicIpMap[i] = struct{}{}
			}
		}
		if v.IntranetIp != "" {
			privateIpMap[v.IntranetIp] = struct{}{}
		}
		if v.IntranetIp != "" {
			privateIpMap[v.IntranetIp] = struct{}{}
		}
		for _, i := range v.InnerIpAddress.IpAddress {
			if i != "" {
				privateIpMap[i] = struct{}{}
			}
		}
		for _, i := range v.NetworkInterfaces.NetworkInterface {
			for _, set := range i.PrivateIpSets.PrivateIpSet {
				if set.PrivateIpAddress != "" {
					privateIpMap[set.PrivateIpAddress] = struct{}{}
				}
			}
		}
		for _, i := range v.VpcAttributes.PrivateIpAddress.IpAddress {
			if i != "" {
				privateIpMap[i] = struct{}{}
			}
		}
		for i, _ := range publicIpMap {
			publicIp = append(publicIp, i)
		}
		for i, _ := range privateIpMap {
			privateIp = append(privateIp, i)
		}
		if v.EipAddress.IpAddress != "" {
			publicIp = append(publicIp, v.EipAddress.IpAddress)
		}
		for _, sg := range v.SecurityGroupIds.SecurityGroupId {
			securityGroupCids = append(securityGroupCids, sg)
		}
		subnetCid = v.VpcAttributes.VSwitchId
		vpcCid = v.VpcAttributes.VpcId
		zoneCid = v.ZoneId
		projectCid = v.ResourceGroupId
		imageCid = v.ImageId
		list = append(list, &CloudServer{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.InstanceId,
				Name:   v.InstanceName,
				Status: cloudrepo.GetCloudEnum(cloudrepo.CloudAli, "cloud_server", v.Status),
			},
			ServerSpec:            v.InstanceType,
			ChangeType:            cloudrepo.GetCloudEnum(cloudrepo.CloudAli, "pay_type", v.InstanceChargeType),
			RenewStatus:           cloudrepo.GetCloudEnum(cloudrepo.CloudAli, "renew_type", renewStatus),
			PrivateIP:             strings.Join(privateIp, ", "),
			PublicIP:              strings.Join(publicIp, ", "),
			ExpireTime:            cloudrepo.TimeTrans(v.ExpiredTime),
			SeverSecurityGroupCid: securityGroupCids,
			SubnetCid:             subnetCid,
			VpcCid:                vpcCid,
			ImageCid:              imageCid,
			ZoneCid:               zoneCid,
			ProjectCid:            projectCid,
		})
	}
	return
}

// RebootServer 重启云服务器
func (c *Client) RebootServer(req *cloudrepo.RebootCloudServerReq) (*cloudrepo.RebootCloudServerReply, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	request := ecs.CreateRebootInstanceRequest()
	request.SetScheme("https")
	request.SetConnectTimeout(time.Second * 14)
	request.SetReadTimeout(time.Second * 16)
	request.DryRun = requests.NewBoolean(false)
	request.ForceStop = requests.NewBoolean(req.IsForce)
	request.RegionId = c.Region
	request.InstanceId = req.InstanceID
	resp, err := cli.RebootInstance(request)
	if err != nil {
		return &cloudrepo.RebootCloudServerReply{}, err
	}
	return &cloudrepo.RebootCloudServerReply{
		RequestID: resp.RequestId,
	}, nil
}

// ReinstallServer 重装云服务器
func (c *Client) ReinstallServer(req *cloudrepo.ReinstallCloudServerReq) (*cloudrepo.ReinstallCloudServerReply, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	request := ecs.CreateReplaceSystemDiskRequest()
	request.SetScheme("https")
	request.SetConnectTimeout(time.Second * 14)
	request.SetReadTimeout(time.Second * 16)
	request.RegionId = c.Region
	request.InstanceId = req.InstanceID
	request.ImageId = req.ImageID
	request.KeyPairName = req.KeyPairName
	resp, err := cli.ReplaceSystemDisk(request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.ReinstallCloudServerReply{
		RequestID: resp.RequestId,
	}, nil
}

// RenameServer 重命名云服务器
func (c *Client) RenameServer(req *cloudrepo.RenameCloudServerReq) (*cloudrepo.RenameCloudServerReply, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	request := ecs.CreateModifyInstanceAttributeRequest()
	request.InstanceId = req.InstanceID
	request.InstanceName = req.NewName
	request.HostName = req.NewName
	request.SetScheme("https")
	resp, err := cli.ModifyInstanceAttribute(request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.RenameCloudServerReply{
		RequestID: resp.RequestId,
	}, nil
}

// ChangeServerConfig 变更云服务器配置
func (c *Client) ChangeServerConfig(req *cloudrepo.ChangeConfigCloudServerReq) (*cloudrepo.ChangeConfigCloudServerReply, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	request := ecs.CreateModifyPrepayInstanceSpecRequest()
	request.Scheme = "https"
	request.InstanceId = req.InstanceID
	request.InstanceType = req.InstanceType
	request.AutoPay = requests.NewBoolean(true)
	request.MigrateAcrossZone = requests.NewBoolean(true)
	resp, err := cli.ModifyPrepayInstanceSpec(request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.ChangeConfigCloudServerReply{
		RequestID: resp.RequestId,
	}, nil
}

// DeleteServer 删除云服务器
func (c *Client) DeleteServer(req *cloudrepo.DeleteCloudServerReq) (*cloudrepo.DeleteCloudServerReply, error) {
	// 检查服务器状态-已停止状态才支持清退
	status, err := c.getInstanceStatus(req.InstanceID)
	if err != nil {
		log.Error("Fail to get instance status on DeleteServer: %s", err)
		return nil, err
	}
	if status != cloudrepo.CloudVMStatusStopped {
		log.Warn("DeleteServer: server status is %s expect %s", status, cloudrepo.CloudVMStatusStopped)
		return nil, fmt.Errorf("server status is not %s", cloudrepo.CloudVMStatusStopped)
	}
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	request := ecs.CreateDeleteInstanceRequest()
	request.Force = requests.NewBoolean(req.IsForce)
	request.InstanceId = req.InstanceID
	request.TerminateSubscription = requests.NewBoolean(true)
	resp, err := cli.DeleteInstance(request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.DeleteCloudServerReply{
		RequestID: resp.RequestId,
	}, nil
}

// getInstanceStatus 获取云服务器状态
func (c *Client) getInstanceStatus(instanceID string) (string, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return "", err
	}
	request := ecs.CreateDescribeInstanceStatusRequest()
	request.RegionId = c.Region
	request.InstanceId = &[]string{instanceID}
	resp, err := cli.DescribeInstanceStatus(request)
	if err != nil {
		return "", err
	}
	if len(resp.InstanceStatuses.InstanceStatus) == 0 {
		return "", errors.New("instance not found")
	}
	stateMp := map[string]string{
		"Pending":  cloudrepo.CloudVMStatusPending,
		"Running":  cloudrepo.CloudVMStatusRunning,
		"Starting": cloudrepo.CloudVMStatusStarting,
		"Stopping": cloudrepo.CloudVMStatusStopping,
		"Stopped":  cloudrepo.CloudVMStatusStopped,
	}
	state := stateMp[resp.InstanceStatuses.InstanceStatus[0].Status]
	return state, nil
}

// StartServer 启动云服务器
func (c *Client) StartServer(req *cloudrepo.StartCloudServerReq) (*cloudrepo.StartCloudServerReply, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	request := ecs.CreateStartInstanceRequest()
	request.SetScheme("https")
	request.SetConnectTimeout(time.Second * 14)
	request.SetReadTimeout(time.Second * 16)
	request.DryRun = requests.NewBoolean(false)
	request.RegionId = c.Region
	request.InstanceId = req.InstanceID
	resp, err := cli.StartInstance(request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.StartCloudServerReply{
		RequestID: resp.RequestId,
	}, nil
}

// StopServer 停止云服务器
func (c *Client) StopServer(req *cloudrepo.StopCloudServerReq) (*cloudrepo.StopCloudServerReply, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	request := ecs.CreateStopInstanceRequest()
	request.SetScheme("https")
	request.SetConnectTimeout(time.Second * 14)
	request.SetReadTimeout(time.Second * 16)
	request.DryRun = requests.NewBoolean(false)
	request.ForceStop = requests.NewBoolean(req.IsForce)
	request.RegionId = c.Region
	request.InstanceId = req.InstanceID
	resp, err := cli.StopInstance(request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.StopCloudServerReply{
		RequestID: resp.RequestId,
	}, nil
}

// DescribeServer 获取云服务器信息
func (c *Client) DescribeServer(req *cloudrepo.DescribeCloudServerReq) (*cloudrepo.DescribeCloudServerReply, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	// 将字符串转换为切片
	slice := strings.Split(req.InstanceID, ",")
	// 将切片转换为JSON
	jsonData, err := json.Marshal(slice)
	if err != nil {
		return nil, err
	}
	instanceIds := string(jsonData)
	request := ecs.CreateDescribeInstancesRequest()
	request.RegionId = c.Region
	request.InstanceIds = instanceIds
	resp, err := cli.DescribeInstances(request)
	if IsNotFoundErr(err) || resp == nil || len(resp.Instances.Instance) == 0 {
		return nil, cloudrepo.NotFoundError
	}
	if err != nil {
		return nil, err
	}
	/**
		ChargeType
		  PrePaid：包年包月。
		  PostPaid：按量付费。
		InstanceState
	      Pending：创建中。
	      Running：运行中。
	      Starting：启动中。
	      Stopping：停止中。
	      Stopped：已停止。
	*/
	stateMp := map[string]string{
		"Pending":  cloudrepo.CloudVMStatusPending,
		"Running":  cloudrepo.CloudVMStatusRunning,
		"Starting": cloudrepo.CloudVMStatusStarting,
		"Stopping": cloudrepo.CloudVMStatusStopping,
		"Stopped":  cloudrepo.CloudVMStatusStopped,
	}
	chargeMp := map[string]string{
		"PrePaid":  cloudrepo.CloudVmChargeTypePrePaid,
		"PostPaid": cloudrepo.CloudVmChargeTypePostPaid,
	}
	state := stateMp[resp.Instances.Instance[0].Status]
	chargeType := chargeMp[resp.Instances.Instance[0].InstanceChargeType]
	var isOperation bool //true: 重装操作系统已完成；false: 重装操作系统未完成
	if resp.Instances.Instance[0].OperationLocks.LockReason != nil && len(resp.Instances.Instance[0].OperationLocks.LockReason) < 1 {
		isOperation = true
	}
	return &cloudrepo.DescribeCloudServerReply{
		InstanceID:    resp.Instances.Instance[0].InstanceId,
		InstanceName:  resp.Instances.Instance[0].InstanceName,
		ImageID:       resp.Instances.Instance[0].ImageId,
		InstanceState: state,
		InstanceType:  resp.Instances.Instance[0].InstanceType,
		ChargeType:    chargeType,
		IsOperation:   isOperation,
	}, nil
}

// ChangeServerChargeType 变更云服务器计费类型
func (c *Client) ChangeServerChargeType(req *cloudrepo.ChangeServerChargeTypeReq) (*cloudrepo.ChangeServerChargeTypeReply, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	request := ecs.CreateModifyInstanceChargeTypeRequest()
	request.SetScheme("https")
	request.SetConnectTimeout(time.Second * 14)
	request.SetReadTimeout(time.Second * 16)
	request.DryRun = requests.NewBoolean(false)
	request.RegionId = c.Region
	request.InstanceIds = fmt.Sprintf(`["%s"]`, req.InstanceID)
	request.InstanceChargeType = req.ChargeType
	resp, err := cli.ModifyInstanceChargeType(request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.ChangeServerChargeTypeReply{
		RequestID: resp.RequestId,
	}, nil
}

// IsNotFoundErr 判断是否为找不到错误
func IsNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	if err.Error() == "NotFoundError" {
		return true
	}
	return false
}
