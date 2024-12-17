// Package worker
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
package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"time"

	"github.com/bilibili/HCP/app/interface/v1/configs"
	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/bilibili/HCP/app/interface/v1/internal/biz/enum"
	"github.com/bilibili/HCP/cloudrepo"
	common "github.com/bilibili/HCP/common/models"
	"github.com/bilibili/HCP/utils/aes"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/pkg/errors"
)

const (
	StartServerWaitTime        = 10 * time.Minute // 启动云服务器等待时间
	StopServerWaitTime         = 10 * time.Minute // 停止云服务器等待时间
	ConfigServerWaitTime       = 30 * time.Minute // 改配云服务器等待时间
	RebootServerWaitTime       = 10 * time.Minute // 重启云服务器等待时间
	RenameServerWaitTime       = 5 * time.Minute  // 改名云服务器等待时间
	ReinstallServerWaitTime    = 10 * time.Minute // 重装云服务器等待时间
	ChangeServerChargeTypeTime = 30 * time.Minute // 变更云服务器计费方式等待时间
	DeleteServerWaitTime       = 10 * time.Minute // 删除云服务器等待时间
)

var (
	NotFoundErrorStr = "-60001"
)

// SyncCloudResource 同步云资源数据
func (oe *ObjectTriggerExecutable) SyncCloudResource(resource *common.SyncCloudResourceRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	err := datasync.SyncCloudResource(ctx, resource)
	if err != nil {
		log.Error("SyncCloudResource instance %s instance_type %s sync failed", resource.ResourceID, resource.ResourceType)
		return fmt.Errorf("SyncCloudResource instance %s instance_type %s sync failed", resource.ResourceID, resource.ResourceType)
	}
	time.Sleep(30 * time.Second)
	return nil
}

// PreOperate 预操作
func (oe *ObjectTriggerExecutable) PreOperate(accountID, serverID int) (*biz.Account, *ResourceOperate, cloudrepo.CloudProvider, error) {
	newAccountID := accountID
	resource, err := oe.GetResource(serverID, enum.ResourceTypeCloudServer)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "cloud_server get region failed")
	}
	if resource.InstanceCID == "" || resource.RegionCID == "" || resource.ProjectCID == "" {
		return nil, nil, nil, errors.New("cloud_server get region and project failed")
	}
	if accountID == 0 {
		newAccountID = resource.AccountID
	}
	accountInfo, err := oe.GetAccountAkSk(newAccountID)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "PreOperate cloud_server get account ak sk failed")
	}
	repo, err := cloudrepo.GetRepo(accountInfo.Provider.Alias)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "cloud_server get repo failed")
	}
	provider, err := repo.GetProvider(&cloudrepo.GetProviderReq{
		Region:    resource.RegionCID,
		SecretId:  accountInfo.OperateSecretID,
		SecretKey: accountInfo.OperateSecretKey,
	})
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "cloud_server get provider failed")
	}
	return accountInfo, resource, provider, nil
}

// GetResource 获取资源数据
func (oe *ObjectTriggerExecutable) GetResource(resourceID int, resourceType string) (*ResourceOperate, error) {
	var (
		accountId int
		region    string
		project   string
		name      string
		cid       string
		imageCid  string
	)
	conditions := map[string]interface{}{
		"id": resourceID,
	}
	switch resourceType {
	case enum.ResourceTypeCloudServer:
		resp, err := cloudServerUseCase.QueryCloudServer(&biz.CloudServerWhere{Conditions: conditions}, &biz.CloudServerOutput{})
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return nil, errors.Wrap(err, "GetResource QueryCloudServer")
		}
		if len(resp) == 0 {
			oe.Status = enum.ExecutableStateFailure
			return nil, errors.New("GetResource QueryCloudServer empty")
		}
		region = resp[0].Region.CID
		project = resp[0].Project.CID
		name = resp[0].Name
		cid = resp[0].CID
		accountId = resp[0].AccountID
	case enum.ResourceTypeCloudServerImage:
		resp, err := cloudServerImageUseCase.QueryCloudServerImage(&biz.CloudServerImageWhere{Conditions: conditions}, &biz.CloudServerImageOutput{})
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return nil, errors.Wrap(err, "GetResource QueryCloudServerImage")
		}
		if len(resp) == 0 {
			oe.Status = enum.ExecutableStateFailure
			return nil, errors.New("GetResource QueryCloudServerImage empty")
		}
		imageCid = resp[0].CID
	default:
		oe.Status = enum.ExecutableStateFailure
		return nil, errors.New("unknown resource type")
	}
	return &ResourceOperate{
		RegionCID:    region,
		ProjectCID:   project,
		InstanceCID:  cid,
		InstanceName: name,
		AccountID:    accountId,
		ImageCID:     imageCid,
	}, nil
}

