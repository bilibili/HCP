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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/bilibili/HCP/cloudrepo"
	"github.com/bilibili/HCP/utils/ip"
	"github.com/go-kratos/kratos/pkg/log"
	bssModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/bss/v2/model"
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	"github.com/pkg/errors"
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

// GetPublicIP 获取云服务器外网IP
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
	// 部分Region实例化client会报错，这里进行异常捕获
	defer func() {
		e := recover()
		if e != nil {
			err := errors.New(fmt.Sprintf("%v", e))
			if err != nil {
				return
			}
		}
	}()
	var cli *ecs.EcsClient
	var data []model.ServerDetail
	ep, ok := ecsEndPointMap[req.Region]
	if ok {
		projects, err := c.ListProjectIam(req)
		if err != nil {
			return nil, err
		}
		for _, p := range projects {
			if p.Name != req.Region {
				continue
			}
			cli = c.clientEcsWithEndPoint(p.Id, ep)
			servers, err := c.queryServer(req, cli)
			if err != nil {
				return nil, err
			}
			for idx, _ := range servers {
				data = append(data, servers[idx])
			}
		}
		return c.list2DoServer(data)
	} else {
		cli = c.clientEcs()
		resp, err := c.queryServer(req, cli)
		if err != nil {
			log.Error("queryServerHuawei error(%s)", err.Error())
			if strings.Contains(err.Error(), "could not be found") || strings.Contains(err.Error(), "itemNotFound") {
				return c.list2DoServer(data)
			}
			return nil, err
		}
		return c.list2DoServer(resp)
	}
}

