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
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/bilibili/HCP/cloudrepo"
	"github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
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

// GetPrivateIP 获取云服务器内网IP
func (c *CloudServer) GetPrivateIP() string {
	return c.PrivateIP
}

// GetPublicIP 获取云服务器公网IP
func (c *CloudServer) GetPublicIP() string {
	return c.PublicIP
}

// GetExpireTime 获取云服务器到期时间
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

// ListServer 获取云服务器列表
func (c *Client) ListServer(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServer, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	var data []Instance
	pageNum, pageSize := 1, 100
	request := cvm.NewDescribeInstancesRequest()
	if req.ResourceID != "" {
		request.InstanceIds = []*string{&req.ResourceID}
	}
	for {
		offset := int64((pageNum - 1) * pageSize)
		limit := int64(pageSize)
		request.Offset = &offset
		request.Limit = &limit
		resp, err := cli.DescribeInstances(request)
		if err != nil {
			return nil, err
		}
		response := &struct {
			Response struct {
				TotalCount  int64      `json:"TotalCount,omitempty" name:"TotalCount"`
				InstanceSet []Instance `json:"InstanceSet,omitempty" name:"InstanceSet"`
				RequestId   string     `json:"RequestId,omitempty" name:"RequestId"`
			} `json:"Response"`
		}{}
		err = json.Unmarshal([]byte(resp.ToJsonString()), response)
		data = append(data, response.Response.InstanceSet...)
		if len(data) >= int(*resp.Response.TotalCount) {
			break
		}
		pageNum += 1
	}
	return list2DoServer(data)
}