// GetAccountAkSk 获取账号ak sk
func (oe *ObjectTriggerExecutable) GetAccountAkSk(accountID int) (*biz.Account, error) {
	conditions := map[string]interface{}{
		"id": accountID,
	}
	resp, err := accountUseCase.QueryAccount(&biz.AccountWhere{Conditions: conditions}, &biz.AccountOutput{})
	if err != nil {
		oe.Status = enum.ExecutableStateFailure
		return nil, errors.Wrap(err, "GetAccountAkSk")
	}
	if len(resp) == 0 {
		oe.Status = enum.ExecutableStateFailure
		return nil, errors.New("account not found")
	}
	accessKey := resp[0].OperateSecretID
	secretKey := resp[0].OperateSecretKey
	// secretKey解密
	if secretKey != "" {
		aesSecretKey, err := aes.Decrypt(configs.Conf.CloudSecret.SecretAesKey, secretKey)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return nil, errors.Wrap(err, "GetAccountAkSk")
		}
		secretKey = aesSecretKey
	}
	return &biz.Account{
		KeyPairID:        resp[0].KeyPairID,
		KeyPairName:      resp[0].KeyPairName,
		OperateSecretID:  accessKey,
		OperateSecretKey: secretKey,
		Provider: &biz.Provider{
			ID:    resp[0].Provider.ID,
			Alias: resp[0].Provider.Alias,
		},
		Alias:      resp[0].Alias,
		ProviderID: resp[0].ProviderID,
		CloudProductCommon: biz.CloudProductCommon{
			ID:   resp[0].CloudProductCommon.ID,
			Name: resp[0].CloudProductCommon.Name,
			CID:  resp[0].CloudProductCommon.CID,
		},
	}, nil
}

// StartCloudServerWait 等待云服务器启动完成
func (oe *ObjectTriggerExecutable) StartCloudServerWait(resource *ResourceOperate, provider cloudrepo.CloudProvider) error {
	ctx, cancel := context.WithTimeout(context.Background(), StartServerWaitTime)
	defer cancel()
	reply, err := provider.DescribeServer(&cloudrepo.DescribeCloudServerReq{
		InstanceID: resource.InstanceCID,
		ProjectID:  resource.ProjectCID,
	})
	if err != nil {
		log.Error("StartCloudServerWait instance %s describe server failed err: %v", resource.InstanceCID, err)
		return errors.Wrap(err, "cloud_server describe server failed")
	}
	if reply.InstanceState == cloudrepo.CloudVMStatusRunning {
		return nil
	}
Loop:
	for {
		if reply.InstanceState == cloudrepo.CloudVMStatusStarting {
			break
		}
		resp, err := provider.StartServer(&cloudrepo.StartCloudServerReq{
			InstanceID: resource.InstanceCID,
			ProjectID:  resource.ProjectCID,
		})
		if err == nil {
			break
		}
		select {
		case <-ctx.Done():
			log.Error("StartServer instance %s start timeout", resource.InstanceCID)
			return errors.Errorf("timeout")
		default:
			if strings.Contains(err.Error(), "is Locked by order") {
				log.Warn("StartServer is Locked by order instance %s err: %v, resp: %v", resource.InstanceCID, err, resp)
				time.Sleep(30 * time.Second)
				continue Loop
			} else {
				log.Error("StartServer instance %s start failed err: %v, resp: %v", resource.InstanceCID, err, resp)
				return errors.Wrap(err, "cloud_server start server failed")
			}
		}
	}
	// 等待云服务器启动完成
forLoop:
	for {
		select {
		case <-ctx.Done():
			log.Error("StartServer instance %s start timeout", resource.InstanceCID)
			return errors.Errorf("timeout")
		default:
			reply, err := provider.DescribeServer(&cloudrepo.DescribeCloudServerReq{
				InstanceID: resource.InstanceCID,
				ProjectID:  resource.ProjectCID,
			})
			if err != nil {
				log.Error("StartServer instance %s describe server failed err: %v", resource.InstanceCID, err)
				return errors.Wrap(err, "cloud_server start server describe server failed")
			}
			if reply.InstanceState != cloudrepo.CloudVMStatusRunning {
				time.Sleep(30 * time.Second)
				continue forLoop
			}
			break
		}
		return nil
	}
}

// StopCloudServer 云服务器关机
func (oe *ObjectTriggerExecutable) StopCloudServer(resource *ResourceOperate, provider cloudrepo.CloudProvider) error {
	// 查询云服务器是否为关机状态
	reply, err := provider.DescribeServer(&cloudrepo.DescribeCloudServerReq{
		InstanceID: resource.InstanceCID,
		ProjectID:  resource.ProjectCID,
	})
	if err != nil {
		log.Error("StopCloudServer instance %s describe server failed err: %v", resource.InstanceCID, err)
		return errors.Wrap(err, "cloud_server stop server describe server failed")
	}
	if reply.InstanceState == cloudrepo.CloudVMStatusStopped {
		return nil
	}
	// 云服务器关机
	_, err = provider.StopServer(&cloudrepo.StopCloudServerReq{
		InstanceID: resource.InstanceCID,
		ProjectID:  resource.ProjectCID,
	})
	if err != nil {
		log.Error("StopServer instance %s stop failed err: %v", resource.InstanceCID, err)
		return errors.Wrap(err, "cloud_server stop server failed")
	}
	return nil
}

// StopCloudServerWait 等待云服务器关机完成
func (oe *ObjectTriggerExecutable) StopCloudServerWait(resource *ResourceOperate, provider cloudrepo.CloudProvider) error {
	ctx, cancel := context.WithTimeout(context.Background(), StopServerWaitTime)
	defer cancel()
forLoop:
	for {
		select {
		case <-ctx.Done():
			log.Error("StopServer instance %s stop timeout", resource.InstanceCID)
			return errors.Errorf("timeout")
		default:
			reply, err := provider.DescribeServer(&cloudrepo.DescribeCloudServerReq{
				InstanceID: resource.InstanceCID,
				ProjectID:  resource.ProjectCID,
			})
			if err != nil {
				log.Error("StopServer instance %s describe server failed err: %v", resource.InstanceCID, err)
				return errors.Wrap(err, "cloud_server stop server describe server failed")
			}
			if reply.InstanceState != cloudrepo.CloudVMStatusStopped {
				time.Sleep(30 * time.Second)
				continue forLoop
			}
			break
		}
		return nil
	}
}