func (c *Client) list2DoServer(resp []model.ServerDetail) (list []cloudrepo.CloudServer, err error) {
	var ids []string
	for _, dt := range resp {
		ids = append(ids, dt.Id)
	}
	for _, v := range resp {
		_ = v
		var (
			privateIp         []string
			publicIp          []string
			securityGroupCids []string
			subnetCid         string
			vpcCid            string
			imageCid          string
			zoneCid           string
			projectCid        string
		)
		instanceInfoResp, err := c.ListPayPerUseCustomerResources(ids)
		if err != nil {
			log.Error("queryServerHuawei ListPayPerUseCustomerResources error(%s)", err.Error())
		}
		var instanceInfoMap = make(map[string]bssModel.OrderInstanceV2)
		for _, info := range instanceInfoResp {
			if *info.ResourceId == "" {
				continue
			}
			instanceInfoMap[*info.ResourceId] = info
		}
		expireAt, chargeType, renewStatus := "", "", ""
		if ct, ok := v.Metadata["charging_mode"]; ok {
			chargeType = cloudrepo.GetCloudEnum(cloudrepo.CloudHuawei, "pay_type", ct)
		}
		if info, ok := instanceInfoMap[v.Id]; ok {
			expireAt = *info.ExpireTime
			renewStatus = cloudrepo.GetCloudEnum(cloudrepo.CloudHuawei, "renew_type", strconv.Itoa(int(*info.ExpirePolicy)))
		}
		for _, sg := range v.SecurityGroups {
			securityGroupCids = append(securityGroupCids, sg.Id)
		}
		if vpcId, ok := v.Metadata["vpc_id"]; ok {
			vpcCid = vpcId
		}
		for _, address := range v.Addresses {
			for _, addr := range address {
				if !ip.Ipv4Parser.CheckIP(addr.Addr) {
					continue
				}
				if ip.Ipv4Parser.IsInnerIp(addr.Addr) {
					privateIp = append(privateIp, addr.Addr)
				} else {
					publicIp = append(publicIp, addr.Addr)
				}
			}
		}
		imageCid = v.Image.Id
		zoneCid = v.OSEXTAZavailabilityZone
		projectCid = *v.EnterpriseProjectId
		list = append(list, &CloudServer{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.Id,
				Name:   v.Name,
				Status: cloudrepo.GetCloudEnum(cloudrepo.CloudHuawei, "cloud_server", v.Status),
			},
			ServerSpec:            v.Flavor.Id,
			ChangeType:            chargeType,
			RenewStatus:           renewStatus,
			PrivateIP:             strings.Join(privateIp, ","),
			PublicIP:              strings.Join(publicIp, ","),
			ExpireTime:            expireAt,
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

// queryServer 查询云服务器
func (c *Client) queryServer(req *cloudrepo.GetCloudProductReq, cli *ecs.EcsClient) (data []model.ServerDetail, err error) {
	pageNum, pageSize := 1, 100
	request := model.ListServersDetailsRequest{
		EnterpriseProjectId: nil,
		Flavor:              nil,
		Ip:                  nil,
		Limit:               nil,
		Name:                nil,
		NotTags:             nil,
		Offset:              nil,
		ReservationId:       nil,
		Status:              nil,
		Tags:                nil,
	}
	if req.ResourceID != "" {
		ResourceName, err := c.queryServerAttribute(req, cli)
		if err != nil {
			return nil, err
		}
		request.Name = &ResourceName
	}
	for {
		limit := int32(pageSize)
		offset := int32(pageNum)
		request.Limit = &limit
		request.Offset = &offset
		resp, err := cli.ListServersDetails(&request)
		if err != nil {
			return nil, err
		}
		data = append(data, *resp.Servers...)
		if len(data) >= int(*resp.Count) {
			break
		}
		pageNum += 1
	}
	return data, nil
}

// queryServerAttribute 查询云服务器属性
func (c *Client) queryServerAttribute(req *cloudrepo.GetCloudProductReq, cli *ecs.EcsClient) (resourceName string, err error) {
	request := model.NovaShowServerRequest{
		ServerId:            req.ResourceID,
		OpenStackAPIVersion: nil,
	}
	resp, err := cli.NovaShowServer(&request)
	if err != nil {
		return "", err
	}
	return resp.Server.Name, nil
}

// RebootServer 重启云服务器
func (c *Client) RebootServer(req *cloudrepo.RebootCloudServerReq) (*cloudrepo.RebootCloudServerReply, error) {
	client := c.clientNewEcs(req.ProjectID)
	request := &model.BatchRebootServersRequest{}
	var listServersReboot = []model.ServerId{
		{
			Id: req.InstanceID,
		},
	}
	rebootbody := &model.BatchRebootSeversOption{
		Servers: listServersReboot,
	}
	if req.IsForce {
		rebootbody.Type = model.GetBatchRebootSeversOptionTypeEnum().HARD
	} else {
		rebootbody.Type = model.GetBatchRebootSeversOptionTypeEnum().SOFT
	}
	request.Body = &model.BatchRebootServersRequestBody{
		Reboot: rebootbody,
	}
	resp, err := client.BatchRebootServers(request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.RebootCloudServerReply{
		RequestID: fmt.Sprintf("%v", *resp.JobId),
	}, nil
}

// ReinstallServer 重装云服务器系统
func (c *Client) ReinstallServer(req *cloudrepo.ReinstallCloudServerReq) (*cloudrepo.ReinstallCloudServerReply, error) {
	client := c.clientNewEcs(req.ProjectID)
	request := &model.ChangeServerOsWithCloudInitRequest{}
	request.ServerId = req.InstanceID
	keyName := req.KeyPairName
	osChangeBody := &model.ChangeServerOsWithCloudInitOption{
		Keyname: &keyName,
		Imageid: req.ImageID,
	}
	request.Body = &model.ChangeServerOsWithCloudInitRequestBody{
		OsChange: osChangeBody,
	}
	resp, err := client.ChangeServerOsWithCloudInit(request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.ReinstallCloudServerReply{
		RequestID: fmt.Sprintf("%v", *resp.JobId),
	}, nil
}

// RenameServer 重命名云服务器
func (c *Client) RenameServer(req *cloudrepo.RenameCloudServerReq) (*cloudrepo.RenameCloudServerReply, error) {
	client := c.clientNewEcs(req.ProjectID)
	request := &model.UpdateServerRequest{}
	request.ServerId = req.InstanceID
	request.Body = &model.UpdateServerRequestBody{
		Server: &model.UpdateServerOption{
			Name:        &req.NewName,
			Description: nil,
		},
	}
	resp, err := client.UpdateServer(request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.RenameCloudServerReply{
		RequestID: resp.String(),
	}, nil
}

// ChangeServerConfig 变更云服务器配置
func (c *Client) ChangeServerConfig(req *cloudrepo.ChangeConfigCloudServerReq) (*cloudrepo.ChangeConfigCloudServerReply, error) {
	client := c.clientNewEcs(req.ProjectID)
	request := &model.ResizeServerRequest{}
	isAutoPay := "true"
	request.ServerId = req.InstanceID
	request.Body = &model.ResizeServerRequestBody{
		Resize: &model.ResizePrePaidServerOption{
			FlavorRef: req.InstanceType,
			Extendparam: &model.ResizeServerExtendParam{
				IsAutoPay: &isAutoPay,
			},
		},
	}
	resp, err := client.ResizeServer(request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.ChangeConfigCloudServerReply{
		RequestID: fmt.Sprintf("%v", *resp.JobId),
	}, nil
}

// DeleteServer 删除云服务器
func (c *Client) DeleteServer(req *cloudrepo.DeleteCloudServerReq) (*cloudrepo.DeleteCloudServerReply, error) {
	client := c.clientNewEcs(req.ProjectID)
	request := &model.ShowServerRequest{}
	request.ServerId = req.InstanceID
	resp, err := client.ShowServer(request)
	if err != nil {
		return nil, err
	}
	var requestID string
	if chargingMode, find := resp.Server.Metadata["charging_mode"]; find {
		switch chargingMode {
		case "0":
			requestID, err = c.DeletePostPaidServer(req)
		case "1":
			requestID, err = c.DeletePrePaidServer(req)
		default:
			err = errors.New("not support the charging mode")
		}
	}
	return &cloudrepo.DeleteCloudServerReply{
		RequestID: requestID,
	}, err
}

// DeletePostPaidServer 删除按需云服务器
func (c *Client) DeletePostPaidServer(req *cloudrepo.DeleteCloudServerReq) (string, error) {
	client := c.clientNewEcs(req.ProjectID)
	DeletePublicIp, DeleteVolume := true, true
	listServersReboot := []model.ServerId{
		{
			Id: req.InstanceID,
		},
	}
	request := &model.DeleteServersRequest{
		Body: &model.DeleteServersRequestBody{
			DeletePublicip: &DeletePublicIp,
			DeleteVolume:   &DeleteVolume,
			Servers:        listServersReboot,
		},
	}
	resp, err := client.DeleteServers(request)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", *resp.JobId), nil
}

// DeletePrePaidServer 删除包年/包月云服务器
func (c *Client) DeletePrePaidServer(req *cloudrepo.DeleteCloudServerReq) (string, error) {
	client := c.clientBss(c.Region)
	var request = &bssModel.CancelResourcesSubscriptionRequest{
		Body: &bssModel.UnsubscribeResourcesReq{
			ResourceIds:     []string{req.InstanceID},
			UnsubscribeType: 1,
		},
	}
	resp, err := client.CancelResourcesSubscription(request)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", resp.String()), nil
}

// StartServer 启动云服务器
func (c *Client) StartServer(req *cloudrepo.StartCloudServerReq) (*cloudrepo.StartCloudServerReply, error) {
	client := c.clientNewEcs(req.ProjectID)
	request := &model.BatchStartServersRequest{}
	var listServersOsStart = []model.ServerId{
		{
			Id: req.InstanceID,
		},
	}
	osStartBody := &model.BatchStartServersOption{
		Servers: listServersOsStart,
	}
	request.Body = &model.BatchStartServersRequestBody{
		OsStart: osStartBody,
	}
	resp, err := client.BatchStartServers(request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.StartCloudServerReply{
		RequestID: fmt.Sprintf("%v", *resp.JobId),
	}, nil
}

// StopServer 停止云服务器
func (c *Client) StopServer(req *cloudrepo.StopCloudServerReq) (*cloudrepo.StopCloudServerReply, error) {
	client := c.clientNewEcs(req.ProjectID)
	request := &model.BatchStopServersRequest{}
	var listServersOsStop = []model.ServerId{
		{
			Id: req.InstanceID,
		},
	}
	var typeOsStop model.BatchStopServersOptionType
	if req.IsForce {
		typeOsStop = model.GetBatchStopServersOptionTypeEnum().SOFT
	} else {
		typeOsStop = model.GetBatchStopServersOptionTypeEnum().SOFT
	}

	osStopBody := &model.BatchStopServersOption{
		Servers: listServersOsStop,
		Type:    &typeOsStop,
	}
	request.Body = &model.BatchStopServersRequestBody{
		OsStop: osStopBody,
	}
	resp, err := client.BatchStopServers(request)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.StopCloudServerReply{
		RequestID: fmt.Sprintf("%v", *resp.JobId),
	}, nil
}

// DescribeServer 查询云服务器信息
func (c *Client) DescribeServer(req *cloudrepo.DescribeCloudServerReq) (*cloudrepo.DescribeCloudServerReply, error) {
	client := c.clientNewEcs(req.ProjectID)
	request := &model.ShowServerRequest{}
	request.ServerId = req.InstanceID
	resp, err := client.ShowServer(request)
	if IsNotFoundErr(err) || (resp != nil && resp.Server.Status == "DELETED") {
		return nil, cloudrepo.NotFoundError
	}
	stateMp := map[string]string{
		"SHUTOFF":     cloudrepo.CloudVMStatusStopped,
		"ACTIVE":      cloudrepo.CloudVMStatusRunning,
		"powering-on": cloudrepo.CloudVMStatusStarting,
		"REBOOT":      cloudrepo.CloudVMStatusStarting,
		"DELETE":      cloudrepo.CloudVMStatueDeleted,
	}
	state := stateMp[resp.Server.Status]
	return &cloudrepo.DescribeCloudServerReply{
		InstanceID:    resp.Server.Id,
		InstanceName:  resp.Server.Name,
		ImageID:       resp.Server.Image.Id,
		InstanceState: state,
		InstanceType:  resp.Server.Flavor.Name,
	}, nil
}

// IsNotFoundErr 判断是否为不存在错误
func IsNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	tmpErr := struct {
		ErrorCode string `json:"error_code,omitempty"`
	}{}
	json.Unmarshal([]byte(err.Error()), &tmpErr)
	if tmpErr.ErrorCode == "Ecs.0114" {
		return true
	}
	return false
}
