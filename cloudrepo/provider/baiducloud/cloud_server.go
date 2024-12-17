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
	"fmt"
	"strconv"
	"strings"

	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/bilibili/HCP/app/interface/v1/configs"
	"github.com/bilibili/HCP/cloudrepo"
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

// GetPrivateIP 获取云服务器内网ip
func (c *CloudServer) GetPrivateIP() string {
	return c.PrivateIP
}

// GetPublicIP 获取云服务器外网ip
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
	cli, err := c.clientBcc()
	if err != nil {
		return nil, err
	}
	maxItems := 1000
	var marker = new(string)
	*marker = ""
	request := &api.ListInstanceArgs{
		MaxKeys:     maxItems,
		InstanceIds: req.ResourceID,
	}
	var data []api.InstanceModel
	for {
		request.Marker = *marker
		resp, err := cli.ListInstances(request)
		if err != nil {
			return nil, err
		}
		data = append(data, resp.Instances...)
		marker = &resp.NextMarker
		if !resp.IsTruncated {
			break
		}
	}
	return list2DoServer(data)
}

func list2DoServer(resp []api.InstanceModel) (list []cloudrepo.CloudServer, err error) {
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
		for _, i := range v.NicInfo.Ips {
			privateIp = append(privateIp, i.PrivateIp)
		}
		for _, i := range v.NicInfo.Ips {
			publicIp = append(publicIp, i.Eip)
		}
		for _, i := range v.NicInfo.SecurityGroups {
			securityGroupCids = append(securityGroupCids, i)
		}
		subnetCid = v.SubnetId
		vpcCid = v.VpcId
		imageCid = v.ImageId
		zoneCid = v.ZoneName
		for _, tag := range v.Tags {
			if tag.TagKey == configs.Conf.CloudConf.TagProjectKey {
				projectCid = fmt.Sprintf("%s-%s", tag.TagKey, tag.TagValue)
			}
		}
		list = append(list, &CloudServer{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.InstanceId,
				Name:   v.InstanceName,
				Status: cloudrepo.GetCloudEnum(cloudrepo.CloudBaidu, "cloud_server", string(v.Status)),
			},
			ServerSpec:            v.Spec,
			ChangeType:            cloudrepo.GetCloudEnum(cloudrepo.CloudBaidu, "pay_type", v.PaymentTiming),
			RenewStatus:           cloudrepo.GetCloudEnum(cloudrepo.CloudBaidu, "renew_type", strconv.FormatBool(v.AutoRenew)),
			PrivateIP:             strings.Join(privateIp, ","),
			PublicIP:              strings.Join(publicIp, ","),
			ExpireTime:            cloudrepo.TimeTrans(v.ExpireTime),
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

func (c *Client) RebootServer(req *cloudrepo.RebootCloudServerReq) (*cloudrepo.RebootCloudServerReply, error) {
	cli, err := c.clientBcc()
	if err != nil {
		return nil, err
	}
	err = cli.RebootInstance(req.InstanceID, req.IsForce)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.RebootCloudServerReply{
		RequestID: req.InstanceID,
	}, nil
}

func (c *Client) ReinstallServer(req *cloudrepo.ReinstallCloudServerReq) (*cloudrepo.ReinstallCloudServerReply, error) {
	cli, err := c.clientBcc()
	if err != nil {
		return nil, err
	}
	IsOpenHostEye := true
	IsPreserveData := true
	// 指定重装时所需密钥对ID
	keyPairCid := req.KeyPairID
	if keyPairCid == "" {
		return nil, errors.New("key pair not found")
	}
	args := &api.RebuildInstanceArgsV2{
		// 指定的镜像ID
		ImageId: req.ImageID,
		// 待重装实例所要绑定的密钥对ID，必须传递adminPass、keypairId其中一个参数
		KeypairId: keyPairCid,
		// 是否开启主机安全，默认开启
		IsOpenHostEye: &IsOpenHostEye,
		// 是否保留数据重装，默认保留
		IsPreserveData: &IsPreserveData,
	}
	err = cli.RebuildInstanceV2(req.InstanceID, args)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.ReinstallCloudServerReply{
		RequestID: req.InstanceID,
	}, nil
}

func (c *Client) RenameServer(req *cloudrepo.RenameCloudServerReq) (*cloudrepo.RenameCloudServerReply, error) {
	cli, err := c.clientBcc()
	if err != nil {
		return nil, err
	}
	args := &api.ModifyInstanceAttributeArgs{
		Name: req.NewName,
	}
	err = cli.ModifyInstanceAttribute(req.InstanceID, args)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.RenameCloudServerReply{
		RequestID: req.InstanceID,
	}, nil
}

func (c *Client) ChangeServerConfig(req *cloudrepo.ChangeConfigCloudServerReq) (*cloudrepo.ChangeConfigCloudServerReply, error) {
	cli, err := c.clientBcc()
	if err != nil {
		return nil, err
	}
	resizeInstanceArgs := &api.ResizeInstanceArgs{
		Spec: req.InstanceType,
	}
	err = cli.ResizeInstanceBySpec(req.InstanceID, resizeInstanceArgs)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.ChangeConfigCloudServerReply{
		RequestID: req.InstanceID,
	}, nil
}

func (c *Client) DeleteServer(req *cloudrepo.DeleteCloudServerReq) (*cloudrepo.DeleteCloudServerReply, error) {
	cli, err := c.clientBcc()
	if err != nil {
		return nil, err
	}
	detail, err := cli.GetInstanceDetail(req.InstanceID)
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, err
	}
	switch detail.Instance.PaymentTiming {
	case "Prepaid":
		deletePrepaidInstanceWithRelateResourceArgs := &api.DeletePrepaidInstanceWithRelateResourceArgs{
			// 批量释放的实例id
			InstanceId: req.InstanceID,
			// 是否关联释放当前实例挂载的EIP和数据盘(只有该字段为true时 DeleteCdsSnapshotFlag字段才会有效，若该字段为false,DeleteCdsSnapshotFlag字段的值无效）
			RelatedReleaseFlag: true,
			// 设置是否释放云磁盘快照
			DeleteCdsSnapshotFlag: true,
			// 设置是否释放弹性网卡
			DeleteRelatedEnisFlag: true,
		}
		resp, err := cli.DeletePrepaidInstanceWithRelateResource(deletePrepaidInstanceWithRelateResourceArgs)
		if err != nil {
			return nil, err
		}
		if !resp.InstanceRefundFlag {
			return nil, errors.New("instance release failed")
		}
	case "Postpaid":
		err = cli.BatchDeleteInstanceWithRelateResource(&api.BatchDeleteInstanceWithRelateResourceArgs{
			RelatedReleaseFlag:    true,
			DeleteCdsSnapshotFlag: true,
			BccRecycleFlag:        true,
			DeleteRelatedEnisFlag: false,
			InstanceIds:           []string{req.InstanceID},
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("DeleteServer should be PrePaid or PostPaid")
	}
	return &cloudrepo.DeleteCloudServerReply{
		RequestID: req.InstanceID,
	}, nil
}

func (c *Client) StartServer(req *cloudrepo.StartCloudServerReq) (*cloudrepo.StartCloudServerReply, error) {
	cli, err := c.clientBcc()
	if err != nil {
		return nil, err
	}
	err = cli.StartInstance(req.InstanceID)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.StartCloudServerReply{
		RequestID: req.InstanceID,
	}, nil
}

func (c *Client) StopServer(req *cloudrepo.StopCloudServerReq) (*cloudrepo.StopCloudServerReply, error) {
	cli, err := c.clientBcc()
	if err != nil {
		return nil, err
	}
	err = cli.StopInstance(req.InstanceID, true)
	if err != nil {
		return nil, err
	}
	return &cloudrepo.StopCloudServerReply{
		RequestID: req.InstanceID,
	}, nil
}

func (c *Client) DescribeServer(req *cloudrepo.DescribeCloudServerReq) (*cloudrepo.DescribeCloudServerReply, error) {
	cli, err := c.clientBcc()
	if err != nil {
		return nil, err
	}
	resp, err := cli.ListInstances(&api.ListInstanceArgs{
		InstanceIds: req.InstanceID,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Instances) == 0 {
		return nil, cloudrepo.NotFoundError
	}
	stateMp := map[string]string{
		"Stopped":  cloudrepo.CloudVMStatusStopped,
		"Running":  cloudrepo.CloudVMStatusRunning,
		"Starting": cloudrepo.CloudVMStatusStarting,
		"Stopping": cloudrepo.CloudVMStatusStopping,
	}
	chargeMp := map[string]string{
		"Prepaid":  cloudrepo.CloudVmChargeTypePrePaid,
		"Postpaid": cloudrepo.CloudVmChargeTypePostPaid,
	}
	instance := resp.Instances[0]
	state := stateMp[string(instance.Status)]
	chargeType := chargeMp[instance.PaymentTiming]
	return &cloudrepo.DescribeCloudServerReply{
		InstanceID:    instance.InstanceId,
		InstanceName:  instance.InstanceName,
		ImageID:       instance.ImageId,
		InstanceType:  instance.Spec,
		InstanceState: state,
		ChargeType:    chargeType,
	}, nil
}