// ConfigCloudServer 云服务器改配
func (oe *ObjectTriggerExecutable) ConfigCloudServer(resource *ResourceOperate, provider cloudrepo.CloudProvider, spec string) error {
	_, err := provider.ChangeServerConfig(&cloudrepo.ChangeConfigCloudServerReq{
		InstanceID:   resource.InstanceCID,
		ProjectID:    resource.ProjectCID,
		InstanceType: spec,
	})
	if err != nil {
		oe.Status = enum.ExecutableStateFailure
		log.Error("ChangeServerConfig instance %s change config failed err: %v", resource.InstanceCID, err)
		return errors.Wrap(err, "cloud_server change server config failed")
	}
	return nil
}

// ConfigCloudServerWait 等待云服务器改配完成
func (oe *ObjectTriggerExecutable) ConfigCloudServerWait(resource *ResourceOperate, provider cloudrepo.CloudProvider, instanceType, cloudID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ConfigServerWaitTime)
	defer cancel()
forLoop:
	for {
		select {
		case <-ctx.Done():
			oe.Status = enum.ExecutableStateFailure
			log.Error("StopServer instance %s stop timeout", resource.InstanceCID)
			return errors.Errorf("timeout")
		default:
			reply, err := provider.DescribeServer(&cloudrepo.DescribeCloudServerReq{
				InstanceID: resource.InstanceCID,
				ProjectID:  resource.ProjectCID,
			})
			if err != nil {
				oe.Status = enum.ExecutableStateFailure
				log.Error("StopServer instance %s describe server failed err: %v", resource.InstanceCID, err)
				return errors.Wrap(err, "cloud_server change server describe server failed")
			}
			if reply.InstanceState == "FAILED" {
				oe.Status = enum.ExecutableStateFailure
				log.Error("ConfigServer instance %s modify failed", resource.InstanceCID)
				return errors.Errorf("cloud_server modify failed")
			}
			if reply.InstanceType != instanceType {
				time.Sleep(30 * time.Second)
				continue forLoop
			}
			if cloudID == cloudrepo.CloudTencent && reply.LatestOperationState != "SUCCESS" { //腾讯云判断重装系统是否已完成
				time.Sleep(30 * time.Second)
				continue forLoop
			}
		}
		break
	}
	return nil
}

// RebootCloudServer 云服务器重启
func (oe *ObjectTriggerExecutable) RebootCloudServer(resource *ResourceOperate, provider cloudrepo.CloudProvider) error {
	_, err := provider.RebootServer(&cloudrepo.RebootCloudServerReq{
		InstanceID: resource.InstanceCID,
		ProjectID:  resource.ProjectCID,
	})
	if err != nil {
		log.Error("RebootCloudServer instance %s stop failed err: %v", resource.InstanceCID, err)
		return errors.Wrap(err, "cloud_server reboot server failed")
	}
	return nil
}

// RebootCloudServerWait 等待云服务器重启完成
func (oe *ObjectTriggerExecutable) RebootCloudServerWait(resource *ResourceOperate, provider cloudrepo.CloudProvider, cloudID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), RebootServerWaitTime)
	defer cancel()
forLoop:
	for {
		select {
		case <-ctx.Done():
			log.Error("RebootServer instance %s reboot timeout", resource.InstanceCID)
			return errors.Errorf("timeout")
		default:
			reply, err := provider.DescribeServer(&cloudrepo.DescribeCloudServerReq{
				InstanceID: resource.InstanceCID,
				ProjectID:  resource.ProjectCID,
			})
			if err != nil {
				log.Error("RebootCloudServerWait instance %s describe server failed err: %v", resource.InstanceCID, err)
				return errors.Wrap(err, "cloud_server reboot server describe server failed")
			}
			if reply.InstanceState != cloudrepo.CloudVMStatusRunning {
				time.Sleep(30 * time.Second)
				continue forLoop
			}
			if cloudID == cloudrepo.CloudTencent && reply.LatestOperationState != "SUCCESS" { //腾讯云判断重装系统是否已完成
				time.Sleep(30 * time.Second)
				continue forLoop
			}
			break
		}
		return nil
	}
}

// RenameCloudServer 云服务器改名
func (oe *ObjectTriggerExecutable) RenameCloudServer(resource *ResourceOperate, provider cloudrepo.CloudProvider) error {
	_, err := provider.RenameServer(&cloudrepo.RenameCloudServerReq{
		InstanceID: resource.InstanceCID,
		ProjectID:  resource.ProjectCID,
		NewName:    resource.InstanceName,
	})
	if err != nil {
		log.Error("RenameCloudServer instance %s rename failed err: %v", resource.InstanceCID, err)
		return errors.Wrap(err, "cloud_server rename server failed")
	}
	return nil
}

// RenameCloudServerWait 等待云服务器改名完成·
func (oe *ObjectTriggerExecutable) RenameCloudServerWait(resource *ResourceOperate, provider cloudrepo.CloudProvider) error {
	ctx, cancel := context.WithTimeout(context.Background(), RenameServerWaitTime)
	defer cancel()
forLoop:
	for {
		select {
		case <-ctx.Done():
			log.Error("RebootServer instance %s rename timeout", resource.InstanceCID)
			return errors.Errorf("timeout")
		default:
			reply, err := provider.DescribeServer(&cloudrepo.DescribeCloudServerReq{
				InstanceID: resource.InstanceCID,
				ProjectID:  resource.ProjectCID,
			})
			if err != nil {
				log.Error("RenameCloudServerWait instance %s describe server failed err: %v", resource.InstanceCID, err)
				return errors.Wrap(err, "cloud_server rename server describe server failed")
			}
			if reply.InstanceName != resource.InstanceName {
				time.Sleep(5 * time.Second)
				continue forLoop
			}
			break
		}
		return nil
	}
}