func list2DoServer(resp []Instance) (list []cloudrepo.CloudServer, err error) {
	for _, v := range resp {
		_ = v
		var (
			intraIps          []string
			extraIps          []string
			securityGroupCids []string
			subnetCid         string
			vpcCid            string
			imageCid          string
			zoneCid           string
			projectCid        string
		)
		for idx, _ := range v.PublicIpAddresses {
			extraIps = append(extraIps, v.PublicIpAddresses[idx])
		}
		for idx, _ := range v.PrivateIpAddresses {
			intraIps = append(intraIps, v.PrivateIpAddresses[idx])
		}
		for _, sg := range v.SecurityGroupIds {
			securityGroupCids = append(securityGroupCids, sg)
		}
		subnetCid = v.VirtualPrivateCloud.SubnetId
		vpcCid = v.VirtualPrivateCloud.VpcId
		imageCid = v.ImageId
		zoneCid = v.Placement.Zone
		projectCid = strconv.FormatInt(v.Placement.ProjectId, 10)
		list = append(list, &CloudServer{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.InstanceId,
				Name:   v.InstanceName,
				Status: cloudrepo.GetCloudEnum(cloudrepo.CloudTencent, "cloud_server", v.InstanceState),
			},
			ServerSpec:            v.InstanceType,
			ChangeType:            cloudrepo.GetCloudEnum(cloudrepo.CloudTencent, "pay_type", v.InstanceChargeType),
			RenewStatus:           cloudrepo.GetCloudEnum(cloudrepo.CloudTencent, "renew_type", v.RenewFlag),
			PrivateIP:             strings.Join(intraIps, ","),
			PublicIP:              strings.Join(extraIps, ","),
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
	request := cvm.NewRebootInstancesRequest()
	request.InstanceIds = common.StringPtrs([]string{req.InstanceID})
	if req.IsForce {
		request.StopType = common.StringPtr("HARD")
	} else {
		request.StopType = common.StringPtr("SOFT")
	}
	resp, err := cli.RebootInstancesWithContext(context.TODO(), request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.RebootCloudServerReply{
		RequestID: fmt.Sprintf("%v", resp.Response.RequestId),
	}, nil
}

// ReinstallServer 重装云服务器
func (c *Client) ReinstallServer(req *cloudrepo.ReinstallCloudServerReq) (*cloudrepo.ReinstallCloudServerReply, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	// 获取秘钥
	keyRequest := cvm.NewDescribeKeyPairsRequest()
	keyRequest.Filters = []*cvm.Filter{
		&cvm.Filter{
			Name:   common.StringPtr("key-name"),
			Values: []*string{common.StringPtr(req.KeyPairName)},
		},
	}
	keyResponse, err := cli.DescribeKeyPairs(keyRequest)
	if err != nil {
		return nil, errors.Errorf("获取秘钥报错：%v", err)
	}
	if len(keyResponse.Response.KeyPairSet) != 1 {
		return nil, errors.Errorf("获取秘钥错误，不存在或不唯一")
	}
	keyId := *(keyResponse.Response.KeyPairSet[0].KeyId)
	// 重装
	request := cvm.NewResetInstanceRequest()
	request.InstanceId = common.StringPtr(req.InstanceID)
	request.ImageId = common.StringPtr(req.ImageID)
	request.LoginSettings = &cvm.LoginSettings{
		KeyIds: common.StringPtrs([]string{keyId}),
	}
	resp, err := cli.ResetInstanceWithContext(context.TODO(), request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.ReinstallCloudServerReply{
		RequestID: fmt.Sprintf("%v", resp.Response.RequestId),
	}, nil
}

// RenameServer 重命名云服务器
func (c *Client) RenameServer(req *cloudrepo.RenameCloudServerReq) (*cloudrepo.RenameCloudServerReply, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	request := cvm.NewModifyInstancesAttributeRequest()
	request.InstanceIds = common.StringPtrs([]string{req.InstanceID})
	request.InstanceName = &req.NewName
	resp, err := cli.ModifyInstancesAttribute(request)
	if err != nil {
		return nil, err
	}
	if resp.Response == nil || resp.Response.RequestId == nil {
		return nil, errors.New("response for modify tencentcloud instance is nil")
	}
	return &cloudrepo.RenameCloudServerReply{
		RequestID: fmt.Sprintf("%v", resp.Response.RequestId),
	}, nil
}

// ChangeServerConfig 变更云服务器配置
func (c *Client) ChangeServerConfig(req *cloudrepo.ChangeConfigCloudServerReq) (*cloudrepo.ChangeConfigCloudServerReply, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	request := cvm.NewResetInstancesTypeRequest()
	request.InstanceIds = common.StringPtrs([]string{req.InstanceID})
	request.InstanceType = &req.InstanceType
	request.ForceStop = common.BoolPtr(true)
	resp, err := cli.ResetInstancesType(request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.ChangeConfigCloudServerReply{
		RequestID: fmt.Sprintf("%v", resp.Response.RequestId),
	}, nil
}

// DeleteServer 删除云服务器
func (c *Client) DeleteServer(req *cloudrepo.DeleteCloudServerReq) (*cloudrepo.DeleteCloudServerReply, error) {
	cli, err := c.clientCvm()
	if err != nil {
		return nil, err
	}
	request := cvm.NewTerminateInstancesRequest()
	request.InstanceIds = common.StringPtrs([]string{req.InstanceID})
	request.ReleasePrepaidDataDisks = common.BoolPtr(true)
	resp, err := cli.TerminateInstancesWithContext(context.TODO(), request)
	return &cloudrepo.DeleteCloudServerReply{
		RequestID: fmt.Sprintf("%v", resp.Response.RequestId),
	}, err
}

// StartServer 启动云服务器
func (c *Client) StartServer(req *cloudrepo.StartCloudServerReq) (*cloudrepo.StartCloudServerReply, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	request := cvm.NewStartInstancesRequest()
	request.InstanceIds = common.StringPtrs([]string{req.InstanceID})
	resp, err := cli.StartInstancesWithContext(context.TODO(), request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.StartCloudServerReply{
		RequestID: fmt.Sprintf("%v", resp.Response.RequestId),
	}, nil
}

// StopServer 停止云服务器
func (c *Client) StopServer(req *cloudrepo.StopCloudServerReq) (*cloudrepo.StopCloudServerReply, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	request := cvm.NewStopInstancesRequest()
	request.InstanceIds = common.StringPtrs([]string{req.InstanceID})
	if req.IsForce {
		request.StopType = common.StringPtr("HARD")
	} else {
		request.StopType = common.StringPtr("SOFT")
	}

	resp, err := cli.StopInstancesWithContext(context.TODO(), request)
	return &cloudrepo.StopCloudServerReply{
		RequestID: fmt.Sprintf("%v", resp.Response.RequestId),
	}, nil
}

// DescribeServer 获取云服务器信息
func (c *Client) DescribeServer(req *cloudrepo.DescribeCloudServerReq) (*cloudrepo.DescribeCloudServerReply, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	request := cvm.NewDescribeInstancesRequest()
	request.InstanceIds = common.StringPtrs([]string{req.InstanceID})
	resp, err := cli.DescribeInstances(request)
	if err != nil {
		return nil, err
	}
	if err != nil || resp.Response == nil || len(resp.Response.InstanceSet) == 0 || *resp.Response.InstanceSet[0].InstanceState == "SHUTDOWN" {
		return nil, cloudrepo.NotFoundError
	}
	stateMp := map[string]string{
		"PENDING":  cloudrepo.CloudVMStatusPending,
		"RUNNING":  cloudrepo.CloudVMStatusRunning,
		"STARTING": cloudrepo.CloudVMStatusStarting,
		"STOPPING": cloudrepo.CloudVMStatusStopping,
		"STOPPED":  cloudrepo.CloudVMStatusStopped,
	}
	chargeMp := map[string]string{
		"PREPAID":          cloudrepo.CloudVmChargeTypePrePaid,
		"POSTPAID_BY_HOUR": cloudrepo.CloudVmChargeTypePostPaid,
	}
	state := stateMp[*resp.Response.InstanceSet[0].InstanceState]
	chargeType := chargeMp[*resp.Response.InstanceSet[0].InstanceChargeType]
	latestOperationState := ""
	if len(resp.Response.InstanceSet) > 0 && resp.Response.InstanceSet[0].LatestOperationState != nil {
		latestOperationState = *resp.Response.InstanceSet[0].LatestOperationState
	}
	return &cloudrepo.DescribeCloudServerReply{
		InstanceID:           *resp.Response.InstanceSet[0].InstanceId,
		InstanceName:         *resp.Response.InstanceSet[0].InstanceName,
		ImageID:              *resp.Response.InstanceSet[0].ImageId,
		InstanceState:        state,
		InstanceType:         *resp.Response.InstanceSet[0].InstanceType,
		ChargeType:           chargeType,
		LatestOperationState: latestOperationState,
	}, nil
}
