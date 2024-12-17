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
	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	dataSync "github.com/bilibili/HCP/app/interface/v1/internal/biz/sync"
)

var (
	chargeTypeUseCase           *biz.ChargeTypeUseCase
	diskTypeUseCase             *biz.DiskTypeUseCase
	serverSpecUseCase           *biz.CloudServerSpecUseCase
	severImageUseCase           *biz.CloudServerImageUseCase
	securityGroupUseCase        *biz.CloudSecurityGroupUseCase
	subnetUseCase               *biz.CloudSubnetUseCase
	vpcUseCase                  *biz.CloudVpcUseCase
	zoneUseCase                 *biz.CloudZoneUseCase
	regionUseCase               *biz.CloudRegionUseCase
	jobUseCase                  *biz.JobUseCase
	providerUseCase             *biz.ProviderUseCase
	accountUseCase              *biz.AccountUseCase
	userUseCase                 *biz.UserUseCase
	projectConfigUseCase        *biz.ProjectConfigUseCase
	projectAccountConfigUseCase *biz.ProjectAccountConfigUseCase
	projectUserConfigUseCase    *biz.ProjectUserConfigUseCase
	cloudServerUseCase          *biz.CloudServerUseCase
	cloudServerSpecUseCase      *biz.CloudServerSpecUseCase
	cloudServerImageUseCase     *biz.CloudServerImageUseCase
	terraformUseCase            *biz.TerraformUseCase
	datasync                    *dataSync.Client
	worker                      *Worker
)

// Worker 任务工作
type Worker struct {
	// Job 任务
	JobPipeline Pipeline
	// JobCreator 任务创建者
	JobCreatorPipeline Pipeline
	// cancel 关闭通道
	cancel chan struct{}
}

// NewWorker 创建任务工作客户端
func NewWorker(job *biz.JobUseCase, provider *biz.ProviderUseCase, account *biz.AccountUseCase, user *biz.UserUseCase, projectConfig *biz.ProjectConfigUseCase,
	projectAccountConfig *biz.ProjectAccountConfigUseCase, projectUserConfig *biz.ProjectUserConfigUseCase, cloudServer *biz.CloudServerUseCase, cloudServerSpec *biz.CloudServerSpecUseCase,
	cloudServerImage *biz.CloudServerImageUseCase, terraform *biz.TerraformUseCase, region *biz.CloudRegionUseCase, zone *biz.CloudZoneUseCase, vpc *biz.CloudVpcUseCase,
	subnet *biz.CloudSubnetUseCase, securityGroup *biz.CloudSecurityGroupUseCase, severImage *biz.CloudServerImageUseCase, serverSpec *biz.CloudServerSpecUseCase,
	diskType *biz.DiskTypeUseCase, chargeType *biz.ChargeTypeUseCase, sync *dataSync.Client) (*Worker, func()) {
	jobUseCase = job
	providerUseCase = provider
	accountUseCase = account
	userUseCase = user
	projectConfigUseCase = projectConfig
	projectAccountConfigUseCase = projectAccountConfig
	projectUserConfigUseCase = projectUserConfig
	cloudServerUseCase = cloudServer
	cloudServerSpecUseCase = cloudServerSpec
	cloudServerImageUseCase = cloudServerImage
	terraformUseCase = terraform
	regionUseCase = region
	zoneUseCase = zone
	vpcUseCase = vpc
	subnetUseCase = subnet
	securityGroupUseCase = securityGroup
	severImageUseCase = severImage
	serverSpecUseCase = serverSpec
	diskTypeUseCase = diskType
	chargeTypeUseCase = chargeType
	datasync = sync
	worker = &Worker{
		//TODO: 由 Job.raw 生成批次和对象
		JobCreatorPipeline: NewPipeline("JobCreatorPipeline", 1000, NewJobCreatorPipelineInput()),
		JobPipeline:        NewPipeline("JobPipeline", 1000, NewJobPipelineInput()),
		cancel:             make(chan struct{}, 1),
	}
	return worker, func() {
		worker.Close()
	}
}

// Start 启动任务工作
func (o *Worker) Start() {
	go o.JobPipeline.Run()        // 启动任务管道
	go o.JobCreatorPipeline.Run() // 启动任务创建管道
}

// Close 关闭任务工作
func (o *Worker) Close() {
	close(o.cancel)
	o.JobPipeline.Shutdown()
	o.JobCreatorPipeline.Shutdown()
}