// ReinstallCloudServer 云服务器重装
func (oe *ObjectTriggerExecutable) ReinstallCloudServer(resource *ResourceOperate, provider cloudrepo.CloudProvider, imageCID, accountCID, keyID, keyName string) error {
	_, err := provider.ReinstallServer(&cloudrepo.ReinstallCloudServerReq{
		InstanceID:  resource.InstanceCID,
		ProjectID:   resource.ProjectCID,
		ImageID:     imageCID,
		AccountID:   accountCID,
		KeyPairID:   keyID,
		KeyPairName: keyName,
	})
	if err != nil {
		log.Error("ReinstallCloudServer instance %s reinstall failed err: %v", resource.InstanceCID, err)
		return errors.Wrap(err, "cloud_server reinstall server failed")
	}
	return nil
}

// ReinstallCloudServerWait 等待云服务器重装完成
func (oe *ObjectTriggerExecutable) ReinstallCloudServerWait(resource *ResourceOperate, provider cloudrepo.CloudProvider, imageCID, cloudID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ReinstallServerWaitTime)
	defer cancel()
forLoop:
	for {
		select {
		case <-ctx.Done():
			log.Error("RebootServer instance %s reinstall timeout", resource.InstanceCID)
			return errors.Errorf("timeout")
		default:
			reply, err := provider.DescribeServer(&cloudrepo.DescribeCloudServerReq{
				InstanceID: resource.InstanceCID,
				ProjectID:  resource.ProjectCID,
			})
			if err != nil {
				log.Error("ReinstallCloudServerWait instance %s describe server failed err: %v", resource.InstanceCID, err)
				return errors.Wrap(err, "cloud_server reinstall server describe server failed")
			}
			if cloudID == cloudrepo.CloudAli && !reply.IsOperation { //阿里云判断重装系统是否已完成
				time.Sleep(30 * time.Second)
				continue forLoop
			}
			if cloudID == cloudrepo.CloudTencent && reply.LatestOperationState != "SUCCESS" { //腾讯云判断重装系统是否已完成
				time.Sleep(30 * time.Second)
				continue forLoop
			}
			if reply.ImageID != imageCID {
				time.Sleep(30 * time.Second)
				continue forLoop
			}
			break
		}
		return nil
	}
}

// SetReservationTime 设置保留时间
func (oe *ObjectTriggerExecutable) SetReservationTime(period int) error {
	err := jobUseCase.UpdateObjectReservationTime(&biz.Object{
		Id:                oe.Id,
		ReservationPeriod: period,
	})
	if err != nil {
		log.Error("SetReservationTime update object reservation time failed err: %v", err)
		return errors.Wrap(err, "set retention time failed")
	}
	return nil
}

// SetReservationTimeWait 等待设置保留时间完成
func (oe *ObjectTriggerExecutable) SetReservationTimeWait() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
forLoop:
	for {
		select {
		case <-ctx.Done():
			log.Error("SetReservationTimeWait object id %s set retention time timeout", oe.Id)
			return errors.Errorf("timeout")
		default:
			object, err := jobUseCase.GetObjectByObjectId(oe.Id)
			if err != nil {
				log.Error("SetReservationTimeWait get object by id %s failed err: %v", oe.Id, err)
				return errors.Wrap(err, "set retention time get object failed")
			}
			optTime := object.StartTime.Add(time.Second * time.Duration(object.ReservationPeriod)) // 保留时间到期时间
			if optTime.After(time.Now()) {                                                         // 保留时间未到期，等待
				time.Sleep(60 * time.Second)
				continue forLoop
			}
			break
		}
		return nil
	}
}

// ChangeCloudServerChargeType 变更云服务器计费类型
func (oe *ObjectTriggerExecutable) ChangeCloudServerChargeType(resource *ResourceOperate, provider cloudrepo.CloudProvider) error {
	_, err := provider.ChangeServerChargeType(&cloudrepo.ChangeServerChargeTypeReq{
		InstanceID: resource.InstanceCID,
		ProjectID:  resource.ProjectCID,
		ChargeType: "PostPaid",
	})
	if err != nil {
		log.Error("ChangeCloudServerChargeType instance %s stop failed err: %v", resource.InstanceCID, err)
		return errors.Wrap(err, "cloud_server change server charge type failed")
	}
	return nil
}

// ChangeCloudServerChargeTypeWait 等待变更云服务器计费类型完成
func (oe *ObjectTriggerExecutable) ChangeCloudServerChargeTypeWait(resource *ResourceOperate, provider cloudrepo.CloudProvider) error {
	ctx, cancel := context.WithTimeout(context.Background(), ChangeServerChargeTypeTime)
	defer cancel()
forLoop:
	for {
		select {
		case <-ctx.Done():
			log.Error("ChangeCloudServerChargeTypeWait instance %s change charge type timeout", resource.InstanceCID)
			return errors.Errorf("timeout")
		default:
			reply, err := provider.DescribeServer(&cloudrepo.DescribeCloudServerReq{
				InstanceID: resource.InstanceCID,
				ProjectID:  resource.ProjectCID,
			})
			if err != nil {
				log.Error("ChangeCloudServerChargeTypeWait instance %s describe server failed err: %v", resource.InstanceCID, err)
				return errors.Wrap(err, "cloud_server change charge type server describe server failed")
			}
			if reply.ChargeType != cloudrepo.CloudVmChargeTypePostPaid {
				time.Sleep(30 * time.Second)
				continue forLoop
			}
			break
		}
		return nil
	}
}

// DeleteCloudServer 删除云服务器
func (oe *ObjectTriggerExecutable) DeleteCloudServer(resource *ResourceOperate, provider cloudrepo.CloudProvider) error {
	// 查询云服务器是否存在
	_, err := provider.DescribeServer(&cloudrepo.DescribeCloudServerReq{
		InstanceID: resource.InstanceCID,
		ProjectID:  resource.ProjectCID,
	})
	if err != nil && strings.Contains(err.Error(), NotFoundErrorStr) {
		log.Warn("DeleteCloudServer instance %s not found", resource.InstanceCID)
		return nil
	}
	// 删除云服务器
	_, err = provider.DeleteServer(&cloudrepo.DeleteCloudServerReq{
		InstanceID: resource.InstanceCID,
		ProjectID:  resource.ProjectCID,
	})
	if err != nil {
		log.Error("DeleteCloudServer instance %s delete failed err: %v", resource.InstanceCID, err)
		return errors.Wrap(err, "cloud_server delete server failed")
	}
	return nil
}

// DeleteCloudServerWait 等待删除云服务器完成
func (oe *ObjectTriggerExecutable) DeleteCloudServerWait(resource *ResourceOperate, provider cloudrepo.CloudProvider) error {
	ctx, cancel := context.WithTimeout(context.Background(), DeleteServerWaitTime)
	defer cancel()
forLoop:
	for {
		select {
		case <-ctx.Done():
			log.Error("DeleteCloudServerWait instance %s delete timeout", resource.InstanceCID)
			return errors.Errorf("timeout")
		default:
			resp, err := provider.DescribeServer(&cloudrepo.DescribeCloudServerReq{
				InstanceID: resource.InstanceCID,
				ProjectID:  resource.ProjectCID,
			})
			if err != nil && strings.Contains(err.Error(), NotFoundErrorStr) {
				log.Warn("DeleteCloudServerWait instance %s describe server has been deleted", resource.InstanceCID)
				return nil
			}
			if resp.InstanceState == cloudrepo.CloudVMStatusTerminated {
				log.Warn("DeleteCloudServerWait instance %s has been deleted", resource.InstanceCID)
				return nil
			}
			time.Sleep(30 * time.Second)
			continue forLoop
		}
	}
}

// CreateTerraformJobEnvironment 创建terraform工作环境
func (oe *ObjectTriggerExecutable) CreateTerraformJobEnvironment(req *Terraform) (statePath string, outputPath string, err error) {
	JobID := req.JobID
	if JobID == 0 {
		return "", "", errors.Errorf("job id is empty")
	}
	stepID := req.StepID
	if stepID == 0 {
		return "", "", errors.Errorf("step id is empty")
	}
	objectID := req.ObjectID
	if objectID == 0 {
		return "", "", errors.Errorf("object id is empty")
	}
	operate := req.Operate
	if operate == "" {
		return "", "", errors.Errorf("operate is empty")
	}
	statePath = fmt.Sprintf("%v/job-%v/step-%v/object-%v/%v/state", configs.Conf.CloudConf.TerraformPath, JobID, stepID, objectID, operate)
	outputPath = fmt.Sprintf("%v/job-%v/step-%v/object-%v/%v/output", configs.Conf.CloudConf.TerraformPath, JobID, stepID, objectID, operate)
	if err := os.MkdirAll(statePath, os.ModePerm); err != nil {
		log.Error("CreateTerraformJobEnvironment mkdir state path failed err: %v", err)
		return "", "", errors.Wrap(err, "mkdir state path failed")
	}
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		log.Error("CreateTerraformJobEnvironment mkdir output path failed err: %v", err)
		return "", "", errors.Wrap(err, "mkdir output path failed")
	}
	return statePath, outputPath, nil
}

// CreateTerraformServerParameterDependent 创建云服务器参数依赖
func (oe *ObjectTriggerExecutable) CreateTerraformServerParameterDependent(req *biz.CreateCloudServer) (*biz.CreateCloudServer, error) {
	var data *biz.CreateCloudServer
	data = req
	if req.Region.ID != 0 { //获取云地域
		resp, err := regionUseCase.QueryCloudRegion(&biz.CloudRegionWhere{Conditions: map[string]interface{}{"id": req.Region.ID}}, &biz.CloudRegionOutput{})
		if err != nil {
			log.Error("CreateTerraformServerParameterDependent query cloud region failed err: %v", err)
			return nil, errors.Wrap(err, "query cloud region failed")
		}
		if len(resp) < 1 {
			return nil, errors.Errorf("cloud region not found")
		}
		data.Region = resp[0]
	}
	if req.Zone.ID != 0 { //获取云可用区
		resp, err := zoneUseCase.QueryCloudZone(&biz.CloudZoneWhere{Conditions: map[string]interface{}{"id": req.Zone.ID}}, &biz.CloudZoneOutput{})
		if err != nil {
			log.Error("CreateTerraformServerParameterDependent query cloud zone failed err: %v", err)
			return nil, errors.Wrap(err, "query cloud zone failed")
		}
		if len(resp) < 1 {
			return nil, errors.Errorf("cloud zone not found")
		}
		data.Zone = resp[0]
	}
	if req.Vpc.ID != 0 { //获取云VPC
		resp, err := vpcUseCase.QueryCloudVpc(&biz.CloudVpcWhere{Conditions: map[string]interface{}{"id": req.Vpc.ID}}, &biz.CloudVpcOutput{})
		if err != nil {
			log.Error("CreateTerraformServerParameterDependent query cloud vpc failed err: %v", err)
			return nil, errors.Wrap(err, "query cloud vpc failed")
		}
		if len(resp) < 1 {
			return nil, errors.Errorf("cloud vpc not found")
		}
		data.Vpc = resp[0]
	}
	if req.Subnet.ID != 0 { //获取云子网
		resp, err := subnetUseCase.QueryCloudSubnet(&biz.CloudSubnetWhere{Conditions: map[string]interface{}{"id": req.Subnet.ID}}, &biz.CloudSubnetOutput{})
		if err != nil {
			log.Error("CreateTerraformServerParameterDependent query cloud subnet failed err: %v", err)
			return nil, errors.Wrap(err, "query cloud subnet failed")
		}
		if len(resp) < 1 {
			return nil, errors.Errorf("cloud subnet not found")
		}
		data.Subnet = resp[0]
	}
	if req.SecurityGroup != nil { //获取云安全组
		sgID, ok := req.SecurityGroup.(map[string]interface{})["id"]
		if !ok {
			return nil, errors.Errorf("security group id not found")
		}
		resp, err := securityGroupUseCase.QueryCloudSecurityGroup(&biz.CloudSecurityGroupWhere{Conditions: map[string]interface{}{"id": sgID}}, &biz.CloudSecurityGroupOutput{})
		if err != nil {
			log.Error("CreateTerraformServerParameterDependent query cloud security group failed err: %v", err)
			return nil, errors.Wrap(err, "query cloud security group failed")
		}
		if len(resp) < 1 {
			return nil, errors.Errorf("cloud security group not found")
		}
		data.SecurityGroup = resp[0]
	}
	if req.ServerImage.ID != 0 { //获取云服务器镜像
		resp, err := severImageUseCase.QueryCloudServerImage(&biz.CloudServerImageWhere{Conditions: map[string]interface{}{"id": req.ServerImage.ID}}, &biz.CloudServerImageOutput{})
		if err != nil {
			log.Error("CreateTerraformServerParameterDependent query cloud server image failed err: %v", err)
			return nil, errors.Wrap(err, "query cloud server image failed")
		}
		if len(resp) < 1 {
			return nil, errors.Errorf("cloud server image not found")
		}
		data.ServerImage = resp[0]
	}
	if req.ServerSpec.ID != 0 { //获取云服务器规格
		resp, err := serverSpecUseCase.QueryCloudServerSpec(&biz.CloudServerSpecWhere{Conditions: map[string]interface{}{"id": req.ServerSpec.ID}}, &biz.CloudServerSpecOutput{})
		if err != nil {
			log.Error("CreateTerraformServerParameterDependent query cloud server spec failed err: %v", err)
			return nil, errors.Wrap(err, "query cloud server spec failed")
		}
		if len(resp) < 1 {
			return nil, errors.Errorf("cloud server spec not found")
		}
		data.ServerSpec = resp[0]
	}
	if req.DiskType.ID != 0 { //获取云服务器磁盘类型
		resp, err := diskTypeUseCase.QueryDiskType(&biz.DiskTypeWhere{Conditions: map[string]interface{}{"id": req.DiskType.ID}}, &biz.DiskTypeOutput{})
		if err != nil {
			log.Error("CreateTerraformServerParameterDependent query disk type failed err: %v", err)
			return nil, errors.Wrap(err, "query disk type failed")
		}
		if len(resp) < 1 {
			return nil, errors.Errorf("disk type not found")
		}
		data.DiskType = resp[0]
	}
	if req.ChargeType.ID != 0 { //获取云服务器计费类型
		resp, err := chargeTypeUseCase.QueryChargeType(&biz.ChargeTypeWhere{Conditions: map[string]interface{}{"id": req.ChargeType.ID}}, &biz.ChargeTypeOutput{})
		if err != nil {
			log.Error("CreateTerraformServerParameterDependent query charge type failed err: %v", err)
			return nil, errors.Wrap(err, "query charge type failed")
		}
		if len(resp) < 1 {
			return nil, errors.Errorf("charge type not found")
		}
		data.ChargeType = resp[0]
	}
	if req.RenewalType.ID != 0 { //获取云服务器续费类型
		resp, err := chargeTypeUseCase.QueryChargeType(&biz.ChargeTypeWhere{Conditions: map[string]interface{}{"id": req.RenewalType.ID}}, &biz.ChargeTypeOutput{})
		if err != nil {
			log.Error("CreateTerraformServerParameterDependent query renewal type failed err: %v", err)
			return nil, errors.Wrap(err, "query renewal type failed")
		}
		if len(resp) < 1 {
			return nil, errors.Errorf("renewal type not found")
		}
		data.RenewalType = resp[0]
	}
	return data, nil
}

// CreateTerraformTemplates 创建terraform模板
func (oe *ObjectTriggerExecutable) CreateTerraformTemplates(req *Terraform, providerID int) error {
	input := req.Input // 输入参数
	if input == "" {
		return errors.Errorf("input is empty")
	}
	KeyPairName := req.KeyPairName // 密钥对名称
	if KeyPairName == "" {
		return errors.Errorf("key pair id is empty")
	}
	secretID := req.SecretID // 账号AK
	if input == "" {
		return errors.Errorf("secret id is empty")
	}
	secretKey := req.SecretKey // 账号SK
	if input == "" {
		return errors.Errorf("secret key is empty")
	}
	operate := req.Operate
	if operate == "" {
		return errors.Errorf("operate name is empty")
	}
	accountID := req.AccountID
	if accountID == 0 {
		return errors.Errorf("account id is empty")
	}
	if providerID == 0 {
		return errors.Errorf("provider id is empty")
	}
	prefix := req.NamePrefix
	if prefix == "" {
		return errors.Errorf("name prefix is empty")
	}
	suffix := req.NameSuffix
	if suffix == "" {
		return errors.Errorf("name suffix is empty")
	}
	counter := req.OperandCounter
	if counter == 0 {
		return errors.Errorf("operand counter is empty")
	}
	statePath := req.StatePath
	if statePath == "" {
		return errors.Errorf("state path is empty")
	}
	rdsArch := req.RdsArch
	redisArch := req.RedisArch
	DependentKey := req.DependentKey
	DependentValue := req.DependentValue
	conditions := make(map[string]interface{})
	conditions["operate"] = operate
	conditions["provider_id"] = providerID
	if req.bindPublicIP != "" {
		conditions["dependent_parameters_REGEX"] = fmt.Sprintf("bindPublicIp:%v", req.bindPublicIP)
	}
	if rdsArch != "" {
		conditions["dependent_parameters_REGEX"] = fmt.Sprintf("dbArch:%v", rdsArch)
	}
	if redisArch != "" {
		conditions["dependent_parameters_REGEX"] = fmt.Sprintf("redisArch:%v", redisArch)
	}
	if DependentKey != "" && DependentValue != "" {
		conditions["dependent_parameters"] = fmt.Sprintf("%v:%v", DependentKey, DependentValue)
	}
	// 获取项目ID
	jobData, err := jobUseCase.QueryJob(&biz.JobWhere{Conditions: map[string]interface{}{"id": req.JobID}}, &biz.JobOutput{})
	if err != nil {
		oe.Status = enum.ExecutableStateFailure
		log.Error("CreateTerraformTemplates query job failed err: %v", err)
		return errors.Wrap(err, "query job failed")
	}
	if len(jobData) == 0 || jobData[0].Raw == "" {
		oe.Status = enum.ExecutableStateFailure
		return errors.Errorf("job data is empty")
	}
	// 解析JSON数据到结构体切片
	var projectAccountConfig *biz.ProjectAccountConfig
	err = json.Unmarshal([]byte(jobData[0].Raw), &projectAccountConfig)
	if err != nil {
		oe.Status = enum.ExecutableStateFailure
		return errors.Wrap(err, "cloud_server create server get project id json unmarshal failed")
	}
	// 获取terraform模版
	resp, err := terraformUseCase.QueryTerraform(&biz.TerraformWhere{Conditions: conditions}, &biz.TerraformOutput{})
	if err != nil {
		log.Error("CreateTerraformTemplates query terraform failed err: %v", err)
		return errors.Wrap(err, "query terraform failed")
	}
	if len(resp) == 0 {
		return errors.Errorf("terraform template is empty")
	}
	// 获取云项目信息
	projectConditions := map[string]interface{}{
		"account_id":        accountID,
		"project_config_id": projectAccountConfig.ProjectConfig.ID,
		"is_delete":         0,
	}
	projectInfo, err := projectAccountConfigUseCase.QueryProjectAccountConfig(&biz.ProjectAccountConfigWhere{Conditions: projectConditions}, &biz.ProjectAccountConfigOutput{})
	if err != nil {
		log.Error("CreateTerraformTemplates query project failed err: %v", err)
		return errors.Wrap(err, "CreateTerraformTemplates query project failed")
	}
	if len(projectInfo) == 0 {
		return errors.Errorf("project info is empty")
	}
	// 生成资源模版
	for _, t := range resp {
		path := fmt.Sprintf("%v/%v.tf", statePath, t.Name)
		file, err := os.Create(path)
		if err != nil {
			log.Error("CreateTerraformTemplates create resource file failed err: %v", err)
			return errors.Wrap(err, "create resource file failed")
		}
		if err := ioutil.WriteFile(path, []byte(t.Data), os.ModePerm); err != nil {
			log.Error("CreateTerraformTemplates write file failed err: %v", err)
			return errors.Wrap(err, "write resource file failed")
		}
		_ = file.Close()
	}
	// 生成变量模版
	inputM := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &inputM); err != nil {
		log.Error("CreateTerraformTemplates unmarshal input failed err: %v", err)
		return errors.Wrap(err, "unmarshal input failed")
	}
	inputM["key_pair_name"] = KeyPairName
	inputM["secret_ak"] = secretID
	inputM["secret_sk"] = secretKey
	inputM["cloud_project_name"] = projectInfo[0].Project.Name
	inputM["cloud_project_cid"] = projectInfo[0].Project.CID
	index := oe.OperandCounter
	GenInputName(prefix, suffix, index, inputM)
	variables := genVariable("input", inputM)
	path := fmt.Sprintf("%v/variable.tf", statePath)
	file, err := os.Create(path)
	if err != nil {
		log.Error("CreateTerraformTemplates create variable file failed err: %v", err)
		return errors.Wrap(err, "create variable file failed")
	}
	if err := ioutil.WriteFile(path, []byte(variables), os.ModePerm); err != nil {
		log.Error("CreateTerraformTemplates write file failed err: %v", err)
		return errors.Wrap(err, "write variable file failed")
	}
	_ = file.Close()
	return nil
}

// ExecTerraformInit 执行terraform init
func (oe *ObjectTriggerExecutable) ExecTerraformInit(req *Terraform) error {
	statePath := req.StatePath
	if statePath == "" {
		return errors.Errorf("state path is empty")
	}
	outputPath := req.OutputPath
	if outputPath == "" {
		return errors.Errorf("output path is empty")
	}
	cmd := fmt.Sprintf("init -plugin-dir=%v", configs.Conf.CloudConf.TerraformPlugin)
	out, err := runTerraformCmd(cmd, statePath, outputPath)
	if err != nil {
		log.Error("ExecTerraformInit run terraform init failed err: %v", err)
		return errors.Wrap(err, "run terraform init failed")
	}
	log.Info("ExecTerraformInit run terraform init success stdout: %v", out)
	return nil
}

// ExecTerraformPlan 执行terraform plan
func (oe *ObjectTriggerExecutable) ExecTerraformPlan(req *Terraform) error {
	statePath := req.StatePath
	if statePath == "" {
		return errors.Errorf("state path is empty")
	}
	outputPath := req.OutputPath
	if outputPath == "" {
		return errors.Errorf("output path is empty")
	}
	out, err := runTerraformCmd("plan", statePath, outputPath)
	if err != nil {
		log.Error("ExecTerraformPlan run terraform plan failed err: %v", err)
		return errors.Wrap(err, "run terraform plan failed")
	}
	log.Info("ExecTerraformPlan run terraform plan success stdout: %v", out)
	return nil
}

// ExecTerraformApply 执行terraform apply
func (oe *ObjectTriggerExecutable) ExecTerraformApply(req *Terraform) error {
	statePath := req.StatePath
	if statePath == "" {
		return errors.Errorf("state path is empty")
	}
	outputPath := req.OutputPath
	if outputPath == "" {
		return errors.Errorf("output path is empty")
	}
	out, err := runTerraformCmd("apply -auto-approve", statePath, outputPath)
	if err != nil {
		log.Error("ExecTerraformApply run terraform apply failed err: %v", err)
		return errors.Wrap(err, "run terraform apply failed")
	}
	log.Info("ExecTerraformApply run terraform apply success stdout: %v", out)
	out, err = runTerraformCmd("output resource_id", statePath, outputPath)
	if err != nil {
		log.Error("ExecTerraformApply run terraform output failed err: %v", err)
		return errors.Wrap(err, "run terraform output failed")
	}
	out, err = runTerraformCmd("output -json", statePath, outputPath)
	if err != nil {
		log.Error("ExecTerraformApply run terraform output json failed err: %v", err)
		return errors.Wrap(err, "run terraform output json failed")
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		log.Error("ExecTerraformApply unmarshal output json failed err: %v", err)
		return errors.Wrap(err, "unmarshal output json failed")
	}
	bytes, _ := json.Marshal(m)
	log.Info("ExecTerraformApply terraform output json success: %v", string(bytes))
	return nil
}

func runTerraformCmd(cmd, statePath, outPath string) (string, error) {
	command := exec.CommandContext(context.Background(), configs.Conf.CloudConf.TerraformBin, strings.Fields(cmd)...)
	command.Dir = statePath
	command.Env = []string{}
	raw, err := command.CombinedOutput()
	if err != nil {
		return string(raw), errors.Errorf("run terraform cmd args[%v], stdout: %v, err: %v", command.Args, string(raw), err)
	}
	return string(raw), nil
}

func GenInputName(prefix, suffix string, index int, input map[string]interface{}) {
	if prefix != "" && suffix != "" {
		if suffix == "standard-01" && index < 10 {
			input["name"] = fmt.Sprintf("%v-0%v", prefix, index)
		} else if suffix == "standard-01" && index >= 10 {
			input["name"] = fmt.Sprintf("%v-%v", prefix, index)
		} else if suffix == "standard-1" {
			input["name"] = fmt.Sprintf("%v-%v", prefix, index)
		} else {
			input["name"] = fmt.Sprintf("%v", prefix)
		}
	}
}

func genVariable(key string, m interface{}) string {
	kind := reflect.TypeOf(m).Kind()
	switch kind {
	case reflect.String, reflect.Float32, reflect.Float64:
		return fmt.Sprintf(`
variable %v {
	type = string 
	default = "%v"
}
`, key, m)
	case reflect.Slice:
		list := []string{}
		switch m.(type) {
		case []interface{}:
			for _, vv := range m.([]interface{}) {
				list = append(list, fmt.Sprintf(`"%v"`, vv))
			}
		case []string:
			for _, vv := range m.([]string) {
				list = append(list, fmt.Sprintf(`"%v"`, vv))
			}
		}
		return fmt.Sprintf(`
variable %v {
	type = list
	default = [%v]
}`, key, strings.Join(list, ","))
	case reflect.Map:
		result := ""
		for k, vv := range m.(map[string]interface{}) {
			if vv == nil {
				continue
			}
			if strings.Contains(k, "_id") || strings.Contains(k, "security_group_cid") || strings.Contains(k, "subnet_cid") {
				// 避免重复声明
				continue
			}
			result += genVariable(fmt.Sprintf("%v_%v", key, k), vv)
		}
		return result
	default:
		return ""
	}
}
